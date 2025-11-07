package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"final_project/pkg/db"
)

func getHandler(w http.ResponseWriter, r *http.Request) {
	//получаю айди
	id := r.URL.Query().Get("id")
	//если айди не указан
	if id == "" {
		writeJson(w, map[string]string{"error": "не указан идентификатор"})
		return
	}
	//отправляю получать таск по айди
	res, err := db.GetTask(id)
	if err != nil {
		writeJson(w, map[string]string{"error": "задача не найдена"})
		return
	}
	//пишу ответ
	writeJson(w, res)

}

func putHandler(w http.ResponseWriter, r *http.Request) {
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
	err = db.UpdateTask(&task)
	if err != nil {
		writeJson(w, map[string]string{"error": "не удалось изменить задачу"})
		return
	}
	writeJson(w, map[string]interface{}{})
}

func doneHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	//если айди не указан
	if id == "" {
		writeJson(w, map[string]string{"error": "не указан идентификатор"})
		return
	}
	task, err := db.GetTask(id)
	if err != nil {
		writeJson(w, map[string]string{"error": "задача не найдена"})
		return
	}
	if task.Repeat == "" {
		err = db.DeleteTask(id)
		if err != nil {
			writeJson(w, map[string]string{"error": "не удалось удалить задачу"})
			return
		}
		writeJson(w, map[string]interface{}{})
	} else {
		now := time.Now()
		newDate, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			writeJson(w, map[string]string{"error": "не удалось рассчитать дату задачи"})
			return
		}

		err = db.UpdateDate(newDate, id)
		if err != nil {
			writeJson(w, map[string]string{"error": "не удалось изменить дату задачи"})
			return
		}
		writeJson(w, map[string]interface{}{})
	}
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	err := db.DeleteTask(id)
	if err != nil {
		writeJson(w, map[string]string{"error": "не удалось удалить задачу"})
		return
	}
	writeJson(w, map[string]interface{}{})
}
