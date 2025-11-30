package daysteps

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	// TODO: реализовать функцию
	// Разделяем введенную строку на слайс строк по запятой
	parts := strings.Split(data, ",")
	// Проверяем длину слайса на равность 2
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("некорректный формат даннных")
	}
	// Преобразование первого элемента слайса в количество шагов
	steps, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return 0, 0, fmt.Errorf("невозможно преобразовать шаги в число")
	}
	// Проверить: количество шагов должно быть больше 0
	if steps <= 0 {
		return 0, 0, fmt.Errorf("количество шагов должно быть больше 0")
	}
	// Преобразовать второй элемент слайса в time.Duration
	durationStr := strings.TrimSpace(parts[1])
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, 0, fmt.Errorf("невозможно распознать продолжительность")
	}
	// Если всё прошло без ошибок, возвращаем количество шагов, продолжительность и nil
	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	// TODO: реализовать функцию
	// Получить данные с помощью parsePackage()
	steps, duration, err := parsePackage(data)
	if err != nil {
		fmt.Println("Ошибка:", err)
		return ""
	}

	// Проверить, чтобы количество шагов было больше 0
	if steps <= 0 {
		return ""
	}

	// Вычислить дистанцию в метрах: шаги × длина шага
	distanceM := float64(steps) * stepLength

	// Перевести дистанцию в километры
	distanceKm := distanceM / float64(mInKm)

	// Вычислить количество калорий
	calories := 
	WalkingSpentCalories(steps, duration, weight, height)

	// Формирование и возврат строки
	return fmt.Sprintf(
		"Количество шагов: %d.\n"+
			"Дистанция составила %.2f км.\n"+
			"Вы сожгли %.2f ккал.",
		steps, distanceKm, calories,
	)
}
