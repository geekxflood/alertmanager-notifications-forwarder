FROM golang:1.20-alpine AS builder

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
RUN apk update && apk add git
RUN go get ./...
RUN mkdir /go/src/build 
RUN go build -a -gcflags=all="-l -B" -ldflags="-w -s" -o build/amnf ./...

# Second stage
FROM alpine:3.17

COPY --from=builder /go/src/build/amnf /usr/local/bin/amnf
CMD ["/usr/local/bin/amnf"]
EXPOSE 9740
