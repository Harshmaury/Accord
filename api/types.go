// Package api defines the compile-time contract between engxd (server)
// and herald (client). Both import this package — field changes that break
// compatibility fail at compile time, not silently at runtime.
//
// HARDENING RULES — never violate:
//   1. Field names and json tags must exactly match engxd handler output.
//      Verify against internal/api/handler/*.go before any change.
//   2. Never add required fields to request types — always omitempty.
//      Breaking change = major version bump (coordinate nexus + herald).
//   3. Response[T] is the only envelope — never parse raw {"ok":bool} outside this package.
//   4. ErrorCode is stable — values are API surface. Never rename, only add.
//   5. Version must match internal/api/server.go X-Nexus-API-Version header value.
//
// Verified against: nexus internal/api/handler/*.go, internal/state/db.go
// Last verified: 2026-03-20
package api

// ── API VERSION ───────────────────────────────────────────────────────────────

// Version is the current Nexus HTTP API version.
// MUST match the value returned by engxd in X-Nexus-API-Version header.
// Increment when any response field is renamed or removed (major change).
const Version = "1"

// VersionHeader is the HTTP response header name carrying the API version.
const VersionHeader = "X-Nexus-API-Version"

// ── ENVELOPE ─────────────────────────────────────────────────────────────────

// Response is the standard JSON envelope for all Nexus API responses.
// Shape verified against internal/api/handler/respond.go:
//   {"ok": bool, "data": T, "error": "string"}
type Response[T any] struct {
	OK    bool   `json:"ok"`
	Data  T      `json:"data,omitempty"`
	Error string `json:"error,omitempty"`
}

// ── ERROR CODES ───────────────────────────────────────────────────────────────

// ErrorCode is a stable machine-readable error identifier.
// Switch on Code — NEVER on Message (messages are human-readable and may change).
// All codes are permanent API surface — never rename or remove.
type ErrorCode string

const (
	ErrNotFound          ErrorCode = "NOT_FOUND"          // 404 — resource does not exist
	ErrAlreadyExists     ErrorCode = "ALREADY_EXISTS"      // 409 — resource already registered
	ErrInvalidInput      ErrorCode = "INVALID_INPUT"       // 400 — malformed request body or params
	ErrUnauthorized      ErrorCode = "UNAUTHORIZED"        // 401 — missing or invalid X-Service-Token
	ErrDaemonUnavailable ErrorCode = "DAEMON_UNAVAILABLE"  // transport — cannot reach engxd
	ErrVersionMismatch   ErrorCode = "VERSION_MISMATCH"    // X-Nexus-API-Version mismatch
	ErrInternal          ErrorCode = "INTERNAL"            // 500 — unexpected server error
)

// Error is a structured API error. Implements the error interface.
type Error struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
}

func (e *Error) Error() string {
	if e == nil {
		return "<nil>"
	}
	return string(e.Code) + ": " + e.Message
}

// ── SERVICE TYPES ─────────────────────────────────────────────────────────────
// Verified against: internal/state/db.go Service struct + GET /services handler output.

// ServiceDTO is the API representation of a managed service.
// Field names match state.Service json tags exactly.
type ServiceDTO struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Project      string `json:"project"`
	DesiredState string `json:"desired_state"`  // matches state.Service.DesiredState
	ActualState  string `json:"actual_state"`   // matches state.Service.ActualState
	Provider     string `json:"provider"`
	Config       string `json:"config"`
	FailCount    int    `json:"fail_count"`
}

// ServiceRegisterRequest is the body for POST /services/register.
// Verified against internal/api/handler/services_register.go registerServiceRequest.
type ServiceRegisterRequest struct {
	ID       string `json:"id"`
	Name     string `json:"name,omitempty"`   // optional — defaults to ID if empty
	Project  string `json:"project"`
	Provider string `json:"provider"`         // "process" | "docker" | "k8s"
	Config   string `json:"config"`           // JSON-encoded provider config
}

// ServiceRegisterResponse is the body for POST /services/register.
// Verified against handler: respondOK(w, map[string]string{"id": req.ID, "name": req.Name})
type ServiceRegisterResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// ServiceResetResponse is the body for POST /services/{id}/reset.
// Verified against handler: respondOK(w, map[string]any{"id": id, "reset": true})
type ServiceResetResponse struct {
	ID    string `json:"id"`
	Reset bool   `json:"reset"`
}

// ── PROJECT TYPES ─────────────────────────────────────────────────────────────
// Verified against: internal/state/db.go Project struct
//                   internal/controllers/project_controller.go ProjectStatus

// ProjectDTO is the API representation of a registered project.
// Note: GET /projects returns []controllers.ProjectStatus, not []state.Project.
// This DTO matches the ProjectStatus shape from controllers.
type ProjectDTO struct {
	ID       string       `json:"id"`
	Name     string       `json:"name"`
	Services []ServiceDTO `json:"services"`
}

// ProjectActionResponse is the body for POST /projects/{id}/start and /stop.
// Verified against internal/api/handler/projects.go Start/Stop handlers.
type ProjectActionResponse struct {
	ProjectID string `json:"project_id"`
	Queued    int    `json:"queued"`
	Message   string `json:"message,omitempty"`
}

// ProjectRegisterRequest is the body for POST /projects/register.
// Verified against internal/api/handler/projects.go Register handler.
type ProjectRegisterRequest struct {
	ID          string `json:"id"`
	Name        string `json:"name,omitempty"`
	Path        string `json:"path"`
	Language    string `json:"language,omitempty"`
	ProjectType string `json:"project_type,omitempty"`
}

// ── AGENT TYPES ───────────────────────────────────────────────────────────────
// Verified against: internal/state/db_agents.go Agent struct
//                   internal/api/handler/agents.go

// AgentDTO is the API representation of a registered remote agent.
// Token field is intentionally omitted — never returned by the API.
type AgentDTO struct {
	ID           string `json:"id"`
	Hostname     string `json:"hostname"`
	Address      string `json:"address"`
	Online       bool   `json:"online"`       // derived: last_seen within 30s
	LastSeen     string `json:"last_seen"`    // RFC3339 string from handler
	RegisteredAt string `json:"registered_at"`
}

// ── EVENT TYPES ───────────────────────────────────────────────────────────────
// Verified against: internal/state/db.go Event struct
//                   internal/api/handler/events.go

// EventDTO is the API representation of a platform event.
type EventDTO struct {
	ID        int64  `json:"id"`
	ServiceID string `json:"service_id"`
	Type      string `json:"type"`
	Source    string `json:"source"`
	TraceID   string `json:"trace_id"`
	Component string `json:"component"`
	Outcome   string `json:"outcome"`
	Payload   string `json:"payload,omitempty"`
	CreatedAt string `json:"created_at"` // RFC3339
}

// ── HEALTH TYPES ──────────────────────────────────────────────────────────────
// Verified against: internal/api/server.go makeHealthHandler

// HealthData is the response body for GET /health.
type HealthData struct {
	Status        string  `json:"status"`
	UptimeSeconds float64 `json:"uptime_seconds"`
	DaemonVersion string  `json:"daemon_version"`
}

// ── INVARIANT CHECKS ─────────────────────────────────────────────────────────
// These constants document the API surface explicitly.
// If engxd handler output changes, these must be updated here AND
// a version bump must be coordinated across nexus + herald + accord.

const (
	// ProviderProcess is the process runtime provider identifier.
	ProviderProcess = "process"
	// ProviderDocker is the Docker runtime provider identifier.
	ProviderDocker = "docker"
	// ProviderK8s is the Kubernetes runtime provider identifier.
	ProviderK8s = "k8s"

	// StateRunning is the running service state string.
	StateRunning = "running"
	// StateStopped is the stopped service state string.
	StateStopped = "stopped"
	// StateMaintenance is the maintenance service state string.
	StateMaintenance = "maintenance"
	// StateCrashed is the crashed service state string.
	StateCrashed = "crashed"
)

// ServiceResetData carries the result of a service reset.
type ServiceResetData struct {
	ID    string `json:"id"`
	Reset bool   `json:"reset"`
}

// ServiceIdentity is a minimal service identifier returned by create operations.
type ServiceIdentity struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// ProjectActionData carries the result of a project lifecycle action.
type ProjectActionData struct {
	ProjectID string `json:"project_id"`
	Queued    int    `json:"queued"`
	Message   string `json:"message,omitempty"`
}
