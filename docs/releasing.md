Releasing
=========

Releases are available on Sysl's [Github releases page](https://github.com/anz-bank/sysl/releases) and on package registries (e.g. [Docker Hub](https://hub.docker.com/u/anzbank)).



GitHub Actions triggers the release workflow on push tag `v*.*.*`. See [sysl GitHub release workflow](https://github.com/anz-bank/sysl/blob/master/.github/workflows/release.yml) for further details.



Releasing is automated via [GoReleaser](https://goreleaser.com/).

GoReleaser creates and deploys `sysl-X.Y.Z-Os-Arch.tar.gz` and `sysl-X.Y.Z-Windows-Arch.zip` to the [Sysl Github Release page](https://github.com/anz-bank/sysl/releases). It also pushes Sysl's Docker Images `anzbank/sysl:latest` and `anzbank/sysl:X.Y.Z` to [Docker Hub](https://hub.docker.com/r/anzbank/sysl).

See [GoReleaser config file](https://github.com/anz-bank/sysl/blob/master/.github/workflows/.goreleaser.yml) for further details.