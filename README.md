# Biathlon-Prototype

Система моделирования и анализа биатлонной гонки на основе событий.

## Установка и запуск проекта

### 1. Клонирование репозитория
```bash
git clone https://github.com/voronkov44/biathlon-prototype2.git
```

### 2. Переход в корневую директорию 
```bash
cd biathlon-prototype
```

### 3. Запуск проекта
```
go run main.go
```

### Сборка и запуск проекта через Docker
Для удобства был собран `dockerfile`

*Требуется установка [docker](https://www.docker.com/products/docker-desktop/), если не установлен, смотрите [зависимости](https://github.com/voronkov44/biathlon-prototype?tab=readme-ov-file#%D0%B7%D0%B0%D0%B2%D0%B8%D1%81%D0%B8%D0%BC%D0%BE%D1%81%D1%82%D0%B8)*

Сборка образа
```
docker build . -t biathlon:v1
```

Запуск контейнера
```
docker run -it --rm biathlon:v1
```

С монтированием файлов событий и логов
```
docker run -it --rm \ -v $(pwd)/input_files:/opt/input_files \ -v $(pwd)/logs:/opt/logs \ biathlon:v1
```

Управление контейнерами
```
# Просмотр запущенных контейнеров
docker ps

# Остановка контейнера
docker stop <container_id>

# Удаление образа
docker rmi biathlon:v1

# Очистка системы
docker system prune
```

## Тесты
Запуск всех тестов
```bash
go test -v ./...
```

Запуск определенных тестов

Пример:
```bash
go test -v ./race
```


## Структура проекта
```
biathlon-prototype/
├── configs/
│ └── config.go # Чтение конфигурации гонки
│  └── config_test.go # Тест файла config
├── events/
│ └── event.go # Парсинг событий гонки
│  └── event_test.go # Тест файла event
├── input_files/
│ ├── config.json # Параметры гонки
│ └── events.txt # Лог событий гонки
├── logs/
│ ├── errors.log # Лог ошибок
│ └── events.log # Лог событий
├── models/
│ └── athlete.go # Модель участника
├── race/
│ ├── race.go # Обработка событий гонки
│  └── race_test.go # Тест файла race
│ ├── results.go # Вывод результатов
│  └── results_test.go # Тест файла results
├── utils/
│ └── time.go # Утилиты для работы со временем
├── main.go # Точка входа
└── Dockerfile # Конфигурация Docker
```

## Конфигурация гонки (config.json)
```
{
    "laps": 2, // Количество кругов
    "lapLen": 3651, // Длина каждого основного круга
    "penaltyLen": 50, // Длина каждого штрафного круга
    "firingLines": 1, // Количество огневых рубежей на круг
    "start": "09:30:00", // Планируемое время старта время первого участника
    "startDelta": "00:00:30" // Планируемый интервал между стартами
}
```

## Формат событий (events.txt)
Каждое событие имеет формат:
```
[время] ID-события ID-участника [параметры]
```
Пример:
```
[09:45:05.000] 6 1 1
```

Коды событий:
```
const (
    EventRegister         = 1  // Регистрация участника
    EventStartTimeLottery = 2  // Жеребьевка времени старта
    EventAtStartLine      = 3  // Участник на стартовой линии
    EventStart            = 4  // Старт гонки
    EventAtFiringLine     = 5  // Огневой рубеж (параметр: номер рубежа)
    EventHitSuccessful    = 6  // Попадание (параметр: номер мишени)
    EventLeaveFiringLine  = 7  // Покинул огневой рубеж
    EventEnterPenalty     = 8  // Вход на штрафной круг
    EventLeavePenalty     = 9  // Выход из штрафа
    EventLapFinish        = 10 // Завершение круга
    EventCantContinue     = 11 // Не может продолжить
    EventDisqualified     = 32 // Дисквалификация
    EventFinished         = 33 // Финиш
    EventHitMissed        = 61 // Промах (параметр: номер мишени)
)
```

## Модель участника
```
type Athlete struct {
    ID               int            // Уникальный идентификатор
    RegisteredAt     time.Time      // Время регистрации
    StartTimePlanned time.Time      // Планируемое время старта
    StartTimeActual  *time.Time     // Фактическое время старта
    FinishTime       *time.Time     // Время финиша
    Status           Status         // Текущий статус
    LapTimes         []time.Duration// Время кругов
    PenaltyTimes     []time.Duration// Штрафное время
    CurrentLap       int            // Текущий круг
    TotalPenalty     int            // Общий штраф (сек)
    Shots            int            // Всего выстрелов
    Hits             int            // Успешные попадания
    FiringLineTimes  map[int]time.Time // Время на огневых рубежах
    LastLapTime      time.Time      // Время последнего круга
    TotalDistance    int            // Общая дистанция (м)
    AvgSpeed         float64        // Средняя скорость (м/с)
    Accuracy         float64        // Точность стрельбы (%)
}
```




## **Зависимости**

Установка пакета [Docker Engine](https://docs.docker.com/engine/install/)
