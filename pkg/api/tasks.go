package api

import (
	"bytes"
	"encoding/json"
	"final_project/pkg/db"
	"net/http"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")

	tasks, err := db.Tasks(50, search)
	if err != nil {
		writeJson(w, map[string]string{"error": "не удалось получить ближайшие задачи"})
		return
	}

	writeJson(w, TasksResp{
		Tasks: tasks})
}

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
	writeJson(w, db.Task{})
}
