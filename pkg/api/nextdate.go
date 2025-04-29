package api

import (
	"fmt"
	"net/http"
	"slices"
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

	case "m":
		if len(parts) < 2 {
			return "", fmt.Errorf("ну указано количество месяцев: %s", repeat)
		}

		partsDay := strings.Split(parts[1], ",")
		days := []int{}
		for _, part := range partsDay {
			day, err := strconv.Atoi(part)
			if err != nil {
				return "", fmt.Errorf("ошибка при распознавании дня: %v", err)
			}
			if day < -31 || day == 0 || day > 31 {
				return "", fmt.Errorf("некорректное значение дня: %d", day)
			}
			days = append(days, day)
		}

		months := []int{}
		if len(parts) >= 3 {
			partsMonth := strings.Split(parts[2], ",")
			for _, part := range partsMonth {
				month, err := strconv.Atoi(part)
				if err != nil {
					return "", fmt.Errorf("ошибка при распознавании месяца: %v", err)
				}
				if month < 1 || month > 12 {
					return "", fmt.Errorf("некорректное значение месяца: %d", month)
				}
				months = append(months, month)
			}
		}

		for {
			if afterNow(date, now) {
				year, month := date.Year(), date.Month()
				currentDay := date.Day()

				if len(months) > 0 {
					matchMonth := slices.Contains(months, int(month))
					if !matchMonth {
						date = date.AddDate(0, 0, 1)
						continue
					}
				}

				lastDayOfMonth := time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
				for _, day := range days {
					if day == -3 {
						return "", fmt.Errorf("некорректное значение дня: %d", day)
					}
					correctDay := day
					if correctDay < 0 {
						correctDay = lastDayOfMonth + day + 1
					}
					if correctDay == currentDay {
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
	if r.Method != http.MethodGet {
		http.Error(w, "method not supported:", http.StatusMethodNotAllowed)
		return
	}

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
