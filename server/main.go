package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	"notebook/api/controllers"
	"notebook/api/util"
)

func main() {
	//gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	db, err := util.GetDB()
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	userController := controllers.UserController{}
	userController.Init(db)

	router.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "OPTIONS", "PUT"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "User-Agent", "Referrer", "Host", "authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowAllOrigins:  false,
		AllowOriginFunc: func(origin string) bool {
			if origin == "http://localhost:3000" || origin == "https://kolluru.herokuapp.com" || origin == "https://sathvika.netlify.com" {
				return true
			}
			return false
		},
		MaxAge: 86400,
	}))

	router.Use(func(ctx *gin.Context) {
		if !util.Contains([]string{"POST", "PUT", "PATCH"}, ctx.Request.Method) {
			return
		}

		if ctx.Request.Header["Content-Length"][0] == "0" {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Payload should not be empty"})
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}

		if len(ctx.Request.Header["Content-Type"]) == 0 ||
			!util.Contains(ctx.Request.Header["Content-Type"], "application/json") {
			ctx.JSON(http.StatusUnsupportedMediaType, gin.H{"message": "Content type should be application/json"})
			ctx.AbortWithStatus(http.StatusUnsupportedMediaType)
			return
		}
	})

	router.Use(static.Serve("/", static.LocalFile("./web", true)))
	router.POST("/api/login", userController.Login)
	router.POST("/api/users", controllers.HandlerF(), userController.CreateUser)
	router.GET("/api/users", controllers.HandlerF(), userController.GetUsers)
	router.GET("/api/users/:id", controllers.HandlerF(), userController.GetUserByID)
	router.GET("/api/users/:id/:id1", controllers.HandlerF(), userController.GetUserByNameAndNumber)
	router.PUT("/api/users/:id", controllers.HandlerF(), userController.UpdateUserForID)
	router.DELETE("/api/users/:id", controllers.HandlerF(), userController.DeleteUserForID)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}
	router.Run(":" + port)
}
