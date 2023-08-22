package controllers

import (
	"car-park/internal/controllers/admin"
	"car-park/internal/controllers/drivers"
	"car-park/internal/controllers/enterprises"
	"car-park/internal/controllers/geo"
	"car-park/internal/controllers/vehicles"
	"car-park/internal/storage"
)

type Controllers struct {
	Enterprises *enterprises.Controller
	Vehicles    *vehicles.Controller
	Drivers     *drivers.Controller
	Geo         *geo.Controller
	Admin       *admin.Controller
}

func New(repository storage.Client) *Controllers {
	vehiclesCtrl := vehicles.NewCtrl(repository)
	enterprisesCtrl := enterprises.NewCtrl(repository)

	return &Controllers{
		Enterprises: enterprisesCtrl,
		Vehicles:    vehiclesCtrl,
		Drivers:     drivers.NewCtrl(repository),
		Geo:         geo.NewCtrl(repository),
		Admin:       admin.NewCtrl(vehiclesCtrl, enterprisesCtrl),
	}
}
