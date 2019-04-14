FROM golang:1.12

RUN echo 'alias ll="ls -lahG"' >> ~/.bashrc

RUN go get -u github.com/golang/dep/cmd/dep
RUN go get github.com/spf13/cobra/cobra
RUN go get github.com/derekparker/delve/cmd/dlv
