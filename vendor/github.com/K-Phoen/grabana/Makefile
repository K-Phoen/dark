WITH_COVERAGE?=false

ifeq ($(WITH_COVERAGE),true)
GOCMD_TEST?=go test -mod=vendor -coverpkg=./... -coverprofile=coverage.txt -covermode=atomic ./...
else
GOCMD_TEST?=go test -mod=vendor
endif

.PHONY: lint
lint:
	docker run --rm -v $(shell pwd):/app -w /app golangci/golangci-lint:v1.23.1 golangci-lint run -c .golangci.yaml

.PHONY: tests
tests:
	$(GOCMD_TEST) ./...

.PHONY: up
up:
	docker run -d \
        --name=grabana_prometheus \
        -p 9090:9090 \
        -v $(shell pwd)/testdata/prometheus.yml:/etc/prometheus/prometheus.yml \
        prom/prometheus
	docker run -d \
      -p 3000:3000 \
      --name=grabana_grafana \
      -e "GF_SECURITY_ADMIN_PASSWORD=secret" \
      grafana/grafana

.PHONY: down
down:
	docker rm -f grabana_grafana
	docker rm -f grabana_prometheus
