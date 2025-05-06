package race

import (
	"biathlon-prototype/configs"
	"biathlon-prototype/events"
	"biathlon-prototype/models"
	"testing"
	"time"
)

func TestNewRace(t *testing.T) {
	cfg := configs.Config{
		Laps:        3,
		LapLen:      4000,
		PenaltyLen:  150,
		FiringLines: 2,
		Start:       "10:00:00",
		StartDelta:  "00:01:00",
	}

	r, err := NewRace(cfg)
	if err != nil {
		t.Fatalf("NewRace() error = %v, want nil", err)
	}

	if r.Config.Laps != 3 {
		t.Errorf("Expected 3 laps, got %d", r.Config.Laps)
	}

	if r.StartTime.Hour() != 10 {
		t.Errorf("Expected start at 10:00, got %v", r.StartTime)
	}

	if r.StartDelta != time.Minute {
		t.Errorf("Expected start delta 1m, got %v", r.StartDelta)
	}
}

func TestHandleEvent_Registration(t *testing.T) {
	r := createTestRace()
	event := createTestEvent(events.EventRegister, "09:00:00.000", 1)

	r.HandleEvent(event)

	if len(r.Athletes) != 1 {
		t.Fatalf("Expected 1 athlete, got %d", len(r.Athletes))
	}

	athlete := r.Athletes[1]
	if athlete.ID != 1 {
		t.Errorf("Expected athlete ID 1, got %d", athlete.ID)
	}

	if !athlete.RegisteredAt.Equal(event.Time) {
		t.Errorf("Expected registration at %v, got %v", event.Time, athlete.RegisteredAt)
	}
}

func TestHandleEvent_Start(t *testing.T) {
	r := createTestRace()
	registerAthlete(r, 1)
	event := createTestEvent(events.EventStart, "10:00:00.000", 1)

	r.HandleEvent(event)

	athlete := r.Athletes[1]
	if athlete.Status != models.StatusRacing {
		t.Errorf("Expected status Racing, got %v", athlete.Status)
	}

	if athlete.StartTimeActual == nil {
		t.Error("Expected StartTimeActual to be set")
	} else if !athlete.StartTimeActual.Equal(event.Time) {
		t.Errorf("Expected start at %v, got %v", event.Time, athlete.StartTimeActual)
	}
}

func TestHandleEvent_FiringLine(t *testing.T) {
	r := createTestRace()
	registerAndStartAthlete(r, 1)
	event := createTestEvent(events.EventAtFiringLine, "10:30:00.000", 1, "1")

	r.HandleEvent(event)

	athlete := r.Athletes[1]
	if len(athlete.FiringLineTimes) != 1 {
		t.Errorf("Expected 1 firing line time, got %d", len(athlete.FiringLineTimes))
	}

	if line := r.CurrentFiring[1]; line != 1 {
		t.Errorf("Expected current firing line 1, got %d", line)
	}
}

func TestHandleEvent_Shooting(t *testing.T) {
	r := createTestRace()
	registerAndStartAthlete(r, 1)
	firingLineEvent := createTestEvent(events.EventAtFiringLine, "10:30:00.000", 1, "1")
	r.HandleEvent(firingLineEvent)

	// Test successful hit
	hitEvent := createTestEvent(events.EventHitSuccessful, "10:30:05.000", 1, "1")
	r.HandleEvent(hitEvent)

	athlete := r.Athletes[1]
	if athlete.Hits != 1 {
		t.Errorf("Expected 1 hit, got %d", athlete.Hits)
	}

	// Test missed shot
	missEvent := createTestEvent(events.EventHitMissed, "10:30:06.000", 1, "2")
	r.HandleEvent(missEvent)

	if athlete.Shots != 2 {
		t.Errorf("Expected 2 shots, got %d", athlete.Shots)
	}

	if len(athlete.PenaltyTimes) != 1 {
		t.Errorf("Expected 1 penalty, got %d", len(athlete.PenaltyTimes))
	}
}

func TestHandleEvent_LapFinish(t *testing.T) {
	r := createTestRace()
	registerAndStartAthlete(r, 1)
	startTime := *r.Athletes[1].StartTimeActual

	lapEvent := createTestEvent(events.EventLapFinish, "10:30:00.000", 1)
	r.HandleEvent(lapEvent)

	athlete := r.Athletes[1]
	if athlete.CurrentLap != 1 {
		t.Errorf("Expected current lap 1, got %d", athlete.CurrentLap)
	}

	if len(athlete.LapTimes) != 1 {
		t.Errorf("Expected 1 lap time, got %d", len(athlete.LapTimes))
	}

	expectedDuration := lapEvent.Time.Sub(startTime)
	if athlete.LapTimes[0] != expectedDuration {
		t.Errorf("Expected lap time %v, got %v", expectedDuration, athlete.LapTimes[0])
	}
}

func TestCalculateStats(t *testing.T) {
	r := createTestRace()
	registerAndStartAthlete(r, 1)
	finishAthlete(r, 1)

	r.CalculateStats()

	athlete := r.Athletes[1]
	if athlete.TotalDistance != r.Config.LapLen*len(athlete.LapTimes) {
		t.Errorf("Incorrect total distance calculation")
	}

	if athlete.AvgSpeed <= 0 {
		t.Error("Expected avg speed > 0")
	}

	if athlete.Accuracy < 0 || athlete.Accuracy > 100 {
		t.Errorf("Invalid accuracy value: %.2f", athlete.Accuracy)
	}
}

// Helper functions
func createTestRace() *Race {
	cfg := configs.Config{
		Laps:        3,
		LapLen:      4000,
		PenaltyLen:  150,
		FiringLines: 2,
		Start:       "10:00:00",
		StartDelta:  "00:01:00",
	}
	r, _ := NewRace(cfg)
	return r
}

func createTestEvent(eventID int, timeStr string, athleteID int, params ...string) events.Event {
	t, _ := time.Parse("15:04:05.000", timeStr)
	return events.Event{
		Time:      t,
		EventID:   eventID,
		AthleteID: athleteID,
		Params:    params,
	}
}

func registerAthlete(r *Race, id int) {
	event := createTestEvent(events.EventRegister, "09:00:00.000", id)
	r.HandleEvent(event)
}

func registerAndStartAthlete(r *Race, id int) {
	registerAthlete(r, id)
	lotteryEvent := createTestEvent(events.EventStartTimeLottery, "09:05:00.000", id, "10:00:00.000")
	r.HandleEvent(lotteryEvent)
	startEvent := createTestEvent(events.EventStart, "10:00:00.000", id)
	r.HandleEvent(startEvent)
}

func finishAthlete(r *Race, id int) {
	// Simulate full race
	lap1 := createTestEvent(events.EventLapFinish, "10:30:00.000", id)
	r.HandleEvent(lap1)
	lap2 := createTestEvent(events.EventLapFinish, "11:00:00.000", id)
	r.HandleEvent(lap2)
	finish := createTestEvent(events.EventFinished, "11:30:00.000", id)
	r.HandleEvent(finish)
}
