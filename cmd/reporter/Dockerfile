FROM alpine:latest

RUN apk add --no-cache ca-certificates

COPY templates templates
COPY static static

ADD build/reporter /bin/reporter

ENTRYPOINT ["/bin/reporter"]