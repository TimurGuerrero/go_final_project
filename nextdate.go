package main

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func NextDate(now time.Time, date string, repeat string) (string, error) {
	startDate, err := time.Parse("20060102", date)
	if err != nil {
		return "", errors.New("invalid date format")
	}

	// Проверяем, если правило не указано
	if repeat == "" {
		return "", errors.New("no repeat rule specified")
	}

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
				break
			}
		}
	case repeat == "y":
		// Правило y
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

	default:
		return "", errors.New("unknown repeat rule")
	}

	// Проверяем, что следующая дата больше текущей
	if nextDate.Before(now) {
		return "", errors.New("next date is not greater than now")
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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte(nextDate))
}
