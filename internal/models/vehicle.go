package models

import "time"

type Vehicle struct {
	ID              int64
	ModelID         int64
	Price           int
	ManufactureYear int
	Mileage         int
	Color           int
	VIN             string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       *time.Time
}
