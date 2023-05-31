package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"car-park/cmd"
	"car-park/internal/controllers/models"
	"car-park/internal/controllers/vehicles"
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
	vehicleProvider := vehicles.New(repository)
	modelsProvider := models.New(repository)

	server := gin.New()

	// MIDDLEWARE
	server.Use(gin.Recovery(), gin.Logger())

	apiVehicle := server.Group("/api/v1/vehicles")
	apiVehicle.GET("", func(c *gin.Context) {
		c.JSONP(http.StatusOK, gin.H{
			"vehicles": vehicleProvider.FetchAll(ctx),
		})
	})

	apiVehicleAdmin := apiVehicle.Group("/admin")
	apiVehicleAdmin.GET("", func(c *gin.Context) {
		c.HTML(http.StatusOK, "vehicles_admin.html",
			gin.H{
				"models": modelsProvider.FetchAll(ctx),
			})
	})
	apiVehicleAdmin.POST("/add", vehicleProvider.Create)
	apiVehicleAdmin.POST("/update", vehicleProvider.Update)
	apiVehicleAdmin.POST("/delete", vehicleProvider.Delete)

	server.GET("/view/vehicle/redirect", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/api/v1/vehicles/admin/")
	})

	server.LoadHTMLGlob("templates/views/*")
	server.GET("/view/vehicles", func(c *gin.Context) {
		c.HTML(http.StatusOK, "vehicles.html", gin.H{
			"title":    "Vehicles",
			"vehicles": vehicleProvider.FetchAll(ctx),
		})
	})
	server.GET("/view/models", func(c *gin.Context) {
		c.HTML(http.StatusOK, "models.html", gin.H{
			"title":  "Models",
			"models": modelsProvider.FetchAll(ctx),
		})
	})

	server.Run()
}
