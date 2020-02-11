Releasing
=========

Releases are available on [Sysl's GitHub releases page](https://github.com/anz-bank/sysl/releases) and on package registries (e.g. [Docker Hub](https://hub.docker.com/u/anzbank)).
&nbsp;
Sysl is using [GitHub Actions](https://help.github.com/en/actions/getting-started-with-github-actions/about-github-actions) as continuous integration (CI) and continuous deployment (CD) tool.
&nbsp;
### Steps to publish new release
1. Anyone who wants to publish a new release has to create a pull request which needs approvement before is merged.
2. A merged PR with commit message contains `#major`, `#minor` or `#patch` will trigger the [Generate-Tag workflow](https://github.com/anz-bank/sysl/blob/master/.github/workflows/generate-tag.yml) to generate and push the respective version tag. 
	> Note: This action will not bump the tag if the HEAD commit has already been tagged. If two or more keywords are present, the highest-ranking one will take precedence. 
	
	> We follow [Semver](https://semver.org/) for versioning.
3. The version tag push event will then trigger the [Release workflow](https://github.com/anz-bank/sysl/blob/master/.github/workflows/release.yml) to publish release to [Sysl's GitHub releases page](https://github.com/anz-bank/sysl/releases).
	> Releasing is automated via [GoReleaser](https://goreleaser.com/). GoReleaser creates and deploys `sysl-X.Y.Z-Os-Arch.tar.gz` and `sysl-X.Y.Z-Windows-Arch.zip` to the [Sysl Github Release page](https://github.com/anz-bank/sysl/releases). It also pushes Sysl's Docker Images `anzbank/sysl:latest` and `anzbank/sysl:X.Y.Z` to [Docker Hub](https://hub.docker.com/r/anzbank/sysl). See [GoReleaser config file](https://github.com/anz-bank/sysl/blob/master/.github/workflows/.goreleaser.yml) for further details.
