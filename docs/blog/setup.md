The Sysl open source set-up
===========================

We, some of the ANZ open source developers, have recently set up what we consider to be good engineering practises for our enterprise-backed open source project [Sysl](https://github.com/anz-bank/sysl). The following is a record of systems and practises we deem essential and the products and platforms we have chosen (in parenthesis):

* Open Source Licence included with source code ([Apache License 2.0](https://github.com/anz-bank/sysl/blob/master/LICENSE))
* Version control system (git & Github)
* (Write) access control (Github organisation & teams)
* One step build (`pip install pytest flake8 -e .`)
* One step test (`pytest`)
* Continuous integration ([Travis CI](https://travis-ci.org/anz-bank/sysl) for Linux and macOS, [Appveyor CI](https://ci.appveyor.com/project/anz-bank/sysl) for Windows)
* Code reviews ([Github Code Reviews](https://github.com/features/code-review))
* Code coverage ([Codecov.io](https://codecov.io/github/anz-bank/sysl/))
* Automated release process (Tags trigger deployment from CI systems)
* Automated quality assurance (Pull requests are blocked until all checks pass)
* Issue tracking ([Github Issue](https://github.com/anz-bank/sysl/issues) tracking with [template](https://github.com/anz-bank/sysl/blob/master/ISSUE_TEMPLATE.md))
* Project Management ([Github projects](https://github.com/anz-bank/sysl/projects))
* Documentation in same repository as source code ([README](https://github.com/anz-bank/sysl/blob/master/README.rst) and [docs/](https://github.com/anz-bank/sysl/blob/master/docs) as starting point)
* Chat ([Gitter](https://gitter.im/anz-bank/sysl))
* Status dashboard ([Badges](https://github.com/anz-bank/sysl/blob/master/README.rst))

The most involved pieces have been automated quality assurance and release automation. These parts of the Sysl project also keep evolving as we add new quality checks and artefact types. They deserve a closer look.

Our branch model is simple - we develop feature branches on our own forks and issue pull requests to `master` on `anz-bank/sysl`. As suggested by [Github Flow](https://guides.github.com/introduction/flow/), "There's only one rule: anything in the master branch is always deployable". Releases are linked to tags of commits in `master` on `anz-bank/sysl`.

Our original (upstream, parent) repository is owned by the Github organization [`anz-bank`](https://github.com/anz-bank), which has a `sysl-developer` team. Only team members have write access to `anz-bank/sysl` and can merge pull requests into `master`. Additionally, the `master` branch of  `anz-bank/sysl` is protected and has restrictions enabled to automate quality checks:

 * Require pull request reviews before merging
 * Dismiss stale pull request approvals when new commits are pushed
 * Require review from Code Owners
 * Require branches to be up to date before merging
 * Require status checks to pass before merging
   - Require passing tests and no linting warning (on Travis CI Linux, Travis CI macOS, Appveyor Windows)
   - Require stable or improved codecoverage

Sysl releases are available on Sysl's [Github releases page](https://github.com/anz-bank/sysl/releases) and on various package registries (e.g. PyPI, BinTray, Docker Hub).  We follow [Semver](https://semver.org/) for versioning.

Our release process has been automated with script and artefact deployment by CI. A new release can be started with `src/scripts/release.sh prepare X.Y.Z`. After the automatically generated pull request is approved and merged, `src/scripts/release.sh deploy X.Y.Z` will create and push the release tag, which will then trigger Travis and Appveyor to deploy the artefacts (see the [Releasing documentation](https://github.com/anz-bank/sysl/blob/master/docs/releasing.md) for more details).

Happy coding and project set-up!

