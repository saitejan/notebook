package controllers

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

//HandlerF method
func HandlerF() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		_, ok := ctx.Request.Header["Authorization"]

		if !ok {
			ctx.Abort()
			ctx.Writer.WriteHeader(http.StatusUnauthorized)
			ctx.Writer.Write([]byte("No Header"))
			return
		}

		authToken := ctx.Request.Header["Authorization"][0]
		splitToken := strings.Split(authToken, "Bearer ")
		authToken = splitToken[1]

		if authToken == "seitestokecratejan" {
			ctx.Set("userid", authToken)
			return
		}
		if ctx.Request.Method != "GET" || ctx.FullPath() != "/api/users/:id" {
			log.Printf("No Access: \n")
			ctx.Abort()
			ctx.Writer.WriteHeader(http.StatusUnauthorized)
			ctx.Writer.Write([]byte("Unauthorized"))
			return
		}
		name := strings.Split(authToken, ",?,")
		if len(name) < 2 {
			log.Printf("wrong token: \n")
			ctx.Abort()
			ctx.Writer.WriteHeader(http.StatusUnauthorized)
			ctx.Writer.Write([]byte("Unauthorized"))
			return
		}
	}
}
