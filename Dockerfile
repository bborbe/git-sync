FROM golang:1.23.4 AS build
COPY . /workspace
WORKDIR /workspace
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -ldflags "-s" -a -installsuffix cgo -o /main
CMD ["/bin/bash"]

FROM alpine:3.21 AS alpine
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
