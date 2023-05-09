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

TOOL_VERSION="v0.12.0"

ROOT_DIR="$(realpath "$(dirname "${BASH_SOURCE[0]}")"/..)"

COMMAND=()
if command -v controller-gen; then
  COMMAND=(controller-gen)
elif command -v "${ROOT_DIR}/bin/controller-gen"; then
  COMMAND=("${ROOT_DIR}/bin/controller-gen")
else
  GOBIN="${ROOT_DIR}/bin/" go install sigs.k8s.io/controller-tools/cmd/controller-gen@${TOOL_VERSION}
  COMMAND=("${ROOT_DIR}/bin/controller-gen")
fi

function update() {
  "${COMMAND[@]}" rbac:roleName=manager-role crd webhook paths="./pkg/apis/v1alpha1/..." output:crd:artifacts:config=config/crd/bases
  # "${COMMAND[@]}" object:headerFile="hack/boilerplate/boilerplate.go.txt" paths="./..."
}

cd "${ROOT_DIR}"
update || exit 1
