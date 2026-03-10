package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/dojo-harvester/backend/internal/storage"
)

func (s *Server) registerTaskRoutes(api *gin.RouterGroup) {
	api.POST("/auth/login", s.login)
	api.GET("/tasks", s.getTasks)
	api.POST("/tasks/:task_id/claim", s.claimTask)
	api.PUT("/tasks/:task_id/progress", s.updateTaskProgress)
	
	// Dev endpoint to seed tasks
	api.POST("/tasks/seed", s.seedTasks)
}

func (s *Server) login(c *gin.Context) {
	// Mock login for T092
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid credentials"})
		return
	}

	// In real app, proxy to data platform
	c.JSON(http.StatusOK, gin.H{
		"token": "mock-jwt-token-for-" + req.Username,
		"user": gin.H{
			"id": "op-1",
			"username": req.Username,
		},
	})
}

func (s *Server) getTasks(c *gin.Context) {
	// In real app, fetch from platform and cache
	tasks, err := s.store.GetAllTasks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tasks"})
		return
	}
	
	c.JSON(http.StatusOK, tasks)
}

func (s *Server) claimTask(c *gin.Context) {
	taskID := c.Param("task_id")
	
	task, err := s.store.GetTask(taskID)
	if err != nil || task == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}
	
	// Mock current user
	assignee := "op-1"
	task.AssigneeID = &assignee
	task.Status = "in_progress"
	task.UpdatedAt = time.Now().UTC()
	
	if err := s.store.SaveTask(task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to claim task"})
		return
	}
	
	c.JSON(http.StatusOK, task)
}

func (s *Server) updateTaskProgress(c *gin.Context) {
	taskID := c.Param("task_id")
	
	var req struct {
		CompletedEpisodes int `json:"completed_episodes"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	
	task, err := s.store.GetTask(taskID)
	if err != nil || task == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}
	
	task.CompletedEpisodes = req.CompletedEpisodes
	task.UpdatedAt = time.Now().UTC()
	
	if task.CompletedEpisodes >= task.RequiredEpisodes {
		task.Status = "completed"
		now := time.Now().UTC()
		task.CompletedAt = &now
	}
	
	if err := s.store.SaveTask(task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
		return
	}
	
	c.JSON(http.StatusOK, task)
}

func (s *Server) seedTasks(c *gin.Context) {
	now := time.Now().UTC()
	task := &storage.HarvestTask{
		TaskID: "task-123",
		Title: "Pick up the apple",
		Description: "Teach the robot to pick up the red apple from the table",
		RequiredEpisodes: 50,
		CompletedEpisodes: 0,
		Status: "pending",
		Priority: "high",
		CreatedBy: "system",
		CreatedAt: now,
		UpdatedAt: now,
		Objectives: []string{"Approach table", "Locate apple", "Grasp apple firmly", "Lift without dropping"},
		RobotConfiguration: map[string]interface{}{
			"model": "alpha-v2",
			"cameras": []string{"head", "left_wrist", "right_wrist"},
		},
	}
	
	s.store.SaveTask(task)
	c.JSON(http.StatusOK, task)
}
