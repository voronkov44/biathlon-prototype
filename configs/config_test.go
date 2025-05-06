package configs

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Создаем временный конфигурационный файл
	const testConfig = `{
		"laps": 3,
		"lapLen": 4000,
		"penaltyLen": 150,
		"firingLines": 2,
		"start": "10:00:00",
		"startDelta": "00:01:00"
	}`

	tmpFile, err := os.CreateTemp("", "config_test_*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(testConfig); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	// Тест 1: Успешная загрузка конфигурации
	t.Run("Valid configs", func(t *testing.T) {
		cfg, err := LoadConfig(tmpFile.Name())
		if err != nil {
			t.Errorf("LoadConfig() error = %v, want nil", err)
			return
		}

		if cfg.Laps != 3 {
			t.Errorf("Laps = %d, want 3", cfg.Laps)
		}

		if cfg.LapLen != 4000 {
			t.Errorf("LapLen = %d, want 4000", cfg.LapLen)
		}

		if cfg.PenaltyLen != 150 {
			t.Errorf("PenaltyLen = %d, want 150", cfg.PenaltyLen)
		}

		if cfg.FiringLines != 2 {
			t.Errorf("FiringLines = %d, want 2", cfg.FiringLines)
		}

		if cfg.Start != "10:00:00" {
			t.Errorf("Start = %s, want 10:00:00", cfg.Start)
		}

		if cfg.StartDelta != "00:01:00" {
			t.Errorf("StartDelta = %s, want 00:01:00", cfg.StartDelta)
		}
	})

	// Тест 2: Файл не существует
	t.Run("File not exists", func(t *testing.T) {
		_, err := LoadConfig("nonexistent_file.json")
		if err == nil {
			t.Error("Expected error for nonexistent file, got nil")
		}
	})

	// Тест 3: Невалидный JSON
	t.Run("Invalid JSON", func(t *testing.T) {
		invalidFile, err := os.CreateTemp("", "config_invalid_*.json")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.Remove(invalidFile.Name())

		if _, err := invalidFile.WriteString("{ invalid json }"); err != nil {
			t.Fatalf("Failed to write to temp file: %v", err)
		}
		invalidFile.Close()

		_, err = LoadConfig(invalidFile.Name())
		if err == nil {
			t.Error("Expected error for invalid JSON, got nil")
		}
	})
}

func TestConfigValidation(t *testing.T) {
	t.Run("Config validation", func(t *testing.T) {
	})
}
