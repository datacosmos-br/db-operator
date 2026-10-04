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

package kube_test

import (
	"context"
	"testing"

	kindav1beta1 "github.com/db-operator/db-operator/v2/api/v1beta1"
	"github.com/db-operator/db-operator/v2/internal/helpers/kube"
	"github.com/db-operator/db-operator/v2/pkg/consts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/events"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func foreignSecret(name string) *corev1.Secret {
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: "default",
			Labels: map[string]string{
				consts.USED_BY_KIND_LABEL_KEY: "DbUser",
				consts.USED_BY_NAME_LABEL_KEY: "someone-else",
			},
		},
	}
}

func fakeHelper(t *testing.T, caller *kindav1beta1.Database, objs ...client.Object) (*kube.KubeHelper, client.Client) {
	t.Helper()
	scheme := runtime.NewScheme()
	require.NoError(t, corev1.AddToScheme(scheme))
	cli := fake.NewClientBuilder().WithScheme(scheme).WithObjects(objs...).Build()
	return kube.NewKubeHelper(cli, events.NewFakeRecorder(8), caller), cli
}

func TestUnitHandleDeleteSkipsUnusedSecret(t *testing.T) {
	t.Parallel()
	caller := database.DeepCopy()
	secret := foreignSecret("funmotors-secret")
	kh, cli := fakeHelper(t, caller, secret)

	err := kh.HandleDelete(context.Background(), secret.DeepCopy())
	require.NoError(t, err)

	got := &corev1.Secret{}
	require.NoError(t, cli.Get(context.Background(), types.NamespacedName{Namespace: "default", Name: "funmotors-secret"}, got))
	assert.Equal(t, "someone-else", got.Labels[consts.USED_BY_NAME_LABEL_KEY])
	assert.Equal(t, "DbUser", got.Labels[consts.USED_BY_KIND_LABEL_KEY])
}

func TestUnitHandleDeleteStillReleasesOwnedSecret(t *testing.T) {
	t.Parallel()
	caller := database.DeepCopy()
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "owned-secret",
			Namespace: "default",
			Labels: map[string]string{
				consts.USED_BY_KIND_LABEL_KEY: caller.Kind,
				consts.USED_BY_NAME_LABEL_KEY: caller.Name,
			},
		},
	}
	kh, cli := fakeHelper(t, caller, secret)

	require.NoError(t, kh.HandleDelete(context.Background(), secret.DeepCopy()))

	got := &corev1.Secret{}
	require.NoError(t, cli.Get(context.Background(), types.NamespacedName{Namespace: "default", Name: "owned-secret"}, got))
	assert.NotContains(t, got.Labels, consts.USED_BY_KIND_LABEL_KEY)
	assert.NotContains(t, got.Labels, consts.USED_BY_NAME_LABEL_KEY)
}

func TestUnitUpdateDoesNotSkipForeignSecret(t *testing.T) {
	t.Parallel()
	caller := database.DeepCopy()
	secret := foreignSecret("funmotors-secret-update")
	kh, cli := fakeHelper(t, caller, secret)

	current := &corev1.Secret{}
	require.NoError(t, cli.Get(context.Background(), types.NamespacedName{Namespace: "default", Name: secret.Name}, current))
	err := kh.ModifyObject(context.Background(), current)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "editing is not possible")

	got := &corev1.Secret{}
	require.NoError(t, cli.Get(context.Background(), types.NamespacedName{Namespace: "default", Name: secret.Name}, got))
	assert.Equal(t, "someone-else", got.Labels[consts.USED_BY_NAME_LABEL_KEY])
}

func TestUnitHandleDeleteStillRefusesForeignConfigMap(t *testing.T) {
	t.Parallel()
	caller := database.DeepCopy()
	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "funmotors-config",
			Namespace: "default",
			Labels: map[string]string{
				consts.USED_BY_KIND_LABEL_KEY: "DbUser",
				consts.USED_BY_NAME_LABEL_KEY: "someone-else",
			},
		},
	}
	kh, cli := fakeHelper(t, caller, cm)

	err := kh.HandleDelete(context.Background(), cm.DeepCopy())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "editing is not possible")

	got := &corev1.ConfigMap{}
	require.NoError(t, cli.Get(context.Background(), types.NamespacedName{Namespace: "default", Name: cm.Name}, got))
	assert.Equal(t, "someone-else", got.Labels[consts.USED_BY_NAME_LABEL_KEY])
}
