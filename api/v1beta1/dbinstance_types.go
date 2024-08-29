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
	"context"
	"errors"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// DbInstanceSpec defines the desired state of DbInstance
type DbInstanceSpec struct {
	// Important: Run "make generate" to regenerate code after modifying this file
	Engine          string                  `json:"engine"`
	AdminUserSecret NamespacedName          `json:"adminSecretRef"`
	Backup          DbInstanceBackup        `json:"backup,omitempty"`
	Monitoring      DbInstanceMonitoring    `json:"monitoring,omitempty"`
	SSLConnection   DbInstanceSSLConnection `json:"sslConnection,omitempty"`
	// A list of privileges that are allowed to be set as Dbuser's extra privileges
	AllowedPrivileges []string `json:"allowedPrivileges,omitempty"`
	DbInstanceSource  `json:",inline"`
}

// DbInstanceSource represents the source of an instance.
// Only one of its members may be specified.
type DbInstanceSource struct {
	Google    *GoogleInstance    `json:"google,omitempty" protobuf:"bytes,1,opt,name=google"`
	Generic   *GenericInstance   `json:"generic,omitempty" protobuf:"bytes,2,opt,name=generic"`
	OracleOCI *OracleOCIInstance `json:"oracleoci,omitempty" protobuf:"bytes,3,opt,name=oracleoci"`
	Azure     *AzureInstance     `json:"azure,omitempty" protobuf:"bytes,4,opt,name=azure"`
	AWS       *AWSInstance       `json:"aws,omitempty" protobuf:"bytes,5,opt,name=aws"`
}

// DbInstanceStatus defines the observed state of DbInstance
type DbInstanceStatus struct {
	// Important: Run "make generate" to regenerate code after modifying this file
	Phase     string            `json:"phase"`
	Status    bool              `json:"status"`
	Info      map[string]string `json:"info,omitempty"`
	Checksums map[string]string `json:"checksums,omitempty"`
}

// GoogleInstance is used when instance type is Google Cloud SQL
// and describes necessary informations to use google API to create sql instances
type GoogleInstance struct {
	InstanceName  string         `json:"instance"`
	ConfigmapName NamespacedName `json:"configmapRef"`
	APIEndpoint   string         `json:"apiEndpoint,omitempty"`
	ClientSecret  NamespacedName `json:"clientSecretRef,omitempty"`
}

// OracleOCIInstance is used when instance type is Oracle OCI
type OracleOCIInstance struct {
	InstanceName  string         `json:"instance"`
	ConfigmapName NamespacedName `json:"configmapRef"`
	APIEndpoint   string         `json:"apiEndpoint,omitempty"`
	ClientSecret  NamespacedName `json:"clientSecretRef,omitempty"`
}

// AzureInstance is used when instance type is Azure
type AzureInstance struct {
	InstanceName  string         `json:"instance"`
	ConfigmapName NamespacedName `json:"configmapRef"`
	APIEndpoint   string         `json:"apiEndpoint,omitempty"`
	ClientSecret  NamespacedName `json:"clientSecretRef,omitempty"`
}

// AWSInstance is used when instance type is AWS
type AWSInstance struct {
	InstanceName  string         `json:"instance"`
	ConfigmapName NamespacedName `json:"configmapRef"`
	APIEndpoint   string         `json:"apiEndpoint,omitempty"`
	ClientSecret  NamespacedName `json:"clientSecretRef,omitempty"`
}

// BackendServer defines backend database server
type BackendServer struct {
	Host          string `json:"host"`
	Port          uint16 `json:"port"`
	MaxConnection uint16 `json:"maxConn"`
	ReadOnly      bool   `json:"readonly,omitempty"`
}

// GenericInstance is used when instance type is generic
// and describes necessary information to use instance
// generic instance can be any backend, it must be reachable by described address and port
type GenericInstance struct {
	Host         string   `json:"host,omitempty"`
	HostFrom     *FromRef `json:"hostFrom,omitempty"`
	Port         uint16   `json:"port,omitempty"`
	PortFrom     *FromRef `json:"portFrom,omitempty"`
	PublicIP     string   `json:"publicIp,omitempty"`
	PublicIPFrom *FromRef `json:"publicIpFrom,omitempty"`
	// BackupHost address will be used for dumping database for backup
	// Usually secondary address for primary-secondary setup or cluster lb address
	// If it's not defined, above Host will be used as backup host address.
	BackupHost string `json:"backupHost,omitempty"`
}

type FromRef struct {
	Kind      string `json:"kind"`
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Key       string `json:"key"`
}

// DbInstanceBackup defines name of google bucket to use for storing database dumps for backup when backup is enabled
type DbInstanceBackup struct {
	Bucket string `json:"bucket"`
}

// DbInstanceMonitoring defines if exporter
type DbInstanceMonitoring struct {
	Enabled bool `json:"enabled"`
}

// DbInstanceSSLConnection defines whether connection from db-operator to instance has to be ssl or not
type DbInstanceSSLConnection struct {
	Enabled bool `json:"enabled"`
	// SkipVerify use SSL connection, but don't check against a CA
	SkipVerify bool `json:"skip-verify"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster,shortName=dbin
//+kubebuilder:printcolumn:name="Phase",type=string,JSONPath=`.status.phase`,description="current phase"
//+kubebuilder:printcolumn:name="Status",type=string,JSONPath=`.status.status`,description="health status"
// +kubebuilder:storageversion

// DbInstance is the Schema for the dbinstances API
type DbInstance struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DbInstanceSpec   `json:"spec,omitempty"`
	Status DbInstanceStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// DbInstanceList contains a list of DbInstance
type DbInstanceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DbInstance `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DbInstance{}, &DbInstanceList{})
}

// ValidateEngine checks if defined engine by DbInstance object is supported by db-operator
func (dbin *DbInstance) ValidateEngine() error {
	if (dbin.Spec.Engine == "mysql") || (dbin.Spec.Engine == "postgres") ||
		(dbin.Spec.Engine == "mongodb") || (dbin.Spec.Engine == "clickhouse") ||
		(dbin.Spec.Engine == "oracle") || (dbin.Spec.Engine == "sqlserver") {
		return nil
	}

	return errors.New("not supported engine type")
}

// ValidateExistingDatabase checks if there's an existing database for the same instance in any namespace
func (dbin *DbInstance) ValidateExistingDatabase(ctx context.Context, c client.Client) error {
	var dbList DbInstanceList
	if err := c.List(ctx, &dbList); err != nil {
		return err
	}

	for _, db := range dbList.Items {
		if db.Spec.AdminUserSecret == dbin.Spec.AdminUserSecret && db.Name != dbin.Name {
			return fmt.Errorf("a database for instance %s already exists in namespace %s", dbin.Spec.AdminUserSecret, db.Namespace)
		}
	}

	return nil
}

// ValidateBackend checks if backend type of instance is defined properly
// returns error when more than one backend types are defined
// or when no backend type is defined
func (dbin *DbInstance) ValidateBackend() error {
	source := dbin.Spec.DbInstanceSource

	if source.Google == nil && source.Generic == nil &&
		source.OracleOCI == nil && source.Azure == nil &&
		source.AWS == nil {
		return errors.New("no instance type defined")
	}

	numSources := 0

	if source.Google != nil {
		numSources++
	}

	if source.Generic != nil {
		numSources++
	}

	if source.OracleOCI != nil {
		numSources++
	}

	if source.Azure != nil {
		numSources++
	}

	if source.AWS != nil {
		numSources++
	}

	if numSources > 1 {
		return errors.New("may not specify more than 1 instance type")
	}

	return nil
}

// GetBackendType returns type of instance infrastructure.
// Infrastructure where database is running ex) google cloud sql, generic instance
func (dbin *DbInstance) GetBackendType() (string, error) {
	err := dbin.ValidateBackend()
	if err != nil {
		return "", err
	}

	source := dbin.Spec.DbInstanceSource

	if source.Google != nil {
		return "google", nil
	}

	if source.OracleOCI != nil {
		return "oracleoci", nil
	}

	if source.Azure != nil {
		return "azure", nil
	}

	if source.AWS != nil {
		return "aws", nil
	}

	if source.Generic != nil {
		return "generic", nil
	}

	return "", errors.New("no backend type defined")
}

// IsMonitoringEnabled returns boolean value if monitoring is enabled for the instance
func (dbin *DbInstance) IsMonitoringEnabled() bool {
	return dbin.Spec.Monitoring.Enabled
}

// DbInstances don't have the cleanup feature
func (dbin *DbInstance) IsCleanup() bool {
	return false
}

func (dbin *DbInstance) IsDeleted() bool {
	return dbin.GetDeletionTimestamp() != nil
}

// This method isn't supported by dbin
func (dbin *DbInstance) GetSecretName() string {
	return ""
}

func (db *DbInstance) Hub() {
	// Function to mark the DbInstance as a hub
}
