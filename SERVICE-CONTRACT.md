# SERVICE-CONTRACT.md — Accord

**Type:** Library (no runtime, no HTTP server)
**Module:** github.com/Harshmaury/Accord
**Version:** 0.1.0

---

## What Accord Is

Accord is the compile-time contract between engxd (server) and herald (client).
It defines:

- All Nexus HTTP API request/response types
- Structured error codes (ErrorCode enum)
- API version constant and header name
- Generic Response[T] envelope

Accord has zero non-stdlib dependencies. Any Go program can import it.

---

## What Accord Is NOT

- Not a client (that is herald)
- Not a server (that is nexus/internal/api)
- Not a runtime (that is anchor)
- No business logic — types and constants only

---

## Import Rules

| Consumer | May import accord? |
|---|---|
| engxd (server) | ✅ Yes — use types in handlers |
| herald (client) | ✅ Yes — decode responses with types |
| engx (CLI) | ✅ Yes — via herald only |
| atlas, forge | ✅ Yes — via herald only |
| Any Go program | ✅ Yes — zero deps |

---

## Versioning

Accord follows semver strictly.
- Patch: fix typos, add omitempty tags
- Minor: add new types or fields (backwards compatible)
- Major: rename/remove fields (breaking — coordinate with nexus + herald)

When accord major version bumps, nexus and herald must bump together.
