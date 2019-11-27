# Contributing

We appreciate pull requests from everyone!

Install the required development tools:

* git
* Python 2.7

Fork, then clone this repo:

    git clone https://github.com/YOUR_USERNAME/sysl.git

Optionally, for an isolated development environment consider using [virtualenv](https://virtualenv.pypa.io/en/stable/) and [virtualenvwrapper](https://virtualenvwrapper.readthedocs.io/en/latest/).

Set up your environment:

	pip install -e ".[dev]"

Test and lint the source code and your changes with:

	pytest
	flake

If you also want to run Java output related tests, install `Java 8` and `gradle` then run:

	cd test/java && gradle test

Finally, push to your fork and [submit a pull request](https://github.com/anz-bank/sysl/compare)

At this stage you are waiting on us to review your pull request. We like to at least comment on pull requests
within three business days (and often within one). We may suggest
some changes or improvements.

Some things that will increase the chance that your pull request is accepted:

* Write tests
* Don't break continuous integration on [Travis](https://travis-ci.org/anz-bank/sysl) or [Appveyor](https://ci.appveyor.com/project/anz-bank/sysl)
* Write a [good commit message](https://chris.beams.io/posts/git-commit/)
