import Link from 'next/link';

export default function Home() {
  return (
    <div className="min-h-screen grid grid-rows-[auto_1fr_auto] p-8">
      <header className="text-center py-8">
        <h1 className="text-4xl font-bold">sg3t41</h1>
      </header>
      <main className="grid place-items-center">
        <div className="text-center">
          <div className="mt-8">
            <Link
              href="/articles"
              className="bg-green-500 hover:bg-green-600 text-white px-6 py-3 rounded-lg font-medium inline-block"
            >
              ブログ
            </Link>
          </div>
        </div>
      </main>
      <footer className="text-center py-4 text-sm text-gray-500">
        <p>© 2025 sg3t41</p>
      </footer>
    </div>
  );
}