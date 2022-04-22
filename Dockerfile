# golang image base smallest alpine
FROM golang:1.17.6-alpine3.15


# set go path and project path in the docker
ENV GOPROXY https://proxy.golang.com.cn,direct
# ENV GOPATH /go, aleardy set by golang image
ENV PROJECTPATH /ctrsploit

VOLUME ["/ctrsploit"]
RUN apk add git
RUN apk add upx
RUN go install github.com/mitchellh/gox@v1.0.1

ADD . $PROJECTPATH

RUN chmod +x $PROJECTPATH/build/build.sh
RUN ln -s $PROJECTPATH/build/build.sh /usr/local/bin/build_ctrsploit

# get ready for building ctrsploit
WORKDIR $PROJECTPATH

ENTRYPOINT ["sh", "-c"]
CMD ["build_ctrsploit"]
# enter the docker, run "build_ctrsploit" to build the ctrsploit binary files
