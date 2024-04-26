/**
 * Run `build` or `dev` with `SKIP_ENV_VALIDATION` to skip env validation. This is especially useful
 * for Docker builds.
 */

/** @type {import("next").NextConfig} */
const config = {
    eslint:{
        ignoreDuringBuilds: true,
        

    },
    reactStrictMode: false,
    images: {
        remotePatterns: [
          {
            protocol: 'https',
            hostname: "**",
         
          },
          {
            protocol: 'http',
            hostname:"**"
          }
        ],
      },
      async headers() {
        return [
          {
            source: '/watch',
            headers: [
              {
                key: 'Cache-Control',
                value: 'no-store',
              },
            ],
          },
        ]
      },
};

export default config;
