Sysl
====
.. image:: https://img.shields.io/travis/anz-bank/sysl/master.svg?branch=master
   :target: https://travis-ci.org/anz-bank/sysl
.. image:: https://img.shields.io/appveyor/ci/anz-bank/sysl/master.svg?logo=appveyor
   :target: https://ci.appveyor.com/project/anz-bank/sysl
.. image:: https://img.shields.io/codecov/c/github/anz-bank/sysl/master.svg
   :target: https://codecov.io/gh/anz-bank/sysl/branch/master
.. image:: https://img.shields.io/gitter/room/nwjs/nw.js.svg
   :target: https://gitter.im/anz-bank/sysl

Sysl (pronounced "sizzle") is a system specification language. Using Sysl, you
can specify systems, endpoints, endpoint behaviour, data models and data
transformations. The Sysl compiler automatically generates sequence diagrams,
integrations, and other views. It also offers a range of code generation
options, all from one common Sysl specification.

The set of outputs is open-ended and allows for your own extensions. Sysl has been created with extensibility in mind and it will grow to support other representations over time.

Installation
------------
Install `Python 2.7 <https://www.python.org/downloads/>`_.
If your specific environment causes problems you might find our `guide <docs/environment_guide.md>`_ helpful.

Install Sysl with ::

  > pip install sysl

Now you can execute Sysl on the command line with ::

  > sysl   textpb -o out/petshop.txt /demo/petshop/petshop
  > reljam model /demo/petshop/petshop PetShopModel

See ``sysl --help`` and ``reljam --help`` for more options.

Development
-----------
Install dependencies and the ``sysl`` package with symlinks ::

  > pip install pytest flake8 -e .

Test and lint the source code and your changes with ::

  > pytest
  > flake8

Consider using `virtualenv <https://virtualenv.pypa.io/en/stable/>`_ and `virtualenvwrapper <https://virtualenvwrapper.readthedocs.io/en/latest/>`_ to isolate your development environment.

For Java tests install `Java 8 <https://docs.oracle.com/javase/8/docs/technotes/guides/install/install_overview.html>`_ and `gradle <https://gradle.org/install/>`_ and run ::

 > gradle test -b test/java/build.gradle

If your corporate environment restricts access to ``jcenter`` `this workaround <docs/gradle.md>`_ might solve the issue. Additionally, our `environment guide <docs/environment_guide.md>`_ has tips on using ``virtualenv`` with ``gradle test``.

We encourage contributions to this project! Please have a look at the `contributing guide <CONTRIBUTING.md>`_ for more information.

If you need to create a release follow the `release documentation <docs/releasing.md>`_.


Extending Sysl
--------------
In order to easily reuse and extend Sysl across systems, the Sysl compiler translates Sysl files
into an intermediate representation expressed as protocol buffer messages. These protobuf messages can be consumed in your favorite programming language and transformed to your desired output. In this way you are creating your own Sysl exporter.

Using the `protoc compiler <https://developers.google.com/protocol-buffers/>`_ you can translate the definition file of the intermediate representation ``src/proto/sysl.proto`` into your preferred programming language in a one-off step or on every build. You can then easily consume Sysl models in your programming language of choice in a typesafe way without having to write a ton of mapping
boilerplate. With that you can create your own tailored output diagrams, source code, views, integrations or other desired outputs.

In this project, several Python based exporters exist under ``src/sysl/exporters`` and the relevant Python protobuf definitions ``sysl_pb2.py`` have been created from ``sysl.proto`` with ::

  > protoc --python_out=src/sysl/proto  --proto_path=src/proto sysl.proto

If ``sysl.proto`` is updated, the above command needs to be re-run to update the corresponding Python definitions in ``sysl_pb2.py``.

Status
------
Sysl is currently targeted at early adopters. The current focus is to improve documentation and usability, especially error messages and warnings.
