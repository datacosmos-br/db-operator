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
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// DbUserSpec defines the desired state of DbUser
type DbUserSpec struct {
	// DatabaseRef should contain a name of a Database to create a user there
	// Database should be in the same namespace with the user, unless NamespaceRef is specified
	DatabaseRef  string   `json:"databaseRef,omitempty"`
	DatabaseRefs []string `json:"databaseRefs,omitempty"`
	NamespaceRef string   `json:"namespaceRef,omitempty"`
	UserName     string   `json:"user,omitempty"`
	// AccessType that should be given to a user
	// Currently only readOnly and readWrite are supported by the operator
	AccessType string `json:"accessType"`
	// SecretName name that should be used to save user's credentials
	SecretName string `json:"secretName"`
	// A list of additional roles that should be added to the user
	ExtraPrivileges []string    `json:"extraPrivileges,omitempty"`
	Credentials     Credentials `json:"credentials,omitempty"`
	Cleanup         bool        `json:"cleanup,omitempty"`
	// Should the user be granted to the admin user
	// For example, it should be set to true on Azure instance,
	// because the admin given by them is not a super user,
	// but should be set to false on AWS, when rds_iam extra
	// privilege is added
	// By default is set to true
	// Only applies to Postgres, doesn't have any effect on Mysql
	// TODO: Default should be false, but not to introduce breaking
	//       changes it's now set to true. It should be changed in
	//       in the next API version
	// +kubebuilder:default=true
	// +optional
	GrantToAdmin bool `json:"grantToAdmin"`
}

// DbUserStatus defines the observed state of DbUser
type DbUserStatus struct {
	Status       bool   `json:"status"`
	DatabaseName string `json:"database"`
	// It's required to let the operator update users
	Created bool `json:"created"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="Status",type=boolean,JSONPath=`.status.status`,description="current dbuser status"
//+kubebuilder:printcolumn:name="DatabaseName",type=string,JSONPath=`.spec.databaseRef`,description="To which database user should have access"
//+kubebuilder:printcolumn:name="DatabaseNames",type=string,JSONPath=`.spec.databaseRefs`,description="To which databases user should have access"
//+kubebuilder:printcolumn:name="AccessType",type=string,JSONPath=`.spec.accessType`,description="A type of access the user has"
//+kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`,description="time since creation of resource"

// DbUser is the Schema for the dbusers API
type DbUser struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DbUserSpec   `json:"spec,omitempty"`
	Status DbUserStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// DbUserList contains a list of DbUser
type DbUserList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DbUser `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DbUser{}, &DbUserList{})
}

// Access types that are supported by the operator
const (
	READONLY  = "readOnly"
	READWRITE = "readWrite"
)

// IsAccessTypeSupported returns an error if access type is not supported
func IsAccessTypeSupported(wantedAccessType string) error {
	supportedAccessTypes := []string{READONLY, READWRITE}
	for _, supportedAccessType := range supportedAccessTypes {
		if supportedAccessType == wantedAccessType {
			return nil
		}
	}
	return fmt.Errorf("the provided access type is not supported by the operator: %s - please choose one of these: %v",
		wantedAccessType,
		supportedAccessTypes,
	)
}

// GetDatabase retrieves the referenced Database and validates the namespace
func (dbu *DbUser) GetDatabase(ctx context.Context, c client.Client) (*Database, error) {
	db := &Database{}
	namespace := dbu.Namespace
	if dbu.Spec.NamespaceRef != "" {
		namespace = dbu.Spec.NamespaceRef
	}

	if dbu.Spec.DatabaseRef != "" && len(dbu.Spec.DatabaseRefs) > 0 {
		return nil, fmt.Errorf("cannot specify both databaseRef and databaseRefs")
	}

	if dbu.Spec.DatabaseRef != "" {
		if err := c.Get(ctx, client.ObjectKey{Name: dbu.Spec.DatabaseRef, Namespace: namespace}, db); err != nil {
			return nil, fmt.Errorf("unable to fetch database: %v", err)
		}
		if err := db.ValidateDbUserNamespace(dbu.Namespace); err != nil {
			return nil, err
		}
		return db, nil
	}

	for _, dbRef := range dbu.Spec.DatabaseRefs {
		if err := c.Get(ctx, client.ObjectKey{Name: dbRef, Namespace: namespace}, db); err != nil {
			return nil, fmt.Errorf("unable to fetch database: %v", err)
		}
		if err := db.ValidateDbUserNamespace(dbu.Namespace); err != nil {
			return nil, err
		}
	}
	return db, nil
}

// ValidateExistingUser checks if the user already exists for the same instance
func (dbu *DbUser) ValidateExistingUser(ctx context.Context, c client.Client) error {
	var dbUserList DbUserList
	if err := c.List(ctx, &dbUserList); err != nil {
		return err
	}

	for _, existingDbUser := range dbUserList.Items {
		if existingDbUser.Spec.UserName == dbu.Spec.UserName && existingDbUser.Name != dbu.Name {
			return fmt.Errorf("a user with name %s already exists in namespace %s", dbu.Spec.UserName, existingDbUser.Namespace)
		}
	}

	return nil
}

// DbUsers don't have cleanup feature implemented
func (dbu *DbUser) IsCleanup() bool {
	return dbu.Spec.Cleanup
}

func (dbu *DbUser) IsDeleted() bool {
	return dbu.GetDeletionTimestamp() != nil
}

func (dbu *DbUser) GetSecretName() string {
	return dbu.Spec.SecretName
}
