Releasing
=========

Releases are available on [Sysl's GitHub releases page](https://github.com/anz-bank/sysl/releases) and on package registries (e.g. [Docker Hub](https://hub.docker.com/u/anzbank)).
&nbsp;
Sysl is using [GitHub Actions](https://help.github.com/en/actions/getting-started-with-github-actions/about-github-actions) as continuous integration (CI) and continuous deployment (CD) system.
&nbsp;

### Steps to publish new release
1. Anyone who wants to publish a new release has to create a pull request which needs approvement before is merged.
2. A merged PR with commit message contains `#major` or `#minor`(we don't use SemVer `patch` versions in Sysl) will trigger the [Generate-Tag workflow](https://github.com/anz-bank/sysl/blob/master/.github/workflows/generate-tag.yml). It will:
	1. Generate and push the respective version tag according to the commit message
	2. Auto-generate the related CHANGELOG
	3. Commit and push CHANGELOG's changes to a new branch and create a PR to merge the CHANGELOG file
	4. _The PR be can be auto-approved and auto-merged. But it isn't right now. The idea is to let others notice that there's a new release and maybe tidy the CHANGELOG if they want_
	> Note: This action will not bump the tag if the HEAD commit has already been tagged. If two or more keywords are present, the highest-ranking one will take precedence. 
	
	> We follow [Semver](https://semver.org/) for versioning.
3. The version tag push event will trigger the [Release workflow](https://github.com/anz-bank/sysl/blob/master/.github/workflows/release.yml) to publish release to [Sysl's GitHub releases page](https://github.com/anz-bank/sysl/releases) and [Docker Hub](https://hub.docker.com/r/anzbank/sysl).
	> Releasing is automated via [GoReleaser](https://goreleaser.com/). GoReleaser creates and deploys `sysl-X.Y.Z-Os-Arch.tar.gz` and `sysl-X.Y.Z-Windows-Arch.zip` to the [Sysl Github Release page](https://github.com/anz-bank/sysl/releases). It also pushes Sysl's Docker Images `anzbank/sysl:latest` and `anzbank/sysl:X.Y.Z` to [Docker Hub](https://hub.docker.com/r/anzbank/sysl). See [GoReleaser config file](https://github.com/anz-bank/sysl/blob/master/.github/workflows/.goreleaser.yml) for further details.

### A tested example
* The new release proposal PR: https://github.com/anz-bank/sysl/pull/617
* The commit which merges new release PR: https://github.com/anz-bank/sysl/commit/edd3c4ee6cab7bcf99d580efef0e14669a597374
* The triggered GitHub workflow which generates tag: https://github.com/anz-bank/sysl/runs/451818569?check_suite_focus=true
* The triggered GitHub workflow which publishes respective release: https://github.com/anz-bank/sysl/runs/451819073?check_suite_focus=true
* The published binary release: https://github.com/anz-bank/sysl/releases/tag/v0.6.2
* The published Docker image: https://hub.docker.com/r/anzbank/sysl/tags?page=1 (image `latest` and `0.6.2`)
* The auto-created PR to merge generated CHANGELOG: https://github.com/anz-bank/sysl/pull/618
