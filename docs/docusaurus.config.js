module.exports = {
  title: "Sysl",
  tagline:
    "Deliver code, data models and visualisations from a single language",
  url: "https://sysl.io",
  baseUrl: "/",
  favicon: "img/logo-blue-net-s.png",
  organizationName: "anz-bank",
  projectName: "sysl",
  themeConfig: {
    navbar: {
      title: "",
      logo: {
        alt: "Sysl Logo",
        src: "img/logo-blue.png",
        srcDark: "img/logo-white.png",
      },
      items: [
        { to: "/", label: "Home", position: "right" },
        { to: "docs", label: "Docs", position: "right" },
        { to: "docs/discussions", label: "Community", position: "right" },
        { to: "blog", label: "Blog", position: "right" },
        {
          href: "https://play.sysl.io/",
          label: "Playground",
          position: "right",
        },
        {
          href: "https://github.com/anz-bank/sysl",
          position: "right",
          className: "header-github-link",
          label: "GitHub",
          "aria-label": "GitHub repository",
        },
      ],
    },
    footer: {
      copyright: `Copyright © ${new Date().getFullYear()} Sysl (Apache-2.0 License)`,
    },
    algolia: {
      apiKey: "e746801b8fd0862b43f994ad13cff9b5",
      indexName: "sysl",
    },
    colorMode: {
      defaultMode: "light",
      disableSwitch: false,
    },
    googleAnalytics: {
      trackingID: "UA-173443254-1", // GA account owner: sysl.usr@gmail.com
      anonymizeIP: true,
    },
  },
  presets: [
    [
      "@docusaurus/preset-classic",
      {
        docs: {
          sidebarPath: require.resolve("./sidebars.js"),
          editUrl: "https://github.com/anz-bank/sysl-website/edit/master/",
          admonitions: {
            infima: true,
            customTypes: {
              right: {
                ifmClass: "success",
                keyword: "right",
                svg: '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><path d="M20.285 2l-11.285 11.567-5.286-5.011-3.714 3.716 9 8.728 15-15.285z"/></svg>',
              },
              wrong: {
                ifmClass: "danger",
                keyword: "wrong",
                svg: '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><path d="M24 20.188l-8.315-8.209 8.2-8.282-3.697-3.697-8.212 8.318-8.31-8.203-3.666 3.666 8.321 8.24-8.206 8.313 3.666 3.666 8.237-8.318 8.285 8.203z"/></svg>',
              },
            },
          },
        },
        community: {
          homePageId: "discussions",
          sidebarPath: require.resolve("./sidebars.js"),
        },
        blog: {
          showReadingTime: true,
          editUrl: "https://github.com/anz-bank/sysl-website/edit/master/",
        },
        theme: {
          customCss: require.resolve("./src/css/custom.css"),
        },
      },
    ],
  ],
  plugins: [],
  scripts: [],
  stylesheets: [
    "https://fonts.googleapis.com/css?family=Lato:wght@400;900|Roboto|Source+Code+Pro",
    "https://at-ui.github.io/feather-font/css/iconfont.css",
  ],
  themes: [],
};
