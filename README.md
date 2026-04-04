# Actions Exporter

Prometheus exporter for GitHub Actions workflow states and API rate limit usage.

Tired of GitHub disabling your scheduled workflows after 60 days of repository inactivity, and only randomly discovering it some time later? This exporter was written for you: a timeseries is added for each of your workflows that's disabled, with a label that informs whether it was disabled manually (by you) or disabled due to inactivity (by GitHub.)

Supports monitoring workflows for GitHub **organizations** (`--github-org`) and/or **users** (`--github-user`). At least one must be provided; both can be used simultaneously.

Example:

```text
# HELP github_rate_limit_limit The maximum number of requests allowed per hour (x-ratelimit-limit).
# TYPE github_rate_limit_limit gauge
github_rate_limit_limit 5000

# HELP github_rate_limit_used The number of requests used in the current rate limit window (x-ratelimit-used).
# TYPE github_rate_limit_used gauge
github_rate_limit_used 96

# HELP github_workflow_state Shows non-active workflow state for workflows belonging to a GitHub user or organization.
# TYPE github_workflow_state gauge
github_workflow_state{owner="FooUser",repository="BarRepository",state="disabled_manually",workflow="Lint + Test"} 1
```

## GitHub PAT Scopes

This app requires a classic Personal Access Token with the following scope:

- **`repo`** — needed to list organization/user repositories and their Actions workflows (includes private repos)

If you only need to monitor public repositories, `public_repo` is sufficient.
