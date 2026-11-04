# as-gitverse-api

🇷🇺 [Русская версия](README_RUS.md)

A CLI for bulk management of [Gitverse](https://gitverse.ru) repositories through the
public API. Authentication uses a Personal Access Token.

## Installation

```
go build -o gitverse ./cmd/gitverse
```

## Token

Precedence: the `--token` flag > the `GITVERSE_TOKEN` environment variable > the credentials file
`~/Library/Application Support/gitverse/credentials` (macOS) /
`~/.config/gitverse/credentials` (Linux), permissions 0600.

```
export GITVERSE_TOKEN=xxxxxxxxxxxxxxxx
```

## Global flags

| Flag | Description |
|---|---|
| `--token` | PAT (overrides `GITVERSE_TOKEN` and the credentials file) |
| `--json` | output as JSON instead of a table |
| `--concurrency` | max parallel requests in bulk operations (default 5) |
| `--verbose` | verbose logging (the token never appears in the log) |

## Commands

### `repo` — repositories

| Command | Description |
|---|---|
| `repo create --names/--file [--private]` | create one or more repositories (bulk) |
| `repo list [--org <org>]` | list repositories of the current user or an organization |
| `repo collaborators list --owner --repo` | list the collaborators of a repository |
| `repo collaborators add --owner --repo --names/--file --permission` | add collaborators (bulk; `--permission`: `read`/`write`/`admin`) |
| `repo commit-stats --owner --repo [--since] [--until]` | commit statistics per author |

```
# create several public repositories in one call
gitverse repo create --names foo,bar,baz

# create private repositories from a file (one name per line)
gitverse repo create --file repos.txt --private

# list your own repositories
gitverse repo list

# list the repositories of an organization
gitverse repo list --org studentlabs

# list the collaborators of a repository
gitverse repo collaborators list --owner studentlabs --repo slide_backend

# add collaborators with write permission
gitverse repo collaborators add --owner studentlabs --repo slide_backend --names alice,bob --permission write

# commit statistics per author
gitverse repo commit-stats --owner studentlabs --repo slide_backend

# commit statistics for a period (--since/--until are dates in YYYY-MM-DD format, --until is inclusive)
gitverse repo commit-stats --owner studentlabs --repo slide_backend --since 2026-08-01 --until 2026-08-28
```

## JSON output

Any command can print JSON instead of a table:

```
gitverse repo list --org studentlabs --json
```

## Intentionally not implemented

- Any `delete` commands for repositories — forbidden by a hard project rule
  (deletion is done manually through the Gitverse UI only).
- `collaborators remove` and `org members` (a full listing of organization
  members) — checked against the public API: these endpoints do not exist
  (`GET /orgs/{org}/members` without a username → `400`; only
  `GET /orgs/{org}/members/{username}` exists, which checks the membership of a single user).
- Webhooks (`webhook`) and organization teams (`team`) — excluded from the project
  scope, not implemented.
