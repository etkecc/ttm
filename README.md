# ttm - a `time` To Matrix [![Matrix](https://img.shields.io/matrix/ttm:etke.cc?logo=matrix&server_fqdn=matrix.org&style=for-the-badge)](https://matrix.to/#/#ttm:etke.cc) [![Buy me a Coffee](https://shields.io/badge/donate-buy%20me%20a%20coffee-green?logo=buy-me-a-coffee&style=for-the-badge)](https://buymeacoffee.com/etkecc) [![coverage report](https://gitlab.com/etke.cc/ttm/badges/main/coverage.svg)](https://gitlab.com/etke.cc/ttm/-/commits/main) [![Go Report Card](https://goreportcard.com/badge/gitlab.com/etke.cc/ttm)](https://goreportcard.com/report/gitlab.com/etke.cc/ttm) [![Go Reference](https://pkg.go.dev/badge/gitlab.com/etke.cc/ttm.svg)](https://pkg.go.dev/gitlab.com/etke.cc/ttm)

A `time`-like command that will send end of an arbitrary command output and some other info (like exit status) to matrix room.

Consider this project a "bash-oneliner" to do some stuff fast. It's not supposed to be beautiful and shiny, it just must work.

## Features

* Run any arbitrary command with colors, tty window size and other fancy stuff
* Collect time information (real/user/sys, like `time command`)
* Send that info to the matrix room
* Project functionality considered final, so no breaking changes are possible, only fixes, tests and minor updates (like new options)

<details>
<summary>How it looks like</summary>

### you run command in terminal...

```bash
$ ttm ansible-playbook --with args
# ... scroll-scroll-scroll
PLAY RECAP *****************************************************************************************************************************
gitlab.com                    : ok=33   changed=0    unreachable=0    failed=0    skipped=147  rescued=0    ignored=0


real	15.166239745s
user	10.330419s
sys		2.213327s
```

### ...and get fancy html-formated message in matrix

**ttm report**

```bash
ansible-playbook --with args
```

```bash
# end of log (if configured)
PLAY RECAP *****************************************************************************************************************************
gitlab.com                    : ok=33   changed=0    unreachable=0    failed=0    skipped=147  rescued=0    ignored=0
```

```bash
real	15.166239745s
user	10.330419s
sys		2.213327s
```

Exit code: `0`

</details>

## How to get

* Arch Linux [AUR](https://aur.archlinux.org/packages/time-to-matrix-git/)
* or [Releases](https://gitlab.com/etke.cc/ttm/-/releases) for freebsd, linux and MacOS
* or `go install gitlab.com/etke.cc/ttm@latest`
* or from source code

## Quick start

```bash
# get the ttm app first, from the "How to get" section

# ttm configuration, check the "Configuration" section below to get the full (pretty impressive) list of available options
export TTM_HOMESERVER="https://matrix.example.com"
export TTM_TOKEN="matrix_access_token"
export TTM_ROOM="!XODRhTLplrymaFicdK:etke.cc"

# synax: ttm <any-command> [with any args]
ttm echo "that command will be sent to matrix, alongside with time stats, exit code and neat html formatting"
```

## Configuration

All ttm config options are env vars

### mandatory

* **TTM_HOMESERVER** - the real address of your matrix HS, not a delegated url, eg: `https://matrix.example.com`
* **TTM_ROOM** - matrix room id or alias, eg: `!XODRhTLplrymaFicdK:etke.cc` or `#ttm:etke.cc`

auth can be done with access token or login/password, that's up to you

**login/password** _if you want to login with matrix user login and password_

* **TTM_LOGIN** - matrix login (localpart), eg: `ttm`, not `@ttm:example.com`
* **TTM_PASSWORD** - matrix password, eg: `superSecurePAssword`

**access token** _if you want to login with existing session_

* **TTM_TOKEN** - matrix session access token

### optional options

* **TTM_LOG** - send full log information to matrix, allowed values: `0`, `1`; default: `0`
* **TTM_MSGTYPE** - message type, allowed values: `m.text`, `m.notice`; default: `m.text`
* **TTM_NOHTML** - do not send html-formatted message, only plaintext (increases allowed log size x2), allowed values: `0`, `1`; default: `0`
* **TTM_NOTEXT** - do not send plaintext message, only html (increases allowed log size x2), allowed values: `0`, `1`; default: `0`
* **TTM_NOTICEFAIL** - send message with `TTM_MSGTYPE="m.notice"` if exit code is not 0, allowed values: `0`, `1`; default: `0`
* **TTM_NOTIME** - do not send time information to matrix, allowed values: `0`, `1`; default: `0`

### deprecated

following env vars still available and work, but were replaced:

* TTM_ROOMID (use **TTM_ROOM** instead) - matrix room id or alias (eg: `!fsafaFSAsf:example.com` or `#ttm:etke.cc`)

## Additional info

ttm developed by [etke.cc](https://etke.cc), say hello in [#ttm:etke.cc](https://matrix.to/#/#ttm:etke.cc) matrix room (or ask questions, if you have any)
