package vehicles

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"car-park/internal/models"
	"car-park/internal/storage"
)

const timeRound = time.Second

const (
	idField      = "id"
	modelField   = "model"
	priceField   = "price"
	yearField    = "year"
	mileageField = "mileage"
	colorField   = "color"
	vinField     = "vin"
)

const redirectPath = "/view/vehicles/"

type Provider struct {
	db storage.Client
}

func New(st storage.Client) *Provider {
	return &Provider{
		db: st,
	}
}

func (p *Provider) Create(c *gin.Context) {
	defer redirectToView(c)

	err := c.Request.ParseForm()
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	modelValue := c.Request.FormValue(modelField)
	colorValue := c.Request.FormValue(colorField)

	toInt := func(key string) (int, error) {
		return strconv.Atoi(c.Request.Form.Get(key))
	}

	modelID, err := strconv.ParseInt(modelValue, 10, 64)
	colorID, err := toInt(colorValue)
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

	query := `INSERT INTO vehicle (model_id, price, manufacture_year, mileage, color, vin)
    VALUES
        ($1, $2, $3, $4, $5, $6)`

	resp, err := p.db.Query(context.Background(), query,
		vehicle.ModelID, vehicle.Price, vehicle.ManufactureYear, vehicle.Mileage, vehicle.Color, vehicle.VIN)

	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	resp.Close()
}

func (p *Provider) Update(c *gin.Context) {
	defer redirectToView(c)

	err := c.Request.ParseForm()
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	modelID, err := strconv.ParseInt(c.Request.Form.Get(modelField), 10, 64)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	toInt := func(key string) (int, error) {
		return strconv.Atoi(c.Request.Form.Get(key))
	}

	id, err := strconv.ParseInt(c.Request.Form.Get(idField), 10, 64)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	price, err := toInt(priceField)
	year, err := toInt(yearField)
	mileage, err := toInt(mileageField)
	color, err := toInt(colorField)

	vehicle := models.Vehicle{
		ModelID:         modelID,
		Price:           price,
		ManufactureYear: year,
		Mileage:         mileage,
		Color:           color,
		VIN:             c.Request.Form.Get(vinField),
	}

	query := `UPDATE vehicle 
SET model_id = $1, price = $2, manufacture_year = $3,
    mileage = $4, color = $5, vin = $6
WHERE id = $7`

	resp, err := p.db.Query(context.Background(), query,
		vehicle.ModelID, vehicle.Price, vehicle.ManufactureYear, vehicle.Mileage, vehicle.Color, vehicle.VIN, id)

	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	resp.Close()
}

func (p *Provider) Delete(c *gin.Context) {
	defer redirectToView(c)

	err := c.Request.ParseForm()
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	b := c.Request.Form
	print(b)

	id, err := strconv.ParseInt(c.Request.Form.Get(idField), 10, 64)
	if err != nil {
		// toodo: handle
		return
	}

	query := `DELETE FROM vehicle
WHERE id=$1`

	resp, err := p.db.Query(context.Background(), query, id)

	resp.Close()
}

func (p *Provider) FetchAll(ctx context.Context) []models.Vehicle {
	query := `SELECT * FROM vehicle`
	resp, err := p.db.Query(ctx, query)
	if err != nil {
		fmt.Fprintf(os.Stderr, "query failed: %v\n", err)
		return []models.Vehicle{}
	}

	var (
		id, modelID                 int64
		enterpriseID                *int64
		price, year, mileage, color int
		vin                         string
		created, updated            time.Time
		deleted                     *time.Time
	)

	vehicles := make([]models.Vehicle, 0)
	for resp.Next() {
		err = resp.Scan(&id, &modelID, &enterpriseID, &price, &year, &mileage, &color, &vin, &created, &updated, &deleted)
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
			CreatedAt:       created.Round(timeRound),
			UpdatedAt:       updated.Round(timeRound),
			DeletedAt:       deleted,
		})
	}

	return vehicles
}

func redirectToView(c *gin.Context) {
	c.Redirect(http.StatusFound, redirectPath)
}
