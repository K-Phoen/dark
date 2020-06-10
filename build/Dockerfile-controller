# Build
FROM golang:1.14-alpine3.11 AS build

RUN apk --update --no-cache add make=4.2.1-r2

WORKDIR /app
COPY . /app

RUN make controller_build

# Run

FROM alpine:3.11
ENTRYPOINT ["/app/dark-controller"]
WORKDIR /app
RUN apk --update --no-cache add ca-certificates=20191127-r2 && update-ca-certificates

COPY --from=build /app/dark-controller /app/dark-controller
