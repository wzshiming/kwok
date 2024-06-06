#!/usr/bin/env bash
# Copyright 2023 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -o errexit
set -o nounset
set -o pipefail

DIR="$(dirname "${BASH_SOURCE[0]}")"

ROOT_DIR="$(realpath "${DIR}/..")"

function deepcopy-gen() {
  go run k8s.io/code-generator/cmd/deepcopy-gen "$@"
}

function defaulter-gen() {
  go run k8s.io/code-generator/cmd/defaulter-gen "$@"
}

function conversion-gen() {
  go run k8s.io/code-generator/cmd/conversion-gen "$@"
}

function client-gen() {
  go run k8s.io/code-generator/cmd/client-gen "$@"
}

function gen() {
  rm -rf \
    "${ROOT_DIR}/pkg/apis/internalversion"/zz_generated.*.go \
    "${ROOT_DIR}/pkg/apis/v1alpha1"/zz_generated.*.go \
    "${ROOT_DIR}/pkg/apis/config/v1alpha1"/zz_generated.*.go \
    "${ROOT_DIR}/pkg/apis/action/v1alpha1"/zz_generated.*.go
  echo "Generating deepcopy"
  deepcopy-gen \
    ./pkg/apis/internalversion/ \
    ./pkg/apis/v1alpha1/ \
    ./pkg/apis/config/v1alpha1/ \
    ./pkg/apis/action/v1alpha1/ \
    ./pkg/apis/operator/v1alpha1/ \
    --output-file zz_generated.deepcopy.go \
    --go-header-file ./hack/boilerplate/boilerplate.generatego.txt
  echo "Generating defaulter"
  defaulter-gen \
    ./pkg/apis/v1alpha1/ \
    ./pkg/apis/config/v1alpha1/ \
    ./pkg/apis/action/v1alpha1/ \
    ./pkg/apis/operator/v1alpha1/ \
    --output-file zz_generated.defaults.go \
    --go-header-file ./hack/boilerplate/boilerplate.generatego.txt
  echo "Generating conversion"
  conversion-gen \
    ./pkg/apis/internalversion/ \
    --output-file zz_generated.conversion.go \
    --go-header-file ./hack/boilerplate/boilerplate.generatego.txt

  rm -rf "${ROOT_DIR}/pkg/client"
  echo "Generating client"
  client-gen \
    --clientset-name versioned \
    --input-base "" \
    --input sigs.k8s.io/kwok/pkg/apis/v1alpha1 \
    --output-pkg sigs.k8s.io/kwok/pkg/client/clientset \
    --output-dir ./pkg/client/clientset \
    --go-header-file ./hack/boilerplate/boilerplate.generatego.txt \
    --plural-exceptions="Logs:Logs,ClusterLogs:ClusterLogs"

  client-gen \
    --clientset-name versioned \
    --input-base "" \
    --input sigs.k8s.io/kwok/pkg/apis/operator/v1alpha1 \
    --output-pkg sigs.k8s.io/kwok/pkg/client/operator/clientset \
    --output-dir ./pkg/client/operator/clientset \
    --go-header-file ./hack/boilerplate/boilerplate.generatego.txt
}

cd "${ROOT_DIR}" && gen
