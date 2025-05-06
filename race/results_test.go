package race

import (
	"biathlon-prototype/configs"
	"biathlon-prototype/models"
	"bytes"
	"io"
	"os"
	"sort"
	"strings"
	"testing"
	"time"
)

func TestPrintResults(t *testing.T) {
	// Создаем тестовую гонку
	race := createTestRaceWithAthletes()

	// Перехватываем вывод
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create pipe: %v", err)
	}
	os.Stdout = w

	// Выводим результаты
	race.PrintResults()

	// Восстанавливаем stdout
	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	_, err = io.Copy(&buf, r)
	if err != nil {
		t.Fatalf("Error reading output: %v", err)
	}

	output := buf.String()

	// Проверяем ключевые элементы вывода
	expected := []string{
		"🏁 Итоговый отчет:",
		"1. Участник 1 - Finished",
		"Общее время:",
		"Круг 1:",
		"Общая дистанция: 8000 м",
		"2. Участник 2 - Finished",
		"3. Участник 3 - Disqualified",
	}

	for _, exp := range expected {
		if !strings.Contains(output, exp) {
			t.Errorf("Expected output to contain %q", exp)
		}
	}
}

func TestResultsSorting(t *testing.T) {
	r := createTestRaceWithAthletes()
	r.CalculateStats()

	results := r.getSortedResults()

	// Проверяем порядок результатов
	if len(results) != 3 {
		t.Fatalf("Expected 3 results, got %d", len(results))
	}

	// Первый должен быть участник 1 (финишировал первым)
	if results[0].ID != 1 {
		t.Errorf("Expected first place to be athlete 1, got %d", results[0].ID)
	}

	// Последний должен быть участник 3 (дисквалифицирован)
	if results[2].ID != 3 {
		t.Errorf("Expected last place to be athlete 3, got %d", results[2].ID)
	}
}

// Вспомогательная функция для создания тестовой гонки с участниками
func createTestRaceWithAthletes() *Race {
	cfg := configs.Config{
		Laps:        2,
		LapLen:      4000,
		PenaltyLen:  150,
		FiringLines: 2,
		Start:       "10:00:00",
		StartDelta:  "00:01:00",
	}
	r, _ := NewRace(cfg)

	// Участник 1 - успешно финишировал
	athlete1 := &models.Athlete{
		ID:              1,
		Status:          models.StatusFinished,
		StartTimeActual: timePtr(time.Date(0, 1, 1, 10, 0, 0, 0, time.UTC)),
		FinishTime:      timePtr(time.Date(0, 1, 1, 11, 30, 0, 0, time.UTC)),
		LapTimes:        []time.Duration{30 * time.Minute, 30 * time.Minute},
		PenaltyTimes:    []time.Duration{2 * time.Minute},
		Shots:           10,
		Hits:            8,
	}
	r.Athletes[1] = athlete1

	// Участник 2 - финишировал позже
	athlete2 := &models.Athlete{
		ID:              2,
		Status:          models.StatusFinished,
		StartTimeActual: timePtr(time.Date(0, 1, 1, 10, 1, 0, 0, time.UTC)),
		FinishTime:      timePtr(time.Date(0, 1, 1, 11, 35, 0, 0, time.UTC)),
		LapTimes:        []time.Duration{32 * time.Minute, 32 * time.Minute},
		PenaltyTimes:    []time.Duration{3 * time.Minute},
		Shots:           10,
		Hits:            7,
	}
	r.Athletes[2] = athlete2

	// Участник 3 - дисквалифицирован
	athlete3 := &models.Athlete{
		ID:     3,
		Status: models.StatusDisqualified,
		Shots:  5,
		Hits:   2,
	}
	r.Athletes[3] = athlete3

	return r
}

// Вспомогательная функция для получения отсортированных результатов
func (r *Race) getSortedResults() []*models.Athlete {
	var results []*models.Athlete
	for _, athlete := range r.Athletes {
		results = append(results, athlete)
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].Status == results[j].Status {
			if results[i].FinishTime != nil && results[j].FinishTime != nil {
				return results[i].FinishTime.Before(*results[j].FinishTime)
			}
			return results[i].ID < results[j].ID
		}
		return results[i].Status < results[j].Status
	})

	return results
}

func timePtr(t time.Time) *time.Time {
	return &t
}
