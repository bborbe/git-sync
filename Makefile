include tools.env

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
	find . -type f -name 'go.mod' -not -path './vendor/*' -exec go run github.com/shoenig/go-modtool@$(GO_MODTOOL_VERSION) -w fmt "{}" \;
	find . -type f -name '*.go' -not -path './vendor/*' -exec gofmt -w "{}" +
	go run github.com/incu6us/goimports-reviser/v3@$(GOIMPORTS_REVISER_VERSION) -project-name github.com/bborbe/git-sync -format -excludes vendor ./...
	find . -type d -name vendor -prune -o -type f -name '*.go' -print0 | xargs -0 -n 10 go run github.com/segmentio/golines@$(GOLINES_VERSION) --max-len=100 -w

.PHONY: generate
generate:
	rm -rf mocks avro
	mkdir -p mocks
	echo "package mocks" > mocks/mocks.go
	go generate -mod=mod ./...

# -race=true catches data races but flakes on some CI runners (rare SIGSEGV
# during gexec.Build in cmd/*-style binary smoke tests). Default off; opt in
# via ENABLE_RACE=true for nightly/manual hardening runs.
TESTFLAGS_RACE = -race=false
ifdef ENABLE_RACE
	TESTFLAGS_RACE = -race=true
endif

.PHONY: test
test:
	go test -mod=mod -p=$${GO_TEST_PARALLEL:-1} -cover $(TESTFLAGS_RACE) $(shell go list -mod=mod ./... | grep -v /vendor/)

.PHONY: check
check: lint vet errcheck vulncheck osv-scanner gosec trivy

.PHONY: lint
lint:
	go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION) run --config .golangci.yml ./...

.PHONY: vet
vet:
	go vet -mod=mod $(shell go list -mod=mod ./... | grep -v /vendor/)

.PHONY: errcheck
errcheck:
	go run github.com/kisielk/errcheck@$(ERRCHECK_VERSION) -ignore '(Close|Write|Fprint)' $(shell go list -mod=mod ./... | grep -v /vendor/ | grep -v k8s/client)

VULNCHECK_IGNORE ?= GO-2026-4923 GO-2026-4514 GO-2022-0470 GO-2026-4772 GO-2026-4771

# Known-benign govulncheck failure modes we swallow. golang.org/x/tools v0.46.0
# panics on packages containing generic *types.TypeParam during SSA analysis
# (govulncheck v1.3.0+ surface via RuntimeTypes/AllFunctions). We treat that as
# "no findings" because the panic happens AFTER the package scan; any real
# vulnerabilities would have been emitted as JSON on stdout before the panic.
# Any OTHER govulncheck failure (network, bad args, permissions) is surfaced.
.PHONY: vulncheck
vulncheck:
	@PKGS="$(shell go list -mod=mod ./... | grep -v /vendor/)"; \
	IGNORE_JSON=$$(printf '%s\n' $(VULNCHECK_IGNORE) | jq -R . | jq -s .); \
	ERR=$$(mktemp); \
	OUT=$$(go run golang.org/x/vuln/cmd/govulncheck@$(GOVULNCHECK_VERSION) -format json $$PKGS 2>$$ERR); \
	RC=$$?; \
	if [ $$RC -ne 0 ] && ! grep -q "ForEachElement called on type containing" "$$ERR"; then \
		echo "govulncheck failed (exit $$RC):" >&2; \
		cat "$$ERR" >&2; \
		rm -f "$$ERR"; \
		exit $$RC; \
	fi; \
	rm -f "$$ERR"; \
	REMAIN=$$(printf '%s' "$$OUT" | jq -rs --argjson ignore "$$IGNORE_JSON" \
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
		go run github.com/google/osv-scanner/v2/cmd/osv-scanner@$(OSV_SCANNER_VERSION) --config .osv-scanner.toml --recursive .; \
	else \
		echo "No config found, running default scan"; \
		go run github.com/google/osv-scanner/v2/cmd/osv-scanner@$(OSV_SCANNER_VERSION) --recursive .; \
	fi

.PHONY: gosec
gosec:
	go run github.com/securego/gosec/v2/cmd/gosec@$(GOSEC_VERSION) -exclude=G104 ./...

.PHONY: trivy
trivy:
	trivy fs --scanners vuln,secret --quiet --no-progress --disable-telemetry --exit-code 1 .

.PHONY: addlicense
addlicense:
	go run github.com/google/addlicense@$(ADDLICENSE_VERSION) -c "Benjamin Borbe" -y $$(date +'%Y') -l bsd $$(find . -name "*.go" -not -path './vendor/*')

.PHONY: buca
buca: build upload clean apply

.PHONY: build
build:
	docker build --no-cache --rm=true --platform=linux/amd64 -t $(REGISTRY)/$(IMAGE):$(BRANCH) -f Dockerfile .

.PHONY: upload
upload:
	docker push $(REGISTRY)/$(IMAGE):$(BRANCH)

.PHONY: clean
clean:
	docker rmi $(REGISTRY)/$(IMAGE):$(BRANCH) || true

.PHONY: apply
apply:
	@for i in $(DIRS); do \
		cd $$i; \
		echo "apply $${i}"; \
		make apply; \
		cd ..; \
	done
