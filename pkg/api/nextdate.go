package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const dateFormat = "20060102"

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", fmt.Errorf("repeat rule is not specified")
	}

	date, err := time.Parse(dateFormat, dstart)
	if err != nil {
		return "", fmt.Errorf("invalid date %q: %w", dstart, err)
	}

	parts := strings.Fields(repeat)

	switch parts[0] {
	case "d":
		interval, err := parseDayInterval(parts)
		if err != nil {
			return "", err
		}
		for {
			date = date.AddDate(0, 0, interval)
			if afterNow(date, now) {
				break
			}
		}

	case "y":
		if len(parts) != 1 {
			return "", fmt.Errorf("invalid repeat format %q", repeat)
		}
		for {
			date = date.AddDate(1, 0, 0)
			if afterNow(date, now) {
				break
			}
		}

	case "w":
		weekdays, err := parseWeekdays(parts)
		if err != nil {
			return "", err
		}
		for {
			date = date.AddDate(0, 0, 1)
			if weekdays[isoWeekday(date)] && afterNow(date, now) {
				break
			}
		}

	case "m":
		days, months, err := parseMonthDays(parts)
		if err != nil {
			return "", err
		}
		for {
			date = date.AddDate(0, 0, 1)
			if monthMatches(date, months) && dayMatches(date, days) && afterNow(date, now) {
				break
			}
		}

	default:
		return "", fmt.Errorf("unsupported repeat format %q", repeat)
	}

	return date.Format(dateFormat), nil
}

func afterNow(date, now time.Time) bool {
	y1, m1, d1 := date.Date()
	y2, m2, d2 := now.Date()
	dateOnly := time.Date(y1, m1, d1, 0, 0, 0, 0, time.UTC)
	nowOnly := time.Date(y2, m2, d2, 0, 0, 0, 0, time.UTC)
	return dateOnly.After(nowOnly)
}

func parseDayInterval(parts []string) (int, error) {
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid repeat format for 'd' rule")
	}
	interval, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, fmt.Errorf("invalid day interval: %w", err)
	}
	if interval < 1 || interval > 400 {
		return 0, fmt.Errorf("day interval must be between 1 and 400")
	}
	return interval, nil
}

func parseWeekdays(parts []string) (map[int]bool, error) {
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid repeat format for 'w' rule")
	}
	values, err := parseIntList(parts[1], 1, 7)
	if err != nil {
		return nil, fmt.Errorf("invalid weekday list: %w", err)
	}
	set := make(map[int]bool, len(values))
	for _, v := range values {
		set[v] = true
	}
	return set, nil
}

func parseMonthDays(parts []string) (days []int, months map[int]bool, err error) {
	if len(parts) < 2 || len(parts) > 3 {
		return nil, nil, fmt.Errorf("invalid repeat format for 'm' rule")
	}
	days, err = parseDayOfMonthList(parts[1])
	if err != nil {
		return nil, nil, fmt.Errorf("invalid day-of-month list: %w", err)
	}
	if len(parts) == 3 {
		values, err := parseIntList(parts[2], 1, 12)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid month list: %w", err)
		}
		months = make(map[int]bool, len(values))
		for _, v := range values {
			months[v] = true
		}
	}
	return days, months, nil
}

func parseIntList(s string, min, max int) ([]int, error) {
	items := strings.Split(s, ",")
	result := make([]int, 0, len(items))
	for _, item := range items {
		n, err := strconv.Atoi(item)
		if err != nil {
			return nil, fmt.Errorf("invalid value %q", item)
		}
		if n < min || n > max {
			return nil, fmt.Errorf("value %d out of range [%d, %d]", n, min, max)
		}
		result = append(result, n)
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("empty value list")
	}
	return result, nil
}

func parseDayOfMonthList(s string) ([]int, error) {
	items := strings.Split(s, ",")
	result := make([]int, 0, len(items))
	for _, item := range items {
		n, err := strconv.Atoi(item)
		if err != nil {
			return nil, fmt.Errorf("invalid value %q", item)
		}
		if n == 0 || n < -2 || n > 31 {
			return nil, fmt.Errorf("day of month %d out of range", n)
		}
		result = append(result, n)
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("empty value list")
	}
	return result, nil
}

func isoWeekday(date time.Time) int {
	wd := int(date.Weekday())
	if wd == 0 {
		wd = 7
	}
	return wd
}

func lastDayOfMonth(date time.Time) int {
	firstOfNextMonth := time.Date(date.Year(), date.Month()+1, 1, 0, 0, 0, 0, date.Location())
	return firstOfNextMonth.AddDate(0, 0, -1).Day()
}

func dayMatches(date time.Time, days []int) bool {
	day := date.Day()
	last := lastDayOfMonth(date)
	for _, v := range days {
		switch {
		case v > 0 && v == day:
			return true
		case v == -1 && day == last:
			return true
		case v == -2 && day == last-1:
			return true
		}
	}
	return false
}

func monthMatches(date time.Time, months map[int]bool) bool {
	if months == nil {
		return true
	}
	return months[int(date.Month())]
}

func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowParam := r.FormValue("now")
	dateParam := r.FormValue("date")
	repeat := r.FormValue("repeat")

	now := time.Now()
	if nowParam != "" {
		parsedNow, err := time.Parse(dateFormat, nowParam)
		if err != nil {
			http.Error(w, fmt.Sprintf("invalid now parameter: %v", err), http.StatusBadRequest)
			return
		}
		now = parsedNow
	}

	next, err := NextDate(now, dateParam, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte(next))
}
