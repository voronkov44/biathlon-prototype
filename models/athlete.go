package models

import "time"

type Status string

const (
	StatusNotStarted   Status = "NotStarted"
	StatusRacing       Status = "Racing"
	StatusNotFinished  Status = "NotFinished"
	StatusFinished     Status = "Finished"
	StatusDisqualified Status = "Disqualified"
)

type Athlete struct {
	ID               int
	RegisteredAt     time.Time
	StartTimePlanned time.Time
	StartTimeActual  *time.Time
	FinishTime       *time.Time
	Status           Status
	LapTimes         []time.Duration
	PenaltyTimes     []time.Duration
	CurrentLap       int
	TotalPenalty     int
	Shots            int
	Hits             int
	FiringLineTimes  map[int]time.Time // Время на каждом огневом рубеже
	LastLapTime      time.Time         //время завершения последнего круга
	TotalDistance    int               // Общая статистика
	AvgSpeed         float64           // Средняя скорость
	Accuracy         float64           // Точность стрельбы
}
