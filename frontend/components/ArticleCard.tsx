import Link from 'next/link';
import { Article } from '../types/article';
import { formatDateTime } from '../utils/dateFormat';
import { getStatusBadge, getStatusText } from '../utils/statusUtils';

interface ArticleCardProps {
  article: Article;
  mounted: boolean;
}

/**
 * 記事カードコンポーネント
 * 記事の基本情報とサムネイルを表示
 */
export default function ArticleCard({ article, mounted }: ArticleCardProps) {
  return (
    <div className="border border-gray-200 rounded-lg p-6 hover:shadow-md transition-shadow">
      <div className="grid gap-4 md:flex md:gap-4">
        {/* サムネイル */}
        <div className="w-full h-48 md:w-40 md:h-28 md:flex-shrink-0">
          <img
            src={mounted ? (article.article_image 
              ? `${process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'}/api/v1/uploads/${article.article_image}` 
              : '/noimage.svg'
            ) : '/noimage.svg'}
            alt={article.title}
            className="w-full h-full object-cover rounded-lg shadow-sm bg-gray-100"
            onError={(e) => {
              e.currentTarget.src = '/noimage.svg';
            }}
          />
        </div>

        {/* 記事情報 */}
        <div className="flex-1 grid gap-3">
          {/* タイトルとステータス */}
          <div className="flex justify-between items-start gap-4">
            <h2 className="text-xl font-semibold text-gray-900 hover:text-blue-600">
              <Link href={`/articles/${article.id}`}>{article.title}</Link>
            </h2>
            <span className={getStatusBadge(article.status)}>
              {getStatusText(article.status)}
            </span>
          </div>
        
          {/* 概要 */}
          {article.summary && (
            <p className="text-gray-600 line-clamp-2">{article.summary}</p>
          )}
        
          {/* タグ */}
          {article.tags.length > 0 && (
            <div className="flex flex-wrap gap-2">
              {article.tags.map((tag, index) => (
                <span
                  key={index}
                  className="px-2 py-1 bg-blue-100 text-blue-800 text-xs rounded-full"
                >
                  {tag}
                </span>
              ))}
            </div>
          )}
        
          {/* 日付情報 */}
          <div className="flex justify-between items-center text-sm text-gray-500">
            <span>作成: {formatDateTime(article.created_at)}</span>
            {article.published_at && (
              <span>公開: {formatDateTime(article.published_at)}</span>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}