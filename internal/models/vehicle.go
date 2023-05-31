package models

import "time"

type Vehicle struct {
	ID              int64      `json:"id"`
	ModelID         int64      `json:"modelId"`
	Price           int        `json:"price"`
	ManufactureYear int        `json:"year"`
	Mileage         int        `json:"mileage"`
	Color           int        `json:"colorId"`
	VIN             string     `json:"vin"`
	CreatedAt       time.Time  `json:"created"`
	UpdatedAt       time.Time  `json:"updated"`
	DeletedAt       *time.Time `json:"deleted"`
}
