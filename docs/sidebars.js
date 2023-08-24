// @ts-check

const { readdirSync } = require('fs');
const capitalize = require('lodash/capitalize');

/** @type {import('@docusaurus/plugin-content-docs').SidebarsConfig} */
const sidebars = {
  // By default, Docusaurus generates a sidebar from the docs folder structure
  tutorialSidebar: [
    {
      type: 'doc',
      id: 'welcome/welcome',
      label: 'Welcome',
    },
    {
      type: 'category',
      label: 'Quickstart',
      collapsed: false,
      items: [
        {
          type: 'doc',
          id: 'quickstart/installation',
          label: 'Installation',
        },
        {
          type: 'doc',
          id: 'quickstart/available-tools',
          label: 'Available Tools',
        }
      ],
    },
    {
      type: 'category',
      label: 'Concepts',
      collapsed: true,
      items: [
        {
          type: 'autogenerated',
          dirName: 'concepts',
        },
      ]
    },
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
    {
      type: 'category',
      label: 'Get Involved',
      items: [
        // 'get-involved/ways-to-get-involved',
        'get-involved/how-to-contribute-a-tool',
      ]
    },
    {
      type: 'category',
      label: 'About Us',
      collapsed: true,
      items: [
        {
          type: 'autogenerated',
          dirName: 'about-us',
        },
      ],
    }
  ],
};

module.exports = sidebars;
