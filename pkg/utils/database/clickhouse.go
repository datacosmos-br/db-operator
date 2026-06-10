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
	"database/sql"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// ClickHouse is a database interface, abstracted object
// represents a database on ClickHouse instance
// can be used to execute queries to ClickHouse database
type ClickHouse struct {
	Host        string
	Port        uint16
	Database    string
	ClusterName string
	// Replicated, when true, creates the database with ENGINE = Replicated so
	// DDL auto-replicates across the cluster via ClickHouse Keeper. Requires
	// ClusterName to be set.
	Replicated bool
	// ZooKeeperPath overrides the Keeper/ZooKeeper znode path for a Replicated
	// database. Empty means the default /clickhouse/databases/<database>.
	ZooKeeperPath string
}

// Internal helpers, these functions are not part for the `Database` interface

// onCluster returns the " ON CLUSTER '<name>'" clause when the database belongs
// to a cluster, or an empty string otherwise. Every DDL statement that must
// reach all nodes appends it. Note: users, roles, quotas and settings profiles
// are server-scoped (not covered by the Replicated database engine), so their
// DDL must keep ON CLUSTER even when the database itself is Replicated.
func (ch ClickHouse) onCluster() string {
	if ch.ClusterName != "" {
		return fmt.Sprintf(" ON CLUSTER '%s'", ch.ClusterName)
	}
	return ""
}

// zooKeeperPath returns the Keeper znode path for a Replicated database,
// defaulting to /clickhouse/databases/<database> when not explicitly set.
func (ch ClickHouse) zooKeeperPath() string {
	if ch.ZooKeeperPath != "" {
		return ch.ZooKeeperPath
	}
	return fmt.Sprintf("/clickhouse/databases/%s", ch.Database)
}

// getDbConn opens a database/sql handle through the clickhouse-go v2 driver.
// OpenDB never returns an error; connection failures surface on first query,
// matching the previous sql.Open behaviour. The (*sql.DB, error) signature is
// kept so callers stay uniform with the postgres/mysql engines.
func (ch ClickHouse) getDbConn(dbname, user, password string) (*sql.DB, error) {
	return clickhouse.OpenDB(&clickhouse.Options{
		Addr: []string{fmt.Sprintf("%s:%d", ch.Host, ch.Port)},
		Auth: clickhouse.Auth{
			Database: dbname,
			Username: user,
			Password: password,
		},
		Protocol:    clickhouse.Native,
		DialTimeout: 30 * time.Second,
	}), nil
}

// executeExec runs a statement as the admin user. The context is honored, so a
// caller may attach per-query ClickHouse settings via clickhouse.Context.
func (ch ClickHouse) executeExec(ctx context.Context, database, query string, admin *DatabaseUser) error {
	log := log.FromContext(ctx)
	db, err := ch.getDbConn(database, admin.Username, admin.Password)
	if err != nil {
		log.Error(err, "failed to open a db connection")
		return err
	}

	defer db.Close()
	_, err = db.ExecContext(ctx, query)

	return err
}

func (ch ClickHouse) isDbExist(ctx context.Context, admin *DatabaseUser) bool {
	check := fmt.Sprintf("SELECT name FROM system.databases WHERE name = '%s'", ch.Database)

	return ch.isRowExist(ctx, "default", check, admin.Username, admin.Password)
}

func (ch ClickHouse) isUserExist(ctx context.Context, admin *DatabaseUser, user *DatabaseUser) bool {
	check := fmt.Sprintf("SELECT name FROM system.users WHERE name = '%s'", user.Username)

	return ch.isRowExist(ctx, "default", check, admin.Username, admin.Password)
}

func (ch ClickHouse) isRowExist(ctx context.Context, database, query, user, password string) bool {
	log := log.FromContext(ctx)
	db, err := ch.getDbConn(database, user, password)
	if err != nil {
		log.Error(err, "failed to open a db connection")
		return false
	}
	defer db.Close()

	var name string
	err = db.QueryRowContext(ctx, query).Scan(&name)
	if err != nil {
		log.V(2).Info("failed executing query", "error", err)
		return false
	}
	return true
}

// Functions that implement the `Database` interface

// CheckStatus checks status of ClickHouse database
// if the connection to database works
func (ch ClickHouse) CheckStatus(ctx context.Context, user *DatabaseUser) error {
	db, err := ch.getDbConn(ch.Database, user.Username, user.Password)
	if err != nil {
		return fmt.Errorf("db conn test failed - couldn't get db conn: %s", err)
	}
	defer db.Close()
	res, err := db.QueryContext(ctx, "SELECT 1")
	if err != nil {
		return fmt.Errorf("db conn test failed - failed to execute query: %s", err)
	}
	res.Close()

	return nil
}

// GetCredentials returns credentials of the ClickHouse database
func (ch ClickHouse) GetCredentials(ctx context.Context, user *DatabaseUser) Credentials {
	return Credentials{
		Name:     ch.Database,
		Username: user.Username,
		Password: user.Password,
	}
}

// ParseAdminCredentials parse admin username and password of ClickHouse database from secret data
func (ch ClickHouse) ParseAdminCredentials(ctx context.Context, data map[string][]byte) (*DatabaseUser, error) {
	admin := &DatabaseUser{}

	if user, ok := data["user"]; ok {
		admin.Username = string(user)
	} else {
		return nil, errors.New("no admin user found")
	}

	if password, ok := data["password"]; ok {
		admin.Password = string(password)
	} else {
		return nil, errors.New("no admin password found")
	}

	return admin, nil
}

func (ch ClickHouse) GetDatabaseAddress(ctx context.Context) DatabaseAddress {
	return DatabaseAddress{
		Host: ch.Host,
		Port: ch.Port,
	}
}

func (ch ClickHouse) QueryAsUser(ctx context.Context, query string, user *DatabaseUser) (string, error) {
	log := log.FromContext(ctx)
	db, err := ch.getDbConn(ch.Database, user.Username, user.Password)
	if err != nil {
		log.Error(err, "failed to open a db connection")
		return "", err
	}
	defer db.Close()

	var result string
	if err := db.QueryRowContext(ctx, query).Scan(&result); err != nil {
		log.Error(err, "failed executing query", "query", query)
		return "", err
	}
	return result, nil
}

func (ch ClickHouse) createDatabase(ctx context.Context, admin *DatabaseUser) error {
	log := log.FromContext(ctx)

	create := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`%s", ch.Database, ch.onCluster())
	execCtx := ctx
	if ch.Replicated {
		if ch.ClusterName == "" {
			return errors.New("clickhouse: replicated database requires clusterName to be set")
		}
		// {shard}/{replica} are ClickHouse macros, substituted per node from
		// each server's <macros> config — emitted as literal strings, not
		// interpolated by the operator.
		create += fmt.Sprintf(" ENGINE = Replicated('%s', '{shard}', '{replica}')", ch.zooKeeperPath())
		// The Replicated database engine is experimental on ClickHouse 24.x and
		// GA from 25.x. Enabling the setting works across the whole supported
		// range — it is an accepted no-op where the engine is already GA. The
		// setting is scoped to this one statement via clickhouse.Context.
		execCtx = clickhouse.Context(ctx, clickhouse.WithSettings(clickhouse.Settings{
			"allow_experimental_database_replicated": 1,
		}))
	}

	if !ch.isDbExist(ctx, admin) {
		if err := ch.executeExec(execCtx, "default", create, admin); err != nil {
			log.Error(err, "failed creating ClickHouse database")
			return err
		}
	}

	return nil
}

func (ch ClickHouse) deleteDatabase(ctx context.Context, admin *DatabaseUser) error {
	log := log.FromContext(ctx)
	drop := fmt.Sprintf("DROP DATABASE IF EXISTS `%s`%s", ch.Database, ch.onCluster())

	if ch.isDbExist(ctx, admin) {
		if err := ch.executeExec(ctx, "default", drop, admin); err != nil {
			log.Error(err, "failed dropping ClickHouse database")
			return err
		}
	}
	return nil
}

func (ch ClickHouse) createOrUpdateUser(ctx context.Context, admin *DatabaseUser, user *DatabaseUser) error {
	// ClickHouse `local_directory` access storage is PER-NODE: a single
	// (load-balanced) isUserExist probe cannot observe a partial/per-node state, so
	// branching create-vs-update on it leaves users missing on some nodes (and loops
	// on ALTER → UNKNOWN_USER). Instead, ALWAYS ensure existence on EVERY node with
	// CREATE USER IF NOT EXISTS ON CLUSTER (idempotent — skips nodes that already have
	// the user, heals nodes that don't), then enforce the identity with ALTER ON
	// CLUSTER. Both are idempotent; together they converge deterministically to
	// "exists on all nodes" regardless of prior partial state.
	if err := ch.createUser(ctx, admin, user); err != nil {
		return err
	}
	if err := ch.updateUser(ctx, admin, user); err != nil {
		return err
	}
	return ch.setUserPermission(ctx, admin, user)
}

// hostClause returns a " HOST REGEXP '<pattern>'" clause restricting which
// client hosts the user may connect from, or an empty string when unrestricted.
func hostClause(user *DatabaseUser) string {
	if user.HostRegexp != "" {
		return fmt.Sprintf(" HOST REGEXP '%s'", user.HostRegexp)
	}
	return ""
}

// identifiedClause returns the IDENTIFIED clause for CREATE/ALTER USER. When the
// user carries a PasswordHash it adopts that exact credential via
// IDENTIFIED WITH sha256_hash (no plaintext is needed or logged); otherwise it
// sets the plaintext password the operator generated/holds.
func identifiedClause(user *DatabaseUser) string {
	if user.PasswordHash != "" {
		return fmt.Sprintf(" IDENTIFIED WITH sha256_hash BY '%s'", user.PasswordHash)
	}
	return fmt.Sprintf(" IDENTIFIED BY '%s'", user.Password)
}

func (ch ClickHouse) createUser(ctx context.Context, admin *DatabaseUser, user *DatabaseUser) error {
	log := log.FromContext(ctx)
	create := fmt.Sprintf("CREATE USER IF NOT EXISTS '%s'%s%s%s",
		user.Username, ch.onCluster(), identifiedClause(user), hostClause(user))

	if err := ch.executeExec(ctx, "default", create, admin); err != nil {
		log.Error(err, "failed creating ClickHouse user")
		return err
	}

	return nil
}

func (ch ClickHouse) updateUser(ctx context.Context, admin *DatabaseUser, user *DatabaseUser) error {
	log := log.FromContext(ctx)
	update := fmt.Sprintf("ALTER USER '%s'%s%s%s",
		user.Username, ch.onCluster(), identifiedClause(user), hostClause(user))

	if err := ch.executeExec(ctx, "default", update, admin); err != nil {
		log.Error(err, "failed updating ClickHouse user")
		return err
	}

	return nil
}

// clickhouseGrantPrivileges maps an access type to the ClickHouse privilege
// list granted on the database. readWrite covers data read/write plus table
// management within the database; ClickHouse expands CREATE/DROP into their
// sub-privileges automatically.
func clickhouseGrantPrivileges(accessType string) (string, error) {
	switch accessType {
	case ACCESS_TYPE_READONLY:
		return "SELECT", nil
	case ACCESS_TYPE_READWRITE:
		return "SELECT, INSERT, ALTER, CREATE, DROP", nil
	case ACCESS_TYPE_MAINUSER:
		return "ALL", nil
	default:
		return "", fmt.Errorf("unknown access type: %s", accessType)
	}
}

func (ch ClickHouse) setUserPermission(ctx context.Context, admin *DatabaseUser, user *DatabaseUser) error {
	log := log.FromContext(ctx)

	privileges, err := clickhouseGrantPrivileges(user.AccessType)
	if err != nil {
		return err
	}
	grant := fmt.Sprintf("GRANT%s %s ON `%s`.* TO '%s'",
		ch.onCluster(), privileges, ch.Database, user.Username)
	if err := ch.executeExec(ctx, ch.Database, grant, admin); err != nil {
		log.Error(err, "failed granting privileges to ClickHouse user")
		return err
	}

	// ExtraPrivileges are pre-existing RBAC roles — the operator grants them
	// but never creates them (consistent with the postgres/mysql engines).
	for _, role := range user.ExtraPrivileges {
		grantRole := fmt.Sprintf("GRANT%s %s TO '%s'", ch.onCluster(), role, user.Username)
		if err := ch.executeExec(ctx, ch.Database, grantRole, admin); err != nil {
			log.Error(err, "failed granting role to ClickHouse user", "role", role)
			return err
		}
	}

	// ExtraGrants are explicit privilege grants beyond the per-database grant,
	// e.g. cluster-wide (REMOTE, CLUSTER ON *.*) or system-table privileges.
	for _, g := range user.ExtraGrants {
		grant := fmt.Sprintf("GRANT%s %s ON %s TO '%s'",
			ch.onCluster(), g.Privileges, g.On, user.Username)
		if err := ch.executeExec(ctx, ch.Database, grant, admin); err != nil {
			log.Error(err, "failed granting extra privilege to ClickHouse user", "on", g.On)
			return err
		}
	}

	if err := ch.applyQuota(ctx, admin, user); err != nil {
		log.Error(err, "failed applying ClickHouse quota")
		return err
	}
	if err := ch.applySettingsProfile(ctx, admin, user); err != nil {
		log.Error(err, "failed applying ClickHouse settings profile")
		return err
	}

	return nil
}

// chQuotaName / chProfileName build deterministic names for the operator-managed
// RBAC objects so they can be created (OR REPLACE) and dropped idempotently.
func chQuotaName(user string) string   { return "dbo_quota_" + user }
func chProfileName(user string) string { return "dbo_profile_" + user }

// applyQuota creates (OR REPLACE) the user's resource quota. No-op when the
// user has no quota configured.
func (ch ClickHouse) applyQuota(ctx context.Context, admin *DatabaseUser, user *DatabaseUser) error {
	if user.Quota == nil {
		return nil
	}
	q := user.Quota
	limits := []string{}
	if q.MaxQueries > 0 {
		limits = append(limits, fmt.Sprintf("MAX queries = %d", q.MaxQueries))
	}
	if q.MaxResultRows > 0 {
		limits = append(limits, fmt.Sprintf("MAX result_rows = %d", q.MaxResultRows))
	}
	if q.MaxExecutionTimeSeconds > 0 {
		limits = append(limits, fmt.Sprintf("MAX execution_time = %d", q.MaxExecutionTimeSeconds))
	}
	stmt := fmt.Sprintf("CREATE QUOTA OR REPLACE %s%s FOR INTERVAL %d SECOND",
		chQuotaName(user.Username), ch.onCluster(), q.IntervalSeconds)
	if len(limits) > 0 {
		stmt += " " + strings.Join(limits, ", ")
	}
	stmt += fmt.Sprintf(" TO '%s'", user.Username)
	return ch.executeExec(ctx, ch.Database, stmt, admin)
}

// applySettingsProfile creates (OR REPLACE) a settings profile from the user's
// settings and assigns it. No-op when the user has no settings configured.
func (ch ClickHouse) applySettingsProfile(ctx context.Context, admin *DatabaseUser, user *DatabaseUser) error {
	if len(user.Settings) == 0 {
		return nil
	}
	parts := make([]string, 0, len(user.Settings))
	for k, v := range user.Settings {
		parts = append(parts, fmt.Sprintf("%s = %s", k, v))
	}
	sort.Strings(parts) // deterministic SQL across reconciles
	create := fmt.Sprintf("CREATE SETTINGS PROFILE OR REPLACE %s%s SETTINGS %s",
		chProfileName(user.Username), ch.onCluster(), strings.Join(parts, ", "))
	if err := ch.executeExec(ctx, ch.Database, create, admin); err != nil {
		return err
	}
	assign := fmt.Sprintf("ALTER USER '%s'%s SETTINGS PROFILE %s",
		user.Username, ch.onCluster(), chProfileName(user.Username))
	return ch.executeExec(ctx, ch.Database, assign, admin)
}

// dropRBACObjects removes the operator-managed quota and settings profile of a
// user. Safe to call unconditionally — both statements use IF EXISTS.
func (ch ClickHouse) dropRBACObjects(ctx context.Context, admin *DatabaseUser, username string) error {
	dropQuota := fmt.Sprintf("DROP QUOTA IF EXISTS %s%s", chQuotaName(username), ch.onCluster())
	if err := ch.executeExec(ctx, ch.Database, dropQuota, admin); err != nil {
		return err
	}
	dropProfile := fmt.Sprintf("DROP SETTINGS PROFILE IF EXISTS %s%s", chProfileName(username), ch.onCluster())
	return ch.executeExec(ctx, ch.Database, dropProfile, admin)
}

func (ch ClickHouse) revokePermissions(ctx context.Context, admin *DatabaseUser, user *DatabaseUser) error {
	log := log.FromContext(ctx)
	if !ch.isUserExist(ctx, admin, user) {
		return nil
	}

	privileges, err := clickhouseGrantPrivileges(user.AccessType)
	if err != nil {
		return err
	}
	revoke := fmt.Sprintf("REVOKE%s %s ON `%s`.* FROM '%s'",
		ch.onCluster(), privileges, ch.Database, user.Username)
	if err := ch.executeExec(ctx, ch.Database, revoke, admin); err != nil {
		log.Error(err, "failed revoking privileges from ClickHouse user")
		return err
	}

	for _, g := range user.ExtraGrants {
		revoke := fmt.Sprintf("REVOKE%s %s ON %s FROM '%s'",
			ch.onCluster(), g.Privileges, g.On, user.Username)
		if err := ch.executeExec(ctx, ch.Database, revoke, admin); err != nil {
			log.Error(err, "failed revoking extra privilege from ClickHouse user", "on", g.On)
			return err
		}
	}

	for _, role := range user.ExtraPrivileges {
		revokeRole := fmt.Sprintf("REVOKE%s %s FROM '%s'", ch.onCluster(), role, user.Username)
		if err := ch.executeExec(ctx, ch.Database, revokeRole, admin); err != nil {
			log.Error(err, "failed revoking role from ClickHouse user", "role", role)
			return err
		}
	}

	if err := ch.dropRBACObjects(ctx, admin, user.Username); err != nil {
		log.Error(err, "failed dropping ClickHouse RBAC objects")
		return err
	}

	return nil
}

func (ch ClickHouse) execAsUser(ctx context.Context, query string, user *DatabaseUser) error {
	log := log.FromContext(ctx)
	db, err := ch.getDbConn(ch.Database, user.Username, user.Password)
	if err != nil {
		log.Error(err, "failed to open a db connection")
		return err
	}

	defer db.Close()
	_, err = db.Exec(query)

	return err
}

func (ch ClickHouse) deleteUser(ctx context.Context, admin *DatabaseUser, user *DatabaseUser) error {
	log := log.FromContext(ctx)

	// DROP USER removes grants but not named quotas/profiles — clean them up
	// explicitly so the generated-user delete path leaves nothing behind.
	if err := ch.dropRBACObjects(ctx, admin, user.Username); err != nil {
		log.Error(err, "failed dropping ClickHouse RBAC objects")
		return err
	}

	drop := fmt.Sprintf("DROP USER IF EXISTS '%s'%s", user.Username, ch.onCluster())
	if err := ch.executeExec(ctx, "default", drop, admin); err != nil {
		log.Error(err, "failed deleting ClickHouse user")
		return err
	}

	return nil
}
