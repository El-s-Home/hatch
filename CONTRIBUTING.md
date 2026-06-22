# Contributing to El Foundation

## One Task = One Branch = One Owner

Every issue gets its own branch. Branch names follow this pattern:

```
owner-identifier/short-description
```

Examples:
- `cto/ELF-4-engineering-foundation`
- `eng-1/add-auth-middleware`

## No Direct Commits to `main`

All changes go through a pull request. No exceptions.

## Pull Request Process

1. **Open a PR** from your branch to `main`.
2. **Fill out the PR template** (risk, rollback, verification).
3. **Request review** from the relevant owner:
   - Code changes → another engineer or CTO
   - UX-facing changes → UXDesigner
   - Security-sensitive changes → SecurityEngineer
4. **Address feedback** or escalate disagreements in writing.
5. **Ship on green.** Once CI passes and review is approved, the owner merges.

## Commit Messages

Commit messages explain **why**, not what. The diff shows what changed; the message explains the reasoning.

Good:
```
Add rotating refresh tokens

Using a rotating refresh token strategy prevents replay attacks
and gives us a clean theft-detection signal. See ADR-003.

Co-Authored-By: Paperclip <noreply@paperclip.ing>
```

Bad:
```
Update auth.ts
```

## Code Style

- TypeScript strict. No `any` unless explicitly approved — use `unknown` + narrowing.
- Keep `lib/` (pure logic) and `services/` (I/O, DB, network) separate.
- Prefer small modules over clever abstractions.
- No comments unless the code is genuinely non-obvious or there is a real `// FIXME`.
- No defensive try/catch around things that should not fail. Let it throw.
- Server Components by default; reach for `"use client"` only when state, effects, or browser APIs are needed.

## Definition of Done

A task is not done until **all** of the following are true:

1. Code is written and reviewed.
2. Tests pass. CI is green.
3. Documentation is updated.
4. No secrets in plain text.
5. User-facing changes are validated.
6. Rollback path is known.
7. Handoff is clean — follow-up work is captured in a new issue.

## Security

- Never commit secrets, credentials, or customer data.
- Security-sensitive changes (auth, crypto, secrets, permissions) require SecurityEngineer review before merging.
- Report vulnerabilities to the CTO immediately.

## Questions?

Open an issue or ask in the project channel. Async-first: write it down.
