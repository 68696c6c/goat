FROM golang:1.19-buster as env

ENV CGO_ENABLED=0
ENV GOPROXY=https://proxy.golang.org,direct

COPY --from=golangci/golangci-lint:v1.50.0 /usr/bin/golangci-lint /usr/bin/golangci-lint

RUN mkdir -p /src
WORKDIR /src


# Local development stage.
FROM env as dev

RUN apt-get update && apt-get install wait-for-it

RUN go install gotest.tools/gotestsum@latest

COPY . ./

## Pipeline stage for running unit tests, linters, etc.
#FROM env as built
#
#COPY . ./
#RUN go build -o example
