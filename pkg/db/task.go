package db

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Task struct {
	ID      string `json:"id" db:"id"`
	Date    string `json:"date" db:"date"`
	Title   string `json:"title" db:"title"`
	Comment string `json:"comment" db:"comment"`
	Repeat  string `json:"repeat" db:"repeat"`
}

func AddTask(task *Task) (int64, error) {
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func GetTask(id string) (*Task, error) {
	var t Task
	err := db.QueryRow(
		`SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`, id).
		Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("задача не найдена")
	}
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func UpdateTask(task *Task) error {
	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("задача не найдена")
	}
	return nil
}

func UpdateDate(next string, id string) error {
	res, err := db.Exec(`UPDATE scheduler SET date = ? WHERE id = ?`, next, id)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("задача не найдена")
	}
	return nil
}

func DeleteTask(id string) error {
	res, err := db.Exec(`DELETE FROM scheduler WHERE id = ?`, id)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("задача не найдена")
	}
	return nil
}

func Tasks(limit int) ([]*Task, error) {
	return selectTasks(
		`SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT ?`,
		limit)
}

func TasksSearch(search string, limit int) ([]*Task, error) {
	if search == "" {
		return Tasks(limit)
	}

	if d, err := time.Parse("02.01.2006", search); err == nil {
		return selectTasks(
			`SELECT id, date, title, comment, repeat FROM scheduler WHERE date = ? ORDER BY date LIMIT ?`,
			d.Format("20060102"), limit)
	}

	like := "%" + search + "%"
	return selectTasks(
		`SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date LIMIT ?`,
		like, like, limit)
}

func selectTasks(query string, args ...any) ([]*Task, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]*Task, 0)
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
			return nil, err
		}
		tasks = append(tasks, &t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}
