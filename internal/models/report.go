package models

import "time"

type Report struct {
	ID        int64
	Mileage   int
	Title     string
	Period    time.Duration
	StartDate time.Time
	EndDate   time.Time
	Result    string
	Type      int
}

type VehicleReport struct {
	ID      int64
	Mileage int
	Report
}
