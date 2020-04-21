FROM golang:1.14.2

ARG GO_TAGS

COPY . $GOPATH/src/oauth2-proxy-nexus3/
RUN \
    cd $GOPATH/src/oauth2-proxy-nexus3 && \
    CGO_ENABLED=0 go build -tags="${GO_TAGS}" -o /tmp/oauth2-proxy-nexus3

FROM scratch

COPY --from=0 /tmp/oauth2-proxy-nexus3 /

ENTRYPOINT [ "/oauth2-proxy-nexus3" ]
