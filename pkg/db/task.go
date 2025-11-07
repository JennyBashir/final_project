package db

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

type Task struct {
	ID      int64  `json:"id,string"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTask(task *Task) (int64, error) {
	var id int64
	if db == nil {
		return 0, fmt.Errorf("база данных недоступна")
	}

	res, err := db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)",
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))
	if err != nil {
		return 0, err
	}
	id, err = res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func Tasks(limit int, search string) ([]*Task, error) {
	if db == nil {
		return nil, fmt.Errorf("база данных недоступна")
	}

	search = strings.TrimSpace(search)
	var tRows *sql.Rows
	var err error

	if search != "" {
		if t, errD := time.Parse("02.01.2006", search); errD == nil {
			dateStr := t.Format("20060102")
			tRows, err = db.Query(
				"SELECT id, date, title, comment, repeat FROM scheduler WHERE date = ? ORDER BY date LIMIT ?",
				dateStr, limit,
			)
		} else {
			a := "%" + search + "%"
			tRows, err = db.Query(
				"SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE ? COLLATE NOCASE OR comment LIKE ? COLLATE NOCASE ORDER BY date LIMIT ?",
				a, a, limit,
			)
		}
	} else {
		tRows, err = db.Query(
			"SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT ?",
			limit,
		)
	}

	if err != nil {
		return nil, err
	}
	defer tRows.Close()

	var tasks []*Task

	for tRows.Next() {
		var id int64
		var date, title, comment, repeat string

		if err := tRows.Scan(&id, &date, &title, &comment, &repeat); err != nil {
			return nil, err
		}

		tasks = append(tasks, &Task{
			ID:      id,
			Date:    date,
			Title:   title,
			Comment: comment,
			Repeat:  repeat,
		})
	}
	if err := tRows.Err(); err != nil {
		return nil, err
	}

	if tasks == nil {
		tasks = []*Task{}
	}

	return tasks, nil
}

func GetTask(id string) (*Task, error) {
	//по указанному id возвратит структуру или указатель
	//на структуру Task. Так как нужно получить только одну
	//запись, то нужно использовать метод
	//QueryRow() — err := db.QueryRow(...).Scan(...).
	row := db.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = :id",
		sql.Named("id", id))
	//переменные куда запишется результат скана
	var iD int64
	var date, title, comment, repeat string
	//сканируем и записываем в переменные
	err := row.Scan(&iD, &date, &title, &comment, &repeat)
	if err != nil {
		return nil, fmt.Errorf("не удалось просканировать данные")
	}
	//заполняем структуру данными
	task := &Task{
		ID:      iD,
		Date:    date,
		Title:   title,
		Comment: comment,
		Repeat:  repeat,
	}
	//возвращаем структуру
	return task, nil
}

func UpdateTask(task *Task) error {
	query := "UPDATE scheduler SET title = :title, date = :date, comment = :comment, repeat = :repeat WHERE id = :id"
	res, err := db.Exec(query,
		sql.Named("title", task.Title),
		sql.Named("date", task.Date),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
		sql.Named("id", task.ID))
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf(`incorrect id for updating task`)
	}
	return nil
}
