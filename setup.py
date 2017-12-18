#!/usr/bin/env python
# -*- coding: utf-8 -*-

import io
import os
import sys
from shutil import rmtree

from setuptools import find_packages, setup, Command
from setuptools.command.test import test as TestCommand

# Package meta-data.
NAME = 'sysl'
DESCRIPTION = 'System specification language with compiler and code generator'
URL = 'https://github.com/anz-bank/sysl'
EMAIL = 'marcelo.cantos@anz.com'
AUTHOR = 'ANZ'

REQUIRED = [
  'httplib2',
  'lxml',
  'openpyxl',
  'plantuml',
  'protobuf',
  'pylint',
  'PyYAML',
  'requests',
  'six'
]

here = os.path.abspath(os.path.dirname(__file__))

with io.open(os.path.join(here, 'README.rst'), encoding='utf-8') as f:
    long_description = '\n' + f.read()

about = {}
with open(os.path.join(here, 'src', NAME, '__version__.py')) as f:
    exec(f.read(), about)

class CleanCommand(Command):
    user_options = []
    def initialize_options(self):
        pass
    def finalize_options(self):
        pass
    def run(self):
        os.system('rm -vrf ./build ./dist ./.eggs ./*.pyc ./*.tgz ./*.egg-info')

class PyTest(TestCommand):
    user_options = [('pytest-args=', 'a', "Arguments to pass to pytest")]
    def initialize_options(self):
        TestCommand.initialize_options(self)
        self.pytest_args = ''
    def run_tests(self):
        import shlex
        #import here, cause outside the eggs aren't loaded
        import pytest
        errno = pytest.main(shlex.split(self.pytest_args))
        sys.exit(errno)

class DistCommand(Command):
    user_options = []
    def initialize_options(self):
        pass
    def finalize_options(self):
        pass
    def run(self):
        try:
            rmtree(os.path.join(here, 'dist'))
        except OSError:
            pass
        os.system('{0} setup.py sdist bdist_wheel --universal'.format(sys.executable))
        sys.exit()


setup(
    name=NAME,
    version=about['__version__'],
    description=DESCRIPTION,
    long_description=long_description,
    author=AUTHOR,
    author_email=EMAIL,
    url=URL,
    package_dir={'':'src'},
    packages=find_packages('src', exclude=('tests',)),
    # entry_points={
    #     'console_scripts': ['mycli=mymodule:cli'],
    # },
    install_requires=REQUIRED,
    include_package_data=True,
    license='Apache 2.0',
    classifiers=[
        'License :: OSI Approved :: Apache Software License',
        'Programming Language :: Python',
        'Programming Language :: Python :: 2.7',
    ],
    setup_requires=['pytest-runner'],
    tests_require=['pytest'],
    cmdclass={
        'clean': CleanCommand,
        'test': PyTest,
        'dist': DistCommand, 
    }
)