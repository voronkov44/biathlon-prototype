package race

import (
	"biathlon-prototype/models"
	"biathlon-prototype/utils"
	"fmt"
	"math"
	"sort"
)

func (r *Race) PrintResults() {
	// Рассчитываем дополнительную статистику перед выводом
	r.CalculateStats()

	var results []*models.Athlete

	// Собираем всех участников
	for _, athlete := range r.Athletes {
		results = append(results, athlete)
	}

	// Сортируем по статусу и времени
	sort.Slice(results, func(i, j int) bool {
		if results[i].Status == results[j].Status {
			if results[i].FinishTime != nil && results[j].FinishTime != nil {
				return results[i].FinishTime.Before(*results[j].FinishTime)
			}
			return results[i].ID < results[j].ID
		}
		return results[i].Status < results[j].Status
	})

	fmt.Println("\n🏁 Итоговый отчет:")
	for pos, athlete := range results {
		fmt.Printf("%d. Участник %d - %s\n", pos+1, athlete.ID, athlete.Status)

		// Основная информация о времени
		if athlete.StartTimeActual != nil {
			if athlete.FinishTime != nil {
				totalTime := athlete.FinishTime.Sub(*athlete.StartTimeActual)
				fmt.Printf("   Общее время: %s\n", utils.FormatDuration(totalTime))
			}

			// Время кругов с расчетом скорости
			for i, lapTime := range athlete.LapTimes {
				if i < len(athlete.LapTimes) {
					speed := float64(r.Config.LapLen) / lapTime.Seconds()
					fmt.Printf("   Круг %d: %s (%.2f м/с)\n",
						i+1, utils.FormatDuration(lapTime), speed)
				}
			}

			// Штрафное время
			totalPenaltySeconds := 0
			for i, penaltyTime := range athlete.PenaltyTimes {
				if penaltyTime > 0 {
					speed := float64(r.Config.PenaltyLen) / penaltyTime.Seconds()
					totalPenaltySeconds += int(penaltyTime.Seconds())
					fmt.Printf("   Штраф %d: %s (%.2f м/с)\n",
						i+1, utils.FormatDuration(penaltyTime), math.Round(speed*100)/100)
				}
			}
			if totalPenaltySeconds > 0 {
				fmt.Printf("   Общее штрафное время: %d сек\n", totalPenaltySeconds)
			}
		}

		// Расширенная статистика
		fmt.Printf("   Общая дистанция: %d м\n", athlete.TotalDistance)
		if athlete.AvgSpeed > 0 {
			fmt.Printf("   Средняя скорость: %.2f м/с\n", athlete.AvgSpeed)
		}
		if athlete.Shots > 0 {
			fmt.Printf("   Точность стрельбы: %.1f%%\n", athlete.Accuracy)
		}
		fmt.Printf("   Стрельба: %d/%d попаданий\n\n",
			athlete.Hits, athlete.Shots)
	}
}
