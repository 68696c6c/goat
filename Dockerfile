FROM golang:1.12

RUN echo 'alias ll="ls -lahG"' >> ~/.bashrc

RUN go get github.com/derekparker/delve/cmd/dlv

WORKDIR /go/goat
ENV GO111MODULE=on
COPY go.mod go.sum ./
RUN go mod download
COPY . /go/goat
