/*
 * Copyright 2023 DB-Operator Authors
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

package consts

// This package exists to avoid cycle import. Put here consts that are used across packages

// Database Related Consts
const (
	POSTGRES_DB         = "POSTGRES_DB"
	POSTGRES_USER       = "POSTGRES_USER"
	POSTGRES_PASSWORD   = "POSTGRES_PASSWORD"
	MYSQL_DB            = "DB"
	MYSQL_USER          = "USER"
	MYSQL_PASSWORD      = "PASSWORD"
	MONGODB_DB          = "MONGODB_DB"
	MONGODB_USER        = "MONGODB_USER"
	MONGODB_PASSWORD    = "MONGODB_PASSWORD"
	CLICKHOUSE_DB       = "CLICKHOUSE_DB"
	CLICKHOUSE_USER     = "CLICKHOUSE_USER"
	CLICKHOUSE_PASSWORD = "CLICKHOUSE_PASSWORD"
	ORACLE_DB           = "ORACLE_DB"
	ORACLE_USER         = "ORACLE_USER"
	ORACLE_PASSWORD     = "ORACLE_PASSWORD"
	SQLSERVER_DB        = "SQLSERVER_DB"
	SQLSERVER_USER      = "SQLSERVER_USER"
	SQLSERVER_PASSWORD  = "SQLSERVER_PASSWORD"
)

// Database engines
const (
	ENGINE_POSTGRES   = "postgres"
	ENGINE_MYSQL      = "mysql"
	ENGINE_MONGODB    = "mongodb"
	ENGINE_CLICKHOUSE = "clickhouse"
	ENGINE_ORACLE     = "oracle"
	ENGINE_SQLSERVER  = "sqlserver"
)

// SSL modes
const (
	SSL_DISABLED  = "disabled"
	SSL_REQUIRED  = "required"
	SSL_VERIFY_CA = "verify_ca"
)

// Kubernetes Annotations
const (
	TEMPLATE_ANNOTATION_KEY       = "kinda.rocks/db-operator-templated-keys"
	SECRET_FORCE_RECONCILE        = "kinda.rocks/secret-force-reconcile"
	DATABASE_FORCE_FULL_RECONCILE = "kinda.rocks/db-force-full-reconcile"
	USED_OBJECTS                  = "kinda.rocks/used-objects"
)

// Kubernetes Labels
const (
	MANAGED_BY_LABEL_KEY   = "app.kubernetes.io/managed-by"
	MANAGED_BY_LABEL_VALUE = "db-operator"
	USED_BY_KIND_LABEL_KEY = "kinda.rocks/used-by-kind"
	USED_BY_NAME_LABEL_KEY = "kinda.rocks/used-by-name"
)

// Privileges

const ALL_PRIVILEGES = "ALL PRIVILEGES"
