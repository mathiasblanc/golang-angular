package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"

	auth0 "github.com/auth0-community/go-auth0"
	"github.com/gin-gonic/gin"
	"github.com/mathiasblanc/golang-angular/handlers"
	jose "gopkg.in/square/go-jose.v2"
)

var (
	audience string
	domain   string
)

func main() {
	setAuth0Variables()
	web := gin.Default()

	web.NoRoute(func(ctx *gin.Context) {
		dir, file := path.Split(ctx.Request.RequestURI)
		ext := filepath.Ext(file)

		if file == "" || ext == "" {
			ctx.File("./ui/dist/ui/index.html")
		} else {
			ctx.File("./ui/dist/ui/" + path.Join(dir, file))
		}
	})

	authorized := web.Group("/")
	authorized.Use(authRequired())

	authorized.GET("/todo", handlers.GetTodoListHandler)
	authorized.POST("/todo", handlers.AddTodoHandler)
	authorized.DELETE("/todo/:id", handlers.DeleteTodoHandler)
	authorized.PUT("/todo", handlers.CompleteTodoHandler)

	err := web.Run(":3000")

	if err != nil {
		panic(err)
	}
}

func setAuth0Variables() {
	audience = os.Getenv("AUTH0_API_IDENTIFIER")
	domain = os.Getenv("AUTH0_DOMAIN")
}

func authRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		var auth0Domain = "https://" + domain + "/"
		client := auth0.NewJWKClient(auth0.JWKClientOptions{URI: auth0Domain + ".well-known/jwks.json"}, nil)
		configuration := auth0.NewConfiguration(client, []string{audience}, auth0Domain, jose.RS256)
		validator := auth0.NewValidator(configuration, nil)

		_, err := validator.ValidateRequest(c.Request)

		if err != nil {
			fmt.Println(err)

			log.Println(err)
			terminateWithError(http.StatusUnauthorized, "token is not valid", c)
			return
		}
		c.Next()

	}
}

func terminateWithError(statusCode int, message string, c *gin.Context) {
	c.JSON(statusCode, gin.H{"error": message})
	c.Abort()
}
