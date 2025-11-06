package api

import (
	"net/http"

	"final_project/pkg/db"
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
