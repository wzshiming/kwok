/*
Copyright The Kubernetes Authors.

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

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	gentype "k8s.io/client-go/gentype"
	v1alpha1 "sigs.k8s.io/kwok/pkg/apis/v1alpha1"
	apisv1alpha1 "sigs.k8s.io/kwok/pkg/client/clientset/versioned/typed/apis/v1alpha1"
)

// fakeClusterExecs implements ClusterExecInterface
type fakeClusterExecs struct {
	*gentype.FakeClientWithList[*v1alpha1.ClusterExec, *v1alpha1.ClusterExecList]
	Fake *FakeKwokV1alpha1
}

func newFakeClusterExecs(fake *FakeKwokV1alpha1) apisv1alpha1.ClusterExecInterface {
	return &fakeClusterExecs{
		gentype.NewFakeClientWithList[*v1alpha1.ClusterExec, *v1alpha1.ClusterExecList](
			fake.Fake,
			"",
			v1alpha1.SchemeGroupVersion.WithResource("clusterexecs"),
			v1alpha1.SchemeGroupVersion.WithKind("ClusterExec"),
			func() *v1alpha1.ClusterExec { return &v1alpha1.ClusterExec{} },
			func() *v1alpha1.ClusterExecList { return &v1alpha1.ClusterExecList{} },
			func(dst, src *v1alpha1.ClusterExecList) { dst.ListMeta = src.ListMeta },
			func(list *v1alpha1.ClusterExecList) []*v1alpha1.ClusterExec {
				return gentype.ToPointerSlice(list.Items)
			},
			func(list *v1alpha1.ClusterExecList, items []*v1alpha1.ClusterExec) {
				list.Items = gentype.FromPointerSlice(items)
			},
		),
		fake,
	}
}
