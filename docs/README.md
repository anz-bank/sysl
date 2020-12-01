# Sysl Website

Sysl website https://sysl.io is built using [Docusaurus 2](https://v2.docusaurus.io/).

## Contributing

### Project structure

```
sysl                        // repo root
└── docs                       // documentation website content
    ├── blog                       // all the blogs
    ├── docs                       // all the docs
    ├── src
    │   ├── css                     // common CSS
    │   └── pages                   // landing page components and CSS
    ├── static                      // website assets
    │   └── img
    ├── docusaurus.config.js        // configuration
    ├── sidebars.js                 // sidebar management
    ├── package.json                // dependencies and scripts
    ├── README.md                   // this file
    └── yarn.lock                   // dependency version manifest
```

### Contribute to Docs

#### Edit Existing Doc File

All documentation is [Markdown](https://daringfireball.net/projects/markdown/syntax) files under `docs/`. Docusaurus 2 can do [more than just parsing Markdown](https://v2.docusaurus.io/docs/markdown-features).

When adding images, [useBaseUrl](https://v2.docusaurus.io/docs/docusaurus-core/#usebaseurl) (instead of standard instead of Markdown image syntax) ensures the links will be correct when published.

#### Add New Doc File

1. Create a new Markdown file in `docs/` with a [Docusaurus header](https://v2.docusaurus.io/docs/markdown-features#markdown-headers).
2. Add it to `sidebars.js`.

### Contribute to Blog

All blog posts are also Markdown files under `blog/`. Following [the Docusaurus instructions](https://v2.docusaurus.io/docs/blog) to contribute to blog.

### Contribute to Homepage

The homepage code of this website is in `src/pages/index.js` with configuration `docusaurus.config.js`. Follow [configuration docs](https://v2.docusaurus.io/docs/configuration) to contribute to Homepage.

## Development

### Requirements

- [Node.js](https://nodejs.org/en/download/)
- [Yarn](https://classic.yarnpkg.com/en/docs/install#mac-stable)
- [pyspell](https://facelessuser.github.io/pyspelling/) `pip3 install pyspell`
- [aspell](http://aspell.net/) `brew install aspell`

### Local Development

```
yarn install
yarn start
```

This command starts a local development server and open up a browser window. Most changes are reflected live without having to restart the server or refresh the browser.

### Build

```
yarn build
```

This command generates static content into the `build` directory and can be served using any static contents hosting service.

### Format

This project uses [Prettier](https://prettier.io/) to format all files.

```
yarn format
```

### Deployment

A development version is deployed using [Netlify](https://www.netlify.com/). Deployment previews are enabled, so each PR has a unique deployment preview link which can be found in the Github Status Checks.

### Linter

<!-- [spellcheck](https://github.com/marketplace/actions/github-spellcheck-action) GitHub Action is used as the English spelling check linter. Add custom terms in `.wordlist.txt` to pass the spelling check. To run the spellcheck locally, run `npm run spellcheck` -->

The [markdown-link-check](https://github.com/marketplace/actions/markdown-link-check) GitHub Action validates that the docs contain no broken links.

### Search

[Algolia DocSearch](https://docsearch.algolia.com/) is used for website searching. The crawler is configured in [docsearch-config](https://github.com/algolia/docsearch-configs/blob/master/configs/sysl.json) and runs daily.
