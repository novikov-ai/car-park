package main

import (
	"car-park/cmd"
	"car-park/internal/controllers"
	"car-park/internal/controllers/auth"
	"car-park/internal/storage"
	"context"
	"github.com/gin-gonic/gin"
	_ "net/http/pprof"
)

func main() {
	ctx := context.Background()

	db := cmd.MustInitDB(ctx)
	defer db.Close(context.Background())

	repository := storage.New(db)
	ctrls := controllers.New(repository)

	server := gin.Default()

	setupRoutes(server, ctrls)

	server.Run()
}

func setupRoutes(server *gin.Engine, ctrl *controllers.Controllers) {
	server.Use(gin.Recovery(), gin.Logger())

	server.LoadHTMLGlob("templates/*.html")

	server.GET("vehicles", ctrl.Vehicles.ShowAll)

	// API
	groupAPI := server.Group("api")
	groupAPI.GET("/vehicles", ctrl.Vehicles.ShowAllJSON)
	groupAPI.GET("/drivers", ctrl.Drivers.ShowAllJSON)
	groupAPI.GET("/enterprises", ctrl.Enterprises.ShowAllJSON)
	groupAPI.GET("/enterprises/report", ctrl.Enterprises.SumMileageJSON)

	groupVehicleAPI := groupAPI.Group("/vehicle")
	groupVehicleAPI.POST("/new", ctrl.Vehicles.Create)
	groupVehicleAPI.POST("/update", ctrl.Vehicles.Update)
	groupVehicleAPI.POST("/delete", ctrl.Vehicles.Delete)

	server.GET("admin/vehicles", ctrl.Admin.ShowEnterpriseAndVehicleByManagerID)

	// API - GPS
	groupGpsAPI := groupAPI.Group("/gps")
	groupGpsAPI.GET("/track", ctrl.Geo.ShowAllTracksJSON)
	groupGpsAPI.GET("/trip", ctrl.Geo.ShowAllTripsJSON)

	// API - REPORTS
	groupReportsAPI := groupAPI.Group("/reports")
	groupReportsAPI.GET("", ctrl.Vehicles.GetReportJSON)

	// REPORTS
	groupReports := server.Group("reports")
	groupReports.GET("", ctrl.Vehicles.GetReport)

	server.GET("reports/vehicle", ctrl.Vehicles.GetReportByVehicle)
	server.GET("report/vehicle/monthly", ctrl.Vehicles.GetReportByVehicleMonthly)
	server.GET("report/vehicle/annual", ctrl.Vehicles.GetReportByVehicleAnnual)

	// LOGIN
	groupLogin := server.Group("/login")
	groupLogin.GET("", auth.ShowLoginPage)
	groupLogin.GET("/auth", auth.Authorize)

	// VIEW
	groupManager := server.Group("/manager")
	groupManager.GET("/vehicles", ctrl.Vehicles.ShowAllByEnterpriseID)

	server.GET("enterprises", ctrl.Enterprises.ShowAllByManagerID)

	server.GET("/gps/trip", ctrl.Geo.ShowTripsByVehicle)
}
