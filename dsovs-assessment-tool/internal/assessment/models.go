package assessment

type ControlScore struct {
ControlID string `json:"control_id"`
Current   int    `json:"current"`
Target    int    `json:"target"`
}

type Summary struct {
AverageCurrent float64 `json:"average_current"`
AverageTarget  float64 `json:"average_target"`
Gap            float64 `json:"gap"`
}
