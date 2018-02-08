Releasing
=========

Releasing is automated a script `src/scripts/releasing.sh`

Travis creates and deploys `sysl-X.Y.Z-py2-none-any.whl` and `sysl-X.Y.Z.tar.gz` to the [Sysl Github Release page](https://github.com/anz-bank/sysl/releases) and Appveyor adds `sysl.exe`. Travis also deploys Sysl's wheel and sdist distributions to PyPI.

A new release can be started with

	src/scripts/release.sh prepare X.Y.Z

and after the automatically generated pull request is approved and merged

	src/scripts/release.sh deploy X.Y.Z

will create and push the tag which will trigger Travis and Appveyor to depoy the new release.
