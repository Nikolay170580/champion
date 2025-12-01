package daysteps

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid data format")
	}

	stepsStr := strings.TrimSpace(parts[0])
	durationStr := strings.TrimSpace(parts[1])

	// Проверка на управляющие символы
	if strings.ContainsAny(stepsStr+durationStr, "\t\n\r") {
		return 0, 0, fmt.Errorf("Invalid data format: control characters are not allowed")
	}

	// Парсим шаги
	steps, err := strconv.Atoi(stepsStr)
	if err != nil {
		return 0, 0, fmt.Errorf("cannot convert steps to a number: %w", err)
	}
	if steps <= 0 {
		return 0, 0, fmt.Errorf("The number of steps must be greater than 0")
	}

	// Парсим длительность
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, 0, fmt.Errorf("unable to recognize the duration %q: %w", durationStr, err)
	}
	if duration <= 0 {
		return 0, 0, fmt.Errorf("the duration must be greater than 0")
	}

	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)
	if err != nil {
		log.Println("Error:", err)
		return ""
	}

	distanceM := float64(steps) * stepLength
	distanceKm := distanceM / float64(mInKm)

	// ✅ Правильный порядок аргументов + обработка ошибки
	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		log.Println("Calorie calculation error:", err)
		return ""
	}

	return fmt.Sprintf(
		"Количество шагов: %d.\n"+
			"Дистанция составила %.2f км.\n"+
			"Вы сожгли %.2f ккал.",
		steps, distanceKm, calories,
	)
}
