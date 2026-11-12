# Publish and deploy the modified Flux panel

## Published images (2026-11-12)

The `linux/amd64` images from this fork were built, smoke-tested, pushed, and verified on Docker Hub:

- `anyi662/springboot-backend:1.4.3-otter.1` (index digest `sha256:60039ca65720c303d2d72264bd95bcc1413ebf588202dbe496d7190c3bca1d6d`)
- `anyi662/vite-frontend:1.4.3-otter.1` (index digest `sha256:4160675a72234518a9495408357939a936a16dfc2e65ab010099e6494f336f8c`)

For this release, set `FLUX_DOCKERHUB_USER=anyi662` and `FLUX_IMAGE_TAG=1.4.3-otter.1` in the production `.env`. The production Compose example with IPv6 is `docker-compose-v6.yml`.

The backend and frontend already have Dockerfiles. This fork's `docker-compose-v4.yml` and `docker-compose-v6.yml` now read `FLUX_DOCKERHUB_USER` and `FLUX_IMAGE_TAG` to select your images. The MySQL image, volume, ports, and network definitions are unchanged. The ingress gost process is separate and needs its own upgrade before PROXY Protocol is enabled.

From a Linux host with Docker Buildx, log in to Docker Hub and build/push both images from this exact fork:

```bash
docker login
bash scripts/publish-panel-images.sh YOUR_DOCKERHUB_USER 1.4.3-otter.1 linux/amd64
```

If the production host is ARM64, pass `linux/arm64`, or `linux/amd64,linux/arm64` for both. The two repositories are `YOUR_DOCKERHUB_USER/springboot-backend` and `YOUR_DOCKERHUB_USER/vite-frontend`. Use a new, immutable tag for every release. Docker Hub authentication and repository access are required; the script does not store credentials.

On the production host, keep the existing Compose project/directory and `.env` values, copy in this fork's matching Compose file and migration SQL, then add:

```dotenv
FLUX_DOCKERHUB_USER=YOUR_DOCKERHUB_USER
FLUX_IMAGE_TAG=1.4.3-otter.1
```

Use the same Compose file and project name as the current deployment, so Compose keeps the existing `mysql_data` volume. For an IPv6 deployment replace `docker-compose-v4.yml` with `docker-compose-v6.yml` below. Back up MySQL before the one-time schema migration:

```bash
docker compose -f docker-compose-v4.yml exec -T mysql sh -c \
  'exec mysqldump -uroot -p"$MYSQL_ROOT_PASSWORD" "$MYSQL_DATABASE"' \
  > flux-before-proxy-protocol.sql
docker compose -f docker-compose-v4.yml exec -T mysql sh -c \
  'exec mysql -uroot -p"$MYSQL_ROOT_PASSWORD" "$MYSQL_DATABASE"' \
  < doc/migrations/20261111_proxy_protocol.sql
docker compose -f docker-compose-v4.yml config --images
docker compose -f docker-compose-v4.yml pull backend frontend
docker compose -f docker-compose-v4.yml up -d --no-deps backend frontend
docker compose -f docker-compose-v4.yml ps
docker compose -f docker-compose-v4.yml logs --tail=100 backend frontend
```

The migration is necessary for an existing `mysql_data` volume; mounting the updated `gost.sql` only initializes a fresh database. Do not run the migration twice. Do not use `down -v`. Keep PROXY Protocol off until the entry gost and landing nodes are upgraded and tested; see the Xboard panel's `docs/strict-device-leases.md` for the full rollout sequence.
