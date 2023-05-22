package storage

import (
	"car-park/internal/models"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"os"
	"time"
)

type Storage struct {
	db *pgx.Conn
}

func New(conn *pgx.Conn) *Storage {
	return &Storage{db: conn}
}

func (s *Storage) FetchAll(ctx context.Context) []models.Vehicle {
	query := `SELECT * FROM vehicle`
	resp, err := s.db.Query(ctx, query)
	if err != nil {
		panic(err)
	}

	var (
		id                          int64
		modelID                     int64
		price, year, mileage, color int
		vin                         string
		created, updated            time.Time
		deleted                     *time.Time
	)

	vehicles := make([]models.Vehicle, 0)
	for resp.Next() {
		err = resp.Scan(&id, &modelID, &price, &year, &mileage, &color, &vin, &created, &updated, &deleted)
		if err != nil {
			panic(err)
		}
		vehicles = append(vehicles, models.Vehicle{
			ID:              id,
			ModelID:         modelID,
			Price:           price,
			ManufactureYear: year,
			Mileage:         mileage,
			Color:           color,
			VIN:             vin,
			CreatedAt:       created,
			UpdatedAt:       updated,
			DeletedAt:       deleted,
		})
	}

	return vehicles
}

func (s *Storage) FetchAllModels(ctx context.Context) []models.Model {
	query := `SELECT * FROM model`
	resp, err := s.db.Query(ctx, query)
	if err != nil {
		panic(err)
	}

	var (
		id                                        int64
		brand                                     string
		vehicleType, seats, tank, vehicleCapacity int
	)

	mm := make([]models.Model, 0)
	for resp.Next() {
		err = resp.Scan(&id, &brand, &vehicleType, &seats, &tank, &vehicleCapacity)
		if err != nil {
			fmt.Fprintf(os.Stderr, "scan failed: %v\n", err)
			continue
		}

		mm = append(mm, models.Model{
			ID:              id,
			BrandName:       brand,
			VehicleType:     vehicleType,
			Seats:           seats,
			Tank:            tank,
			VehicleCapacity: vehicleCapacity,
		})
	}

	return mm
}
