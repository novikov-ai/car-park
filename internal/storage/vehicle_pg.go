package storage

import (
	"car-park/internal/models"
	"context"
	"github.com/jackc/pgx/v5"
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
		price, year, mileage, color int
		vin                         string
		created, updated            time.Time
		deleted                     *time.Time
	)

	vehicles := make([]models.Vehicle, 0)
	for resp.Next() {
		err = resp.Scan(&id, &price, &year, &mileage, &color, &vin, &created, &updated, &deleted)
		if err != nil {
			panic(err)
		}
		vehicles = append(vehicles, models.Vehicle{
			ID:              id,
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
