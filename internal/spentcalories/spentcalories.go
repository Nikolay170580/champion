package spentcalories

import (
	"errors"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"time"
)

// Основные константы.
const (
	lenStep                    = 0.65 // средняя длина шага (fallback при неизвестном росте)
	mInKm                      = 1000 // метров в километре
	minInH                     = 60   // минут в часе
	stepLengthCoefficient      = 0.45 // коэффициент длины шага от роста
	walkingCaloriesCoefficient = 0.5  // коэффициент калорий при ходьбе
)

// parseDurationFlexible поддерживает "1.5h", "30.5m" и стандартные форматы.
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

	return 0, fmt.Errorf("invalid duration format %q", s)
}

// parseTraining парсит "шаги,Тип,длительность".
func parseTraining(data string) (int, string, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 3 {
		return 0, "", 0, fmt.Errorf("incorrect format: 3 fields expected, received %d", len(parts))
	}

	stepsStr := strings.TrimSpace(parts[0])
	activity := strings.TrimSpace(parts[1])
	durationStr := strings.TrimSpace(parts[2])

	if strings.ContainsAny(stepsStr+activity+durationStr, "\t\n\r") {
		return 0, "", 0, fmt.Errorf("invalid format: control characters are not allowed")
	}

	if stepsStr == "+" || stepsStr == "-" {
		return 0, "", 0, errors.New("parsing steps error: invalid syntax")
	}

	steps, err := strconv.Atoi(stepsStr)
	if err != nil {
		return 0, "", 0, fmt.Errorf("parsing steps error: %w", err)
	}
	if steps <= 0 {
		return 0, "", 0, fmt.Errorf("the number of steps must be > 0, received: %d", steps)
	}

	duration, err := parseDurationFlexible(durationStr)
	if err != nil {
		return 0, "", 0, fmt.Errorf("duration parsing error %q: %w", durationStr, err)
	}
	if duration <= 0 {
		return 0, "", 0, fmt.Errorf("duration must be > 0, received: %v", duration)
	}

	return steps, activity, duration, nil
}

// distance вычисляет дистанцию в км.
// Использует: шаги × (рост × 0.45), но если рост ≤ 0 — fallback на 0.65 м.
func distance(steps int, height float64) float64 {
	var stepLength float64
	if height > 0 {
		stepLength = height * stepLengthCoefficient
	} else {
		stepLength = lenStep
	}
	return float64(steps) * stepLength / float64(mInKm)
}

// meanSpeed вычисляет скорость в км/ч.
func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0.0
	}
	dist := distance(steps, height)
	hours := duration.Hours()
	if hours == 0 {
		return 0.0
	}
	return dist / hours
}

// TrainingInfo возвращает отформатированную строку с результатами.
func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		log.Println("Parsing error:", err)
		return "", err
	}

	var calories float64
	var calErr error

	switch activity {
	case "Бег":
		calories, calErr = RunningSpentCalories(steps, weight, height, duration)
	case "Ходьба":
		calories, calErr = WalkingSpentCalories(steps, weight, height, duration)
	default:
		return "", fmt.Errorf("неизвестный тип тренировки: %s", activity)
	}

	if calErr != nil {
		log.Println("Calorie calculation error:", calErr)
		return "", calErr
	}

	dist := distance(steps, height)
	speed := meanSpeed(steps, height, duration)
	durationHours := duration.Hours()

	// Отсечение до 2 знаков (не округление!) — важно для 590.625 → 590.62
	dist = math.Round(dist*100) / 100
	speed = math.Round(speed*100) / 100
	calories = math.Round(calories*100) / 100

	return fmt.Sprintf(
		"Тип тренировки: %s\n"+
			"Длительность: %.2f ч.\n"+
			"Дистанция: %.2f км.\n"+
			"Скорость: %.2f км/ч\n"+
			"Сожгли калорий: %.2f\n",
		activity,
		durationHours,
		dist,
		speed,
		calories,
	), nil
}

// RunningSpentCalories — бег. Требует height > 0.
func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, errors.New("Invalid parameters: steps, weight, height, and duration must be greater than 0")
	}

	speed := meanSpeed(steps, height, duration)
	durationInMinutes := duration.Minutes()
	return (weight * speed * durationInMinutes) / float64(minInH), nil
}

// WalkingSpentCalories — ходьба. Требует height > 0.
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, errors.New("Incorrect parameters: steps, weight, height, and duration must be greater than 0")
	}

	speed := meanSpeed(steps, height, duration)
	durationInMinutes := duration.Minutes()
	calories := (weight * speed * durationInMinutes) / float64(minInH)
	return calories * walkingCaloriesCoefficient, nil
}
