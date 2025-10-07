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
	stepLength = 0.65
	mInKm      = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	// проверка на пробелы в строке
	if strings.ContainsAny(data, " \t\n\r") {
		return 0, 0, fmt.Errorf("Пробелы запрещены")
	}

	parts := strings.Split(data, ",")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("формат данных неверный: ожидалось 'число,время'")
	}

	stepsStr := parts[0]
	durationStr := parts[1]

	// Проверка на минус
	if strings.Contains(stepsStr, "-") {
		return 0, 0, fmt.Errorf("отрицательные числа запрещены")
	}

	// Убираем плюс если есть
	if strings.HasPrefix(stepsStr, "+") {
		stepsStr = stepsStr[1:]
		// Проверяем что не осталась пустая строка
		if stepsStr == "" {
			return 0, 0, fmt.Errorf("неверный формат числа")
		}
	}

	// Проверка что остались только цифры
	for _, char := range stepsStr {
		if char < '0' || char > '9' {
			return 0, 0, fmt.Errorf("неверный формат числа")
		}
	}

	steps, err := strconv.Atoi(stepsStr)
	if err != nil {
		return 0, 0, fmt.Errorf("ошибка парсинга числа: %v", err)
	}

	if steps <= 0 {
		return 0, 0, fmt.Errorf("количество шагов должно быть больше 0")
	}

	// Проверка времени
	if strings.Contains(durationStr, "-") {
		return 0, 0, fmt.Errorf("отрицательное время запрещено")
	}

	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, 0, fmt.Errorf("ошибка парсинга времени: %v", err)
	}

	if duration <= 0 {
		return 0, 0, fmt.Errorf("время должно быть больше 0")
	}

	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)
	if err != nil {
		log.Println(err)
		return ""
	}

	if steps <= 0 || height <= 0 {
		return ""
	}

	distanceKm := (float64(steps) * stepLength) / mInKm

	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)

	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n",
		steps, distanceKm, calories)
}
