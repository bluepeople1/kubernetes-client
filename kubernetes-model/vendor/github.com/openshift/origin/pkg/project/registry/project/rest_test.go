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
package project

import (
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apirequest "k8s.io/apiserver/pkg/endpoints/request"

	projectapi "github.com/openshift/origin/pkg/project/apis/project"
)

func TestProjectStrategy(t *testing.T) {
	ctx := apirequest.NewDefaultContext()
	if Strategy.NamespaceScoped() {
		t.Errorf("Projects should not be namespace scoped")
	}
	if Strategy.AllowCreateOnUpdate() {
		t.Errorf("Projects should not allow create on update")
	}
	project := &projectapi.Project{
		ObjectMeta: metav1.ObjectMeta{Name: "foo", ResourceVersion: "10"},
	}
	Strategy.PrepareForCreate(ctx, project)
	if len(project.Spec.Finalizers) != 1 || project.Spec.Finalizers[0] != projectapi.FinalizerOrigin {
		t.Errorf("Prepare For Create should have added project finalizer")
	}
	errs := Strategy.Validate(ctx, project)
	if len(errs) != 0 {
		t.Errorf("Unexpected error validating %v", errs)
	}
	invalidProject := &projectapi.Project{
		ObjectMeta: metav1.ObjectMeta{Name: "bar", ResourceVersion: "4"},
	}
	// ensure we copy spec.finalizers from old to new
	Strategy.PrepareForUpdate(ctx, invalidProject, project)
	if len(invalidProject.Spec.Finalizers) != 1 || invalidProject.Spec.Finalizers[0] != projectapi.FinalizerOrigin {
		t.Errorf("PrepareForUpdate should have preserved old.spec.finalizers")
	}
	errs = Strategy.ValidateUpdate(ctx, invalidProject, project)
	if len(errs) == 0 {
		t.Errorf("Expected a validation error")
	}
	if invalidProject.ResourceVersion != "4" {
		t.Errorf("Incoming resource version on update should not be mutated")
	}
}
