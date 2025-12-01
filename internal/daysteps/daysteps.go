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

	stepsStr := parts[0]
	durationStr := parts[1]

	// Проверка на управляющие символы
	if strings.ContainsAny(stepsStr+durationStr, "\t\n\r") {
		return 0, 0, fmt.Errorf("Invalid data format: control characters are not allowed")
	}

	// Запрет пробелов в начале/конце (включая обычные пробелы)
	if strings.TrimSpace(stepsStr) != stepsStr || strings.TrimSpace(durationStr) != durationStr {
		return 0, 0, fmt.Errorf("The number of steps must be greater than 0")
	}

	// Проверка на одиночные '+' и '-'
	if stepsStr == "+" || stepsStr == "-" {
		return 0, 0, fmt.Errorf("The number of steps must be greater than 0")
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
	duration, err := parseDurationFlexible(durationStr)
	if err != nil {
		return 0, 0, fmt.Errorf("unable to recognize the duration %q: %w", durationStr, err)
	}
	if duration <= 0 {
		return 0, 0, fmt.Errorf("the duration must be greater than 0")
	}

	return steps, duration, nil
}

func parseDurationFlexible(s string) (time.Duration, error) {
	if d, err := time.ParseDuration(s); err == nil {
		return d, nil
	}

	if strings.HasSuffix(s, "h") {
		numStr := strings.TrimSuffix(s, "h")
		if v, err := strconv.ParseFloat(numStr, 64); err == nil && v > 0 {
			return time.Duration(v * float64(time.Hour)), nil
		}
	} else if strings.HasSuffix(s, "m") {
		numStr := strings.TrimSuffix(s, "m")
		if v, err := strconv.ParseFloat(numStr, 64); err == nil && v > 0 {
			return time.Duration(v * float64(time.Minute)), nil
		}
	}

	return 0, fmt.Errorf("invalid duration format")
}

func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)
	if err != nil {
		log.Println("Error:", err)
		return ""
	}

	distanceM := float64(steps) * stepLength
	distanceKm := distanceM / float64(mInKm)

	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		log.Println("Calorie calculation error:", err)
		return ""
	}

	return fmt.Sprintf(
		"Количество шагов: %d.\n"+
			"Дистанция составила %.2f км.\n"+
			"Вы сожгли %.2f ккал.\n",
		steps, distanceKm, calories,
	)
}
