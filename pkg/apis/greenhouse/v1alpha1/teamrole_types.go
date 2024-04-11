// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and Greenhouse contributors
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	greenhouseapis "github.com/cloudoperators/greenhouse/pkg/apis"
)

// TeamRoleSpec defines the desired state of a TeamRole
type TeamRoleSpec struct {
	// Rules is a list of rbacv1.PolicyRules used on a managed RBAC (Cluster)Role
	Rules []rbacv1.PolicyRule `json:"rules,omitempty"`
}

// TeamRoleStatus defines the observed state of a TeamRole
type TeamRoleStatus struct{}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// TeamRole is the Schema for the TeamRoles API
type TeamRole struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TeamRoleSpec   `json:"spec,omitempty"`
	Status TeamRoleStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// TeamRoleList contains a list of Role
type TeamRoleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TeamRole `json:"items"`
}

func init() {
	SchemeBuilder.Register(&TeamRole{}, &TeamRoleList{})
}

// GetRBACName returns the name of the rbacv1.ClusterRole that will be created on the remote cluster
func (tr *TeamRole) GetRBACName() string {
	return greenhouseapis.RBACPrefix + tr.GetName()
}
