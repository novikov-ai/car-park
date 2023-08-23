package main

import (
	"bytes"
	"car-park/cmd"
	"car-park/internal/controllers/geo"
	"car-park/internal/controllers/trips"
	"car-park/internal/models"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

var (
	vehicleID              = int64(0)
	trackLength            = int64(0)
	vehicleMaxSpeed        = int64(0)
	vehicleMaxAcceleration = int64(0)
	dotsStep               = int64(0)

	writeEvery = time.Second * 10
)

func init() {
	flag.Int64Var(&vehicleID, "vehicle", 0, "vehicle id")
	flag.Int64Var(&trackLength, "length", 0, "track length")
	flag.Int64Var(&vehicleMaxSpeed, "speed", 0, "vehicle max speed")
	flag.Int64Var(&vehicleMaxAcceleration, "acceleration", 0, "vehicle max acceleration")
	flag.Int64Var(&dotsStep, "step", 0, "dots step")
}

func main() {
	flag.Parse()

	//if vehicleID == 0 {
	//	panic("must provide an integer vehicle id")
	//}

	ctx := context.Background()

	repository := cmd.MustInitDB(ctx)
	defer repository.Close(ctx)

	geoClient := geo.NewGeoClient(repository)
	provider := trips.New(repository)

	rand.Seed(time.Now().UnixNano())

	for i := 1; i <= 10; i++ {
		i := i
		fmt.Println("Vehicle # ", i, " started to record...")

		for j := 0; j <= rand.Intn(20); j++ {
			start, end := genStartEnd()
			coords, err := getCoordinatesByStartEnd(start, end)
			if coords == nil || err != nil {
				fmt.Println("failed fetch coords, more: ", err)
				continue
			}

			day := rand.Intn(365)
			month := rand.Intn(12)
			year := rand.Intn(20)

			tripStartedShifted := time.Now().Add(-time.Hour * 24 * time.Duration(day))
			tripStartedShifted.Add(-time.Hour * 24 * 31 * time.Duration(month))
			tripStartedShifted.Add(-time.Hour * 24 * 31 * 365 * time.Duration(year))

			tripStarted := tripStartedShifted.UTC()

			startPoint := fmt.Sprintf("%f,%f", start.Longitude, start.Latitude)
			endPoint := fmt.Sprintf("%f,%f", end.Longitude, end.Latitude)

			tripID := provider.Create(ctx, models.Trip{
				VehicleID:    int64(i),
				StartedPoint: &startPoint,
				EndedPoint:   &endPoint,
				StartedAt:    &tripStarted,
				EndedAt:      nil,
				TrackLength: DistanceBetweenPoints(geo.GPSPoint{
					Longitude: start.Longitude,
					Latitude:  start.Latitude,
				}, geo.GPSPoint{
					Longitude: end.Longitude,
					Latitude:  end.Latitude,
				}),
				MaxVelocity:     rand.Intn(100) + 50,
				MaxAcceleration: rand.Intn(100) + 50,
			})

			if tripID == 0 {
				fmt.Println("trip hasn't been created")
				continue
			}

			for _, coord := range coords {
				coord := coord

				lng := coord[0]
				lat := coord[1]

				id := geoClient.AddPointToTrack(ctx, geo.GPSPoint{
					Longitude: lng,
					Latitude:  lat,
					Created:   tripStartedShifted.Add(time.Minute).UTC(),
				}, geo.VehicleTrack{
					ID:        tripID,
					VehicleID: int64(i),
				})

				if id == 0 {
					fmt.Println("error ocurried: ")
					return
				}

				// fmt.Printf("#%v: wrote point #%v with lng: %v, lat: %v\n", i, id, lng, lat)

				// time.Sleep(writeEvery)
			}

			tripEnded := tripStartedShifted.Add(time.Hour * time.Duration(rand.Intn(24))).UTC()
			err = provider.EndTrip(ctx, tripID, tripEnded)
			if err != nil {
				fmt.Println("Error happened:", err)
			}
		}
		fmt.Println("Vehicle # ", i, " ended!")

		time.Sleep(writeEvery)
	}

	fmt.Println("success!")
}

const (
	earthRadiusKm = 6371 // Earth's radius in kilometers
)

func degToRad(deg float64) float64 {
	return deg * (math.Pi / 180)
}

func DistanceBetweenPoints(point1, point2 geo.GPSPoint) int {
	lat1 := degToRad(point1.Latitude)
	lon1 := degToRad(point1.Longitude)
	lat2 := degToRad(point2.Latitude)
	lon2 := degToRad(point2.Longitude)

	dLat := lat2 - lat1
	dLon := lon2 - lon1

	a := math.Pow(math.Sin(dLat/2), 2) + math.Cos(lat1)*math.Cos(lat2)*math.Pow(math.Sin(dLon/2), 2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	distance := earthRadiusKm * c
	return int(distance)
}

func getCoordinatesByStartEnd(start, end geo.GPSPoint) ([][]float64, error) {
	requestBody := map[string]interface{}{
		"points": [][]float64{
			{
				start.Longitude,
				start.Latitude,
			},
			{
				end.Longitude,
				end.Latitude,
			},
		},
		"snap_preventions": []string{
			"motorway",
			"ferry",
			"tunnel",
		},
		"details": []string{
			"road_class",
			"surface",
		},
		"profile":        "car",
		"locale":         "en",
		"instructions":   false,
		"calc_points":    true,
		"points_encoded": false,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(http.MethodPost,
		"https://graphhopper.com/api/1/route",
		bytes.NewBuffer(jsonBody))

	query := request.URL.Query()
	query.Add("key", "6761ff87-c061-41b6-99c2-d2de3cc15a8b")

	request.URL.RawQuery = query.Encode()

	request.Header.Add("Content-Type", "application/json")

	cl := http.DefaultClient
	response, err := cl.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var coords Coordinates
	err = json.Unmarshal(body, &coords)
	if err != nil {
		return nil, err
	}

	if len(coords.Paths) == 0 {
		return nil, errors.New("coords paths is empty")
	}

	if len(coords.Paths[0].Points.Coordinates) == 0 {
		return nil, errors.New("coords paths points is empty")
	}

	return coords.Paths[0].Points.Coordinates, nil
}

func genStartEnd() (geo.GPSPoint, geo.GPSPoint) {
	lts := fetchLocations()

	rand.Seed(time.Now().UnixNano())

	p1 := lts.Geonames[rand.Intn(len(lts.Geonames))]
	p2 := lts.Geonames[rand.Intn(len(lts.Geonames))]

	lng1, _ := strconv.ParseFloat(p1.Lng, 64)
	lat1, _ := strconv.ParseFloat(p1.Lat, 64)

	lng2, _ := strconv.ParseFloat(p2.Lng, 64)
	lat2, _ := strconv.ParseFloat(p2.Lat, 64)

	return geo.GPSPoint{
			Longitude: lng1,
			Latitude:  lat1,
		}, geo.GPSPoint{
			Longitude: lng2,
			Latitude:  lat2,
		}
}

func fetchLocations() Locations {
	request, err := http.NewRequest(http.MethodGet,
		"http://api.geonames.org/searchJSON?country=RU&username=al.novikov08",
		nil)

	request.Header.Add("Content-Type", "application/json")

	cl := http.DefaultClient
	response, err := cl.Do(request)
	if err != nil {
		return Locations{}
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return Locations{}
	}

	var l Locations
	err = json.Unmarshal(body, &l)
	if err != nil {
		return Locations{}
	}

	return l
}

type Locations struct {
	TotalResultsCount int `json:"totalResultsCount"`
	Geonames          []struct {
		AdminCode1  string `json:"adminCode1"`
		Lng         string `json:"lng"`
		GeonameId   int    `json:"geonameId"`
		ToponymName string `json:"toponymName"`
		CountryId   string `json:"countryId"`
		Fcl         string `json:"fcl"`
		Population  int    `json:"population"`
		CountryCode string `json:"countryCode"`
		Name        string `json:"name"`
		FclName     string `json:"fclName"`
		CountryName string `json:"countryName"`
		FcodeName   string `json:"fcodeName"`
		AdminName1  string `json:"adminName1"`
		Lat         string `json:"lat"`
		Fcode       string `json:"fcode"`
		AdminCodes1 struct {
			ISO31662 string `json:"ISO3166_2"`
		} `json:"adminCodes1,omitempty"`
	} `json:"geonames"`
}

type Coordinates struct {
	Paths []struct {
		Distance  float64   `json:"distance"`
		Weight    float64   `json:"weight"`
		Time      int       `json:"time"`
		Transfers int       `json:"transfers"`
		Bbox      []float64 `json:"bbox"`
		Points    struct {
			Type        string      `json:"type"`
			Coordinates [][]float64 `json:"coordinates"`
		} `json:"points"`
		Instructions []struct {
			Distance             float64     `json:"distance"`
			Heading              float64     `json:"heading,omitempty"`
			Sign                 int         `json:"sign"`
			Interval             []int       `json:"interval"`
			Text                 string      `json:"text"`
			Time                 int         `json:"time"`
			StreetName           string      `json:"street_name"`
			Points               [][]float64 `json:"points"`
			StreetRef            string      `json:"street_ref,omitempty"`
			StreetDestination    string      `json:"street_destination,omitempty"`
			StreetDestinationRef string      `json:"street_destination_ref,omitempty"`
			LastHeading          float64     `json:"last_heading,omitempty"`
		} `json:"instructions"`
		Legs    []interface{} `json:"legs"`
		Details struct {
			Surface   [][]interface{} `json:"surface"`
			RoadClass [][]interface{} `json:"road_class"`
		} `json:"details"`
		Ascend           float64 `json:"ascend"`
		Descend          float64 `json:"descend"`
		SnappedWaypoints struct {
			Type        string      `json:"type"`
			Coordinates [][]float64 `json:"coordinates"`
		} `json:"snapped_waypoints"`
	} `json:"paths"`
}
