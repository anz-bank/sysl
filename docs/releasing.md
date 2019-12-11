Releasing
=========

Releases are available on Sysl's [Github releases page](https://github.com/anz-bank/sysl/releases) and on various package registries (e.g. PyPI, BinTray).

Releasing is automated via `pkg/scripts/release.sh`

A new release can be started with

	pkg/scripts/release.sh prepare X.Y.Z

and after the automatically generated pull request is approved and merged

	pkg/scripts/release.sh deploy X.Y.Z

will create and push the release tag, which will then trigger Travis and Appveyor to deploy the artefacts.

The generated [Github Release]((https://github.com/anz-bank/sysl/releases)) is created as a draft and needs to be manually published after adding release notes.


## Artefact deployment in detail

Travis CI creates and deploys `sysl-X.Y.Z-py2-none-any.whl`, `sysl-X.Y.Z.tar.gz` and `sysl-lib-X.Y.Z.jar` to the [Sysl Github Release page](https://github.com/anz-bank/sysl/releases) and Appveyor CI adds `sysl.exe`.

Travis also deploys Sysl's wheel and sdist distributions to [PyPI](https://pypi.python.org/pypi/sysl) and the Sysl Java library `sysl-lib-X.Y.Z.jar` to [BinTray](https://bintray.com/anz-bank/maven/sysl-lib/).

Travis also pushes Sysl's Docker Images `anzbank/sysl` and `anzbank/sysl:X.Y.Z` to Docker Hub. See `https://hub.docker.com/r/anzbank/sysl/` for more details.

Sysl generated Java code might have a dependency on `sysl-lib-X.Y.Z.jar`. To resolve this dependency with `gradle`, add `compile 'io.sysl:sysl-lib:X.Y.Z'` to your `build.gradle` file and for `maven` use:

```
<dependency>
  <groupId>io.sysl</groupId>
  <artifactId>sysl-lib</artifactId>
  <version>0.1.5</version>
  <type>pom</type>
</dependency>
```
See [sysl-lib on BinTray](https://bintray.com/anz-bank/maven/sysl-lib/) for further details.
