package app

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"dsovs-assessment-tool/internal/storage"
)

func TestScoreSaveUpdatesAssessmentJSON(t *testing.T) {
	s := newTestServer(t)
	projectID := seedProject(t, s)
	assessmentID := seedAssessment(t, s, projectID)

	form := url.Values{}
	form.Add("control_id", "C-1")
	form.Set("current_level[C-1]", "2")
	form.Set("target_level[C-1]", "3")
	form.Set("evidence_notes[C-1]", "Updated evidence")
	form.Set("action_notes[C-1]", "Updated action")
	form.Set("priority[C-1]", "high")
	form.Set("confidence[C-1]", "medium")

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/assessments/"+assessmentID+"/scores", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	s.http.Handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusSeeOther {
		t.Fatalf("POST scores: got %d want %d", rec.Code, http.StatusSeeOther)
	}

	saved, err := s.store.GetAssessment(assessmentID)
	if err != nil {
		t.Fatalf("GetAssessment: %v", err)
	}
	if len(saved.Scores) != 1 {
		t.Fatalf("scores length = %d, want 1", len(saved.Scores))
	}
	sc := saved.Scores[0]
	if sc.CurrentLevel == nil || *sc.CurrentLevel != 2 {
		t.Fatalf("current level = %v, want 2", sc.CurrentLevel)
	}
	if sc.TargetLevel == nil || *sc.TargetLevel != 3 {
		t.Fatalf("target level = %v, want 3", sc.TargetLevel)
	}
	if sc.EvidenceNotes != "Updated evidence" {
		t.Fatalf("evidence = %q, want updated", sc.EvidenceNotes)
	}
	if sc.UpdatedAt.IsZero() {
		t.Fatal("score UpdatedAt was not set")
	}

	b, err := os.ReadFile(filepath.Join(s.cfg.DataDir, "assessments", assessmentID+".json"))
	if err != nil {
		t.Fatalf("ReadFile assessment json: %v", err)
	}
	var persisted storage.Assessment
	if err := json.Unmarshal(b, &persisted); err != nil {
		t.Fatalf("Unmarshal assessment json: %v", err)
	}
	if len(persisted.Scores) != 1 || persisted.Scores[0].CurrentLevel == nil || *persisted.Scores[0].CurrentLevel != 2 {
		t.Fatal("persisted assessment JSON did not include saved score")
	}

	eventBytes, err := os.ReadFile(filepath.Join(s.cfg.DataDir, "events.ndjson"))
	if err != nil {
		t.Fatalf("ReadFile events: %v", err)
	}
	if !strings.Contains(string(eventBytes), `"type":"score.updated"`) {
		t.Fatal("score.updated event was not appended")
	}
}

func TestSavedScoresReloadOnAssessmentPage(t *testing.T) {
	s := newTestServer(t)
	projectID := seedProject(t, s)
	assessmentID := seedAssessment(t, s, projectID)

	form := url.Values{}
	form.Add("control_id", "C-1")
	form.Set("current_level[C-1]", "0")
	form.Set("target_level[C-1]", "3")
	form.Set("evidence_notes[C-1]", "Evidence A")

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/assessments/"+assessmentID+"/scores", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	s.http.Handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusSeeOther {
		t.Fatalf("POST scores: got %d want %d", rec.Code, http.StatusSeeOther)
	}

	page := getPageBody(t, s, "/assessments/"+assessmentID)
	if !strings.Contains(page, `name="current_level[C-1]"`) || !strings.Contains(page, `id="current-C-1-0"`) || !strings.Contains(page, `id="current-C-1-0" type="radio" name="current_level[C-1]" value="0" checked`) {
		t.Fatal("assessment page did not reload saved current level")
	}
	if !strings.Contains(page, `name="target_level[C-1]"`) || !strings.Contains(page, `id="target-C-1-3" type="radio" name="target_level[C-1]" value="3" checked`) {
		t.Fatal("assessment page did not reload saved target level")
	}
	if !strings.Contains(page, "Evidence A") {
		t.Fatal("assessment page did not reload saved evidence notes")
	}
}

func TestInvalidScoreValuesAreRejected(t *testing.T) {
	s := newTestServer(t)
	projectID := seedProject(t, s)
	assessmentID := seedAssessment(t, s, projectID)

	form := url.Values{}
	form.Add("control_id", "C-1")
	form.Set("current_level[C-1]", "9")

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/assessments/"+assessmentID+"/scores", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	s.http.Handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("invalid level status = %d, want %d", rec.Code, http.StatusBadRequest)
	}

	badControlForm := url.Values{}
	badControlForm.Add("control_id", "UNKNOWN")
	rec = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/assessments/"+assessmentID+"/scores", strings.NewReader(badControlForm.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	s.http.Handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("unknown control status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestResultsUseSavedScores(t *testing.T) {
	s := newTestServer(t)
	projectID := seedProject(t, s)
	assessmentID := seedAssessment(t, s, projectID)

	form := url.Values{}
	form.Add("control_id", "C-1")
	form.Set("current_level[C-1]", "1")
	form.Set("target_level[C-1]", "3")
	form.Set("priority[C-1]", "high")

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/assessments/"+assessmentID+"/scores", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	s.http.Handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusSeeOther {
		t.Fatalf("POST scores: got %d want %d", rec.Code, http.StatusSeeOther)
	}

	results := getPageBody(t, s, "/assessments/"+assessmentID+"/results")
	for _, want := range []string{
		"100%",
		"<td>Plan</td>",
		"<td>1</td>",
		"<td>3</td>",
		"<code>C-1</code>",
	} {
		if !strings.Contains(results, want) {
			t.Fatalf("results page missing %q", want)
		}
	}
}

func TestAssessmentPageRendersControlScoringCardFields(t *testing.T) {
	s := newTestServer(t)
	projectID := seedProject(t, s)
	assessmentID := seedAssessment(t, s, projectID)

	page := getPageBody(t, s, "/assessments/"+assessmentID)
	for _, want := range []string{
		`class="wcag-control-card"`,
		`name="current_level[C-1]"`,
		`name="target_level[C-1]"`,
		`name="evidence_notes[C-1]"`,
		`name="action_notes[C-1]"`,
	} {
		if !strings.Contains(page, want) {
			t.Fatalf("assessment page missing %q", want)
		}
	}
}

func TestReportIncludesSavedControlData(t *testing.T) {
	s := newTestServer(t)
	projectID := seedProject(t, s)
	assessmentID := seedAssessment(t, s, projectID)

	form := url.Values{}
	form.Add("control_id", "C-1")
	form.Set("current_level[C-1]", "1")
	form.Set("target_level[C-1]", "2")
	form.Set("evidence_notes[C-1]", "Evidence for report")
	form.Set("action_notes[C-1]", "Action for report")

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/assessments/"+assessmentID+"/scores", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	s.http.Handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusSeeOther {
		t.Fatalf("POST scores: got %d want %d", rec.Code, http.StatusSeeOther)
	}

	report := getPageBody(t, s, "/assessments/"+assessmentID+"/report")
	for _, want := range []string{
		"<code>C-1</code>",
		"Evidence for report",
		"Action for report",
		"<td>1</td>",
		"<td>2</td>",
	} {
		if !strings.Contains(report, want) {
			t.Fatalf("report page missing %q", want)
		}
	}
}
