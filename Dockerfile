# First stage container
FROM golang:1.24-alpine as builder

ENV GO111MODULE=on
RUN apk add --no-cache make git

RUN mkdir -p /src
WORKDIR /src

COPY go.mod .
COPY go.sum .

RUN go mod download
COPY . .

RUN make && mkdir -p /build/bin && mv build/bin/* /build/bin

# Second stage container
FROM alpine:3.22

RUN apk add --no-cache ca-certificates
COPY --from=builder /build/bin/* /usr/local/bin/

# Define your entrypoint or command
# ENTRYPOINT [""]
