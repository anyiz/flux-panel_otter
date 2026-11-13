# Flux Otter 1.4.3-otter.3

Fixes TCP real-IP forwarding that appeared enabled in the panel but emitted no PROXY header. The Java backend sent a JSON number, which Go decoded as float64; the bundled GOST metadata reader did not recognize that type. The backend now sends the version as a string, compatible with existing Otter agents. The database and public API continue to use integer values 0/1/2.

Two regression tests cover serialized v1/v2 settings for direct and tunnel TCP forwarding, disabled rules, and UDP exclusion. The release workflow runs these tests before publishing binaries.

Docker Hub images (linux/amd64):

- `anyi662/springboot-backend:1.4.3-otter.3`
- `anyi662/vite-frontend:1.4.3-otter.3`

The frontend is unchanged and reuses the .2 image content. GOST source is unchanged; existing Otter .2 agents do not need a restart for this backend fix. Install scripts, Compose defaults, and the panel-generated agent command point to this release.

Upgrade the backend, then re-save only the intended PROXY-enabled test forwards so they receive the corrected metadata. Do not enable PROXY on unrelated production forwards. No new database migration is required; retain the existing MySQL volume. Existing Xboard lease configuration and test allowlists are outside this release and must remain unchanged.

The release includes AMD64/ARM64 GOST binaries, installers, Compose samples, SQL files and checksums. The panel Docker images are AMD64 only.
