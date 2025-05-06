package race

import (
	"biathlon-prototype/configs"
	"biathlon-prototype/events"
	"biathlon-prototype/models"
	"biathlon-prototype/utils"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Race struct {
	Config        configs.Config
	StartTime     time.Time
	StartDelta    time.Duration
	Athletes      map[int]*models.Athlete
	EventLog      []string
	CurrentFiring map[int]int
}

func NewRace(cfg configs.Config) (*Race, error) {
	startTime, err := time.Parse("15:04:05", cfg.Start)
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга времени старта: %v", err)
	}

	startDeltaStr := cfg.StartDelta
	if strings.Contains(startDeltaStr, ":") {
		parts := strings.Split(startDeltaStr, ":")
		if len(parts) == 3 {
			h, _ := strconv.Atoi(parts[0])
			m, _ := strconv.Atoi(parts[1])
			s, _ := strconv.Atoi(parts[2])
			startDeltaStr = fmt.Sprintf("%dh%dm%ds", h, m, s)
		}
	}

	startDelta, err := time.ParseDuration(startDeltaStr)
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга стартового интервала: %v", err)
	}

	return &Race{
		Config:        cfg,
		StartTime:     startTime,
		StartDelta:    startDelta,
		Athletes:      make(map[int]*models.Athlete),
		EventLog:      make([]string, 0),
		CurrentFiring: make(map[int]int),
	}, nil
}

func (r *Race) logEvent(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	r.EventLog = append(r.EventLog, msg)
	fmt.Println(msg)
}

func (r *Race) HandleEvent(event events.Event) {
	athlete, exists := r.Athletes[event.AthleteID]
	if !exists {
		athlete = &models.Athlete{
			ID:              event.AthleteID,
			Status:          models.StatusNotStarted,
			FiringLineTimes: make(map[int]time.Time),
			LapTimes:        make([]time.Duration, 0),
			PenaltyTimes:    make([]time.Duration, 0),
		}
		r.Athletes[event.AthleteID] = athlete
	}

	switch event.EventID {
	case events.EventRegister:
		athlete.RegisteredAt = event.Time
		athlete.Status = models.StatusNotStarted
		r.logEvent("[%s] Участник(%d) зарегистрирован",
			utils.FormatTime(event.Time), athlete.ID)

	case events.EventStartTimeLottery:
		if len(event.Params) > 0 {
			startTime, err := time.Parse("15:04:05.000", event.Params[0])
			if err == nil {
				athlete.StartTimePlanned = startTime
				r.logEvent("[%s] Время старта для участника(%d) установлено жеребьевкой на %s",
					utils.FormatTime(event.Time), athlete.ID, event.Params[0])
			}
		}

	case events.EventAtStartLine:
		athlete.Status = models.StatusRacing
		r.logEvent("[%s] Участник(%d) на стартовой линии",
			utils.FormatTime(event.Time), athlete.ID)

	case events.EventStart:
		now := event.Time
		athlete.StartTimeActual = &now
		athlete.Status = models.StatusRacing
		r.logEvent("[%s] Участник(%d) начал гонку",
			utils.FormatTime(event.Time), athlete.ID)

	case events.EventAtFiringLine:
		if len(event.Params) > 0 {
			firingLine, err := strconv.Atoi(event.Params[0])
			if err == nil {
				athlete.FiringLineTimes[firingLine] = event.Time
				r.CurrentFiring[athlete.ID] = firingLine
				r.logEvent("[%s] Участник(%d) на огневом рубеже(%d)",
					utils.FormatTime(event.Time), athlete.ID, firingLine)
			}
		}

	case events.EventHitSuccessful:
		if len(event.Params) > 0 {
			athlete.Hits++
			athlete.Shots++
			r.logEvent("[%s] Участник(%d) попал в мишень %s",
				utils.FormatTime(event.Time), athlete.ID, event.Params[0])
		}

	case events.EventHitMissed:
		if len(event.Params) > 0 {
			athlete.Shots++ // Только счетчик выстрелов
			r.logEvent("[%s] Участник(%d) промахнулся по мишени %s",
				utils.FormatTime(event.Time), athlete.ID, event.Params[0])
			// Автоматический штрафной круг за каждый промах
			athlete.PenaltyTimes = append(athlete.PenaltyTimes, 0)
		}

	case events.EventLeaveFiringLine:
		firingLine := r.CurrentFiring[athlete.ID]
		if startTime, exists := athlete.FiringLineTimes[firingLine]; exists {
			timeSpent := event.Time.Sub(startTime)
			r.logEvent("[%s] Участник(%d) покинул огневой рубеж(%d) (время: %v)",
				utils.FormatTime(event.Time), athlete.ID, firingLine, timeSpent)
		}

	case events.EventEnterPenalty:
		athlete.PenaltyTimes = append(athlete.PenaltyTimes, 0)
		r.logEvent("[%s] Участник(%d) вошел на штрафные круги",
			utils.FormatTime(event.Time), athlete.ID)

	case events.EventLeavePenalty:
		if len(athlete.PenaltyTimes) > 0 {
			penaltyIdx := len(athlete.PenaltyTimes) - 1
			if firingLine, ok := r.CurrentFiring[athlete.ID]; ok {
				if startTime, exists := athlete.FiringLineTimes[firingLine]; exists {
					penaltyTime := event.Time.Sub(startTime)
					athlete.PenaltyTimes[penaltyIdx] = penaltyTime
					// Расчет общего штрафа
					athlete.TotalPenalty += int(penaltyTime.Seconds())
					r.logEvent("[%s] Участник(%d) покинул штрафные круги (время: %v, общий штраф: %d сек)",
						utils.FormatTime(event.Time), athlete.ID, penaltyTime, athlete.TotalPenalty)
				}
			}
		}

	case events.EventLapFinish:
		athlete.CurrentLap++
		if athlete.CurrentLap <= r.Config.Laps && athlete.StartTimeActual != nil {
			var lapTime time.Duration
			if len(athlete.LapTimes) > 0 {
				// Время текущего круга = текущее время - время завершения предыдущего круга
				lapTime = event.Time.Sub(athlete.LastLapTime)
			} else {
				// Для первого круга = текущее время - время старта
				lapTime = event.Time.Sub(*athlete.StartTimeActual)
			}
			athlete.LapTimes = append(athlete.LapTimes, lapTime)
			athlete.LastLapTime = event.Time // Добавляем новое поле для хранения времени последнего круга
			r.logEvent("[%s] Участник(%d) завершил круг %d (время круга: %v, общее время: %v)",
				utils.FormatTime(event.Time), athlete.ID, athlete.CurrentLap,
				lapTime, event.Time.Sub(*athlete.StartTimeActual))
		}

	case events.EventCantContinue:
		athlete.Status = models.StatusNotFinished
		reason := "без указания причины"
		if len(event.Params) > 0 {
			reason = event.Params[0]
		}
		r.logEvent("[%s] Участник(%d) не может продолжить: %s",
			utils.FormatTime(event.Time), athlete.ID, reason)

	case events.EventDisqualified:
		athlete.Status = models.StatusDisqualified
		r.logEvent("[%s] Участник(%d) дисквалифицирован",
			utils.FormatTime(event.Time), athlete.ID)

	case events.EventFinished:
		now := event.Time
		athlete.FinishTime = &now
		athlete.Status = models.StatusFinished
		r.logEvent("[%s] Участник(%d) финишировал",
			utils.FormatTime(event.Time), athlete.ID)
	}

	// Автоматическая дисквалификация за опоздание на старт
	if athlete.Status == models.StatusNotStarted &&
		event.Time.After(athlete.StartTimePlanned.Add(r.StartDelta)) {
		athlete.Status = models.StatusDisqualified
		r.logEvent("[%s] Участник(%d) дисквалифицирован (не стартовал вовремя)",
			utils.FormatTime(event.Time), athlete.ID)
	}
}

// CalculateStats вычисляет дополнительную статистику по участникам
func (r *Race) CalculateStats() {
	for _, a := range r.Athletes {
		// Общая дистанция (только завершенные круги)
		a.TotalDistance = r.Config.LapLen * len(a.LapTimes)

		// Средняя скорость (если гонка завершена)
		if a.StartTimeActual != nil && a.FinishTime != nil {
			totalTime := a.FinishTime.Sub(*a.StartTimeActual).Seconds()
			if totalTime > 0 {
				a.AvgSpeed = float64(a.TotalDistance) / totalTime
			}
		}

		// Точность стрельбы
		if a.Shots > 0 {
			a.Accuracy = float64(a.Hits) / float64(a.Shots) * 100
		}
	}
}
