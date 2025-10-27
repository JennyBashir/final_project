package api

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"
)

var now = time.Now()

var flag bool

func afterNow(date, now time.Time) bool {
	if date.Second() > now.Second() {
		flag = true
	}
	return flag
}

// принимает "время сейчас", исходное время от которого начинается
// отсчет повторений "20060102" и правило повторений в формате
// w7 / d 1 / y / m 3 1,2,6 и тд
func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	date, err := time.Parse("20060102", dstart)
	if err != nil {
		fmt.Errorf("не удалось распознать дату %w", err)
	}
	//разделяю репит на части
	parts := strings.Split(repeat, " ")
	//части репита
	a := parts[0]
	b := parts[1:]
	c := parts[2:]

	//разделяю части на части
	bSlice := strings.Split(b[0], ",")
	cSlice := strings.Split(c[0], ",")
	//новые слайсы для конвертированных частей
	bConv := make([]int, 0, 31)
	cConv := make([]int, 0, 12)

	//конвертация частей
	for _, num := range bSlice {
		num, err := strconv.Atoi(num)
		if err != nil {
			fmt.Errorf("ошибка конвертации строки в число (дни)")
		}
		bConv = append(bConv, num)
	}
	for _, nums := range cSlice {
		nums, err := strconv.Atoi(nums)
		if err != nil {
			fmt.Errorf("ошибка конвертации строки в число (месяцы)")
		}
		cConv = append(cConv, nums)
	}
	switch {
	case repeat == "":
		_, err := db.Exec("DELETE FROM scheduler WHERE repeat = :repeat", sql.Named("repeat", repeat))
		if err != nil {
			fmt.Errorf("не удалось удалить задачу")
		}
		return "", nil
	case a == "d":
		if bConv == nil {
			fmt.Errorf("не указан интервал в днях")
		}
		for _, n := range bConv {
			if n > 400 {
				fmt.Errorf("превышен максимально допустимый интервал")
			}
			for {
				date = date.AddDate(0, 0, n)
				if afterNow(date, now) {
					break
				}

			}
			return "", nil
		}
	case a == "y":
		for {
			date = date.AddDate(1, 0, 0)
			if afterNow(date, now) {
				break
			}
			return "", nil
		}
	case a == "w":
		if bConv == nil {
			fmt.Errorf("не указан интервал в днях недели")
		}
		for _, n := range bConv {
			if n > 7 {
				fmt.Errorf("недопустимое значение дня недели")
			}
			if len(bConv) == 1 {
				//переносим на конкретный день недели
				// 1-пн, 2-вт, 3-ср, 4-чт, 5-пт, 6-сб, 7-вс
				weekDays := []time.Weekday{time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday, time.Saturday, time.Sunday}
				day := weekDays[0]
				switch {
				case bConv[0] == 1:
					day = weekDays[0]
				case bConv[0] == 2:
					day = weekDays[1]
				case bConv[0] == 3:
					day = weekDays[2]
				case bConv[0] == 4:
					day = weekDays[3]
				case bConv[0] == 5:
					day = weekDays[4]
				case bConv[0] == 6:
					day = weekDays[5]
				case bConv[0] == 7:
					day = weekDays[6]
				}
				res := (int(day) - int(date.Weekday()) + 7) % 7
				if res == 0 {
					res = 7
				}
				date = date.AddDate(0, 0, res)
			}
		}
		if len(bConv) > 1 {
			//переносим на любой из указанных дней
		}

	case a == "m":
		if bConv == nil {
			fmt.Errorf("не указан интервал в днях")
		}
		for _, n := range bConv {
			if n > 31 {
				fmt.Errorf("недопустимый день месяца")
			}
		}
		for _, m := range cConv {
			if m > 12 {
				fmt.Errorf("недопустимый месяц")
			}
		}
		switch {
		case len(bConv) == 1 || len(cConv) == 0:
			//переносим на указанное число КАЖДОГО месяца

		case len(bConv) == 1 || len(cConv) == 1:
			//переносим на указанный день УКАЗАННОГО месяца
			//-1 - последний день месяца
			//-2 - предпоследний день месяца
			//1 - первый день
			//2 - второй день и тд

		case len(bConv) > 1 || len(cConv) == 0:
			//переносим на все указанные даты КАЖДОГО месяца

		case len(bConv) > 1 || len(cConv) > 1:
			//переносим на все УКАЗАННЫЕ ДНИ всех указанных месяцев

			//dafault?
		case a != "d" || a != "y" || a != "w" || a != "m":
			fmt.Errorf("недопустимый символ %s: %w", a, err)

		}
	}
	return "", nil
}
