'use client';

import { useState, useEffect } from 'react';
import { FiTwitter, FiInstagram, FiGithub, FiExternalLink, FiYoutube, FiLinkedin } from 'react-icons/fi';
import { SiLine, SiTiktok } from 'react-icons/si';

interface Link {
  id: number;
  title: string;
  url: string;
  description?: string;
  platform: string;
  icon_name?: string;
  background_color: string;
  text_color: string;
  order_index: number;
  is_active: boolean;
}

interface LinkResponse {
  links: Link[];
  total: number;
}

const platformIcons: Record<string, React.ComponentType<{ className?: string }>> = {
  twitter: FiTwitter,
  instagram: FiInstagram,
  github: FiGithub,
  line: SiLine,
  youtube: FiYoutube,
  tiktok: SiTiktok,
  linkedin: FiLinkedin,
  website: FiExternalLink,
};

export default function LinksPage() {
  const [links, setLinks] = useState<Link[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    fetchLinks();
  }, []);

  const fetchLinks = async () => {
    try {
      const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/v1/links?active=true`);
      if (!response.ok) {
        throw new Error('リンクの取得に失敗しました');
      }
      const data: LinkResponse = await response.json();
      setLinks(data.links);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'エラーが発生しました');
    } finally {
      setLoading(false);
    }
  };

  const getIconComponent = (platform: string) => {
    return platformIcons[platform] || FiExternalLink;
  };

  const handleLinkClick = (url: string) => {
    window.open(url, '_blank', 'noopener,noreferrer');
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-purple-100 via-pink-50 to-orange-100 flex items-center justify-center">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-purple-600"></div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-purple-100 via-pink-50 to-orange-100 flex items-center justify-center">
        <div className="text-center p-8 bg-white rounded-xl shadow-lg">
          <div className="text-red-500 text-lg font-medium mb-2">エラー</div>
          <div className="text-gray-600">{error}</div>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-purple-100 via-pink-50 to-orange-100 py-8 px-4">
      <div className="max-w-md mx-auto">
        {/* プロフィールヘッダー */}
        <div className="text-center mb-8">
          <h1 className="text-2xl font-bold text-gray-800 mb-2">sg3t41</h1>
        </div>

        {/* リンクリスト */}
        <div className="space-y-4">
          {links.length === 0 ? (
            <div className="text-center py-12">
              <div className="text-gray-500 text-lg">リンクがありません</div>
            </div>
          ) : (
            links.map((link) => {
              const IconComponent = getIconComponent(link.platform);
              return (
                <button
                  key={link.id}
                  onClick={() => handleLinkClick(link.url)}
                  className="w-full p-4 rounded-2xl shadow-lg hover:shadow-xl transform hover:-translate-y-1 transition-all duration-200 text-left group"
                  style={{
                    backgroundColor: link.background_color,
                    color: link.text_color,
                  }}
                >
                  <div className="flex items-center space-x-4">
                    <div className="flex-shrink-0">
                      <IconComponent className="w-6 h-6" />
                    </div>
                    <div className="flex-1 min-w-0">
                      <div className="font-semibold text-lg truncate">
                        {link.title}
                      </div>
                      {link.description && (
                        <div className="text-sm opacity-90 truncate mt-1">
                          {link.description}
                        </div>
                      )}
                    </div>
                    <div className="flex-shrink-0 opacity-70 group-hover:opacity-100 transition-opacity">
                      <FiExternalLink className="w-5 h-5" />
                    </div>
                  </div>
                </button>
              );
            })
          )}
        </div>

        {/* フッター */}
        <div className="text-center mt-12 text-gray-500 text-sm">
          <p>© 2025 sg3t41. All rights reserved.</p>
        </div>
      </div>
    </div>
  );
}