/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  // Molstar makes prod build fail with swcMinify enabled
  // https://github.com/molstar/molstar/issues/1046
  // https://github.com/vercel/next.js/issues/52373
  swcMinify: false,
  output: "standalone",
  async redirects() {
    return [
      {
        source: "/",
        destination: "/tasks/protein-binder-design",
        permanent: false,
      },
      {
        source: "/tasks",
        destination: "/tasks/protein-binder-design",
        permanent: false,
      },
    ];
  },
};

export default nextConfig;
