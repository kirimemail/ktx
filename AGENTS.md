# ktx CLI Agent Guide

This document describes how to update and maintain the `ktx` CLI.

## Project Structure

```
ktx/
├── main.go              # CLI entry point, argument parsing
├── cmd/                 # Command implementations
│   ├── domains.go
│   ├── credentials.go
│   ├── send.go
│   ├── validate.go
│   ├── logs.go
│   ├── suppressions.go
│   ├── webhooks.go
│   ├── quota.go
│   └── cmd_test.go      # Unit tests
├── .github/workflows/
│   └── release.yml      # Multi-platform release build
├── install.sh           # Installation script
├── README.md            # User documentation
├── SKILLS.md            # User quick reference
└── go.mod / go.sum      # Dependencies
```

## Adding a New Command

1. Create new file in `cmd/` (e.g., `cmd/newcmd.go`):

```go
package cmd

import (
    "fmt"

    smtpsdk "github.com/kirimemail/kirimemail-smtp-go-sdk"
)

func NewCmdCmd(client *smtpsdk.Client, defaultDomain string, args []string) error {
    // Implementation
    return nil
}
```

2. Add case in `main.go`:

```go
case "newcmd":
    return cmd.NewCmdCmd(client, *domain, subArgs)
```

3. Add tests in `cmd/cmd_test.go`:

```go
func TestNewCmdCmd(t *testing.T) {
    client := smtpsdk.NewClient("testuser", "testtoken")
    err := NewCmdCmd(client, "example.com", []string{})
    // assertions
}
```

4. Update help text in `printUsage()` and README.md

## Updating SDK Version

If the KirimEmail SDK is updated:

```bash
# Update to latest
go get github.com/kirimemail/kirimemail-smtp-go-sdk@latest

# Or specific version/commit
go get github.com/kirimemail/kirimemail-smtp-go-sdk@<commit-hash>

# Tidy dependencies
go mod tidy

# Verify build
go build -o ktx .
go test ./...
```

## Releasing New Version

1. Ensure all tests pass:
```bash
go test ./...
```

2. Update version if needed (no formal versioning file)

3. Tag and push:
```bash
git add .
git commit -m "Description of changes"
git tag v1.x.x
git push origin main
git push origin v1.x.x
```

4. GitHub Actions will:
   - Build for linux/darwin/windows × amd64/arm64
   - Create checksums
   - Create GitHub release with artifacts

## Testing Changes Locally

```bash
# Build
go build -o ktx .

# Test with fake credentials (will fail at API call)
./ktx domains list -username test -token test

# With real credentials
KIRIM_USERNAME=user KIRIM_TOKEN=token ./ktx quota
```

## Release Artifacts

After tagging, the following artifacts are created:

| File | Platform |
|------|----------|
| ktx-linux-amd64 | Linux x86_64 |
| ktx-linux-arm64 | Linux ARM64 |
| ktx-darwin-amd64 | macOS Intel |
| ktx-darwin-arm64 | macOS Apple Silicon |
| ktx-windows-amd64.exe | Windows x86_64 |
| ktx-windows-arm64.exe | Windows ARM64 |

Each has a corresponding `.sha256sum` checksum file.

## Common Issues

### Build fails with SDK import
```bash
go mod tidy
go get github.com/kirimemail/kirimemail-smtp-go-sdk@latest
```

### Tests fail
Check SDK types match the implementation. Run:
```bash
go vet ./...
go build ./...
```

### Release artifact missing
Check `.github/workflows/release.yml` artifact names match `ktx-*` pattern.

## Workflow Updates

To modify the release process, edit `.github/workflows/release.yml`:

- Add/remove platforms in the matrix
- Change artifact naming
- Modify upload patterns
