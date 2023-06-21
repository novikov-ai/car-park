package fixtures

import (
	"context"
	"fmt"
	"os"

	"car-park/internal/models"
	"car-park/internal/storage"
)

type Maker struct {
	db storage.Client
}

func New(st storage.Client) *Maker {
	return &Maker{
		db: st,
	}
}

func (m *Maker) CreateVehicle(vehicle models.Vehicle, enterpriseID int64) int64 {
	query := `INSERT INTO vehicle (model_id, enterprise_id, price, manufacture_year, mileage, color, vin)
    VALUES
        ($1, $2, $3, $4, $5, $6, $7)
RETURNING id`

	row := m.db.QueryRow(context.Background(), query,
		vehicle.ModelID, enterpriseID, vehicle.Price, vehicle.ManufactureYear, vehicle.Mileage, vehicle.Color, vehicle.VIN)

	var insertedID int64
	if err := row.Scan(&insertedID); err != nil {
		fmt.Fprintf(os.Stderr, "can't create fixtures: %v", err)
		return 0
	}

	return insertedID
}

func (m *Maker) CreateDrivers(drivers []models.Driver) {
	query := `INSERT INTO driver (enterprise_id, active_car_id, age, salary, experience)
VALUES
    ($1, $2, $3, $4, $5)`

	for _, driver := range drivers {
		resp, err := m.db.Query(context.Background(), query,
			driver.EnterpriseID, driver.ActiveCarID, driver.Age, driver.Salary, driver.Experience)

		if err != nil {
			fmt.Fprintf(os.Stderr, "can't create fixtures: %v", err)
			return
		}

		resp.Close()
	}
}
