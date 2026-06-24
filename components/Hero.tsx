export default function Hero() {
  return (
    <section className="pt-32 pb-20 px-6">
      <div className="max-w-4xl mx-auto text-center">
        <div className="inline-block mb-6 px-4 py-2 rounded-full bg-blue-500/10 border border-blue-500/20 text-blue-400 text-sm font-medium">
          Open source · MIT licensed
        </div>
        
        <h1 className="text-5xl md:text-6xl font-bold tracking-tight mb-6 leading-tight">
          Self-hostable HTTP request<br />
          inspector&nbsp;+&nbsp;mocker
        </h1>
        
        <p className="text-xl text-zinc-400 mb-8 max-w-2xl mx-auto">
          One Go binary. SQLite under the hood.
          <code className="mx-1 px-2 py-0.5 bg-zinc-800 border border-zinc-700 rounded text-fuchsia-400 text-base">
            docker compose up
          </code>
          , and you have an inspection endpoint and a live feed in under 30 seconds.
        </p>
        
        <div className="flex flex-col sm:flex-row gap-4 justify-center mb-12">
          <a
            href="https://github.com/elfoundation/hatch"
            className="inline-flex items-center justify-center gap-2 px-6 py-3 bg-blue-500 hover:bg-blue-600 text-white font-medium rounded-lg transition-colors"
          >
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
              <path d="M9 19c-5 1.5-5-2.5-7-3m14 6v-3.87a3.37 3.37 0 0 0-.94-2.61c3.14-.35 6.44-1.54 6.44-7A5.44 5.44 0 0 0 20 4.77 5.07 5.07 0 0 0 19.91 1S18.73.65 16 2.48a13.38 13.38 0 0 0-7 0C6.27.65 5.09 1 5.09 1A5.07 5.07 0 0 0 5 4.77a5.44 5.44 0 0 0-1.5 3.78c0 5.42 3.3 6.61 6.44 7A3.37 3.37 0 0 0 9 18.13V22"></path>
            </svg>
            View on GitHub
          </a>
          <a
            href="https://github.com/elfoundation/hatch#quick-start"
            className="inline-flex items-center justify-center gap-2 px-6 py-3 bg-zinc-800 hover:bg-zinc-700 text-white font-medium rounded-lg border border-zinc-700 transition-colors"
          >
            Quick Start →
          </a>
        </div>
        
        {/* Terminal */}
        <div className="max-w-2xl mx-auto rounded-xl overflow-hidden border border-zinc-800 bg-zinc-900/50 backdrop-blur-sm">
          <div className="flex items-center gap-2 px-4 py-3 bg-zinc-800/50 border-b border-zinc-800">
            <div className="w-3 h-3 rounded-full bg-zinc-600"></div>
            <div className="w-3 h-3 rounded-full bg-zinc-600"></div>
            <div className="w-3 h-3 rounded-full bg-zinc-600"></div>
            <span className="ml-2 text-sm text-zinc-500">terminal</span>
          </div>
          <pre className="p-4 text-left text-sm font-mono overflow-x-auto">
            <code>
              <span className="text-green-400">$</span>{' '}
              <span className="text-zinc-300">docker compose up -d</span>{'\n'}
              <span className="text-zinc-500">✓ hatch  started on :8080</span>{'\n'}
              <span className="text-green-400">$</span>{' '}
              <span className="text-zinc-300">curl -X POST https://your-bin.hatch.surf/test -d {'{\'hello\':\'world\'}'}</span>{'\n'}
              <span className="text-zinc-500">✓ captured — view at https://your-bin.hatch.surf/inspect</span>
            </code>
          </pre>
        </div>
      </div>
    </section>
  )
}
