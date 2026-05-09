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

//+kubebuilder:webhook:path=/validate-kinda-rocks-v1beta1-dbuser,mutating=false,failurePolicy=fail,sideEffects=None,groups=kinda.rocks,resources=dbusers,verbs=create;update,versions=v1beta1,name=vdbuser.kb.io,admissionReviewVersions=v1

// DbUserCustomValidator implements admission.Validator[*DbUser].
type DbUserCustomValidator struct{}

var _ admission.Validator[*DbUser] = &DbUserCustomValidator{}

// ValidateCreate validates the DbUser on creation.
func (v *DbUserCustomValidator) ValidateCreate(ctx context.Context, obj *DbUser) (admission.Warnings, error) {
	warnings := []string{}
	if err := TestExtraPrivileges(obj.Spec.ExtraPrivileges); err != nil {
		return nil, err
	}
	if len(obj.Spec.ExtraPrivileges) > 0 {
		warnings = append(warnings,
			"extra privileges is an experimental feature, please use at your own risk and feel free to open GitHub issues.")
	}

	dbuserlog.Info("validate create", "name", obj.Name)
	if err := IsAccessTypeSupported(obj.Spec.AccessType); err != nil {
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
func (v *DbUserCustomValidator) ValidateUpdate(ctx context.Context, oldObj, newObj *DbUser) (admission.Warnings, error) {
	warnings := []string{}
	dbuserlog.Info("validate update", "name", newObj.Name)

	if len(newObj.Spec.ExtraPrivileges) > 0 {
		warnings = append(warnings,
			"extra privileges is an experimental feature, please use at your own risk and feel free to open GitHub issues.")
	}
	if err := TestExtraPrivileges(newObj.Spec.ExtraPrivileges); err != nil {
		return nil, err
	}
	if err := IsAccessTypeSupported(newObj.Spec.AccessType); err != nil {
		return nil, err
	}
	if newObj.Spec.Credentials.Templates != nil {
		if err := ValidateTemplates(newObj.Spec.Credentials.Templates); err != nil {
			return nil, err
		}
	}
	if oldObj.Spec.Postgres.GrantToAdmin != newObj.Spec.Postgres.GrantToAdmin {
		return nil, errors.New("grantToAdmin is an immutable field")
	}
	for _, role := range oldObj.Spec.ExtraPrivileges {
		if !slices.Contains(newObj.Spec.ExtraPrivileges, role) {
			warnings = append(
				warnings,
				fmt.Sprintf("extra privileges can't be removed by the operator, please manually revoke %s from the user %s",
					role, newObj.Name),
			)
		}
	}

	return warnings, nil
}

// ValidateDelete validates the DbUser on deletion.
func (v *DbUserCustomValidator) ValidateDelete(ctx context.Context, obj *DbUser) (admission.Warnings, error) {
	dbuserlog.Info("validate delete", "name", obj.Name)
	return nil, nil
}
