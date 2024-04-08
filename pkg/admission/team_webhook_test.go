// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and Greenhouse contributors
// SPDX-License-Identifier: Apache-2.0

package admission_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	greenhousev1alpha1 "github.com/cloudoperators/greenhouse/pkg/apis/greenhouse/v1alpha1"
	"github.com/cloudoperators/greenhouse/pkg/test"
)

var _ = Describe("Validate Team Creation", func() {

	teamStub := greenhousev1alpha1.Team{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "test-org",
		},
		Spec: greenhousev1alpha1.TeamSpec{
			Description:    "Test Team",
			MappedIDPGroup: "IDP_GROUP_NAME_MATCHING_TEAM",
		},
	}

	BeforeEach(func() {
		pluginDefinition := greenhousev1alpha1.PluginDefinition{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "test-org",
				Name:      "test-plugindefinition-1",
			},
			Spec: greenhousev1alpha1.PluginDefinitionSpec{
				Description: "Test PluginDefinition 1",
				HelmChart: &greenhousev1alpha1.HelmChartReference{
					Name:       "./../../test/fixtures/myChart",
					Repository: "dummy",
					Version:    "1.0.0",
				},
			},
		}
		err := test.K8sClient.Create(test.Ctx, &pluginDefinition)
		Expect(err).ToNot(HaveOccurred(), "There should be no error when creating a pluginDefinition")

		pluginDefinition2 := greenhousev1alpha1.PluginDefinition{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-plugindefinition-2",
				Namespace: "test-org",
			},
			Spec: greenhousev1alpha1.PluginDefinitionSpec{
				Description: "Test PluginDefinition 2",
				HelmChart: &greenhousev1alpha1.HelmChartReference{
					Name:       "./../../test/fixtures/myChart",
					Repository: "dummy",
					Version:    "1.0.0",
				},
			},
		}
		err = test.K8sClient.Create(test.Ctx, &pluginDefinition2)
		Expect(err).ToNot(HaveOccurred(), "There should be no error when creating a pluginDefinition")
	})

	It("should correctly validate a team upon creation", func() {
		teamNoGreenhouseLabels := teamStub
		teamNoGreenhouseLabels.SetName("non-greenhouse-labels")

		teamGreenhouseLabels := teamStub
		teamGreenhouseLabels.SetName("greenhouse-labels")

		teamFaultyGreenhouseLabels := teamStub
		teamFaultyGreenhouseLabels.SetName("faulty-greenhouse-labels")

		By("correctly allowing creation of a team with non-greenhouse labels", func() {

			teamNoGreenhouseLabels.SetLabels(map[string]string{
				"some-key": "some-value",
			})
			err := test.K8sClient.Create(test.Ctx, &teamNoGreenhouseLabels)
			Expect(err).ToNot(HaveOccurred(), "There should be no error when creating a team with non-greenhouse labels")
		})

		By("correctly allowing creation of a tema with greenhouse labels that use whitelabeled labels and/or existing pluginDefinition names", func() {

			teamGreenhouseLabels.SetLabels(map[string]string{
				"greenhouse.sap/test-plugindefinition-1": "true",
				"greenhouse.sap/support-group":           "true",
			})
			err := test.K8sClient.Create(test.Ctx, &teamGreenhouseLabels)
			Expect(err).ToNot(HaveOccurred(), "There should be no error when creating a team with greenhouse labels that use existing pluginDefinition names")
		})

		By("correctly denying creation of a team with greenhouse labels that use non-existing pluginDefinition names", func() {

			teamFaultyGreenhouseLabels.SetLabels(map[string]string{
				"greenhouse.sap/test-plugindefinition-3": "true",
			})
			err := test.K8sClient.Create(test.Ctx, &teamFaultyGreenhouseLabels)
			Expect(err).To(HaveOccurred(), "There should be an error when creating a team with greenhouse labels that use non-existing pluginDefinition names")
			Expect(err.Error()).To(ContainSubstring("Only pluginDefinition names as greenhouse labels allowed."))
		})

		By("correctly allowing update of a team with non-greenhouse labels", func() {
			teamNoGreenhouseLabels.SetLabels(map[string]string{
				"some-key":       "some-value",
				"some-other-key": "some-other-value",
			})
			err := test.K8sClient.Update(test.Ctx, &teamNoGreenhouseLabels)
			Expect(err).ToNot(HaveOccurred(), "There should be no error when updating a team with non-greenhouse labels")
		})

		By("correctly allowing update of a team with greenhouse labels that use whitelisted labels and/or existing pluginDefinition names", func() {
			teamGreenhouseLabels.SetLabels(map[string]string{
				"greenhouse.sap/test-plugindefinition-1": "true",
				"greenhouse.sap/test-plugindefinition-2": "true",
				"greenhouse.sap/support-group":           "true",
			})
			err := test.K8sClient.Update(test.Ctx, &teamGreenhouseLabels)
			Expect(err).ToNot(HaveOccurred(), "There should be no error when updating a team with greenhouse labels that use existing pluginDefinition names")
		})

		By("correctly denying update of a team with greenhouse labels that use non-existing pluginDefinition names", func() {
			teamGreenhouseLabels.SetLabels(map[string]string{
				"greenhouse.sap/test-plugindefinition-3": "true",
			})
			err := test.K8sClient.Update(test.Ctx, &teamGreenhouseLabels)
			Expect(err).To(HaveOccurred(), "There should be an error when updating a team with greenhouse labels that use non-existing pluginDefinition names")
			Expect(err.Error()).To(ContainSubstring("Only pluginDefinition names as greenhouse labels allowed."))
		})

	})

})
