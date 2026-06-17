package storage

import "time"

type CatalogueRecord struct {
Version   string         `json:"version"`
SHA256    string         `json:"sha256"`
FetchedAt time.Time      `json:"fetched_at"`
Body      map[string]any `json:"body"`
}

type Project struct {
ID        string    `json:"id"`
Name      string    `json:"name"`
CreatedAt time.Time `json:"created_at"`
UpdatedAt time.Time `json:"updated_at"`
}

type Assessment struct {
ID        string    `json:"id"`
ProjectID string    `json:"project_id"`
Name      string    `json:"name"`
CreatedAt time.Time `json:"created_at"`
UpdatedAt time.Time `json:"updated_at"`
}

type Improvement struct {
ID           string    `json:"id"`
AssessmentID string    `json:"assessment_id"`
Title        string    `json:"title"`
Status       string    `json:"status"`
UpdatedAt    time.Time `json:"updated_at"`
}
