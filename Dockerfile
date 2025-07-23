FROM golang:1.23-alpine3.21 AS build

ARG GO_TAGS

COPY . $GOPATH/src/oauth2-proxy-nexus3/

WORKDIR $GOPATH/src/oauth2-proxy-nexus3

RUN CGO_ENABLED=0 go build -tags="${GO_TAGS}" -o /tmp/oauth2-proxy-nexus3

FROM scratch

COPY --from=build /tmp/oauth2-proxy-nexus3 /

ENTRYPOINT [ "/oauth2-proxy-nexus3" ]
