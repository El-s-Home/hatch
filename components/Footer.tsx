export default function Footer() {
  return (
    <footer className="py-8 px-6 border-t border-zinc-800">
      <div className="max-w-4xl mx-auto">
        <div className="flex flex-col md:flex-row justify-between items-center gap-4">
          <span className="text-white font-bold flex items-center gap-2">
            <span className="text-blue-400">⬡</span> Hatch
          </span>
          <p className="text-zinc-500 text-sm text-center md:text-right">
            © 2026 El Foundation.{' '}
            <a
              href="https://github.com/elfoundation/hatch"
              className="text-zinc-400 hover:text-white transition-colors"
            >
              Hatch
            </a>{' '}
            is released under the{' '}
            <a
              href="https://github.com/elfoundation/hatch/blob/main/LICENSE"
              className="text-zinc-400 hover:text-white transition-colors"
            >
              MIT License
            </a>.
          </p>
        </div>
      </div>
    </footer>
  )
}
