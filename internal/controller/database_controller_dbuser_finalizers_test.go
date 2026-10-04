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
	"context"
	"testing"

	kindav1beta1 "github.com/db-operator/db-operator/v2/api/v1beta1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func dbUserFinalizerReconciler(t *testing.T, objs ...client.Object) *DatabaseReconciler {
	t.Helper()
	scheme := runtime.NewScheme()
	require.NoError(t, kindav1beta1.AddToScheme(scheme))
	cli := fake.NewClientBuilder().WithScheme(scheme).WithObjects(objs...).Build()
	return &DatabaseReconciler{Client: cli, Scheme: scheme}
}

func dbUserFinalizerDatabase(finalizers ...string) *kindav1beta1.Database {
	return &kindav1beta1.Database{
		ObjectMeta: metav1.ObjectMeta{Name: "otel-db", Namespace: "tenant", Finalizers: finalizers},
	}
}

func existingDbUser(name string) *kindav1beta1.DbUser {
	return &kindav1beta1.DbUser{ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: "tenant"}}
}

func storedFinalizers(t *testing.T, r *DatabaseReconciler, dbcr *kindav1beta1.Database) []string {
	t.Helper()
	stored := &kindav1beta1.Database{}
	require.NoError(t, r.Get(context.Background(), client.ObjectKeyFromObject(dbcr), stored))
	return stored.Finalizers
}

func TestUnitReleaseOrphanDbUserFinalizersKeepsExistingDbUser(t *testing.T) {
	dbcr := dbUserFinalizerDatabase("db.otel-db", "dbuser.reader", "dbuser.writer")
	r := dbUserFinalizerReconciler(t, dbcr, existingDbUser("writer"))

	referenced, err := r.releaseOrphanDbUserFinalizers(context.Background(), dbcr)

	require.NoError(t, err)
	assert.True(t, referenced)
	assert.Equal(t, []string{"db.otel-db", "dbuser.writer"}, dbcr.Finalizers)
	assert.Equal(t, []string{"db.otel-db", "dbuser.writer"}, storedFinalizers(t, r, dbcr))
}

func TestUnitReleaseOrphanDbUserFinalizersReleasesOnlyOrphans(t *testing.T) {
	dbcr := dbUserFinalizerDatabase("db.otel-db", "dbuser.reader")
	r := dbUserFinalizerReconciler(t, dbcr)

	referenced, err := r.releaseOrphanDbUserFinalizers(context.Background(), dbcr)

	require.NoError(t, err)
	assert.False(t, referenced)
	assert.Equal(t, []string{"db.otel-db"}, storedFinalizers(t, r, dbcr))
}

func TestUnitReleaseOrphanDbUserFinalizersLeavesExistingUntouched(t *testing.T) {
	dbcr := dbUserFinalizerDatabase("db.otel-db", "dbuser.writer")
	r := dbUserFinalizerReconciler(t, dbcr, existingDbUser("writer"))
	before := dbcr.ResourceVersion

	referenced, err := r.releaseOrphanDbUserFinalizers(context.Background(), dbcr)

	require.NoError(t, err)
	assert.True(t, referenced)
	assert.Equal(t, before, dbcr.ResourceVersion)
	assert.Equal(t, []string{"db.otel-db", "dbuser.writer"}, storedFinalizers(t, r, dbcr))
}
