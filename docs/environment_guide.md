Environment specific tips
=========================

OSX
---
On OSX we recommend installing Python 2.7 with [Homebrew](https://brew.sh/) and adding it to the `PATH` in your `.bash_profile`.

	brew install python
	echo export PATH=/usr/local/opt/python/libexec/bin:\$PATH >> ~/.bash_profile


Corporate environment
---------------------
If you are behind an NTLM proxy, you might need to setup `cntlm` to work effectively from your command line. Your specific environment and OS will determine this setup.

In case of SSL certificate interception you might need to add `--trusted-host pypi.python.org` to your pip commands or setup a global `pip.conf`
file (see [official docs](https://pip.pypa.io/en/stable/user_guide/#config-file) or [Stackoverflow](https://stackoverflow.com/a/46410817)).
