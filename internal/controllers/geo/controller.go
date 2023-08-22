package geo

import (
	"car-park/internal/controllers/tools/query"
	"car-park/internal/storage"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

const iso8601 = "2006-01-02T15:04:05Z07:00"

type Controller struct {
	client Client
}

func NewCtrl(storage storage.Client) *Controller {
	return &Controller{
		client: NewGeoClient(storage),
	}
}

func (ctrl *Controller) ShowAllTracksJSON(c *gin.Context) {
	_, start, end := query.ParseQueryParamsIdStartEndUnix(c)

	trackParam := c.Request.FormValue("id")
	trackID, err := strconv.Atoi(trackParam)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	var track interface{}
	trackPoints := ctrl.client.GetTrackByTrip(c, int64(trackID), float64(start), float64(end))

	if c.Request.URL.Query().Has("geoJSON") {
		track = ctrl.client.ToGeoJSON(trackPoints)
	} else {
		track = trackPoints
	}

	c.JSON(http.StatusOK, gin.H{
		"track": track,
	})
}

func (ctrl *Controller) ShowAllTripsJSON(c *gin.Context) {
	vehicle, start, end := query.ParseQueryParamsIdStartEndDates(c)
	startValue, endValue := "", ""
	if start != nil {
		startValue = start.Format(iso8601)
	}
	if end != nil {
		endValue = end.Format(iso8601)
	}

	trips, err := ctrl.client.GetTripsByVehicle(c, vehicle, startValue, endValue)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"trips": trips,
	})
}

func (ctrl *Controller) ShowTripsByVehicle(c *gin.Context) {
	vehicle, start, end := query.ParseQueryParamsIdStartEndDates(c)
	startValue, endValue := "", ""
	if start != nil {
		startValue = start.Format(iso8601)
	}
	if end != nil {
		endValue = end.Format(iso8601)
	}

	trips, err := ctrl.client.GetTripsByVehicle(c, vehicle, startValue, endValue)
	if err != nil {
		return
	}

	c.HTML(http.StatusOK, "vehicle_trips.html", gin.H{
		"title": "Trips",
		"trips": trips,
	})
}
