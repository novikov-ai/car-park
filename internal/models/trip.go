package models

import "time"

type Trip struct {
	ID              int64      `json:"id"`
	VehicleID       int64      `json:"vehicle_id"`
	StartedPoint    *string    `json:"started_point"`
	EndedPoint      *string    `json:"ended_point"`
	LocationStart   *Location  `json:"location_start"`
	LocationEnd     *Location  `json:"location_end"`
	StartedAt       *time.Time `json:"started_at"`
	EndedAt         *time.Time `json:"ended_at"`
	ScheduledAt     *time.Time `json:"scheduled_at"`
	TrackLength     int        `json:"track_length"`
	MaxVelocity     int        `json:"max_velocity"`
	MaxAcceleration int        `json:"max_acceleration"`
}
