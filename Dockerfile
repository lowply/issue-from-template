FROM golang:1 AS builder
WORKDIR /go/src
COPY src .
RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/main .

LABEL "com.github.actions.name"="Issue From Template"
LABEL "com.github.actions.description"="Issue From Template"
LABEL "com.github.actions.icon"="alert-circle"
LABEL "com.github.actions.color"="green"
LABEL "repository"="https://github.com/lowply/issue-from-template"
LABEL "homepage"="https://github.com/lowply/issue-from-template"
LABEL "maintainer"="Sho Mizutani <lowply@github.com>"

FROM alpine:latest AS runner
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/bin/main /usr/local/bin/

ADD entrypoint.sh /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]
