package drivers

import (
	"car-park/internal/models"
	"car-park/internal/storage"
	"context"
	"fmt"
	"os"
)

type Provider struct {
	db storage.Client
}

func New(st storage.Client) *Provider {
	return &Provider{
		db: st,
	}
}

func (p *Provider) FetchAll(ctx context.Context) []models.Driver {
	query := `SELECT * FROM driver`
	resp, err := p.db.Query(ctx, query)
	if err != nil {
		panic(err)
	}

	var (
		id, enterpriseID        int64
		activeCarID             *int64
		age, salary, experience int
	)

	drivers := make([]models.Driver, 0)
	for resp.Next() {
		err = resp.Scan(&id, &enterpriseID, &activeCarID, &age, &salary, &experience)
		if err != nil {
			fmt.Fprintf(os.Stderr, "scan failed: %v\n", err)
			return []models.Driver{}
		}

		drivers = append(drivers, models.Driver{
			ID:           id,
			EnterpriseID: enterpriseID,
			ActiveCarID:  activeCarID,
			Age:          age,
			Salary:       salary,
			Experience:   salary,
		})
	}

	return drivers
}

type vehicleDriver struct {
	VehicleID int64 `json:"vehicleId"`
	DriverID  int64 `json:"driverId"`
}

func (p *Provider) FetchAllVehicleDrivers(ctx context.Context) []vehicleDriver {
	query := `SELECT * FROM driver_vehicle`
	resp, err := p.db.Query(ctx, query)
	if err != nil {
		panic(err)
	}

	var (
		vehicle_id, driver_id int64
	)

	vehicleDrivers := make([]vehicleDriver, 0)
	for resp.Next() {
		err = resp.Scan(&vehicle_id, &driver_id)
		if err != nil {
			fmt.Fprintf(os.Stderr, "scan failed: %v\n", err)
			return []vehicleDriver{}
		}

		vehicleDrivers = append(vehicleDrivers, vehicleDriver{
			VehicleID: vehicle_id,
			DriverID:  driver_id,
		})
	}

	return vehicleDrivers
}
