/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: false,
  // Molstar makes prod build fail with swcMinify enabled
  // https://github.com/molstar/molstar/issues/1046
  // https://github.com/vercel/next.js/issues/52373
  swcMinify: false,
  output: "standalone",
  async redirects() {
    return [
      {
        source: "/",
        destination: "/experiments/new/protein-binder-design",
        permanent: false,
      },
      {
        source: "/tasks/:path*",
        destination: "/experiments/new/:path*",
        permanent: false,
      },
      {
        source: "/experiments/new",
        destination: "/experiments/new/protein-binder-design",
        permanent: false,
      },
    ];
  },
};

export default nextConfig;
