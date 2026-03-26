# Definition of Done

After completing your implementation, review your own changes against each criterion below. These are quality checks you perform by inspecting your work — not commands to run (linting and tests already ran via `validationCommand`). Report any unmet criterion as a blocker.

## Code Quality

- Exported types, functions, and interfaces have doc comments
- Error handling uses proper context wrapping
- No debug output (print statements, fmt.Printf) — use glog for logging
- No `os.Exit` in library functions — only in `main()`
- Follow idiomatic Go patterns

## Testing

- New code has good test coverage (target >= 80%)
- Changes to existing code have tests covering at least the changed behavior
- Tests use Ginkgo v2 / Gomega with Counterfeiter mocks

## Install

- `go install github.com/bborbe/git-sync@latest` works
- No `exclude` or `replace` directives in go.mod (break remote install)

## Documentation

- README.md is updated if the change affects usage, configuration, or setup
- CHANGELOG.md has an entry under `## Unreleased`
