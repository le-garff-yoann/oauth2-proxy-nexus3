FROM rnix/openssl-gost:3

WORKDIR /tmp

RUN openssl req \
    -new -newkey rsa:2048 -days 3650 -nodes -x509 \
    -subj '/C=FR/ST=Bre/L=Rennes/O=Acme Inc./OU=IT/CN=*.localhost' \
    -keyout global.key -out global.crt

FROM nginx:1.17.9

COPY --from=0 /tmp/global.* /etc/nginx/ssl/
COPY nginx/ /etc/nginx/conf.d
