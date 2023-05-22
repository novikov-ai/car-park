package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"car-park/cmd"
	"car-park/internal/storage"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to load environment-file: %v\n", err)
		os.Exit(1)
	}
}

func main() {
	ctx := context.Background()

	db := cmd.MustInitDB(ctx)
	defer db.Close(context.Background())

	repository := storage.New(db)
	vehicles := repository.FetchAll(ctx)

	server := gin.New()

	// MIDDLEWARE
	server.Use(gin.Recovery(), gin.Logger())

	apiRoutes := server.Group("/api")
	{
		apiRoutes.GET("/vehicles", func(c *gin.Context) {
			c.JSONP(http.StatusOK, gin.H{
				"vehicles": vehicles,
			})
		})
	}

	server.LoadHTMLGlob("templates/views/*")
	viewRoutes := server.Group("/view")
	{
		viewRoutes.GET("/vehicles", func(c *gin.Context) {
			c.HTML(http.StatusOK, "vehicles.html", gin.H{
				"title":    "Vehicles",
				"vehicles": repository.FetchAll(ctx),
			})
		})
		viewRoutes.GET("/models", func(c *gin.Context) {
			c.HTML(http.StatusOK, "models.html", gin.H{
				"title":  "Models",
				"models": repository.FetchAllModels(ctx),
			})
		})
	}

	server.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
