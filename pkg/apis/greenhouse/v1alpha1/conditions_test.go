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

package v1alpha1_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	greenhousev1alpha1 "github.com/cloudoperators/greenhouse/pkg/apis/greenhouse/v1alpha1"
)

var _ = Describe("Test conditions util functions", func() {

	var (
		timeNow = metav1.NewTime(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC))
	)

	DescribeTable("should correctly identify conditions",
		func(condition1 greenhousev1alpha1.Condition, condition2 greenhousev1alpha1.Condition, expected bool) {
			Expect(condition1.Equal(condition2)).To(Equal(expected))
		},
		Entry("should correctly identify equal conditions", greenhousev1alpha1.Condition{
			Type:               greenhousev1alpha1.ReadyCondition,
			Status:             metav1.ConditionTrue,
			LastTransitionTime: timeNow,
			Message:            "test",
		}, greenhousev1alpha1.Condition{
			Type:               greenhousev1alpha1.ReadyCondition,
			Status:             metav1.ConditionTrue,
			LastTransitionTime: timeNow,
			Message:            "test",
		}, true),
		Entry("should correctly identify conditions differing in the message", greenhousev1alpha1.Condition{
			Type:               greenhousev1alpha1.ReadyCondition,
			Status:             metav1.ConditionTrue,
			LastTransitionTime: timeNow,
			Message:            "test",
		}, greenhousev1alpha1.Condition{
			Type:               greenhousev1alpha1.ReadyCondition,
			Status:             metav1.ConditionTrue,
			LastTransitionTime: timeNow,
			Message:            "test2",
		}, false),
		Entry("should correctly identify conditions differing in the status", greenhousev1alpha1.Condition{
			Type:               greenhousev1alpha1.ReadyCondition,
			Status:             metav1.ConditionTrue,
			LastTransitionTime: timeNow,
			Message:            "test",
		}, greenhousev1alpha1.Condition{
			Type:               greenhousev1alpha1.ReadyCondition,
			Status:             metav1.ConditionFalse,
			LastTransitionTime: timeNow,
			Message:            "test",
		}, false),
		Entry("should correctly identify conditions differing in the type", greenhousev1alpha1.Condition{
			Type:               greenhousev1alpha1.ReadyCondition,
			Status:             metav1.ConditionTrue,
			LastTransitionTime: timeNow,
			Message:            "test",
		}, greenhousev1alpha1.Condition{
			Type:               greenhousev1alpha1.ReadyCondition,
			Status:             metav1.ConditionTrue,
			LastTransitionTime: timeNow,
			Message:            "test",
		}, true),
		Entry("should correctly ingore differing in the last transition time", greenhousev1alpha1.Condition{
			Type:               greenhousev1alpha1.ReadyCondition,
			Status:             metav1.ConditionTrue,
			LastTransitionTime: timeNow,
			Message:            "test",
		}, greenhousev1alpha1.Condition{
			Type:               greenhousev1alpha1.ReadyCondition,
			Status:             metav1.ConditionTrue,
			LastTransitionTime: metav1.NewTime(metav1.Now().AddDate(0, 0, -1)),
			Message:            "test",
		}, true),
	)

	DescribeTable("should correctly get the condition Status",
		func(condition greenhousev1alpha1.Condition, expected bool) {
			Expect(condition.IsTrue()).To(Equal(expected))
		},
		Entry("should correctly identify a true condition", greenhousev1alpha1.Condition{
			Type:               greenhousev1alpha1.ReadyCondition,
			Status:             metav1.ConditionTrue,
			LastTransitionTime: timeNow,
			Message:            "test",
		}, true),
		Entry("should correctly identify a false condition", greenhousev1alpha1.Condition{
			Type:               greenhousev1alpha1.ReadyCondition,
			Status:             metav1.ConditionFalse,
			LastTransitionTime: timeNow,
			Message:            "test",
		}, false),
	)

	DescribeTable("should correctly use SetCondition on StatusConditions",
		func(
			initialStatusConditions greenhousev1alpha1.StatusConditions,
			expected greenhousev1alpha1.StatusConditions,
			conditions ...greenhousev1alpha1.Condition,

		) {
			initialStatusConditions.SetConditions(conditions...)
			Expect(initialStatusConditions).To(Equal(expected))
		},
		Entry(
			"should correctly add a condition to empty StatusConditions",
			greenhousev1alpha1.StatusConditions{},
			greenhousev1alpha1.StatusConditions{
				Conditions: []greenhousev1alpha1.Condition{
					{
						Type:               greenhousev1alpha1.ReadyCondition,
						Status:             metav1.ConditionTrue,
						LastTransitionTime: timeNow,
						Message:            "test",
					},
				},
			},
			greenhousev1alpha1.Condition{
				Type:               greenhousev1alpha1.ReadyCondition,
				Status:             metav1.ConditionTrue,
				Message:            "test",
				LastTransitionTime: timeNow,
			}),
		Entry(
			"should correctly add a condition to StatusConditions with an existing condition",
			greenhousev1alpha1.StatusConditions{
				Conditions: []greenhousev1alpha1.Condition{
					{
						Type:               greenhousev1alpha1.HeadscaleReady,
						Status:             metav1.ConditionFalse,
						LastTransitionTime: timeNow,
						Message:            "test",
					},
				},
			},
			greenhousev1alpha1.StatusConditions{
				Conditions: []greenhousev1alpha1.Condition{
					{
						Type:               greenhousev1alpha1.HeadscaleReady,
						Status:             metav1.ConditionFalse,
						LastTransitionTime: timeNow,
						Message:            "test",
					},
					{
						Type:               greenhousev1alpha1.ReadyCondition,
						Status:             metav1.ConditionTrue,
						LastTransitionTime: timeNow,
						Message:            "test",
					},
				},
			},
			greenhousev1alpha1.Condition{
				Type:               greenhousev1alpha1.ReadyCondition,
				Status:             metav1.ConditionTrue,
				LastTransitionTime: timeNow,
				Message:            "test",
			}),
		Entry(
			"should correctly update a condition with matching Type in StatusConditions with a different condition",
			greenhousev1alpha1.StatusConditions{
				Conditions: []greenhousev1alpha1.Condition{
					{
						Type:               greenhousev1alpha1.HeadscaleReady,
						Status:             metav1.ConditionFalse,
						LastTransitionTime: timeNow,
						Message:            "test",
					},
					{
						Type:               greenhousev1alpha1.ReadyCondition,
						Status:             metav1.ConditionFalse,
						LastTransitionTime: timeNow,
						Message:            "test",
					},
				},
			},
			greenhousev1alpha1.StatusConditions{
				Conditions: []greenhousev1alpha1.Condition{
					{
						Type:               greenhousev1alpha1.HeadscaleReady,
						Status:             metav1.ConditionFalse,
						LastTransitionTime: timeNow,
						Message:            "test",
					},
					{
						Type:               greenhousev1alpha1.ReadyCondition,
						Status:             metav1.ConditionTrue,
						LastTransitionTime: timeNow,
						Message:            "test2",
					},
				},
			},
			greenhousev1alpha1.Condition{
				Type:               greenhousev1alpha1.ReadyCondition,
				Status:             metav1.ConditionTrue,
				LastTransitionTime: timeNow,
				Message:            "test2",
			}),
		Entry(
			"should ignore updating a condition with matching Type but differing LastTransitionTime in StatusConditions with a different condition",
			greenhousev1alpha1.StatusConditions{
				Conditions: []greenhousev1alpha1.Condition{
					{
						Type:               greenhousev1alpha1.HeadscaleReady,
						Status:             metav1.ConditionFalse,
						LastTransitionTime: timeNow,
						Message:            "test",
					},
					{
						Type:               greenhousev1alpha1.ReadyCondition,
						Status:             metav1.ConditionFalse,
						LastTransitionTime: timeNow,
						Message:            "test",
					},
				},
			},
			greenhousev1alpha1.StatusConditions{
				Conditions: []greenhousev1alpha1.Condition{
					{
						Type:               greenhousev1alpha1.HeadscaleReady,
						Status:             metav1.ConditionFalse,
						LastTransitionTime: timeNow,
						Message:            "test",
					},
					{
						Type:               greenhousev1alpha1.ReadyCondition,
						Status:             metav1.ConditionFalse,
						LastTransitionTime: timeNow,
						Message:            "test",
					},
				},
			},
			greenhousev1alpha1.Condition{
				Type:               greenhousev1alpha1.ReadyCondition,
				Status:             metav1.ConditionFalse,
				LastTransitionTime: metav1.NewTime(metav1.Now().AddDate(0, 0, -1)),
				Message:            "test",
			}),
		Entry(
			"should not update a conditions LastTransitionTime if only the message changes",
			greenhousev1alpha1.StatusConditions{
				Conditions: []greenhousev1alpha1.Condition{
					{
						Type:               greenhousev1alpha1.HeadscaleReady,
						Status:             metav1.ConditionFalse,
						LastTransitionTime: timeNow,
						Message:            "test",
					},
				},
			},
			greenhousev1alpha1.StatusConditions{
				Conditions: []greenhousev1alpha1.Condition{
					{
						Type:               greenhousev1alpha1.HeadscaleReady,
						Status:             metav1.ConditionFalse,
						LastTransitionTime: timeNow,
						Message:            "test2",
					},
				},
			},
			greenhousev1alpha1.Condition{
				Type:               greenhousev1alpha1.HeadscaleReady,
				Status:             metav1.ConditionFalse,
				LastTransitionTime: metav1.NewTime(metav1.Now().AddDate(0, 0, -1)),
				Message:            "test2",
			},
		),
		Entry(
			"should set and update multiple conditions",
			greenhousev1alpha1.StatusConditions{
				Conditions: []greenhousev1alpha1.Condition{
					{
						Type:               greenhousev1alpha1.HeadscaleReady,
						Status:             metav1.ConditionFalse,
						LastTransitionTime: timeNow,
						Message:            "test",
					},
				},
			},
			greenhousev1alpha1.StatusConditions{
				Conditions: []greenhousev1alpha1.Condition{
					{
						Type:               greenhousev1alpha1.HeadscaleReady,
						Status:             metav1.ConditionFalse,
						LastTransitionTime: timeNow,
						Message:            "test2",
					},
					{
						Type:               greenhousev1alpha1.ReadyCondition,
						Status:             metav1.ConditionTrue,
						LastTransitionTime: timeNow,
						Message:            "test",
					},
				},
			},
			greenhousev1alpha1.Condition{

				Type:               greenhousev1alpha1.HeadscaleReady,
				Status:             metav1.ConditionFalse,
				LastTransitionTime: metav1.NewTime(metav1.Now().AddDate(0, 0, -1)),
				Message:            "test2",
			},
			greenhousev1alpha1.Condition{
				Type:               greenhousev1alpha1.ReadyCondition,
				Status:             metav1.ConditionTrue,
				LastTransitionTime: timeNow,
				Message:            "test",
			},
		),
	)

	It("should correctly identify equal conditions", func() {
		By("identifying equal conditions")
		condition1 := greenhousev1alpha1.Condition{
			Type:               greenhousev1alpha1.ReadyCondition,
			Status:             metav1.ConditionTrue,
			LastTransitionTime: timeNow,
			Message:            "test",
		}
		condition2 := greenhousev1alpha1.Condition{
			Type:               greenhousev1alpha1.ReadyCondition,
			Status:             metav1.ConditionTrue,
			LastTransitionTime: metav1.NewTime(metav1.Now().AddDate(0, 0, -1)),
			Message:            "test",
		}
		Expect(condition1.Equal(condition2)).To(BeTrue())

	})
})
