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

func postgresGrantFixture(t *testing.T, schemas []string, dropPublic bool, react func(query string) error) (*Postgres, *DatabaseUser, *DatabaseUser, *[]string) {
	t.Helper()
	queries := []string{}
	admin := &DatabaseUser{Username: "admin", Password: "admin"}
	user := &DatabaseUser{
		Username:   "cosmos-funmotors-cosmos-funmotors-pmmgrafana-pg-db",
		AccessType: ACCESS_TYPE_MAINUSER,
	}
	p := &Postgres{
		Database:         "cosmos-funmotors-pmmgrafana-pg-db",
		DropPublicSchema: dropPublic,
		Schemas:          schemas,
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

func TestUnitPostgresGrantCreatesRequestedSchemaOn3F000(t *testing.T) {
	t.Parallel()
	schemaGrants := 0
	p, user, admin, queries := postgresGrantFixture(t, []string{"reporting"}, true, func(query string) error {
		if strings.HasPrefix(query, "GRANT ALL ON SCHEMA") {
			schemaGrants++
			if schemaGrants == 1 {
				return &pq.Error{Code: "3F000", Message: "schema reporting does not exist"}
			}
		}
		return nil
	})

	err := p.setUserPermission(context.Background(), admin, user)
	require.NoError(t, err)

	var creates []string
	var schemaGrantsSeen []string
	for _, query := range *queries {
		assert.NotContains(t, query, "CASCADE")
		assert.NotContains(t, query, "DROP DATABASE")
		assert.NotContains(t, query, "DROP SCHEMA")
		if strings.HasPrefix(query, "CREATE SCHEMA") {
			creates = append(creates, query)
		}
		if strings.HasPrefix(query, "GRANT ALL ON SCHEMA") {
			schemaGrantsSeen = append(schemaGrantsSeen, query)
		}
	}
	require.Equal(t, []string{`CREATE SCHEMA IF NOT EXISTS "reporting";`}, creates)
	require.Equal(t, []string{
		`GRANT ALL ON SCHEMA "reporting" TO "cosmos-funmotors-cosmos-funmotors-pmmgrafana-pg-db"`,
		`GRANT ALL ON SCHEMA "reporting" TO "cosmos-funmotors-cosmos-funmotors-pmmgrafana-pg-db"`,
	}, schemaGrantsSeen)
}

func TestUnitPostgresGrantCreatesImplicitPublicSchemaOn3F000(t *testing.T) {
	t.Parallel()
	schemaGrants := 0
	p, user, admin, queries := postgresGrantFixture(t, nil, false, func(query string) error {
		if strings.HasPrefix(query, "GRANT ALL ON SCHEMA") {
			schemaGrants++
			if schemaGrants == 1 {
				return &pq.Error{Code: "3F000", Message: "schema public does not exist"}
			}
		}
		return nil
	})

	err := p.setUserPermission(context.Background(), admin, user)
	require.NoError(t, err)

	foundCreate := false
	for _, query := range *queries {
		if query == `CREATE SCHEMA IF NOT EXISTS "public";` {
			foundCreate = true
		}
	}
	assert.True(t, foundCreate)
	assert.Equal(t, 2, schemaGrants)
}

func TestUnitPostgresGrantReturnsNon3F000(t *testing.T) {
	t.Parallel()
	p, user, admin, queries := postgresGrantFixture(t, nil, false, func(query string) error {
		if strings.HasPrefix(query, "GRANT ALL ON SCHEMA") {
			return &pq.Error{Code: "42501", Message: "permission denied"}
		}
		return nil
	})

	err := p.setUserPermission(context.Background(), admin, user)
	require.Error(t, err)
	assert.False(t, postgresInvalidSchemaName(err))

	for _, query := range *queries {
		assert.NotContains(t, query, "CREATE SCHEMA")
		assert.NotContains(t, query, "CASCADE")
	}
}

func TestUnitPostgresGrantRetriesSchemaOnce(t *testing.T) {
	t.Parallel()
	p, user, admin, queries := postgresGrantFixture(t, nil, false, func(query string) error {
		if strings.HasPrefix(query, "GRANT ALL ON SCHEMA") {
			return &pq.Error{Code: "3F000", Message: "schema public does not exist"}
		}
		return nil
	})

	err := p.setUserPermission(context.Background(), admin, user)
	require.Error(t, err)
	assert.True(t, postgresInvalidSchemaName(err))

	creates := 0
	schemaGrants := 0
	for _, query := range *queries {
		if strings.HasPrefix(query, "CREATE SCHEMA") {
			creates++
		}
		if strings.HasPrefix(query, "GRANT ALL ON SCHEMA") {
			schemaGrants++
		}
	}
	assert.Equal(t, 1, creates)
	assert.Equal(t, 2, schemaGrants)
}
