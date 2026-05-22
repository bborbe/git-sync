# Changelog

All notable changes to this project will be documented in this file.

## v1.5.11

- security: bump github.com/containerd/containerd to v1.7.32 (CVE-2026-46680, GHSA-fqw6-gf59-qr4w)

## v1.5.10

- security: bump github.com/go-git/go-git/v5 to v5.19.1 (CVE-2026-45570, CVE-2026-45571)

## v1.5.9

- security: suppress docker/docker CVE-2026-41567 / GHSA-x86f-5xw2-fm2r, CVE-2026-42306 / GHSA-rg2x-37c3-w2rh, CVE-2026-41568 / GHSA-vp62-88p7-qqf5 in .trivyignore/.osv-scanner.toml (no upstream fix; docker/docker module path unpatched, latest is v28.5.2)

## v1.5.8

- security: bump github.com/go-git/go-git/v5 to v5.19.0 (CVE-2026-45022)
- security: bump Go to 1.26.3 (GO-2026-4918, GO-2026-4971)
- chore: remove stale unused ignore entries from .osv-scanner.toml
- security: suppress docker/docker advisories GHSA-pxq6-2prw-chj9, GHSA-x744-4wpc-v9h2 in .trivyignore/.osv-scanner.toml (no upstream fix; latest is v28.5.2, advisory wants >= v29.3.1)

## v1.5.7

- chore: bump golangci-lint/v2 v2.11.4 → v2.12.1
- chore: bump osv-scanner/v2 v2.3.5 → v2.3.6
- chore: bump ginkgo/v2 v2.28.2 → v2.28.3, gomega v1.39.1 → v1.40.0
- fix: replace deprecated ioutil.TempDir with os.MkdirTemp in system_test.go (new govet inline check)
- chore: github.com/go-git/go-git/v5 confirmed at v5.18.0 (CVE-2026-41506 already patched)
- chore: github.com/docker/docker remains at v28.5.2+incompatible (v29.x not yet available on module proxy — advisory unfixable today)

## v1.5.6

- bump ginkgo/v2 v2.28.2, gosec v2.26.1, vuln v1.3.0
- bump anthropic-sdk-go v1.38.0, openai-go v3.32.0, genai v1.54.0
- add bahlo/generic-list-go, buger/jsonparser, invopop/jsonschema, mailru/easyjson, wk8/go-ordered-map deps
- remove anthropic-sdk-go replace directive

## v1.5.5

- bump golang.org/x/crypto, net, mod, tools, text, term, telemetry
- bump go-git/go-git v5.18.0
- bump golang.org/x/vuln v1.2.0
- add cache mounts for go-build and golangci-lint in dark-factory config

## v1.5.4

- bump Go toolchain to 1.26.2
- update counterfeiter to v6.12.2 and golang.org/x/sys to v0.43.0
- add OSV/Trivy ignores for new CVEs with no available fixes

## v1.5.3

- Update dependencies to fix security vulnerabilities (go-git/v5 v5.17.2)

## v1.5.2

- Update go-git/go-git to v5.17.1 (fix security vulnerabilities)

## v1.5.1

- Bump Go base image from 1.24.2 to 1.26.1
- Bump Alpine base image from 3.21 to 3.23
- Update addlicense, goimports-reviser, counterfeiter dependencies
- Add opencontainers/runtime-spec replace directive
- Enable parallel golangci-lint runners

## v1.5.0

- align Makefile with service pattern (lint, osv-scanner, gosec, trivy, golines, go-modtool)
- update tools.go with new tools, remove obsolete gogen-avro and golint
- add .golangci.yml
- fix lint issues: callbackUrl → callbackURL, handle stdin.Close error
- fix buildkit vulnerability (v0.23.2 → v0.28.1)

## v1.4.9

- fix: Preserve proxy environment variables in system test to allow git clone via container proxy

Please choose versions by [Semantic Versioning](http://semver.org/).

* MAJOR version when you make incompatible API changes,
* MINOR version when you add functionality in a backwards-compatible manner, and
* PATCH version when you make backwards-compatible bug fixes.

## v1.4.8

- remove vendor
- go mod update

## v1.4.7

- go mod update

## v1.4.6

- go mod update

## v1.4.5

- go mod update

## v1.4.4

- go mod update

## v1.4.3

- go mod update

## v1.4.2

- go mod upgrade
- Upload golang and alpine

## v1.4.1

- Upload golang and alpine

## v1.4.0

- Use go modules
- Add multi Docker tags

## v1.3.0

- Change build to multistage Dockerfile 

## v1.2.1

- Read callback-url from env CALLBACK_URL

## v1.2.0

- Add Parameter callback-url will be called after each git-sync
- Add Ginkgo
- Use Deps

## v1.1.6

- Update Alpine to 3.7

## v1.1.5

- add missign ssh_config

## v1.1.4

- Use Go 1.10

## v1.1.3

- Refactoring
- Use deps instead glide
- Add Dockerfiles

## v1.1.1

- Add vendor

## v1.1.0

- Update readme

## v1.0.0

- Initial Version
