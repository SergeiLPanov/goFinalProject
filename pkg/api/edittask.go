package api

import (
	"net/http"
	"strconv"

	"goFinalProject/pkg/db"
)

func editTaskHandler(w http.ResponseWriter, r *http.Request) {
	task, err := decodeTask(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if task.ID == "" {
		writeError(w, http.StatusBadRequest, "не указан идентификатор")
		return
	}
	if _, err := strconv.ParseInt(task.ID, 10, 64); err != nil {
		writeError(w, http.StatusBadRequest, "некорректный идентификатор")
		return
	}

	if err := db.UpdateTask(&task); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{})
}
