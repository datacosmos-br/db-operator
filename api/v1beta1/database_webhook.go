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
	"regexp"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/strings/slices"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// log is for logging in this package.
var databaselog = logf.Log.WithName("database-resource")

var databaseMgr ctrl.Manager

func (r *Database) SetupWebhookWithManager(mgr ctrl.Manager) error {
	databaseMgr = mgr
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-kinda-rocks-v1beta1-database,mutating=true,failurePolicy=fail,sideEffects=None,groups=kinda.rocks,resources=databases,verbs=create;update,versions=v1beta1,name=mdatabase.kb.io,admissionReviewVersions=v1

const (
	DEFAULT_TEMPLATE_VALUE = "{{ .Protocol }}://{{ .Username }}:{{ .Password }}@{{ .Hostname }}:{{ .Port }}/{{ .Database }}"
	DEFAULT_TEMPLATE_NAME  = "CONNECTION_STRING"
)

var _ webhook.CustomDefaulter = &Database{}

// Default implements webhook.CustomDefaulter so a webhook will be registered for the type
func (r *Database) Default(ctx context.Context, obj runtime.Object) error {
	databaselog.Info("default", "name", r.Name)
	if len(r.Spec.SecretsTemplates) == 0 && len(r.Spec.Credentials.Templates) == 0 {
		r.Spec.Credentials = Credentials{
			Templates: Templates{
				&Template{
					Name:     DEFAULT_TEMPLATE_NAME,
					Template: DEFAULT_TEMPLATE_VALUE,
					Secret:   true,
				},
			},
		}
	}
	return nil
}

//+kubebuilder:webhook:path=/validate-kinda-rocks-v1beta1-database,mutating=false,failurePolicy=fail,sideEffects=None,groups=kinda.rocks,resources=databases,verbs=create;update,versions=v1beta1,name=vdatabase.kb.io,admissionReviewVersions=v1

var _ webhook.CustomValidator = &Database{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *Database) ValidateCreate(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	databaselog.Info("validate create", "name", r.Name)

	if r.Spec.SecretsTemplates != nil && r.Spec.Credentials.Templates != nil {
		return nil, errors.New("using both: secretsTemplates and templates, is not allowed")
	}

	if r.Spec.SecretsTemplates != nil {
		databaselog.Info("secretsTemplates are deprecated, it will be removed in the next API version. Please, consider switching to templates")
		if err := ValidateSecretTemplates(r.Spec.SecretsTemplates); err != nil {
			return nil, err
		}
	}

	if r.Spec.Credentials.Templates != nil {
		if err := ValidateTemplates(r.Spec.Credentials.Templates, true); err != nil {
			return nil, err
		}
	}

	if err := r.ValidateNamespace(); err != nil {
		return nil, err
	}
	if err := r.ValidateExistingDatabase(context.Background(), databaseMgr.GetClient()); err != nil {
		return nil, err
	}
	return nil, nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *Database) ValidateUpdate(ctx context.Context, obj runtime.Object, old runtime.Object) (admission.Warnings, error) {
	databaselog.Info("validate update", "name", r.Name)

	if r.Spec.SecretsTemplates != nil && r.Spec.Credentials.Templates != nil {
		return nil, errors.New("using both: secretsTemplates and templates, is not allowed")
	}

	if r.Spec.SecretsTemplates != nil {
		databaselog.Info("secretsTemplates are deprecated, it will be removed in the next API version. Please, consider switching to templates")
		err := ValidateSecretTemplates(r.Spec.SecretsTemplates)
		if err != nil {
			return nil, err
		}
	}

	if r.Spec.Credentials.Templates != nil {
		if err := ValidateTemplates(r.Spec.Credentials.Templates, true); err != nil {
			return nil, err
		}
	}

	// Ensure fields are immutable
	immutableErr := "cannot change %s, the field is immutable"
	oldDatabase, ok := old.(*Database)
	if !ok {
		return nil, fmt.Errorf("failed to cast old object to Database")
	}

	if r.Spec.Instance != oldDatabase.Spec.Instance {
		return nil, fmt.Errorf(immutableErr, "spec.instance")
	}

	if r.Spec.Postgres.Template != oldDatabase.Spec.Postgres.Template {
		return nil, fmt.Errorf(immutableErr, "spec.postgres.template")
	}

	if err := r.ValidateNamespace(); err != nil {
		return nil, err
	}
	if err := r.ValidateExistingDatabase(context.Background(), databaseMgr.GetClient()); err != nil {
		return nil, err
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

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *Database) ValidateDelete(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	databaselog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil, nil
}
