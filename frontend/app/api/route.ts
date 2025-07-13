
import { NextResponse } from 'next/server';

export async function GET() {
  try {
    // Docker Compose環境では、サービス名でバックエンドにアクセスできます
    const backendUrl = 'http://backend:8080';
    const response = await fetch(backendUrl);

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    const text = await response.text();
    return new NextResponse(text);

  } catch (error) {
    console.error('Failed to fetch from backend:', error);
    // エラーが発生した場合は、500エラーとメッセージを返す
    return new NextResponse('Internal Server Error', { status: 500 });
  }
}
