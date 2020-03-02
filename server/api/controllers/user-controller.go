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

// UserController struct
type UserController struct {
	userRepository *repositories.UserRepository
}

// Credentials struct
type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

// Init method
func (c *UserController) Init(db *sql.DB) {
	c.userRepository = &repositories.UserRepository{}
	c.userRepository.Init(db)
}

// CreateUser method
func (c *UserController) CreateUser(ctx *gin.Context) {
	var user models.User
	er1 := ctx.BindJSON(&user)
	if er1 != nil {
		ctx.JSON(400, gin.H{
			"error": "number should not be string",
		})
		return
	}
	if user.Name == "" {
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
	// user.Number = userid
	createdUser, err := c.userRepository.CreateUser(user)
	if err != nil {
		log.Printf("Error: %v\n", err)
		ctx.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(201, gin.H{
		"user": createdUser,
	})
	return
}

// GetUsers method
func (c *UserController) GetUsers(ctx *gin.Context) {
	users := []models.User{}
	var err error
	users, err = c.userRepository.GetAllUsers()
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"users": users,
	})
}

// GetUserByID method
func (c *UserController) GetUserByID(ctx *gin.Context) {
	idstr := ctx.Param("id")
	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": fmt.Sprintf("%s is not a valid number", idstr),
		})
		return
	}

	user, err := c.userRepository.GetUserByID(id)
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"user": user,
	})
}

// GetUserByNameAndNumber method
func (c *UserController) GetUserByNameAndNumber(ctx *gin.Context) {
	idstr := ctx.Param("id")
	// number := idstr.(string)

	useridi := ctx.Param("id1")
	// name := useridi.(string)

	user, err := c.userRepository.GetUserByNameAndNumber(idstr, useridi)
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"user": user,
	})
}

// UpdateUserForID method
func (c *UserController) UpdateUserForID(ctx *gin.Context) {
	idstr := ctx.Param("id")
	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": fmt.Sprintf("%s is not a valid number", idstr),
		})
		return
	}

	existingUser, err := c.userRepository.GetUserByID(id)
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}

	var receive struct {
		Amount       float64       `json:"amount"`
		Transactions []interface{} `json:"transactions"`
	}
	ctx.BindJSON(&receive)
	existingUser.Amount = receive.Amount
	maxlen := 10
	if len(existingUser.Transactions) < 10 {
		maxlen = len(existingUser.Transactions)
	}
	// var a []interface{}
	// a = append(a, receive.Transactions)
	existingUser.Transactions = append(receive.Transactions, existingUser.Transactions[:maxlen]...)
	err = c.userRepository.UpdateUser(id, existingUser)
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

// DeleteUserForID method
func (c *UserController) DeleteUserForID(ctx *gin.Context) {
	idstr := ctx.Param("id")
	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": fmt.Sprintf("%s is not a valid number", idstr),
		})
		return
	}

	err = c.userRepository.DeleteUser(id)
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

func (c *UserController) Login(ctx *gin.Context) {
	var users = map[string]string{
		"supersai":   "lakshmisat",
		"saitejan11": "saiTEja",
	}
	role := "user"

	var user Credentials
	ctx.BindJSON(&user)

	_, ok := users[user.Username]
	if ok {
		if users[user.Username] == user.Password {
			role = "admin"
			ctx.JSON(200, gin.H{
				"token": "seitestokecratejan",
				"role":  role,
			})
			return
		}
	}

	ussr, err := c.userRepository.GetUserByNameAndNumber(user.Username, user.Password)
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": "Wrong Credentials",
		})
		return
	}

	ctx.JSON(200, gin.H{
		"token": ussr.Name + ",?," + ussr.Number,
		"role":  role,
	})
}
