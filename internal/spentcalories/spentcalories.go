package spentcalories

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

// parseTraining парсит строку вида "3456,Ходьба,3h00m".
func parseTraining(data string) (int, string, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 3 {
		return 0, "", 0, fmt.Errorf("incorrect format: 3 fields expected, received %d", len(parts))
	}

	stepsStr := strings.TrimSpace(parts[0])
	activity := strings.TrimSpace(parts[1])
	durationStr := strings.TrimSpace(parts[2])

	// Проверка на запрещённые символы (табуляция, переводы строк)
	if strings.ContainsAny(stepsStr+activity+durationStr, "\t\n\r") {
		return 0, "", 0, fmt.Errorf("invalid format: control characters are not allowed")
	}

	// Парсим шаги
	steps, err := strconv.Atoi(stepsStr)
	if err != nil {
		return 0, "", 0, fmt.Errorf("parsing steps error: %w", err)
	}
	if steps <= 0 {
		return 0, "", 0, fmt.Errorf("the number of steps must be > 0, received: %d", steps)
	}

	// Парсим длительность
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, "", 0, fmt.Errorf("duration parsing error %q: %w", durationStr, err)
	}
	if duration <= 0 {
		return 0, "", 0, fmt.Errorf("duration must be > 0, received: %v", duration)
	}

	return steps, activity, duration, nil
}

// distance вычисляет дистанцию в километрах.
func distance(steps int, height float64) float64 {
	stepLength := height * stepLengthCoefficient
	distanceM := float64(steps) * stepLength
	return distanceM / float64(mInKm)
}

// meanSpeed вычисляет среднюю скорость в км/ч.
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

// TrainingInfo возвращает информацию о тренировке.
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

	result := fmt.Sprintf(
		"Тип тренировки: %s\n"+
			"Длительность: %.2f ч.\n"+
			"Дистанция: %.2f км.\n"+
			"Скорость: %.2f км/ч\n"+
			"Сожгли калорий: %.2f",
		activity,
		durationHours,
		dist,
		speed,
		calories,
	)

	return result, nil
}

// RunningSpentCalories рассчитывает калории для бега.
func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, errors.New("Invalid parameters: steps, weight, height, and duration must be greater than 0")
	}

	speed := meanSpeed(steps, height, duration)
	durationInMinutes := duration.Minutes()
	calories := (weight * speed * durationInMinutes) / float64(minInH)
	return calories, nil
}

// WalkingSpentCalories рассчитывает калории для ходьбы.
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, errors.New("Incorrect parameters: steps, weight, height, and duration must be greater than 0")
	}

	speed := meanSpeed(steps, height, duration)
	durationInMinutes := duration.Minutes()
	calories := (weight * speed * durationInMinutes) / float64(minInH)
	calories *= walkingCaloriesCoefficient
	return calories, nil
}
