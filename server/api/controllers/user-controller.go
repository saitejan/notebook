package controllers

import (
	"database/sql"
	"fmt"
	"log"
	"notebook/api/models"
	"notebook/api/repositories"
	"strconv"

	"github.com/gin-gonic/gin"
)

// TaskController struct
type TaskController struct {
	userRepository *repositories.UserRepository
}

// Init method
func (c *TaskController) Init(db *sql.DB) {
	c.userRepository = &repositories.UserRepository{}
	c.userRepository.Init(db)
}

// CreateTask method
func (c *TaskController) CreateTask(ctx *gin.Context) {
	var task models.User
	ctx.BindJSON(&task)
	if task.Name == "" {
		ctx.JSON(400, gin.H{
			"error": "name should not be empty",
		})
		return
	}
	// useridi, exists := ctx.Get("number")
	// if !exists {
	// 	ctx.JSON(400, gin.H{
	// 		"error": "number not found in request context",
	// 	})
	// 	return
	// }
	// userid := useridi.(string)
	// task.Number = userid
	createdTask, err := c.userRepository.CreateTask(task)
	if err != nil {
		log.Printf("Error: %v\n", err)
		ctx.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(201, gin.H{
		"task": createdTask,
	})
}

// GetTasks method
func (c *TaskController) GetTasks(ctx *gin.Context) {
	tasks := []models.User{}
	var err error
	tasks, err = c.userRepository.GetAllTasks()
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"tasks": tasks,
	})
}

// GetTaskByID method
func (c *TaskController) GetTaskByID(ctx *gin.Context) {
	idstr := ctx.Param("id")
	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": fmt.Sprintf("%s is not a valid number", idstr),
		})
		return
	}

	task, err := c.userRepository.GetTaskByID(id)
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"task": task,
	})
}

// GetTaskByNameAndNumber method
func (c *TaskController) GetTaskByNameAndNumber(ctx *gin.Context) {
	idstr := ctx.Param("id")
	// number := idstr.(string)

	useridi := ctx.Param("id1")
	// name := useridi.(string)

	task, err := c.userRepository.GetTaskByNameAndNumber(idstr, useridi)
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"task": task,
	})
}

// UpdateTaskForID method
func (c *TaskController) UpdateTaskForID(ctx *gin.Context) {
	idstr := ctx.Param("id")
	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": fmt.Sprintf("%s is not a valid number", idstr),
		})
		return
	}

	existingTask, err := c.userRepository.GetTaskByID(id)
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}

	var receive struct {
		Transactions []interface{} `json:"transactions"`
	}
	ctx.BindJSON(&receive)
	maxlen := 10
	if len(existingTask.Transactions) < 10 {
		maxlen = len(existingTask.Transactions)
	}
	// var a []interface{}
	// a = append(a, receive.Transactions)
	existingTask.Transactions = append(receive.Transactions, existingTask.Transactions[:maxlen]...)
	err = c.userRepository.UpdateTask(id, existingTask)
	if err != nil {
		ctx.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(200, gin.H{
		"message": fmt.Sprintf("%d updated", id),
	})
}

// DeleteTaskForID method
func (c *TaskController) DeleteTaskForID(ctx *gin.Context) {
	idstr := ctx.Param("id")
	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": fmt.Sprintf("%s is not a valid number", idstr),
		})
		return
	}

	err = c.userRepository.DeleteTask(id)
	if err != nil {
		ctx.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"message": fmt.Sprintf("%d deleted", id),
	})
}
