#!/usr/bin/env bash
set -euo pipefail

# Run after `docker login` from the repository root.
namespace="${1:?Usage: scripts/publish-panel-images.sh DOCKERHUB_USER [TAG] [PLATFORM]}"
tag="${2:-1.4.3-otter.2}"
platform="${3:-linux/amd64}"
cd "$(dirname "$0")/.."

docker buildx build --platform "$platform" \
  --tag "$namespace/springboot-backend:$tag" --push ./springboot-backend
docker buildx build --platform "$platform" \
  --tag "$namespace/vite-frontend:$tag" --push ./vite-frontend

printf 'Published %s/springboot-backend:%s and %s/vite-frontend:%s (%s)\n' \
  "$namespace" "$tag" "$namespace" "$tag" "$platform"
