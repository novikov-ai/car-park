package drivers

import (
	"car-park/internal/storage"
	"github.com/gin-gonic/gin"
	"net/http"
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
		"drivers": ctrl.provider.FetchAll(c),
	})
}
