# Actions Exporter

Prometheus exporter for GitHub Actions workflow states and API rate limit usage.

## GitHub PAT Scopes

This app requires a classic Personal Access Token with the following scope:

- **`repo`** — needed to list organization repositories and their Actions workflows (includes private repos)

If you only need to monitor public repositories, `public_repo` is sufficient.
