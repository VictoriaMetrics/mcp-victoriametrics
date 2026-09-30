set -e
set -o pipefail

# check source code by linter
gofmt -l -w -s ./cmd
go vet ./cmd/...
which golangci-lint || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(go env GOPATH)/bin v2.14.0
golangci-lint run
