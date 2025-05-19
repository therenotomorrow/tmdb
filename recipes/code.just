# ---- golangci-lint

GOLANGCI_LINT_VERSION := 'v2.1.6'
GOLANGCI_LINT_PATH := BIN / 'golangci-lint'
GOLANGCI_LINT := GOLANGCI_LINT_PATH + '@' + GOLANGCI_LINT_VERSION

[private]
install-golangci-lint:
    curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b {{ BIN }} {{ GOLANGCI_LINT_VERSION }}
    mv {{ GOLANGCI_LINT_PATH }} {{ GOLANGCI_LINT }}

[group('code')]
[doc('Run static analysis using `golangci-lint` to detect code issues')]
lint:
    @if test ! -e {{ GOLANGCI_LINT }}; then just install-golangci-lint; fi
    {{ GOLANGCI_LINT }} run ./...

# ---- mockery

MOCKERY_VERSION := 'v3.2.5'
MOCKERY_PATH := BIN / 'mockery'
MOCKERY := MOCKERY_PATH + '@' + MOCKERY_VERSION

[private]
install-mockery:
    GOBIN={{ BIN }} go install github.com/vektra/mockery/v3@{{ MOCKERY_VERSION }}
    mv {{ MOCKERY_PATH }} {{ MOCKERY }}

[group('code')]
[doc('Generate Go interfaces mocks using `mockery`')]
mock:
    @if test ! -e {{ MOCKERY }}; then just install-mockery; fi
    {{ MOCKERY }}

# ---- fieldaligment

FIELDALIGNMENT_VERSION := 'v0.33.0'
FIELDALIGNMENT_PATH := BIN / 'fieldalignment'
FIELDALIGNMENT := FIELDALIGNMENT_PATH + '@' + FIELDALIGNMENT_VERSION

[private]
install-fieldaligment:
    GOBIN={{ BIN }} go install golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment@{{ FIELDALIGNMENT_VERSION }}
    mv {{ FIELDALIGNMENT_PATH }} {{ FIELDALIGNMENT }}

[group('code')]
[doc('Reorder struct fields to improve memory layout and reduce padding')]
align:
    @if test ! -e {{ FIELDALIGNMENT }}; then just install-fieldaligment; fi
    {{ FIELDALIGNMENT }} --fix ./...

# ---- build

BUILD_PATH := BIN / 'tmdb'
BUILD_MAIN := invocation_directory() / 'cmd' / 'tmdb' / 'main.go'

[group('code')]
[doc('Build the tmdb binary with an optional `TMDB_TOKEN` embedded at compile time')]
build token='':
    go build -o {{ BUILD_PATH }} -ldflags "-X 'main.TMDBToken={{ token }}'" {{ BUILD_MAIN }}

# ---- default

[group('code')]
[doc('Run all code quality tools: mock generation, struct alignment, and static analysis')]
code: mock align lint
