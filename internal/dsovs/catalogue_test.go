package dsovs

import "testing"

func TestParsePhasesGroupsFlatControlsByPhase(t *testing.T) {
	body := map[string]any{
		"controls": []any{
			map[string]any{
				"id":        "ORG-001",
				"title":     "Risk Assessment",
				"objective": "Identify risk",
				"phase":     "Plan",
			},
			map[string]any{
				"id":        "IMP-001",
				"title":     "Improve Things",
				"objective": "Improve risk handling",
				"phase":     "Improve",
			},
		},
	}

	phases := ParsePhases(body)
	if len(phases) != 2 {
		t.Fatalf("phase count = %d, want 2", len(phases))
	}
	if phases[0].ID != "plan" || phases[0].Name != "Plan" {
		t.Fatalf("first phase = %#v, want Plan grouped phase", phases[0])
	}
	if len(phases[0].Controls) != 1 {
		t.Fatalf("plan control count = %d, want 1", len(phases[0].Controls))
	}
	if phases[0].Controls[0].ID != "ORG-001" || phases[0].Controls[0].Title != "Risk Assessment" {
		t.Fatalf("first control = %#v, want ORG-001 Risk Assessment", phases[0].Controls[0])
	}
}
