# Flux Otter 1.4.3-otter.1

This release packages this fork's Linux GOST agent, panel installers, Compose files, and database schema. The GOST binaries are built from the tagged `go-gost` source for AMD64 and ARM64; `SHA256SUMS` lists their hashes.

The matching panel images are available on Docker Hub:

- `anyi662/springboot-backend:1.4.3-otter.1`
- `anyi662/vite-frontend:1.4.3-otter.1`

The node installer selects `gost-amd64` or `gost-arm64` for the local architecture and downloads it from this release. The panel installer downloads this release's Compose file and `gost.sql`, then selects the matching Docker images.

For an existing panel database, back up MySQL and apply the one-time `20261111_proxy_protocol.sql` migration before enabling PROXY Protocol. Upgrade the panel, relay GOST agents, and landing Xboard nodes as described in `doc/dockerhub-upgrade.md` and the Xboard panel's `docs/strict-device-leases.md`. Leave PROXY Protocol disabled until the relay and landing nodes are ready.
