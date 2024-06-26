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
	"k8s.io/apimachinery/pkg/types"
)

// NamespacedName is a fork of the kubernetes api type of the same name.
// Sadly this is required because CRD structs must have all fields json tagged and the kubernetes type is not tagged.
type NamespacedName struct {
	Namespace string `json:"Namespace"`
	Name      string `json:"Name"`
}

// ToKubernetesType converts our local type to the kubernetes API equivalent.
func (nn *NamespacedName) ToKubernetesType() types.NamespacedName {
	if nn == nil {
		return types.NamespacedName{}
	}

	return types.NamespacedName{
		Name:      nn.Name,
		Namespace: nn.Namespace,
	}
}

// Tempaltes to add custom entries to ConfigMaps and Secrets
type Template struct {
	Name     string `json:"name"`
	Template string `json:"template"`
	Secret   bool   `json:"secret"`
}

type Templates []*Template

// Credentials should be used to setup everything relates to k8s secrets and configmaps
// TODO(@allanger): Field .spec.secretName should be moved here in the v1beta2 version
type Credentials struct {
	// Templates to add custom entries to ConfigMaps and Secrets
	Templates Templates `json:"templates,omitempty"`
}
