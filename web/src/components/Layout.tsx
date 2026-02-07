import { Outlet, Link, useLocation } from 'react-router-dom';

export default function Layout() {
  const location = useLocation();

  return (
    <div className="min-h-screen flex flex-col">
      {/* Header */}
      <header className="bg-[var(--bg-secondary)] border-b border-[var(--border)] px-6 py-3 flex items-center justify-between">
        <Link to="/" className="flex items-center gap-3 no-underline">
          <span className="text-2xl">🚀</span>
          <h1 className="text-lg font-bold text-[var(--text-primary)] m-0">
            API Tester
          </h1>
        </Link>
        <nav className="flex gap-4">
          <Link
            to="/"
            className={`text-sm no-underline px-3 py-1.5 rounded-md transition-colors ${
              location.pathname === '/'
                ? 'bg-[var(--accent)] text-white'
                : 'text-[var(--text-secondary)] hover:text-[var(--text-primary)]'
            }`}
          >
            Dashboard
          </Link>
        </nav>
      </header>

      {/* Main content */}
      <main className="flex-1 p-6 max-w-7xl mx-auto w-full">
        <Outlet />
      </main>
    </div>
  );
}
