// @ts-check

const { readdirSync } = require('fs');
const capitalize = require('lodash/capitalize');

/** @type {import('@docusaurus/plugin-content-docs').SidebarsConfig} */
const sidebars = {
  tutorialSidebar: [
    {
      type: 'doc',
      id: 'tutorials/tutorials',
      label: 'Tutorials',
    },
    {
      type: 'category',
      label: 'Reference',
      collapsed: false,
      items: [
        {
          type: 'autogenerated',
          dirName: 'reference',
        },
      ],
    },
  ],
};

module.exports = sidebars;