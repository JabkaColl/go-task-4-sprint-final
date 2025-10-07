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

	if strings.ContainsAny(data, " \t\n\r") {
		return 0, "", 0, fmt.Errorf("Пробелы запрещены")
	}

	partsTwo := strings.Split(data, ",")
	if len(partsTwo) != 3 {
		return 0, "", 0, fmt.Errorf("формат данных неверный: ожидалось 'количество шагов, вид активности, продолжительность'")
	}

	stepsTwoStr := partsTwo[0]
	activityTwoStr := partsTwo[1]
	durationTwoStr := partsTwo[2]

	if strings.Contains(stepsTwoStr, "-") {
		return 0, "", 0, fmt.Errorf("отрицательные числа запрещены")
	}

	if strings.HasPrefix(stepsTwoStr, "+") {
		stepsTwoStr = stepsTwoStr[1:]
		if stepsTwoStr == "" {
			return 0, "", 0, fmt.Errorf("неверный формат числа")
		}
	}

	stepsTwo, err := strconv.Atoi(stepsTwoStr)
	if err != nil {
		return 0, "", 0, err
	}

	if stepsTwo <= 0 {
		return 0, "", 0, fmt.Errorf("количество шагов должно быть больше 0")
	}

	if strings.Contains(durationTwoStr, "-") {
		return 0, "", 0, fmt.Errorf("Отрицательное время запрещено")
	}

	durationTwo, err := time.ParseDuration(durationTwoStr)
	if err != nil {
		return 0, "", 0, fmt.Errorf("ошибка парсинга времени: %v", err)
	}

	if durationTwo <= 0 {
		return 0, "", 0, errors.New("продолжительность должна быть больше нуля")
	}

	return stepsTwo, activityTwoStr, durationTwo, nil
}

func distance(steps int, height float64) float64 {
	strideLength := height * stepLengthCoefficient
	interval := (float64(steps) * strideLength) / mInKm
	return interval
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	distanceInterval := distance(steps, height)
	speed := distanceInterval / duration.Hours()
	return speed
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		log.Println("Ошибка парсинга тренировки:", err)
		return "", err
	}

	if weight <= 0 || height <= 0 {
		return "", errors.New("рост и вес должны быть больше 0")
	}

	distanceInterval := distance(steps, height)
	meanSped := meanSpeed(steps, height, duration)
	durationH := duration.Hours()

	calories := 0.0

	switch activity {
	case "Бег":
		calories, err = RunningSpentCalories(steps, weight, height, duration)
	case "Ходьба":
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
	default:
		return "", errors.New("неизвестный тип тренировки")
	}

	if err != nil {

		return "", err
	}

	return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
		activity, durationH, distanceInterval, meanSped, calories), nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if weight <= 0 || height <= 0 {
		return 0, errors.New("рост и вес должны быть больше 0")
	}

	if steps <= 0 {
		return 0, errors.New("шаги должны быть положительным числом")
	}

	if duration <= 0 {
		return 0, errors.New("продолжительность должна быть больше нуля")
	}

	speed := meanSpeed(steps, height, duration)
	durationMinyts := duration.Minutes()
	calories := (weight * speed * durationMinyts) / minInH // * walkingCaloriesCoefficient
	return calories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if weight <= 0 || height <= 0 {
		return 0, errors.New("рост и вес должны быть больше 0")
	}

	if steps <= 0 {
		return 0, errors.New("шаги должны быть положительным числом")
	}

	if duration <= 0 {
		return 0, errors.New("продолжительность должна быть больше нуля")
	}

	speed := meanSpeed(steps, height, duration)
	durationMinyts := duration.Minutes()
	calories := ((weight * speed * durationMinyts) / minInH) * walkingCaloriesCoefficient
	return calories, nil
}
