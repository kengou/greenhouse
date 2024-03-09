// Copyright 2024-2026 SAP SE or an SAP affiliate company and Greenhouse contributors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// User specifies a human person.
type User struct {
	// ID is the unique identifier of the user.
	ID string `json:"id"`
	// FirstName of the user.
	FirstName string `json:"firstName"`
	// LastName of the user.
	LastName string `json:"lastName"`
	// Email of the user.
	Email string `json:"email"`
}

// TeamMembershipSpec defines the desired state of TeamMembership
type TeamMembershipSpec struct {
	// Members list users that are part of a team.
	// +optional
	Members []User `json:"members,omitempty"`
}

// TeamMembershipStatus defines the observed state of TeamMembership
type TeamMembershipStatus struct {
	// LastSyncedTime is the information when was the last time the membership was synced
	// +optional
	LastSyncedTime *metav1.Time `json:"lastSyncedTime,omitempty"`
	// LastChangedTime is the information when was the last time the membership was actually changed
	// +optional
	LastChangedTime *metav1.Time `json:"lastUpdateTime,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// TeamMembership is the Schema for the teammemberships API
type TeamMembership struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TeamMembershipSpec   `json:"spec,omitempty"`
	Status TeamMembershipStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// TeamMembershipList contains a list of TeamMembership
type TeamMembershipList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TeamMembership `json:"items"`
}

func init() {
	SchemeBuilder.Register(&TeamMembership{}, &TeamMembershipList{})
}
