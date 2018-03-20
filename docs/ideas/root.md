Sysl root
=========

When executing `sysl` one needs to specify a root directory with `--root DIR`, which serves as root for specifying `sysl modules` on the command line and in import statement:

* `sysl textpb --root sysl_dir /town`
* `import /transit/bus` at top of `town.sysl` file

If `--root` isn't set it defaults to `.` - the current working directory.
Modules, similar to Python modules, are derived from the file name stripping the `.sysl` extension and adding a leading `/`.

### Future enhancement ideas
* allow for file name extension and relative paths in import statements:
	`import ../utils.sysl`
* specify paths relative to `root` with `//`
	`import //utils`
* allow for filenames rather than modules on the command line:
	`sysl textpb demo/simple.sysl`
