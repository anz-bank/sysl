# Sysl UI

Built in React, serves up your defined services in a Web UI.

## Develop

Type `npm run` in the terminal to start up the UI locally

## Updating

To update the UI shipped with Sysl, the following steps need to be done:

### Automated

`make resources`

### Manual

- Create a release build using `npm run build`
- Generate a static go file from these assets using `pkger`
- Currently, we aren't able to get pkger to put the gofile in _pkg/ui due to a bug with an open PR.
  - The workaround is to manually copy it over to _pkg/ui_ and rename the package from `main` to `ui`
  - Execute this using `make resources`
- After that, we can run `make build` or `make release` as usual, to build the sysl releases. The static files are now just embedded in the Sysl binary, maintaining the portability advantages of Go.
