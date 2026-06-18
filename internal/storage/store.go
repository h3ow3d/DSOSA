package storage

import "time"

type CatalogueRecord struct {
	Version   string         `json:"version"`
	SHA256    string         `json:"sha256"`
	FetchedAt time.Time      `json:"fetched_at"`
	Body      map[string]any `json:"body"`
}

type Project struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	ClientName  string    `json:"client_name"`
	Owner       string    `json:"owner"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ScoreEntry struct {
	ControlID     string    `json:"control_id"`
	CurrentLevel  *int      `json:"current_level,omitempty"`
	TargetLevel   *int      `json:"target_level,omitempty"`
	NotApplicable bool      `json:"not_applicable"`
	EvidenceNotes string    `json:"evidence_notes"`
	ActionNotes   string    `json:"action_notes"`
	Priority      string    `json:"priority"`
	Confidence    string    `json:"confidence"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Assessment struct {
	ID              string       `json:"id"`
	ProjectID       string       `json:"project_id"`
	StandardVersion string       `json:"standard_version"`
	CatalogueHash   string       `json:"catalogue_hash"`
	Name            string       `json:"name"`
	AssessmentDate  string       `json:"assessment_date"`
	Assessor        string       `json:"assessor"`
	Scope           string       `json:"scope"`
	CreatedAt       time.Time    `json:"created_at"`
	UpdatedAt       time.Time    `json:"updated_at"`
	Scores          []ScoreEntry `json:"scores,omitempty"`
}

type Improvement struct {
	ID           string    `json:"id"`
	AssessmentID string    `json:"assessment_id"`
	Title        string    `json:"title"`
	Status       string    `json:"status"`
	UpdatedAt    time.Time `json:"updated_at"`
}
