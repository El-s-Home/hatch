import { useRequestStream } from "./useRequestStream";
import { RequestRow } from "./RequestRow";

interface Props {
  endpointId: string;
  apiBase?: string;
}

export function RequestList({ endpointId, apiBase = "" }: Props) {
  const { requests, connected } = useRequestStream(endpointId, { apiBase });

  return (
    <div>
      <div className="mb-6 flex items-baseline justify-between">
        <div>
          <h1 className="text-lg font-semibold text-zinc-100">Hatch</h1>
          <div className="font-mono text-sm text-zinc-500">
            Endpoint: <span className="text-zinc-300">{endpointId}</span>
          </div>
        </div>
        <span className="flex items-center gap-1.5 text-xs text-zinc-500">
          <span
            className={`inline-block h-2 w-2 rounded-full ${connected ? "bg-emerald-400" : "bg-zinc-600"}`}
          />
          {connected ? "live" : "offline"}
        </span>
      </div>

      {requests.length === 0 ? (
        <div className="py-20 text-center text-zinc-500">
          <p className="mb-1 text-zinc-300">Waiting for requests…</p>
          <p className="text-sm">
            Send a request to{" "}
            <code className="rounded bg-zinc-800 px-1.5 py-0.5 text-zinc-300">/{endpointId}</code> to
            see it appear here.
          </p>
        </div>
      ) : (
        <div className="space-y-3">
          {requests.map((req) => (
            <RequestRow key={req.id} req={req} apiBase={apiBase} />
          ))}
        </div>
      )}
    </div>
  );
}
