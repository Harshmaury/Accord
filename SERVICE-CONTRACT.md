# SERVICE-CONTRACT.md — Accord
# @version: 1.0.0
# @updated: 2026-03-20

**Module:** `github.com/Harshmaury/Accord`
**Type:** Library — zero runtime, zero dependencies
**Role:** Compile-time contract between engxd (server) and herald (client)

---

## Contract Definition

Accord owns the answer to: *"what does the Nexus HTTP API look like?"*

Every field name, every JSON tag, every error code is defined here once.
engxd encodes with these types. Herald decodes with these types.
If they disagree — the build fails. Not the runtime. The build.

---

## Verification Requirement (MANDATORY)

Before any change to accord types, verify against the actual engxd handlers:

| Accord type | Verify against |
|---|---|
| `ServiceDTO` | `internal/state/db.go` `Service` struct json tags |
| `ServiceRegisterRequest` | `internal/api/handler/services_register.go` `registerServiceRequest` |
| `ServiceResetResponse` | `internal/api/handler/services_reset.go` `respondOK(w, map[string]any{...})` |
| `ProjectDTO` | `internal/controllers/project_controller.go` `ProjectStatus` |
| `ProjectActionResponse` | `internal/api/handler/projects.go` `Start`/`Stop` handlers |
| `AgentDTO` | `internal/state/db_agents.go` `Agent` struct + handler serialisation |
| `EventDTO` | `internal/state/db.go` `Event` struct json tags |
| `HealthData` | `internal/api/server.go` `makeHealthHandler` response shape |

**If a verification step is skipped, the change is rejected.**

---

## Versioning Rules

| Change | Version bump | Action required |
|---|---|---|
| Add new optional field (omitempty) | patch | update accord only |
| Add new type | minor | update accord only |
| Rename field / change json tag | **major** | coordinate nexus + herald + accord simultaneously |
| Remove field | **major** | coordinate nexus + herald + accord simultaneously |
| Add new ErrorCode | minor | update accord only |
| Rename ErrorCode | **major** | breaking — callers switch on Code |

A major version bump requires:
1. ADR documenting the breaking change
2. nexus, herald, accord all updated in a single coordinated release
3. All consumers (engx, atlas, forge) updated before the old version is removed

---

## Import Rules

| Consumer | Rule |
|---|---|
| engxd handlers | MAY import accord to use request/response types |
| herald | MUST import accord — all types come from here |
| engx CLI | MUST NOT import accord directly — use herald |
| Atlas, Forge | MUST NOT import accord directly — use herald |
| External tools | MAY import accord directly — zero deps, stable API |

---

## What Accord Is NOT

- Not a client (that is herald)
- Not a server (that is nexus)
- No HTTP calls
- No goroutines
- No state

---

## Breakage Prevention Checklist

Before merging any PR to accord:

- [ ] All changed types verified against current nexus handler output
- [ ] No non-stdlib imports added
- [ ] No business logic added
- [ ] Version constant updated if breaking change
- [ ] herald go.mod updated to new accord version
- [ ] nexus go.mod updated if handlers migrate to accord types
