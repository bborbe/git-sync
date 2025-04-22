REGISTRY ?= docker.io
IMAGE ?= bborbe/git-sync
BRANCH ?= $(shell git rev-parse --abbrev-ref HEAD)
DIRS += $(shell find */* -maxdepth 0 -name Makefile -exec dirname "{}" \;)

default: precommit

precommit: ensure format generate test check addlicense
	@echo "ready to commit"

ensure:
	go mod tidy
	go mod verify
	rm -rf vendor

format:
	go run -mod=mod github.com/incu6us/goimports-reviser/v3 -project-name github.com/bborbe/git-sync -format -excludes vendor ./...

generate:
	rm -rf mocks avro
	go generate -mod=mod ./...

test:
	go test -mod=mod -p=1 -cover -race $(shell go list -mod=mod ./... | grep -v /vendor/)

check: vet errcheck vulncheck

vet:
	go vet -mod=mod $(shell go list -mod=mod ./... | grep -v /vendor/)

errcheck:
	go run -mod=mod github.com/kisielk/errcheck -ignore '(Close|Write|Fprint)' $(shell go list -mod=mod ./... | grep -v /vendor/)

addlicense:
	go run -mod=mod github.com/google/addlicense -c "Benjamin Borbe" -y $$(date +'%Y') -l bsd $$(find . -name "*.go" -not -path './vendor/*')

vulncheck:
	go run -mod=mod golang.org/x/vuln/cmd/govulncheck $(shell go list -mod=mod ./... | grep -v /vendor/)

.PHONY: build
build:
	go mod vendor
	docker build --no-cache --rm=true --platform=linux/amd64 -t $(REGISTRY)/$(IMAGE):$(BRANCH) -f Dockerfile .

.PHONY: upload
upload:
	docker push $(REGISTRY)/$(IMAGE):$(BRANCH)

.PHONY: clean
clean:
	docker rmi $(REGISTRY)/$(IMAGE):$(BRANCH) || true
	rm -rf vendor

.PHONY: apply
apply:
	@for i in $(DIRS); do \
		cd $$i; \
		echo "apply $${i}"; \
		make apply; \
		cd ..; \
	done

.PHONY: buca
buca: build upload clean apply
