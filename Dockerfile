FROM golang:1.19-alpine AS builder

LABEL maintainer="Christophe Rime <christopherime@me.com>"

ARG ARCH=amd64

ENV GOROOT /usr/local/go
ENV GOPATH /go
ENV PATH $GOPATH/bin:$GOROOT/bin:$PATH
ENV GO_VERSION 1.19
ENV GO111MODULE on
ENV CGO_ENABLED=0

# Build dependencies
WORKDIR /go/src/
COPY . .
RUN apk update && apk add make git
RUN make build

# Second stage
FROM alpine:latest

COPY --from=builder /go/src/build/amnf /usr/local/bin/amnf
RUN mkdir /etc/amnf
COPY config.yml /etc/amnf/config.yml
CMD ["/usr/local/bin/amnf", "-c", "/etc/amnf/config.yml"]
