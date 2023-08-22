package enterprises

import (
	"car-park/internal/controllers/tools/query"
	"car-park/internal/models"
	"car-park/internal/storage"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type Controller struct {
	provider *Provider
}

func NewCtrl(storage storage.Client) *Controller {
	return &Controller{
		provider: New(storage),
	}
}

func (ctrl *Controller) ShowAllJSON(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"enterprises": ctrl.provider.FetchAll(c),
	})
}

func (ctrl *Controller) SumMileageJSON(c *gin.Context) {
	c.Request.ParseForm()

	id, start, end := query.ParseQueryParams(c)

	c.JSON(http.StatusOK, map[string]interface{}{
		"mileage": ctrl.provider.SumMileageByVehicle(c, id, start, end),
	})
}

func (ctrl *Controller) ShowAllByManagerID(c *gin.Context) {
	id := c.Query("manager")
	v, err := strconv.Atoi(id)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request")
		c.Abort()
		return
	}

	ents := ctrl.provider.FetchAllByManagerID(c, int64(v))

	if len(ents) == 0 {
		c.String(http.StatusBadRequest, "bad request")
		c.Abort()
		return
	}

	c.HTML(http.StatusOK, "manager_enterprises.html", gin.H{
		"title":       "Enterprises",
		"enterprises": ents,
	})
}

func (ctrl *Controller) FetchRawByManagerID(c *gin.Context, id int64) []models.Enterprise {
	return ctrl.provider.FetchAllByManagerID(c, id)
}
