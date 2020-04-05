FROM golang:1.14.2

COPY . $GOPATH/src/nexus3-proxy/
RUN \
    cd $GOPATH/src/nexus3-proxy && \
    CGO_ENABLED=0 go build -o /tmp/nexus3-proxy

FROM scratch

COPY --from=0 /tmp/nexus3-proxy /

ENTRYPOINT [ "/nexus3-proxy" ]
