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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// KwokKind is the kind of the Logs.
	KWOKKind = "Kwok"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +genclient
// +kubebuilder:subresource:status
// +kubebuilder:rbac:groups=kwok.x-k8s.io,resources=Kwoks,verbs=create;get;list;watch
// +kubebuilder:rbac:groups=kwok.x-k8s.io,resources=Kwoks/status,verbs=update;patch

// Kwok provides Kwok configuration for a single pod.
type Kwok struct {
	//+k8s:conversion-gen=false
	metav1.TypeMeta `json:",inline"`
	// Standard list metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	metav1.ObjectMeta `json:"metadata"`
	// Spec holds spec for Kwok
	Spec KwokSpec `json:"spec"`
	// Status holds status for Kwok
	//+k8s:conversion-gen=false
	Status KwokStatus `json:"status,omitempty"`
}

// KwokStatus holds status for Kwok
type KwokStatus struct {
	// Conditions holds conditions for Kwok
	// +patchMergeKey=type
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=type
	Conditions []Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`
}

// KwokSpec holds spec for Kwok.
type KwokSpec struct {
	Image        string                `json:"image"`
	Selector     *metav1.LabelSelector `json:"selector"`
	NodeSelector *NodeSelector         `json:"nodeSelector,omitempty"`
}

type NodeSelector struct {
	MatchLabels      map[string]string `json:"matchLabels,omitempty"`
	MatchAnnotations map[string]string `json:"matchAnnotations,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true

// KwokList contains a list of Kwok
type KwokList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Kwok `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Kwok{}, &KwokList{})
}
