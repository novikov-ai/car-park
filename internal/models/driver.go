package models

type Driver struct {
	ID           int64  `json:"id"`
	EnterpriseID int64  `json:"enterpriseId"`
	ActiveCarID  *int64 `json:"activeCarId"`
	Age          int    `json:"age"`
	Salary       int    `json:"salary"`
	Experience   int    `json:"experience"`
}
