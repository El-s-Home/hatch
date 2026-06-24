const features = [
  {
    icon: (
      <svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round">
        <circle cx="12" cy="12" r="10"></circle>
        <polyline points="12 6 12 12 16 14"></polyline>
      </svg>
    ),
    title: 'Capture',
    description: 'Method, path, headers, query, body. Persists across restarts because the storage is SQLite on disk, not a hosted queue.',
    color: 'blue',
  },
  {
    icon: (
      <svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round">
        <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"></path>
        <circle cx="12" cy="12" r="3"></circle>
      </svg>
    ),
    title: 'Inspect',
    description: 'A live SSE feed of incoming requests. Click any captured request to see the headers, the body, the timing.',
    color: 'purple',
  },
  {
    icon: (
      <svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round">
        <polyline points="16 18 22 12 16 6"></polyline>
        <polyline points="8 6 2 12 8 18"></polyline>
      </svg>
    ),
    title: 'Mock',
    description: 'Return a 200, a 500, or a custom JSON payload. For testing your own retry, backoff, and error-handling logic.',
    color: 'cyan',
  },
]

const colorClasses = {
  blue: {
    icon: 'bg-blue-500/15 text-blue-400',
    hover: 'hover:border-blue-500/40 hover:shadow-blue-500/10',
  },
  purple: {
    icon: 'bg-purple-500/15 text-purple-400',
    hover: 'hover:border-purple-500/40 hover:shadow-purple-500/10',
  },
  cyan: {
    icon: 'bg-cyan-500/15 text-cyan-400',
    hover: 'hover:border-cyan-500/40 hover:shadow-cyan-500/10',
  },
}

export default function Features() {
  return (
    <section className="py-20 px-6">
      <div className="max-w-4xl mx-auto">
        <h2 className="text-3xl md:text-4xl font-bold text-center mb-16">
          Three things, and nothing else
        </h2>
        
        <div className="grid md:grid-cols-3 gap-8">
          {features.map((feature) => {
            const colors = colorClasses[feature.color as keyof typeof colorClasses]
            return (
              <div
                key={feature.title}
                className={`p-6 rounded-xl bg-gradient-to-b from-zinc-900/80 to-zinc-900/40 border border-zinc-800 ${colors.hover} transition-all hover:-translate-y-1 shadow-lg hover:shadow-xl`}
              >
                <div className={`w-14 h-14 flex items-center justify-center rounded-xl ${colors.icon} mb-4 shadow-lg`}>
                  {feature.icon}
                </div>
                <h3 className="text-xl font-semibold mb-2">{feature.title}</h3>
                <p className="text-zinc-400 leading-relaxed">
                  {feature.description}
                </p>
              </div>
            )
          })}
        </div>
      </div>
    </section>
  )
}
