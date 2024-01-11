// These could eventually be fetched from the backend, but for now are hardcoded
// Adding a task? Add it's image to /public/images with filename task-<slug>.png

export const tasks = [
  {
    name: "Protein Binder Design",
    slug: "protein-binder-design",
    available: true,
    default_tool: {
      CID: "QmYmKL9fCzxoEQAEvCtLZHXXwjJWdn7bbRE9afTXV5cB3v",
    },
  },
  {
    name: "Protein Folding",
    slug: "protein-folding",
    available: false,
  },
  {
    name: "Protein Docking",
    slug: "protein-docking",
    available: false,
  },
  {
    name: "Small Molecule Docking",
    slug: "small-molecule-docking",
    available: false,
  },
];
