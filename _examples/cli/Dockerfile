FROM golang:1.24-bookworm

ENV CGO_ENABLED=0
ENV GOPROXY=https://proxy.golang.org,direct

COPY --from=golangci/golangci-lint:v1.50.0 /usr/bin/golangci-lint /usr/bin/golangci-lint

RUN go install gotest.tools/gotestsum@latest

RUN mkdir -p /app
WORKDIR /app/src
