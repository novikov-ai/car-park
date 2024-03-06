package vehicles

import (
	"car-park/internal/controllers/tools/qparser"
	"car-park/internal/models"
	"car-park/internal/storage"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	provider *Provider
}

func NewCtrl(storage storage.Client) *Controller {
	return &Controller{
		provider: New(storage),
	}
}

func (ctrl *Controller) ShowAll(c *gin.Context) {
	c.HTML(http.StatusOK, "vehicles.html", gin.H{
		"vehicles": ctrl.provider.FetchAll(c),
	})
}

func (ctrl *Controller) ShowAllJSON(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"vehicles": ctrl.provider.FetchAll(c),
	})
}

func (ctrl *Controller) Create(c *gin.Context) {
	insertedID := ctrl.provider.Create(c)
	if insertedID == 0 {
		c.HTML(http.StatusBadRequest, "bad_request.html", gin.H{})
		c.Abort()
		return
	}

	c.HTML(http.StatusOK, "success.html", gin.H{})
}

func (ctrl *Controller) Update(c *gin.Context) {
	err := ctrl.provider.Update(c)
	if err != nil {
		c.HTML(http.StatusBadRequest, "bad_request.html", gin.H{})
		c.Abort()
		return
	}

	c.HTML(http.StatusOK, "success.html", gin.H{})
}

func (ctrl *Controller) Delete(c *gin.Context) {
	err := ctrl.provider.Delete(c)

	if err != nil {
		c.HTML(http.StatusBadRequest, "bad_request.html", gin.H{})
		c.Abort()
		return
	}

	c.HTML(http.StatusOK, "success.html", gin.H{})
}

func (ctrl *Controller) GetReportJSON(c *gin.Context) {
	queryParser, err := qparser.New(c.Request.URL)
	if err != nil {
		return
	}

	vehicle := queryParser.GetVehicle()
	start, end := queryParser.GetStartTimeUnix(), queryParser.GetEndTimeUnix()

	reports := ctrl.provider.GetVehicleReportDaily(c, vehicle, start, end)

	c.JSON(http.StatusOK, gin.H{
		"report": reports,
	})
}

func (ctrl *Controller) GetReport(c *gin.Context) {
	c.HTML(http.StatusOK, "reports.html", gin.H{
		"vehicles": ctrl.provider.FetchAll(c),
	})
}

func (ctrl *Controller) GetReportByVehicle(c *gin.Context) {
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

	queryParser, err := qparser.New(c.Request.URL)
	if err != nil {
		return
	}

	vehicleID := queryParser.GetVehicle()
	start, end := queryParser.GetStartTimeUnix(), queryParser.GetEndTimeUnix()

	var vr []models.VehicleReport

	switch c.Request.Form.Get("period") {
	case "daily":
		fallthrough
	case "monthly":
		fallthrough
	case "yearly":
		vr = ctrl.provider.GetVehicleReportDaily(c, vehicleID, start, end)
	default:
		vr = ctrl.provider.GetVehicleReportDaily(c, vehicleID, start, end)
	}

	c.HTML(http.StatusOK, "reports_detailed.html", gin.H{
		"title": fmt.Sprintf("Пробег автомобиля за период (%s - %s)", startDate, endDate),
		"data":  vr,
	})
}

func (ctrl *Controller) GetReportByVehicleMonthly(c *gin.Context) {
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
}

func (ctrl *Controller) GetReportByVehicleAnnual(c *gin.Context) {
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
}

func (ctrl *Controller) ShowAllByEnterpriseID(c *gin.Context) {
	entID := c.Request.FormValue("enterprise")
	id, err := strconv.ParseInt(entID, 10, 64)
	if err != nil {
		return
	}

	c.HTML(http.StatusOK, "manager_vehicles.html", gin.H{
		"title":    "Vehicles",
		"vehicles": ctrl.provider.FetchAllByEnterpriseID(c, id),
	})
}

func (ctrl *Controller) FetchRawByManagerID(c *gin.Context, id int64) []models.Vehicle {
	return ctrl.provider.FetchAllByManagerID(c, id)
}

func (ctrl *Controller) SumMileageJSON(c *gin.Context) {
	c.Request.ParseForm()

	queryParser, err := qparser.New(c.Request.URL)
	if err != nil {
		return
	}

	id := queryParser.GetVehicle()
	start, end := queryParser.GetStartTime(), queryParser.GetEndTime()

	c.JSON(http.StatusOK, map[string]interface{}{
		"mileage": ctrl.provider.SumMileageByVehicle(c, id, start, end),
	})
}
