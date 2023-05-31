package models

import (
	"context"
	"fmt"
	"os"

	"car-park/internal/models"
	"car-park/internal/storage"
)

type Provider struct {
	db storage.Client
}

func New(st storage.Client) *Provider {
	return &Provider{
		db: st,
	}
}

func (p *Provider) FetchAll(ctx context.Context) []models.Model {
	query := `SELECT * FROM model`
	resp, err := p.db.Query(ctx, query)
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
			return []models.Model{}
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
