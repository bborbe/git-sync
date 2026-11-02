FROM golang:1.27.3 AS build
COPY . /workspace
WORKDIR /workspace
# No explicit -mod flag on purpose. vendor/ is gitignored, so a bare clone (the
# cluster Job clones the repo and has none) would fail a -mod=vendor build
# outright; with no flag Go uses vendor/ when it is present and the module cache
# otherwise, which keeps the local arm on vendor and lets the cluster arm build.
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s" -a -installsuffix cgo -o /main
CMD ["/bin/bash"]

FROM alpine:3.23 AS alpine
RUN apk --no-cache add \
	ca-certificates \
	rsync \
	openssh-client \
	tzdata \
	&& rm -rf /var/cache/apk/*

COPY files/ssh_config /root/.ssh/config
COPY --from=build /main /main

ENV GIT_SYNC_DEST /git
VOLUME ["/git"]

ENTRYPOINT ["/main"]
