/** @type {import('next').NextConfig} */
const nextConfig = {
  async rewrites() {
    const rewriteRules = [];

    // process.envからPROXY_..._SOURCE_PATHというキーを探す
    for (const key in process.env) {
      if (key.startsWith('PROXY_') && key.endsWith('_SOURCE_PATH')) {
        // キーからプレフィックス部分（例: PROXY_00）を取得
        const prefix = key.replace('_SOURCE_PATH', '');
        
        const source = process.env[key];
        const destination = process.env[`${prefix}_DESTINATION_URL`];

        if (source && destination) {
          console.log(`[プロキシ] ルールを追加: ${source} -> ${destination}`);
          rewriteRules.push({
            source: source,
            destination: destination,
          });
        }
      }
    }

    return rewriteRules;
  },
};

export default nextConfig;
