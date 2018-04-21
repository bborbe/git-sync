IMAGE ?= bborbe/git-sync
REGISTRY ?= docker.io
ifeq ($(VERSION),)
	VERSION = $(shell git describe --tags `git rev-list --tags --max-count=1`)
endif

all: test install

prepare:
	go get -u golang.org/x/tools/cmd/goimports
	go get -u github.com/golang/lint/golint
	go get -u github.com/kisielk/errcheck
	go get -u github.com/golang/dep/cmd/dep
	go get -u github.com/Masterminds/glide
	go get -u github.com/bborbe/docker-utils/cmd/docker-remote-tag-exists

test:
	go test -cover -race $(shell go list ./... | grep -v /vendor/)

install:
	GOBIN=$(GOPATH)/bin GO15VENDOREXPERIMENT=1 go install *.go

goimports:
	go get golang.org/x/tools/cmd/goimports

format: goimports
	find . -type f -name '*.go' -not -path './vendor/*' -exec gofmt -w "{}" +
	find . -type f -name '*.go' -not -path './vendor/*' -exec goimports -w "{}" +

buildgo:
	CGO_ENABLED=0 GOOS=linux go build -ldflags "-s" -a -installsuffix cgo -o git-sync ./go/src/github.com/bborbe/git-sync

build:
	docker build --no-cache --rm=true -t $(REGISTRY)/$(IMAGE)-build:$(VERSION) -f ./Dockerfile.build .
	docker run -t $(REGISTRY)/$(IMAGE)-build:$(VERSION) /bin/true
	docker cp `docker ps -q -n=1 -f ancestor=$(REGISTRY)/$(IMAGE)-build:$(VERSION) -f status=exited`:/git-sync .
	docker rm `docker ps -q -n=1 -f ancestor=$(REGISTRY)/$(IMAGE)-build:$(VERSION) -f status=exited` || true
	docker build --no-cache --rm=true --tag=$(REGISTRY)/$(IMAGE):$(VERSION) -f Dockerfile.static .
	rm -f git-sync

upload:
	docker push $(REGISTRY)/$(IMAGE):$(VERSION)

clean:
	docker rmi $(REGISTRY)/$(IMAGE):$(VERSION) || true

trigger:
	@go get github.com/bborbe/docker-utils/cmd/docker-remote-tag-exists
	@exists=`docker-remote-tag-exists \
		-registry=${REGISTRY} \
		-repository="${IMAGE}" \
		-credentialsfromfile \
		-tag="${VERSION}" \
		-alsologtostderr \
		-v=0`; \
	trigger="build"; \
	if [ "$${exists}" = "true" ]; then \
		trigger="skip"; \
	fi; \
	echo $${trigger}
