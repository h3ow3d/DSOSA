package app

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log/slog"
	"math"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"dsovs-assessment-tool/internal/dsovs"
	"dsovs-assessment-tool/internal/storage"
	"dsovs-assessment-tool/internal/web"
)

// ─── view models ──────────────────────────────────────────────────────────────

type controlRow struct {
	Control dsovs.Control
	Score   storage.ScoreEntry
}

func (s *Server) handleReports(w http.ResponseWriter, r *http.Request) {
	s.handleAssessments(w, r)
}

type phaseRow struct {
	Phase    dsovs.Phase
	Controls []controlRow
}

type phaseResult struct {
	Name    string
	Current float64
	Target  float64
	Gap     float64
}

type checklistItem struct {
	Label string
	Done  bool
}

// ─── route registration ───────────────────────────────────────────────────────

func (s *Server) registerRoutes(mux *http.ServeMux) {
	mux.Handle("GET /static/", http.StripPrefix("/static/", web.StaticHandler()))

	// Dashboard / catalogue
	mux.HandleFunc("GET /{$}", s.handleDashboard)
	mux.HandleFunc("GET /dashboard", s.handleDashboard)
	mux.HandleFunc("GET /catalogue", s.handleDashboard)
	mux.HandleFunc("POST /catalogue/sync", s.handleCatalogueSync)

	// Projects
	mux.HandleFunc("GET /projects", s.handleProjects)
	mux.HandleFunc("GET /projects/new", s.handleProjectNew)
	mux.HandleFunc("POST /projects", s.handleProjectCreate)
	mux.HandleFunc("GET /projects/{id}", s.handleProjectDetail)

	// Assessments
	mux.HandleFunc("GET /assessments", s.handleAssessments)
	mux.HandleFunc("GET /reports", s.handleReports)
	mux.HandleFunc("GET /projects/{id}/assessments/new", s.handleAssessmentNew)
	mux.HandleFunc("POST /projects/{id}/assessments", s.handleAssessmentCreate)
	mux.HandleFunc("GET /assessments/{id}", s.handleAssessmentDetail)
	mux.HandleFunc("POST /assessments/{id}/scores", s.handleScoreSave)
	mux.HandleFunc("GET /assessments/{id}/results", s.handleResults)
	mux.HandleFunc("GET /assessments/{id}/report", s.handleReport)
}

// ─── helpers ──────────────────────────────────────────────────────────────────

func newID() string {
	return fmt.Sprintf("%016x", time.Now().UnixNano())
}

func firstOrNil[T any](items []T) any {
	if len(items) == 0 {
		return nil
	}
	return items[0]
}

func syncMessage(changed bool) string {
	if changed {
		return "Catalogue synced successfully"
	}
	return "Catalogue already up to date"
}

func renderErr(w http.ResponseWriter, err error, msg string) {
	slog.Error(msg, "error", err)
	http.Error(w, msg, http.StatusInternalServerError)
}

func buildPhaseRows(phases []dsovs.Phase, scores []storage.ScoreEntry) []phaseRow {
	scoreMap := make(map[string]storage.ScoreEntry, len(scores))
	for _, sc := range scores {
		scoreMap[sc.ControlID] = sc
	}
	rows := make([]phaseRow, 0, len(phases))
	for _, ph := range phases {
		cr := make([]controlRow, 0, len(ph.Controls))
		for _, ctrl := range ph.Controls {
			sc := scoreMap[ctrl.ID]
			if sc.ControlID == "" {
				sc.ControlID = ctrl.ID
			}
			cr = append(cr, controlRow{Control: ctrl, Score: sc})
		}
		rows = append(rows, phaseRow{Phase: ph, Controls: cr})
	}
	return rows
}

func buildPhaseResults(phases []dsovs.Phase, scores []storage.ScoreEntry) []phaseResult {
	scoreMap := make(map[string]storage.ScoreEntry, len(scores))
	for _, sc := range scores {
		scoreMap[sc.ControlID] = sc
	}
	results := make([]phaseResult, 0, len(phases))
	for _, ph := range phases {
		var sumCur, sumTgt float64
		n := 0
		for _, ctrl := range ph.Controls {
			sc, ok := scoreMap[ctrl.ID]
			if !ok || sc.NotApplicable {
				continue
			}
			sumCur += float64(sc.Current)
			sumTgt += float64(sc.Target)
			n++
		}
		var cur, tgt float64
		if n > 0 {
			cur = math.Round(sumCur/float64(n)*100) / 100
			tgt = math.Round(sumTgt/float64(n)*100) / 100
		}
		results = append(results, phaseResult{
			Name:    ph.Name,
			Current: cur,
			Target:  tgt,
			Gap:     math.Round((tgt-cur)*100) / 100,
		})
	}
	return results
}

func overallScore(prs []phaseResult) (cur, tgt float64) {
	if len(prs) == 0 {
		return 0, 0
	}
	var sc, st float64
	for _, pr := range prs {
		sc += pr.Current
		st += pr.Target
	}
	n := float64(len(prs))
	return math.Round(sc/n*100) / 100, math.Round(st/n*100) / 100
}

func completionPct(phases []dsovs.Phase, scores []storage.ScoreEntry) int {
	scoreMap := make(map[string]bool, len(scores))
	for _, sc := range scores {
		if sc.Current > 0 || sc.NotApplicable {
			scoreMap[sc.ControlID] = true
		}
	}
	total := 0
	scored := 0
	for _, ph := range phases {
		for _, c := range ph.Controls {
			total++
			if scoreMap[c.ID] {
				scored++
			}
		}
	}
	if total == 0 {
		return 0
	}
	return scored * 100 / total
}

type gapItem struct {
	ControlID string
	Title     string
	Gap       int
	Priority  string
}

func topGaps(phases []dsovs.Phase, scores []storage.ScoreEntry, n int) []gapItem {
	scoreMap := make(map[string]storage.ScoreEntry)
	for _, sc := range scores {
		scoreMap[sc.ControlID] = sc
	}
	titleMap := make(map[string]string)
	for _, ph := range phases {
		for _, c := range ph.Controls {
			titleMap[c.ID] = c.Title
		}
	}
	items := make([]gapItem, 0)
	for _, sc := range scores {
		if sc.NotApplicable {
			continue
		}
		gap := sc.Target - sc.Current
		if gap > 0 {
			items = append(items, gapItem{
				ControlID: sc.ControlID,
				Title:     titleMap[sc.ControlID],
				Gap:       gap,
				Priority:  sc.Priority,
			})
		}
	}
	sort.Slice(items, func(i, j int) bool { return items[i].Gap > items[j].Gap })
	if n > len(items) {
		n = len(items)
	}
	return items[:n]
}

func missingEvidence(phases []dsovs.Phase, scores []storage.ScoreEntry) []controlRow {
	scoreMap := make(map[string]storage.ScoreEntry)
	for _, sc := range scores {
		scoreMap[sc.ControlID] = sc
	}
	var out []controlRow
	for _, ph := range phases {
		for _, c := range ph.Controls {
			sc := scoreMap[c.ID]
			if sc.Current > 0 && !sc.NotApplicable && strings.TrimSpace(sc.Evidence) == "" {
				out = append(out, controlRow{Control: c, Score: sc})
			}
		}
	}
	return out
}

// radarSVG generates a dependency-free SVG radar/star chart.
func radarSVG(prs []phaseResult) template.HTML {
	if len(prs) == 0 {
		return `<svg viewBox="0 0 400 400" width="400" height="400"><text x="200" y="200" text-anchor="middle" fill="#718096">No data yet</text></svg>`
	}
	const (
		w    = 400
		h    = 400
		cx   = 200.0
		cy   = 200.0
		maxR = 140.0
		maxV = 3.0
	)
	n := len(prs)
	var b strings.Builder

	fmt.Fprintf(&b, `<svg viewBox="0 0 %d %d" width="%d" height="%d" xmlns="http://www.w3.org/2000/svg">`, w, h, w, h)

	// Grid polygons at L1, L2, L3
	for level := 1; level <= 3; level++ {
		r := (float64(level) / maxV) * maxR
		var pts []string
		for i := 0; i < n; i++ {
			angle := 2*math.Pi*float64(i)/float64(n) - math.Pi/2
			x := cx + r*math.Cos(angle)
			y := cy + r*math.Sin(angle)
			pts = append(pts, fmt.Sprintf("%.1f,%.1f", x, y))
		}
		fmt.Fprintf(&b, `<polygon points="%s" fill="none" stroke="#cbd5e0" stroke-width="1"/>`, strings.Join(pts, " "))
		// Level label
		angle0 := -math.Pi / 2
		lx := cx + (float64(level)/maxV)*maxR*math.Cos(angle0) + 4
		ly := cy + (float64(level)/maxV)*maxR*math.Sin(angle0)
		fmt.Fprintf(&b, `<text x="%.1f" y="%.1f" font-size="9" fill="#a0aec0">%d</text>`, lx, ly, level)
	}

	// Axes and labels
	for i, pr := range prs {
		angle := 2*math.Pi*float64(i)/float64(n) - math.Pi/2
		x2 := cx + maxR*math.Cos(angle)
		y2 := cy + maxR*math.Sin(angle)
		fmt.Fprintf(&b, `<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="#cbd5e0" stroke-width="1"/>`, cx, cy, x2, y2)

		lx := cx + (maxR+22)*math.Cos(angle)
		ly := cy + (maxR+22)*math.Sin(angle)
		anchor := "middle"
		cosA := math.Cos(angle)
		if cosA > 0.15 {
			anchor = "start"
		} else if cosA < -0.15 {
			anchor = "end"
		}
		// Wrap long labels
		label := pr.Name
		if len(label) > 18 {
			label = label[:16] + "…"
		}
		fmt.Fprintf(&b, `<text x="%.1f" y="%.1f" text-anchor="%s" dominant-baseline="middle" font-size="10" fill="#4a5568">%s</text>`,
			lx, ly, anchor, template.HTMLEscapeString(label))
	}

	// Target polygon (dashed, blue)
	var tgtPts []string
	for i, pr := range prs {
		angle := 2*math.Pi*float64(i)/float64(n) - math.Pi/2
		r := (pr.Target / maxV) * maxR
		x := cx + r*math.Cos(angle)
		y := cy + r*math.Sin(angle)
		tgtPts = append(tgtPts, fmt.Sprintf("%.1f,%.1f", x, y))
	}
	fmt.Fprintf(&b, `<polygon points="%s" fill="#bee3f8" fill-opacity="0.4" stroke="#3182ce" stroke-width="2" stroke-dasharray="5,3"/>`, strings.Join(tgtPts, " "))

	// Current polygon (solid, green)
	var curPts []string
	for i, pr := range prs {
		angle := 2*math.Pi*float64(i)/float64(n) - math.Pi/2
		r := (pr.Current / maxV) * maxR
		x := cx + r*math.Cos(angle)
		y := cy + r*math.Sin(angle)
		curPts = append(curPts, fmt.Sprintf("%.1f,%.1f", x, y))
	}
	fmt.Fprintf(&b, `<polygon points="%s" fill="#c6f6d5" fill-opacity="0.7" stroke="#38a169" stroke-width="2"/>`, strings.Join(curPts, " "))

	// Legend
	fmt.Fprintf(&b, `<rect x="10" y="10" width="14" height="4" fill="#38a169"/>`)
	fmt.Fprintf(&b, `<text x="28" y="16" font-size="11" fill="#4a5568">Current</text>`)
	fmt.Fprintf(&b, `<rect x="10" y="24" width="14" height="4" fill="#3182ce"/>`)
	fmt.Fprintf(&b, `<text x="28" y="30" font-size="11" fill="#4a5568">Target</text>`)

	fmt.Fprintf(&b, `</svg>`)
	return template.HTML(b.String()) //nolint:gosec // computed SVG
}

// ─── handlers ─────────────────────────────────────────────────────────────────

func (s *Server) handleDashboard(w http.ResponseWriter, r *http.Request) {
	catalogues, _ := s.store.ListCatalogueVersions()
	sort.Slice(catalogues, func(i, j int) bool { return catalogues[i].FetchedAt.After(catalogues[j].FetchedAt) })

	var (
		hasCat     = len(catalogues) > 0
		catVer     = ""
		lastCat    any
		phases     []dsovs.Phase
		phaseCount int
		ctrlCount  int
		lastSynced time.Time
	)
	if hasCat {
		lastCat = catalogues[0]
		phases = dsovs.ParsePhases(catalogues[0].Body)
		phaseCount = len(phases)
		for _, ph := range phases {
			ctrlCount += len(ph.Controls)
		}
		catVer = dsovs.CatalogueVersion(catalogues[0].Body)
		if catVer == "" {
			catVer = catalogues[0].Version
		}
		lastSynced = catalogues[0].FetchedAt
	}

	projects := s.store.ListProjects()
	assessments := s.store.ListAssessments()

	hasScores := false
	for _, a := range assessments {
		if len(a.Scores) > 0 {
			hasScores = true
			break
		}
	}

	checklist := []checklistItem{
		{"Sync catalogue", hasCat},
		{"Create project", len(projects) > 0},
		{"Create assessment", len(assessments) > 0},
		{"Score controls", hasScores},
		{"View results", hasScores},
		{"Print report", false},
	}

	data := map[string]any{
		"Title":            "Dashboard",
		"Nav":              "dashboard",
		"HasCatalogue":     hasCat,
		"CatalogueVersion": catVer,
		"PhaseCount":       phaseCount,
		"ControlCount":     ctrlCount,
		"LastSynced":       lastSynced,
		"LastCatalogue":    lastCat,
		"ProjectCount":     len(projects),
		"AssessmentCount":  len(assessments),
		"SyncMessage":      r.URL.Query().Get("synced"),
		"Checklist":        checklist,
	}
	if err := s.renderer.Render(w, "dashboard", data); err != nil {
		renderErr(w, err, "dashboard render failed")
	}
}

func (s *Server) handleCatalogueSync(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	result, err := dsovs.Sync(ctx, s.client, s.store)
	if err != nil {
		slog.Error("catalogue sync failed", "error", err)
		http.Error(w, "sync failed", http.StatusBadGateway)
		return
	}
	slog.Info("catalogue synced", "version", result.Version, "changed", result.Changed)
	http.Redirect(w, r, "/dashboard?synced="+url.QueryEscape(syncMessage(result.Changed)), http.StatusSeeOther)
}

// ─── projects ─────────────────────────────────────────────────────────────────

func (s *Server) handleProjects(w http.ResponseWriter, r *http.Request) {
	projects := s.store.ListProjects()
	sort.Slice(projects, func(i, j int) bool { return projects[i].UpdatedAt.After(projects[j].UpdatedAt) })
	data := map[string]any{
		"Title":    "Projects",
		"Nav":      "projects",
		"Projects": projects,
	}
	if err := s.renderer.Render(w, "projects", data); err != nil {
		renderErr(w, err, "projects render failed")
	}
}

func (s *Server) handleProjectNew(w http.ResponseWriter, r *http.Request) {
	data := map[string]any{
		"Title":        "New Project",
		"Nav":          "projects",
		"ErrorSummary": strings.TrimSpace(r.URL.Query().Get("error")),
		"Form": map[string]string{
			"name":        strings.TrimSpace(r.URL.Query().Get("name")),
			"client_name": strings.TrimSpace(r.URL.Query().Get("client_name")),
			"owner":       strings.TrimSpace(r.URL.Query().Get("owner")),
			"description": strings.TrimSpace(r.URL.Query().Get("description")),
		},
	}
	if err := s.renderer.Render(w, "project_new", data); err != nil {
		renderErr(w, err, "project_new render failed")
	}
}

func (s *Server) handleProjectCreate(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}
	now := time.Now().UTC()
	p := storage.Project{
		ID:          newID(),
		Name:        strings.TrimSpace(r.FormValue("name")),
		ClientName:  strings.TrimSpace(r.FormValue("client_name")),
		Owner:       strings.TrimSpace(r.FormValue("owner")),
		Description: strings.TrimSpace(r.FormValue("description")),
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if p.Name == "" {
		q := make(url.Values)
		q.Set("error", "Project name is required")
		q.Set("name", p.Name)
		q.Set("client_name", p.ClientName)
		q.Set("owner", p.Owner)
		q.Set("description", p.Description)
		http.Redirect(w, r, "/projects/new?"+q.Encode(), http.StatusSeeOther)
		return
	}
	if err := s.store.SaveProject(p); err != nil {
		renderErr(w, err, "save project failed")
		return
	}
	_ = s.store.AppendEvent(storage.Event{
		Type:    "project.created",
		Time:    now,
		Payload: map[string]any{"id": p.ID, "name": p.Name},
	})
	slog.Info("project created", "id", p.ID, "name", p.Name)
	http.Redirect(w, r, "/projects/"+p.ID, http.StatusSeeOther)
}

func (s *Server) handleProjectDetail(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	p, err := s.store.GetProject(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	assessments := s.store.ListAssessmentsByProject(id)
	sort.Slice(assessments, func(i, j int) bool { return assessments[i].CreatedAt.After(assessments[j].CreatedAt) })
	data := map[string]any{
		"Title":       p.Name,
		"Nav":         "projects",
		"Project":     p,
		"Assessments": assessments,
	}
	if err := s.renderer.Render(w, "project", data); err != nil {
		renderErr(w, err, "project render failed")
	}
}

// ─── assessments ──────────────────────────────────────────────────────────────

func (s *Server) handleAssessments(w http.ResponseWriter, r *http.Request) {
	assessments := s.store.ListAssessments()
	sort.Slice(assessments, func(i, j int) bool { return assessments[i].CreatedAt.After(assessments[j].CreatedAt) })

	// Build project lookup map
	projects := s.store.ListProjects()
	projMap := make(map[string]storage.Project, len(projects))
	for _, p := range projects {
		projMap[p.ID] = p
	}

	data := map[string]any{
		"Title":       "Assessments",
		"Nav":         "assessments",
		"Assessments": assessments,
		"Projects":    projMap,
	}
	if err := s.renderer.Render(w, "assessments", data); err != nil {
		renderErr(w, err, "assessments render failed")
	}
}

func (s *Server) handleAssessmentNew(w http.ResponseWriter, r *http.Request) {
	pid := r.PathValue("id")
	p, err := s.store.GetProject(pid)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	catalogues, _ := s.store.ListCatalogueVersions()
	data := map[string]any{
		"Title":        "New Assessment",
		"Nav":          "assessments",
		"Project":      p,
		"Catalogues":   catalogues,
		"ErrorSummary": strings.TrimSpace(r.URL.Query().Get("error")),
		"Form": map[string]string{
			"name":            strings.TrimSpace(r.URL.Query().Get("name")),
			"assessment_date": strings.TrimSpace(r.URL.Query().Get("assessment_date")),
			"assessor":        strings.TrimSpace(r.URL.Query().Get("assessor")),
			"scope":           strings.TrimSpace(r.URL.Query().Get("scope")),
		},
	}
	if err := s.renderer.Render(w, "assessment_new", data); err != nil {
		renderErr(w, err, "assessment_new render failed")
	}
}

func (s *Server) handleAssessmentCreate(w http.ResponseWriter, r *http.Request) {
	pid := r.PathValue("id")
	if _, err := s.store.GetProject(pid); err != nil {
		http.NotFound(w, r)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}

	cat, _ := s.store.ReadCurrentCatalogue()
	catHash := ""
	catVer := ""
	if cat != nil {
		catHash = cat.SHA256
		catVer = dsovs.CatalogueVersion(cat.Body)
		if catVer == "" {
			catVer = cat.Version
		}
	}

	now := time.Now().UTC()
	a := storage.Assessment{
		ID:              newID(),
		ProjectID:       pid,
		StandardVersion: catVer,
		CatalogueHash:   catHash,
		Name:            strings.TrimSpace(r.FormValue("name")),
		AssessmentDate:  strings.TrimSpace(r.FormValue("assessment_date")),
		Assessor:        strings.TrimSpace(r.FormValue("assessor")),
		Scope:           strings.TrimSpace(r.FormValue("scope")),
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	if a.Name == "" {
		q := make(url.Values)
		q.Set("error", "Assessment name is required")
		q.Set("name", a.Name)
		q.Set("assessment_date", a.AssessmentDate)
		q.Set("assessor", a.Assessor)
		q.Set("scope", a.Scope)
		http.Redirect(w, r, "/projects/"+pid+"/assessments/new?"+q.Encode(), http.StatusSeeOther)
		return
	}
	if err := s.store.SaveAssessment(a); err != nil {
		renderErr(w, err, "save assessment failed")
		return
	}
	_ = s.store.AppendEvent(storage.Event{
		Type:    "assessment.created",
		Time:    now,
		Payload: map[string]any{"id": a.ID, "project_id": pid, "name": a.Name},
	})
	slog.Info("assessment created", "id", a.ID, "project_id", pid)
	http.Redirect(w, r, "/assessments/"+a.ID, http.StatusSeeOther)
}

func (s *Server) handleAssessmentDetail(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	a, err := s.store.GetAssessment(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	p, _ := s.store.GetProject(a.ProjectID)

	cat, _ := s.store.ReadCurrentCatalogue()
	var phases []dsovs.Phase
	if cat != nil {
		phases = dsovs.ParsePhases(cat.Body)
	}

	phaseRows := buildPhaseRows(phases, a.Scores)

	total := 0
	for _, ph := range phases {
		total += len(ph.Controls)
	}

	data := map[string]any{
		"Title":        a.Name,
		"Nav":          "assessments",
		"Assessment":   a,
		"Project":      p,
		"Phases":       phaseRows,
		"Levels":       []int{0, 1, 2, 3},
		"Priorities":   []string{"", "high", "medium", "low"},
		"Confidences":  []string{"", "high", "medium", "low"},
		"ScoreCount":   len(a.Scores),
		"ControlCount": total,
		"HasCatalogue": cat != nil,
	}
	if err := s.renderer.Render(w, "assessment", data); err != nil {
		renderErr(w, err, "assessment render failed")
	}
}

func (s *Server) handleScoreSave(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	a, err := s.store.GetAssessment(id)
	if err != nil {
		http.Error(w, "assessment not found", http.StatusNotFound)
		return
	}

	var scores []storage.ScoreEntry
	if err := json.NewDecoder(r.Body).Decode(&scores); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	a.Scores = scores
	a.UpdatedAt = time.Now().UTC()
	if err := s.store.SaveAssessment(*a); err != nil {
		renderErr(w, err, "save scores failed")
		return
	}

	_ = s.store.AppendEvent(storage.Event{
		Type:    "score.updated",
		Time:    a.UpdatedAt,
		Payload: map[string]any{"assessment_id": id, "count": len(scores)},
	})
	slog.Info("scores saved", "assessment_id", id, "count", len(scores))

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"message": "Scores saved",
		"count":   len(scores),
	})
}

// ─── results ──────────────────────────────────────────────────────────────────

func (s *Server) handleResults(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	a, err := s.store.GetAssessment(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	p, _ := s.store.GetProject(a.ProjectID)

	cat, _ := s.store.ReadCurrentCatalogue()
	var phases []dsovs.Phase
	if cat != nil {
		phases = dsovs.ParsePhases(cat.Body)
	}

	prs := buildPhaseResults(phases, a.Scores)
	curScore, tgtScore := overallScore(prs)
	pct := completionPct(phases, a.Scores)
	gaps := topGaps(phases, a.Scores, 10)
	missing := missingEvidence(phases, a.Scores)
	chart := radarSVG(prs)

	data := map[string]any{
		"Title":           "Results: " + a.Name,
		"Nav":             "assessments",
		"Assessment":      a,
		"Project":         p,
		"OverallCurrent":  curScore,
		"OverallTarget":   tgtScore,
		"CompletionPct":   pct,
		"PhaseResults":    prs,
		"TopGaps":         gaps,
		"MissingEvidence": missing,
		"RadarSVG":        chart,
		"HasScores":       len(a.Scores) > 0,
	}
	if err := s.renderer.Render(w, "results", data); err != nil {
		renderErr(w, err, "results render failed")
	}
}

// ─── report ───────────────────────────────────────────────────────────────────

func (s *Server) handleReport(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	a, err := s.store.GetAssessment(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	p, _ := s.store.GetProject(a.ProjectID)

	cat, _ := s.store.ReadCurrentCatalogue()
	var phases []dsovs.Phase
	if cat != nil {
		phases = dsovs.ParsePhases(cat.Body)
	}

	prs := buildPhaseResults(phases, a.Scores)
	curScore, tgtScore := overallScore(prs)
	gaps := topGaps(phases, a.Scores, 10)
	phaseRows := buildPhaseRows(phases, a.Scores)
	chart := radarSVG(prs)

	_ = s.store.AppendEvent(storage.Event{
		Type:    "report.viewed",
		Time:    time.Now().UTC(),
		Payload: map[string]any{"assessment_id": id},
	})

	data := map[string]any{
		"Title":          "Report: " + a.Name,
		"Nav":            "assessments",
		"Assessment":     a,
		"Project":        p,
		"OverallCurrent": curScore,
		"OverallTarget":  tgtScore,
		"PhaseResults":   prs,
		"TopGaps":        gaps,
		"Phases":         phaseRows,
		"RadarSVG":       chart,
		"HasScores":      len(a.Scores) > 0,
		"PrintedAt":      time.Now().UTC().Format("2 Jan 2006 15:04 UTC"),
	}
	if err := s.renderer.Render(w, "report", data); err != nil {
		renderErr(w, err, "report render failed")
	}
}
