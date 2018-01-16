Releasing
=========

We intend to automate the releasing process and move it to the CI systems but at the following manual steps are required:

1. Update `src/sysl/__version__.py` based on the [latest released version](https://github.com/anz-bank/sysl/releases) and [SemVer](https://semver.org/):

		VERSION = (X, Y, Z)
2. Ensure you have a clean repo:
	- commit all un-committed changes
	- merge with `upstream master` via pull request
	- confirm passing CI and local tests
3. Create and push tag `vX.Y.Z`
4. Create a wheel distribution file:

		> python setup.py bdist_wheel

5. For first time usage only install `twine`
6. For first time usage only set up your `.pypirc` with credentials provided by Sysl core developers:

		> pip install twine
		> cat `~/.pypirc`
		[distutils]
		index-servers = pypi
		[pypi]
		username = <SYSL-CORE-DEV-PROVIDED>
		password = <SYSL-CORE-DEV-PROVIDED>


6. Upload to PyPI

		> twine upload dist/sysl-X.Y.Z-py2-none-any.whl


