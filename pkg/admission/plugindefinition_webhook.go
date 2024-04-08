// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and Greenhouse contributors
// SPDX-License-Identifier: Apache-2.0

package admission

import (
	"context"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	greenhousev1alpha1 "github.com/cloudoperators/greenhouse/pkg/apis/greenhouse/v1alpha1"
)

// Webhook for the PluginDefinition custom resource.

func SetupPluginDefinitionWebhookWithManager(mgr ctrl.Manager) error {
	return setupWebhook(mgr,
		&greenhousev1alpha1.PluginDefinition{},
		webhookFuncs{
			defaultFunc:        DefaultPluginDefinition,
			validateCreateFunc: ValidateCreatePluginDefinition,
			validateUpdateFunc: ValidateUpdatePluginDefinition,
			validateDeleteFunc: ValidateDeletePluginDefinition,
		},
	)
}

//+kubebuilder:webhook:path=/mutate-greenhouse-sap-v1alpha1-plugindefinition,mutating=true,failurePolicy=fail,sideEffects=None,groups=greenhouse.sap,resources=plugindefinitions,verbs=create;update,versions=v1alpha1,name=mplugindefinition.kb.io,admissionReviewVersions=v1

func DefaultPluginDefinition(_ context.Context, _ client.Client, _ runtime.Object) error {
	return nil
}

//+kubebuilder:webhook:path=/validate-greenhouse-sap-v1alpha1-plugindefinition,mutating=false,failurePolicy=fail,sideEffects=None,groups=greenhouse.sap,resources=plugindefinitions,verbs=create;update,versions=v1alpha1,name=vplugindefinition.kb.io,admissionReviewVersions=v1

func ValidateCreatePluginDefinition(_ context.Context, _ client.Client, o runtime.Object) (admission.Warnings, error) {
	pluginDefinition, ok := o.(*greenhousev1alpha1.PluginDefinition)
	if !ok {
		return nil, nil
	}
	if err := validatePluginDefinitionMustSpecifyHelmChartOrUIApplication(pluginDefinition); err != nil {
		return nil, err
	}
	return nil, validatePluginDefinitionOptionValueAndType(pluginDefinition)
}

func ValidateUpdatePluginDefinition(_ context.Context, _ client.Client, _, o runtime.Object) (admission.Warnings, error) {
	pluginDefinition, ok := o.(*greenhousev1alpha1.PluginDefinition)
	if !ok {
		return nil, nil
	}
	if err := validatePluginDefinitionMustSpecifyHelmChartOrUIApplication(pluginDefinition); err != nil {
		return nil, err
	}
	return nil, validatePluginDefinitionOptionValueAndType(pluginDefinition)
}

func ValidateDeletePluginDefinition(_ context.Context, _ client.Client, _ runtime.Object) (admission.Warnings, error) {
	return nil, nil
}

func validatePluginDefinitionMustSpecifyHelmChartOrUIApplication(pluginDefinition *greenhousev1alpha1.PluginDefinition) error {
	if pluginDefinition.Spec.HelmChart == nil && pluginDefinition.Spec.UIApplication == nil {
		return apierrors.NewInvalid(pluginDefinition.GroupVersionKind().GroupKind(), pluginDefinition.GetName(), field.ErrorList{
			field.Required(field.NewPath("spec").Child("helmChart", "uiApplication"),
				"A PluginDefinition without both spec.helmChart and spec.uiApplication is invalid."),
		})
	}
	return nil
}

// validatePluginDefinitionOptionValueAndType validates that the type and value of each PluginOption matches.
func validatePluginDefinitionOptionValueAndType(pluginDefinition *greenhousev1alpha1.PluginDefinition) error {
	for _, option := range pluginDefinition.Spec.Options {
		if err := option.IsValid(); err != nil {
			return apierrors.NewInvalid(pluginDefinition.GroupVersionKind().GroupKind(), pluginDefinition.GetName(), field.ErrorList{
				field.Invalid(field.NewPath("spec").Child("options").Child("name"), option.Name,
					"A PluginOption Default must match the specified Type."),
			})
		}
	}
	return nil
}
