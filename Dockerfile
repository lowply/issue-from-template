FROM golang:1

LABEL "com.github.actions.name"="Issue From Template"
LABEL "com.github.actions.description"="Issue From Template"
LABEL "com.github.actions.icon"="alert-circle"
LABEL "com.github.actions.color"="green"
LABEL "repository"="https://github.com/lowply/issue-from-template"
LABEL "homepage"="https://github.com/lowply/issue-from-template"
LABEL "maintainer"="Sho Mizutani <lowply@github.com>"

WORKDIR /go/src
COPY src .
RUN GO111MODULE=on go build -o /go/bin/main

ADD entrypoint.sh /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]
