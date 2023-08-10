FROM alpine:latest
WORKDIR /
COPY . .

ENTRYPOINT ["/bin/main"]
