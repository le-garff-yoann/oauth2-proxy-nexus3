FROM alpine/openssl:3.5.1 AS openssl

WORKDIR /tmp

RUN openssl req \
    -new -newkey rsa:2048 -days 3650 -nodes -x509 \
    -subj "/C=FR/ST=Bre/L=Rennes/O=Acme Inc./OU=IT/CN=*.localhost" \
    -keyout global.key -out global.crt

FROM nginx:1.29

COPY --from=openssl /tmp/global.* /etc/nginx/ssl/
COPY nginx/ /etc/nginx/conf.d
