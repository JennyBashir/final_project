package api

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func afterNow(date, now time.Time) bool {
	return date.Format("20060102") > now.Format("20060102")
}

func sliceConvert(slice []string) ([]int, error) {
	var resSlice []int

	for _, s := range slice {
		if s == "" {
			continue
		}

		n, err := strconv.Atoi(s)
		if err != nil {
			return nil, fmt.Errorf("error converting a string to a number %w", err)
		}

		resSlice = append(resSlice, n)
	}

	return resSlice, nil
}

func nextDayHandler(w http.ResponseWriter, req *http.Request) {

	layout := "20060102"
	now := time.Now()

	if n := req.FormValue("now"); n != "" {
		parsedNow, err := time.Parse(layout, n)
		if err != nil {
			http.Error(w, "invalid now parameter", http.StatusBadRequest)
			return
		}
		now = parsedNow
	}

	stringResponse, err := NextDate(now, req.FormValue("date"), req.FormValue("repeat"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if _, err := io.WriteString(w, stringResponse); err != nil {
		http.Error(w, "failed to write response", http.StatusInternalServerError)
		return
	}
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {

	if repeat == "" {
		return "", fmt.Errorf("repeat not specified")
	}

	date, err := time.Parse("20060102", dstart)
	if err != nil {
		return "", fmt.Errorf("couldn't recognize the date %w", err)
	}

	parts := strings.Split(repeat, " ")
	if len(parts) == 0 {
		return "", fmt.Errorf("incorrect repeat format")
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
		return "", fmt.Errorf("incorrect input format %s", repeat)
	}

	bSlice := strings.Split(b, ",")
	cSlice := strings.Split(c, ",")

	bConv, err := sliceConvert(bSlice)
	if err != nil {
		return "", fmt.Errorf("invalid days format: %w", err)
	}

	cConv, err := sliceConvert(cSlice)
	if err != nil {
		return "", fmt.Errorf("invalid months format: %w", err)
	}

	switch a {

	case "d":
		if len(bConv) == 0 {
			return "", fmt.Errorf("the interval in days is not specified")
		}

		if bConv[0] > 400 {
			return "", fmt.Errorf("the maximum allowed interval has been exceeded")
		}

		if bConv[0] == 1 && !afterNow(date, now) {
			return now.Format("20060102"), nil
		}

		for {
			date = date.AddDate(0, 0, bConv[0])
			if afterNow(date, now) {
				break
			}
		}
		return date.Format("20060102"), nil

	case "y":
		for {
			date = date.AddDate(1, 0, 0)
			if afterNow(date, now) {
				break
			}
		}
		return date.Format("20060102"), nil

	case "w":
		if len(bConv) == 0 {
			return "", fmt.Errorf("the interval in days of the week is not specified")
		}

		for _, n := range bConv {
			if n < 1 || n > 7 {
				return "", fmt.Errorf("invalid day of the week value %d", n)
			}
		}

		for {
			date = date.AddDate(0, 0, 1)
			if afterNow(date, now) {
				curW := int(date.Weekday())
				if curW == 0 {
					curW = 7
				}

				for _, n := range bConv {
					if n == curW {
						return date.Format("20060102"), nil
					}
				}
			}
		}

	case "m":
		if len(bConv) == 0 {
			return "", fmt.Errorf("the interval in days is not specified %w", err)
		}

		var day [32]bool
		var month [13]bool

		for _, n := range bConv {
			if n < -2 || n > 31 {
				return "", fmt.Errorf("invalid day of the month: %d", n)
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
					return "", fmt.Errorf("invalid month: %d", m)
				}
				month[m] = true
			}
		}

		for {
			date = date.AddDate(0, 0, 1)
			if afterNow(date, now) {
				d := date.Day()
				m := int(date.Month())

				if !month[m] {
					continue
				}
				lastDay := time.Date(date.Year(), date.Month()+1, 0, 0, 0, 0, 0, date.Location()).Day()
				secondLastDay := lastDay - 1

				isRegularDay := d < len(day) && day[d]
				isLastDay := day[len(day)-1] && d == lastDay
				isSecondLastDay := day[len(day)-2] && d == secondLastDay

				if day[31] && d == lastDay && lastDay < 31 {
					date = time.Date(date.Year(), date.Month()+1, 0, 0, 0, 0, 0, date.Location())
					continue
				}

				if isRegularDay || isLastDay || isSecondLastDay {
					break
				}
			}
		}
		return date.Format("20060102"), nil

	default:
		return "", fmt.Errorf("invalid character %s", a)
	}
}
