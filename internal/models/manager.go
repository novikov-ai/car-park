package models

type Manager struct {
	ID       int64  `json:"id"`
	FullName string `json:"name"`
	Age      int    `json:"age"`
	Salary   int    `json:"salary"`
}
