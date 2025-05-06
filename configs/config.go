package configs

import (
	"encoding/json"
	"os"
)

type Config struct {
	Laps        int    `json:"laps"`        //Количество кругов основной дистанции
	LapLen      int    `json:"lapLen"`      //Длина каждого основного круга
	PenaltyLen  int    `json:"penaltyLen"`  //Длина каждого штрафного круга
	FiringLines int    `json:"firingLines"` //Количество огневых рубежей на круг
	Start       string `json:"start"`       //Планируемое время старта первого участника
	StartDelta  string `json:"startDelta"`  //Планируемый интервал между стартами
}

// LoadConfig читает конфигурационный JSON-файл и возвращает структуру Config
func LoadConfig(filename string) (Config, error) {
	var config Config

	data, err := os.ReadFile(filename)
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(data, &config)
	return config, err
}
