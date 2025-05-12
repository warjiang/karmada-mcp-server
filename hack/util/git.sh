#!/bin/bash
# Copyright 2024 The Karmada Authors.
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

set -euo pipefail


function util::git::get_version() {
  # git describe --tags --dirty
  # GIT_VERSION=$(git rev-parse --abbrev-ref HEAD)
  v=$(git rev-parse --abbrev-ref HEAD)
  v=${v//\//-}
  echo "${v}"
}

function util::git::version_ldflags() {
  # Git information
  GIT_VERSION=$(git rev-parse --abbrev-ref HEAD)
  GIT_COMMIT_HASH=$(git rev-parse HEAD)
  if git_status=$(git status --porcelain 2>/dev/null) && [[ -z ${git_status} ]]; then
    GIT_TREE_STATE="clean"
  else
    GIT_TREE_STATE="dirty"
  fi
  BUILD_DATE=$(date -u +'%Y-%m-%dT%H:%M:%SZ')
  LDFLAGS="-X github.com/warjiang/karmada-mcp-server/pkg/environment.gitVersion=${GIT_VERSION} \
           -X github.com/warjiang/karmada-mcp-server/pkg/environment.gitCommit=${GIT_COMMIT_HASH} \
           -X github.com/warjiang/karmada-mcp-server/pkg/environment.gitTreeState=${GIT_TREE_STATE} \
           -X github.com/warjiang/karmada-mcp-server/pkg/environment.buildDate=${BUILD_DATE}"
  echo "$LDFLAGS"
}
