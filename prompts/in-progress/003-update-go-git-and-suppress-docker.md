---
status: committing
summary: Bumped go-git/v5 to v5.19.0 and Go to 1.26.3, removed 6 stale osv-scanner ignore entries, confirmed existing docker/docker ignores cover all trivy/osv findings, make precommit exits 0.
container: git-sync-003-update-go-git-and-suppress-docker
dark-factory-version: v0.156.1-1-g04f3863-dirty
created: "2026-05-12T13:00:00Z"
queued: "2026-05-12T17:14:40Z"
started: "2026-05-12T17:58:33Z"
completed: "2026-05-12T17:53:34Z"
lastFailReason: 'execute prompt: docker run failed: wait command: exit status 143'
---

<summary>
- Bumps `github.com/go-git/go-git/v5` from v5.18.0 to v5.19.0 (CVE-2026-45022 High)
- Bumps Go from 1.26.2 to 1.26.3 (stdlib CVEs GO-2026-4918, GO-2026-4971)
- Resolves Dependabot go-git advisory for bborbe/git-sync
- For `github.com/docker/docker` (no upstream fix yet — latest is v28.5.2, advisory wants >= v29.3.1): keeps the existing advisory IDs in `.trivyignore` and `.osv-scanner.toml`; appends any NEW IDs that surface during scan
- Removes stale "unused ignore" entries from `.osv-scanner.toml` (osv-scanner errors on these)
- `make precommit` exits 0 after the change
- CHANGELOG `## Unreleased` documents the bumps and ignore-list cleanup
</summary>

<objective>
Patch CVE-2026-45022 in `github.com/go-git/go-git/v5` by upgrading to v5.19.0, and unblock `make precommit` for the docker/docker advisory by adding the unresolvable IDs to the project's existing security-scanner ignore lists.
</objective>

<context>
Read `CLAUDE.md` for project conventions.
Read `docs/dod.md` for the Definition of Done.

Current `go.mod` (both indirect):
- `github.com/go-git/go-git/v5 v5.18.0` — has fix: bump to v5.19.0
- `github.com/docker/docker v28.5.2+incompatible` — no upstream fix (latest available is v28.5.2; advisory wants >= v29.3.1)

Existing ignore-file patterns to mirror:

`.trivyignore` already contains entries like:
```
# github.com/docker/docker indirect dep, no fix available via Go modules
CVE-2026-34040
CVE-2026-33997
```

`.osv-scanner.toml` already contains entries like:
```toml
[[IgnoredVulns]]
id = "GHSA-pxq6-2prw-chj9"
reason = "github.com/docker/docker indirect dep, no fix available"
```

These files MUST grow if new advisory IDs surface during `make precommit`.
</context>

<requirements>
1. Bump go-git/v5:
   ```bash
   go get github.com/go-git/go-git/v5@v5.19.0
   go mod tidy
   ```

2. Bump Go version 1.26.2 → 1.26.3 (patches stdlib CVEs GO-2026-4918, GO-2026-4971):
   - Edit `go.mod`: change `go 1.26.2` to `go 1.26.3`
   - Edit `Dockerfile`: change `FROM golang:1.26.2` to `FROM golang:1.26.3`
   - Run `go mod tidy` again

3. Remove stale "unused ignore" entries from `.osv-scanner.toml`. After running `make precommit`, osv-scanner reports which `[[IgnoredVulns]]` IDs are unused. Delete those blocks. As of 2026-05-12 the unused IDs are: `GO-2026-4923`, `GHSA-6jwv-w5xf-7j27`, `GO-2022-0470`, `GO-2026-4772`, `GO-2026-4771`, `GHSA-xmrv-pmrh-hhx2`. Verify against the actual scanner output before deleting in case the list shifted.

4. Run `make precommit`. If it fails because trivy or osv-scanner reports NEW advisory IDs (CVE-* or GHSA-*) on `github.com/docker/docker` that are NOT yet in the ignore files, add them:
   - Append CVE-IDs to `.trivyignore` under the existing `# github.com/docker/docker indirect dep, no fix available via Go modules` block.
   - Append GHSA-IDs as new `[[IgnoredVulns]]` blocks in `.osv-scanner.toml` with `reason = "github.com/docker/docker indirect dep, no fix available"`.
   - Re-run `make precommit` until it exits 0.

5. Do NOT add ignore entries for any advisory that has a fix available (those must be patched, not suppressed). Only docker/docker is allowed to be suppressed in this prompt.

6. Update `CHANGELOG.md` under `## Unreleased`:
   ```
   - security: bump github.com/go-git/go-git/v5 to v5.19.0 (CVE-2026-45022)
   - security: bump Go to 1.26.3 (GO-2026-4918, GO-2026-4971)
   - chore: remove stale unused ignore entries from .osv-scanner.toml
   - security: suppress docker/docker advisories <list-IDs> in .trivyignore/.osv-scanner.toml (no upstream fix; latest is v28.5.2, advisory wants >= v29.3.1)
   ```
   Replace `<list-IDs>` with the IDs you actually added (or write "no new IDs" if the existing ignore lists already covered everything).

7. Verify:
   - `go list -m github.com/go-git/go-git/v5` reports `v5.19.0`
   - `grep '^go ' go.mod` shows `go 1.26.3`
   - `make precommit` exits 0
</requirements>

<constraints>
- Only edit: `go.mod`, `go.sum`, `Dockerfile`, `.trivyignore`, `.osv-scanner.toml`, `CHANGELOG.md`
- Do NOT bump unrelated deps
- Do NOT add docker/docker as a direct dep (must stay `// indirect`)
- Do NOT add a `replace` or `exclude` directive
- Do NOT commit — dark-factory handles git
- Existing tests must still pass
</constraints>

<verification>
```bash
go list -m github.com/go-git/go-git/v5     # must print v5.19.0
make precommit                              # must exit 0
grep -c "docker/docker" .trivyignore        # must be >= 1
```
</verification>
