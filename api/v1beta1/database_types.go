/*
 * Copyright 2021 kloeckner.i GmbH
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

package v1beta1

import (
	"fmt"

	"github.com/db-operator/db-operator/v2/pkg/consts"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// DatabaseSpec defines the desired state of Database
type DatabaseSpec struct {
	SecretName        string            `json:"secretName"`
	Instance          string            `json:"instance"`
	DeletionProtected bool              `json:"deletionProtected"`
	Backup            DatabaseBackup    `json:"backup,omitempty"`
	SecretsTemplates  map[string]string `json:"secretsTemplates,omitempty"`
	Postgres          Postgres          `json:"postgres,omitempty"`
	Clickhouse        Clickhouse        `json:"clickhouse,omitempty"`
	Cleanup           bool              `json:"cleanup,omitempty"`
	Credentials       Credentials       `json:"credentials,omitempty"`
	ExtraGrants       []*ExtraGrant     `json:"extraGrants,omitempty"`
	// If specified, DB Operator will try to use an existing user to assign permissions
	// User will not be removed, when a database is removed, but the permissions added by the
	// operator will be cleaned up
	ExistingUser string `json:"existingUser,omitempty"`
	// DatabaseName overrides the name of the database that gets created. When
	// empty, the operator generates it as "<namespace>-<name>". Set it to
	// adopt a pre-existing database or to keep an externally-defined name.
	// +optional
	DatabaseName string `json:"databaseName,omitempty"`
}

type ExtraGrant struct {
	User       string `json:"user"`
	AccessType string `json:"accessType"`
}

// Postgres struct should be used to provide resource that only applicable to postgres
type Postgres struct {
	Extensions []string `json:"extensions,omitempty"`
	// If set to true, the public schema will be dropped after the database creation
	DropPublicSchema bool `json:"dropPublicSchema,omitempty"`
	// Specify schemas to be created. The user created by db-operator will have all access on them.
	Schemas []string `json:"schemas,omitempty"`
	// Let user create database from template
	Template string `json:"template,omitempty"`
	// Owner is the role name that becomes the database owner via
	// ALTER DATABASE ... OWNER TO. The operator also reassigns the public
	// schema owner so PG15+ rules (CREATE on public follows schema
	// ownership, not PUBLIC role) let the application user create
	// objects directly. Empty preserves the legacy behaviour where the
	// admin role keeps ownership.
	// +optional
	Owner string `json:"owner,omitempty"`
	// PostInitSQL is an ordered list of idempotent SQL statements executed
	// as the admin against this database AFTER creation/extensions/
	// schemas/owner. Use it for bootstrap DDL such as
	// "ALTER SCHEMA public OWNER TO ...", "GRANT CREATE ON SCHEMA public
	// TO ...", or product-specific seed grants. Each statement runs in
	// its own exec; a failure aborts reconciliation so callers can fix
	// and retry.
	// +optional
	PostInitSQL []string `json:"postInitSQL,omitempty"`
	// CrossDatabaseGrants applies the per-user permission flow (the same
	// READONLY/READWRITE/MAIN grant logic used for the owning database)
	// against OTHER databases on the same instance. Use when an
	// application user must access a sibling database it does not own —
	// e.g. airbyte-db needs READWRITE on airbyte-visibility for temporal
	// schema setup. The target user MUST already exist; this field only
	// grants existing roles. Grants apply on the target database's
	// public schema only.
	// +optional
	CrossDatabaseGrants []CrossDatabaseGrant `json:"crossDatabaseGrants,omitempty"`
}

// CrossDatabaseGrant grants a known role access on sibling databases.
// Mirrors DbUser AccessType semantics: "readOnly", "readWrite", "main".
type CrossDatabaseGrant struct {
	// Username is the role that receives the grants. Must exist on the
	// instance — the operator does not create or rotate it here.
	Username string `json:"username"`
	// Databases is the list of sibling database names on which to grant
	// access. Each one runs the standard per-database permission flow.
	Databases []string `json:"databases"`
	// AccessType is one of "readOnly", "readWrite", or "main" — matches
	// the DbUser AccessType vocabulary used by setUserPermission.
	AccessType string `json:"accessType"`
}

// Clickhouse struct should be used to provide resource that only applicable to ClickHouse
type Clickhouse struct {
	// ClusterName is the name of the ClickHouse cluster (used for ON CLUSTER queries)
	ClusterName string `json:"clusterName,omitempty"`
	// Replicated, when true and clusterName is set, creates the database with
	// ENGINE = Replicated so DDL auto-replicates across the cluster via ClickHouse Keeper.
	// +kubebuilder:default=false
	Replicated bool `json:"replicated,omitempty"`
	// ZooKeeperPath overrides the Keeper/ZooKeeper znode path for a Replicated
	// database. Defaults to /clickhouse/databases/<database> when empty.
	// +optional
	ZooKeeperPath string `json:"zooKeeperPath,omitempty"`
	// HostRegexp restricts the database's main user to client hosts whose name
	// matches this ClickHouse HOST REGEXP pattern. Empty means no restriction.
	// +optional
	HostRegexp string `json:"hostRegexp,omitempty"`
	// ExtraGrants are additional privilege grants applied to the main user
	// beyond the per-database grant — e.g. cluster-wide privileges
	// (REMOTE, CLUSTER ON *.*) or system-table grants.
	// +optional
	ExtraGrants []ClickhouseGrant `json:"extraGrants,omitempty"`
}

// ClickhouseGrant is an extra ClickHouse privilege grant for a user, rendered
// as: GRANT <privileges> ON <on> TO <user>.
type ClickhouseGrant struct {
	// Privileges is the comma-separated ClickHouse privilege list,
	// e.g. "REMOTE, CLUSTER" or "SELECT(cluster)".
	Privileges string `json:"privileges"`
	// On is the grant target, e.g. "*.*" or "system.clusters".
	On string `json:"on"`
}

// DatabaseStatus defines the observed state of Database
type DatabaseStatus struct {
	// Important: Run "make generate" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	Status                bool                `json:"status"`
	MonitorUserSecretName string              `json:"monitorUserSecret,omitempty"`
	ProxyStatus           DatabaseProxyStatus `json:"proxyStatus,omitempty"`
	DatabaseName          string              `json:"database"`
	UserName              string              `json:"user"`
	Engine                string              `json:"engine"`
	OperatorVersion       string              `json:"operatorVersion,omitempty"`
	ExtraGrants           []*ExtraGrant       `json:"extraGrants,omitempty"`
}

// DatabaseProxyStatus defines whether proxy for database is enabled or not
// if so, provide information
type DatabaseProxyStatus struct {
	Status      bool   `json:"status"`
	ServiceName string `json:"serviceName"`
	SQLPort     int32  `json:"sqlPort"`
}

// DatabaseBackup defines the desired state of backup and schedule
type DatabaseBackup struct {
	Enable bool   `json:"enable"`
	Cron   string `json:"cron"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=db
// +kubebuilder:printcolumn:name="Status",type=boolean,JSONPath=`.status.status`,description="current db status"
// +kubebuilder:printcolumn:name="Protected",type=boolean,JSONPath=`.spec.deletionProtected`,description="If database is protected to not get deleted."
// +kubebuilder:printcolumn:name="DBInstance",type=string,JSONPath=`.spec.instance`,description="instance reference"
// +kubebuilder:printcolumn:name="OperatorVersion",type=string,JSONPath=`.status.operatorVersion`,description="db-operator version of last full reconcile"
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`,description="time since creation of resource"
// +kubebuilder:storageversion
// Database is the Schema for the databases API
type Database struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DatabaseSpec   `json:"spec,omitempty"`
	Status DatabaseStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// DatabaseList contains a list of Database
type DatabaseList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Database `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Database{}, &DatabaseList{})
}

// GetProtocol returns the protocol that is required for connection (postgresql, mysql or clickhouse)
func (db *Database) GetProtocol() (string, error) {
	switch db.Status.Engine {
	case consts.ENGINE_POSTGRES:
		return "postgresql", nil
	case consts.ENGINE_MYSQL:
		return db.Status.Engine, nil
	case consts.ENGINE_CLICKHOUSE:
		// "clickhouse" is also the DSN scheme of the clickhouse-go v2 driver.
		return db.Status.Engine, nil
	default:
		return "", fmt.Errorf("unknown engine %s", db.Status.Engine)
	}
}

func (db *Database) IsCleanup() bool {
	return db.Spec.Cleanup
}

func (db *Database) IsDeleted() bool {
	return db.GetDeletionTimestamp() != nil
}

func (db *Database) GetSecretName() string {
	return db.Spec.SecretName
}

func (db *Database) ToClientObject() client.Object {
	return db
}

// AccessSecretName returns string value to define name of the secret resource for accessing instance
func (db *Database) InstanceAccessSecretName() string {
	return "dbin-" + db.Spec.Instance + "-access-secret"
}

func (extraGrant *ExtraGrant) IsExtraGrant(extraGrants []*ExtraGrant) bool {
	for _, existingExtraGrant := range extraGrants {
		if existingExtraGrant.User == extraGrant.User &&
			existingExtraGrant.AccessType == extraGrant.AccessType {
			return true
		}
	}
	return false
}

// Function to mark the Database as a hub
func (db *Database) Hub() {}
