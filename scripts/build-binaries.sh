set -e
set -o pipefail

export BUILD_DATE=$(date -u +'%Y%m%d-%H%M%S')
export VERSION_FROM_COMMIT=$(git describe --long --all | tr '/' '-')
export VERSION_FROM_TAG=$(git tag -l --points-at HEAD)
export BUILD_VERSION=${VERSION_FROM_TAG:-$VERSION_FROM_COMMIT}
go build -ldflags "-X main.version='$BUILD_VERSION' -X main.date='$BUILD_DATE'" -o ./mcp-victoriametrics ./cmd/mcp-victoriametrics
