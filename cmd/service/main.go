package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	csrf "github.com/utrack/gin-csrf"

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

var creds = gin.H{
	"ismirnov": map[string]interface{}{
		"id": int64(1),
	},
	"mgreen": map[string]interface{}{
		"id": int64(2),
	},
}

func main() {
	ctx := context.Background()

	db := cmd.MustInitDB(ctx)
	defer db.Close(context.Background())

	repository := storage.New(db)
	vehicleProvider := vehicles.New(repository)
	modelsProvider := models.New(repository)
	driversProvider := drivers.New(repository)
	enterprisesProvider := enterprises.New(repository)

	server := gin.Default()

	// MIDDLEWARE
	server.Use(gin.Recovery(), gin.Logger())

	// AUTH endpoints
	authorized := server.Group("/admin", gin.BasicAuth(gin.Accounts{
		"ismirnov": "one",
		"mgreen":   "two",
	}),
		func(c *gin.Context) {
			token := c.Request.Header.Get("X-Csrf-Token")
			if token == "" || !validToken(token) {
				c.String(http.StatusBadRequest, "access denied")
				c.Abort()
				return
			}

			c.Set("csrfToken", token)
			c.String(http.StatusOK, csrf.GetToken(c))
		},
	)

	authorized.GET("enterprises", func(c *gin.Context) {
		user := c.MustGet(gin.AuthUserKey).(string)
		if userData, ok := creds[user]; ok {
			ud, ok := userData.(map[string]interface{})
			if !ok {
				return
			}

			userID, ok := ud["id"]
			if v, ok := userID.(int64); ok {
				ee := enterprisesProvider.FetchAllByManagerID(c, v)
				c.JSON(http.StatusOK, gin.H{
					"enterprises": ee,
				})
			}
		} else {
			c.JSON(http.StatusOK, gin.H{"user": user, "secret": "NO SECRET :("})
		}
	})

	authorized.GET("vehicles", func(c *gin.Context) {
		user := c.MustGet(gin.AuthUserKey).(string)
		if userData, ok := creds[user]; ok {
			ud, ok := userData.(map[string]interface{})
			if !ok {
				return
			}

			userID, ok := ud["id"]
			if v, ok := userID.(int64); ok {
				vv := vehicleProvider.FetchAllByManagerID(c, v)
				c.JSON(http.StatusOK, gin.H{
					"vehicles": vv,
				})
			}
		} else {
			c.JSON(http.StatusOK, gin.H{"user": user, "secret": "NO SECRET :("})
		}
	})

	// OTHER endpoints
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
	apiVehicleAdmin.POST("/add",
		gin.BasicAuth(gin.Accounts{
			"ismirnov": "one",
			"mgreen":   "two",
		}),
		checkCSRF(),
		vehicleProvider.Create)

	apiVehicleAdmin.PUT("/update",
		gin.BasicAuth(gin.Accounts{
			"ismirnov": "one",
			"mgreen":   "two",
		}),
		checkCSRF(),
		vehicleProvider.Update)

	apiVehicleAdmin.DELETE("/delete",
		gin.BasicAuth(gin.Accounts{
			"ismirnov": "one",
			"mgreen":   "two",
		}),
		checkCSRF(),
		vehicleProvider.Delete)

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

	server.Run("127.0.0.1:2023")
}

func checkCSRF() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("X-Csrf-Token")
		if token == "" || !validToken(token) {
			c.String(http.StatusBadRequest, "access denied")
			c.Abort()
			return
		}
	}
}

func validToken(token string) bool {
	// todo: implement
	return true
}
