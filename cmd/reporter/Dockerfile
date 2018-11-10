FROM alpine:latest

RUN apk add --no-cache ca-certificates

RUN mkdir -p templates
COPY templates/index.html templates/

ADD build/reporter /bin/reporter

ENTRYPOINT ["/bin/reporter"]