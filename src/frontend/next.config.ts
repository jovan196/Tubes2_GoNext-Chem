import { NextConfig } from 'next'

const nextConfig: NextConfig = {
  // Jika backend Go Anda berjalan di localhost:8080,
  // kita "proxy" /api/* ke sana
  async rewrites() {
    return [
      {
        source: '/api/:path*',
        destination: 'http://localhost:8080/api/:path*',
      },
    ]
  },
  // (opsional) agar Next.js bisa mengenali .ts pada config
  typescript: {
    ignoreBuildErrors: false,
  },
}

export default nextConfig;
