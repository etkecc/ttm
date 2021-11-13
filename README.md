# ttm - a `time` To Matrix [![Matrix](https://img.shields.io/matrix/ttm:etke.cc?logo=matrix&server_fqdn=matrix.org&style=for-the-badge)](https://matrix.to/#/#ttm:etke.cc) [![Buy me a Coffee](https://shields.io/badge/donate-buy%20me%20a%20coffee-green?logo=buy-me-a-coffee&style=for-the-badge)](https://buymeacoffee.com/etkecc) [![coverage report](https://gitlab.com/etke.cc/ttm/badges/main/coverage.svg)](https://gitlab.com/etke.cc/ttm/-/commits/main) [![Go Report Card](https://goreportcard.com/badge/gitlab.com/etke.cc/ttm)](https://goreportcard.com/report/gitlab.com/etke.cc/ttm) [![Go Reference](https://pkg.go.dev/badge/gitlab.com/etke.cc/ttm.svg)](https://pkg.go.dev/gitlab.com/etke.cc/ttm)

A `time`-like command that will send end of an arbitrary command output and some other info (like exit status) to matrix room.

Consider this project a "bash-oneliner" to do some stuff fast. It's not supposed to be beautiful and shiny, it just must work.

## Features

* Run any arbitrary command with colors, tty window size and other fancy stuff
* Collect time information (real/user/sys, like `time command`)
* Send that info to the matrix room

## How it looks like

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


## Stability and project future

* Project functionality considered final
* No breaking changes
* Only bug fixes, tests and minor internal changes are possible

So, feel free to use `latest`, it works. It will work. That will not change.

## TODO

* unit tests

## Configuration

done via env vars:

* **TTM_HOMESERVER** - the real address of your matrix HS, not a delegated url (eg: `https://matrix.example.com`)

* **TTM_LOGIN** - _(only for password auth)_ matrix login (localpart) (eg: `ttm`, not `@ttm:example.com`)
* **TTM_PASSWORD** - _(only for password auth)_ matrix password
* **TTM_TOKEN** - _(only for access token)_ matrix session access token

* **TTM_ROOM** - matrix room id or alias (eg: `!fsafaFSAsf:example.com` or `#ttm:etke.cc`)
* **TTM_MSGTYPE** - message type, default: `m.text`, can be set to `m.notice`

* **TTM_NOTIME** - do not send time information to matrix, default: false
* **TTM_NOHTML** - do not send html-formatted message, only plaintext (increases allowed log size x2), default: `0`
* **TTM_NOTEXT** - do not send plaintext message, only html (increases allowed log size x2), default: `0`
* **TTM_NOTICEFAIL** - send message with `TTM_MSGTYPE="m.notice"` if exit code is not 0, default: `0`
* **TTM_LOG** - send full log information to matrix, default: `0`

_following env vars still available and work, but were replaced_:

* TTM_ROOMID (use **TTM_ROOM** instead) - matrix room id or alias (eg: `!fsafaFSAsf:example.com` or `#ttm:etke.cc`)

### Examples

**minimal** (access token)

```bash
# auth
export TTM_HOMESERVER=https://matrix.example.com
export TTM_TOKEN=your_access_token

# room
export TTM_ROOMID=!ttmroom:example.com
```

**login/password auth with full log and without time info**

```bash
# auth
export TTM_HOMESERVER=https://matrix.example.com
export TTM_LOGIN=ttm
export TTM_PASSWORD=thatsecure

# room
export TTM_ROOMID=!ttmroom:example.com

# options
export TTM_NOTIME=1
export TTM_LOG=1
```

**access token auth without html formatting**

```bash
# auth
export TTM_HOMESERVER=https://matrix.example.com
export TTM_TOKEN=your_access_token

# room
export TTM_ROOMID=!ttmroom:example.com

# options
export TTM_NOHTML=1
```

## How to get

1. Arch Linux [AUR](https://aur.archlinux.org/packages/time-to-matrix-git/)
2. or [Releases](https://gitlab.com/etke.cc/ttm/-/releases) for freebsd, linux and MacOS
3. or `go install gitlab.com/etke.cc/ttm@latest`
4. or from source code
