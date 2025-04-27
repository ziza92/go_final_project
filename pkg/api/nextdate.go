package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const ConstDate = "20060102"

func afterNow(date, now time.Time) bool {

	dateYear, dateMonth, dateDay := date.Date()
	nowYear, nowMonth, nowDay := now.Date()
	dateTruncated := time.Date(dateYear, dateMonth, dateDay, 0, 0, 0, 0, time.UTC)
	nowTruncated := time.Date(nowYear, nowMonth, nowDay, 0, 0, 0, 0, time.UTC)
	return dateTruncated.After(nowTruncated)
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {

	if repeat == "" {
		return "", fmt.Errorf("повтор не задан")
	}

	date, err := time.Parse(ConstDate, dstart)
	if err != nil {
		return "", fmt.Errorf("неверный формат даты: %s", dstart)
	}

	parts := strings.Split(repeat, " ")

	switch parts[0] {
	case "d":
		if len(parts) != 2 {
			return "", fmt.Errorf("неверный формат даты: %s", repeat)
		}
		interval, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", fmt.Errorf("неверный формат даты: %s", repeat)
		} else if interval < 1 || interval > 400 {
			return "", fmt.Errorf("задан недопустимый интервал повторения: %d", interval)
		}

		for {
			date = date.AddDate(0, 0, interval)
			if afterNow(date, now) {
				return date.Format(ConstDate), nil
			}
		}

	case "y":
		for {
			date = date.AddDate(1, 0, 0)
			if afterNow(date, now) {
				return date.Format(ConstDate), nil
			}
		}
	case "w":
		if len(parts) < 2 {
			return "", fmt.Errorf("нужен аргумент после w")
		}

		days := []int{}
		numbers := strings.Split(parts[1], ",")
		for _, number := range numbers {
			day, err := strconv.Atoi(number)
			if err != nil || day < 1 || day > 7 {
				return "", fmt.Errorf("ошибка парсинга дня недели: %w", err)
			}
			days = append(days, day)
		}

		// ищем ближайший день недели
		for {
			if afterNow(date, now) {
				currentDay := weekdayIso(date.Weekday())
				for _, d := range days {
					if d == currentDay {
						return date.Format(ConstDate), nil
					}
				}
			}
			date = date.AddDate(0, 0, 1)
		}

	default:
		return "", fmt.Errorf("неподдерживаемый формат даты: %s", repeat)
	}
}

func nextDayHandler(w http.ResponseWriter, r *http.Request) {

	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")

	var now time.Time
	if nowStr == "" {
		now = time.Now()
	} else {
		var err error
		now, err = time.Parse(ConstDate, nowStr)
		if err != nil {
			http.Error(w, fmt.Sprintf("invalid now format: %v", err), http.StatusBadRequest)
			return
		}
	}

	if dateStr == "" {
		http.Error(w, "date parameter is required", http.StatusBadRequest)
		return
	}
	if repeat == "" {
		http.Error(w, "repeat parameter is required", http.StatusBadRequest)
		return
	}

	nextDate, err := NextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprint(w, nextDate)
}
func weekdayIso(w time.Weekday) int {
	if w == time.Sunday {
		return 7
	}
	return int(w)
}
