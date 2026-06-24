export default function Header() {
  return (
    <header className="fixed top-0 left-0 right-0 z-50 py-4 bg-zinc-950/80 backdrop-blur-lg border-b border-zinc-800">
      <div className="max-w-4xl mx-auto px-6 flex justify-between items-center">
        <a href="/" className="text-xl font-bold text-white flex items-center gap-2 hover:text-white transition-colors">
          <span className="text-blue-400">⬡</span> Hatch
        </a>
        <nav className="flex items-center gap-6">
          <a
            href="https://github.com/elfoundation/hatch"
            className="text-zinc-400 hover:text-white transition-colors"
          >
            GitHub
          </a>
          <a
            href="/blog/"
            className="text-zinc-400 hover:text-white transition-colors"
          >
            Blog
          </a>
        </nav>
      </div>
    </header>
  )
}
