# Changelog

All notable changes to this project will be documented in this file.

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
