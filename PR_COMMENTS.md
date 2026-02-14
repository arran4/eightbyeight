# Issues Found with go-subcommand v0.0.17

## 1. Generated code requires Go 1.16+ due to `go:embed`
**Description:**
The generated `cmd/eightbyeight/templates/templates.go` file uses the `go:embed` directive, which was introduced in Go 1.16. The project `go.mod` specified `go 1.14`, causing compilation errors when running the generated CLI.

**Steps to Reproduce:**
1.  Run `go generate` using `go-subcommand` v0.0.17.
2.  Try to build or run the generated CLI (`go build ./cmd/eightbyeight`) in a project targeting Go < 1.16.
3.  Observe error: `go:embed requires go1.16 or later (-lang was set to go1.14; check go.mod)`

**Impact:**
Breaking change for projects targeting Go versions older than 1.16. The generated code is not backward compatible with the project's original Go version (1.14).

**Possible Solutions:**
-   **Upgrade Project:** Upgrade `go.mod` to `go 1.16` (This was done in this PR to proceed).
-   **Conditional Generation:** Update `go-subcommand` to detect the Go version in `go.mod` and generate code compatible with it (e.g., using string constants instead of `embed` for older versions), or fail with a clear error message.
-   **Documentation:** Explicitly document that `go-subcommand` v0.0.17+ requires Go 1.16+ for generated code.

## 2. Regression Tests Failing (Environment/Pre-existing)
**Description:**
`TestReproducePatterns` in `repro_test.go` fails with significant pixel differences when comparing generated images against legacy BMPs in `exampledata/`. This failure persists regardless of the Go version (1.14 or 1.16) or the CLI changes.

**Status:**
This issue appears to be environment-related (likely font rendering differences) or pre-existing, and is unrelated to the `go-subcommand` integration.
