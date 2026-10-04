/*
 * Copyright 2026 DB-Operator Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
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
	"strings"
	"testing"

	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func postgresDeleteFixture(t *testing.T, react func(query string) error) (*Postgres, *DatabaseUser, *DatabaseUser, *[]string) {
	t.Helper()
	queries := []string{}
	admin := &DatabaseUser{Username: "admin", Password: "admin"}
	user := &DatabaseUser{
		Username:   "funmotors-pg-grafana",
		AccessType: ACCESS_TYPE_READWRITE,
	}
	p := &Postgres{
		Database: "funmotors",
		MainUser: &DatabaseUser{Username: "main"},
		testUserExists: func(context.Context, *DatabaseUser, *DatabaseUser) bool {
			return true
		},
		testExec: func(_ context.Context, _ string, query string, _ *DatabaseUser) error {
			queries = append(queries, query)
			if react == nil {
				return nil
			}
			return react(query)
		},
	}
	return p, user, admin, &queries
}

func TestUnitPostgresDeleteUserDropsDependentsOn2BP01(t *testing.T) {
	t.Parallel()
	cascaded := false
	p, user, admin, queries := postgresDeleteFixture(t, func(query string) error {
		if strings.Contains(query, "CASCADE") {
			cascaded = true
			return nil
		}
		if strings.Contains(query, "DROP OWNED BY") || strings.HasPrefix(query, "DROP USER") {
			if !cascaded {
				return &pq.Error{
					Code:    "2BP01",
					Message: "cannot drop objects because other objects depend on them",
				}
			}
		}
		return nil
	})

	err := DeleteUser(context.Background(), p, user, admin)
	require.NoError(t, err)

	var dropOwned []string
	dropUser := 0
	for _, query := range *queries {
		if strings.Contains(query, "DROP OWNED BY") {
			dropOwned = append(dropOwned, query)
		}
		if strings.HasPrefix(query, "DROP USER") {
			dropUser++
		}
	}
	require.Len(t, dropOwned, 2)
	assert.Equal(t, `DROP OWNED BY "funmotors-pg-grafana";`, dropOwned[0])
	assert.Equal(t, `DROP OWNED BY "funmotors-pg-grafana" CASCADE;`, dropOwned[1])
	assert.Equal(t, 2, dropUser)
}

func TestUnitPostgresRevokePermissionsDoesNotCascade(t *testing.T) {
	t.Parallel()
	p, user, admin, queries := postgresDeleteFixture(t, func(query string) error {
		if strings.Contains(query, "DROP OWNED BY") {
			return &pq.Error{
				Code:    "2BP01",
				Message: "cannot drop objects because other objects depend on them",
			}
		}
		return nil
	})

	err := RevokePermissions(context.Background(), p, user, admin)
	require.Error(t, err)
	assert.True(t, postgresDependentObjectsStillExist(err))

	foundDropOwned := false
	for _, query := range *queries {
		assert.NotContains(t, query, "CASCADE")
		if strings.Contains(query, "DROP OWNED BY") {
			foundDropOwned = true
			assert.Equal(t, `DROP OWNED BY "funmotors-pg-grafana";`, query)
		}
	}
	assert.True(t, foundDropOwned)
}

func TestUnitPostgresDeleteUserDoesNotCascadeOnOtherErrors(t *testing.T) {
	t.Parallel()
	p, user, admin, queries := postgresDeleteFixture(t, func(query string) error {
		if strings.HasPrefix(query, "DROP USER") {
			return &pq.Error{Code: "42501", Message: "permission denied"}
		}
		return nil
	})

	err := DeleteUser(context.Background(), p, user, admin)
	require.Error(t, err)
	assert.False(t, postgresDependentObjectsStillExist(err))
	for _, query := range *queries {
		assert.NotContains(t, query, "CASCADE")
	}
}

func TestUnitPostgresDeleteUserReturnsErrorWhenDependentsRemain(t *testing.T) {
	t.Parallel()
	p, user, admin, queries := postgresDeleteFixture(t, func(query string) error {
		if strings.Contains(query, "DROP OWNED BY") || strings.HasPrefix(query, "DROP USER") {
			return &pq.Error{
				Code:    "2BP01",
				Message: "cannot drop objects because other objects depend on them",
			}
		}
		return nil
	})

	err := DeleteUser(context.Background(), p, user, admin)
	require.Error(t, err)
	assert.True(t, postgresDependentObjectsStillExist(err))

	foundCascade := false
	for _, query := range *queries {
		if strings.Contains(query, "CASCADE") {
			foundCascade = true
		}
	}
	assert.True(t, foundCascade)
}
