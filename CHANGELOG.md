# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.8.0] - 2026-11-05

### Added

- `?` toggles full keyboard help, including `home`/`g` and `end`/`G` for
  playlist navigation. The playlist shrinks to make room for expanded help.

## [0.7.0] - 2026-11-04

### Added

- Volume is saved on exit and restored before playback on the next launch.
  Preferences are stored in `$XDG_CONFIG_HOME/tplay/state.json`, defaulting to
  `~/.config/tplay/state.json`. Without valid saved state, volume starts at
  100%; read and write failures warn without preventing playback.

## [0.6.0] - 2026-11-03

### Added

- `make uninstall` removes the installed `tplay` binary without removing its
  directory or other files.
- `DESTDIR` support for staging source installations when packaging.

### Changed

- `make install` builds locally and installs with mode `755` into `$(PREFIX)/bin`.
  The default `PREFIX` is now `/usr/local`; use `PREFIX="$HOME/.local"` for the
  previous per-user location. Install and uninstall support paths with spaces.

## [0.5.1] - 2026-11-02

### Changed

- Expanded CLI help to list `--help`/`-h` and `--version` and explain that
  omitting the target plays the current directory.
- Documented help and version commands, build version sources, and the
  `VERSION` override for source builds.

## [0.5.0] - 2026-11-01

### Added

- `r` cycles repeat off, all and one. Repeat-all wraps at playlist boundaries;
  repeat-one replays the current track when it finishes, while `n`/`p` still
  navigate between tracks. Failed files are skipped, and playback stops if
  every reachable track fails, even with repeat enabled.
- `s` toggles shuffle, playing each track once per pass without reordering the
  displayed list. Previous retraces the shuffled order within a pass;
  repeat-all starts a fresh shuffled pass. Toggling shuffle keeps the current
  track playing.
- The playlist header shows repeat and shuffle modes; both are off by default.

## [0.4.0] - 2026-10-30

### Added

- `tplay -r [directory]` includes audio files from all subdirectories, sorted
  naturally by path. Without a directory argument, it scans the current
  directory. Recursive scanning is disabled by default; opening an individual
  file keeps the playlist limited to its siblings.
- `tplay --help` prints usage and describes the `-r` option.

## [0.3.0] - 2026-10-29

### Added

- `install.sh` installs the latest GitHub release into `~/.local/bin` on Linux
  and macOS (amd64 and arm64) via `curl | sh`, verifies SHA-256 checksums, and
  warns if `ffmpeg`/`ffprobe` or the install directory in `PATH` is missing.

### Changed

- CI checks the installer with offline integration tests on pushes to `main`
  and on pull requests.

## [0.2.0] - 2026-10-28

### Added

- Prebuilt binaries for Linux and macOS (amd64 and arm64) on GitHub Releases,
  built by GoReleaser when a `v*` tag is pushed; the release notes are taken
  from this changelog.
- `tplay` without an argument plays the current directory instead of printing
  the usage.

## [0.1.1] - 2026-10-28

### Changed

- The module moved to `github.com/dzhenetl/tplay` after the GitHub username
  change; install with `go install github.com/dzhenetl/tplay/cmd/tplay@latest`.

## [0.1.0] - 2026-10-26

### Added

- Playback of any audio format `ffmpeg` can decode; `tplay <file or directory>`
  starts with the given file or the first file of the directory.
- Playlist of every audio file in the same directory, sorted by name with
  numbers in natural order (`2` before `10`).
- Keys: `space` pause, `←`/`→` (`h`/`l`) seek 5 seconds, `n`/`p` next and
  previous track (`p` restarts the track after 3 seconds), `↑`/`↓` (`k`/`j`,
  `g`/`G`) move in the playlist, `enter` play the selection, `+`/`-` volume,
  `q`/`esc` quit.
- The next track starts when one ends; a file that fails to play is skipped in
  the direction you last moved.
- `tplay --version` prints the version embedded at build time.
- Install with `go install github.com/dzhenetl/tplay/cmd/tplay@latest`.
- CI on GitHub Actions runs `make lint` and `make test` on pushes to `main` and
  on pull requests.

[0.8.0]: https://github.com/dzhenetl/tplay/compare/v0.7.0...v0.8.0
[0.7.0]: https://github.com/dzhenetl/tplay/compare/v0.6.0...v0.7.0
[0.6.0]: https://github.com/dzhenetl/tplay/compare/v0.5.1...v0.6.0
[0.5.1]: https://github.com/dzhenetl/tplay/compare/v0.5.0...v0.5.1
[0.5.0]: https://github.com/dzhenetl/tplay/compare/v0.4.0...v0.5.0
[0.4.0]: https://github.com/dzhenetl/tplay/compare/v0.3.0...v0.4.0
[0.3.0]: https://github.com/dzhenetl/tplay/compare/v0.2.0...v0.3.0
[0.2.0]: https://github.com/dzhenetl/tplay/compare/v0.1.1...v0.2.0
[0.1.1]: https://github.com/dzhenetl/tplay/compare/v0.1.0...v0.1.1
[0.1.0]: https://github.com/dzhenetl/tplay/releases/tag/v0.1.0
