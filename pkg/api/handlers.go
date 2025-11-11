package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"final_project/pkg/db"
)

func getHandler(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("id")
	if id == "" {
		writeJson(w, map[string]string{"error": "ID not specified"})
		return
	}

	res, err := db.GetTask(id)
	if err != nil {
		writeJson(w, map[string]string{"error": "task not found"})
		return
	}

	writeJson(w, res)
}

func putHandler(w http.ResponseWriter, r *http.Request) {
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

	err = db.UpdateTask(&task)
	if err != nil {
		writeJson(w, map[string]string{"error": "couldn't change task"})
		return
	}

	writeJson(w, map[string]interface{}{})
}

func doneHandler(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("id")
	if id == "" {
		writeJson(w, map[string]string{"error": "ID not specified"})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJson(w, map[string]string{"error": "task not found"})
		return
	}

	if task.Repeat == "" {
		err = db.DeleteTask(id)
		if err != nil {
			writeJson(w, map[string]string{"error": "couldn't delete task"})
			return
		}
		writeJson(w, map[string]interface{}{})
	} else {
		now := time.Now()
		newDate, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			writeJson(w, map[string]string{"error": "couldn't calculate the task date"})
			return
		}

		err = db.UpdateDate(newDate, id)
		if err != nil {
			writeJson(w, map[string]string{"error": "couldn't change the task date"})
			return
		}
		writeJson(w, map[string]interface{}{})
	}
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("id")
	if id == "" {
		writeJson(w, map[string]string{"error": "ID not specified"})
	}

	err := db.DeleteTask(id)
	if err != nil {
		writeJson(w, map[string]string{"error": "couldn't delete task"})
		return
	}

	writeJson(w, map[string]interface{}{})
}
