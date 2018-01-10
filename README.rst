Sysl
====

.. image:: https://img.shields.io/travis/anz-bank/sysl/master.svg?label=Linux%20build%20%40%20Travis%20CI
   :target: http://travis-ci.org/anz-bank/sysl
.. image:: https://img.shields.io/appveyor/ci/anz-bank/sysl.svg?label=Windows%20build%20%40%20Appveyor%20CI
   :target: https://ci.appveyor.com/project/anz-bank/sysl

Sysl (pronounced "sizzle") is a system specification language. Using Sysl, you
can specify systems, endpoints, endpoint behaviour, data models and data
transformations. The Sysl compiler automatically generates sequence diagrams,
integrations, and other views. It also offers a range of code generation
options, all from one common Sysl specification.

The set of outputs is open-ended and allows for your own extensions. Sysl has been created with extensibility in mind and it will grow to support other representations over time.

Installation
------------
If you are interested in trying out Sysl, you will need to install `Python 2.7 <https://www.python.org/downloads/>`_ and the Python package manager `Pip <https://pip.pypa.io/en/stable/installing/>`_. You can then build sysl from source with ::

  > git clone https://github.com/anz-bank/sysl.git
  > cd sysl
  > pip install .

Now you can execute Sysl as command line tool with ::

  > sysl    --root demo/petshop textpb -o out/petshop.txt /petshop
  > reljam  --root demo/petshop model /petshop PetShopModel

which is equivalent to ::

  > sysl textpb -o out/petshop.txt /demo/petshop/petshop
  > reljam model /demo/petshop/petshop PetShopModel

See ``sysl --help`` and ``reljam --help`` from more options.
If you are behind a corporate proxy you might want to consider setting up a global ``pip.conf``
file (details on `Stackoverflow <https://stackoverflow.com/a/46410817>`_ or `official docs <https://pip.pypa.io/en/stable/user_guide/#config-file>`_).

Development
-----------
Install dependencies and the ``sysl`` package with symlinks ::

  > pip install -e ".[dev]"

Test and lint the source code and your changes with ::

  > pytest
  > flake

and create a distribution with ::

  > python setup.py bdist_wheel --universal

Consider using `virtualenv <https://virtualenv.pypa.io/en/stable/>`_ and `virtualenvwrapper <https://virtualenvwrapper.readthedocs.io/en/latest/>`_ to isolate your environment.

Extending Sysl
--------------
In order to easily reuse and extend Sysl across systems, the Sysl compiler translates Sysl files
into an intermediate representation expressed as protocol buffer messages. These protobuf messages can be consumed in your favorite programming language and transformed to your desired output. In this way you are creating your own Sysl exporter.

Using the `protoc compiler <https://developers.google.com/protocol-buffers/>`_ you can translate the definition file of the intermediate representation ``src/proto/sysl.proto`` into your preferred programming language in a one of step or on every build. You can then easily consume Sysl models in your programming language of choice  in a typesafe way without having to write a ton of mapping
boilerplate. With that you can create your own tailored output diagrams, source code, views, integrations or other desired outputs.

In this project, several Python based exporters exist under ``src/sysl/exporters`` and the relevant Python protobuf definitions ``sysl_pb2.py`` have been created from ``sysl.proto`` with ::

  > protoc --python_out=src/sysl/proto  --proto_path=src/proto sysl.proto

If ``sysl.proto`` is updated the above command needs to be re-run to update the corresponding Python definitions in ``sysl_pb2.py``.

Status
------
Sysl is currently targeted at early adopters. The current focus is to improve documentation and usability, especially error messages and warnings.
