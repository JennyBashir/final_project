package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"final_project/pkg/db"
)

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:

		addTaskHandler(w, r)
	case http.MethodGet:
		getHandler(w, r)
	case http.MethodPut:
		putHandler(w, r)
	case http.MethodDelete:
		deleteHandler(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func checkDate(task *db.Task) error {
	now := time.Now()
	today := now.Format("20060102")
	if task.Date == "" || task.Date == "today" {
		task.Date = today
		return nil
	}

	t, err := time.Parse("20060102", task.Date)
	if err != nil {
		return fmt.Errorf("incorrect date, format differs from 20060102")
	}

	if t.Format("20060102") < today {
		if task.Repeat == "" {
			task.Date = today
		} else {
			next, err := NextDate(now, task.Date, task.Repeat)
			if err != nil {
				return err
			}
			task.Date = next
		}
	}
	return nil
}

func writeJson(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	resp, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "couldn't convert the response to json", http.StatusInternalServerError)
		return
	}
	_, err = w.Write(resp)
	if err != nil {
		http.Error(w, "couldn't send reply", http.StatusInternalServerError)
		return
	}
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		writeJson(w, map[string]string{"error": "couldn't read the request body"})
		return
	}

	err = json.Unmarshal(buf.Bytes(), &task)
	if err != nil {
		writeJson(w, map[string]string{"error": "couldn't read the request"})
		return
	}

	if task.Title == "" {
		writeJson(w, map[string]string{"error": "the title is not specified"})
		return
	}

	err = checkDate(&task)
	if err != nil {
		writeJson(w, map[string]string{"error": "incorrect request format"})
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		writeJson(w, map[string]string{"error": "couldn't add a task to the database"})
		return
	}

	writeJson(w, map[string]string{"id": fmt.Sprintf("%d", id)})
}
