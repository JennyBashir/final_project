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
