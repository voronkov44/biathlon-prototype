package events

import (
	"biathlon-prototype/utils"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// константы идентификаторов событий
const (
	EventRegister         = 1
	EventStartTimeLottery = 2
	EventAtStartLine      = 3
	EventStart            = 4
	EventAtFiringLine     = 5
	EventHitSuccessful    = 6
	EventLeaveFiringLine  = 7
	EventEnterPenalty     = 8
	EventLeavePenalty     = 9
	EventLapFinish        = 10
	EventCantContinue     = 11

	EventDisqualified = 32
	EventFinished     = 33
	EventHitMissed    = 61
)

type Event struct {
	Time      time.Time
	EventID   int
	AthleteID int
	Params    []string
	Raw       string
}

// ParseEvent парсит строку события в структуру Event
func ParseEvent(line string) (Event, error) {
	var event Event
	event.Raw = line

	// Пропускаем пустые строки
	if strings.TrimSpace(line) == "" {
		return event, fmt.Errorf("пустая строка")
	}

	// Разбиваем строку на части
	parts := strings.Fields(line)
	if len(parts) < 3 {
		return event, fmt.Errorf("Некорректная строка события: %s", line)
	}

	// Парсим время события
	t, err := utils.ParseTime(strings.Trim(parts[0], "[]"))
	if err != nil {
		return event, fmt.Errorf("Ошибка парсинга времени: %v", err)
	}
	event.Time = t

	// парсим ID события
	eventID, err := strconv.Atoi(parts[1])
	if err != nil {
		return event, fmt.Errorf("Ошибка парсинга ID события: %v", err)
	}
	event.EventID = eventID

	// парсим ID участника
	athleteID, err := strconv.Atoi(parts[2])
	if err != nil {
		return event, fmt.Errorf("Ошибка парсинга ID участника: %v", err)
	}
	event.AthleteID = athleteID

	// Сохраняем остальные параметры, если есть
	if len(parts) > 3 {
		event.Params = parts[3:]
	}

	return event, nil
}
