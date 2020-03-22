REGISTRY?=kphoen
IMAGE=$(REGISTRY)/dark

GO_MAIN_SRC=main.go controller.go
GOCMD?=CGO_ENABLED=0 go

VERSION?=latest

.PHONY: build
build:
	$(GOCMD) build -mod vendor -o dark $(GO_MAIN_SRC)

.PHONY: image
image:
	docker build -f build/Dockerfile -t $(IMAGE):$(VERSION) .

.PHONY: image
push: image
	docker push $(IMAGE):$(VERSION)

.PHONY: clean
clean:
	rm dark
