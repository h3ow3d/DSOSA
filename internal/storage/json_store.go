package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type JSONStore struct {
	root         string
	catalogueDir string
	projectsDir  string
	assessDir    string
	improveDir   string
	eventsFile   string
}

func NewJSONStore(root string) (*JSONStore, error) {
	s := &JSONStore{
		root:         root,
		catalogueDir: filepath.Join(root, "catalogue"),
		projectsDir:  filepath.Join(root, "projects"),
		assessDir:    filepath.Join(root, "assessments"),
		improveDir:   filepath.Join(root, "improvements"),
		eventsFile:   filepath.Join(root, "events.ndjson"),
	}
	for _, dir := range []string{s.catalogueDir, s.projectsDir, s.assessDir, s.improveDir} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, fmt.Errorf("create %s: %w", dir, err)
		}
	}
	return s, nil
}

func (s *JSONStore) SaveCatalogue(record CatalogueRecord) error {
	path := filepath.Join(s.catalogueDir, record.Version+".json")
	return writeJSON(path, record)
}

func (s *JSONStore) ListCatalogueVersions() ([]CatalogueRecord, error) {
	entries, err := os.ReadDir(s.catalogueDir)
	if err != nil {
		return nil, err
	}
	out := make([]CatalogueRecord, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}
		var rec CatalogueRecord
		if err := readJSON(filepath.Join(s.catalogueDir, entry.Name()), &rec); err == nil {
			out = append(out, rec)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].FetchedAt.After(out[j].FetchedAt) })
	return out, nil
}

func (s *JSONStore) ReadCurrentCatalogue() (*CatalogueRecord, error) {
	items, err := s.ListCatalogueVersions()
	if err != nil || len(items) == 0 {
		return nil, err
	}
	return &items[0], nil
}

func (s *JSONStore) ReadCatalogueByHash(sha256 string) (*CatalogueRecord, error) {
	items, err := s.ListCatalogueVersions()
	if err != nil {
		return nil, err
	}
	for _, item := range items {
		if item.SHA256 == sha256 {
			return &item, nil
		}
	}
	return nil, errors.New("catalogue not found")
}

// Projects

func (s *JSONStore) SaveProject(p Project) error {
	return writeJSON(filepath.Join(s.projectsDir, p.ID+".json"), p)
}

func (s *JSONStore) GetProject(id string) (*Project, error) {
	var p Project
	if err := readJSON(filepath.Join(s.projectsDir, id+".json"), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func (s *JSONStore) ListProjects() []Project {
	return readCollection[Project](s.projectsDir)
}

// Assessments

func (s *JSONStore) SaveAssessment(a Assessment) error {
	return writeJSON(filepath.Join(s.assessDir, a.ID+".json"), a)
}

func (s *JSONStore) GetAssessment(id string) (*Assessment, error) {
	var a Assessment
	if err := readJSON(filepath.Join(s.assessDir, id+".json"), &a); err != nil {
		return nil, err
	}
	return &a, nil
}

func (s *JSONStore) ListAssessments() []Assessment {
	return readCollection[Assessment](s.assessDir)
}

func (s *JSONStore) ListAssessmentsByProject(projectID string) []Assessment {
	all := s.ListAssessments()
	out := make([]Assessment, 0)
	for _, a := range all {
		if a.ProjectID == projectID {
			out = append(out, a)
		}
	}
	return out
}

func (s *JSONStore) ListImprovements() []Improvement {
	return readCollection[Improvement](s.improveDir)
}

func (s *JSONStore) AppendEvent(event Event) error {
	return appendNDJSON(s.eventsFile, event)
}

func (s *JSONStore) CurrentControlCount() int {
	cat, err := s.ReadCurrentCatalogue()
	if err != nil || cat == nil {
		return 0
	}
	return countControlLikeMaps(cat.Body)
}

func countControlLikeMaps(value any) int {
	switch v := value.(type) {
	case []any:
		total := 0
		for _, item := range v {
			total += countControlLikeMaps(item)
		}
		return total
	case map[string]any:
		total := 0
		if _, hasID := v["id"]; hasID {
			total++
		}
		for _, item := range v {
			total += countControlLikeMaps(item)
		}
		return total
	default:
		return 0
	}
}

func writeJSON(path string, value any) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(value)
}

func readJSON(path string, target any) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewDecoder(f).Decode(target)
}

func readCollection[T any](dir string) []T {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	out := make([]T, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}
		var item T
		if err := readJSON(filepath.Join(dir, entry.Name()), &item); err == nil {
			out = append(out, item)
		}
	}
	return out
}
