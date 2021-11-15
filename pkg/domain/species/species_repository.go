package species

import (
	"context"

	"go.uber.org/zap"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/luanapp/gin-example/config/database"
	"github.com/luanapp/gin-example/pkg/logger"
)

type (
	repository struct {
		conn *pgxpool.Pool
	}
)

const (
	insertSpSQL = `
INSERT INTO natural_history_museum.species
(id, scientific_name, genus, family, "order", class, phylum, kingdom)
VALUES
($1, $2, $3, $4, $5, $6, $7, $8);
`
	updateSpSQL = `
UPDATE natural_history_museum.species
SET scientific_name = $1, genus = $2,
    family = $3, "order" = $4, class = $5,
    phylum = $6, kingdom = $7
WHERE id = $8;
`
	deleteSpSQL = `
DELETE FROM natural_history_museum.species
WHERE id = $1;
`

	getByIdSQL = `
SELECT
	sp.id, sp.scientific_name, sp.genus, sp.family, sp."order", sp.class, sp.phylum, sp.kingdom
FROM
	natural_history_museum.species sp WHERE sp.id = $1;
`
)

var (
	sugar *zap.SugaredLogger
)

func init() {
	sugar = logger.New()
}

func DefaultRepository() *repository {
	return &repository{
		conn: database.GetConnection(),
	}
}

func (r *repository) getAll() ([]Species, error) {
	sql := "SELECT sp.id, sp.scientific_name, sp.family FROM natural_history_museum.species sp;"
	rows, err := r.conn.Query(context.Background(), sql)
	if err != nil {
		sugar.Errorw("failed to get species from database", "error", err.Error())
		return nil, err
	}
	defer rows.Close()

	sps := make([]Species, 0)
	for rows.Next() {
		var sp Species
		err = rows.Scan(&sp.Id, &sp.ScientificName, &sp.Family)
		if err != nil {
			sugar.Errorw("failed to get species from database", "error", err.Error())
			return nil, err
		}
		sps = append(sps, sp)
	}

	return sps, nil
}

func (r *repository) getById(id string) (*Species, error) {
	rows := r.conn.QueryRow(context.Background(), getByIdSQL, id)

	var sp Species
	err := rows.Scan(&sp.Id, &sp.ScientificName, &sp.Genus, &sp.Family, &sp.Order, &sp.Class, &sp.Phylum, &sp.Kingdom)
	if err != nil {
		sugar.Errorw("failed to get species by id from database", "error", err.Error())
		return nil, err
	}

	return &sp, nil
}

func (r *repository) save(sp *Species) (*Species, error) {
	sp.Id = uuid.NewString()

	_, err := r.conn.Query(context.Background(), insertSpSQL,
		sp.Id, sp.ScientificName, sp.Genus, sp.Family, sp.Order, sp.Class, sp.Phylum, sp.Kingdom)
	if err != nil {
		sugar.Errorw("failed to save species into database", "error", err.Error())
		return nil, err
	}

	if err != nil {
		sugar.Errorw("failed to save species into database", "error", err.Error())
		return nil, err
	}

	sp, _ = r.getById(sp.Id)
	return sp, nil
}

func (r *repository) update(sp *Species) error {
	_, err := r.conn.Exec(context.Background(), updateSpSQL, sp.ScientificName, sp.Genus, sp.Family,
		sp.Order, sp.Class, sp.Phylum, sp.Kingdom, sp.Id)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) delete(id string) error {
	_, err := r.conn.Exec(context.Background(), deleteSpSQL, id)
	if err != nil {
		return err
	}
	return nil
}
