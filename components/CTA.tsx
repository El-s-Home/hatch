export default function CTA() {
  return (
    <section className="py-24 px-6 bg-gradient-to-b from-zinc-900/30 to-zinc-900/60">
      <div className="max-w-4xl mx-auto text-center">
        <h2 className="text-4xl md:text-5xl font-extrabold mb-4 bg-gradient-to-r from-white to-blue-400 bg-clip-text text-transparent">
          Start inspecting in 30&nbsp;seconds
        </h2>
        <p className="text-xl text-zinc-400 mb-10">
          One binary. Your data stays yours.
        </p>
        <a
          href="https://github.com/elfoundation/hatch"
          className="inline-flex items-center justify-center gap-2 px-10 py-5 bg-blue-500 hover:bg-blue-600 text-white font-semibold rounded-lg transition-all shadow-lg shadow-blue-500/30 hover:shadow-xl hover:shadow-blue-500/40 hover:-translate-y-0.5 text-lg"
        >
          Get Hatch →
        </a>
      </div>
    </section>
  )
}
