FROM alpine:latest
RUN mkdir -p /app
ADD index.html /app/index.html
ADD gopath/bin/web /app/web
WORKDIR /app
ENTRYPOINT ["/app/web"]
