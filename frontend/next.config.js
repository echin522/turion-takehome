/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  // https://nextjs.org/docs/app/api-reference/next-config-js/output#automatically-copying-traced-files
  output: "standalone",
  // TODO: Fix webpack cache warning for large strings.
  experimental: {
    optimizePackageImports: ["@mantine/core", "@mantine/hooks"],
  },
  staticPageGenerationTimeout: 180,
  swcMinify: false,
};

module.exports = nextConfig;
