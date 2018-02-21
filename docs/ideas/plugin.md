Plugin architecture
===================

Create front-end/driver that runs core plugins (e.g. reljam, textpb) or custom plugins (e.g. mykotlin-exporter), similar to git or protobuf. The Front-end looks for an executable `sysl-myplugin` if called as `sysl myplugin`.

There are two types of plugins:

* Export (default)
* Import

#### Export plugin interface:
The front-end provides additional command line arguments and \*.sysl file serialized as proto-message on stdin to the export plugin exe.

#### Import plugin interface:
The front-end provides additional command line arguments.
An Import Plugin is expected to produce \*.sysl files or sysl proto-message on stdout.
There's a special command line argument required to specify that a plugin is a _import_ plugin: e.g. `sysl my-import -I #...`
