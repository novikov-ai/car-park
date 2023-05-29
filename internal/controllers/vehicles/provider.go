package vehicles

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"car-park/internal/models"
	"car-park/internal/storage"
)

type Provider struct {
	db storage.Client
}

func New(st storage.Client) *Provider {
	return &Provider{
		db: st,
	}
}

func (p *Provider) Create(c *gin.Context) {
	err := c.Request.ParseForm()
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	modelID, err := strconv.ParseInt(c.Request.Form.Get("model"), 10, 64)
	if err != nil {
		// toodo: handle
		panic(err)
	}

	toInt := func(key string) (int, error) {
		return strconv.Atoi(c.Request.Form.Get(key))
	}

	price, err := toInt("price")
	year, err := toInt("year")
	mileage, err := toInt("mileage")
	color, err := toInt("color")

	vehicle := models.Vehicle{
		ModelID:         modelID,
		Price:           price,
		ManufactureYear: year,
		Mileage:         mileage,
		Color:           color,
		VIN:             c.Request.Form.Get("vin"),
	}

	query := `INSERT INTO vehicle (model_id, price, manufacture_year, mileage, color, vin)
    VALUES
        ($1, $2, $3, $4, $5, $6)`

	resp, err := p.db.Query(context.Background(), query,
		vehicle.ModelID, vehicle.Price, vehicle.ManufactureYear, vehicle.Mileage, vehicle.Color, vehicle.VIN)

	resp.Close()

	c.Redirect(http.StatusFound, "/view/vehicles/")
}

func (p *Provider) Update(c *gin.Context) {
	c.Redirect(http.StatusFound, "/view/vehicles/")

	err := c.Request.ParseForm()
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	modelID, err := strconv.ParseInt(c.Request.Form.Get("model"), 10, 64)
	if err != nil {
		// toodo: handle
		panic(err)
	}

	toInt := func(key string) (int, error) {
		return strconv.Atoi(c.Request.Form.Get(key))
	}

	id, err := strconv.ParseInt(c.Request.Form.Get("id"), 10, 64)
	if err != nil {
		// toodo: handle
		panic(err)
	}

	price, err := toInt("price")
	year, err := toInt("year")
	mileage, err := toInt("mileage")
	color, err := toInt("color")

	vehicle := models.Vehicle{
		ModelID:         modelID,
		Price:           price,
		ManufactureYear: year,
		Mileage:         mileage,
		Color:           color,
		VIN:             c.Request.Form.Get("vin"),
	}

	query := `UPDATE vehicle 
SET model_id = $1, price = $2, manufacture_year = $3,
    mileage = $4, color = $5, vin = $6
WHERE id = $7`

	resp, err := p.db.Query(context.Background(), query,
		vehicle.ModelID, vehicle.Price, vehicle.ManufactureYear, vehicle.Mileage, vehicle.Color, vehicle.VIN,
		id)

	resp.Close()
}

func (p *Provider) Delete(c *gin.Context) {
	defer c.Redirect(http.StatusFound, "/view/vehicles/")

	err := c.Request.ParseForm()
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	b := c.Request.Form
	print(b)

	id, err := strconv.ParseInt(c.Request.Form.Get("id"), 10, 64)
	if err != nil {
		// toodo: handle
		return
	}

	// todo: maybe just update delete field??
	query := `DELETE FROM vehicle
WHERE id=$1`

	resp, err := p.db.Query(context.Background(), query, id)

	resp.Close()
}

func (p *Provider) FetchAll(ctx context.Context) []models.Vehicle {
	query := `SELECT * FROM vehicle`
	resp, err := p.db.Query(ctx, query)
	if err != nil {
		panic(err)
	}

	var (
		id                          int64
		modelID                     int64
		price, year, mileage, color int
		vin                         string
		created, updated            time.Time
		deleted                     *time.Time
	)

	vehicles := make([]models.Vehicle, 0)
	for resp.Next() {
		err = resp.Scan(&id, &modelID, &price, &year, &mileage, &color, &vin, &created, &updated, &deleted)
		if err != nil {
			panic(err)
		}
		vehicles = append(vehicles, models.Vehicle{
			ID:              id,
			ModelID:         modelID,
			Price:           price,
			ManufactureYear: year,
			Mileage:         mileage,
			Color:           color,
			VIN:             vin,
			CreatedAt:       created,
			UpdatedAt:       updated,
			DeletedAt:       deleted,
		})
	}

	return vehicles
}
