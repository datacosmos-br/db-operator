/*
 * Copyright 2021 kloeckner.i GmbH
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

import "context"

// CreateDatabase executes queries to create database
func CreateDatabase(ctx context.Context, db Database, admin *DatabaseUser) error {
	err := db.createDatabase(ctx, admin)
	if err != nil {
		return err
	}
	return nil
}

// DeleteDatabase executes queries to delete database and user
func DeleteDatabase(ctx context.Context, db Database, admin *DatabaseUser) error {
	err := db.deleteDatabase(ctx, admin)
	if err != nil {
		return err
	}

	return nil
}

// CreateOrUpdateUser executes queries to create or update user
func CreateOrUpdateUser(ctx context.Context, db Database, dbuser *DatabaseUser, admin *DatabaseUser) error {
	err := db.createOrUpdateUser(ctx, admin, dbuser)
	if err != nil {
		return err
	}
	return nil
}

// CreateUser executes queries to a create user
func CreateUser(ctx context.Context, db Database, dbuser *DatabaseUser, admin *DatabaseUser) error {
	err := db.createUser(ctx, admin, dbuser)
	if err != nil {
		return err
	}
	return nil
}

func UpdateUser(ctx context.Context, db Database, dbuser *DatabaseUser, admin *DatabaseUser) error {
	err := db.createOrUpdateUser(ctx, admin, dbuser)
	if err != nil {
		return err
	}
	return nil
}

func DeleteUser(ctx context.Context, db Database, dbuser *DatabaseUser, admin *DatabaseUser) error {
	err := db.deleteUser(ctx, admin, dbuser)
	if err != nil {
		return err
	}

	return nil
}

// New returns database interface according to engine type
func New(engine string) Database {
	switch engine {
	case "postgres":
		return &Postgres{}
	case "mysql":
		return &Mysql{}
	case "clickhouse":
		return &ClickHouse{}
	case "dummy":
		return &Dummy{}
	}

	return nil
}
