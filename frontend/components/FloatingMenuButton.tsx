'use client';

import { useState } from 'react';
import Link from 'next/link';
import { useAuth } from '../contexts/AuthContext';
import { useScrollHeader } from '../hooks/useScrollHeader';

export default function FloatingMenuButton() {
  const { user, isAuthenticated, isAdmin, logout, isLoading } = useAuth();
  const [isMenuOpen, setIsMenuOpen] = useState(false);
  const [isAnimating, setIsAnimating] = useState(false);
  const isVisible = useScrollHeader();

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
    closeMenu();
  };

  const openMenu = () => {
    setIsMenuOpen(true);
    setTimeout(() => setIsAnimating(true), 10);
  };

  const closeMenu = () => {
    setIsAnimating(false);
    setTimeout(() => setIsMenuOpen(false), 300);
  };

  if (isLoading) return null;

  return (
    <>
      {/* ハンバーガーメニューボタン */}
      <div className="fixed top-4 left-0 right-0 z-50 pointer-events-none">
        <div className="max-w-6xl mx-auto px-4">
          <button
            onClick={openMenu}
            className={`p-3 bg-white/90 backdrop-blur-sm rounded-full shadow-lg border border-gray-200 hover:bg-white hover:shadow-xl transition-all duration-300 pointer-events-auto ${
              isVisible ? 'translate-y-0 opacity-100' : '-translate-y-16 opacity-0'
            }`}
            aria-label="メニューを開く"
          >
            <svg className="w-6 h-6 text-gray-700" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16M4 18h16" />
            </svg>
          </button>
        </div>
      </div>

      {/* メニューオーバーレイ */}
      {isMenuOpen && (
        <>
          <div
            className="fixed inset-0 bg-black/50 z-50 transition-opacity duration-300"
            onClick={closeMenu}
          />
          
          {/* スライドメニュー */}
          <div
            className={`fixed top-0 left-0 h-full w-80 bg-white shadow-2xl z-50 transform transition-transform duration-300 ease-out ${
              isAnimating ? 'translate-x-0' : '-translate-x-full'
            }`}
          >
            {/* メニューヘッダー */}
            <div className="flex items-center justify-between p-6 border-b border-gray-200">
              <h2 className="text-xl font-bold text-gray-900">メニュー</h2>
              <button
                onClick={closeMenu}
                className="p-2 rounded-lg hover:bg-gray-100 transition-colors"
                aria-label="メニューを閉じる"
              >
                <svg className="w-6 h-6 text-gray-700" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>

            {/* メニューコンテンツ */}
            <div className="p-6">
              {/* ユーザー情報 */}
              {isAuthenticated && (
                <div className="mb-6 pb-6 border-b border-gray-200">
                  <div className="flex items-center gap-3">
                    {user?.profile_image && (
                      <img
                        src={user.profile_image}
                        alt={user.name}
                        className="w-12 h-12 rounded-full"
                      />
                    )}
                    <div>
                      <p className="font-medium text-gray-900">{user?.name}</p>
                      {isAdmin && (
                        <p className="text-sm text-blue-600">管理者</p>
                      )}
                    </div>
                  </div>
                </div>
              )}

              {/* ナビゲーションリンク */}
              <nav className="space-y-1">
                <Link
                  href="/"
                  onClick={closeMenu}
                  className="flex items-center justify-between px-4 py-3 rounded-lg text-gray-700 hover:bg-blue-50 hover:text-blue-600 transition-colors border border-transparent hover:border-blue-200"
                >
                  <div className="flex items-center gap-3">
                    <span className="text-lg">🏠</span>
                    <span>ホーム</span>
                  </div>
                  <svg className="w-4 h-4 opacity-50" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                  </svg>
                </Link>
                <Link
                  href="/articles"
                  onClick={closeMenu}
                  className="flex items-center justify-between px-4 py-3 rounded-lg text-gray-700 hover:bg-blue-50 hover:text-blue-600 transition-colors border border-transparent hover:border-blue-200"
                >
                  <div className="flex items-center gap-3">
                    <span className="text-lg">📝</span>
                    <span>ブログ</span>
                  </div>
                  <svg className="w-4 h-4 opacity-50" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                  </svg>
                </Link>
                <Link
                  href="/projects"
                  onClick={closeMenu}
                  className="flex items-center justify-between px-4 py-3 rounded-lg text-gray-700 hover:bg-purple-50 hover:text-purple-600 transition-colors border border-transparent hover:border-purple-200"
                >
                  <div className="flex items-center gap-3">
                    <span className="text-lg">🛠️</span>
                    <span>制作物</span>
                  </div>
                  <svg className="w-4 h-4 opacity-50" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                  </svg>
                </Link>
                <Link
                  href="/finance"
                  onClick={closeMenu}
                  className="flex items-center justify-between px-4 py-3 rounded-lg text-gray-700 hover:bg-green-50 hover:text-green-600 transition-colors border border-transparent hover:border-green-200"
                >
                  <div className="flex items-center gap-3">
                    <span className="text-lg">💰</span>
                    <span>収支表</span>
                  </div>
                  <svg className="w-4 h-4 opacity-50" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                  </svg>
                </Link>
                <Link
                  href="/shared"
                  onClick={closeMenu}
                  className="flex items-center justify-between px-4 py-3 rounded-lg text-gray-700 hover:bg-orange-50 hover:text-orange-600 transition-colors border border-transparent hover:border-orange-200"
                >
                  <div className="flex items-center gap-3">
                    <span className="text-lg">🔗</span>
                    <span>共有</span>
                  </div>
                  <svg className="w-4 h-4 opacity-50" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                  </svg>
                </Link>
                <Link
                  href="/photos"
                  onClick={closeMenu}
                  className="flex items-center justify-between px-4 py-3 rounded-lg text-gray-700 hover:bg-pink-50 hover:text-pink-600 transition-colors border border-transparent hover:border-pink-200"
                >
                  <div className="flex items-center gap-3">
                    <span className="text-lg">📸</span>
                    <span>写真</span>
                  </div>
                  <svg className="w-4 h-4 opacity-50" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                  </svg>
                </Link>
                <Link
                  href="/notes"
                  onClick={closeMenu}
                  className="flex items-center justify-between px-4 py-3 rounded-lg text-gray-700 hover:bg-yellow-50 hover:text-yellow-600 transition-colors border border-transparent hover:border-yellow-200"
                >
                  <div className="flex items-center gap-3">
                    <span className="text-lg">📔</span>
                    <span>メモ</span>
                  </div>
                  <svg className="w-4 h-4 opacity-50" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                  </svg>
                </Link>
                <Link
                  href="/ranking"
                  onClick={closeMenu}
                  className="flex items-center justify-between px-4 py-3 rounded-lg text-gray-700 hover:bg-red-50 hover:text-red-600 transition-colors border border-transparent hover:border-red-200"
                >
                  <div className="flex items-center gap-3">
                    <span className="text-lg">🏆</span>
                    <span>ランキング</span>
                  </div>
                  <svg className="w-4 h-4 opacity-50" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                  </svg>
                </Link>
                <Link
                  href="/links"
                  onClick={closeMenu}
                  className="flex items-center justify-between px-4 py-3 rounded-lg text-gray-700 hover:bg-indigo-50 hover:text-indigo-600 transition-colors border border-transparent hover:border-indigo-200"
                >
                  <div className="flex items-center gap-3">
                    <span className="text-lg">🌐</span>
                    <span>リンク</span>
                  </div>
                  <svg className="w-4 h-4 opacity-50" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                  </svg>
                </Link>
                
                {isAdmin && (
                  <>
                    <div className="my-4 border-t border-gray-200"></div>
                    <p className="px-4 py-2 text-sm font-medium text-gray-500">管理機能</p>
                    <Link
                      href="/articles/new"
                      onClick={closeMenu}
                      className="flex items-center justify-between px-4 py-3 rounded-lg text-gray-700 hover:bg-blue-50 hover:text-blue-600 transition-colors border border-transparent hover:border-blue-200"
                    >
                      <div className="flex items-center gap-3">
                        <span className="text-lg">✏️</span>
                        <span>新規記事作成</span>
                      </div>
                      <svg className="w-4 h-4 opacity-50" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                      </svg>
                    </Link>
                  </>
                )}
              </nav>

              {/* アクションボタン */}
              <div className="mt-8">
                {isAuthenticated ? (
                  <button
                    onClick={handleLogout}
                    className="w-full bg-red-500 hover:bg-red-600 text-white px-4 py-3 rounded-lg font-medium transition-colors"
                  >
                    ログアウト
                  </button>
                ) : (
                  <button
                    onClick={handleLineLogin}
                    className="w-full bg-green-500 hover:bg-green-600 text-white px-4 py-3 rounded-lg font-medium transition-colors"
                  >
                    ログイン
                  </button>
                )}
              </div>
            </div>
          </div>
        </>
      )}
    </>
  );
}