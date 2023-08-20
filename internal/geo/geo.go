package geo

import (
	"bytes"
	"car-park/internal/models"
	"car-park/internal/storage"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type VehicleTrack struct {
	ID, VehicleID int64
}

type GeoParsingRequest struct {
	Start, End *GPSPoint
}

type GPSPoint struct {
	Longitude float64   `json:"longitude"`
	Latitude  float64   `json:"latitude"`
	Created   time.Time `json:"created"`
}

type GeoPoint struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

type redisClient struct {
	r  *redis.Client
	st storage.Client
}

type gpsClient struct {
	st storage.Client
}

func NewGeoClient(storage storage.Client) Client {
	return &gpsClient{st: storage}
}

func (gc *gpsClient) AddPointToTrack(ctx context.Context, point GPSPoint, track VehicleTrack) int64 {
	row := gc.st.QueryRow(ctx, `INSERT INTO gps_point (trip_id, longitude, latitude, created_at)
VALUES ($1, $2, $3, $4)
RETURNING id`, track.ID, point.Longitude, point.Latitude, point.Created)

	var id int64
	err := row.Scan(&id)
	if err != nil {
		return 0
	}

	return id
}
func (gc *gpsClient) GetTrackByTrip(c context.Context,
	tripID int64,
	start, end float64) map[string]interface{} {

	//track := make([]GPSPoint, 0, 0)

	gpsPointsRows, err := gc.st.Query(c, `SELECT longitude, latitude, created_at
FROM gps_point
WHERE trip_id = $1 AND 
      EXTRACT(EPOCH FROM created_at) >= $2 AND 
      EXTRACT(EPOCH FROM created_at) <= $3`,
		tripID, start, end)

	if err != nil {
		return nil
	}

	//resp := gc.st.QueryRow(c, `SELECT utc
	//FROM enterprise as e
	//JOIN vehicle as v on v.enterprise_id = e.id
	//WHERE v.id = $1`, vehicleID)
	//
	//var utc int
	//err = resp.Scan(&utc)
	//if err != nil {
	//	return nil
	//}

	var coordinatesList []map[string]float64

	points := make([]GPSPoint, 0, 0)
	for gpsPointsRows.Next() {
		var point GPSPoint
		gpsPointsRows.Scan(&point.Longitude, &point.Latitude, &point.Created)

		//point.Created.Add(time.Second * time.Duration(60*60*utc)) // enterprise timeshift
		points = append(points, point)

		coordinatesList = append(coordinatesList, map[string]float64{"lgn": point.Longitude, "lat": point.Latitude})
	}

	//for _, point := range result {
	//	enterpriseTimeShift := float64(60 * 60 * utc)
	//
	//	parsed, err := parseGPSPoint(point.Member, point.Score+enterpriseTimeShift)
	//	if err != nil {
	//		break
	//	}
	//	track = append(track, parsed)
	//
	//	coordinatesList = append(coordinatesList, map[string]float64{"lgn": parsed.Longitude, "lat": parsed.Latitude})
	//}

	response := map[string]interface{}{
		"coordinates": coordinatesList,
	}

	return response
}
func (gc *gpsClient) ToGeoJSON(pp map[string]interface{}) []GeoPoint {
	points := make([]GeoPoint, 0, len(pp))

	for _, v := range pp {
		vv, ok := v.(map[string]float64)
		if !ok {
			return nil
		}

		points = append(points, GeoPoint{
			Type:        "Point",
			Coordinates: []float64{vv["lgn"], vv["lat"]},
		})
	}

	return points
}

func (gc *gpsClient) GetTripsByVehicle(ctx context.Context,
	vehicleID int64, start, end string) ([]models.Trip, error) {

	query := `SELECT id, vehicle_id, started_point, ended_point, started_at, ended_at, track_length, max_velocity, max_acceleration
FROM trip
WHERE vehicle_id = $1 AND started_at >= $2 AND ended_at <= $3
`
	resp, err := gc.st.Query(ctx, query, vehicleID, start, end)
	if err != nil {
		fmt.Fprintf(os.Stderr, "query failed: %v\n", err)
		return []models.Trip{}, nil
	}

	locationsRequest := map[int64]GeoParsingRequest{}

	trips := make([]models.Trip, 0, 0)
	for resp.Next() {
		trip := models.Trip{
			ID:              0,
			VehicleID:       0,
			StartedPoint:    nil,
			EndedPoint:      nil,
			StartedAt:       nil,
			EndedAt:         nil,
			ScheduledAt:     nil,
			TrackLength:     0,
			MaxVelocity:     0,
			MaxAcceleration: 0,
		}

		var (
			startPoint, endPoint *string
			started, ended       *time.Time
		)

		var length, velocity, acceleration *int

		err = resp.Scan(&trip.ID, &trip.VehicleID, &startPoint, &endPoint, &started, &ended,
			&length, &velocity, &acceleration)

		if err != nil {
			fmt.Fprintf(os.Stderr, "query failed: %v\n", err)
			return []models.Trip{}, nil
		}

		if length != nil {
			trip.TrackLength = *length
		}
		if velocity != nil {
			trip.MaxVelocity = *velocity
		}
		if acceleration != nil {
			trip.MaxAcceleration = *acceleration
		}

		trip.StartedPoint = startPoint
		trip.EndedPoint = endPoint
		trip.StartedAt = started
		trip.EndedAt = ended

		startedGpsPoint := parsePointToGpsPoint(trip.StartedPoint)
		endedGpsPoint := parsePointToGpsPoint(trip.EndedPoint)

		locationsRequest[trip.ID] = GeoParsingRequest{
			Start: startedGpsPoint, End: endedGpsPoint,
		}

		trips = append(trips, trip)
	}

	if len(trips) == 0 {
		return []models.Trip{}, errors.New("error")
	}

	type LatLng struct {
		Lat float64 `json:"lat"`
		Lng float64 `json:"lng"`
	}

	type locationLatLng struct {
		start LatLng
		end   LatLng
	}

	type geoData struct {
		Data LatLng `json:"latLng"`
	}

	locs := make([]locationLatLng, 0, len(trips))
	locsTest := make([]geoData, 0, len(trips))

	for _, trip := range trips {
		point := locationsRequest[trip.ID]

		locs = append(locs, locationLatLng{
			start: LatLng{
				Lat: point.Start.Latitude,
				Lng: point.Start.Longitude,
			},
			end: LatLng{
				Lat: point.End.Latitude,
				Lng: point.End.Longitude,
			},
		})

		locsTest = append(locsTest, geoData{
			Data: LatLng{
				Lat: point.Start.Latitude,
				Lng: point.Start.Longitude,
			},
		})

		locsTest = append(locsTest, geoData{
			Data: LatLng{
				Lat: point.End.Latitude,
				Lng: point.End.Longitude,
			},
		})
	}

	requestBody := map[string]interface{}{
		"locations": locsTest,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return []models.Trip{}, err
	}

	request, err := http.NewRequest(http.MethodPost,
		"https://www.mapquestapi.com/geocoding/v1/batch?key=7Qpwe6gyKAxgzpNrNV0bsZ4LB1YiBsrt",
		bytes.NewBuffer(jsonBody))

	request.Header.Set("Content-Type", "application/json")

	cl := http.DefaultClient
	respGeo, err := cl.Do(request)
	if err != nil {
		return []models.Trip{}, err
	}
	defer respGeo.Body.Close()

	if respGeo.StatusCode == http.StatusOK {
		body, err := io.ReadAll(respGeo.Body)
		if err != nil {
			return []models.Trip{}, err
		}

		var decoder models.GeoDecoder
		err = json.Unmarshal(body, &decoder)
		if err != nil {
			return []models.Trip{}, err
		}

		k := 0
		for i, r := range decoder.Results {
			k = i / 2

			if i%2 == 0 {
				for _, rl := range r.Locations {
					trips[k].LocationStart = &models.Location{
						Street:     rl.Street,
						Area:       rl.AdminArea4,
						Area2:      rl.AdminArea5,
						City:       rl.AdminArea6,
						PostalCode: rl.PostalCode,
					}
				}
			} else {
				for _, rl := range r.Locations {
					trips[k].LocationEnd = &models.Location{
						Street:     rl.Street,
						Area:       rl.AdminArea4,
						Area2:      rl.AdminArea5,
						City:       rl.AdminArea6,
						PostalCode: rl.PostalCode,
					}
				}
			}
		}
	}

	return trips, nil
}

func parsePointToGpsPoint(point *string) *GPSPoint {
	if point == nil {
		return nil
	}

	ll := strings.Split(*point, ",")
	if len(ll) != 2 {
		return nil
	}

	longitude, err := strconv.ParseFloat(ll[0], 64)
	if err != nil {
		return nil
	}

	latitude, err := strconv.ParseFloat(ll[1], 64)
	if err != nil {
		return nil
	}

	return &GPSPoint{
		Longitude: longitude,
		Latitude:  latitude,
	}
}

//
//func New(redis *redis.Client, storage storage.Client) Client {
//	return &redisClient{
//		r:  redis,
//		st: storage,
//	}
//}

// todo: refactor below
//
//func (rc *redisClient) AddPointToTrack(ctx context.Context,
//	point GPSPoint, track VehicleTrack) error {
//
//	err := rc.r.ZAdd(ctx, compoundTrackKey(track), redis.Z{
//		Score:  float64(time.Now().UTC().Unix()),
//		Member: fmt.Sprintf("%f,%f", point.Latitude, point.Longitude),
//	}).Err()
//
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func compoundTrackKey(track VehicleTrack) string {
//	//return fmt.Sprintf("%d:%d", track.VehicleID, track.ID)
//	return fmt.Sprintf("%v", track.VehicleID)
//}
//
//func (rc *redisClient) GetTrack(ctx context.Context, vehicleID int64, start, end float64) map[string]interface{} {
//	//func (rc *redisClient) GetTrack(ctx context.Context, vehicleID int64, start, end float64) ([]GPSPoint, error) {
//	track := make([]GPSPoint, 0, 0)
//
//	result, err := rc.r.ZRangeByScoreWithScores(ctx,
//		//fmt.Sprintf("%s", strconv.FormatInt(vehicleID, 10)), // get all keys starting with vehicleID
//		strconv.FormatInt(vehicleID, 10),
//		&redis.ZRangeBy{
//			Min: fmt.Sprintf("%f", float64(0)),
//			Max: fmt.Sprintf("%f", end),
//		}).Result()
//
//	if err != nil {
//		return nil
//	}
//
//	resp := rc.st.QueryRow(ctx, `SELECT utc
//FROM enterprise as e
//JOIN vehicle as v on v.enterprise_id = e.id
//WHERE v.id = $1`, vehicleID)
//
//	var utc int
//	err = resp.Scan(&utc)
//	if err != nil {
//		return nil
//	}
//
//	var coordinatesList []map[string]float64
//	for _, point := range result {
//		enterpriseTimeShift := float64(60 * 60 * utc)
//
//		parsed, err := parseGPSPoint(point.Member, point.Score+enterpriseTimeShift)
//		if err != nil {
//			break
//		}
//		track = append(track, parsed)
//
//		coordinatesList = append(coordinatesList, map[string]float64{"lgn": parsed.Longitude, "lat": parsed.Latitude})
//	}
//
//	response := map[string]interface{}{
//		"coordinates": coordinatesList,
//	}
//
//	return response
//}
//
//func (rc *redisClient) ToGeoJSON(pp []GPSPoint) []GeoPoint {
//	points := make([]GeoPoint, 0, len(pp))
//
//	for _, p := range pp {
//		points = append(points, GeoPoint{
//			Type:        "Point",
//			Coordinates: []float64{p.Longitude, p.Latitude},
//		})
//	}
//
//	return points
//}
//
//func parseGPSPoint(pointZ interface{}, score float64) (GPSPoint, error) {
//	pointStr, ok := pointZ.(string)
//	if !ok {
//		return GPSPoint{}, errors.New("can't parse gps point")
//	}
//
//	coords := strings.Split(pointStr, ",")
//	if len(coords) != 2 {
//		return GPSPoint{}, fmt.Errorf("invalid GPS point format: %s", pointZ)
//	}
//
//	latitude, err := strconv.ParseFloat(coords[0], 64)
//	if err != nil {
//		return GPSPoint{}, fmt.Errorf("failed to parse latitude: %v", err)
//	}
//
//	longitude, err := strconv.ParseFloat(coords[1], 64)
//	if err != nil {
//		return GPSPoint{}, fmt.Errorf("failed to parse longitude: %v", err)
//	}
//
//	return GPSPoint{
//		Latitude:  latitude,
//		Longitude: longitude,
//		Created:   time.Unix(int64(score), 0),
//	}, nil
//}
//
//func (rc *redisClient) GetTripsByVehicle(ctx context.Context,
//	vehicleID int64, start, end string) ([]models.Trip, error) {
//
//	query := `SELECT id, vehicle_id, started_point, ended_point, started_at, ended_at
//FROM trip
//WHERE vehicle_id = $1 AND started_at >= $2 AND ended_at <= $3
//`
//	resp, err := rc.st.Query(ctx, query, vehicleID, start, end)
//	if err != nil {
//		fmt.Fprintf(os.Stderr, "query failed: %v\n", err)
//		return []models.Trip{}, nil
//	}
//
//	locationsRequest := map[int64]GeoParsingRequest{}
//
//	trips := make([]models.Trip, 0, 0)
//	for resp.Next() {
//		trip := models.Trip{}
//
//		err = resp.Scan(&trip.ID, &trip.VehicleID, &trip.StartedPoint, &trip.EndedPoint, &trip.StartedAt, &trip.EndedAt)
//
//		if err != nil {
//			fmt.Fprintf(os.Stderr, "query failed: %v\n", err)
//			return []models.Trip{}, nil
//		}
//
//		startedPoint := parsePointToGpsPoint(trip.StartedPoint)
//		endedPoint := parsePointToGpsPoint(trip.EndedPoint)
//
//		locationsRequest[trip.ID] = GeoParsingRequest{
//			Start: startedPoint, End: endedPoint,
//		}
//
//		trips = append(trips, trip)
//	}
//
//	if len(trips) == 0 {
//		return []models.Trip{}, errors.New("error")
//	}
//
//	type LatLng struct {
//		Lat float64 `json:"lat"`
//		Lng float64 `json:"lng"`
//	}
//
//	type locationLatLng struct {
//		start LatLng
//		end   LatLng
//	}
//
//	type geoData struct {
//		Data LatLng `json:"latLng"`
//	}
//
//	locs := make([]locationLatLng, 0, len(trips))
//	locsTest := make([]geoData, 0, len(trips))
//
//	for _, trip := range trips {
//		point := locationsRequest[trip.ID]
//
//		locs = append(locs, locationLatLng{
//			start: LatLng{
//				Lat: point.Start.Latitude,
//				Lng: point.Start.Longitude,
//			},
//			end: LatLng{
//				Lat: point.End.Latitude,
//				Lng: point.End.Longitude,
//			},
//		})
//
//		locsTest = append(locsTest, geoData{
//			Data: LatLng{
//				Lat: point.Start.Latitude,
//				Lng: point.Start.Longitude,
//			},
//		})
//
//		locsTest = append(locsTest, geoData{
//			Data: LatLng{
//				Lat: point.End.Latitude,
//				Lng: point.End.Longitude,
//			},
//		})
//	}
//
//	requestBody := map[string]interface{}{
//		"locations": locsTest,
//	}
//
//	jsonBody, err := json.Marshal(requestBody)
//	if err != nil {
//		return []models.Trip{}, err
//	}
//
//	request, err := http.NewRequest(http.MethodPost,
//		"https://www.mapquestapi.com/geocoding/v1/batch?key=7Qpwe6gyKAxgzpNrNV0bsZ4LB1YiBsrt",
//		bytes.NewBuffer(jsonBody))
//
//	request.Header.Set("Content-Type", "application/json")
//
//	cl := http.DefaultClient
//	respGeo, err := cl.Do(request)
//	if err != nil {
//		return []models.Trip{}, err
//	}
//	defer respGeo.Body.Close()
//
//	if respGeo.StatusCode == http.StatusOK {
//		body, err := io.ReadAll(respGeo.Body)
//		if err != nil {
//			return []models.Trip{}, err
//		}
//
//		var decoder models.GeoDecoder
//		err = json.Unmarshal(body, &decoder)
//		if err != nil {
//			return []models.Trip{}, err
//		}
//
//		k := 0
//		for i, r := range decoder.Results {
//			k = i / 2
//
//			if i%2 == 0 {
//				for _, rl := range r.Locations {
//					trips[k].LocationStart = &models.Location{
//						Street:     rl.Street,
//						Area:       rl.AdminArea4,
//						Area2:      rl.AdminArea5,
//						City:       rl.AdminArea6,
//						PostalCode: rl.PostalCode,
//					}
//				}
//			} else {
//				for _, rl := range r.Locations {
//					trips[k].LocationEnd = &models.Location{
//						Street:     rl.Street,
//						Area:       rl.AdminArea4,
//						Area2:      rl.AdminArea5,
//						City:       rl.AdminArea6,
//						PostalCode: rl.PostalCode,
//					}
//				}
//			}
//		}
//	}
//
//	return trips, nil
//}
//
//func parsePointToGpsPoint(point *string) *GPSPoint {
//	if point == nil {
//		return nil
//	}
//
//	ll := strings.Split(*point, ",")
//	if len(ll) != 2 {
//		return nil
//	}
//
//	longitude, err := strconv.ParseFloat(ll[0], 64)
//	if err != nil {
//		return nil
//	}
//
//	latitude, err := strconv.ParseFloat(ll[1], 64)
//	if err != nil {
//		return nil
//	}
//
//	return &GPSPoint{
//		Longitude: longitude,
//		Latitude:  latitude,
//	}
//}
