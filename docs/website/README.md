# Sysl website

This directory contains the source for the [sysl.io](https://sysl.io) website built with [Hugo](https://gohugo.io/), a static site generator which produces html from markdown. In this directory, `<repo>/docs/website` start `hugo serve` and view contents on [http://localhost:1313/](http://localhost:1313/) for updating docs and live reloading.

The website can be built with `hugo` which puts all content into the `public` directory.
On every merge into `upstream` `master` the website gets updated with the Netflify-Hugo-Github integration (see `<repo>/netlify.tom`, [Netlify docs](https://gohugo.io/hosting-and-deployment/hosting-on-netlify/)).

## Updating CSS

In order to update CSS, work with the unminified CSS in `static/css/` and change `layouts/_default/baseof.html` to use these unminified files (commented out there). When done with the changes, minify the updated css and revert to using it:

1. Remove references to `styles.min.css` in `layouts/_default/baseof.html` and reference unminified files instead (see comment there).
2. Update the uniminfied CSS files
3. `rm static/css/styles.min.css`
4. `hugo`
5. `npm install -g purify-css` (first time only)
6. `purifycss static/css/*.css public/**/*.html public/*.html static/js/jquery-2.1.4.min.js static/js/kube.min.js -im -o static/css/styles.min.css`
7. Revert `layouts/_default/baseof.html` to use `styles.min.css` again
