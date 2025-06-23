import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  output: 'standalone',
  
  // 画像の最適化を無効化（外部APIから配信するため）
  images: {
    unoptimized: true,
  },
  
  // ハイドレーションエラーの抑制
  experimental: {
    optimizeCss: false,
  },
};

export default nextConfig;
