package vehicles

import (
	"car-park/internal/constants"
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"car-park/internal/models"
	"car-park/internal/storage"
)

const (
	timeRound = time.Second
)

const (
	idField      = "id"
	modelField   = "model"
	priceField   = "price"
	yearField    = "year"
	mileageField = "mileage"
	colorField   = "color"
	vinField     = "vin"
)

const redirectPath = "/admin/vehicles"

type Provider struct {
	db storage.Client
}

func New(st storage.Client) *Provider {
	return &Provider{
		db: st,
	}
}

func (p *Provider) Create(c *gin.Context) int {
	err := c.Request.ParseForm()
	if err != nil {
		return 0
	}

	modelValue := c.Request.FormValue(modelField)
	colorValue := c.Request.FormValue(colorField)

	modelID := constants.Models[modelValue]
	colorID := constants.Colors[colorValue]

	toInt := func(key string) (int, error) {
		return strconv.Atoi(c.Request.Form.Get(key))
	}

	price, err := toInt(priceField)
	year, err := toInt(yearField)
	mileage, err := toInt(mileageField)

	vehicle := models.Vehicle{
		ModelID:         modelID,
		Price:           price,
		ManufactureYear: year,
		Mileage:         mileage,
		Color:           colorID,
		VIN:             c.Request.Form.Get(vinField),
	}

	query := `INSERT INTO vehicle (model_id, price, manufacture_year, mileage, color, vin, purchased_at)
    VALUES
        ($1, $2, $3, $4, $5, $6, $7)
RETURNING id`

	inserted := 0
	err = p.db.QueryRow(context.Background(), query,
		vehicle.ModelID, vehicle.Price, vehicle.ManufactureYear,
		vehicle.Mileage, vehicle.Color, vehicle.VIN, time.Now()).Scan(&inserted)

	if err != nil {
		return 0
	}

	return inserted
}

func (p *Provider) Update(c *gin.Context) error {
	err := c.Request.ParseForm()
	if err != nil {
		return err
	}

	modelID := constants.Models[c.Request.FormValue(modelField)]
	colorID := constants.Colors[c.Request.FormValue(colorField)]

	toInt := func(key string) (int, error) {
		return strconv.Atoi(c.Request.Form.Get(key))
	}

	id, err := strconv.ParseInt(c.Request.Form.Get(idField), 10, 64)
	if err != nil {
		return err
	}

	price, err := toInt(priceField)
	year, err := toInt(yearField)
	mileage, err := toInt(mileageField)

	vehicle := models.Vehicle{
		ModelID:         modelID,
		Price:           price,
		ManufactureYear: year,
		Mileage:         mileage,
		Color:           colorID,
		VIN:             c.Request.Form.Get(vinField),
	}

	query := `UPDATE vehicle 
SET model_id = $1, price = $2, manufacture_year = $3,
    mileage = $4, color = $5, vin = $6
WHERE id = $7`

	resp, err := p.db.Query(context.Background(), query,
		vehicle.ModelID, vehicle.Price, vehicle.ManufactureYear, vehicle.Mileage, vehicle.Color, vehicle.VIN, id)

	if err != nil {
		return err
	}

	resp.Close()

	return nil
}

func (p *Provider) Delete(c *gin.Context) error {
	err := c.Request.ParseForm()
	if err != nil {
		return err
	}

	id, err := strconv.ParseInt(c.Request.Form.Get(idField), 10, 64)
	if err != nil {
		return err
	}

	query := `DELETE FROM vehicle
WHERE id=$1`

	resp, err := p.db.Query(context.Background(), query, id)

	resp.Close()

	return nil
}

func (p *Provider) FetchAll(c *gin.Context) []models.Vehicle {
	query := `
SELECT
    v.id, v.model_id, v.enterprise_id, v.price, v.manufacture_year, v.mileage, v.color, v.vin, v.purchased_at,
    e.utc
FROM vehicle v
FULL JOIN enterprise e ON v.enterprise_id = e.id
WHERE v.id IS NOT NULL`

	offset := c.Query("offset")
	o, err := strconv.Atoi(offset)
	queryPagination := fmt.Sprintf("OFFSET %v", o)

	limit := c.Query("limit")
	l, err := strconv.Atoi(limit)

	if l > 0 {
		query += fmt.Sprintf("\nLIMIT %v %s", l, queryPagination)
	}

	resp, err := p.db.Query(c, query)
	if err != nil {
		fmt.Fprintf(os.Stderr, "query failed: %v\n", err)
		return []models.Vehicle{}
	}

	var (
		id, modelID                 int64
		enterpriseID                *int64
		price, year, mileage, color int
		vin                         string
		purchased                   time.Time
		utc                         *int
	)

	vehicles := make([]models.Vehicle, 0)
	for resp.Next() {
		err = resp.Scan(&id, &modelID, &enterpriseID, &price, &year, &mileage, &color, &vin, &purchased, &utc)

		if err != nil {
			fmt.Fprintf(os.Stderr, "scan failed: %v\n", err)
			return []models.Vehicle{}
		}

		loc := time.UTC
		if utc != nil {
			loc = time.FixedZone(fmt.Sprintf("UTC %v", *utc), *utc*60*60)
		}

		vehicles = append(vehicles, models.Vehicle{
			ID:              id,
			ModelID:         modelID,
			EnterpriseID:    enterpriseID,
			Price:           price,
			ManufactureYear: year,
			Mileage:         mileage,
			Color:           color,
			VIN:             vin,
			PurchasedAt:     purchased.In(loc).Format(time.RFC3339),
		})
	}

	return vehicles
}

func (p *Provider) FetchAllByManagerID(ctx context.Context, managerID int64) []models.Vehicle {
	query := `SELECT v.*
FROM enterprise as e
JOIN manager_enterprise as me on me.enterprise_id = e.id
JOIN vehicle as v on v.enterprise_id = e.id
WHERE manager_id = $1
`
	resp, err := p.db.Query(ctx, query, managerID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "query failed: %v\n", err)
		return []models.Vehicle{}
	}

	var (
		id, modelID                 int64
		enterpriseID                *int64
		price, year, mileage, color int
		vin                         string
		purchased, created, updated time.Time
		deleted                     *time.Time
	)

	vehicles := make([]models.Vehicle, 0)
	for resp.Next() {
		err = resp.Scan(&id, &modelID, &enterpriseID, &price, &year, &mileage, &color, &vin, &purchased, &created, &updated, &deleted)
		if deleted != nil {
			rounded := (*deleted).Round(timeRound)
			deleted = &rounded
		}

		if err != nil {
			fmt.Fprintf(os.Stderr, "scan failed: %v\n", err)
			return []models.Vehicle{}
		}
		vehicles = append(vehicles, models.Vehicle{
			ID:              id,
			ModelID:         modelID,
			EnterpriseID:    enterpriseID,
			Price:           price,
			ManufactureYear: year,
			Mileage:         mileage,
			Color:           color,
			VIN:             vin,
			PurchasedAt:     purchased.String(),
			CreatedAt:       created.Round(timeRound),
			UpdatedAt:       updated.Round(timeRound),
			DeletedAt:       deleted,
		})
	}

	return vehicles
}

func (p *Provider) FetchAllByEnterpriseID(ctx context.Context, entID int64) []models.Vehicle {
	query := `SELECT DISTINCT v.* 
FROM enterprise as e
JOIN manager_enterprise as me on me.enterprise_id = e.id
JOIN vehicle as v on v.enterprise_id = e.id
WHERE me.enterprise_id = $1
ORDER BY v.id 
`
	resp, err := p.db.Query(ctx, query, entID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "query failed: %v\n", err)
		return []models.Vehicle{}
	}

	var (
		id, modelID                 int64
		enterpriseID                *int64
		price, year, mileage, color int
		vin                         string
		purchased, created, updated time.Time
		deleted                     *time.Time
	)

	vehicles := make([]models.Vehicle, 0)
	for resp.Next() {
		err = resp.Scan(&id, &modelID, &enterpriseID, &price, &year, &mileage, &color, &vin, &purchased, &created, &updated, &deleted)
		if deleted != nil {
			rounded := (*deleted).Round(timeRound)
			deleted = &rounded
		}

		if err != nil {
			fmt.Fprintf(os.Stderr, "scan failed: %v\n", err)
			return []models.Vehicle{}
		}
		vehicles = append(vehicles, models.Vehicle{
			ID:              id,
			ModelID:         modelID,
			EnterpriseID:    enterpriseID,
			Price:           price,
			ManufactureYear: year,
			Mileage:         mileage,
			Color:           color,
			VIN:             vin,
			PurchasedAt:     purchased.String(),
			CreatedAt:       created.Round(timeRound),
			UpdatedAt:       updated.Round(timeRound),
			DeletedAt:       deleted,
		})
	}

	return vehicles
}

func (p *Provider) GetVehicleReportDaily(ctx context.Context, id int64,
	start, end int64,
) []models.VehicleReport {

	query := `SELECT v.vin, tr.started_at, tr.ended_at, tr.track_length, tr.max_velocity, tr.max_acceleration
	FROM trip as tr
	JOIN vehicle as v ON v.id = tr.vehicle_id
	WHERE tr.vehicle_id = $1 AND 
      EXTRACT(EPOCH FROM tr.started_at) >= $2 AND 
      EXTRACT(EPOCH FROM tr.ended_at) <= $3`

	resp, err := p.db.Query(ctx, query, id, start, end)
	if err != nil {
		fmt.Fprintf(os.Stderr, "query failed: %v\n", err)
		return []models.VehicleReport{}
	}

	var (
		vehicleVIN                     string
		startTime, endTime             time.Time
		length, velocity, acceleration int
	)

	type dayReports struct {
		started, ended time.Time
		mileage        int
	}

	dr := make([]dayReports, 0, 0)

	var dayOne *dayReports

	i := 0
	for resp.Next() {
		err = resp.Scan(&vehicleVIN, &startTime, &endTime, &length, &velocity, &acceleration)
		if err != nil {
			continue
		}

		if dayOne == nil {
			dayOne = &dayReports{started: startTime, ended: endTime, mileage: length}
			dr = append(dr, *dayOne)
			continue
		}

		if endTime.Before(dayOne.started.Add(time.Hour * 24)) {
			dr[i].mileage += length
		} else {
			dayOne.ended = endTime
			dayOne = nil
			i++
		}
	}

	reports := make([]models.VehicleReport, 0)
	for _, d := range dr {
		reports = append(reports, models.VehicleReport{
			ID:      time.Now().Unix(),
			Mileage: d.mileage,
			Report: models.Report{
				Title:     fmt.Sprintf("Отчет о поездках машины VIN: %s", vehicleVIN),
				StartDate: d.started,
				EndDate:   d.ended,
				Result:    "Груз доставлен",
				Type:      0,
			},
		})
	}

	return reports
}
