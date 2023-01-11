#!/bin/sh

set -e

DIR=$(dirname "$BASH_SOURCE")
echo "Getting into directory ${DIR}/dapr"
cd "${DIR}/dapr"

CGO_ENABLED=0 \
GOOS=linux \
  go build \
  -ldflags="-X github.com/dapr/dapr/pkg/buildinfo.features=AppHealthCheck -X github.com/dapr/dapr/pkg/buildinfo.gitversion=v1-dirty -X github.com/dapr/dapr/pkg/buildinfo.version=dev -X github.com/dapr/kit/logger.DaprVersion=dev -s -w" \
  -o ./dist/ \
  ./cmd/daprd/;
echo "Built dapr to ./dist/daprd"

docker build \
  --build-arg PKG_FILES=daprd \
  -f docker/Dockerfile \
  -t dapr-firewall:latest \
  ./dist/
echo "Built image dapr-firewall:latest"
