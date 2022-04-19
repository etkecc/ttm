# ttm v1.4.1 (2021-11-26)

### Misc :zzz:

* added human-readable changelog
* rewrite the README
* re-configure goreleaser to exclude archives and checksums (generated by gitlab anyway)
* updated Makefile to pass `GOFLAGS` and set the value only if nothing is set in env
* added version information to the help message
* change released package name to `time-to-matrix` to avoid name clashes (project and binary name did not change)

> **NOTE** the package name change isn't marked as breaking change because this project is not present in any package repository at the moment of that release.

# ttm v1.4.0 (2021-11-15)

### Features :sparkles:

* `TTM_ROOM` and automatic room alias resolving (`TTM_ROOMID` is deprecated, but works the same as the new option)
* Arch Linux AUR repo
* Help message

### Misc :zzz:

* refactoring
* tests
* updating scripts and configs with Arch Linux guidelines

# ttm v1.3.0 (2021-11-12)

### Features :sparkles:

* `TTM_MSGTYPE` to send message with different type (`m.text` or `m.notice`)
* `TTM_NOTEXT` to skip plaintext message and send html-formatted only
* `TTM_NOTICEFAIL` to override msgtype to `m.notice` automatically if exit code != 0

# ttm v1.2.0 (2021-11-06)

### Features :sparkles:

* `TTM_NOHTML` to skip html-formatted message and send plaintext only
* `TTM_TOKEN` to use matrix acces token / session, i.e. for SSO envs
* automatic log truncation if message is too big

### Misc :zzz:

* unit tests

# ttm v1.1.0 (2021-10-19)

### Features :sparkles:

* `TTM_NOTIME` to skip time output
* `TTM_LOG` to send full log in message
* binary releases for all major platform and architectures

# ttm v1.0.0 (2021-10-18)

### Features :sparkles:

* send end of an arbitrary command output and some other info (like exit status) to matrix room.

### Bugfixes :bug:

N/A _that's the first version, bugs are creating on that step!_

### Breaking changes :warning:

ttm has been created.