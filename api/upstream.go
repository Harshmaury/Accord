// @accord-project: Accord
// @accord-path: api/upstream.go
// Navigator types added: NavigatorNodeDTO, NavigatorEdgeDTO, NavigatorSummaryDTO,
// NavigatorGraphDTO — verified against navigator/internal/topology/model.go json tags.
// Upstream DTOs for Atlas, Forge, Guardian, and Nexus /metrics.
// These types are the compile-time contract between Herald clients and
// the services they call. Verified against each service's handler output.
//
// Verified against:
//   Atlas   internal/api/handler/workspace.go  projectResponse shape
//   Atlas   internal/api/handler/graph.go       GraphEdge shape
//   Forge   internal/store/storer.go            ExecutionRecord shape
//   Guardian internal/policy/model.go           Finding + Report shapes
//   Nexus   internal/telemetry/metrics.go       Snapshot shape
//
// Last verified: 2026-03-21
package api

import "time"

// ── ATLAS TYPES ───────────────────────────────────────────────────────────────

// AtlasProjectDTO is the API shape returned by Atlas GET /workspace/projects.
// Matches atlas/internal/api/handler/workspace.go projectResponse exactly.
type AtlasProjectDTO struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Path         string   `json:"path"`
	Language     string   `json:"language"`
	Type         string   `json:"type"`
	Source       string   `json:"source"`
	Status       string   `json:"status"` // "verified" | "unverified"
	Capabilities []string `json:"capabilities"`
	DependsOn    []string `json:"depends_on"`
	IndexedAt    string   `json:"indexed_at"`
}

// AtlasEdgeDTO is one graph edge from Atlas GET /workspace/graph.
// Matches the edge element shape in the handler response.
type AtlasEdgeDTO struct {
	FromID   string `json:"from_id"`
	ToID     string `json:"to_id"`
	EdgeType string `json:"edge_type"`
	Source   string `json:"source"`
}

// AtlasGraphDTO is the full response body of Atlas GET /workspace/graph.
type AtlasGraphDTO struct {
	Total  int            `json:"total"`
	Edges  []AtlasEdgeDTO `json:"edges"`
}

// ── FORGE TYPES ───────────────────────────────────────────────────────────────

// ForgeExecutionDTO is one record from Forge GET /history or GET /history/:trace_id.
// Matches forge/internal/store/storer.go ExecutionRecord json tags exactly.
type ForgeExecutionDTO struct {
	ID         string    `json:"id"`
	CommandID  string    `json:"command_id"`
	Intent     string    `json:"intent"`
	Target     string    `json:"target"`
	TraceID    string    `json:"trace_id"`
	Status     string    `json:"status"`
	Output     string    `json:"output,omitempty"`
	Error      string    `json:"error,omitempty"`
	DurationMS int64     `json:"duration_ms"`
	StartedAt  time.Time `json:"started_at"`
	FinishedAt time.Time `json:"finished_at"`
}

// ── GUARDIAN TYPES ────────────────────────────────────────────────────────────

// GuardianFindingDTO is one finding from Guardian GET /guardian/findings.
// Matches guardian/internal/policy/model.go Finding json tags exactly.
type GuardianFindingDTO struct {
	RuleID    string    `json:"rule_id"`
	Severity  string    `json:"severity"`
	Target    string    `json:"target"`
	Message   string    `json:"message"`
	Count     int       `json:"count"`
	FirstSeen time.Time `json:"first_seen"`
	LastSeen  time.Time `json:"last_seen"`
}

// GuardianSummaryDTO is the findings summary from GET /guardian/findings.
type GuardianSummaryDTO struct {
	Total    int `json:"total"`
	Warnings int `json:"warnings"`
	Errors   int `json:"errors"`
}

// GuardianReportDTO is the full response body of GET /guardian/findings.
// Matches guardian/internal/policy/model.go Report json tags exactly.
type GuardianReportDTO struct {
	Findings    []GuardianFindingDTO `json:"findings"`
	Summary     GuardianSummaryDTO   `json:"summary"`
	EvaluatedAt time.Time            `json:"evaluated_at"`
}

// ── NAVIGATOR TYPES ───────────────────────────────────────────────────────────

// NavigatorNodeDTO is one project node from Navigator GET /topology/graph.
// Matches navigator/internal/topology/model.go Node json tags exactly.
type NavigatorNodeDTO struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Type         string   `json:"type"`
	Language     string   `json:"language"`
	Status       string   `json:"status"` // "verified" | "unverified"
	Capabilities []string `json:"capabilities"`
	DependsOn    []string `json:"depends_on"`
	Source       string   `json:"source"`
	Path         string   `json:"path"`
}

// NavigatorEdgeDTO is one directed edge from Navigator GET /topology/graph.
// Matches navigator/internal/topology/model.go Edge json tags exactly.
// Note: Navigator uses json:"type" for EdgeType (not "edge_type" as in AtlasEdgeDTO).
type NavigatorEdgeDTO struct {
	From     string `json:"from"`
	To       string `json:"to"`
	EdgeType string `json:"type"`
	Source   string `json:"source"`
}

// NavigatorSummaryDTO is the count summary embedded in the topology graph response.
type NavigatorSummaryDTO struct {
	TotalProjects   int `json:"total_projects"`
	VerifiedCount   int `json:"verified_count"`
	UnverifiedCount int `json:"unverified_count"`
	TotalEdges      int `json:"total_edges"`
}

// NavigatorGraphDTO is the full response body of Navigator GET /topology/graph.
// Matches navigator/internal/topology/model.go Graph json tags exactly.
type NavigatorGraphDTO struct {
	CollectedAt string               `json:"collected_at"` // RFC3339
	Nodes       []NavigatorNodeDTO   `json:"nodes"`
	Edges       []NavigatorEdgeDTO   `json:"edges"`
	Summary     NavigatorSummaryDTO  `json:"summary"`
}

// ── NEXUS METRICS TYPES ───────────────────────────────────────────────────────

// NexusMetricsDTO is the response body of Nexus GET /metrics.
// Matches nexus/internal/telemetry/metrics.go Snapshot json tags exactly.
type NexusMetricsDTO struct {
	UptimeSeconds         float64 `json:"uptime_seconds"`
	ReconcileCyclesTotal  int64   `json:"reconcile_cycles_total"`
	ServicesStartedTotal  int64   `json:"services_started_total"`
	ServicesStoppedTotal  int64   `json:"services_stopped_total"`
	ServicesCrashedTotal  int64   `json:"services_crashed_total"`
	ServicesDeferredTotal int64   `json:"services_deferred_total"`
	ReconcileErrorsTotal  int64   `json:"reconcile_errors_total"`
	ServicesRunning       int64   `json:"services_running"`
	ServicesInMaintenance int64   `json:"services_in_maintenance"`
}

// ── GATE TYPES ────────────────────────────────────────────────────────────────
// ADR-042: Gate identity authority types.
// Verified against gate/internal/identity/token.go claim structure.

// IdentityClaimDTO is the validated actor identity returned by Gate POST /gate/validate.
// Embedded in execution records and event attribution.
// sub format: "<github_login>" for developers, "agent:<agent-id>" for agents.
type IdentityClaimDTO struct {
	Subject   string   `json:"sub"`
	Scopes    []string `json:"scp"`
	ExpiresAt int64    `json:"exp"`
	TokenID   string   `json:"jti"`
}

// HasScope returns true if the claim includes the given scope.
// Import canon.ScopeExecute etc. — never pass raw scope strings.
func (c *IdentityClaimDTO) HasScope(scope string) bool {
	for _, s := range c.Scopes {
		if s == scope {
			return true
		}
	}
	return false
}

// GateValidateRequest is the body for POST /gate/validate.
type GateValidateRequest struct {
	Token string `json:"token"`
}

// GateValidateResponse is the response body for POST /gate/validate.
// Valid=false carries a human-readable Reason — do not switch on Reason string.
type GateValidateResponse struct {
	Valid  bool              `json:"valid"`
	Claim  *IdentityClaimDTO `json:"claim,omitempty"`
	Reason string            `json:"reason,omitempty"`
}

// GatePublicKeyDTO is the response body for GET /gate/public-key.
// Key is a base64-encoded DER-encoded Ed25519 public key.
// Alg is always "Ed25519" — reject tokens if this value is unexpected.
type GatePublicKeyDTO struct {
	Key string `json:"key"`
	Alg string `json:"alg"`
}

// GateTokenResponse is returned by POST /gate/tokens/* and the OAuth callback.
type GateTokenResponse struct {
	Token     string `json:"token"`
	Subject   string `json:"sub"`
	ExpiresAt int64  `json:"exp"`
}
