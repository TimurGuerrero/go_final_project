package main

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func сontains(slice []int, item int) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}

func NextDate(now time.Time, date string, repeat string) (string, error) {
	startDate, err := time.Parse("20060102", date)

	if err != nil {
		return "", errors.New("invalid date format")
	}

	// Функция для проверки наличия месяца в списке

	var nextDate time.Time

	switch {
	// Правило d <число>
	case strings.HasPrefix(repeat, "d "):
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
				break
			}
		}

		// Правило y
	case repeat == "y":
		nextDate = startDate.AddDate(1, 0, 0)

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
		// Правило w <через запятую от 1 до 7>
		daysOfWeekStr := strings.TrimPrefix(repeat, "w ")
		daysOfWeek := strings.Split(daysOfWeekStr, ",")
		daysMap := make(map[int]bool)

		for _, day := range daysOfWeek {
			d, err := strconv.Atoi(strings.TrimSpace(day)) // Удаляем пробелы
			if err == nil && d >= 1 && d <= 7 {
				daysMap[d] = true
			} else {
				return "", errors.New("invalid day of the week: " + day)
			}
		}

		if len(daysMap) == 0 {
			return "", errors.New("no valid days of the week specified")
		}

		nextDate = startDate
		for {
			nextDate = nextDate.AddDate(0, 0, 1)    // Переходим к следующему дню
			if daysMap[int(nextDate.Weekday())+1] { // +1, так как 1 = понедельник
				break
			}
		}

	case strings.HasPrefix(repeat, "m "):
		// Правило m <через запятую от 1 до 31,-1,-2> [через запятую от 1 до 12]
		monthDaysStr := strings.TrimPrefix(repeat, "m ")
		parts := strings.Split(monthDaysStr, " ")

		if len(parts) == 0 {
			return "", errors.New("no days specified")
		}

		// Обработка дней месяца
		daysOfMonth := strings.Split(parts[0], ",")
		months := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12} // Все месяцы по умолчанию

		// Если указаны месяцы
		if len(parts) > 1 {
			monthStrs := strings.Split(parts[1], ",")
			months = []int{}
			for _, monthStr := range monthStrs {
				month, err := strconv.Atoi(monthStr)
				if err == nil && month >= 1 && month <= 12 {
					months = append(months, month)
				} else {
					return "", errors.New("unknown month in repeat rule")
				}
			}
		}

		nextDate = startDate
		found := false
		for {
			nextDate = nextDate.AddDate(0, 0, 1) // Переходим к следующему дню
			for _, dayStr := range daysOfMonth {
				day, err := strconv.Atoi(dayStr)
				if err == nil {
					// Проверяем, является ли день последним или предпоследним
					if (day == -1 && nextDate.Day() == nextDate.AddDate(0, 1, 0).Day()-1) || // Последний день месяца
						(day == -2 && nextDate.Day() == nextDate.AddDate(0, 1, 0).Day()-2) { // Предпоследний день месяца
						found = true
						break
					}
					// Проверяем, если день больше 0
					if day > 0 && nextDate.Day() == day {
						// Проверяем, что месяц в списке
						if len(months) == 0 || сontains(months, int(nextDate.Month())) {
							found = true
							break
						}
					}
				}
			}
			if found {
				break
			}
		}

	default:
		return "", errors.New("unknown repeat rule")
	}

	// Проверяем, что следующая дата больше текущей и не равна now
	if nextDate.Before(now) {
		return "", errors.New("next date is not greater than now")
	}
	if nextDate.Equal(now) {
		return "", errors.New("next date must be greater than now")
	}

	return nextDate.Format("20060102"), nil
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
