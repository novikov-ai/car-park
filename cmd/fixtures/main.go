package main

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/joho/godotenv"

	"car-park/cmd"
	"car-park/internal/models"
	"car-park/internal/storage"
	"car-park/internal/usecases/fixtures"
)

var (
	enterprisesIDs      string
	vehiclesN, driversN string
)

/*
Сгенерируйте для 3 предприятий по 3-5 тысяч машин и водителей (активный водитель для каждой 10-й машины например),
и разберитесь, как через REST API получать их в режиме пагинации -- не все разом, а листать страничками по 20-50 машинок.
*/

func init() {
	flag.StringVar(&enterprisesIDs, "enterprises", "1;2;3", "enterprises ID with ; delimiter")
	flag.StringVar(&vehiclesN, "vehicles", "4000", "vehicles quantity")
	flag.StringVar(&driversN, "drivers", "5000", "drivers quantity")

	err := godotenv.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to load environment-file: %v\n", err)
		os.Exit(1)
	}
}

func main() {
	flag.Parse()

	ctx := context.Background()

	db := cmd.MustInitDB(ctx)
	defer db.Close(context.Background())

	enterprises := mustParseEnterprises()

	v := mustParseValue(vehiclesN)
	d := mustParseValue(driversN)

	repository := storage.New(db)
	factory := fixtures.New(repository)

	for _, entID := range enterprises {
		vehicles := generateVehicles(entID, v)

		vehiclesIDs := make([]int64, 0, len(vehicles))
		for _, vehicle := range vehicles {
			insertedID := factory.CreateVehicle(vehicle, entID)
			if insertedID != 0 {
				vehiclesIDs = append(vehiclesIDs, insertedID)
			}
		}

		drivers := generateDrivers(entID, vehiclesIDs, d)
		factory.CreateDrivers(drivers)

		fmt.Printf(`Enterprise #%v

Generated vehicles:
%v

Generated drivers:
%v

`, entID, len(vehicles), len(drivers))
	}
}

func mustParseEnterprises() []int64 {
	ids := strings.Split(enterprisesIDs, ";")
	parsedIDs := make([]int64, 0, len(ids))

	for _, id := range ids {
		v, err := strconv.Atoi(id)
		if err != nil {
			fmt.Fprintf(os.Stderr, "wrong arguments format")
			os.Exit(1)
		}

		parsedIDs = append(parsedIDs, int64(v))
	}

	return parsedIDs
}

func mustParseValue(value string) int {
	v, err := strconv.Atoi(value)
	if err != nil {
		fmt.Fprintf(os.Stderr, "wrong arguments format")
		os.Exit(1)
	}

	return v
}

func generateVehicles(entID int64, n int) []models.Vehicle {
	vv := make([]models.Vehicle, 0, n)

	wg := sync.WaitGroup{}
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			v := models.Vehicle{
				ID:              rand.Int63(),
				ModelID:         int64(rand.Intn(4)), // todo according to db
				EnterpriseID:    &entID,
				Price:           15000 + rand.Intn(200_000_000),
				ManufactureYear: 1967 + rand.Intn(56),
				Mileage:         rand.Intn(300_000),
				Color:           rand.Intn(4), // todo according to db
				VIN:             randomVIN(),
			}

			vv = append(vv, v)
		}()
	}
	wg.Wait()

	return vv
}

func randomVIN() string {
	dummyVIN := func() string {
		s := rand.New(rand.NewSource(time.Now().UnixNano()))
		return strconv.Itoa(s.Int())
	}

	//resp, err := http.Get("https://randomvin.com/getvin.php?type=real")
	//if err != nil {
	//	dummyVIN()
	//}
	//defer resp.Body.Close()
	//
	//body, err := io.ReadAll(resp.Body)
	//if err != nil {
	//	dummyVIN()
	//}

	//return string(body)

	return dummyVIN()
}

func generateDrivers(entID int64, vehicles []int64, n int) []models.Driver {
	dd := make([]models.Driver, 0, n)

	for i := 0; i < n; i++ {
		d := models.Driver{
			ID:           rand.Int63(),
			EnterpriseID: entID,
			ActiveCarID:  nil, // every 10 with an active driver
			Age:          18 + rand.Intn(82),
			Salary:       10_000 + rand.Intn(140_000),
			Experience:   rand.Intn(30),
		}

		dd = append(dd, d)
	}

	for i, v := range vehicles {
		v := v

		if i%9 != 0 {
			continue
		}

		if len(dd) == 0 {
			break
		}

		dd[rand.Intn(len(dd))].ActiveCarID = &v
	}

	return dd
}
