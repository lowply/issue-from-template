FROM golang:1

LABEL "com.github.actions.name"="Issue From Template"
LABEL "com.github.actions.description"="Issue From Template"
LABEL "com.github.actions.icon"="alert-circle"
LABEL "com.github.actions.color"="green"
LABEL "repository"="https://github.com/lowply/issue-from-template"
LABEL "homepage"="https://github.com/lowply/issue-from-template"
LABEL "maintainer"="Sho Mizutani <lowply@github.com>"

COPY src src
ENV GO111MODULE=on
WORKDIR /go/src
RUN go build -o ../bin/main

ADD entrypoint.sh /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]
