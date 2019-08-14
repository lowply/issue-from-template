FROM golang:1 as builder

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
WORKDIR /go/src
COPY src .
RUN GO111MODULE=on go build -o /go/bin/main

FROM alpine
RUN apk add --no-cache ca-certificates
RUN update-ca-certificates
COPY --from=builder /go/bin/main /bin/main
ENTRYPOINT main
