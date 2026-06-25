# Deployment handoff — Hatch Core

> Read this before touching deploy. Written 2026-06-25 on branch `rearchitect/core-split`.
> It explains what the repo split changed and **exactly what is broken** as a result.

## TL;DR

The repo was split into an open-core layout. Everything that builds or deploys still
points at the **old paths** and is therefore broken. There is also a **new build
dependency**: the Go binary embeds the compiled web UI, so the frontend must be built
*before* the Go build. The old pipeline never did this.

Nothing here is committed yet. This is the state of the branch.

## What the layout is now

```
core/
  api/                      # Go service (the old root project, moved here)
    go.mod  go.sum          # module path UNCHANGED: github.com/elfoundation/hatch
    cmd/hatch/              # main entrypoint  (was ./cmd/hatch, now ./core/api/cmd/hatch)
    internal/
      handler/ store/ ...
      web/
        web.go              # //go:embed all:dist  -> bundles the SPA into the binary
        dist/              # build OUTPUT of core/web. Vite writes here.
  web/                      # Vite + React 19 + TS + Tailwind v4 SPA  (@hatch/core-web)
    vite.config.ts         # build.outDir = ../api/internal/web/dist
packages/
  inspector-ui/             # shared React component lib (@hatch/inspector-ui), source-only
package.json                # root, pnpm workspace scripts
pnpm-workspace.yaml         # packages: core/web, packages/*
```

site/ (the marketing/landing app) is untouched and out of scope for Core deploy.

## The one thing that is easy to get wrong

**The Go binary does NOT build the frontend. It embeds an already-built copy of it.**

`core/api/internal/web/web.go` has `//go:embed all:dist`. If `core/api/internal/web/dist`
is empty or missing, **`go build` fails** with an embed error. So the correct build order
is always:

1. `pnpm install`
2. `pnpm build:web`   → Vite compiles and writes into `core/api/internal/web/dist`
3. `CGO_ENABLED=0 go build ./core/api/cmd/hatch`  (run from repo root)

Open decision (pick one and make the pipeline match):
- **(A) Build dist in CI/Docker** — add a Node/pnpm stage before the Go stage. Keeps the
  repo clean. `dist` stays a build artifact. *Recommended.*
- **(B) Commit `core/api/internal/web/dist`** — then `go build` works standalone with no
  Node toolchain. Note: `dist` is currently **not** gitignored (root `/dist` rule is
  anchored to root, so it does not match the nested path), so committing it would "just
  work" but you'd be checking in build output.

## What is broken right now (and the fix)

### 1. `Dockerfile` — broken, two reasons
Current contents assume the old root layout:
```dockerfile
COPY go.mod go.sum ./        # these are now at core/api/
RUN ... go build ... ./cmd/hatch   # now ./core/api/cmd/hatch
```
It also has **no frontend stage**, so even with paths fixed it would fail the `go:embed`
(empty dist) — see above.

Fix: make it multi-stage:
- stage 1 `node:22-alpine`: `pnpm install` + `pnpm build:web` (needs the whole repo:
  root package.json, pnpm-workspace.yaml, core/web, packages/inspector-ui).
- stage 2 `golang:1.25-alpine`: copy source incl. the freshly built
  `core/api/internal/web/dist`, then `cd core/api && CGO_ENABLED=0 go build -o /bin/hatch ./cmd/hatch`
  (or build from root with `./core/api/cmd/hatch`).
- stage 3 `scratch`: copy `/bin/hatch`. EXPOSE 8080. (unchanged)

### 2. `.github/workflows/ci.yml` — broken
Every Go step runs from repo root but the module now lives in `core/api/`:
- `gofmt -l .`, `go vet ./...`, `go test ./... -race` → run from `core/api/` (or add a
  `working-directory: core/api`).
- `go build ... ./cmd/hatch` (line ~35) → `./core/api/cmd/hatch`, **and** must run
  `pnpm build:web` first or it fails on embed.
- `docker` job builds context `.` → fine once the Dockerfile is fixed.
- Add Node + pnpm setup to whatever job builds the binary.

### 3. `.github/workflows/release.yml` — broken
`go build ... ./cmd/hatch` (~line 82) → `./core/api/cmd/hatch`, same pnpm-build-first
requirement. Cross-compiled release binaries each need dist present at build time (dist
is platform-independent, so one `pnpm build:web` up front covers all targets).

### 4. `docker-compose.yml` / `Caddyfile` — paths fine, no change needed
compose just does `build: .` and sets env (`HATCH_PORT`, `HATCH_DB_PATH=/data/hatch.db`,
`HATCH_BASE_URL`). Caddy reverse-proxies `hatch:8080`. Both work once the Dockerfile
builds. The `caddy` service is behind the `with-caddy` profile (opt-in for TLS).

## Runtime contract (unchanged by the split)

- Single static binary, listens on `HATCH_PORT` (default 8080), serves API + SPA + SSE.
- `HATCH_DB_PATH` → SQLite file (WAL, single writer). Mount a volume at `/data`.
- `HATCH_BASE_URL` → used when displaying capture URLs.
- `HATCH_ALLOW_PRIVATE_REPLAY=true` → lets replay hit loopback/private IPs. **Off by
  default on purpose** (SSRF guard). Do not enable in a shared/hosted deploy.
- Health: `GET /healthz` (liveness), `GET /readyz` (DB ping). Compose smoke test curls
  `/healthz`.
- No auth, no bin-ID unguessability, no retention sweep, no body-size cap yet. Fine for
  self-host on a trusted network; **not** safe as an open public service as-is.

## Sanity check before you trust a build

```bash
pnpm install && pnpm build:web
test -f core/api/internal/web/dist/index.html   # must exist
(cd core/api && CGO_ENABLED=0 go build -o /tmp/hatch ./cmd/hatch)
(cd core/api && go test ./...)                   # all 3 packages should pass
/tmp/hatch & sleep 1; curl -fsS localhost:8080/healthz; curl -fsSI localhost:8080/ | head -1
```
The binary should serve the dark-theme SPA at `/` and `/e/{id}`, and stream live captures
over SSE at `/e/{id}/events`. If `/` returns "UI not built", the embed step was skipped.

## Out of scope here

Hatch Cloud (Next.js marketing + ngrok-like subdomain tunneling) is a separate system and
is not built yet. This doc covers Core only.
