package main

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type Task struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Done        bool      `json:"done"`
	DueDate     string    `json:"due_date"`
	Priority    string    `json:"priority"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type TaskInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Done        *bool  `json:"done,omitempty"`
	DueDate     string `json:"due_date"`
	Priority    string `json:"priority"`
	Status      string `json:"status"`
}

func normalizeInput(input TaskInput, previous *Task) (Task, error) {
	t := Task{}
	if previous != nil {
		t = *previous
	}
	t.Title = strings.TrimSpace(input.Title)
	t.Description = strings.TrimSpace(input.Description)
	t.DueDate = strings.TrimSpace(input.DueDate)
	if t.Title == "" || len([]rune(t.Title)) > 160 {
		return Task{}, errors.New("Название должно содержать от 1 до 160 символов")
	}
	if len([]rune(t.Description)) > 2000 {
		return Task{}, errors.New("Описание слишком длинное (максимум 2000 символов)")
	}
	if t.DueDate != "" {
		if _, err := time.Parse("2006-01-02", t.DueDate); err != nil {
			return Task{}, fmt.Errorf("Неверная дата, нужен формат YYYY-MM-DD")
		}
	}
	priority := strings.ToLower(strings.TrimSpace(input.Priority))
	if priority == "" {
		if previous != nil {
			priority = previous.Priority
		} else {
			priority = "medium"
		}
	}
	if priority != "low" && priority != "medium" && priority != "high" {
		return Task{}, errors.New("Приоритет: low, medium или high")
	}
	t.Priority = priority
	status := strings.ToLower(strings.TrimSpace(input.Status))
	if status == "" {
		if input.Done != nil {
			if *input.Done {
				status = "done"
			} else {
				status = "todo"
			}
		} else if previous != nil {
			status = previous.Status
		} else {
			status = "todo"
		}
	}
	if status != "todo" && status != "in_progress" && status != "done" {
		return Task{}, errors.New("Статус: todo, in_progress или done")
	}
	t.Status = status
	t.Done = status == "done"
	return t, nil
}

type Store interface {
	List() ([]Task, error)
	Get(id int64) (Task, bool, error)
	Create(input Task) (Task, error)
	Update(id int64, input Task) (Task, bool, error)
	Delete(id int64) (bool, error)
}
