import Link from 'next/link';

interface TileData {
  title: string;
  href: string;
  color: string;
  icon: string;
  size: 'small' | 'medium' | 'large';
  description?: string;
}

const tiles: TileData[] = [
  {
    title: 'ブログ',
    href: '/articles',
    color: 'bg-blue-500 hover:bg-blue-600',
    icon: '📝',
    size: 'large',
    description: '記事・投稿管理'
  },
  {
    title: '制作物',
    href: '/projects',
    color: 'bg-purple-500 hover:bg-purple-600',
    icon: '🛠️',
    size: 'medium',
    description: 'プロジェクト管理'
  },
  {
    title: '収支表',
    href: '/finance',
    color: 'bg-green-500 hover:bg-green-600',
    icon: '💰',
    size: 'medium',
    description: '家計・収支管理'
  },
  {
    title: 'ランキング',
    href: '/ranking',
    color: 'bg-red-500 hover:bg-red-600',
    icon: '🏆',
    size: 'medium',
    description: 'ランキング管理'
  },
  {
    title: '共有',
    href: '/shared',
    color: 'bg-orange-500 hover:bg-orange-600',
    icon: '🔗',
    size: 'small',
    description: 'ファイル共有'
  },
  {
    title: '写真',
    href: '/photos',
    color: 'bg-pink-500 hover:bg-pink-600',
    icon: '📸',
    size: 'small',
    description: 'フォトギャラリー'
  },
  {
    title: 'メモ',
    href: '/notes',
    color: 'bg-yellow-500 hover:bg-yellow-600',
    icon: '📔',
    size: 'small',
    description: '簡単メモ'
  },
  {
    title: 'リンク',
    href: '/links',
    color: 'bg-indigo-500 hover:bg-indigo-600',
    icon: '🌐',
    size: 'small',
    description: 'リンク集'
  }
];

function Tile({ tile }: { tile: TileData }) {
  const getSizeClasses = (size: string) => {
    switch (size) {
      case 'small':
        return 'col-span-1 row-span-1 h-28 sm:h-32';
      case 'medium':
        return 'col-span-2 row-span-1 h-28 sm:col-span-1 sm:row-span-2 sm:h-64';
      case 'large':
        return 'col-span-2 row-span-2 h-56 sm:h-64';
      default:
        return 'col-span-1 row-span-1 h-32';
    }
  };

  return (
    <Link
      href={tile.href}
      className={`${tile.color} ${getSizeClasses(tile.size)} rounded-lg text-white p-3 sm:p-6 transition-all duration-200 transform hover:scale-105 hover:shadow-xl flex flex-col justify-between group relative overflow-hidden`}
    >
      <div className="flex-1">
        <div className="text-2xl sm:text-3xl mb-1 sm:mb-2">{tile.icon}</div>
        <h3 className="text-lg sm:text-xl font-bold mb-1">{tile.title}</h3>
        {tile.description && (
          <p className="text-xs sm:text-sm opacity-90 hidden sm:block">{tile.description}</p>
        )}
      </div>
      <div className="absolute bottom-3 right-3 sm:bottom-6 sm:right-6 opacity-0 group-hover:opacity-100 transition-opacity">
        <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
        </svg>
      </div>
    </Link>
  );
}

export default function Home() {
  return (
    <div className="min-h-screen bg-gray-50 pt-24 pb-8">
      <div className="max-w-6xl mx-auto px-4">
        <header className="text-center py-6 sm:py-8">
          <h1 className="text-3xl sm:text-4xl font-bold text-gray-900">sg3t41</h1>
        </header>
        
        <main className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 gap-3 sm:gap-4 auto-rows-fr">
          {tiles.map((tile, index) => (
            <Tile key={index} tile={tile} />
          ))}
        </main>
      </div>
    </div>
  );
}