FROM alpine:3.7
EXPOSE 8000 9000 80 443

RUN apk add -U --no-cache ca-certificates

ADD release/shop_server /bin/

ENTRYPOINT ["/bin/shop_server"]