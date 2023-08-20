package trips

import (
	"car-park/internal/models"
	"car-park/internal/storage"
	"context"
	"time"
)

type Provider struct {
	db storage.Client
}

func New(st storage.Client) *Provider {
	return &Provider{
		db: st,
	}
}

func (p *Provider) Create(ctx context.Context, trip models.Trip) int64 {
	query := `INSERT INTO trip 
    	(vehicle_id, started_point, ended_point, started_at, ended_at, track_length, max_velocity, max_acceleration)
    VALUES
        ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id`

	resp := p.db.QueryRow(ctx, query,
		trip.VehicleID, trip.StartedPoint, trip.EndedPoint, trip.StartedAt, trip.EndedAt,
		trip.TrackLength, trip.MaxVelocity, trip.MaxAcceleration)

	var id int64
	err := resp.Scan(&id)
	if err != nil {
		return 0
	}

	return id
}

func (p *Provider) EndTrip(ctx context.Context, id int64, endTime time.Time) error {
	query := `UPDATE trip 
SET ended_at = $1
WHERE id = $2`

	_, err := p.db.Query(ctx, query, endTime, id)

	return err
}
