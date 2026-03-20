# ADR-033 — Accord: Shared API Types and Error Codes

**Status:** Accepted
**Date:** 2026-03-20
**Author:** Harsh Maury
**Scope:** New project — github.com/Harshmaury/Accord
**Depends on:** ADR-003 (HTTP/JSON protocol)

---

## Context

engxd and engx share no types. Every CLI command defines anonymous inline structs
to parse API responses. When a field is added or renamed in engxd, the CLI silently
receives nil/zero — no compile-time error. Error responses are free-form strings,
making `engx ci` and `guard` fragile under error conditions.

---

## Decision

Create `github.com/Harshmaury/Accord` — a zero-dependency library defining:

1. All Nexus HTTP API request/response types
2. Generic `Response[T]` envelope matching `{"ok":bool,"data":T}`
3. Structured `ErrorCode` enum — stable machine-readable identifiers
4. `api.Version` constant and `X-Nexus-API-Version` header name

### Rules
1. Zero non-stdlib dependencies
2. No business logic — types and constants only
3. Semver strictly — minor = additive, major = breaking
4. engxd migrates handlers to accord types in a subsequent phase

---

## Compliance

| ADR | Status |
|-----|--------|
| ADR-003 | ✅ Formalises the existing HTTP/JSON protocol |

---

## Next ADR
ADR-034 — Herald: typed HTTP client for the Nexus API
