package admin

import (
	"car-park/internal/controllers/enterprises"
	"car-park/internal/controllers/vehicles"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Controller struct {
	ctrlVehicle    *vehicles.Controller
	ctrlEnterprise *enterprises.Controller
}

func NewCtrl(v *vehicles.Controller, e *enterprises.Controller) *Controller {
	return &Controller{
		ctrlVehicle:    v,
		ctrlEnterprise: e,
	}
}

func (ctrl *Controller) ShowEnterpriseAndVehicleByManagerID(c *gin.Context) {
	c.HTML(http.StatusOK, "admin_vehicles.html", gin.H{
		// todo: add id handle
		"enterprises": ctrl.ctrlEnterprise.FetchRawByManagerID(c, 1),
		"vehicles":    ctrl.ctrlVehicle.FetchRawByManagerID(c, 1),
	})
}
