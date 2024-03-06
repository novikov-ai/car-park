package enterprises

import (
	"car-park/internal/models"
	"car-park/internal/storage"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
	"time"
)

type Provider struct {
	db storage.Client
}

// Postcondition:
// Создан новый провайдер с указанным хранилищем
func New(st storage.Client) *Provider {
	return &Provider{
		db: st,
	}
}

// Postcondition:
// Из БД получены все предприятия
func (p *Provider) FetchAll(ctx context.Context) []models.Enterprise {
	query := `SELECT * FROM enterprise`
	resp, err := p.db.Query(ctx, query)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to proceed query: %v\n", err)
		return []models.Enterprise{}
	}

	var (
		id          int64
		title, city string
		established time.Time
		utc         int
	)

	enterprises := make([]models.Enterprise, 0)
	for resp.Next() {
		err = resp.Scan(&id, &title, &city, &established, &utc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "scan failed: %v\n", err)
			return []models.Enterprise{}
		}

		enterprises = append(enterprises, models.Enterprise{
			ID:          id,
			Title:       title,
			City:        city,
			Established: established.Round(time.Second),
			UTC:         utc,
		})
	}

	return enterprises
}

// Precondition:
// Менеджер существует в системе

// Postcondition:
// Получены все предприятия, принадлежащие менеджеру
func (p *Provider) FetchAllByManagerID(ctx *gin.Context, managerID int64) []models.Enterprise {
	query := `SELECT e.*
FROM enterprise as e
JOIN manager_enterprise as me on me.enterprise_id = e.id
WHERE me.manager_id = $1;
`

	resp, err := p.db.Query(ctx, query, managerID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to proceed query: %v\n", err)
		return []models.Enterprise{}
	}

	var (
		id          int64
		title, city string
		utc         int
		established time.Time
	)

	enterprises := make([]models.Enterprise, 0)
	for resp.Next() {
		err = resp.Scan(&id, &title, &city, &established, &utc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "scan failed: %v\n", err)
			return []models.Enterprise{}
		}

		enterprises = append(enterprises, models.Enterprise{
			ID:          id,
			Title:       title,
			City:        city,
			Established: established.Round(time.Second),
			UTC:         utc,
		})
	}

	return enterprises
}
