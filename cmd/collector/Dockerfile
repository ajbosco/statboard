FROM alpine:latest

RUN apk add --no-cache ca-certificates

ADD build/collector /bin/collector

ENTRYPOINT ["/bin/collector"]