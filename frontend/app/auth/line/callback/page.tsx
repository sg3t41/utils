'use client';

import { useEffect, useState } from 'react';
import { useSearchParams, useRouter } from 'next/navigation';
import Link from 'next/link';
import { useAuth } from '../../../../contexts/AuthContext';

export default function LineCallbackPage() {
  const searchParams = useSearchParams();
  const router = useRouter();
  const { login } = useAuth();
  const [loading, setLoading] = useState(true);
  const [result, setResult] = useState<any>(null);
  const [error, setError] = useState<string | null>(null);
  const [loginSuccess, setLoginSuccess] = useState(false);

  useEffect(() => {
    const code = searchParams.get('code');
    const state = searchParams.get('state');
    const errorParam = searchParams.get('error');

    if (errorParam) {
      setError(`認証エラー: ${errorParam}`);
      setLoading(false);
      return;
    }

    if (!code) {
      setError('認証コードが取得できませんでした');
      setLoading(false);
      return;
    }

    handleCallback(code, state);
  }, [searchParams]);

  const handleCallback = async (code: string, state: string | null) => {
    try {
      const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';
      const response = await fetch(`${apiUrl}/api/v1/auth/line/callback`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          code: code,
          state: state,
        }),
      });

      if (!response.ok) {
        throw new Error('コールバック処理でエラーが発生しました');
      }

      const data = await response.json();
      setResult(data);
      
      // ログイン処理
      if (data.access_token && data.user) {
        // APIから取得したユーザー情報を使用
        const user = {
          id: data.user.id,
          name: data.user.name,
          email: data.user.email || '', 
          line_user_id: data.user.line_user_id,
          profile_image: data.user.profile_image || '',
          created_at: data.user.created_at,
          updated_at: data.user.updated_at,
        };
        
        // refresh_tokenがない場合は空文字列を使用
        const refreshToken = data.refresh_token || '';
        
        login(data.access_token, refreshToken, user);
        setLoginSuccess(true);
        
        // 3秒後に記事一覧ページにリダイレクト
        setTimeout(() => {
          router.push('/articles');
        }, 3000);
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : 'エラーが発生しました');
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <div className="min-h-screen grid place-items-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-green-500 mx-auto"></div>
          <p className="mt-4 text-gray-600">LINEログイン処理中...</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="min-h-screen grid place-items-center">
        <div className="text-center max-w-md mx-auto">
          <h1 className="text-2xl font-bold text-red-600 mb-4">ログインエラー</h1>
          <p className="text-gray-600 mb-6">{error}</p>
          <Link
            href="/articles"
            className="bg-blue-500 hover:bg-blue-600 text-white px-6 py-3 rounded-lg font-medium inline-block"
          >
            記事一覧に戻る
          </Link>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50 py-8">
      <div className="max-w-2xl mx-auto px-4">
        <div className="bg-white rounded-lg shadow-sm p-6">
          <h1 className="text-2xl font-bold text-green-600 mb-6">
            {loginSuccess ? 'ログイン成功！' : 'LINEログイン成功！'}
          </h1>
          
          {loginSuccess && (
            <div className="bg-green-50 border border-green-200 rounded-lg p-4 mb-6">
              <p className="text-green-800">
                ログインが完了しました。3秒後に記事一覧ページに移動します...
              </p>
            </div>
          )}
          
          {result && (
            <div className="space-y-4">
              <div className="bg-green-50 border border-green-200 rounded-lg p-4">
                <h2 className="font-semibold text-green-800 mb-2">認証結果</h2>
                <p className="text-green-700">LINEログインが正常に完了しました</p>
              </div>

              {result.user && (
                <div className="bg-gray-50 border border-gray-200 rounded-lg p-4">
                  <h2 className="font-semibold text-gray-800 mb-2">ユーザー情報</h2>
                  <div className="space-y-2">
                    <p><span className="font-medium">表示名:</span> {result.user.name}</p>
                    <p><span className="font-medium">ユーザーID:</span> {result.user.id}</p>
                    <p><span className="font-medium">LINE ID:</span> {result.user.line_user_id}</p>
                    {result.user.profile_image && (
                      <div>
                        <span className="font-medium">プロフィール画像:</span>
                        <img 
                          src={result.user.profile_image} 
                          alt="プロフィール画像" 
                          className="w-16 h-16 rounded-full mt-2"
                        />
                      </div>
                    )}
                  </div>
                </div>
              )}

              <div className="bg-blue-50 border border-blue-200 rounded-lg p-4">
                <h2 className="font-semibold text-blue-800 mb-2">デバッグ情報</h2>
                <pre className="text-sm text-blue-700 overflow-x-auto">
                  {JSON.stringify(result, null, 2)}
                </pre>
              </div>
            </div>
          )}

          <div className="mt-6 flex gap-4">
            <Link
              href="/articles"
              className="bg-blue-500 hover:bg-blue-600 text-white px-6 py-3 rounded-lg font-medium"
            >
              記事一覧に戻る
            </Link>
            <button
              onClick={() => router.push('/articles')}
              className="bg-gray-500 hover:bg-gray-600 text-white px-6 py-3 rounded-lg font-medium"
            >
              ホームに戻る
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}