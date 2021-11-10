package species

import (
	"context"

	"github.com/jackc/pgx/v4"
	"luana.com/gin-example/config/database"
)

type (
	repository struct {
		conn *pgx.Conn
	}
)

func DefaultRepository() *repository {
	return &repository{
		conn: database.GetConnection(),
	}
}

func (r *repository) getAll() ([]Species, error) {
	rows, err := r.conn.Query(context.Background(), "SELECT sp.id, sp.scientific_name, sp.family FROM natural_history_museum.species sp")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sps := make([]Species, 0)
	for rows.Next() {
		var sp Species
		err = rows.Scan(&sp.Id, &sp.ScientificName, &sp.Family)
		if err != nil {
			return nil, err
		}
		sps = append(sps, sp)
	}

	return sps, nil
}
