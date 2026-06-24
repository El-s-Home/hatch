const reasons = [
  {
    number: '01',
    title: 'Compliance and privacy',
    description: 'Some teams cannot legally send webhook payloads to a hosted SaaS. Hatch keeps the data on your own network.',
  },
  {
    number: '02',
    title: 'Cost',
    description: 'A hosted inspector charges per request, per seat, or per retention day. Hatch is one Go binary on a $5 VPS. There is no per-request fee because there is no one to charge it.',
  },
  {
    number: '03',
    title: 'Speed of setup',
    description: (
      <>
        <code className="px-2 py-0.5 bg-zinc-800 border border-zinc-700 rounded text-fuchsia-400 text-sm">
          docker compose up
        </code>{' '}
        is faster than signing up for a SaaS, verifying your email, configuring your first bin, and pasting the URL into your webhook config.
      </>
    ),
  },
]

export default function Why() {
  return (
    <section className="py-20 px-6 bg-zinc-900/30">
      <div className="max-w-4xl mx-auto">
        <h2 className="text-3xl md:text-4xl font-bold text-center mb-16">
          Why a single binary, not a SaaS
        </h2>
        
        <div className="space-y-8">
          {reasons.map((reason) => (
            <div
              key={reason.number}
              className="flex gap-6 items-start"
            >
              <div className="text-4xl font-bold text-zinc-700 font-mono shrink-0">
                {reason.number}
              </div>
              <div>
                <h3 className="text-xl font-semibold mb-2">{reason.title}</h3>
                <p className="text-zinc-400 leading-relaxed">
                  {reason.description}
                </p>
              </div>
            </div>
          ))}
        </div>
      </div>
    </section>
  )
}
