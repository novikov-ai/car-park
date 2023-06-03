package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"car-park/cmd"
	"car-park/internal/controllers/drivers"
	"car-park/internal/controllers/enterprises"
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

const apiPath = "/api/v1"

func main() {
	ctx := context.Background()

	db := cmd.MustInitDB(ctx)
	defer db.Close(context.Background())

	repository := storage.New(db)
	vehicleProvider := vehicles.New(repository)
	modelsProvider := models.New(repository)
	driversProvider := drivers.New(repository)
	enterprisesProvider := enterprises.New(repository)

	server := gin.New()

	// MIDDLEWARE
	server.Use(gin.Recovery(), gin.Logger())

	apiEnterprises := server.Group(apiPath + "/enterprises")
	apiEnterprises.GET("", func(c *gin.Context) {
		c.JSONP(http.StatusOK, gin.H{
			"enterprises": enterprisesProvider.FetchAll(c),
		})
	})

	apiDrivers := server.Group(apiPath + "/drivers")
	apiDrivers.GET("", func(c *gin.Context) {
		c.JSONP(http.StatusOK, gin.H{
			"drivers": driversProvider.FetchAll(c),
		})
	})
	server.GET(apiPath+"/drivers-vehicles", func(c *gin.Context) {
		c.JSONP(http.StatusOK, gin.H{
			"drivers-vehicles": driversProvider.FetchAllVehicleDrivers(c),
		})
	})

	apiVehicles := server.Group(apiPath + "/vehicles")
	apiVehicles.GET("", func(c *gin.Context) {
		c.JSONP(http.StatusOK, gin.H{
			"vehicles": vehicleProvider.FetchAll(c),
		})
	})

	apiModels := server.Group(apiPath + "/models")
	apiModels.GET("", func(c *gin.Context) {
		c.JSONP(http.StatusOK, gin.H{
			"models": modelsProvider.FetchAll(c),
		})
	})

	apiVehicleAdmin := apiVehicles.Group("/admin")
	apiVehicleAdmin.GET("", func(c *gin.Context) {
		c.HTML(http.StatusOK, "vehicles_admin.html",
			gin.H{
				"models": modelsProvider.FetchAll(c),
			})
	})
	apiVehicleAdmin.POST("/add", vehicleProvider.Create)
	apiVehicleAdmin.POST("/update", vehicleProvider.Update)
	apiVehicleAdmin.POST("/delete", vehicleProvider.Delete)

	server.GET("/view/vehicle/redirect", func(c *gin.Context) {
		c.Redirect(http.StatusFound, apiPath+"/vehicles/admin/")
	})

	server.LoadHTMLGlob("templates/views/*")
	server.GET("/view/vehicles", func(c *gin.Context) {
		c.HTML(http.StatusOK, "vehicles.html", gin.H{
			"title":    "Vehicles",
			"vehicles": vehicleProvider.FetchAll(c),
		})
	})
	server.GET("/view/models", func(c *gin.Context) {
		c.HTML(http.StatusOK, "models.html", gin.H{
			"title":  "Models",
			"models": modelsProvider.FetchAll(c),
		})
	})

	server.Run()
}
