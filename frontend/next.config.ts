import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  // Vercel デプロイ用設定
  ...(process.env.NODE_ENV === 'production' ? {} : { output: 'standalone' }),
  
  // 画像の最適化を無効化（外部APIから配信するため）
  images: {
    unoptimized: true,
    domains: [
      'localhost',
      // 本番APIドメインを後で追加
    ],
  },
  
  // ハイドレーションエラーの抑制
  experimental: {
    optimizeCss: false,
  },
  
  // 環境変数の設定
  env: {
    NEXT_PUBLIC_API_URL: process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080',
  },
};

export default nextConfig;
