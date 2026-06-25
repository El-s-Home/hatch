import { useState } from "react";
import type { HatchRequest, ReplayResult } from "./types";
import { relativeTime, methodStyle, statusStyle, decodeBody, prettyJSON } from "./format";

interface Props {
  req: HatchRequest;
  apiBase?: string;
}

function Section({ label, children }: { label: string; children: React.ReactNode }) {
  return (
    <div className="mb-3 last:mb-0">
      <div className="mb-1 text-[0.7rem] font-semibold uppercase tracking-wide text-zinc-500">{label}</div>
      <pre className="overflow-x-auto whitespace-pre-wrap break-all rounded-md border border-zinc-800 bg-zinc-950 p-3 text-xs text-zinc-300">
        {children}
      </pre>
    </div>
  );
}

export function RequestRow({ req, apiBase = "" }: Props) {
  const [open, setOpen] = useState(false);
  const [replayOpen, setReplayOpen] = useState(false);
  const [target, setTarget] = useState("");
  const [sending, setSending] = useState(false);
  const [error, setError] = useState("");
  const [result, setResult] = useState<ReplayResult | null>(null);

  const body = decodeBody(req.body);

  async function doReplay() {
    if (!target.trim()) {
      setError("Enter a target URL.");
      return;
    }
    setSending(true);
    setError("");
    setResult(null);
    try {
      const resp = await fetch(
        `${apiBase}/e/${encodeURIComponent(req.endpoint_id)}/requests/${req.id}/replay`,
        {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ target_url: target.trim() }),
        },
      );
      const data: ReplayResult = await resp.json();
      if (data.error) setError(data.error);
      else setResult(data);
    } catch (e) {
      setError(`Network error: ${(e as Error).message}`);
    } finally {
      setSending(false);
    }
  }

  return (
    <div className="overflow-hidden rounded-lg border border-zinc-800 bg-zinc-900/50">
      <div className="flex items-center gap-3 px-4 py-3">
        <button
          onClick={() => setOpen((v) => !v)}
          className="flex flex-1 items-center gap-3 text-left"
        >
          <span
            className={`rounded px-2 py-0.5 text-xs font-bold ring-1 ring-inset ${methodStyle(req.method)}`}
          >
            {req.method}
          </span>
          <span className="flex-1 truncate font-mono text-sm text-zinc-200">{req.path}</span>
          <span className="whitespace-nowrap text-xs text-zinc-500" title={req.created_at}>
            {relativeTime(req.created_at)}
          </span>
        </button>
        <button
          onClick={() => setReplayOpen((v) => !v)}
          className="rounded-md bg-indigo-600 px-3 py-1 text-xs font-medium text-white hover:bg-indigo-500"
        >
          ↺ Replay
        </button>
      </div>

      {open && (
        <div className="border-t border-zinc-800 px-4 py-3">
          {req.headers && req.headers !== "{}" && (
            <Section label="Headers">{prettyJSON(req.headers)}</Section>
          )}
          {req.query && <Section label="Query">{req.query}</Section>}
          {body && <Section label="Body">{body}</Section>}
        </div>
      )}

      {replayOpen && (
        <div className="border-t border-zinc-800 px-4 py-3">
          <input
            value={target}
            onChange={(e) => setTarget(e.target.value)}
            placeholder="https://myapp.local/webhook"
            className="mb-2 w-full rounded-md border border-zinc-700 bg-zinc-950 px-3 py-2 font-mono text-sm text-zinc-200 outline-none focus:border-indigo-500"
          />
          <div className="flex items-center gap-2">
            <button
              onClick={doReplay}
              disabled={sending}
              className="rounded-md bg-emerald-600 px-4 py-1.5 text-sm font-medium text-white hover:bg-emerald-500 disabled:opacity-50"
            >
              {sending ? "Sending…" : "Send"}
            </button>
            <button
              onClick={() => setReplayOpen(false)}
              className="rounded-md border border-zinc-700 px-3 py-1.5 text-sm text-zinc-300 hover:bg-zinc-800"
            >
              Cancel
            </button>
            {error && <span className="text-sm text-rose-400">{error}</span>}
          </div>
          {result && (
            <div className="mt-3">
              <span
                className={`inline-block rounded px-2 py-0.5 text-xs font-semibold ring-1 ring-inset ${statusStyle(result.status ?? 0)}`}
              >
                {result.status}
              </span>
              {result.body && <Section label="Response Body">{result.body}</Section>}
            </div>
          )}
        </div>
      )}
    </div>
  );
}
