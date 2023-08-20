package cmd

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"os"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to load environment-file: %v\n", err)
		os.Exit(1)
	}
}

func MustInitDB(ctx context.Context) *pgx.Conn {
	conn, err := pgx.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	err = conn.Ping(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fail to ping db: %v", err)
		os.Exit(1)
	}

	return conn
}

func MustInitCache(ctx context.Context) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}

/*
package main

import (
 "encoding/json"
 "log"
 "net/http"
 "strings"

 "github.com/go-redis/redis"
)

var redisClient *redis.Client

type Trip struct {
 ID           string json:"ID"
 StartedPoint string json:"StartedPoint"
 EndedPoint   string json:"EndedPoint"
 StartedAt    string json:"StartedAt"
 EndedAt      string json:"EndedAt"
}

func main() {
 redisClient = redis.NewClient(&redis.Options{
  Addr: "localhost:6379",
 })

 http.HandleFunc("/api/v1/coordinates", getCoordinates)
 http.HandleFunc("/api/v1/trips", getTrips)
 http.Handle("/", http.FileServer(http.Dir("./public")))

 log.Println("Server is running on http://localhost:8080")
 http.ListenAndServe(":8080", nil)
}

func getCoordinates(w http.ResponseWriter, r *http.Request) {
 tripID := r.URL.Query().Get("id")
 if tripID == "" {
  http.Error(w, "Missing 'id' parameter", http.StatusBadRequest)
  return
 }

 // Fetch coordinates from Redis (replace "your_redis_key" with the actual key storing coordinates)
 coordinates, err := redisClient.LRange("your_redis_key", 0, -1).Result()
 if err != nil {
  http.Error(w, "Error fetching coordinates from Redis", http.StatusInternalServerError)
  return
 }

 // Convert coordinates to a slice of maps (lgn, lat)
 var coordinatesList []map[string]float64
 for _, coord := range coordinates {
  coords := strings.Split(coord, ",")
  if len(coords) == 2 {
   lgn := parseFloat(coords[0])
   lat := parseFloat(coords[1])
   coordinatesList = append(coordinatesList, map[string]float64{"lgn": lgn, "lat": lat})
  }
 }

 response := map[string]interface{}{
  "coordinates": coordinatesList,
 }

 w.Header().Set("Content-Type", "application/json")
 json.NewEncoder(w).Encode(response)
}

// Helper function to convert a string to a float64.
func parseFloat(s string) float64 {
 val, err := strconv.ParseFloat(s, 64)
 if err != nil {
  return 0.0
 }
 return val
}
*/
