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

package v1beta2

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/db-operator/db-operator/v2/pkg/consts"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// log is for logging in this package.
var dbuserlog = logf.Log.WithName("dbuser-resource")

func (r *DbUser) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr, r).
		WithValidator(&DbUserCustomValidator{}).
		Complete()
}

//+kubebuilder:webhook:path=/validate-kinda-rocks-v1beta2-dbuser,mutating=false,failurePolicy=fail,sideEffects=None,groups=kinda.rocks,resources=dbusers,verbs=create;update,versions=v1beta2,name=vdbuser-v1beta2.kb.io,admissionReviewVersions=v1

// DbUserCustomValidator handles validation for DbUser resources.
type DbUserCustomValidator struct{}

func (v *DbUserCustomValidator) ValidateCreate(_ context.Context, obj *DbUser) (admission.Warnings, error) {
	return obj.ValidateCreate()
}

func (v *DbUserCustomValidator) ValidateUpdate(_ context.Context, newObj *DbUser, oldObj *DbUser) (admission.Warnings, error) {
	return newObj.ValidateUpdate(oldObj)
}

func (v *DbUserCustomValidator) ValidateDelete(_ context.Context, obj *DbUser) (admission.Warnings, error) {
	return obj.ValidateDelete()
}

// ValidateCreate validates the DbUser on creation.
func (r *DbUser) ValidateCreate() (admission.Warnings, error) {
	warnings := []string{}
	if err := TestExtraPrivileges(r.Spec.ExtraPrivileges); err != nil {
		return nil, err
	}
	if len(r.Spec.ExtraPrivileges) > 0 {
		warnings = append(warnings,
			"extra privileges is an experimental feature, please use at your own risk and feel free to open GitHub issues.")
	}

	dbuserlog.Info("validate create", "name", r.Name)
	if err := IsAccessTypeSupported(r.Spec.AccessType); err != nil {
		return nil, err
	}

	return warnings, nil
}

func TestExtraPrivileges(privileges []string) error {
	for _, privilege := range privileges {
		if strings.ToUpper(privilege) == consts.ALL_PRIVILEGES {
			return errors.New("it's not allowed to grant ALL PRIVILEGES")
		}
	}
	return nil
}

// ValidateUpdate validates the DbUser on update.
func (r *DbUser) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	warnings := []string{}
	dbuserlog.Info("validate update", "name", r.Name)

	if len(r.Spec.ExtraPrivileges) > 0 {
		warnings = append(warnings,
			"extra privileges is an experimental feature, please use at your own risk and feel free to open GitHub issues.")
	}
	if err := TestExtraPrivileges(r.Spec.ExtraPrivileges); err != nil {
		return nil, err
	}
	if err := IsAccessTypeSupported(r.Spec.AccessType); err != nil {
		return nil, err
	}
	_, ok := old.(*DbUser)
	if !ok {
		return nil, fmt.Errorf("couldn't get the previous version of %s", r.Name)
	}
	if r.Spec.Credentials.Templates != nil {
		if err := ValidateTemplates(r.Spec.Credentials.Templates); err != nil {
			return nil, err
		}
	}
	if old.(*DbUser).Spec.Postgres.GrantToAdmin != r.Spec.Postgres.GrantToAdmin {
		return nil, errors.New("grantToAdmin is an immutable field")
	}
	for _, role := range old.(*DbUser).Spec.ExtraPrivileges {
		if !slices.Contains(r.Spec.ExtraPrivileges, role) {
			warnings = append(
				warnings,
				fmt.Sprintf("extra privileges can't be removed by the operator, please manually revoke %s from the user %s",
					role, r.Name),
			)
		}
	}

	return warnings, nil
}

// ValidateDelete validates the DbUser on deletion.
func (r *DbUser) ValidateDelete() (admission.Warnings, error) {
	dbuserlog.Info("validate delete", "name", r.Name)
	return nil, nil
}
