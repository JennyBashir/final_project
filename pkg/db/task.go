package db

import (
	"database/sql"
	"fmt"
)

type Task struct {
	ID      int64  `json:"id"`
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
