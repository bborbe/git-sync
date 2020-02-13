FROM golang:1.13.7 AS build
COPY . /go/src/github.com/bborbe/git-sync
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s" -a -installsuffix cgo -o /git-sync ./src/github.com/bborbe/git-sync
CMD ["/bin/bash"]

FROM alpine:3.11 as alpine
MAINTAINER Benjamin Borbe <bborbe@rocketnews.de>

RUN apk add --update ca-certificates git bash openssh && rm -rf /var/cache/apk/*

COPY files/ssh_config /root/.ssh/config
COPY --from=build /git-sync /git-sync

ENV GIT_SYNC_DEST /git
VOLUME ["/git"]

ENTRYPOINT ["/git-sync"]
