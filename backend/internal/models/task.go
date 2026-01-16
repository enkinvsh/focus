package models

import "time"

type Task struct {
	ID            int64      `json:"id"`
	UserID        int64      `json:"user_id"`
	Title         string     `json:"title"`
	OriginalInput string     `json:"original,omitempty"`
	TaskType      string     `json:"type"`
	Priority      int        `json:"priority"`
	Completed     bool       `json:"completed"`
	CompletedAt   *time.Time `json:"completed_at,omitempty"`
	DueAt         *time.Time `json:"due_at,omitempty"`
	ReminderSent  bool       `json:"reminder_sent"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type CreateTaskRequest struct {
	Title    string `json:"title" binding:"required"`
	Type     string `json:"type" binding:"required,oneof=Task Long Routine"`
	Priority int    `json:"priority" binding:"min=1,max=3"`
	Original string `json:"original"`
}

type UpdateTaskRequest struct {
	Title     *string `json:"title"`
	Priority  *int    `json:"priority"`
	Completed *bool   `json:"completed"`
}
