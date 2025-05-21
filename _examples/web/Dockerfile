FROM golang:1.24-bookworm AS env

ENV CGO_ENABLED=0
ENV GOPROXY=https://proxy.golang.org,direct

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

COPY --from=golangci/golangci-lint:v1.50.0 /usr/bin/golangci-lint /usr/bin/golangci-lint

RUN mkdir -p /app
WORKDIR /app/src


# Local development stage.
FROM env AS dev

RUN apt-get update && apt-get install wait-for-it

RUN go install gotest.tools/gotestsum@latest
RUN go install github.com/air-verse/air@latest


# Pipeline stage for running unit tests, linters, etc.
FROM env AS built

COPY . ./
RUN go build -o web


# Production stage.
# TODO: extend from alpine and copy files from base stage.
FROM env AS final
ARG aws_cli_arch=x86_64

RUN apt-get update && apt-get install -y jq wget unzip zip && apt-get clean

RUN wget "https://awscli.amazonaws.com/awscli-exe-linux-${aws_cli_arch}.zip" -O "awscliv2.zip"
RUN unzip awscliv2.zip
RUN ./aws/install

COPY --from=built /app /app

COPY ops/scripts/pre-deploy.sh /usr/local/bin/
RUN chmod a+x /usr/local/bin/pre-deploy.sh

CMD ["./web", "server", "80"]
