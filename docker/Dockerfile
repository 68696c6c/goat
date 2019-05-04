FROM golang:1.12 as builder

RUN echo 'alias ll="ls -lahG"' >> ~/.bashrc
RUN go get github.com/derekparker/delve/cmd/dlv

FROM builder as base
WORKDIR /go/goat
ENV GO111MODULE=on
COPY go.mod go.sum ./
RUN go mod download
COPY . /go/goat

FROM base as app
RUN go build -o goat
#CMD ["/go/goat/goat"]
