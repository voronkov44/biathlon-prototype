package main

import (
	"biathlon-prototype/config"
	"biathlon-prototype/events"
	"biathlon-prototype/race"
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// Создаем папку для логов, если ее нет
	if err := os.MkdirAll("logs", 0755); err != nil {
		log.Fatalf("Ошибка создания папки для логов: %v", err)
	}

	// Настройка логгера для событий
	eventsLogFile, err := os.Create(filepath.Join("logs", "events.log"))
	if err != nil {
		log.Fatalf("Ошибка создания файла логов событий: %v", err)
	}
	defer eventsLogFile.Close()
	eventsLogger := log.New(eventsLogFile, "", 0) // Без временных меток для чистого вывода

	// Настройка логгера для ошибок
	errorLogFile, err := os.Create(filepath.Join("logs", "errors.log"))
	if err != nil {
		log.Fatalf("Ошибка создания файла логов ошибок: %v", err)
	}
	defer errorLogFile.Close()
	errorLogger := log.New(errorLogFile, "", log.LstdFlags|log.Lshortfile)

	// Загрузка конфигурации
	cfg, err := config.LoadConfig("input_files/config.json")
	if err != nil {
		errorLogger.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	// Создание гонки
	r, err := race.NewRace(cfg)
	if err != nil {
		errorLogger.Fatalf("Ошибка создания гонки: %v", err)
	}

	// Обработка событий
	file, err := os.Open("input_files/events.txt")
	if err != nil {
		errorLogger.Fatalf("Ошибка открытия файла событий: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue // Пропускаем пустые строки
		}

		event, err := events.ParseEvent(line)
		if err != nil {
			errorLogger.Printf("Строка %d: %v (содержимое: %q)", lineNumber, err, line)
			continue
		}
		r.HandleEvent(event)
	}

	for _, event := range r.EventLog {
		eventsLogger.Println(event)
	}

	fmt.Println("\n=== РЕЗУЛЬТАТЫ ГОНКИ ===")
	r.PrintResults()

	// Информация о логах
	fmt.Printf("\nЛоги сохранены в папке logs:\n")
	fmt.Printf("- События: logs/events.log\n")
	if stat, err := errorLogFile.Stat(); err == nil && stat.Size() > 0 {
		fmt.Printf("- Ошибки: logs/errors.log\n")
	}
}
