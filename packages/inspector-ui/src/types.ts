// Shape of a captured request as returned by the Hatch API.
// `headers` is a JSON-encoded string; `body` is base64 (Go marshals []byte as base64).
export interface HatchRequest {
  id: string;
  endpoint_id: string;
  method: string;
  path: string;
  headers: string;
  query: string;
  body: string | null;
  created_at: string;
}

export interface ReplayResult {
  status?: number;
  headers?: Record<string, string>;
  body?: string;
  error?: string;
}
