/**
 * Copyright (C) 2015 Red Hat, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *         http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package builds

import (
	"fmt"

	g "github.com/onsi/ginkgo"
	o "github.com/onsi/gomega"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	exutil "github.com/openshift/origin/test/extended/util"
)

var _ = g.Describe("[Feature:Builds] build with empty source", func() {
	defer g.GinkgoRecover()
	var (
		buildFixture = exutil.FixturePath("testdata", "builds", "test-nosrc-build.json")
		oc           = exutil.NewCLI("cli-build-nosrc", exutil.KubeConfigPath())
		exampleBuild = exutil.FixturePath("testdata", "builds", "test-build-app")
	)

	g.Context("", func() {

		g.BeforeEach(func() {
			exutil.DumpDockerInfo()
		})

		g.JustBeforeEach(func() {
			g.By("waiting for builder service account")
			err := exutil.WaitForBuilderAccount(oc.KubeClient().Core().ServiceAccounts(oc.Namespace()))
			o.Expect(err).NotTo(o.HaveOccurred())
			oc.Run("create").Args("-f", buildFixture).Execute()
		})

		g.AfterEach(func() {
			if g.CurrentGinkgoTestDescription().Failed {
				exutil.DumpPodStates(oc)
				exutil.DumpPodLogsStartingWith("", oc)
			}
		})

		g.Describe("started build", func() {
			g.It("should build even with an empty source in build config", func() {
				g.By("starting the empty source build")
				br, err := exutil.StartBuildAndWait(oc, "nosrc-build", fmt.Sprintf("--from-dir=%s", exampleBuild))
				br.AssertSuccess()

				g.By(fmt.Sprintf("verifying the status of %q", br.BuildPath))
				build, err := oc.BuildClient().Build().Builds(oc.Namespace()).Get(br.Build.Name, metav1.GetOptions{})
				o.Expect(err).NotTo(o.HaveOccurred())
				o.Expect(build.Spec.Source.Dockerfile).To(o.BeNil())
				o.Expect(build.Spec.Source.Git).To(o.BeNil())
				o.Expect(build.Spec.Source.Images).To(o.BeNil())
				o.Expect(build.Spec.Source.Binary).NotTo(o.BeNil())
			})
		})
	})
})
