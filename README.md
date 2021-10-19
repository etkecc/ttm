# ttm - a `time` To Matrix [![Buy me a Coffee](https://shields.io/badge/donate-buy%20me%20a%20coffee-green?logo=buy-me-a-coffee&style=for-the-badge)](https://buymeacoffee.com/etkecc) [![coverage report](https://gitlab.com/etke.cc/ttm/badges/main/coverage.svg)](https://gitlab.com/etke.cc/ttm/-/commits/main) [![Go Report Card](https://goreportcard.com/badge/gitlab.com/etke.cc/ttm)](https://goreportcard.com/report/gitlab.com/etke.cc/ttm) [![godocs.io](http://godocs.io/gitlab.com/etke.cc/ttm?status.svg)](http://godocs.io/gitlab.com/etke.cc/ttm)

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
* **TTM_LOGIN** - matrix login (localpart) (eg: `ttm`, not `@ttm:example.com`)
* **TTM_PASSWORD** - matrix password
* **TTM_ROOMID** - matrix room id (eg: `!fsafaFSAsf:example.com`)
* **TTM_NOTIME** - do not send time information to matrix, default: false
* **TTM_LOG** - send full log information to matrix, default: false

## How to get

1. [Releases](https://gitlab.com/etke.cc/ttm/-/releases) for freebsd, linux and MacOS
2. or `go install gitlab.com/etke.cc/ttm@latest`
2. or from source code
