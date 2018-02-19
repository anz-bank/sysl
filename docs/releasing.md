Releasing
=========

Releases are available on Sysl's [Github releases page](https://github.com/anz-bank/sysl/releases) and on various package registries (e.g. PyPI, BinTray).

Releasing is automated via `src/scripts/release.sh`

A new release can be started with

	src/scripts/release.sh prepare X.Y.Z

and after the automatically generated pull request is approved and merged

	src/scripts/release.sh deploy X.Y.Z

will create and push the release tag which will then trigger Travis and Appveyor to deploy the artefacts.

The generated [Github Release]((https://github.com/anz-bank/sysl/releases)) is created as a draft and needs to be manually published after adding release notes.


Artefact deployment in detail
-----------------------------
Travis CI creates and deploys `sysl-X.Y.Z-py2-none-any.whl`, `sysl-X.Y.Z.tar.gz` and `sysl-lib-X.Y.Z.jar` to the [Sysl Github Release page](https://github.com/anz-bank/sysl/releases) and Appveyor CI adds `sysl.exe`.

Travis also deploys Sysl's wheel and sdist distributions to [PyPI](https://pypi.python.org/pypi/sysl) and the Sysl Java library `sysl-lib-X.Y.Z.jar` to [BinTray](https://bintray.com/anz-bank/maven/sysl-lib/).

Sysl generated Java code might have a dependency on `sysl-lib-X.Y.Z.jar`. In order to resole this dependency for instance with gradle, add `compile 'io.sysl:sysl-lib:X.Y.Z'` to your `build.gradle` file (see [sysl-lib on BinTray](https://bintray.com/anz-bank/maven/sysl-lib/) for further details).
