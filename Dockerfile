# First stage container
FROM golang:1.12-alpine as builder

ENV GO111MODULE=on
RUN apk add --no-cache make git

RUN mkdir -p /src
WORKDIR /src

COPY go.mod .
COPY go.sum .

RUN go mod download
COPY . .

ARG APP
RUN make ${APP} && mkdir -p /build/bin && mv build/bin/* /build/bin

# Second stage container
FROM alpine:3.9

RUN apk add --no-cache ca-certificates
COPY --from=builder /build/bin/* /usr/local/bin/

# Define your entrypoint or command
# ENTRYPOINT [""]
