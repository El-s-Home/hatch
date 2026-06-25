import { useEffect, useRef, useState } from "react";
import type { HatchRequest } from "./types";

interface Options {
  // Base URL of the Hatch API. Empty string = same origin.
  apiBase?: string;
  // Max rows to keep in the DOM (prevents unbounded growth in long-lived tabs).
  cap?: number;
}

// Subscribes to a Hatch endpoint: loads the initial page, then live-prepends
// new requests via SSE. One data source for both initial and live rows.
export function useRequestStream(endpointId: string, opts: Options = {}) {
  const { apiBase = "", cap = 500 } = opts;
  const [requests, setRequests] = useState<HatchRequest[]>([]);
  const [connected, setConnected] = useState(false);
  const seen = useRef<Set<string>>(new Set());

  useEffect(() => {
    if (!endpointId) return;
    let cancelled = false;
    seen.current = new Set();

    fetch(`${apiBase}/v1/endpoints/${encodeURIComponent(endpointId)}/requests?limit=100`)
      .then((r) => r.json())
      .then((rows: HatchRequest[] | null) => {
        if (cancelled || !rows) return;
        rows.forEach((r) => seen.current.add(r.id));
        setRequests(rows);
      })
      .catch(() => {});

    const es = new EventSource(`${apiBase}/e/${encodeURIComponent(endpointId)}/events`);
    es.onopen = () => !cancelled && setConnected(true);
    es.onerror = () => !cancelled && setConnected(false);
    es.onmessage = (e) => {
      if (cancelled) return;
      try {
        const req: HatchRequest = JSON.parse(e.data);
        if (seen.current.has(req.id)) return;
        seen.current.add(req.id);
        setRequests((prev) => [req, ...prev].slice(0, cap));
      } catch {
        /* ignore malformed frames */
      }
    };

    return () => {
      cancelled = true;
      es.close();
    };
  }, [endpointId, apiBase, cap]);

  return { requests, connected };
}
