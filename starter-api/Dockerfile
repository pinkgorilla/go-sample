# Copyright 2018 The Payfazz Agent Authors.

FROM golang:1.11-alpine as builder
WORKDIR /go/src/github.com/pinkgorilla/go-sample/starter-api
ADD . .
RUN CGO_ENABLED=0 go install -installsuffix "static" ./...

FROM alpine
RUN apk add --no-cache ca-certificates
COPY --from=builder /go/bin /bin
USER nobody:nobody
ENTRYPOINT ["/bin/api-server"]
