package api

import (
	"context"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/enkinvsh/focus-backend/internal/db"
	"github.com/enkinvsh/focus-backend/internal/models"
	"github.com/enkinvsh/focus-backend/internal/services"
	"github.com/gin-gonic/gin"
)

const (
	defaultLimit = 50
	maxLimit     = 200
	maxAudioSize = 5 * 1024 * 1024 // 5MB
)

func GetTasks(c *gin.Context) {
	user := GetUser(c)
	taskType := c.DefaultQuery("type", "Task")
	completed := c.DefaultQuery("completed", "false") == "true"

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", strconv.Itoa(defaultLimit)))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if limit <= 0 || limit > maxLimit {
		limit = defaultLimit
	}
	if offset < 0 {
		offset = 0
	}

	rows, err := db.Pool.Query(context.Background(), `
		SELECT id, title, original_input, task_type, priority, completed, created_at
		FROM tasks 
		WHERE user_id = $1 AND task_type = $2 AND completed = $3
		ORDER BY priority ASC, created_at DESC
		LIMIT $4 OFFSET $5
	`, user.ID, taskType, completed, limit, offset)
	if err != nil {
		log.Printf("GetTasks error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch tasks"})
		return
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var t models.Task
		if err := rows.Scan(&t.ID, &t.Title, &t.OriginalInput, &t.TaskType, &t.Priority, &t.Completed, &t.CreatedAt); err != nil {
			log.Printf("GetTasks scan error: %v", err)
			continue
		}
		t.UserID = user.ID
		tasks = append(tasks, t)
	}

	if tasks == nil {
		tasks = []models.Task{}
	}
	c.JSON(http.StatusOK, gin.H{"tasks": tasks, "limit": limit, "offset": offset})
}

func CreateTask(c *gin.Context) {
	user := GetUser(c)

	var req models.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if req.Priority == 0 {
		req.Priority = 2
	}

	var task models.Task
	err := db.Pool.QueryRow(context.Background(), `
		INSERT INTO tasks (user_id, title, original_input, task_type, priority)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`, user.ID, req.Title, req.Original, req.Type, req.Priority).Scan(&task.ID, &task.CreatedAt)

	if err != nil {
		log.Printf("CreateTask error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create task"})
		return
	}

	task.UserID = user.ID
	task.Title = req.Title
	task.OriginalInput = req.Original
	task.TaskType = req.Type
	task.Priority = req.Priority
	task.Completed = false

	c.JSON(http.StatusCreated, task)
}

func UpdateTask(c *gin.Context) {
	user := GetUser(c)
	taskID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}

	var req models.UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	var completedAt *time.Time
	if req.Completed != nil && *req.Completed {
		now := time.Now()
		completedAt = &now
	}

	result, err := db.Pool.Exec(context.Background(), `
		UPDATE tasks 
		SET 
			title = COALESCE($1, title),
			priority = COALESCE($2, priority),
			completed = COALESCE($3, completed),
			completed_at = $4,
			updated_at = NOW()
		WHERE id = $5 AND user_id = $6
	`, req.Title, req.Priority, req.Completed, completedAt, taskID, user.ID)

	if err != nil {
		log.Printf("UpdateTask error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update task"})
		return
	}

	if result.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func DeleteTask(c *gin.Context) {
	user := GetUser(c)
	taskID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}

	result, err := db.Pool.Exec(context.Background(), `
		DELETE FROM tasks WHERE id = $1 AND user_id = $2
	`, taskID, user.ID)

	if err != nil {
		log.Printf("DeleteTask error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete task"})
		return
	}

	if result.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func GetPreferences(c *gin.Context) {
	user := GetUser(c)

	var u models.User
	err := db.Pool.QueryRow(context.Background(), `
		SELECT language, timezone, theme_index FROM users WHERE id = $1
	`, user.ID).Scan(&u.Language, &u.Timezone, &u.ThemeIndex)

	if err != nil {
		u = models.User{Language: "en", Timezone: "UTC", ThemeIndex: 0}
	}

	c.JSON(http.StatusOK, u)
}

func UpdatePreferences(c *gin.Context) {
	user := GetUser(c)

	var req struct {
		Language   *string `json:"language"`
		Timezone   *string `json:"timezone"`
		ThemeIndex *int    `json:"theme_index"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	_, err := db.Pool.Exec(context.Background(), `
		INSERT INTO users (id, first_name, username, language, timezone, theme_index)
		VALUES ($1, $2, $3, COALESCE($4, 'en'), COALESCE($5, 'UTC'), COALESCE($6, 0))
		ON CONFLICT (id) DO UPDATE SET
			language = COALESCE($4, users.language),
			timezone = COALESCE($5, users.timezone),
			theme_index = COALESCE($6, users.theme_index),
			updated_at = NOW()
	`, user.ID, user.FirstName, user.Username, req.Language, req.Timezone, req.ThemeIndex)

	if err != nil {
		log.Printf("UpdatePreferences error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update preferences"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func CreateTaskFromAudio(c *gin.Context) {
	user := GetUser(c)

	file, header, err := c.Request.FormFile("audio")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "audio file required"})
		return
	}
	defer file.Close()

	if header.Size > maxAudioSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "audio file too large (max 5MB)"})
		return
	}

	audioData, err := io.ReadAll(io.LimitReader(file, maxAudioSize+1))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read audio"})
		return
	}
	if int64(len(audioData)) > maxAudioSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "audio file too large (max 5MB)"})
		return
	}

	mimeType := header.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = "audio/webm"
	}

	taskType := c.DefaultPostForm("type", "Task")
	language := c.DefaultPostForm("language", "en")

	parsedTasks, err := services.TranscribeAndParseTasks(audioData, mimeType, taskType, language)
	if err != nil {
		log.Printf("TranscribeAndParseTasks error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to process audio"})
		return
	}

	var tasks []models.Task
	for _, pt := range parsedTasks {
		title, _ := pt["title"].(string)
		priority := 2
		if p, ok := pt["priority"].(float64); ok {
			priority = int(p)
		}
		tType := taskType
		if t, ok := pt["type"].(string); ok {
			tType = t
		}

		var task models.Task
		err := db.Pool.QueryRow(context.Background(), `
			INSERT INTO tasks (user_id, title, original_input, task_type, priority)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id, created_at
		`, user.ID, title, "[voice]", tType, priority).Scan(&task.ID, &task.CreatedAt)

		if err != nil {
			log.Printf("CreateTaskFromAudio insert error: %v", err)
			continue
		}

		task.UserID = user.ID
		task.Title = title
		task.OriginalInput = "[voice]"
		task.TaskType = tType
		task.Priority = priority
		task.Completed = false
		tasks = append(tasks, task)
	}

	if tasks == nil {
		tasks = []models.Task{}
	}

	c.JSON(http.StatusCreated, gin.H{"tasks": tasks})
}
