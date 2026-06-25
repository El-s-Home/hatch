import { RequestList } from "@hatch/inspector-ui";

// Endpoint id comes from the path: /e/{id}. Fallback to ?e= for dev.
function endpointFromLocation(): string {
  const m = window.location.pathname.match(/^\/e\/([^/]+)/);
  if (m) return decodeURIComponent(m[1]);
  return new URLSearchParams(window.location.search).get("e") ?? "";
}

export function App() {
  const endpointId = endpointFromLocation();

  return (
    <div className="min-h-screen bg-zinc-950 text-zinc-100">
      <div className="mx-auto max-w-3xl px-4 py-8">
        {endpointId ? (
          <RequestList endpointId={endpointId} />
        ) : (
          <div className="py-20 text-center text-zinc-400">
            <h1 className="mb-2 text-lg font-semibold text-zinc-100">Hatch</h1>
            <p className="text-sm">
              Open an endpoint at{" "}
              <code className="rounded bg-zinc-800 px-1.5 py-0.5 text-zinc-300">/e/&lt;name&gt;</code>
            </p>
          </div>
        )}
      </div>
    </div>
  );
}
