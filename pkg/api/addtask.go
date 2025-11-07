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
		//
	default:
		http.Error(w, "метод не определен", http.StatusMethodNotAllowed)
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
		return fmt.Errorf("некорректная дата, формат отличается от 20060102")
	}
	var next string
	if task.Repeat != "" {
		next, err = NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return fmt.Errorf("не удалось просчитать следующую дату для задачи %s %w", task.Title, err)
		}
	}
	if afterNow(now, t) {
		if task.Repeat == "" {
			task.Date = now.Format("20060102")
		} else {
			task.Date = next
		}
	}
	return nil
}

func writeJson(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	resp, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "не удалось преобразовать ответ в json", http.StatusInternalServerError)
		return
	}
	_, err = w.Write(resp)
	if err != nil {
		http.Error(w, "не удалось отправить ответ", http.StatusInternalServerError)
		return
	}
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		writeJson(w, map[string]string{"error": "не удалось прочитать тело запроса"})
		return
	}
	err = json.Unmarshal(buf.Bytes(), &task)
	if err != nil {
		writeJson(w, map[string]string{"error": "не удалось прочитать запрос"})
		return
	}
	if task.Title == "" {
		writeJson(w, map[string]string{"error": "не указан заголовок"})
		return
	}
	err = checkDate(&task)
	if err != nil {
		writeJson(w, map[string]string{"error": "некорректный формат запроса"})
		return
	}
	id, err := db.AddTask(&task)
	if err != nil {
		writeJson(w, map[string]string{"error": "не удалось добавить задачу в базу данных"})
		return
	}
	writeJson(w, map[string]string{"id": fmt.Sprintf("%d", id)})
}
