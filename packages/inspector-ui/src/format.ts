export function relativeTime(iso: string): string {
  const d = new Date(iso);
  if (isNaN(d.getTime())) return iso;
  const diff = Math.floor((Date.now() - d.getTime()) / 1000);
  if (diff < 5) return "just now";
  if (diff < 60) return `${diff}s ago`;
  if (diff < 3600) return `${Math.floor(diff / 60)}m ago`;
  if (diff < 86400) return `${Math.floor(diff / 3600)}h ago`;
  if (diff < 604800) return `${Math.floor(diff / 86400)}d ago`;
  return d.toLocaleDateString("en-US", { month: "short", day: "numeric" });
}

const METHOD_STYLES: Record<string, string> = {
  GET: "bg-emerald-500/15 text-emerald-400 ring-emerald-500/30",
  POST: "bg-sky-500/15 text-sky-400 ring-sky-500/30",
  PUT: "bg-amber-500/15 text-amber-400 ring-amber-500/30",
  PATCH: "bg-amber-500/15 text-amber-400 ring-amber-500/30",
  DELETE: "bg-rose-500/15 text-rose-400 ring-rose-500/30",
};

export function methodStyle(method: string): string {
  return METHOD_STYLES[method.toUpperCase()] ?? "bg-zinc-500/15 text-zinc-400 ring-zinc-500/30";
}

export function statusStyle(status: number): string {
  if (status >= 500) return "bg-rose-500/15 text-rose-400 ring-rose-500/30";
  if (status >= 400) return "bg-rose-500/15 text-rose-400 ring-rose-500/30";
  if (status >= 300) return "bg-amber-500/15 text-amber-400 ring-amber-500/30";
  return "bg-emerald-500/15 text-emerald-400 ring-emerald-500/30";
}

// Decode a base64 body to text for display. Falls back to the raw value if not decodable.
export function decodeBody(body: string | null): string {
  if (!body) return "";
  try {
    return decodeURIComponent(escape(atob(body)));
  } catch {
    try {
      return atob(body);
    } catch {
      return body;
    }
  }
}

export function prettyJSON(raw: string): string {
  try {
    return JSON.stringify(JSON.parse(raw), null, 2);
  } catch {
    return raw;
  }
}
