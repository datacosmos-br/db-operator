/*
 * Copyright 2024 Datacosmos
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package database

import (
	"context"
	"testing"

	"github.com/db-operator/db-operator/v2/pkg/test"
	"github.com/stretchr/testify/assert"
)

// Compile-time check that ClickHouse implements the Database interface.
var _ Database = (*ClickHouse)(nil)

func testClickhouse() (ClickHouse, *DatabaseUser) {
	ch := ClickHouse{
		Host:     test.GetClickhouseHost(),
		Port:     test.GetClickhousePort(),
		Database: "testdb",
	}
	user := &DatabaseUser{
		Username:   "testuser",
		Password:   "testpwd",
		AccessType: ACCESS_TYPE_MAINUSER,
	}
	return ch, user
}

func getClickhouseAdmin() *DatabaseUser {
	return &DatabaseUser{
		Username: test.GetClickhouseAdminUsername(),
		Password: test.GetClickhouseAdminPassword(),
	}
}

// cleanupClickhouse drops the database and user so a test starts from a known state.
func cleanupClickhouse(ch ClickHouse, admin, user *DatabaseUser) {
	_ = ch.deleteUser(context.TODO(), admin, user)
	_ = ch.deleteDatabase(context.TODO(), admin)
}

func TestClickhouseParseAdminCredentials(t *testing.T) {
	ch, _ := testClickhouse()

	invalidData := map[string][]byte{"unknownkey": []byte("wrong")}
	_, err := ch.ParseAdminCredentials(context.TODO(), invalidData)
	assert.Error(t, err, "expected error for missing admin keys")

	missingPassword := map[string][]byte{"user": []byte("admin")}
	_, err = ch.ParseAdminCredentials(context.TODO(), missingPassword)
	assert.Error(t, err, "expected error for missing password")

	validData := map[string][]byte{
		"user":     []byte("admin"),
		"password": []byte("s3cret"),
	}
	cred, err := ch.ParseAdminCredentials(context.TODO(), validData)
	assert.NoError(t, err)
	assert.Equal(t, "admin", cred.Username)
	assert.Equal(t, "s3cret", cred.Password)
}

func TestClickhouseGetCredentials(t *testing.T) {
	ch, user := testClickhouse()

	cred := ch.GetCredentials(context.TODO(), user)
	assert.Equal(t, ch.Database, cred.Name)
	assert.Equal(t, user.Username, cred.Username)
	assert.Equal(t, user.Password, cred.Password)
}

func TestClickhouseGetDatabaseAddress(t *testing.T) {
	ch, _ := testClickhouse()

	addr := ch.GetDatabaseAddress(context.TODO())
	assert.Equal(t, ch.Host, addr.Host)
	assert.Equal(t, ch.Port, addr.Port)
}

func TestClickhouseCreateDatabase(t *testing.T) {
	ch, user := testClickhouse()
	ch.Database = "ch_create_db"
	admin := getClickhouseAdmin()
	cleanupClickhouse(ch, admin, user)
	defer cleanupClickhouse(ch, admin, user)

	assert.False(t, ch.isDbExist(context.TODO(), admin))

	assert.NoError(t, ch.createDatabase(context.TODO(), admin))
	assert.True(t, ch.isDbExist(context.TODO(), admin))

	// createDatabase is idempotent (CREATE DATABASE IF NOT EXISTS).
	assert.NoError(t, ch.createDatabase(context.TODO(), admin))
	assert.True(t, ch.isDbExist(context.TODO(), admin))
}

func TestClickhouseCreateOrUpdateUser(t *testing.T) {
	ch, user := testClickhouse()
	ch.Database = "ch_user_db"
	admin := getClickhouseAdmin()
	cleanupClickhouse(ch, admin, user)
	defer cleanupClickhouse(ch, admin, user)

	assert.NoError(t, ch.createDatabase(context.TODO(), admin))
	assert.False(t, ch.isUserExist(context.TODO(), admin, user))

	assert.NoError(t, ch.createOrUpdateUser(context.TODO(), admin, user))
	assert.True(t, ch.isUserExist(context.TODO(), admin, user))

	// Second call updates the existing user instead of failing.
	assert.NoError(t, ch.createOrUpdateUser(context.TODO(), admin, user))
	assert.True(t, ch.isUserExist(context.TODO(), admin, user))
}

func TestClickhouseCheckStatus(t *testing.T) {
	ch, user := testClickhouse()
	ch.Database = "ch_status_db"
	admin := getClickhouseAdmin()
	cleanupClickhouse(ch, admin, user)
	defer cleanupClickhouse(ch, admin, user)

	// Neither database nor user exist yet.
	assert.Error(t, ch.CheckStatus(context.TODO(), user))

	assert.NoError(t, ch.createDatabase(context.TODO(), admin))
	assert.NoError(t, ch.createOrUpdateUser(context.TODO(), admin, user))
	assert.NoError(t, ch.CheckStatus(context.TODO(), user))

	// Wrong password must fail the status check.
	bad := &DatabaseUser{Username: user.Username, Password: "wrongpwd"}
	assert.Error(t, ch.CheckStatus(context.TODO(), bad))
}

func TestClickhouseQueryAsUser(t *testing.T) {
	ch, user := testClickhouse()
	ch.Database = "ch_query_db"
	admin := getClickhouseAdmin()
	cleanupClickhouse(ch, admin, user)
	defer cleanupClickhouse(ch, admin, user)

	assert.NoError(t, ch.createDatabase(context.TODO(), admin))
	assert.NoError(t, ch.createOrUpdateUser(context.TODO(), admin, user))

	createTable := "CREATE TABLE ch_query_db.items (id Int32, name String) ENGINE = MergeTree ORDER BY id"
	assert.NoError(t, ch.execAsUser(context.TODO(), createTable, user))

	// QueryAsUser is the only data-plane path the operator uses (templated
	// SELECTs in internal/utils/templates); it reads a single scalar value.
	res, err := ch.QueryAsUser(context.TODO(),
		"SELECT name FROM system.tables WHERE database = 'ch_query_db' AND name = 'items'", user)
	assert.NoError(t, err)
	assert.Equal(t, "items", res)

	// A malformed query must surface an error.
	_, err = ch.QueryAsUser(context.TODO(), "SELECT name FROM ch_query_db.does_not_exist", user)
	assert.Error(t, err)
}

func TestClickhouseUserLifecycle(t *testing.T) {
	ch, user := testClickhouse()
	ch.Database = "ch_lifecycle_db"
	admin := getClickhouseAdmin()
	cleanupClickhouse(ch, admin, user)
	defer cleanupClickhouse(ch, admin, user)

	assert.NoError(t, ch.createDatabase(context.TODO(), admin))

	// Create user, then a repeated create is a no-op (CREATE USER IF NOT EXISTS).
	assert.NoError(t, ch.createUser(context.TODO(), admin, user))
	assert.NoError(t, ch.createUser(context.TODO(), admin, user))
	assert.True(t, ch.isUserExist(context.TODO(), admin, user))

	// Rotating the password keeps the user usable.
	user.Password = "rotatedpwd"
	assert.NoError(t, ch.updateUser(context.TODO(), admin, user))

	// Without a grant the user cannot operate on the database.
	createTable := "CREATE TABLE ch_lifecycle_db.t (id Int32) ENGINE = MergeTree ORDER BY id"
	assert.Error(t, ch.execAsUser(context.TODO(), createTable, user))

	// After the grant the user gains full access to the database.
	assert.NoError(t, ch.setUserPermission(context.TODO(), admin, user))
	assert.NoError(t, ch.execAsUser(context.TODO(), createTable, user))
	assert.NoError(t, ch.execAsUser(context.TODO(), "DROP TABLE ch_lifecycle_db.t", user))

	// Revoking removes the access again.
	assert.NoError(t, ch.revokePermissions(context.TODO(), admin, user))
	assert.Error(t, ch.execAsUser(context.TODO(), createTable, user))

	// Revoking from a non-existent user is a no-op.
	ghost := &DatabaseUser{Username: "ch_ghost_user", Password: "x"}
	assert.NoError(t, ch.revokePermissions(context.TODO(), admin, ghost))

	assert.NoError(t, ch.deleteUser(context.TODO(), admin, user))
	assert.False(t, ch.isUserExist(context.TODO(), admin, user))
}

func TestClickhouseDeleteUser(t *testing.T) {
	ch, user := testClickhouse()
	ch.Database = "ch_deluser_db"
	admin := getClickhouseAdmin()
	cleanupClickhouse(ch, admin, user)
	defer cleanupClickhouse(ch, admin, user)

	assert.NoError(t, ch.createDatabase(context.TODO(), admin))
	assert.NoError(t, ch.createOrUpdateUser(context.TODO(), admin, user))
	assert.True(t, ch.isUserExist(context.TODO(), admin, user))

	assert.NoError(t, ch.deleteUser(context.TODO(), admin, user))
	assert.False(t, ch.isUserExist(context.TODO(), admin, user))

	// deleteUser is idempotent (DROP USER IF EXISTS).
	assert.NoError(t, ch.deleteUser(context.TODO(), admin, user))
}

func TestClickhouseDeleteDatabase(t *testing.T) {
	ch, user := testClickhouse()
	ch.Database = "ch_deldb_db"
	admin := getClickhouseAdmin()
	cleanupClickhouse(ch, admin, user)
	defer cleanupClickhouse(ch, admin, user)

	assert.NoError(t, ch.createDatabase(context.TODO(), admin))
	assert.True(t, ch.isDbExist(context.TODO(), admin))

	assert.NoError(t, ch.deleteDatabase(context.TODO(), admin))
	assert.False(t, ch.isDbExist(context.TODO(), admin))

	// deleteDatabase is idempotent (DROP DATABASE IF EXISTS).
	assert.NoError(t, ch.deleteDatabase(context.TODO(), admin))
}

// TestUnitClickhouseGrantPrivileges is a pure unit test (no ClickHouse needed).
func TestUnitClickhouseGrantPrivileges(t *testing.T) {
	cases := map[string]string{
		ACCESS_TYPE_READONLY:  "SELECT",
		ACCESS_TYPE_READWRITE: "SELECT, INSERT, ALTER, CREATE, DROP",
		ACCESS_TYPE_MAINUSER:  "ALL",
	}
	for accessType, want := range cases {
		got, err := clickhouseGrantPrivileges(accessType)
		assert.NoError(t, err)
		assert.Equal(t, want, got)
	}
	_, err := clickhouseGrantPrivileges("bogus")
	assert.Error(t, err)
}

func TestClickhouseAccessTypes(t *testing.T) {
	ch, _ := testClickhouse()
	ch.Database = "ch_access_db"
	admin := getClickhouseAdmin()

	roUser := &DatabaseUser{Username: "ch_ro_user", Password: "p", AccessType: ACCESS_TYPE_READONLY}
	rwUser := &DatabaseUser{Username: "ch_rw_user", Password: "p", AccessType: ACCESS_TYPE_READWRITE}
	cleanup := func() {
		_ = ch.deleteUser(context.TODO(), admin, roUser)
		_ = ch.deleteUser(context.TODO(), admin, rwUser)
		_ = ch.deleteDatabase(context.TODO(), admin)
	}
	cleanup()
	defer cleanup()

	assert.NoError(t, ch.createDatabase(context.TODO(), admin))
	assert.NoError(t, ch.createOrUpdateUser(context.TODO(), admin, roUser))
	assert.NoError(t, ch.createOrUpdateUser(context.TODO(), admin, rwUser))

	createTable := "CREATE TABLE ch_access_db.t (id Int32) ENGINE = MergeTree ORDER BY id"
	// readWrite has CREATE — readOnly does not.
	assert.NoError(t, ch.execAsUser(context.TODO(), createTable, rwUser))
	assert.Error(t, ch.execAsUser(context.TODO(),
		"CREATE TABLE ch_access_db.t2 (id Int32) ENGINE = MergeTree ORDER BY id", roUser))

	// Both access types can SELECT.
	for _, u := range []*DatabaseUser{roUser, rwUser} {
		res, err := ch.QueryAsUser(context.TODO(),
			"SELECT name FROM system.tables WHERE database = 'ch_access_db' AND name = 't'", u)
		assert.NoError(t, err)
		assert.Equal(t, "t", res)
	}

	// After revoke the readWrite user loses CREATE.
	assert.NoError(t, ch.revokePermissions(context.TODO(), admin, rwUser))
	assert.Error(t, ch.execAsUser(context.TODO(),
		"CREATE TABLE ch_access_db.t3 (id Int32) ENGINE = MergeTree ORDER BY id", rwUser))
}

func TestClickhouseRBAC(t *testing.T) {
	ch, _ := testClickhouse()
	ch.Database = "ch_rbac_db"
	admin := getClickhouseAdmin()
	user := &DatabaseUser{
		Username:   "ch_rbac_user",
		Password:   "p",
		AccessType: ACCESS_TYPE_READWRITE,
		CHQuota:    &CHQuota{IntervalSeconds: 3600, MaxQueries: 500, MaxResultRows: 100000},
		CHSettings: map[string]string{"max_memory_usage": "2000000"},
	}
	cleanup := func() {
		_ = ch.deleteUser(context.TODO(), admin, user)
		_ = ch.deleteDatabase(context.TODO(), admin)
	}
	cleanup()
	defer cleanup()

	quotaExists := func() bool {
		return ch.isRowExist(context.TODO(), "default",
			"SELECT name FROM system.quotas WHERE name = 'dbo_quota_ch_rbac_user'",
			admin.Username, admin.Password)
	}
	profileExists := func() bool {
		return ch.isRowExist(context.TODO(), "default",
			"SELECT name FROM system.settings_profiles WHERE name = 'dbo_profile_ch_rbac_user'",
			admin.Username, admin.Password)
	}

	assert.NoError(t, ch.createDatabase(context.TODO(), admin))
	assert.NoError(t, ch.createOrUpdateUser(context.TODO(), admin, user))
	assert.True(t, quotaExists(), "quota should exist after createOrUpdateUser")
	assert.True(t, profileExists(), "settings profile should exist after createOrUpdateUser")

	// CREATE ... OR REPLACE keeps a re-reconcile idempotent.
	assert.NoError(t, ch.createOrUpdateUser(context.TODO(), admin, user))
	assert.True(t, quotaExists())

	// revokePermissions drops the operator-managed RBAC objects.
	assert.NoError(t, ch.revokePermissions(context.TODO(), admin, user))
	assert.False(t, quotaExists(), "quota should be dropped after revoke")
	assert.False(t, profileExists(), "settings profile should be dropped after revoke")
}
