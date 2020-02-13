The Sysl open source setup
==========================

We, some of the ANZ open source developers, have recently set up what we consider to be good engineering practises for our enterprise-backed open source project [Sysl](https://github.com/anz-bank/sysl). The following is a record of systems and practises we deem essential and the products and platforms we have chosen (in parenthesis):

* Open Source Licence included with source code ([Apache License 2.0](https://github.com/anz-bank/sysl/blob/master/LICENSE))
* Version control system (git & Github)
* (Write) access control (Github organisation & teams)
* One step build (GitHub Actions & GoReleaser)(`pip install pytest flake8 -e .` Pysysl)
* One step test (go test)(`pytest` for pysysl)
* Continuous integration (GitHub Actions)
* Code reviews ([Github Code Reviews](https://github.com/features/code-review))
* Code coverage ([Codecov.io](https://codecov.io/github/anz-bank/sysl/))
* Automated release process (Specified master branch commit message triggers tag generation and CHANGELOG generation. Generated tags trigger deployment from CI systems. See the [Releasing documentation](https://github.com/anz-bank/sysl/blob/master/docs/releasing.md) for more details)
* Automated quality assurance (Pull requests are blocked until all checks pass)
* Issue tracking ([Github Issue](https://github.com/anz-bank/sysl/issues) tracking with [template](https://github.com/anz-bank/sysl/tree/master/.github/ISSUE_TEMPLATE))
* Project Management ([Github projects](https://github.com/anz-bank/sysl/projects))
* Documentation in same repository as source code ([README](https://github.com/anz-bank/sysl/blob/master/README.md) and [docs/](https://github.com/anz-bank/sysl/blob/master/docs) as starting point)
* Chat (Slack & <del>[Gitter](https://gitter.im/anz-bank/sysl)</del>)
* Status dashboard ([Badges](https://github.com/anz-bank/sysl/blob/master/README.md))

&nbsp;
The most involved pieces have been automated quality assurance and release automation. These parts of the Sysl project also keep evolving as we add new quality checks and artefact types. They deserve a closer look.

Our branch model is simple - we develop feature branches on our own forks and issue pull requests to `master` on `anz-bank/sysl`. As suggested by [Github Flow](https://guides.github.com/introduction/flow/), "There's only one rule: **anything in the master branch is always deployable**". Releases are linked to tags of commits in `master` on `anz-bank/sysl`.

Our original (upstream, parent) repository is owned by the Github organization [`anz-bank`](https://github.com/anz-bank), which has a `sysl-developer` team. Only team members have write access to `anz-bank/sysl` and can merge pull requests into `master`. Additionally, the `master` branch of  `anz-bank/sysl` is protected and has restrictions enabled to automate quality checks:

 * Require pull request reviews before merging
 * Dismiss stale pull request approvals when new commits are pushed
 * Require review from Code Owners
 * Require branches to be up to date before merging
 * Require status checks to pass before merging
   - Require passing tests and no linting warning (on GitHub Actions tests, GolangCI lint, deploy/netlify docs website)
   - Require stable or improved codecoverage
 * Require pull request reviews before publishing new version tag and release

&nbsp;
Happy coding and project setup!
