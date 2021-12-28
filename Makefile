REGISTRY?=kphoen
CONTROLLER_IMAGE=$(REGISTRY)/dark
CONVERTER_IMAGE=$(REGISTRY)/dark-converter

VERSION?=latest

# ENVTEST_K8S_VERSION refers to the version of kubebuilder assets to be downloaded by envtest binary.
ENVTEST_K8S_VERSION = 1.22

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

WITH_COVERAGE?=false
ifeq ($(WITH_COVERAGE),true)
GOCMD_TEST?=go test -coverpkg=./... -coverprofile=coverage.txt -covermode=atomic ./...
else
GOCMD_TEST?=go test
endif

# Setting SHELL to bash allows bash commands to be executed by recipes.
# This is a requirement for 'setup-envtest.sh' in the test target.
# Options are set to exit when a recipe line exits non-zero or a piped command fails.
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

.PHONY: all
all: build-manager

##@ General

# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk commands is responsible for reading the
# entire set of makefiles included in this invocation, looking for lines of the
# file as xyz: ## something, and then pretty-format the target and help. Then,
# if there's a line with ##@ something, that gets pretty-printed as a category.
# More info on the usage of ANSI control characters for terminal formatting:
# https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_parameters
# More info on the awk command:
# http://linuxcommand.org/lc3_adv_awk.php

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

.PHONY: manifests
manifests: controller-gen ## Generate WebhookConfiguration, ClusterRole and CustomResourceDefinition objects.
	$(CONTROLLER_GEN) rbac:roleName=manager-role crd webhook paths="./..." output:crd:artifacts:config=config/crd/bases

.PHONY: generate
generate: controller-gen ## Generate code containing DeepCopy, DeepCopyInto, and DeepCopyObject method implementations.
	$(CONTROLLER_GEN) object paths="./..."

.PHONY: fmt
fmt: ## Run go fmt against code.
	go fmt ./...

.PHONY: vet
vet: ## Run go vet against code.
	go vet ./...

.PHONY: test
test: envtest ## Run tests.
	KUBEBUILDER_ASSETS="$(shell $(ENVTEST) use $(ENVTEST_K8S_VERSION) -p path)" $(GOCMD_TEST) ./...

.PHONY: lint
lint: ## Lints the code base.
	docker run --rm -v $(shell pwd):/app -w /app golangci/golangci-lint:v1.43.0 golangci-lint run -c build/golangci.yaml

##@ Build

.PHONY: build
build-manager: generate fmt vet ## Build manager binary.
	go build -o bin/manager cmd/controller/main.go

.PHONY: build
build-converter: fmt vet ## Build converter binary.
	go build -o bin/converter cmd/converter/main.go

.PHONY: run
run: ## Run a controller from your host.
	go run cmd/controller/main.go

.PHONY: docker-build-manager
docker-build-manager: ## Build docker image with the manager.
	DOCKER_BUILDKIT=1 docker build -f build/Dockerfile-controller -t ${CONTROLLER_IMAGE}:${VERSION} .

.PHONY: docker-build-converter
docker-build-converter: ## Build docker image with the converter.
	DOCKER_BUILDKIT=1 docker build -f build/Dockerfile-converter -t ${CONVERTER_IMAGE}:${VERSION} .

.PHONY: docker-build
docker-build: docker-build-manager docker-build-converter ## Build all docker images.

.PHONY: docker-push-manager
docker-push-manager: docker-build-manager ## Push docker image with the manager.
	docker push ${CONTROLLER_IMAGE}:${VERSION}

.PHONY: docker-push-converter
docker-push-converter: docker-build-converter ## Push docker image with the converter.
	docker push ${CONVERTER_IMAGE}:${VERSION}

.PHONY: docker-push
docker-push: docker-push-manager docker-push-converter ## Push all docker images.

##@ Development Environment
DEV_ENV_GRAFANA_URL=http://grafana.dark.localhost
DEV_ENV_GRAFANA_ADMIN_PASSWORD=$(shell kubectl get secret loki-grafana -o go-template='{{ index . "data" "admin-password" | base64decode }}')
DEV_ENV_GRAFANA_API_KEY=$(shell  curl --fail -XPOST -H "Content-Type: application/json" -d '{"name": "dark-dev-api-key-$(shell date +%s)", "role": "Admin"}' http://admin:$(DEV_ENV_GRAFANA_ADMIN_PASSWORD)@grafana.dark.localhost/api/auth/keys | jq .key)

.PHONY: dev-env-start
dev-env-start: dev-env-check-binaries dev-env-create-cluster dev-env-provision dev-env-grafana-credentials

.PHONY: dev-env-create-cluster
dev-env-create-cluster:
	k3d cluster create \
		--image="rancher/k3s:v1.21.7-k3s1" \
		-p "80:80@loadbalancer" \
		dark-dev

.PHONY: dev-env-delete-cluster
dev-env-delete-cluster:
	k3d cluster delete dark-dev

.PHONY: dev-env-provision
dev-env-provision:
	helm repo add grafana https://grafana.github.io/helm-charts
	helm repo update
	helm upgrade \
		--install loki grafana/loki-stack \
		--set grafana.enabled=true,prometheus.enabled=true,prometheus.alertmanager.persistentVolume.enabled=false,prometheus.server.persistentVolume.enabled=false
	kubectl apply -f config/crd/bases
	kubectl apply -f config/dev-env

.PHONY: dev-env-grafana-credentials
dev-env-grafana-credentials:
	@echo "==============="
	@echo "Grafana available at $(DEV_ENV_GRAFANA_URL)"
	@kubectl get secret loki-grafana -o go-template='{{range $$k,$$v := .data}}{{printf "%s: " $$k}}{{if not $$v}}{{$$v}}{{else}}{{$$v | base64decode}}{{end}}{{"\n"}}{{end}}'

.PHONY: dev-env-run-controller
dev-env-run-controler:
	GRAFANA_HOST=$(DEV_ENV_GRAFANA_URL) \
	GRAFANA_TOKEN=$(DEV_ENV_GRAFANA_API_KEY) \
	go run ./cmd/controller

.PHONY: dev-env-check-binaires
dev-env-check-binaries:
	@helm version >/dev/null 2>&1 || (echo "ERROR: helm is required."; exit 1)
	@k3d version >/dev/null 2>&1 || (echo "ERROR: k3d is required."; exit 1)
	@kubectl version --client >/dev/null 2>&1 || (echo "ERROR: kubectl is required."; exit 1)
	@jq --version >/dev/null 2>&1 ||(echo "ERROR: jq is required."; exit 1)

##@ Deployment

ifndef ignore-not-found
  ignore-not-found = false
endif

.PHONY: install
install: manifests kustomize ## Install CRDs into the K8s cluster specified in ~/.kube/config.
	$(KUSTOMIZE) build config/crd | kubectl apply -f -

.PHONY: uninstall
uninstall: manifests kustomize ## Uninstall CRDs from the K8s cluster specified in ~/.kube/config. Call with ignore-not-found=true to ignore resource not found errors during deletion.
	$(KUSTOMIZE) build config/crd | kubectl delete --ignore-not-found=$(ignore-not-found) -f -

.PHONY: deploy
deploy: manifests kustomize ## Deploy controller to the K8s cluster specified in ~/.kube/config.
	cd config/manager && $(KUSTOMIZE) edit set image controller=${IMG}
	$(KUSTOMIZE) build config/default | kubectl apply -f -

.PHONY: undeploy
undeploy: ## Undeploy controller from the K8s cluster specified in ~/.kube/config. Call with ignore-not-found=true to ignore resource not found errors during deletion.
	$(KUSTOMIZE) build config/default | kubectl delete --ignore-not-found=$(ignore-not-found) -f -

CONTROLLER_GEN = $(shell pwd)/bin/controller-gen
.PHONY: controller-gen
controller-gen: ## Download controller-gen locally if necessary.
	$(call go-get-tool,$(CONTROLLER_GEN),sigs.k8s.io/controller-tools/cmd/controller-gen@v0.7.0)

KUSTOMIZE = $(shell pwd)/bin/kustomize
.PHONY: kustomize
kustomize: ## Download kustomize locally if necessary.
	$(call go-get-tool,$(KUSTOMIZE),sigs.k8s.io/kustomize/kustomize/v3@v3.8.7)

ENVTEST = $(shell pwd)/bin/setup-envtest
.PHONY: envtest
envtest: ## Download envtest-setup locally if necessary.
	$(call go-get-tool,$(ENVTEST),sigs.k8s.io/controller-runtime/tools/setup-envtest@latest)

# go-get-tool will 'go get' any package $2 and install it to $1.
PROJECT_DIR := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))
define go-get-tool
@[ -f $(1) ] || { \
set -e ;\
TMP_DIR=$$(mktemp -d) ;\
cd $$TMP_DIR ;\
go mod init tmp ;\
echo "Downloading $(2)" ;\
GOBIN=$(PROJECT_DIR)/bin go get $(2) ;\
rm -rf $$TMP_DIR ;\
}
endef
