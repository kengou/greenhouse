// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and Greenhouse contributors
// SPDX-License-Identifier: Apache-2.0

package admission

import (
	"context"
	"reflect"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	greenhousev1alpha1 "github.com/cloudoperators/greenhouse/pkg/apis/greenhouse/v1alpha1"
)

// Webhook for the RoleBinding custom resource.

func SetupTeamRoleBindingWebhookWithManager(mgr ctrl.Manager) error {
	return setupWebhook(mgr,
		&greenhousev1alpha1.TeamRoleBinding{},
		webhookFuncs{
			defaultFunc:        DefaultRoleBinding,
			validateCreateFunc: ValidateCreateRoleBinding,
			validateUpdateFunc: ValidateUpdateRoleBinding,
			validateDeleteFunc: ValidateDeleteRoleBinding,
		},
	)
}

//+kubebuilder:webhook:path=/mutate-greenhouse-sap-v1alpha1-teamrolebinding,mutating=true,failurePolicy=fail,sideEffects=None,groups=greenhouse.sap,resources=teamrolebindings,verbs=create;update,versions=v1alpha1,name=mrolebinding.kb.io,admissionReviewVersions=v1

func DefaultRoleBinding(_ context.Context, _ client.Client, _ runtime.Object) error {
	return nil
}

//+kubebuilder:webhook:path=/validate-greenhouse-sap-v1alpha1-teamrolebinding,mutating=false,failurePolicy=fail,sideEffects=None,groups=greenhouse.sap,resources=teamrolebindings,verbs=create;update,versions=v1alpha1,name=vrolebinding.kb.io,admissionReviewVersions=v1

func ValidateCreateRoleBinding(ctx context.Context, c client.Client, o runtime.Object) (admission.Warnings, error) {
	rb, ok := o.(*greenhousev1alpha1.TeamRoleBinding)
	if !ok {
		return nil, nil
	}

	// check if the referenced role exists
	var r greenhousev1alpha1.TeamRole
	if err := c.Get(ctx, client.ObjectKey{Namespace: rb.Namespace, Name: rb.Spec.TeamRoleRef}, &r); err != nil {
		if apierrors.IsNotFound(err) {
			return nil, apierrors.NewInvalid(rb.GroupVersionKind().GroupKind(), rb.Name, field.ErrorList{field.Invalid(field.NewPath("spec", "roleRef"), rb.Spec.TeamRoleRef, "role does not exist")})
		}
		return nil, apierrors.NewInternalError(err)
	}

	// check if the referenced team exists
	var t greenhousev1alpha1.Team
	if err := c.Get(ctx, client.ObjectKey{Namespace: rb.Namespace, Name: rb.Spec.TeamRef}, &t); err != nil {
		if apierrors.IsNotFound(err) {
			return nil, apierrors.NewInvalid(rb.GroupVersionKind().GroupKind(), rb.Name, field.ErrorList{field.Invalid(field.NewPath("spec", "teamRef"), rb.Spec.TeamRef, "team does not exist")})
		}
		return nil, apierrors.NewInternalError(err)
	}

	err := validateClusterNameOrSelector(ctx, c, rb)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func ValidateUpdateRoleBinding(ctx context.Context, c client.Client, old, cur runtime.Object) (admission.Warnings, error) {
	oldRB, ok := old.(*greenhousev1alpha1.TeamRoleBinding)
	if !ok {
		return nil, nil
	}
	curRB, ok := cur.(*greenhousev1alpha1.TeamRoleBinding)
	if !ok {
		return nil, nil
	}
	switch {
	case validateClusterNameOrSelector(ctx, c, curRB) != nil:
		return nil, apierrors.NewForbidden(
			schema.GroupResource{
				Group:    oldRB.GroupVersionKind().Group,
				Resource: oldRB.Kind,
			}, oldRB.Name, field.Forbidden(field.NewPath("spec"), "must contain either spec.clusterName or spec.clusterSelector"))
	case hasNamespacesChanged(oldRB, curRB):
		return nil, apierrors.NewForbidden(schema.GroupResource{Group: oldRB.GroupVersionKind().Group, Resource: oldRB.Kind}, oldRB.Name, field.Forbidden(field.NewPath("spec", "namespaces"), "cannot be changed"))
	default:
		return nil, nil
	}
}

func ValidateDeleteRoleBinding(_ context.Context, _ client.Client, _ runtime.Object) (admission.Warnings, error) {
	return nil, nil
}

// hasNamespacesChanged returns true if the namespaces in the old and current RoleBinding are different.
func hasNamespacesChanged(old, cur *greenhousev1alpha1.TeamRoleBinding) bool {
	return !reflect.DeepEqual(old.Spec.Namespaces, cur.Spec.Namespaces)
}

// validateClusterNameOrSelector checks if the TeamRoleBinding has a valid clusterName or clusterSelector but not both.
func validateClusterNameOrSelector(ctx context.Context, c client.Client, rb *greenhousev1alpha1.TeamRoleBinding) error {
	if rb.Spec.ClusterName != "" && (len(rb.Spec.ClusterSelector.MatchLabels) > 0 || len(rb.Spec.ClusterSelector.MatchExpressions) > 0) {
		return apierrors.NewInvalid(rb.GroupVersionKind().GroupKind(), rb.Name, field.ErrorList{field.Invalid(field.NewPath("spec", "clusterName"), rb.Spec.ClusterName, "cannot specify both spec.clusterName and spec.clusterSelector")})
	}

	if rb.Spec.ClusterName == "" && (len(rb.Spec.ClusterSelector.MatchLabels) == 0 && len(rb.Spec.ClusterSelector.MatchExpressions) == 0) {
		return apierrors.NewInvalid(rb.GroupVersionKind().GroupKind(), rb.Name, field.ErrorList{field.Invalid(field.NewPath("spec", "clusterName"), rb.Spec.ClusterName, "must specify either spec.clusterName or spec.clusterSelector")})
	}
	return nil
}
