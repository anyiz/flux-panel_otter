# Flux Otter 1.4.3-otter.2

The panel backend now generates ingress GOST installation commands that download `install.sh` from `anyiz/flux-panel_otter` Release `1.4.3-otter.2`. The matching installer downloads GOST binaries from the same release.

Matching Docker Hub images:

- `anyi662/springboot-backend:1.4.3-otter.2`
- `anyi662/vite-frontend:1.4.3-otter.2`

The frontend code is unchanged from `1.4.3-otter.1`; the `.2` frontend tag points to the same image content. The backend image is rebuilt with the corrected command.

Existing deployments must update the backend container to this tag for the panel-generated command to change. Keep the existing MySQL volume. The one-time `20261111_proxy_protocol.sql` migration remains required for databases created before the Otter release; do not apply it twice. Keep PROXY Protocol off until the ingress and landing nodes are ready.
