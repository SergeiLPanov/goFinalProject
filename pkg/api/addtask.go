package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"goFinalProject/pkg/db"
)

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	task, err := decodeTask(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"id": fmt.Sprint(id)})
}

func decodeTask(r *http.Request) (db.Task, error) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		return task, fmt.Errorf("некорректный JSON: %w", err)
	}

	if task.Title == "" {
		return task, fmt.Errorf("не указан заголовок задачи")
	}

	if err := checkDate(&task); err != nil {
		return task, err
	}

	return task, nil
}

func checkDate(task *db.Task) error {
	now := time.Now()

	if task.Date == "" {
		task.Date = now.Format(dateFormat)
	}

	t, err := time.Parse(dateFormat, task.Date)
	if err != nil {
		return fmt.Errorf("некорректный формат даты: %w", err)
	}

	var next string
	if task.Repeat != "" {
		next, err = NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return fmt.Errorf("некорректное правило повторения: %w", err)
		}
	}

	if afterNow(now, t) {
		if task.Repeat == "" {
			task.Date = now.Format(dateFormat)
		} else {
			task.Date = next
		}
	}

	return nil
}
