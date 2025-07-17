import { NextResponse } from 'next/server';

export async function GET() {
  try {
    const backendUrl = process.env.BACKEND_INTERNAL_URL;
    if (!backendUrl) {
      throw new Error('環境変数にBACKEND_INTERNAL_URLが定義されていません。');
    }

    // 環境変数から取得したURLでバックエンドにアクセス
    const response = await fetch(backendUrl);

    if (!response.ok) {
      throw new Error(`HTTPエラー ステータス: ${response.status}`);
    }

    const text = await response.text();
    return new NextResponse(text);
  } catch (error) {
    console.error('バックエンドからのデータ取得に失敗しました:', error);
    // エラーが発生した場合は、500エラーとメッセージを返す
    return new NextResponse('内部サーバーエラー', { status: 500 });
  }
}
