---
status: completed
summary: Bumped github.com/go-git/go-git/v5 from v5.19.0 to v5.19.1 to patch CVE-2026-45570 and CVE-2026-45571; updated CHANGELOG.md with Unreleased entry; make precommit exited 0 with no new advisory IDs requiring suppression.
container: git-sync-exec-005-update-go-git-5-19-1
dark-factory-version: v0.162.0
created: "2026-05-19T20:15:55Z"
queued: "2026-05-19T20:15:55Z"
started: "2026-05-19T20:16:01Z"
completed: "2026-05-19T20:17:50Z"
---

<summary>
- Bumps `github.com/go-git/go-git/v5` from v5.19.0 to v5.19.1 (CVE-2026-45570 Low, CVE-2026-45571 Moderate)
- Resolves Dependabot go-git advisories for bborbe/git-sync (2026-05-19 digest)
- For `github.com/docker/docker` (still no upstream fix — latest v28.5.2, advisory wants >= v29.3.1): keeps existing ignore-list entries; appends any NEW IDs that surface during scan
- `make precommit` exits 0 after the change
- CHANGELOG `## Unreleased` documents the bump
</summary>

<objective>
Patch CVE-2026-45570 / CVE-2026-45571 in `github.com/go-git/go-git/v5` by upgrading to v5.19.1.
</objective>

<context>
Read `CLAUDE.md` for project conventions.
Read `docs/dod.md` for the Definition of Done.

Current `go.mod` (both indirect):
- `github.com/go-git/go-git/v5 v5.19.0` — has fix: bump to v5.19.1
- `github.com/docker/docker v28.5.2+incompatible` — no upstream fix; existing ignore entries already cover known IDs

Existing ignore-file patterns to extend if scanner surfaces new IDs:

`.trivyignore`:
```
# github.com/docker/docker indirect dep, no fix available via Go modules
CVE-2026-34040
CVE-2026-33997
CVE-2026-41567
...
```

`.osv-scanner.toml`:
```toml
[[IgnoredVulns]]
id = "GHSA-x86f-5xw2-fm2r"
reason = "github.com/docker/docker indirect dep, no fix available"
```

Advisories:
- https://github.com/advisories?query=CVE-2026-45570
- https://github.com/advisories?query=CVE-2026-45571
</context>

<requirements>
1. Bump go-git/v5:
   ```bash
   go get github.com/go-git/go-git/v5@v5.19.1
   go mod tidy
   ```

2. Run `make precommit`. If it fails because trivy or osv-scanner reports NEW advisory IDs (CVE-* or GHSA-*) on `github.com/docker/docker` that are NOT yet in the ignore files, append them:
   - CVE-IDs → `.trivyignore` under the existing `# github.com/docker/docker indirect dep, no fix available via Go modules` block.
   - GHSA-IDs → new `[[IgnoredVulns]]` blocks in `.osv-scanner.toml` with `reason = "github.com/docker/docker indirect dep, no fix available"`.
   - Re-run `make precommit`. Cap at 3 iterations. If still failing after 3 iterations, stop and report blocker with the remaining unsuppressed advisory IDs in the prompt summary.

3. If the scanner reports advisories on any package OTHER than `github.com/docker/docker` or `github.com/go-git/go-git/v5`, stop and report blocker — do NOT suppress.

4. Do NOT add ignore entries for any advisory that has a fix available. Only docker/docker is allowed to be suppressed.

5. Update `CHANGELOG.md` under `## Unreleased` (create the section as the first one below the top-of-file header if it does not exist):
   ```
   - security: bump github.com/go-git/go-git/v5 to v5.19.1 (CVE-2026-45570, CVE-2026-45571)
   ```
   If additional docker/docker IDs were added in step 2, append a second line listing them.

6. Verify:
   - `go list -m github.com/go-git/go-git/v5` reports `v5.19.1`
   - `make precommit` exits 0
</requirements>

<constraints>
- Only edit: `go.mod`, `go.sum`, `.trivyignore`, `.osv-scanner.toml`, `CHANGELOG.md`
- Do NOT bump unrelated deps
- Do NOT add docker/docker as a direct dep (must stay `// indirect`)
- Do NOT add a `replace` or `exclude` directive
- Do NOT commit — dark-factory handles git
- Existing tests must still pass
</constraints>

<verification>
```bash
go list -m github.com/go-git/go-git/v5     # must print v5.19.1
make precommit                              # must exit 0
```
</verification>
