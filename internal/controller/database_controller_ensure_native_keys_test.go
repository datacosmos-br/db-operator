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

package controller

import (
	"testing"

	kindav1beta1 "github.com/db-operator/db-operator/v2/api/v1beta1"
	"github.com/db-operator/db-operator/v2/pkg/consts"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestEnsureDatabaseSecretNativeKeysRejectsTemplatedAliases(t *testing.T) {
	dbcr := &kindav1beta1.Database{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "cosmos",
			Name:      "cosmos-db",
		},
		Spec: kindav1beta1.DatabaseSpec{
			DatabaseName: "cosmos",
		},
		Status: kindav1beta1.DatabaseStatus{
			Engine: consts.ENGINE_POSTGRES,
		},
	}

	secret := &corev1.Secret{
		Data: map[string][]byte{
			"username": []byte("{{ .Username }}"),
			"password": []byte("{{ .Password }}"),
			"dbname":   []byte("{{ .Database }}"),
		},
	}

	r := &DatabaseReconciler{}
	if err := r.ensureDatabaseSecretNativeKeys(nil, dbcr, secret); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if string(secret.Data[consts.POSTGRES_USER]) == "{{ .Username }}" {
		t.Errorf("POSTGRES_USER must not be set from an unrendered template alias")
	}
	if string(secret.Data[consts.POSTGRES_PASSWORD]) == "{{ .Password }}" {
		t.Errorf("POSTGRES_PASSWORD must not be set from an unrendered template alias")
	}
	if string(secret.Data[consts.POSTGRES_DB]) == "{{ .Database }}" {
		t.Errorf("POSTGRES_DB must not be set from an unrendered template alias")
	}

	expectedUser := "cosmos-cosmos-db"
	if string(secret.Data[consts.POSTGRES_USER]) != expectedUser {
		t.Errorf("POSTGRES_USER = %q, want %q", secret.Data[consts.POSTGRES_USER], expectedUser)
	}
	if len(secret.Data[consts.POSTGRES_PASSWORD]) == 0 {
		t.Errorf("POSTGRES_PASSWORD should have been generated")
	}
	if string(secret.Data[consts.POSTGRES_DB]) != "cosmos" {
		t.Errorf("POSTGRES_DB = %q, want %q", secret.Data[consts.POSTGRES_DB], "cosmos")
	}
}

func TestEnsureDatabaseSecretNativeKeysPreservesRenderedAliases(t *testing.T) {
	dbcr := &kindav1beta1.Database{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "app",
			Name:      "db",
		},
		Spec: kindav1beta1.DatabaseSpec{
			DatabaseName: "app-db",
		},
		Status: kindav1beta1.DatabaseStatus{
			Engine: consts.ENGINE_POSTGRES,
		},
	}

	secret := &corev1.Secret{
		Data: map[string][]byte{
			"username": []byte("app-db-user"),
			"password": []byte("app-db-pass"),
		},
	}

	r := &DatabaseReconciler{}
	if err := r.ensureDatabaseSecretNativeKeys(nil, dbcr, secret); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if string(secret.Data[consts.POSTGRES_USER]) != "app-db-user" {
		t.Errorf("POSTGRES_USER = %q, want %q", secret.Data[consts.POSTGRES_USER], "app-db-user")
	}
	if string(secret.Data[consts.POSTGRES_PASSWORD]) != "app-db-pass" {
		t.Errorf("POSTGRES_PASSWORD = %q, want %q", secret.Data[consts.POSTGRES_PASSWORD], "app-db-pass")
	}
}

func TestEnsureDatabaseSecretNativeKeysNeverWritesEmptyDatabaseName(t *testing.T) {
	cases := []struct {
		name       string
		engine     string
		key        string
		specName   string
		statusName string
		data       map[string][]byte
		want       string
	}{
		{"postgres missing key defaults", consts.ENGINE_POSTGRES, consts.POSTGRES_DB, "", "", map[string][]byte{}, "keycloak-db"},
		{"postgres empty key heals from status", consts.ENGINE_POSTGRES, consts.POSTGRES_DB, "", "adopted", map[string][]byte{consts.POSTGRES_DB: {}}, "adopted"},
		{"status wins over spec", consts.ENGINE_POSTGRES, consts.POSTGRES_DB, "override", "adopted", map[string][]byte{consts.POSTGRES_DB: {}}, "adopted"},
		{"spec used before status is recorded", consts.ENGINE_POSTGRES, consts.POSTGRES_DB, "override", "", map[string][]byte{consts.POSTGRES_DB: {}}, "override"},
		{"mysql default normalized", consts.ENGINE_MYSQL, consts.MYSQL_DB, "", "", map[string][]byte{consts.MYSQL_DB: {}}, "keycloak_db"},
		{"mysql spec normalized", consts.ENGINE_MYSQL, consts.MYSQL_DB, "My-App", "", map[string][]byte{consts.MYSQL_DB: {}}, "my_app"},
		{"clickhouse empty key defaults", consts.ENGINE_CLICKHOUSE, consts.CLICKHOUSE_DB, "", "", map[string][]byte{consts.CLICKHOUSE_DB: {}}, "keycloak-db"},
		{"existing value preserved", consts.ENGINE_POSTGRES, consts.POSTGRES_DB, "", "adopted", map[string][]byte{consts.POSTGRES_DB: []byte("kept")}, "kept"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dbcr := &kindav1beta1.Database{
				ObjectMeta: metav1.ObjectMeta{Namespace: "keycloak", Name: "db"},
				Spec:       kindav1beta1.DatabaseSpec{DatabaseName: tc.specName},
				Status:     kindav1beta1.DatabaseStatus{Engine: tc.engine, DatabaseName: tc.statusName},
			}
			secret := &corev1.Secret{Data: tc.data}

			r := &DatabaseReconciler{}
			if err := r.ensureDatabaseSecretNativeKeys(nil, dbcr, secret); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got := string(secret.Data[tc.key]); got != tc.want {
				t.Errorf("%s = %q, want %q", tc.key, got, tc.want)
			}
		})
	}
}
