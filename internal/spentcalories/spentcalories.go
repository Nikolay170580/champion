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

func parseTraining(data string) (int, string, time.Duration, error) {
	// TODO: реализовать функцию
	parts := strings.Split(data, ",")
	if len(parts) != 3 {
		return 0, "", 0, fmt.Errorf("некорректный формат: ожидалось 3 поля")
	}

	steps, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return 0, "", 0, fmt.Errorf("ошибка парсинга шагов: %w", err)
	}
	stepsStr := parts[0]
	activity := parts[1]
	durationStr := parts[2]

	if strings.ContainsAny(stepsStr,"\t\n\r") || strings.ContainsAny(activity,"\t\n\r") || strings.ContainsAny(durationStr,"\t\n\r") {
	return 0,"",0,fmt.Errorf("некоректный форматЖ пробелы недопустимы")
  	
	steps, err := strconv.Atoi(stepsStr)
	if err != nil {
		return 0,"",0,fmt.Errorf("ошибка парсинга шагов: %w",err)
	}
	if steps <= 0 {
		return 0,"",0,fmt.Errorf("количество шагов должно быть больше 0")
	}
	
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0,"",0,fmt.Errorf("ошибка парсинга длительности %q: %w", durationStr, err)
	}
	if duration <= 0 {
		return 0,"",0,fmt.Errorf("продолжительность должна быть больше 0")
	}

	return steps, activity, duration, nil
}

func distance(steps int, height float64) float64 {
	// TODO: реализовать функцию
	stepLength := height * stepLengthCoefficient
	distanceM := float64(steps) * stepLength
	return distanceM / float64(mInKm)
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	// TODO: реализовать функцию
	if duration <= 0 {
		return 0.0
	}
	dist := distance(steps, height)
	hours := duration.Hours()
	if hoyrs == 0 {
		return 0.0
	
	}

	return dist / hours
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	// TODO: реализовать функцию
	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		log.Println("Ошибка парсинга:", err)
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
		log.Println("Ошибка расчёта калорий:", calErr)
		return "", calErr
	}

	dist := distance(steps, height)
	speed := meanSpeed(steps, height, duration)
	durationHours := duration.Hours()

	// Форматируем строку в точности как в ТЗ (с точностью до 2 знаков)
	result := fmt.Sprintf(
		"Тип тренировки: %s\n"+
			"Длительность: %.2fч.\n"+
			"Дистанция: %.2fкм.\n"+
			"Скорость: %.2fкм/ч\n"+
			"Сожгли калорий: %.2f\n",
		activity,
		durationHours,
		dist,
		speed,
		calories,
	)

	return result, nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// TODO: реализовать функцию
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, errors.New("некорректные параметры")
	}
	
	speed := meanSpeed(steps, height, duration)
	durationInMinutes := duration.Minutes()

	calories := (weight * speed * durationInMinutes) / float64(minInH)
	return calories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// TODO: реализовать функцию
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, errors.New("некорректные параметры")
	}
	
	speed := meanSpeed(steps, height, duration)
	durationInMinutes := duration.Minutes()

	calories := (weight * speed * durationInMinutes) / float64(minInH)
	calories *= walkingCaloriesCoefficient
	return calories, nil
}
