// @accord-project: accord
// @accord-path: SERVICE-CONTRACT.md
# SERVICE-CONTRACT.md — Accord
# @version: 0.1.2
# @updated: 2026-03-25

**Type:** Library · **Module:** `github.com/Harshmaury/Accord` · **Domain:** Shared protocol

---

## Code

```
api/types.go       Response[T], ErrorCode, ServiceDTO, ProjectDTO, AgentDTO, EventDTO, HealthData
api/upstream.go    AtlasProjectDTO, ForgeExecutionDTO, GuardianReportDTO, NavigatorGraphDTO,
                   NexusMetricsDTO, GateValidateResponse, PlanSpanDTO, IdentityClaimDTO
api/payload.go     DecodePayload[T any], EventPayloadType registry
api/errors.go      Error struct — implements error interface
```

---

## Contract

No HTTP surface. Compile-time only.

**Envelope:** all platform HTTP responses use `Response[T]{OK bool, Data T, Error string}`.

**ErrorCode values (permanent — never rename):**
`NOT_FOUND` · `ALREADY_EXISTS` · `INVALID_INPUT` · `UNAUTHORIZED` · `FORBIDDEN` · `DAEMON_UNAVAILABLE` · `VERSION_MISMATCH` · `INTERNAL`

**Payload decoding:**
```go
alert, err := accord.DecodePayload[canonevents.SystemAlertPayload](event.Payload)
```

**Verification requirement:** before changing any type, verify json tags match the actual handler output in the service it mirrors. See verification table in `api/types.go` header comment.

**Versioning:**

| Change | Bump | Action |
|--------|------|--------|
| Add optional field | patch | Accord only |
| Add new type | minor | Accord only |
| Rename field / json tag | **major** | Nexus + Herald + Accord simultaneously |
| Remove field | **major** | Nexus + Herald + Accord simultaneously |
| Rename ErrorCode | **major** | All callers switch on Code |

**Import rules:**

| Consumer | Rule |
|----------|------|
| engxd handlers | MAY import |
| Herald | MUST import — all types from here |
| engx CLI | MUST NOT — use Herald |
| Atlas, Forge | MUST NOT — use Herald |

---

## Control

No runtime behavior. Zero dependencies beyond Go stdlib.

---

## Context

Compile-time contract between Nexus (encoder) and Herald (decoder). Disagreement fails the build, not the runtime.
