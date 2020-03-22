# Build
FROM golang:1.13-alpine AS build

RUN apk --update --no-cache add make=4.2.1-r2

WORKDIR /app
COPY . /app

RUN make build

# Run

FROM alpine:3.11
ENTRYPOINT ["/app/go-app"]
WORKDIR /app
RUN apk --update --no-cache add ca-certificates=20191127-r1 && update-ca-certificates

COPY --from=build /app/dark /app/dark
