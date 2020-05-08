package main

import (
	"github.com/ajdelgados/golang-graphql/controllers"
	"github.com/ajdelgados/golang-graphql/models"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	db := models.SetupModels()
	defer db.Close()

	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	r.GET("/graphql", controllers.Graphql)

	r.Run()
}
