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

.PHONY: controller_build
controller_build:
	$(GOCMD) build -mod vendor -o dark-controller $(CONTROLLER_MAIN_SRC)

.PHONY: converter_build
converter_build:
	$(GOCMD) build -mod vendor -o dark-converter $(CONVERTER_MAIN_SRC)

.PHONY: build
build: controller_build converter_build

.PHONY: controller_image
controller_image:
	docker build -f build/Dockerfile-controller -t $(CONTROLLER_IMAGE):$(VERSION) .

.PHONY: converter_image
converter_image:
	docker build -f build/Dockerfile-converter -t $(CONVERTER_IMAGE):$(VERSION) .

.PHONY: images
images: converter_image controller_image

.PHONY: controller_push
controller_push: controller_image
	docker push $(CONTROLLER_IMAGE):$(VERSION)

.PHONY: converter_push
converter_push: converter_image
	docker push $(CONVERTER_IMAGE):$(VERSION)

.PHONY: push
push: images

.PHONY: tests
tests:
	$(GOCMD_TEST) ./...

.PHONY: clean
clean:
	rm dark-controller
	rm dark-converter
