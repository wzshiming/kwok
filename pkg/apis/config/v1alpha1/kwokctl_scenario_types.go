/*
Copyright 2023 The Kubernetes Authors.

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
	"encoding/json"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// KwokctlScenarioKind is the kind of the kwokctl scenario.
	KwokctlScenarioKind = "KwokctlScenario"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KwokctlScenario provides scenario  definition for kwokctl.
type KwokctlScenario struct {
	//+k8s:conversion-gen=false
	metav1.TypeMeta `json:",inline"`
	// Standard list metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	metav1.ObjectMeta `json:"metadata,omitempty"`
	// Parameters is the parameters for the scenario.
	Parameters json.RawMessage `json:"parameters,omitempty"`
	// Steps is the steps for the scenario.
	Steps []ScenarioStep `json:"steps"`
	// Plays is the plays for the scenario.
	Plays []ScenarioPlay `json:"plays"`
}

// ScenarioStep is the step for the scenario.
type ScenarioStep struct {
	Name string `json:"name,omitempty"`
	Run  string `json:"run"`
}

// ScenarioPlay is the play for the scenario.
type ScenarioPlay struct {
	Name       string          `json:"name"`
	Parameters json.RawMessage `json:"parameters,omitempty"`
}
