package dsovs

import (
"context"
"crypto/sha256"
"encoding/json"
"fmt"
"time"

"dsovs-assessment-tool/internal/storage"
)

func Sync(ctx context.Context, client *Client, store *storage.JSONStore) (SyncResult, error) {
raw, err := client.Fetch(ctx)
if err != nil {
return SyncResult{}, err
}

sum := sha256.Sum256(raw)
sha := fmt.Sprintf("%x", sum[:])
version := time.Now().UTC().Format("20060102T150405Z") + "-" + sha[:12]
current, _ := store.ReadCurrentCatalogue()
if current != nil && current.SHA256 == sha {
return SyncResult{Version: current.Version, Changed: false}, nil
}

body := map[string]any{}
if err := json.Unmarshal(raw, &body); err != nil {
return SyncResult{}, fmt.Errorf("invalid dsovs json: %w", err)
}

catalogue := storage.CatalogueRecord{Version: version, SHA256: sha, FetchedAt: time.Now().UTC(), Body: body}
if err := store.SaveCatalogue(catalogue); err != nil {
return SyncResult{}, err
}

_ = store.AppendEvent(storage.Event{Type: "catalogue.synced", Time: time.Now().UTC(), Payload: map[string]any{"version": version, "sha256": sha}})
return SyncResult{Version: version, Changed: true}, nil
}
