---
status: completed
summary: Bump github.com/containerd/containerd from v1.7.30 to v1.7.32 to patch CVE-2026-46680
container: git-sync-exec-006-update-containerd-1-7-32
dark-factory-version: v0.164.0
created: "2026-05-22T18:00:00Z"
queued: "2026-05-22T17:11:42Z"
started: "2026-05-22T17:11:54Z"
completed: "2026-05-22T17:14:54Z"
---

<summary>
- Bumps `github.com/containerd/containerd` from v1.7.30 to v1.7.32 (CVE-2026-46680 / GHSA-fqw6-gf59-qr4w, High)
- Resolves Dependabot containerd advisory for bborbe/git-sync (2026-05-21)
- containerd is an indirect dep — must stay `// indirect`
- `make precommit` exits 0 after the change
- CHANGELOG `## Unreleased` documents the bump
</summary>

<objective>
Patch CVE-2026-46680 (containerd user-ID handling bypass allows runAsNonRoot evasion) by upgrading `github.com/containerd/containerd` to v1.7.32.
</objective>

<context>
Read `CLAUDE.md` for project conventions.
Read `docs/dod.md` for the Definition of Done.

Current `go.mod`:
- `github.com/containerd/containerd v1.7.30 // indirect` — has fix: bump to v1.7.32 (vulnerable range: >= 1.7.27, < 1.7.32)

Advisory:
- https://github.com/advisories/GHSA-fqw6-gf59-qr4w
- CVE: CVE-2026-46680, severity High, CVSS 7.3
</context>

<requirements>
1. Bump containerd:
   ```bash
   go get github.com/containerd/containerd@v1.7.32
   go mod tidy
   ```

2. Run `make precommit`. If it fails because trivy or osv-scanner reports NEW advisory IDs on `github.com/docker/docker` (existing suppression pattern) that are NOT yet in the ignore files, append them:
   - CVE-IDs → `.trivyignore` under the existing `# github.com/docker/docker indirect dep, no fix available via Go modules` block.
   - GHSA-IDs → new `[[IgnoredVulns]]` blocks in `.osv-scanner.toml` with `reason = "github.com/docker/docker indirect dep, no fix available"`.
   - Re-run `make precommit`. Cap at 3 iterations. If still failing after 3 iterations, stop and report blocker with the remaining unsuppressed advisory IDs in the prompt summary.

3. If the scanner reports advisories on any package OTHER than `github.com/docker/docker` or `github.com/containerd/containerd`, stop and report blocker — do NOT suppress.

4. Do NOT add ignore entries for any advisory that has a fix available. Only docker/docker is allowed to be suppressed.

5. Update `CHANGELOG.md` under `## Unreleased` (create the section as the first one below the top-of-file header if it does not exist):
   ```
   - security: bump github.com/containerd/containerd to v1.7.32 (CVE-2026-46680, GHSA-fqw6-gf59-qr4w)
   ```
   If additional docker/docker IDs were added in step 2, append a second line listing them.

6. Verify:
   - `go list -m github.com/containerd/containerd` reports `v1.7.32`
   - `make precommit` exits 0
</requirements>

<constraints>
- Only edit: `go.mod`, `go.sum`, `.trivyignore`, `.osv-scanner.toml`, `CHANGELOG.md`
- Do NOT bump unrelated deps
- Do NOT add containerd or docker/docker as a direct dep (must stay `// indirect`)
- Do NOT add a `replace` or `exclude` directive
- Do NOT commit — dark-factory handles git
- Existing tests must still pass
</constraints>

<verification>
```bash
go list -m github.com/containerd/containerd   # must print v1.7.32
make precommit                                 # must exit 0
```
</verification>
