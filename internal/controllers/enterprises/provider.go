package enterprises

import (
	"car-park/internal/models"
	"car-park/internal/storage"
	"context"
	"fmt"
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
		panic(err)
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
