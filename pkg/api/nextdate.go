package api

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func afterNow(date, now time.Time) bool {
	return date.Unix() > now.Unix()
}

func sliceConvert(slice []string) ([]int, error) {
	var resSlice []int
	for _, s := range slice {
		if s == "" {
			continue
		}
		n, err := strconv.Atoi(s)
		if err != nil {
			return nil, fmt.Errorf("ошибка конвертации строки в число %w", err)
		}
		resSlice = append(resSlice, n)
	}
	return resSlice, nil
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", fmt.Errorf("не указан repeat")
	}

	date, err := time.Parse("20060102", dstart)
	if err != nil {
		return "", fmt.Errorf("не удалось распознать дату %w", err)
	}

	parts := strings.Split(repeat, " ")
	if len(parts) == 0 {
		return "", fmt.Errorf("некорректный формат repeat")
	}

	var a, b, c string
	if len(parts) == 1 {
		a = parts[0]
	}
	if len(parts) == 2 {
		a = parts[0]
		b = parts[1]
	}
	if len(parts) == 3 {
		a = parts[0]
		b = parts[1]
		c = parts[2]
	}
	if len(parts) > 3 {
		return "", fmt.Errorf("некорректный формат ввода %s", repeat)
	}

	bSlice := strings.Split(b, ",")
	cSlice := strings.Split(c, ",")

	bConv, err := sliceConvert(bSlice)
	if err != nil {
		return "", fmt.Errorf("days %w", err)
	}

	cConv, err := sliceConvert(cSlice)
	if err != nil {
		return "", fmt.Errorf("months %w", err)
	}

	switch a {
	case "d":
		if len(bConv) == 0 {
			return "", fmt.Errorf("не указан интервал в днях")
		}
		if bConv[0] > 400 {
			return "", fmt.Errorf("превышен максимально допустимый интервал")
		}

		for !afterNow(date, now) {
			date = date.AddDate(0, 0, bConv[0])
		}
		return date.Format("20060102"), nil

	case "y":
		for !afterNow(date, now) {
			date = date.AddDate(1, 0, 0)
		}
		return date.Format("20060102"), nil

	case "w":
		if len(bConv) == 0 {
			return "", fmt.Errorf("не указан интервал в днях недели")
		}
		for _, n := range bConv {
			if n < 1 || n > 7 {
				return "", fmt.Errorf("недопустимое значение дня недели %d", n)
			}
		}
		for !afterNow(date, now) {
			curW := int(date.Weekday())
			if curW == 0 {
				curW = 7
			}
			found := false
			for _, n := range bConv {
				if n == curW {
					found = true
				}
			}
			if found {
				break
			}
			date = date.AddDate(0, 0, 1)
		}
		return date.Format("20060102"), nil

	case "m":
		if len(bConv) == 0 {
			return "", fmt.Errorf("не указан интервал в днях %w", err)
		}

		var day [32]bool
		var month [13]bool

		for _, n := range bConv {
			if n < -2 || n > 31 {
				return "", fmt.Errorf("недопустимый день месяца: %d", n)
			}
			switch n {
			case -1:
				day[len(day)-1] = true
			case -2:
				day[len(day)-2] = true
			default:
				day[n] = true
			}
		}

		if len(cConv) == 0 {
			for i := 1; i <= 12; i++ {
				month[i] = true
			}
		} else {
			for _, m := range cConv {
				if m < 1 || m > 12 {
					return "", fmt.Errorf("недопустимый месяц: %d", m)
				}
				month[m] = true
			}
		}
		for !afterNow(date, now) {
			date = date.AddDate(0, 0, 1)
			d := date.Day()
			m := int(date.Month())
			if d < len(day) && day[d] && month[m] {
				break
			}
		}
		return date.Format("20060102"), nil

	default:
		return "", fmt.Errorf("недопустимый символ %s", a)
	}
}
