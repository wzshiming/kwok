#!/usr/bin/env bash
# Copyright 2024 The Kubernetes Authors.
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

function blob_base() {
  local tag=${1}
  local repo
  repo="$(gh repo view --json url --jq '.url')"
  local url="${repo}/releases/download/${tag}"
  echo "${url}"
}

function update_index() {
  local index_dir="${1}"
  local tag=${2}
  local filepath="${3}"

  if ! gh release view "${tag}" >/dev/null; then
    echo "Please create release ${tag}"
    exit 1
  fi

  gh release upload "${tag}" "${filepath}"

  local url
  url="$(blob_base "${tag}")"

  helm repo index "${index_dir}" --merge "${index_dir}/index.yaml" --url "${url}"
}

function package_and_index() {
  local index_dir="${1}"
  local chart_dir="${2}"
  local chart_alias="${3}"

  local chart_name
  local chart_verison
  local chart_app_verison

  chart_name="$(yq eval '.name' "${chart_dir}/Chart.yaml")"
  chart_verison="$(yq eval '.version' "${chart_dir}/Chart.yaml")"
  chart_app_verison="$(yq eval '.appVersion' "${chart_dir}/Chart.yaml")"

  if yq eval ".entries.${chart_name}.[] | select(.version == \"${chart_verison}\")" "${index_dir}/index.yaml">/dev/null; then
    echo "Version ${chart_verison} already exists"
    return 0
  fi
  local guass_name="${index_dir}/${chart_name}-${chart_verison}.tgz"
  if [[ -f "${guass_name}" ]]; then
    echo "File ${guass_name} already exists"
    return 0
  fi

  helm package "${chart_dir}" --destination "${index_dir}"

  if [[ "${chart_alias}" != "" ]]; then
    local tmp_file="${index_dir}/${chart_alias}-${chart_verison}.tgz"
    mv "${guass_name}" "${tmp_file}"
    if ! update_index "${index_dir}" "${chart_app_verison}" "${tmp_file}"; then
      mv "${tmp_file}" "${guass_name}"
      exit 1
    fi
    mv "${tmp_file}" "${guass_name}"
  else
    update_index "${index_dir}" "${chart_app_verison}" "${guass_name}"
  fi
}

chart_dir="./charts"
index_dir="${ROOT_DIR}/site/static/charts"

package_and_index "${index_dir}" "${chart_dir}/kwok" "kwok-chart" || :
package_and_index "${index_dir}" "${chart_dir}/stage-fast" "kwok-stage-fast-chart" || :
