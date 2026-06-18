package dsovs

// Phase represents a DSOVS phase (group/category) of controls.
type Phase struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	Controls []Control `json:"controls"`
}

// Control represents a single DSOVS security control.
type Control struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Summary string `json:"summary"`
	DocURL  string `json:"doc_url"`
	Type    string `json:"type"`
	Status  string `json:"status"`
	Level0  string `json:"level_0"`
	Level1  string `json:"level_1"`
	Level2  string `json:"level_2"`
	Level3  string `json:"level_3"`
}

// ParsePhases extracts phases and controls from a catalogue body map.
// It tries several common DSOVS JSON structures in order.
func ParsePhases(body map[string]any) []Phase {
	// Candidate top-level array keys, in preference order
	for _, key := range []string{"phases", "groups", "categories", "sections", "domains"} {
		if raw, ok := body[key]; ok {
			if phases := extractPhases(raw); len(phases) > 0 {
				return phases
			}
		}
	}
	// Fallback: look for any top-level []any whose items have controls
	for _, v := range body {
		if phases := extractPhases(v); len(phases) > 0 {
			return phases
		}
	}
	return nil
}

// CatalogueVersion attempts to read a version string from the catalogue body.
func CatalogueVersion(body map[string]any) string {
	for _, metaKey := range []string{"document", "metadata", "info"} {
		if meta, ok := body[metaKey].(map[string]any); ok {
			for _, vk := range []string{"version", "revision", "release"} {
				if v, ok := meta[vk].(string); ok && v != "" {
					return v
				}
			}
		}
	}
	for _, vk := range []string{"version", "revision", "release"} {
		if v, ok := body[vk].(string); ok && v != "" {
			return v
		}
	}
	return ""
}

func extractPhases(raw any) []Phase {
	items, ok := raw.([]any)
	if !ok {
		return nil
	}
	phases := make([]Phase, 0, len(items))
	for _, item := range items {
		m, ok := item.(map[string]any)
		if !ok {
			continue
		}
		phase := Phase{
			ID:   strField(m, "id", "phase_id", "code"),
			Name: strField(m, "name", "title", "phase_name", "label"),
		}
		// Extract controls from known sub-keys
		for _, ck := range []string{"controls", "items", "requirements", "practices"} {
			if cs, ok := m[ck].([]any); ok {
				phase.Controls = extractControls(cs)
				break
			}
		}
		if len(phase.Controls) > 0 || phase.Name != "" {
			phases = append(phases, phase)
		}
	}
	return phases
}

func extractControls(items []any) []Control {
	controls := make([]Control, 0, len(items))
	for _, item := range items {
		m, ok := item.(map[string]any)
		if !ok {
			continue
		}
		c := Control{
			ID:      strField(m, "id", "code", "control_id"),
			Title:   strField(m, "title", "name", "control", "heading"),
			Summary: strField(m, "summary", "objective", "description", "statement"),
			DocURL:  docURL(m),
			Type:    strField(m, "type", "control_type", "category"),
			Status:  strField(m, "status", "state"),
			Level0:  strField(m, "L0", "level_0", "level0", "l0", "maturity_0"),
			Level1:  strField(m, "L1", "level_1", "level1", "l1", "maturity_1"),
			Level2:  strField(m, "L2", "level_2", "level2", "l2", "maturity_2"),
			Level3:  strField(m, "L3", "level_3", "level3", "l3", "maturity_3"),
		}
		if c.ID != "" || c.Title != "" {
			controls = append(controls, c)
		}
	}
	return controls
}

func strField(m map[string]any, keys ...string) string {
	for _, k := range keys {
		if v, ok := m[k].(string); ok && v != "" {
			return v
		}
	}
	return ""
}

func docURL(m map[string]any) string {
	// Direct string URL fields
	for _, k := range []string{"docURL", "doc_url", "url", "link", "href", "reference"} {
		if v, ok := m[k].(string); ok && v != "" {
			return v
		}
	}
	// References might be []any of strings or maps
	if refs, ok := m["references"].([]any); ok && len(refs) > 0 {
		switch v := refs[0].(type) {
		case string:
			return v
		case map[string]any:
			return strField(v, "url", "href", "link")
		}
	}
	return ""
}
