REGISTRY?=kphoen
CONTROLLER_IMAGE=$(REGISTRY)/dark
CONVERTER_IMAGE=$(REGISTRY)/dark-converter

CONTROLLER_MAIN_SRC=cmd/controller/main.go
CONVERTER_MAIN_SRC=cmd/converter/main.go
GOCMD?=CGO_ENABLED=0 go

VERSION?=latest

WITH_COVERAGE?=false
ifeq ($(WITH_COVERAGE),true)
GOCMD_TEST?=go test -mod=vendor -coverpkg=./... -coverprofile=coverage.txt -covermode=atomic ./...
else
GOCMD_TEST?=go test -mod=vendor
endif

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

help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Build

.PHONY: controller_build
controller_build: ## Build controller binary.
	$(GOCMD) build -mod vendor -o dark-controller $(CONTROLLER_MAIN_SRC)

.PHONY: converter_build
converter_build: ## Build converter binary.
	$(GOCMD) build -mod vendor -o dark-converter $(CONVERTER_MAIN_SRC)

.PHONY: build
build: controller_build converter_build ## Build all binaries.

.PHONY: controller_image
controller_image: ## Build docker image with the controller.
	docker build -f build/Dockerfile-controller -t $(CONTROLLER_IMAGE):$(VERSION) .

.PHONY: converter_image
converter_image: ## Build docker image with the converter.
	docker build -f build/Dockerfile-converter -t $(CONVERTER_IMAGE):$(VERSION) .

.PHONY: images
images: converter_image controller_image ## Build all docker images.

.PHONY: controller_push
controller_push: controller_image ## Push docker image with the controller.
	docker push $(CONTROLLER_IMAGE):$(VERSION)

.PHONY: converter_push
converter_push: converter_image ## Push docker image with the converter.
	docker push $(CONVERTER_IMAGE):$(VERSION)

.PHONY: push
push: converter_push controller_push ## Push docker all images.

##@ Development

.PHONY: tests
tests: ## Run tests.
	$(GOCMD_TEST) ./...

.PHONY: clean
clean: ## Remove compiled binaries.
	rm dark-controller
	rm dark-converter
