// Package api defines the shared request/response types and error codes
// for the Nexus HTTP API. Both engxd (server) and herald (client) import
// this package — it is the compile-time contract between them.
//
// Rule: no business logic lives here. Types and constants only.
// Rule: no imports from Nexus, Herald, or any platform service.
// Rule: zero non-stdlib dependencies — accord must be importable by anyone.
package api

// ── API VERSION ───────────────────────────────────────────────────────────────

// Version is the current Nexus HTTP API version.
// Sent by engxd as X-Nexus-API-Version response header.
// Checked by herald on every response.
const Version = "1"

// VersionHeader is the HTTP header name carrying the API version.
const VersionHeader = "X-Nexus-API-Version"

// ── ENVELOPE ─────────────────────────────────────────────────────────────────

// Response is the standard JSON envelope for all Nexus API responses.
type Response[T any] struct {
	OK    bool   `json:"ok"`
	Data  T      `json:"data,omitempty"`
	Error *Error `json:"error,omitempty"`
}

// Error is the structured error payload returned when ok=false.
type Error struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
}

func (e *Error) Error() string {
	if e == nil {
		return ""
	}
	return string(e.Code) + ": " + e.Message
}

// ── ERROR CODES ───────────────────────────────────────────────────────────────

// ErrorCode is a stable machine-readable identifier for API errors.
// Callers switch on Code — never on Message (messages may change).
type ErrorCode string

const (
	// ErrNotFound is returned when a resource does not exist.
	ErrNotFound ErrorCode = "NOT_FOUND"

	// ErrAlreadyExists is returned when creating a resource that already exists.
	ErrAlreadyExists ErrorCode = "ALREADY_EXISTS"

	// ErrInvalidInput is returned when the request body or params are malformed.
	ErrInvalidInput ErrorCode = "INVALID_INPUT"

	// ErrUnauthorized is returned when X-Service-Token is missing or invalid.
	ErrUnauthorized ErrorCode = "UNAUTHORIZED"

	// ErrDaemonUnavailable is returned when engxd cannot be reached.
	ErrDaemonUnavailable ErrorCode = "DAEMON_UNAVAILABLE"

	// ErrVersionMismatch is returned when client and server API versions differ.
	ErrVersionMismatch ErrorCode = "VERSION_MISMATCH"

	// ErrInternal is returned for unexpected server-side errors.
	ErrInternal ErrorCode = "INTERNAL"
)

// ── SERVICE TYPES ─────────────────────────────────────────────────────────────

// ServiceDTO is the API representation of a managed service.
type ServiceDTO struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Project string `json:"project"`
	Desired string `json:"desired"`
	Actual  string `json:"actual"`
	Config  string `json:"config"`
}

// ServiceListResponse is the response body for GET /services.
type ServiceListResponse = Response[[]ServiceDTO]

// ServiceRegisterRequest is the request body for POST /services/register.
type ServiceRegisterRequest struct {
	ID       string `json:"id"`
	Name     string `json:"name,omitempty"`
	Project  string `json:"project"`
	Provider string `json:"provider"`
	Config   string `json:"config"`
}

// ServiceRegisterResponse is the response body for POST /services/register.
type ServiceRegisterResponse = Response[ServiceIdentity]

// ServiceResetResponse is the response body for POST /services/{id}/reset.
type ServiceResetResponse = Response[ServiceResetData]

// ServiceResetData carries the result of a service reset.
type ServiceResetData struct {
	ID    string `json:"id"`
	Reset bool   `json:"reset"`
}

// ── PROJECT TYPES ─────────────────────────────────────────────────────────────

// ProjectDTO is the API representation of a registered project.
type ProjectDTO struct {
	ID       string       `json:"id"`
	Name     string       `json:"name"`
	Services []ServiceDTO `json:"services"`
	Status   string       `json:"status"`
}

// ProjectListResponse is the response body for GET /projects.
type ProjectListResponse = Response[[]ProjectDTO]

// ProjectStartResponse is the response body for POST /projects/{id}/start.
type ProjectStartResponse = Response[ProjectActionData]

// ProjectStopResponse is the response body for POST /projects/{id}/stop.
type ProjectStopResponse = Response[ProjectActionData]

// ProjectActionData carries the result of a project lifecycle action.
type ProjectActionData struct {
	ProjectID string `json:"project_id"`
	Queued    int    `json:"queued"`
	Message   string `json:"message,omitempty"`
}

// ── AGENT TYPES ───────────────────────────────────────────────────────────────

// AgentDTO is the API representation of a registered remote agent.
type AgentDTO struct {
	ID       string `json:"id"`
	Addr     string `json:"addr"`
	Online   bool   `json:"online"`
	LastSeen string `json:"last_seen"`
}

// AgentListResponse is the response body for GET /agents.
type AgentListResponse = Response[[]AgentDTO]

// ── EVENT TYPES ───────────────────────────────────────────────────────────────

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
	CreatedAt string `json:"created_at"`
}

// EventListResponse is the response body for GET /events.
type EventListResponse = Response[[]EventDTO]

// ── HEALTH TYPES ──────────────────────────────────────────────────────────────

// HealthResponse is the response body for GET /health.
type HealthResponse = Response[HealthData]

// HealthData carries daemon health information.
type HealthData struct {
	Status        string  `json:"status"`
	UptimeSeconds float64 `json:"uptime_seconds"`
	DaemonVersion string  `json:"daemon_version"`
	ServicesTotal int     `json:"services_total"`
}

// ── SHARED ────────────────────────────────────────────────────────────────────

// ServiceIdentity is a minimal service identifier returned by create operations.
type ServiceIdentity struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
