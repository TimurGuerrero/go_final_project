package main

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func contains(slice []int, item int) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}

func NextDate(now time.Time, date string, repeat string) (string, error) {
	// Парсинг исходной даты
	startDate, err := time.Parse("20060102", date)
	if err != nil {
		return "", errors.New("неверный формат даты")
	}

	// Проверка на пустое правило
	if repeat == "" {
		return "", errors.New("пустое правило повторения")
	}

	// Обработка правил
	var nextDate time.Time

	switch {
	case strings.HasPrefix(repeat, "d "):
		// Правило d <число>
		daysStr := strings.TrimPrefix(repeat, "d ")
		days, err := strconv.Atoi(daysStr)
		if err != nil || days < 1 || days > 400 {
			return "", errors.New("invalid number of days")
		}

		// Добавляем дни до тех пор, пока не найдем подходящую дату
		nextDate = startDate
		for {
			nextDate = nextDate.AddDate(0, 0, days) // Добавляем дни
			if nextDate.After(now) {                // Проверяем, что следующая дата больше текущей
				return nextDate.Format("20060102"), nil
			}

		}

	case repeat == "y":
		// Правило y
		nextDate := startDate.AddDate(1, 0, 0)

		if nextDate.After(now) {
			return nextDate.Format("20060102"), nil
		}
		// Если следующая дата не превышает текущую, добавляем еще один год
		for {
			nextDate = nextDate.AddDate(1, 0, 0)
			if nextDate.After(now) {
				return nextDate.Format("20060102"), nil
			}
		}

	case strings.HasPrefix(repeat, "w "):
		// Правило w <числа>
		parts := strings.Split(repeat, " ")
		if len(parts) != 2 {
			return "", errors.New("неверный формат правила w")
		}
		daysOfWeek := strings.Split(parts[1], ",")
		weekDays := make([]int, len(daysOfWeek))
		for i, day := range daysOfWeek {
			dayInt, err := strconv.Atoi(day)
			if err != nil || dayInt < 1 || dayInt > 7 {
				return "", errors.New("неверное значение дня недели")
			}
			weekDays[i] = dayInt
		}
		nextDate := startDate
		for {
			nextDate = nextDate.AddDate(0, 0, 1)
			for _, day := range weekDays {
				if nextDate.Weekday() == time.Weekday(day%7) {
					if nextDate.After(now) {
						return nextDate.Format("20060102"), nil
					}
				}
			}
		}

	case strings.HasPrefix(repeat, "m "):
		// Правило m <числа>
		parts := strings.Split(repeat, " ")
		if len(parts) < 2 {
			return "", errors.New("неверный формат правила m")
		}
		days := strings.Split(parts[1], ",")
		months := []int{}
		if len(parts) == 3 {
			monthsStr := strings.Split(parts[2], ",")
			for _, month := range monthsStr {
				monthInt, err := strconv.Atoi(month)
				if err != nil || monthInt < 1 || monthInt > 12 {
					return "", errors.New("неверный месяц")
				}
				months = append(months, monthInt)
			}
		}
		nextDate := startDate
		for {
			nextDate = nextDate.AddDate(0, 0, 1)
			for _, day := range days {
				dayInt, err := strconv.Atoi(day)
				if err != nil || dayInt < -2 || dayInt > 31 {
					return "", errors.New("неверный день месяца")
				}
				if dayInt == -1 && nextDate.Day() == lastDayOfMonth(nextDate) {
					if nextDate.After(now) {
						return nextDate.Format("20060102"), nil
					}
				} else if dayInt == -2 && nextDate.Day() == lastDayOfMonth(nextDate)-1 {
					if nextDate.After(now) {
						return nextDate.Format("20060102"), nil
					}
				} else if nextDate.Day() == dayInt {
					if len(months) == 0 || contains(months, int(nextDate.Month())) {
						if nextDate.After(now) {
							return nextDate.Format("20060102"), nil
						}
					}
				}
			}
		}

	default:
		return "", errors.New("неверный формат правила повторения")
	}
}

func lastDayOfMonth(t time.Time) int {
	return time.Date(t.Year(), t.Month()+1, 0, 0, 0, 0, 0, t.Location()).Day()
}

func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	nowStr := r.FormValue("now")
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	now, err := time.Parse("20060102", nowStr)
	if err != nil {
		http.Error(w, "неверный формат даты 'now'", http.StatusBadRequest)
		return
	}

	nextDate, err := NextDate(now, date, repeat)
	if err != nil {
		http.Error(w, "неверный формат данных", http.StatusBadRequest)
		return
	}

	w.Write([]byte(nextDate))
}
