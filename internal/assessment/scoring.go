package assessment

import "math"

func Summarize(scores []ControlScore) Summary {
	if len(scores) == 0 {
		return Summary{}
	}
	var totalCurrent, totalTarget float64
	for _, score := range scores {
		totalCurrent += float64(score.Current)
		totalTarget += float64(score.Target)
	}
	count := float64(len(scores))
	avgCurrent := totalCurrent / count
	avgTarget := totalTarget / count
	return Summary{
		AverageCurrent: round(avgCurrent),
		AverageTarget:  round(avgTarget),
		Gap:            round(avgTarget - avgCurrent),
	}
}

func round(v float64) float64 {
	return math.Round(v*100) / 100
}
