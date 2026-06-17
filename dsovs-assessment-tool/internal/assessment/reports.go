package assessment

import "sort"

type GapItem struct {
ControlID string  `json:"control_id"`
Gap       float64 `json:"gap"`
}

func TopGaps(scores []ControlScore, n int) []GapItem {
items := make([]GapItem, 0, len(scores))
for _, score := range scores {
items = append(items, GapItem{ControlID: score.ControlID, Gap: float64(score.Target - score.Current)})
}
sort.Slice(items, func(i, j int) bool { return items[i].Gap > items[j].Gap })
if n > len(items) {
n = len(items)
}
if n < 0 {
n = 0
}
return items[:n]
}
