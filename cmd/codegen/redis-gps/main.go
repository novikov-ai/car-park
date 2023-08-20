package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

func main() {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// Generate Redis commands for adding GPS points
	for i := 0; i < 30; i++ {
		vehicleID := rand.Intn(5) + 1
		latitude := generateRandomCoordinate(-90, 90)
		longitude := generateRandomCoordinate(-180, 180)
		timestamp := generateRandomTimestamp(1625241600, 1625328000)

		redisCommand := fmt.Sprintf("ZADD gps_points %d \"%s,%s\"", timestamp, strconv.Itoa(vehicleID), formatCoordinate(latitude), formatCoordinate(longitude))
		fmt.Println(redisCommand)
	}
}

func generateRandomCoordinate(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func generateRandomTimestamp(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func formatCoordinate(coord float64) string {
	return strconv.FormatFloat(coord, 'f', -1, 64)
}
