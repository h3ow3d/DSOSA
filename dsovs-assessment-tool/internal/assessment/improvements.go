package assessment

type Recommendation struct {
ControlID string `json:"control_id"`
Priority  string `json:"priority"`
Action    string `json:"action"`
}

func BuildRecommendations(gaps []GapItem) []Recommendation {
recs := make([]Recommendation, 0, len(gaps))
for _, gap := range gaps {
priority := "low"
switch {
case gap.Gap >= 2:
priority = "high"
case gap.Gap >= 1:
priority = "medium"
}
recs = append(recs, Recommendation{ControlID: gap.ControlID, Priority: priority, Action: "Increase maturity controls"})
}
return recs
}
