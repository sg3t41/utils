'use client';

import { useEffect, useState } from 'react';

export default function HomePage() {
  const [message, setMessage] = useState('読み込み中...');

  useEffect(() => {
    // バックエンドのAPIエンドポイントを直接指定
    const apiUrl = '/api/health';
    // 環境変数からAPIキーを取得
    const apiKey = process.env.NEXT_PUBLIC_API_KEY;

    if (!apiKey) {
      setMessage(
        'APIキーが定義されていません。.envファイルにNEXT_PUBLIC_API_KEYを設定してください。',
      );
      return;
    }

    fetch(apiUrl, {
      headers: {
        // 'Authorization'ヘッダーに 'ApiKey {key}' の形式で設定
        Authorization: `ApiKey ${apiKey}`,
      },
    })
      .then((res) => {
        if (!res.ok) {
          // 認証エラーなどでレスポンスが成功でない場合
          throw new Error(`HTTPエラー ステータス: ${res.status}`);
        }
        return res.text();
      })
      .then((text) => setMessage(text))
      .catch((error) => {
        console.error('データ取得エラー: ', error);
        setMessage(`データの取得に失敗しました: ${error.message}`);
      });
  }, []);

  return (
    <div>
      <h1>Next.js フロントエンド</h1>
      <p>APIからの応答: {message}</p>
    </div>
  );
}
