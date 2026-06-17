package storage

import (
	"encoding/json"
	"os"
)

type Event struct {
	Type    string         `json:"type"`
	Time    any            `json:"time"`
	Payload map[string]any `json:"payload,omitempty"`
}

func appendNDJSON(path string, event any) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	return enc.Encode(event)
}
