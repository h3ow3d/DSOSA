package dsovs

import "time"

type Catalogue struct {
	Version   string         `json:"version"`
	SHA256    string         `json:"sha256"`
	FetchedAt time.Time      `json:"fetched_at"`
	Raw       []byte         `json:"raw"`
	Body      map[string]any `json:"body"`
}

type SyncResult struct {
	Version string `json:"version"`
	Changed bool   `json:"changed"`
}
