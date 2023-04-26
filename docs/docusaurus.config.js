// @ts-check
// Note: type annotations allow type checking and IDEs autocompletion

const lightCodeTheme = require('prism-react-renderer/themes/github');
const darkCodeTheme = require('prism-react-renderer/themes/dracula');

const createConfig = async () => {
  const mdxMermaid = await import('mdx-mermaid');

  /** @type {import('@docusaurus/types').Config} */
  return {
    title: 'LabDAO Documentation',
    tagline: '',
    url: 'https://docs.labdao.com',
    baseUrl: '/',
    onBrokenLinks: 'throw',
    onBrokenMarkdownLinks: 'warn',
    favicon: 'img/LabDAO_Favicon_Teal.png',
    organizationName: 'labdao', // Usually your GitHub org/user name.
    projectName: 'docs', // Usually your repo name.
  
    presets: [
      [
        'classic',
        /** @type {import('@docusaurus/preset-classic').Options} */
        ({
          docs: {
            sidebarPath: require.resolve('./sidebars.js'),
            remarkPlugins: [mdxMermaid.default],
            editUrl: ({ docPath }) => {
              const pathArr = docPath.split('/');
              let repo = 'docs';
              let pathSliceIndex = 0;
              if (pathArr.length > 1 && pathArr[0] === '_projects') {
                repo = pathArr[1];
                pathSliceIndex = 3;
              }
              return `https://github.com/labdao/${repo}/edit/main/docs/${pathArr.slice(pathSliceIndex).join('/')}`
            },
            routeBasePath: '/',
            exclude: [
              'projects/*/*.{md,mdx}',
              '**/_*.{md,mdx}',
            ],
            include: [
              '_projects/*/docs/**/*.{md,mdx}',
              '*/**/*.{md,mdx}',
            ],
          },
          blog: false,
          theme: {
            customCss: require.resolve('./src/css/custom.css'),
          },
          gtag: {
            trackingID: 'G-5MQBDBEY04',
            anonymizeIP: true,
          },
        }),
      ],
    ],
  
    themeConfig:
      /** @type {import('@docusaurus/preset-classic').ThemeConfig} */
      ({
        navbar: {
          title: 'LabDAO',
          logo: {
            alt: 'LabDAO Logo',
            src: 'img/labdaologo_brandmark_Teal.png',
          },
          items: [
            {
              href: 'https://github.com/labdao/docs',
              position: 'right',
              className: 'header-github-link',
            },
          ],
        },
        algolia: {
          appId: 'I8J1DZKSGR',
          apiKey: 'd78d134e15f8f366b04ee89599fe233a',
          indexName: 'labdao',
          debug: false,
        },
        footer: {
          style: 'dark',
          links: [
            {
              title: "Community",
              items: [
                {
                  label: "LabDAO.xyz",
                  href: "https://labdao.xyz/community",
                },
                {
                  label: "Discord",
                  href: "https://discordapp.com/invite/labdao",
                },
              ]
            },
            {
              title: "Develop/Contribute",
              items: [
                {
                  label: "GitHub",
                  href: "https://github.com/labdao",
                }
              ],
            },
            {
              title: "Socials",
              items: [
                {
                  label: "Twitter",
                  href: "https://twitter.com/lab_dao",
                },
              ],
            },
          ],
          copyright: `Copyright © ${new Date().getFullYear()} LabDAO.`,
        },
        prism: {
          theme: lightCodeTheme,
          darkTheme: darkCodeTheme,
        },
      }),
  };
}

module.exports = createConfig;
