package models

type SensorData struct {
	Serial    string  `json:"Serial"`
	Type      string  `json:"Type"`
	Timestamp string  `json:"Timestamp"`
	Reading1  float64 `json:"Reading1"`
	Reading2  float64 `json:"Reading2"`
}
