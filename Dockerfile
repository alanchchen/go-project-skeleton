# First stage container
FROM golang:1.12-alpine as builder

ARG APP
RUN apk add --no-cache make

ADD . $GOPATH/src/github.com/alanchchen/go-project-skeleton
RUN cd $GOPATH/src/github.com/alanchchen/go-project-skeleton && make ${APP} && mkdir -p /build/bin && mv build/bin/* /build/bin

# Second stage container
FROM alpine:3.9

RUN apk add --no-cache ca-certificates
COPY --from=builder /build/bin/* /usr/local/bin/

# Define your entrypoint or command
# ENTRYPOINT [""]
