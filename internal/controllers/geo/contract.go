package geo

import (
	"car-park/internal/models"
	"context"
)

type Client interface {
	AddPointToTrack(ctx context.Context, point GPSPoint, track VehicleTrack) int64
	GetTrackByTrip(c context.Context,
		tripID int64,
		start, end float64) map[string]interface{}
	ToGeoJSON(points map[string]interface{}) []GeoPoint

	GetTripsByVehicle(c context.Context, vehicleID int64, start, end string) ([]models.Trip, error)
}
