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
	lenStep                    = 0.65 // средняя длина шага (на случай, если рост неизвестен — но в коде не используется напрямую)
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

// parseDurationFlexible расширяет time.ParseDuration, поддерживая дробные значения вида "1.5h", "30.5m"
func parseDurationFlexible(s string) (time.Duration, error) {
	// Сначала пробуем стандартный парсер (для "1h30m", "2h", "45m" и т.п.)
	if d, err := time.ParseDuration(s); err == nil {
		return d, nil
	}

	// Поддержка дробного формата: "X.Yh", "X.Ym"
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

	// Если ничего не подошло — ошибка
	return 0, fmt.Errorf("invalid duration format %q", s)
}

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
	if stepsStr == "+" || stepsStr == "-" {
		return 0, "", 0, fmt.Errorf("parsing steps error: %w", strconv.ErrSyntax)
	}
	steps, err := strconv.Atoi(stepsStr)
	if err != nil {
		return 0, "", 0, fmt.Errorf("parsing steps error: %w", err)
	}
	if steps <= 0 {
		return 0, "", 0, fmt.Errorf("the number of steps must be > 0, received: %d", steps)
	}

	// Парсим длительность — используем гибкий парсер
	duration, err := parseDurationFlexible(durationStr)
	if err != nil {
		return 0, "", 0, fmt.Errorf("duration parsing error %q: %w", durationStr, err)
	}
	if duration <= 0 {
		return 0, "", 0, fmt.Errorf("duration must be > 0, received: %v", duration)
	}

	return steps, activity, duration, nil
}

// distance вычисляет дистанцию в километрах.
// Используется формула: шаги × (рост × 0.45) / 1000
func distance(steps int, height float64) float64 {
	if height <= 0 {
		// Если рост недопустим, используем среднюю длину шага 0.65 м (lenStep), как fallback
		return float64(steps) * lenStep / float64(mInKm)
	}
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

	// Округляем до 2 знаков как в тестах
	dist = float64(int(dist*100+0.5)) / 100
	speed = float64(int(speed*100+0.5)) / 100
	calories = float64(int(calories*100+0.5)) / 100

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
// Формула: (вес × скорость × минуты) / 60
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
// Формула: (вес × скорость × минуты) / 60 × 0.5
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
