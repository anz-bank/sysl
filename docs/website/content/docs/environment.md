---
title: "Environment tips"
description: "Trouble setting up a Sysl dev environment inside a restrictive enterprise network? Find some tips and workarounds here."
date: 2018-02-27T15:51:34+11:00
weight: 100
draft: false
bref: "Tips and workarounds for enterprise and other environments"
toc: true
---

Python on OSX
-------------
On OSX we recommend installing Python 2.7 with [Homebrew](https://brew.sh/) and adding it to the `PATH` in your `.bash_profile`.

	brew install python
	echo export PATH=/usr/local/opt/python/libexec/bin:\$PATH >> ~/.bash_profile

Corporate env
-------------
If you are behind an NTLM proxy, you might need to setup `cntlm` to work effectively from your command line. Your specific environment and OS will determine this setup.

In case of SSL certificate interception you might need to add `--trusted-host pypi.python.org` to your pip commands or setup a global `pip.conf`
file (see [official docs](https://pip.pypa.io/en/stable/user_guide/#config-file) or [Stackoverflow](https://stackoverflow.com/a/46410817)).

Gradle setup
------------
In order to run Java tests install [Java 8](https://docs.oracle.com/javase/8/docs/technotes/guides/install/install_overview.html) and [gradle](https://gradle.org/install/). Download or clone the [Sysl repository](https://github.com/anz-bank/sysl) and execute:

	> cd test/java && gradle test

If you get the following error message

	Could not resolve all files for configuration ':compileClasspath'.

this might be related to your corporate environment not being able to access `jcenter`.
Try creating a `<GRADLE-OVERRIDES-FILE>` file with the following content

```
allprojects {
    repositories {
        maven {
            url 'https://<LOCAL_DOMAIN_AND_PATH>/jcenter'
        }
    }
}
```

Replace `<LOCAL_DOMAIN_AND_PATH>` with your local domain and path mirroring `jcenter`.

Then run:

	> gradle test -I <GRADLE-OVERRIDES-FILE>

Java test
---------
If you work with `virtualenv` and `gradle` reports it cannot find `sysl`, try using the `--no-daemon` option:

	gradle --no-daemon test -b test/java/build.gradle

