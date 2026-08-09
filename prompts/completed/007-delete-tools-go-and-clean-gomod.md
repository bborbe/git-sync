---
status: completed
summary: Deleted tools.go and rewrote go.mod to keep only 4 real application deps; go.mod shrank from 448 to 26 lines; five replace workarounds eliminated; go-git CVEs gone; osv-scanner passes cleanly
execution_id: git-sync-drop-tools-go-exec-007-delete-tools-go-and-clean-gomod
dark-factory-version: v0.192.9
created: "2026-08-09T14:45:19Z"
queued: "2026-08-09T14:45:19Z"
started: "2026-08-09T14:49:06Z"
completed: "2026-08-09T14:51:02Z"
---

# Delete tools.go and remove tool-dependency pollution from go.mod

<summary>
- Tool CLIs are no longer declared as Go module dependencies of this project
- The project's dependency list shrinks to only the four packages the code actually uses
- Linters, scanners, and code generators keep running at exactly the same pinned versions
- Long-standing version-conflict workarounds are no longer needed and go away
- The dependency graph no longer drags in unrelated tooling internals
- Two live CVEs disappear because the package carrying them was never a real dependency
- No application behavior changes — this is dependency hygiene only
</summary>

<objective>
Delete `tools.go` and rewrite `go.mod` so tool CLIs are no longer module dependencies, because importing them pulls every tool's transitive dependency tree into this project — inflating `go.mod` to 448 lines, forcing five `replace` workarounds, and dragging in `github.com/go-git/go-git/v5` with two live CVEs (GHSA-hc8v-wwc9-vgxm, GHSA-qgq7-7hm3-q39j) that currently make `make precommit` fail at the `osv-scanner` target. Tool versions stay pinned via `tools.env` and `go run pkg@$(VERSION)`, which already works in the Makefile.
</objective>

<context>
Read the coding plugin's `docs/go-tools-versioning-guide.md` (in-container path: `/home/node/.claude/plugins/marketplaces/coding/docs/go-tools-versioning-guide.md`) — it is the canonical source for this migration. Requirement 5's pollution grep is taken verbatim from its "Migration Steps" §7. Read its "Pitfalls" section too, in particular that `go mod tidy -e` can truncate `go.mod`.

Read `CLAUDE.md` for project conventions.

Read `tools.go` — the file to delete. It imports 11 CLI tools under a `//go:build tools` tag.

Read `Makefile` — it already has `include tools.env` and already invokes every tool as `go run pkg@$(VERSION)`. **No Makefile change is needed.** The remaining `-mod=mod` usages (`go generate -mod=mod`, `go test -mod=mod`, `go vet -mod=mod`, `go list -mod=mod`) are correct and must stay.

Read `tools.env` — the pinned tool versions. Leave it untouched.

Read `go.mod` — note the `replace` block at the top and the direct-require list, which mixes four real application dependencies with ten tool-only ones.

This repo has **no `//go:generate` directives** and **no `github.com/bborbe/*` dependencies**, so there is no mock-regeneration step and no dependency-cascade bump to perform here.
</context>

<requirements>
1. Delete `tools.go` entirely.

2. Rewrite `go.mod` to a minimal form: keep `module github.com/bborbe/git-sync` and `go 1.26.5`, delete the entire `replace` block, and keep ONLY these four direct requires — each one is verified as genuinely imported by non-tool code:
   ```
   github.com/golang/glog
   github.com/onsi/ginkgo/v2
   github.com/onsi/gomega
   github.com/pkg/errors
   ```
   Keep each dep's existing version. Drop every tool-only direct require: `golangci-lint/v2`, `google/addlicense`, `osv-scanner/v2`, `goimports-reviser/v3`, `kisielk/errcheck`, `counterfeiter/v6`, `gosec/v2`, `segmentio/golines`, `shoenig/go-modtool`, `golang.org/x/vuln`.

   Delete the entire indirect require block — `go mod tidy` repopulates it.

3. Run `go mod tidy` to repopulate legitimate indirect requires. Do NOT use `go mod tidy -e` — it can truncate `go.mod` when resolution partially fails.

4. If `go mod tidy` surfaces any `github.com/bborbe/*` dependency (none exist today, but a transitive one could appear), bump it to `@latest` — an unbumped bborbe dep carrying its own `tools.go` re-drags the whole pollution cascade back in:
   ```
   grep '^	github.com/bborbe/' go.mod | awk '{print $1}' | xargs -I {} go get {}@latest
   ```
   If the grep matches nothing, this is a no-op — that is the expected outcome here.

5. Confirm zero tools.go-era pollution remains:
   ```
   grep -E '(cellbuf|go-header|go-diskfs|golangci-lint|osv-scanner|ginkgolinter|charmbracelet/x|denis-tingaikin)' go.mod
   ```
   This must return no matches. If any appear, run `go mod why <package>` to find what still pulls the old cascade.

6. Confirm `go-git` is gone entirely: `grep go-git go.mod` must return no matches. It was only ever reachable through the tool imports — no `.go` source file in this repo imports a git library.
</requirements>

<constraints>
- Do NOT commit — dark-factory handles git.
- Do NOT run `go mod vendor`, and never use `-mod=vendor` in any command. This repo has no `vendor/` directory. If any generic dependency-fix guidance you encounter suggests a `go mod vendor` step, ignore it — this constraint takes precedence.
- Do NOT change any application code. This change touches only `tools.go` (deleted), `go.mod`, and `go.sum`.
- Do NOT modify `tools.env` — no version changes, no reordering.
- Do NOT modify any tool invocation in the `Makefile` — every one is already correct.
- Do NOT edit `.osv-scanner.toml`. Some ignore entries will become unused once the pollution is gone; that is expected and is not an error.
- Existing tests must still pass unchanged.
- If one of the four retained direct deps in requirement 2 turns out to be unused, let `go mod tidy` drop it rather than forcing it back in.
</constraints>

<verification>
Run `make precommit` — must pass (exit 0). It currently FAILS at the `osv-scanner` target on the two go-git CVEs; removing the pollution is what makes it green, so a passing run is the primary signal this worked.

Then confirm the pollution is gone:

```
grep -E '(cellbuf|go-header|go-diskfs|golangci-lint|osv-scanner|ginkgolinter|charmbracelet/x|denis-tingaikin)' go.mod
grep go-git go.mod
ls tools.go
```

The two greps must produce no output, and `ls tools.go` must report that the file does not exist.

Confirm `go.mod` shrank substantially — `wc -l go.mod` should report well under 100 lines (it was 448).
</verification>
