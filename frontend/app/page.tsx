export default function Home() {
  return (
    <div className="min-h-screen grid grid-rows-[auto_1fr_auto] p-8">
      <header className="text-center py-8">
        <h1 className="text-4xl font-bold">Utils App</h1>
        <p className="text-gray-600 mt-2">API + Frontend + PostgreSQL</p>
      </header>
      <main className="grid place-items-center">
        <div className="text-center">
          <p className="text-xl mb-4">Welcome to Utils Application</p>
          <div className="grid gap-4 mt-8">
            <a
              href="/users"
              className="bg-blue-500 hover:bg-blue-600 text-white px-6 py-3 rounded-lg font-medium"
            >
              View Users
            </a>
          </div>
        </div>
      </main>
      <footer className="text-center py-4 text-sm text-gray-500">
        <p>© 2025 Utils App</p>
      </footer>
    </div>
  );
}