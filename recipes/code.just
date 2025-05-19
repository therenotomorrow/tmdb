GOLANGCI_LINT_VERSION := 'v2.1.6'
GOLANGCI_LINT_PATH := BIN / 'golangci-lint'
GOLANGCI_LINT := GOLANGCI_LINT_PATH + '@' + GOLANGCI_LINT_VERSION

[private]
install-golangci-lint:
    curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b {{BIN}} {{GOLANGCI_LINT_VERSION}}
    mv {{GOLANGCI_LINT_PATH}} {{GOLANGCI_LINT}}

[doc('Lint the codebase and check for code integrity issues')]
lint:
    @if test ! -e {{GOLANGCI_LINT}}; then just install-golangci-lint; fi
    {{GOLANGCI_LINT}} run ./...
