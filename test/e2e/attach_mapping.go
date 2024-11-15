/*
Copyright 2024 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package e2e

import (
	"bytes"
	"context"
	"testing"

	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/e2e-framework/klient/k8s/resources"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"

	"sigs.k8s.io/kwok/pkg/apis/v1alpha1"
	"sigs.k8s.io/kwok/pkg/utils/exec"
	"sigs.k8s.io/kwok/test/e2e/helper"
)

// CaseAttachMapping creates a feature that tests attach mapping
func CaseAttachMapping(nodeName, namespace string) *features.FeatureBuilder {
	node := helper.NewNodeBuilder(nodeName).
		Build()
	pod0 := helper.NewPodBuilder("pod0").
		WithNamespace(namespace).
		WithNodeName(nodeName).
		Build()
	role := &rbacv1.Role{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "kube-system",
			Name:      "kwok-controller-attach-role",
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{"pods/attach"},
				Verbs:     []string{"get", "watch", "create"},
			},
		},
	}
	roleBinding := &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "kube-system",
			Name:      "kwok-controller-attach-role-binding",
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      "kwok-controller",
				Namespace: "kube-system",
			},
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "Role",
			Name:     "kwok-controller-attach-role",
		},
	}
	return features.New("Pod Attach Mapping").
		Setup(helper.CreateNode(node)).
		Setup(helper.CreatePod(pod0)).
		Teardown(helper.DeletePod(pod0)).
		Teardown(helper.DeleteNode(node)).
		Assess("test attach", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			client, err := resources.New(cfg.Client().RESTConfig())
			if err != nil {
				t.Fatal(err)
			}

			var pods corev1.PodList
			err = client.List(ctx, &pods, resources.WithLabelSelector("app=kwok-controller"))
			if err != nil {
				t.Fatal(err)
			}

			err = v1alpha1.AddToScheme(client.GetScheme())
			if err != nil {
				t.Fatal(err)
			}

			var ca = &v1alpha1.ClusterAttach{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: pods.Items[0].Namespace,
					Name:      pods.Items[0].Name,
				},
				Spec: v1alpha1.ClusterAttachSpec{
					Attaches: []v1alpha1.AttachConfig{
						{
							Mapping: &v1alpha1.MappingTarget{
								Namespace: pods.Items[0].Namespace,
								Name:      pods.Items[0].Name,
								Container: pods.Items[0].Spec.Containers[0].Name,
							},
						},
					},
				},
			}

			err = client.Create(ctx, ca)
			if err != nil {
				t.Fatal(err)
			}
			defer func() {
				_ = client.Delete(ctx, ca)
			}()

			err = client.Create(ctx, role)
			if err != nil {
				t.Fatal(err)
			}
			defer func() {
				_ = client.Delete(ctx, role)
			}()

			err = client.Create(ctx, roleBinding)
			if err != nil {
				t.Fatal(err)
			}
			defer func() {
				_ = client.Delete(ctx, roleBinding)
			}()

			buf := bytes.NewBuffer(nil)
			cmd, err := exec.Command(exec.WithAllWriteTo(exec.WithFork(ctx, true), buf), "kubectl", "attach", "-n", namespace, "pod0")
			if err != nil {
				t.Fatal(err)
			}
			defer func() {
				err = cmd.Process.Kill()
				if err != nil {
					t.Fatal(err)
				}
			}()

			t.Fatal(buf.String())

			return ctx
		})
}
