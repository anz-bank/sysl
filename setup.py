#!/usr/bin/env python
# -*- coding: utf-8 -*-

import io
import os
import sys
import platform
from shutil import rmtree

from setuptools import find_packages, setup, Command

if platform.system() == 'Windows':
    import py2exe

NAME = 'sysl'
DESCRIPTION = 'System specification language with compiler and code generator'
URL = 'https://github.com/anz-bank/sysl'
EMAIL = 'marcelo.cantos@anz.com'
AUTHOR = 'ANZ'

REQUIRED = [
    'httplib2',
    'urllib3',
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

setup(
    name=NAME,
    version=about['__version__'],
    description=DESCRIPTION,
    long_description=long_description,
    author=AUTHOR,
    author_email=EMAIL,
    url=URL,
    package_dir={'': 'src'},
    packages=find_packages('src', exclude=('tests',)),
    entry_points={
        'console_scripts': [
            'sysl=sysl.core.__main__:main',
            'reljam=sysl.reljam.reljam:main',
        ],
    },
    install_requires=REQUIRED,
    include_package_data=True,
    license='Apache 2.0',
    classifiers=[
        'License :: OSI Approved :: Apache Software License',
        'Programming Language :: Python',
        'Programming Language :: Python :: 2.7',
    ],
    extras_require={
        'dev': [
            'pytest',
            'flake8',
        ]
    },
    # py2exe
    options={'py2exe': {
        'bundle_files': 1,
        'dll_excludes': ['w9xpopen.exe', 'libstdc++-6.dll', 'libgcc_s_dw2-1.dll']
    }},
    console=[{
        'script': 'src/sysl/core/__main__.py',
        'dest_base': 'sysl',
        'icon_resources': [(1, 'docs/favicon.ico')]
    }, {
        'script': 'src/sysl/reljam/__main__.py',
        'dest_base': 'reljam',
        'icon_resources': [(1, 'docs/favicon.ico')]
    }],
    zipfile=None
)
