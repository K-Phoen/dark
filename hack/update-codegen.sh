#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
CODEGEN_PKG=${CODEGEN_PKG:-$(cd "${SCRIPT_ROOT}"; ls -d -1 ./vendor/k8s.io/code-generator 2>/dev/null || echo ../code-generator)}

bash "${CODEGEN_PKG}"/generate-groups.sh "deepcopy,client,informer,lister" \
  github.com/K-Phoen/dark/pkg/generated github.com/K-Phoen/dark/pkg/apis \
  controller:v1 \
  --output-base "$(dirname "${BASH_SOURCE[0]}")/.." \
  --go-header-file "${SCRIPT_ROOT}"/hack/boilerplate.go.txt

mv github.com/K-Phoen/dark/pkg/apis/controller/v1/zz_generated.deepcopy.go pkg/apis/controller/v1/zz_generated.deepcopy.go

rm -rf pkg/generated
mv github.com/K-Phoen/dark/pkg/generated pkg/generated
rm -rf github.com/