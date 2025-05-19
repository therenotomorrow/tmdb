[private]
smoke:
    go test ./...

[private]
cover:
    go test -race -coverprofile=coverage.out $(find . -name '*.go' ! -name '*_gen.go' ! -name 'main.go' -exec dirname {} \; | sort -u)
    go tool cover -func coverage.out

[group('test')]
[doc('Run tests by type: `smoke` for quick checks (default), `cover` for detailed analysis')]
test type='smoke':
    just {{ type }}
