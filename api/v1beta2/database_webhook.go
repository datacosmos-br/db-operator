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

package v1beta2

import (
	"context"
	"fmt"
	"regexp"
	"slices"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// log is for logging in this package.
var databaselog = logf.Log.WithName("database-resource")

func (r *Database) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr, r).
		WithDefaulter(&DatabaseCustomDefaulter{}).
		WithValidator(&DatabaseCustomValidator{}).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-kinda-rocks-v1beta2-database,mutating=true,failurePolicy=fail,sideEffects=None,groups=kinda.rocks,resources=databases,verbs=create;update,versions=v1beta2,name=mdatabase-v1beta2.kb.io,admissionReviewVersions=v1

const (
	DEFAULT_TEMPLATE_VALUE = "{{ .Protocol }}://{{ .Username }}:{{ .Password }}@{{ .Hostname }}:{{ .Port }}/{{ .Database }}"
	DEFAULT_TEMPLATE_NAME  = "CONNECTION_STRING"
)

// DatabaseCustomDefaulter handles defaulting for Database resources.
type DatabaseCustomDefaulter struct{}

func (d *DatabaseCustomDefaulter) Default(_ context.Context, obj *Database) error {
	obj.Default()
	return nil
}

// Default sets default values on the Database resource.
func (r *Database) Default() {
	databaselog.Info("default", "name", r.Name)
	if len(r.Spec.Credentials.Templates) == 0 {
		r.Spec.Credentials.Templates = Templates{
			&Template{
				Name:     DEFAULT_TEMPLATE_NAME,
				Template: DEFAULT_TEMPLATE_VALUE,
			},
		}
	}
}

//+kubebuilder:webhook:path=/validate-kinda-rocks-v1beta2-database,mutating=false,failurePolicy=fail,sideEffects=None,groups=kinda.rocks,resources=databases,verbs=create;update,versions=v1beta2,name=vdatabase-v1beta2.kb.io,admissionReviewVersions=v1

// DatabaseCustomValidator handles validation for Database resources.
type DatabaseCustomValidator struct{}

func (v *DatabaseCustomValidator) ValidateCreate(_ context.Context, obj *Database) (admission.Warnings, error) {
	return obj.ValidateCreate()
}

func (v *DatabaseCustomValidator) ValidateUpdate(_ context.Context, newObj *Database, oldObj *Database) (admission.Warnings, error) {
	return newObj.ValidateUpdate(oldObj)
}

func (v *DatabaseCustomValidator) ValidateDelete(_ context.Context, obj *Database) (admission.Warnings, error) {
	return obj.ValidateDelete()
}

// ValidateCreate validates the Database on creation.
func (r *Database) ValidateCreate() (admission.Warnings, error) {
	databaselog.Info("validate create", "name", r.Name)

	if r.Spec.Credentials.Templates != nil {
		if err := ValidateTemplates(r.Spec.Credentials.Templates); err != nil {
			return nil, err
		}
	}

	return nil, nil
}

// ValidateUpdate validates the Database on update.
func (r *Database) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	databaselog.Info("validate update", "name", r.Name)

	if r.Spec.Credentials.Templates != nil {
		if err := ValidateTemplates(r.Spec.Credentials.Templates); err != nil {
			return nil, err
		}
	}

	// Ensure fields are immutable
	immutableErr := "cannot change %s, the field is immutable"
	oldDatabase, _ := old.(*Database)
	if r.Spec.Instance != oldDatabase.Spec.Instance {
		return nil, fmt.Errorf(immutableErr, "spec.instance")
	}

	if r.Spec.Postgres.Params.Template != oldDatabase.Spec.Postgres.Params.Template {
		return nil, fmt.Errorf(immutableErr, "spec.postgres.template")
	}

	return nil, nil
}

func ValidateSecretTemplates(templates map[string]string) error {
	for _, template := range templates {
		allowedFields := []string{".Protocol", ".DatabaseHost", ".DatabasePort", ".UserName", ".Password", ".DatabaseName"}
		// This regexp is getting fields from mustache templates so then they can be compared to allowed fields
		reg := "{{\\s*([\\w\\.]+)\\s*}}"
		r, _ := regexp.Compile(reg)
		fields := r.FindAllStringSubmatch(template, -1)
		for _, field := range fields {
			if !slices.Contains(allowedFields, field[1]) {
				err := fmt.Errorf("%v is a field that is not allowed for templating, please use one of these: %v", field[1], allowedFields)
				return err
			}
		}
	}
	return nil
}

// ValidateDelete validates the Database on deletion.
func (r *Database) ValidateDelete() (admission.Warnings, error) {
	databaselog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil, nil
}
