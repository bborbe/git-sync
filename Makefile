REGISTRY ?= docker.io
IMAGE ?= bborbe/git-sync
BRANCH ?= $(shell git rev-parse --abbrev-ref HEAD)
DIRS += $(shell find */* -maxdepth 0 -name Makefile -exec dirname "{}" \;)

default: precommit

build:
	docker build --no-cache --rm=true --platform=linux/amd64 -t $(REGISTRY)/$(IMAGE):$(BRANCH) -f Dockerfile .

upload:
	docker push $(REGISTRY)/$(IMAGE):$(BRANCH)

clean:
	docker rmi $(REGISTRY)/$(IMAGE):$(BRANCH) || true


precommit: ensure format generate test check addlicense
	@echo "ready to commit"

ensure:
	go mod verify
	go mod vendor

format:
	go run -mod=vendor github.com/incu6us/goimports-reviser/v3 -project-name github.com/bborbe/git-sync -format -excludes vendor ./...

generate:
	rm -rf mocks avro
	go generate -mod=vendor ./...

test:
	go test -mod=vendor -p=1 -cover -race $(shell go list -mod=vendor ./... | grep -v /vendor/)

check: vet errcheck vulncheck

vet:
	go vet -mod=vendor $(shell go list -mod=vendor ./... | grep -v /vendor/)

errcheck:
	go run -mod=vendor github.com/kisielk/errcheck -ignore '(Close|Write|Fprint)' $(shell go list -mod=vendor ./... | grep -v /vendor/)

addlicense:
	go run -mod=vendor github.com/google/addlicense -c "Benjamin Borbe" -y $$(date +'%Y') -l bsd $$(find . -name "*.go" -not -path './vendor/*')

vulncheck:
	go run -mod=vendor golang.org/x/vuln/cmd/govulncheck $(shell go list -mod=vendor ./... | grep -v /vendor/)
