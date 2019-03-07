FROM golang:1.9

RUN echo 'alias ll="ls -lahG"' >> ~/.bashrc

RUN go get -u github.com/golang/dep/cmd/dep
RUN go get github.com/spf13/cobra/cobra
RUN go get github.com/derekparker/delve/cmd/dlv

ADD https://github.com/go-swagger/go-swagger/releases/download/0.13.0/swagger_linux_amd64 /usr/local/bin/swagger
RUN chmod +x /usr/local/bin/swagger
