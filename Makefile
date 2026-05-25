REGISTRY ?= docker.io
IMAGE ?= bborbe/git-sync
BRANCH ?= $(shell git rev-parse --abbrev-ref HEAD)
DIRS += $(shell find */* -maxdepth 0 -name Makefile -exec dirname "{}" \;)

.PHONY: default
default: precommit

.PHONY: precommit
precommit: ensure format generate test check addlicense
	@echo "ready to commit"

.PHONY: ensure
ensure:
	go mod tidy
	go mod verify
	rm -rf vendor

.PHONY: format
format:
	find . -type f -name 'go.mod' -not -path './vendor/*' -exec go run -mod=mod github.com/shoenig/go-modtool -w fmt "{}" \;
	find . -type f -name '*.go' -not -path './vendor/*' -exec gofmt -w "{}" +
	go run -mod=mod github.com/incu6us/goimports-reviser/v3 -project-name github.com/bborbe/git-sync -format -excludes vendor ./...
	find . -type d -name vendor -prune -o -type f -name '*.go' -print0 | xargs -0 -n 10 go run -mod=mod github.com/segmentio/golines --max-len=100 -w

.PHONY: generate
generate:
	rm -rf mocks avro
	mkdir -p mocks
	echo "package mocks" > mocks/mocks.go
	go generate -mod=mod ./...

.PHONY: test
test:
	go test -mod=mod -p=$${GO_TEST_PARALLEL:-1} -cover -race $(shell go list -mod=mod ./... | grep -v /vendor/)

.PHONY: check
check: lint vet vulncheck osv-scanner trivy

.PHONY: lint
lint:
	go run -mod=mod github.com/golangci/golangci-lint/v2/cmd/golangci-lint run --allow-parallel-runners --config .golangci.yml ./...

.PHONY: vet
vet:
	go vet -mod=mod $(shell go list -mod=mod ./... | grep -v /vendor/)

GOVULNCHECK_VERSION ?= v1.3.0
VULNCHECK_IGNORE ?= GO-2026-4923 GO-2026-4514 GO-2022-0470 GO-2026-4772 GO-2026-4771

.PHONY: vulncheck
vulncheck:
	@PKGS="$(shell go list -mod=mod ./... | grep -v /vendor/)"; \
	IGNORE_JSON=$$(printf '%s\n' $(VULNCHECK_IGNORE) | jq -R . | jq -s .); \
	REMAIN=$$(go run golang.org/x/vuln/cmd/govulncheck@$(GOVULNCHECK_VERSION) -format json $$PKGS 2>/dev/null | \
		jq -rs --argjson ignore "$$IGNORE_JSON" \
			'(map(select(.osv != null)) | map({key: .osv.id, value: (.osv.summary // "")}) | from_entries) as $$sum | \
			 map(select(.finding != null) | .finding) | \
			 map(select(.osv as $$o | $$ignore | index($$o) | not)) | \
			 map("\(.osv)\t\(.trace[-1].module)@\(.trace[-1].version) -> \(.fixed_version)\t\($$sum[.osv] // "")") | \
			 unique | .[]'); \
	if [ -n "$$REMAIN" ]; then \
		echo "Unexpected vulnerabilities (ignored: $(VULNCHECK_IGNORE)):"; \
		printf '%s\n' "$$REMAIN" | column -t -s "$$(printf '\t')"; \
		exit 1; \
	else \
		echo "No unignored vulnerabilities found"; \
	fi

.PHONY: osv-scanner
osv-scanner:
	@if [ -f .osv-scanner.toml ]; then \
		echo "Using .osv-scanner.toml"; \
		go run -mod=mod github.com/google/osv-scanner/v2/cmd/osv-scanner --config .osv-scanner.toml --recursive .; \
	else \
		echo "No config found, running default scan"; \
		go run -mod=mod github.com/google/osv-scanner/v2/cmd/osv-scanner --recursive .; \
	fi

.PHONY: trivy
trivy:
	trivy fs \
	--db-repository ghcr.io/aquasecurity/trivy-db \
	--scanners vuln,secret \
	--quiet \
	--no-progress \
	--disable-telemetry \
	--exit-code 1 .

.PHONY: addlicense
addlicense:
	go run -mod=mod github.com/google/addlicense -c "Benjamin Borbe" -y $$(date +'%Y') -l bsd $$(find . -name "*.go" -not -path './vendor/*')

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
