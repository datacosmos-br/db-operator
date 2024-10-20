//go:build !ignore_autogenerated

/*
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

// Code generated by controller-gen. DO NOT EDIT.

package v1beta1

import (
	"k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AWSInstance) DeepCopyInto(out *AWSInstance) {
	*out = *in
	out.ConfigmapName = in.ConfigmapName
	out.ClientSecret = in.ClientSecret
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AWSInstance.
func (in *AWSInstance) DeepCopy() *AWSInstance {
	if in == nil {
		return nil
	}
	out := new(AWSInstance)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AzureInstance) DeepCopyInto(out *AzureInstance) {
	*out = *in
	out.ConfigmapName = in.ConfigmapName
	out.ClientSecret = in.ClientSecret
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AzureInstance.
func (in *AzureInstance) DeepCopy() *AzureInstance {
	if in == nil {
		return nil
	}
	out := new(AzureInstance)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BackendServer) DeepCopyInto(out *BackendServer) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BackendServer.
func (in *BackendServer) DeepCopy() *BackendServer {
	if in == nil {
		return nil
	}
	out := new(BackendServer)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Clickhouse) DeepCopyInto(out *Clickhouse) {
	*out = *in
	if in.Settings != nil {
		in, out := &in.Settings, &out.Settings
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Clickhouse.
func (in *Clickhouse) DeepCopy() *Clickhouse {
	if in == nil {
		return nil
	}
	out := new(Clickhouse)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Credentials) DeepCopyInto(out *Credentials) {
	*out = *in
	if in.Templates != nil {
		in, out := &in.Templates, &out.Templates
		*out = make(Templates, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(Template)
				**out = **in
			}
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Credentials.
func (in *Credentials) DeepCopy() *Credentials {
	if in == nil {
		return nil
	}
	out := new(Credentials)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Database) DeepCopyInto(out *Database) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Database.
func (in *Database) DeepCopy() *Database {
	if in == nil {
		return nil
	}
	out := new(Database)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Database) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DatabaseBackup) DeepCopyInto(out *DatabaseBackup) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DatabaseBackup.
func (in *DatabaseBackup) DeepCopy() *DatabaseBackup {
	if in == nil {
		return nil
	}
	out := new(DatabaseBackup)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DatabaseList) DeepCopyInto(out *DatabaseList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Database, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DatabaseList.
func (in *DatabaseList) DeepCopy() *DatabaseList {
	if in == nil {
		return nil
	}
	out := new(DatabaseList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *DatabaseList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DatabaseProxyStatus) DeepCopyInto(out *DatabaseProxyStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DatabaseProxyStatus.
func (in *DatabaseProxyStatus) DeepCopy() *DatabaseProxyStatus {
	if in == nil {
		return nil
	}
	out := new(DatabaseProxyStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DatabaseSpec) DeepCopyInto(out *DatabaseSpec) {
	*out = *in
	out.Backup = in.Backup
	if in.SecretsTemplates != nil {
		in, out := &in.SecretsTemplates, &out.SecretsTemplates
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	in.Postgres.DeepCopyInto(&out.Postgres)
	in.MongoDB.DeepCopyInto(&out.MongoDB)
	in.Clickhouse.DeepCopyInto(&out.Clickhouse)
	in.Oracle.DeepCopyInto(&out.Oracle)
	in.SQLServer.DeepCopyInto(&out.SQLServer)
	in.Credentials.DeepCopyInto(&out.Credentials)
	if in.AllowedNamespaces != nil {
		in, out := &in.AllowedNamespaces, &out.AllowedNamespaces
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DatabaseSpec.
func (in *DatabaseSpec) DeepCopy() *DatabaseSpec {
	if in == nil {
		return nil
	}
	out := new(DatabaseSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DatabaseStatus) DeepCopyInto(out *DatabaseStatus) {
	*out = *in
	out.ProxyStatus = in.ProxyStatus
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DatabaseStatus.
func (in *DatabaseStatus) DeepCopy() *DatabaseStatus {
	if in == nil {
		return nil
	}
	out := new(DatabaseStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DbInstance) DeepCopyInto(out *DbInstance) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DbInstance.
func (in *DbInstance) DeepCopy() *DbInstance {
	if in == nil {
		return nil
	}
	out := new(DbInstance)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *DbInstance) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DbInstanceBackup) DeepCopyInto(out *DbInstanceBackup) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DbInstanceBackup.
func (in *DbInstanceBackup) DeepCopy() *DbInstanceBackup {
	if in == nil {
		return nil
	}
	out := new(DbInstanceBackup)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DbInstanceList) DeepCopyInto(out *DbInstanceList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]DbInstance, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DbInstanceList.
func (in *DbInstanceList) DeepCopy() *DbInstanceList {
	if in == nil {
		return nil
	}
	out := new(DbInstanceList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *DbInstanceList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DbInstanceMonitoring) DeepCopyInto(out *DbInstanceMonitoring) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DbInstanceMonitoring.
func (in *DbInstanceMonitoring) DeepCopy() *DbInstanceMonitoring {
	if in == nil {
		return nil
	}
	out := new(DbInstanceMonitoring)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DbInstanceSSLConnection) DeepCopyInto(out *DbInstanceSSLConnection) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DbInstanceSSLConnection.
func (in *DbInstanceSSLConnection) DeepCopy() *DbInstanceSSLConnection {
	if in == nil {
		return nil
	}
	out := new(DbInstanceSSLConnection)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DbInstanceSource) DeepCopyInto(out *DbInstanceSource) {
	*out = *in
	if in.Google != nil {
		in, out := &in.Google, &out.Google
		*out = new(GoogleInstance)
		**out = **in
	}
	if in.Generic != nil {
		in, out := &in.Generic, &out.Generic
		*out = new(GenericInstance)
		(*in).DeepCopyInto(*out)
	}
	if in.OracleOCI != nil {
		in, out := &in.OracleOCI, &out.OracleOCI
		*out = new(OracleOCIInstance)
		**out = **in
	}
	if in.Azure != nil {
		in, out := &in.Azure, &out.Azure
		*out = new(AzureInstance)
		**out = **in
	}
	if in.AWS != nil {
		in, out := &in.AWS, &out.AWS
		*out = new(AWSInstance)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DbInstanceSource.
func (in *DbInstanceSource) DeepCopy() *DbInstanceSource {
	if in == nil {
		return nil
	}
	out := new(DbInstanceSource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DbInstanceSpec) DeepCopyInto(out *DbInstanceSpec) {
	*out = *in
	out.AdminUserSecret = in.AdminUserSecret
	out.Backup = in.Backup
	out.Monitoring = in.Monitoring
	out.SSLConnection = in.SSLConnection
	if in.AllowedPrivileges != nil {
		in, out := &in.AllowedPrivileges, &out.AllowedPrivileges
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	in.DbInstanceSource.DeepCopyInto(&out.DbInstanceSource)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DbInstanceSpec.
func (in *DbInstanceSpec) DeepCopy() *DbInstanceSpec {
	if in == nil {
		return nil
	}
	out := new(DbInstanceSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DbInstanceStatus) DeepCopyInto(out *DbInstanceStatus) {
	*out = *in
	if in.Info != nil {
		in, out := &in.Info, &out.Info
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Checksums != nil {
		in, out := &in.Checksums, &out.Checksums
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DbInstanceStatus.
func (in *DbInstanceStatus) DeepCopy() *DbInstanceStatus {
	if in == nil {
		return nil
	}
	out := new(DbInstanceStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DbUser) DeepCopyInto(out *DbUser) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DbUser.
func (in *DbUser) DeepCopy() *DbUser {
	if in == nil {
		return nil
	}
	out := new(DbUser)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *DbUser) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DbUserList) DeepCopyInto(out *DbUserList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]DbUser, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DbUserList.
func (in *DbUserList) DeepCopy() *DbUserList {
	if in == nil {
		return nil
	}
	out := new(DbUserList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *DbUserList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DbUserSpec) DeepCopyInto(out *DbUserSpec) {
	*out = *in
	if in.DatabaseRefs != nil {
		in, out := &in.DatabaseRefs, &out.DatabaseRefs
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.ExtraPrivileges != nil {
		in, out := &in.ExtraPrivileges, &out.ExtraPrivileges
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	in.Credentials.DeepCopyInto(&out.Credentials)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DbUserSpec.
func (in *DbUserSpec) DeepCopy() *DbUserSpec {
	if in == nil {
		return nil
	}
	out := new(DbUserSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DbUserStatus) DeepCopyInto(out *DbUserStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DbUserStatus.
func (in *DbUserStatus) DeepCopy() *DbUserStatus {
	if in == nil {
		return nil
	}
	out := new(DbUserStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FromRef) DeepCopyInto(out *FromRef) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FromRef.
func (in *FromRef) DeepCopy() *FromRef {
	if in == nil {
		return nil
	}
	out := new(FromRef)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GenericInstance) DeepCopyInto(out *GenericInstance) {
	*out = *in
	if in.HostFrom != nil {
		in, out := &in.HostFrom, &out.HostFrom
		*out = new(FromRef)
		**out = **in
	}
	if in.PortFrom != nil {
		in, out := &in.PortFrom, &out.PortFrom
		*out = new(FromRef)
		**out = **in
	}
	if in.PublicIPFrom != nil {
		in, out := &in.PublicIPFrom, &out.PublicIPFrom
		*out = new(FromRef)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GenericInstance.
func (in *GenericInstance) DeepCopy() *GenericInstance {
	if in == nil {
		return nil
	}
	out := new(GenericInstance)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GoogleInstance) DeepCopyInto(out *GoogleInstance) {
	*out = *in
	out.ConfigmapName = in.ConfigmapName
	out.ClientSecret = in.ClientSecret
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GoogleInstance.
func (in *GoogleInstance) DeepCopy() *GoogleInstance {
	if in == nil {
		return nil
	}
	out := new(GoogleInstance)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MongoDB) DeepCopyInto(out *MongoDB) {
	*out = *in
	if in.Collections != nil {
		in, out := &in.Collections, &out.Collections
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MongoDB.
func (in *MongoDB) DeepCopy() *MongoDB {
	if in == nil {
		return nil
	}
	out := new(MongoDB)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NamespacedName) DeepCopyInto(out *NamespacedName) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NamespacedName.
func (in *NamespacedName) DeepCopy() *NamespacedName {
	if in == nil {
		return nil
	}
	out := new(NamespacedName)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Oracle) DeepCopyInto(out *Oracle) {
	*out = *in
	if in.Tablespaces != nil {
		in, out := &in.Tablespaces, &out.Tablespaces
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Profiles != nil {
		in, out := &in.Profiles, &out.Profiles
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Oracle.
func (in *Oracle) DeepCopy() *Oracle {
	if in == nil {
		return nil
	}
	out := new(Oracle)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *OracleOCIInstance) DeepCopyInto(out *OracleOCIInstance) {
	*out = *in
	out.ConfigmapName = in.ConfigmapName
	out.ClientSecret = in.ClientSecret
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new OracleOCIInstance.
func (in *OracleOCIInstance) DeepCopy() *OracleOCIInstance {
	if in == nil {
		return nil
	}
	out := new(OracleOCIInstance)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Postgres) DeepCopyInto(out *Postgres) {
	*out = *in
	if in.Extensions != nil {
		in, out := &in.Extensions, &out.Extensions
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Schemas != nil {
		in, out := &in.Schemas, &out.Schemas
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Postgres.
func (in *Postgres) DeepCopy() *Postgres {
	if in == nil {
		return nil
	}
	out := new(Postgres)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SQLServer) DeepCopyInto(out *SQLServer) {
	*out = *in
	if in.Schemas != nil {
		in, out := &in.Schemas, &out.Schemas
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Roles != nil {
		in, out := &in.Roles, &out.Roles
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SQLServer.
func (in *SQLServer) DeepCopy() *SQLServer {
	if in == nil {
		return nil
	}
	out := new(SQLServer)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Template) DeepCopyInto(out *Template) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Template.
func (in *Template) DeepCopy() *Template {
	if in == nil {
		return nil
	}
	out := new(Template)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in Templates) DeepCopyInto(out *Templates) {
	{
		in := &in
		*out = make(Templates, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(Template)
				**out = **in
			}
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Templates.
func (in Templates) DeepCopy() Templates {
	if in == nil {
		return nil
	}
	out := new(Templates)
	in.DeepCopyInto(out)
	return *out
}
