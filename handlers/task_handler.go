package handlers

import (
	"errors"
	"log"
	"net/http"
	"time"

	storepkg "hw2hard/store"

	"github.com/gin-gonic/gin"
)

func RegisterTaskRoutes(r *gin.Engine, ts storepkg.TaskStore) {
	r.POST("/task", taskHandler(ts))
	r.GET("/status/:task_id", statusHandler(ts))
	r.GET("/result/:task_id", resultHandler(ts))
}

func taskHandler(ts storepkg.TaskStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		taskID, err := ts.CreateTask()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create task"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"task_id": taskID, "status": "created"})

		go func(id string) {
			log.Printf("processing task %s", id)
			time.Sleep(2 * time.Second)
			_ = ts.SetResult(id, "some_junk_payload")
			_ = ts.SetStatus(id, storepkg.StatusReady)
			log.Printf("task %s ready", id)
		}(taskID)
	}
}

func statusHandler(ts storepkg.TaskStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("task_id")
		status, err := ts.GetStatus(id)
		if err != nil {
			if errors.Is(err, storepkg.ErrNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get status"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": status})
	}
}

func resultHandler(ts storepkg.TaskStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("task_id")
		result, err := ts.GetResult(id)
		if err != nil {
			switch {
			case errors.Is(err, storepkg.ErrNotFound):
				c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
			case errors.Is(err, storepkg.ErrNotReady):
				c.JSON(http.StatusConflict, gin.H{"error": "task not ready"})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get result"})
			}
			return
		}
		c.JSON(http.StatusOK, gin.H{"result": result})
	}
}
