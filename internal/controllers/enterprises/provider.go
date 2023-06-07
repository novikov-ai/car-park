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

func New(st storage.Client) *Provider {
	return &Provider{
		db: st,
	}
}

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
	)

	enterprises := make([]models.Enterprise, 0)
	for resp.Next() {
		err = resp.Scan(&id, &title, &city, &established)
		if err != nil {
			fmt.Fprintf(os.Stderr, "scan failed: %v\n", err)
			return []models.Enterprise{}
		}

		enterprises = append(enterprises, models.Enterprise{
			ID:          id,
			Title:       title,
			City:        city,
			Established: established.Round(time.Second),
		})
	}

	return enterprises
}

func (p *Provider) FetchManagersOnly(ctx *gin.Context, managerID int64) []models.Enterprise {
	query := `SELECT id, title, city, established
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
		established time.Time
	)

	enterprises := make([]models.Enterprise, 0)
	for resp.Next() {
		err = resp.Scan(&id, &title, &city, &established)
		if err != nil {
			fmt.Fprintf(os.Stderr, "scan failed: %v\n", err)
			return []models.Enterprise{}
		}

		enterprises = append(enterprises, models.Enterprise{
			ID:          id,
			Title:       title,
			City:        city,
			Established: established.Round(time.Second),
		})
	}

	return enterprises
}
