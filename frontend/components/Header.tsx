'use client';

import { useAuth } from '../contexts/AuthContext';
import Link from 'next/link';

export default function Header() {
  const { user, isAuthenticated, logout, isLoading } = useAuth();
  
  console.log('Header render:', { user, isAuthenticated, isLoading });

  const handleLineLogin = async () => {
    try {
      const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';
      const response = await fetch(`${apiUrl}/api/v1/auth/line/url`);
      
      if (!response.ok) {
        throw new Error('Failed to get LINE auth URL');
      }

      const data = await response.json();
      window.location.href = data.auth_url;
    } catch (error) {
      console.error('LINE login error:', error);
      alert('LINEログインでエラーが発生しました');
    }
  };

  const handleLogout = () => {
    logout();
    alert('ログアウトしました');
  };

  if (isLoading) {
    return (
      <header className="bg-white border-b border-gray-200 px-4 py-3">
        <div className="max-w-6xl mx-auto flex justify-between items-center">
          <Link href="/" className="text-xl font-bold text-gray-900">
            ブログ管理システム
          </Link>
          <div className="animate-pulse bg-gray-200 h-8 w-24 rounded"></div>
        </div>
      </header>
    );
  }

  return (
    <header className="bg-white border-b border-gray-200 px-4 py-3">
      <div className="max-w-6xl mx-auto flex justify-between items-center">
        <Link href="/" className="text-xl font-bold text-gray-900">
          ブログ管理システム
        </Link>
        
        <nav className="flex items-center gap-4">
          <Link 
            href="/articles" 
            className="text-gray-600 hover:text-gray-900 font-medium"
          >
            記事一覧
          </Link>
          
          {isAuthenticated ? (
            <div className="flex items-center gap-4">
              <div className="flex items-center gap-2">
                {user?.profile_image && (
                  <img
                    src={user.profile_image}
                    alt={user.name}
                    className="w-8 h-8 rounded-full"
                  />
                )}
                <span className="text-sm text-gray-700">
                  {user?.name}でログイン中
                </span>
              </div>
              <button
                onClick={handleLogout}
                className="bg-red-500 hover:bg-red-600 text-white px-3 py-1 rounded text-sm font-medium"
              >
                ログアウト
              </button>
            </div>
          ) : (
            <button
              onClick={handleLineLogin}
              className="bg-green-500 hover:bg-green-600 text-white px-4 py-2 rounded-lg font-medium"
            >
              LINEログイン
            </button>
          )}
        </nav>
      </div>
    </header>
  );
}