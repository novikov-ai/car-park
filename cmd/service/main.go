package main

import (
	"car-park/cmd"
	"car-park/internal/auth"
	"car-park/internal/controllers/drivers"
	"car-park/internal/controllers/enterprises"
	"car-park/internal/controllers/vehicles"
	"car-park/internal/geo"
	"car-park/internal/models"
	"car-park/internal/storage"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"strconv"
	"time"

	_ "net/http/pprof"
)

const (
	//apiPath = "/api/v1"

	iso8601 = "2006-01-02T15:04:05Z07:00"
)

func main() {
	ctx := context.Background()

	db := cmd.MustInitDB(ctx)
	defer db.Close(context.Background())

	repository := storage.New(db)

	//cache := cmd.MustInitCache(ctx)
	geoClient := geo.NewGeoClient(repository)

	enterpriseProvider := enterprises.New(repository)
	vehicleProvider := vehicles.New(repository)
	driverProvider := drivers.New(repository)

	server := gin.Default()

	// MIDDLEWARE
	server.Use(gin.Recovery(), gin.Logger())

	server.LoadHTMLGlob("templates/*.html")

	server.GET("vehicles", func(c *gin.Context) {
		c.HTML(http.StatusOK, "vehicles.html", gin.H{
			"vehicles": vehicleProvider.FetchAll(c),
		})
	})

	// API
	groupAPI := server.Group("api")

	groupAPI.GET("/vehicles", func(c *gin.Context) {
		c.JSON(http.StatusOK, vehicleProvider.FetchAll(c))
	})

	groupAPI.GET("/drivers", func(c *gin.Context) {
		c.JSON(http.StatusOK, driverProvider.FetchAll(c))
	})

	groupAPI.GET("/enterprises", func(c *gin.Context) {
		c.JSON(http.StatusOK, enterpriseProvider.FetchAll(c))
	})

	groupAPI.GET("/enterprises/report", func(c *gin.Context) {
		c.Request.ParseForm()

		id, start, end := parseQueryParams(c)

		c.JSON(http.StatusOK, map[string]interface{}{
			"mileage": enterpriseProvider.SumMileageByVehicle(c, id, start, end),
		})
	})

	// API - CRUD
	groupVehicleAPI := groupAPI.Group("/vehicle")
	groupVehicleAPI.POST("/new", func(c *gin.Context) {
		insertedID := vehicleProvider.Create(c)
		if insertedID == 0 {
			c.HTML(http.StatusBadRequest, "bad_request.html", gin.H{})
			c.Abort()
			return
		}

		c.HTML(http.StatusOK, "success.html", gin.H{})
	})
	groupVehicleAPI.POST("/update", func(c *gin.Context) {
		err := vehicleProvider.Update(c)
		if err != nil {
			c.HTML(http.StatusBadRequest, "bad_request.html", gin.H{})
			c.Abort()
			return
		}

		c.HTML(http.StatusOK, "success.html", gin.H{})
	})
	groupVehicleAPI.POST("/delete", func(c *gin.Context) {
		err := vehicleProvider.Delete(c)

		if err != nil {
			c.HTML(http.StatusBadRequest, "bad_request.html", gin.H{})
			c.Abort()
			return
		}

		c.HTML(http.StatusOK, "success.html", gin.H{})
	})

	server.GET("admin/vehicles", func(c *gin.Context) {
		c.HTML(http.StatusOK, "admin_vehicles.html", gin.H{
			"enterprises": enterpriseProvider.FetchAllByManagerID(c, 1),
			"vehicles":    vehicleProvider.FetchAllByManagerID(c, 1),
		})
	})

	// API - GPS
	groupGpsAPI := groupAPI.Group("/gps")
	groupGpsAPI.GET("/track", func(c *gin.Context) {
		_, start, end := parseQueryParamsIdStartEndUnix(c)

		trackParam := c.Request.FormValue("id")
		trackID, err := strconv.Atoi(trackParam)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		var track interface{}
		trackPoints := geoClient.GetTrackByTrip(ctx, int64(trackID), float64(start), float64(end))

		if c.Request.URL.Query().Has("geoJSON") {
			track = geoClient.ToGeoJSON(trackPoints)
		} else {
			track = trackPoints
		}

		c.JSON(http.StatusOK, gin.H{
			"track": track,
		})
	})

	groupGpsAPI.GET("/trip", func(c *gin.Context) {
		vehicle, start, end := parseQueryParamsIdStartEndDates(c)
		startValue, endValue := "", ""
		if start != nil {
			startValue = start.Format(iso8601)
		}
		if end != nil {
			endValue = end.Format(iso8601)
		}

		trips, err := geoClient.GetTripsByVehicle(ctx, vehicle, startValue, endValue)
		if err != nil {
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"trips": trips,
		})
	})

	// API - REPORTS
	groupReportsAPI := groupAPI.Group("/reports")
	groupReportsAPI.GET("", func(c *gin.Context) {
		vehicle, start, end := parseQueryParamsIdStartEndUnix(c)

		reports := vehicleProvider.GetVehicleReportDaily(ctx, vehicle, start, end)

		c.JSON(http.StatusOK, gin.H{
			"report": reports,
		})
	})

	// REPORTS

	groupReports := server.Group("reports")
	groupReports.GET("", func(c *gin.Context) {
		c.HTML(http.StatusOK, "reports.html", gin.H{
			"vehicles": vehicleProvider.FetchAll(c),
		})
	})

	server.GET("reports/vehicle", func(c *gin.Context) {
		err := c.Request.ParseForm()
		if err != nil {
			c.Status(http.StatusBadRequest)
			c.Abort()
			return
		}

		//vehicle := c.Request.Form.Get("vehicle")
		startDate := c.Request.Form.Get("start")
		endDate := c.Request.Form.Get("end")
		//period := c.Request.Form.Get("period")
		//report := c.Request.Form.Get("report")

		vehicleID, start, end := getVehicleStartEndTimeUnix(c)

		var vr []models.VehicleReport

		switch c.Request.Form.Get("period") {
		case "daily":
			fallthrough
		case "monthly":
			fallthrough
		case "yearly":
			vr = vehicleProvider.GetVehicleReportDaily(ctx, vehicleID, start, end)
		default:
			vr = vehicleProvider.GetVehicleReportDaily(ctx, vehicleID, start, end)
		}

		c.HTML(http.StatusOK, "reports_detailed.html", gin.H{
			"title": fmt.Sprintf("Пробег автомобиля за период (%s - %s)", startDate, endDate),
			"data":  vr,
		})
	})

	server.GET("report/vehicle/monthly", func(c *gin.Context) {
		type daily struct {
			Date    string
			Mileage int
		}

		c.HTML(http.StatusOK, "reports_monthly.html", gin.H{
			"title": "Vehicle #1 - Daily Report",
			"data": []daily{
				{
					Date:    "2023-05-05",
					Mileage: 150,
				},
				{
					Date:    "2023-05-06",
					Mileage: 175,
				},
				{
					Date:    "2023-05-07",
					Mileage: 200,
				},
				{
					Date:    "2023-05-08",
					Mileage: 225,
				},
				{
					Date:    "2023-05-09",
					Mileage: 250,
				},
			}})
	})

	server.GET("report/vehicle/annual", func(c *gin.Context) {
		type daily struct {
			Date    string
			Mileage int
		}

		c.HTML(http.StatusOK, "reports_annual.html", gin.H{
			"title": "Vehicle #1 - Daily Report",
			"data": []daily{
				{
					Date:    "2023-05-05",
					Mileage: 150,
				},
				{
					Date:    "2023-05-06",
					Mileage: 175,
				},
				{
					Date:    "2023-05-07",
					Mileage: 200,
				},
				{
					Date:    "2023-05-08",
					Mileage: 225,
				},
				{
					Date:    "2023-05-09",
					Mileage: 250,
				},
			}})
	})

	// LOGIN
	groupLogin := server.Group("/login")
	groupLogin.GET("", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{})
	})
	groupLogin.GET("/auth", func(c *gin.Context) {
		err := c.Request.ParseForm()
		if err != nil {
			c.String(http.StatusBadRequest, "bad request")
			c.Abort()
			return
		}

		userID := auth.GetUserIdWhenAuth(c)
		if userID == 0 {
			c.Redirect(http.StatusPermanentRedirect, "/login")
			c.Abort()
			return
		}

		c.Redirect(http.StatusMovedPermanently, "../enterprises?manager="+strconv.Itoa(int(userID)))
	})

	// VIEW

	groupManager := server.Group("/manager")
	groupManager.GET("/vehicles", func(c *gin.Context) {
		entID := c.Request.FormValue("enterprise")
		id, err := strconv.ParseInt(entID, 10, 64)
		if err != nil {
			return
		}

		c.HTML(http.StatusOK, "manager_vehicles.html", gin.H{
			"title":    "Vehicles",
			"vehicles": vehicleProvider.FetchAllByEnterpriseID(c, id),
		})
	})

	server.GET("enterprises", func(c *gin.Context) {
		id := c.Query("manager")
		v, err := strconv.Atoi(id)
		if err != nil {
			c.String(http.StatusBadRequest, "bad request")
			c.Abort()
			return
		}

		ents := enterpriseProvider.FetchAllByManagerID(c, int64(v))

		if len(ents) == 0 {
			c.String(http.StatusBadRequest, "bad request")
			c.Abort()
			return
		}

		c.HTML(http.StatusOK, "manager_enterprises.html", gin.H{
			"title":       "Enterprises",
			"enterprises": ents,
		})
	})

	server.GET("/gps/trip", func(c *gin.Context) {
		vehicle, start, end := parseQueryParamsIdStartEndDates(c)
		startValue, endValue := "", ""
		if start != nil {
			startValue = start.Format(iso8601)
		}
		if end != nil {
			endValue = end.Format(iso8601)
		}

		trips, err := geoClient.GetTripsByVehicle(ctx, vehicle, startValue, endValue)
		if err != nil {
			return
		}

		c.HTML(http.StatusOK, "vehicle_trips.html", gin.H{
			"title": "Trips",
			"trips": trips,
		})
	})

	server.Run()
}

func removeQueryParam(query, param string) string {
	u, err := url.Parse(query)
	if err != nil {
		fmt.Println("error parsing URL: ", err)
		return query
	}

	values := u.Query()
	values.Del(param)

	u.RawQuery = values.Encode()

	return u.String()
}

func parseQueryTimeParam(value string) *time.Time {
	parsed, err := time.Parse(time.RFC3339, fmt.Sprintf("%sT00:00:00.000Z", value))
	if err != nil {
		return nil
	}

	return &parsed
}

func getVehicleStartEndTimeUnix(c *gin.Context) (int64, int64, int64) {
	query := c.Request.URL.Query()
	id := query.Get("vehicle")

	idValue, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		idValue = 0
	}

	start := query.Get("start")
	startUnix := parseQueryTimeParam(start).Unix()

	startLocalTime := time.Unix(startUnix, 0)
	startUTC := startLocalTime.UTC()

	end := query.Get("end")
	endUnix := parseQueryTimeParam(end).Unix()

	endLocalTime := time.Unix(endUnix, 0)
	endUTC := endLocalTime.UTC()

	if err != nil || end == "" {
		endUTC = time.Now().UTC()
	}

	return idValue, startUTC.Unix(), endUTC.Unix()
}

func parseQueryParamsIdStartEndUnix(c *gin.Context) (int64, int64, int64) {
	query := c.Request.URL.Query()
	id := query.Get("vehicle")

	idValue, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		idValue = 0
	}

	start := query.Get("start")
	startUnix, err := strconv.ParseInt(start, 10, 64)
	if err != nil {
		startUnix = 0
	}

	startLocalTime := time.UnixMilli(startUnix)
	startUTC := startLocalTime.UTC()

	end := query.Get("end")
	endUnix, err := strconv.ParseInt(end, 10, 64)

	endLocalTime := time.UnixMilli(endUnix)
	endUTC := endLocalTime.UTC()

	if err != nil || end == "" {
		endUTC = time.Now().UTC()
	}

	return idValue, startUTC.Unix(), endUTC.Unix()
}

func parseQueryParamsIdStartEndDates(c *gin.Context) (int64, *time.Time, *time.Time) {
	query := c.Request.URL.Query()
	id := query.Get("vehicle")

	idValue, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		idValue = 0
	}

	var startValue *time.Time
	start := query.Get("start")
	if start == "" {
		zeroUnix := time.Unix(0, 0)
		startValue = &zeroUnix
	} else {
		startValue = parseQueryTimeParam(start)
	}

	var endValue *time.Time
	end := query.Get("end")
	if end == "" {
		timeNow := time.Now()
		endValue = &timeNow
	} else {
		endValue = parseQueryTimeParam(end)
	}

	return idValue, startValue, endValue
}

func parseQueryParams(c *gin.Context) (int64, *time.Time, *time.Time) {
	query := c.Request.URL.Query()
	id := query.Get("id")

	idValue, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		idValue = 0
	}

	var startValue *time.Time
	start := query.Get("start")
	if start == "" {
		zeroUnix := time.Unix(0, 0)
		startValue = &zeroUnix
	} else {
		startValue = parseQueryTimeParam(start)
	}

	var endValue *time.Time
	end := query.Get("end")
	if end == "" {
		timeNow := time.Now()
		endValue = &timeNow
	} else {
		endValue = parseQueryTimeParam(end)
	}

	return idValue, startValue, endValue
}
