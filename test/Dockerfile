FROM golang:1.13-alpine

# Install git
RUN set -ex; \
    apk update; \
    apk add --no-cache git

# Set working directory
WORKDIR /go/src/github.com/mediabuyerbot/go-webdriver

# Run tests
CMD CGO_ENABLED=0 go test ./test/... -v