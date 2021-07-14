FROM golang:1.16-alpine
#ENV GOPROXY="https://proxy.golang.org"
ENV GOPROXY="https://goproxy.io,direct"
COPY . /test
WORKDIR /test
RUN GO111MODULE="on" GOPROXY=$GOPROXY CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go test -count=1 -parallel 1 -v ./...