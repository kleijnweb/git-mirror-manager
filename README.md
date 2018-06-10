# Git Mirror Manager

Dead-simple service for managing Git Mirrors, written in Go.

# TODO

Tests, tests, tests
Lock git operations using mutexes
Better logging
Recover from panics
Create zip files

## Features

Exposes a super simple API to add and delete git mirrors. Updates them periodically using cron syntax, and builds zip files of tags.

### Prerequisites

Requires the git binaries to be installed. Make sure git can find a private key to clone the upstreams.

### API

Add mirror:

```
POST /repo
git@github.com/some/repo-name.git
```

Note that the client is expected to wait for a quick test using `git ls-remote`. The clone is done outside of the request/response scope.

Remove mirror:

```
DELETE /repo/some/repo-name
```

Returns a 400 if the repo doesn't exist. Any `.git` suffix is stripped.

Health-check:

```
GET /ping
```

Everything else returns a 404 or 405, unless some processing failed, in which case, an empty 500. Check the logs for details.

### Logging

Dumps everything it does or went wrong to STDOUT and STDERR respectively.

### Persistence

There is no extra persistence, config files, or the like. On boot, the root mirror directory is scanned for Git repositories.

### Limitations

Plenty, but notably:

 - Since the names or inferred by the repo URI but do not contain a hostname, on collisions POSTs will be rejected.
 - Names are limited to `{namespace}/{name}`
 

## Configuration

| Env Name  |  Default |  Description |
|---|---|---|
|  `GIT_MIRROR_UPDATE_INTERVAL` |  `0 * * * *` |  update frequency using cron notation |
|  `GIT_MIRROR_MANAGER_ADDR` |  `:8080` |  API bind address |
|  `GIT_MIRROR_BASEDIR` |  `/opt/data/mirrors` |  where git mirrors repositories are cloned to |
|  `GIT_MIRROR_DISTDIR` |  `/opt/data/dist` |  where zip files are written to |

## Running

Use the included Dockerfile, mount a private key under `/home/manager/.ssh`, create a data volume.
