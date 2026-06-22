# Technology Stack

## Overview

These choices were made by the CTO on 2026-06-22, aligned with the company charter and operating principles. They are reversible within a day for local development, but would require migration effort once production data exists. All choices default to boring, well-supported technology over novelty.

## Frontend

| Layer | Choice | Rationale |
|---|---|---|
| Framework | Next.js (App Router) | Full-stack React with SSR/SSG, API routes, and a large ecosystem. App Router is the stable future path. |
| Language | TypeScript (strict) | Catches entire classes of bugs at build time. Strict mode is non-negotiable. |
| Styling | Tailwind CSS | Utility-first, colocated with components, no separate CSS files to maintain. |
| UI Components | shadcn/ui pattern | Copy-paste components we own and can customize. No opaque UI library dependencies. |

## Backend & Data

| Layer | Choice | Rationale |
|---|---|---|
| Database | PostgreSQL | Proven, feature-rich, great ORM support. Our relational data model fits it well. |
| ORM | Prisma | Type-safe queries, excellent migration tooling, good DX. |
| Auth | NextAuth.js | Battle-tested, supports many providers, integrates cleanly with Next.js. |
| Storage | Cloudflare R2 | S3-compatible, zero egress fees, good for file uploads and static assets. |
| Payments | TBD | Will evaluate when we have a revenue model. |

## Tooling

| Layer | Choice | Rationale |
|---|---|---|
| Package Manager | pnpm | Fast, disk-space efficient, strict `node_modules` layout avoids phantom dependencies. |
| Monorepo | Turborepo (when needed) | Caching and task orchestration. Only adopt when we have >1 app or shared package. |
| CI/CD | GitHub Actions | Native GitHub integration, free for public repos, cheap for private. |
| Lint | ESLint + Prettier | Standard, autofixable, low-friction. |

## Principles

- **Server Components by default.** Reach for `"use client"` only when the component actually needs state, effects, or browser APIs.
- **Pure logic in `lib/`, I/O in `services/`.** Business logic does not import from `next/server` or call `fetch` directly.
- **Store money as integers.** MNT (Mongolian Tugrik) in smallest unit. Format on display.
- **Observability before optimization.** Measure before fixing. No tuning without metrics.
- **Idempotency.** Operations should be safe to retry. Infrastructure changes should be reproducible.

## Open Decisions

| Decision | Status | Owner | Blocker |
|---|---|---|---|
| Hosting provider (Vercel, Fly, AWS?) | Open | CTO | Need product requirements and traffic estimates |
| Monitoring / alerting stack | Open | CTO | Need hosting decision |
| CDN / edge strategy | Open | CTO | Need hosting decision |
