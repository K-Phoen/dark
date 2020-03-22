# Build
FROM golang:1.13-alpine AS build

RUN apk --update --no-cache add make

WORKDIR /app
ADD . /app

RUN make build

# Run

FROM alpine:3.11
ENTRYPOINT ["/app/go-app"]
WORKDIR /app
RUN apk --update --no-cache add ca-certificates && update-ca-certificates

COPY --from=build /app/dark /app/dark
