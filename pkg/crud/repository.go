package crud

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"go.uber.org/zap"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/luanapp/gin-example/config/database"
	"github.com/luanapp/gin-example/pkg/logger"
	"github.com/luanapp/gin-example/pkg/model"
)

type (
	Repository[T model.Model] interface {
		GetAll() ([]T, error)
		GetById(id string) (*T, error)
		Save(entity *T) (*T, error)
		Update(entity *T) error
		Delete(id string) error
	}

	repository[T model.Model] struct {
		conn *pgxpool.Pool
	}
)

const (
	dbTag        = "db"
	selectAllSQL = `
SELECT %s FROM natural_history_museum.%s %s;
`
	getByIdSQL = `
SELECT %s FROM natural_history_museum.%s %s WHERE sp.id = $1;
`
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
)

var (
	sugar      *zap.SugaredLogger
	tableNames = map[reflect.Type][]string{
		reflect.TypeOf(model.Species{}): {"species", "sp"},
	}
)

func init() {
	sugar = logger.New()
}

func defaultRepository[T model.Model]() Repository[T] {
	return &repository[T]{
		conn: database.GetConnection(),
	}
}

func (r *repository[T]) GetAll() ([]T, error) {
	fields, tableName, tablePrefix := r.getAllSelectFields()
	selectAll := fmt.Sprintf(selectAllSQL, fields, tableName, tablePrefix)
	sugar.Debugf("select all: %s", selectAll)

	rows, err := r.conn.Query(context.Background(), selectAll)
	if err != nil {
		sugar.Errorw("failed to get data from database", "error", err.Error(), "table", tableNames)
		return nil, err
	}
	defer rows.Close()

	entities := make([]T, 0)
	for rows.Next() {
		var t T
		fieldsAddrs := r.getAllSelectFieldsAddrs(t)
		err = rows.Scan(fieldsAddrs...)
		if err != nil {
			sugar.Errorw("failed to get data from database", "error", err.Error(), "table", tableName)
			return nil, err
		}
		r.setFields(&t, fieldsAddrs)
		entities = append(entities, t)
	}

	return entities, nil
}

func (r *repository[T]) GetById(id string) (*T, error) {
	fields, tableName, tablePrefix := r.getAllSelectFields()
	selectById := fmt.Sprintf(getByIdSQL, fields, tableName, tablePrefix)
	sugar.Debugf("select by id: %s", selectById)

	rows := r.conn.QueryRow(context.Background(), selectById, id)

	var t T
	fieldsAddrs := r.getAllSelectFieldsAddrs(t)
	err := rows.Scan(fieldsAddrs...)
	if err != nil {
		sugar.Errorw("failed to get data by id from database", "error", err.Error(), "table", tableName, "id", id)
		return nil, err
	}
	r.setFields(&t, fieldsAddrs)

	return &t, nil
}

func (r *repository[T]) Save(e *T) (*T, error) {
	(*e).SetId(uuid.NewString())

	_, err := r.conn.Query(context.Background(), insertSpSQL) //e.Id, e.ScientificName, e.Genus, e.Family, e.Order, e.Class, e.Phylum, e.Kingdom)
	if err != nil {
		sugar.Errorw("failed to save species into database", "error", err.Error())
		return nil, err
	}

	if err != nil {
		sugar.Errorw("failed to save species into database", "error", err.Error())
		return nil, err
	}

	e, _ = r.GetById((*e).GetId())
	return e, nil
}

func (r *repository[T]) Update(sp *T) error {
	_, err := r.conn.Exec(context.Background(), updateSpSQL, nil) /*sp.ScientificName, sp.Genus, sp.Family,
	sp.Order, sp.Class, sp.Phylum, sp.Kingdom, sp.Id)*/
	if err != nil {
		return err
	}
	return nil
}

func (r *repository[T]) Delete(id string) error {
	_, err := r.conn.Exec(context.Background(), deleteSpSQL, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository[T]) getAllSelectFields() (string, string, string) {
	var e T
	var fieldsNames []string

	v := reflect.ValueOf(e)
	for i := 0; i < v.NumField(); i++ {
		tag := v.Type().Field(i).Tag.Get(dbTag)
		if tag == "" {
			continue
		}

		fieldsNames = append(fieldsNames, tag)
	}

	var selectFields strings.Builder
	tableData := tableNames[v.Type()]
	for _, field := range fieldsNames {
		selectFields.WriteString(tableData[1])
		selectFields.WriteString(".\"")
		selectFields.WriteString(field)
		selectFields.WriteString("\", ")
	}

	return strings.TrimRight(selectFields.String(), ", "), tableData[0], tableData[1]
}

func (r *repository[T]) getAllSelectFieldsAddrs(e T) []any {
	var fieldsAddrs []any
	v := reflect.ValueOf(&e)

	for i := 0; i < v.Elem().NumField(); i++ {
		tag := v.Elem().Type().Field(i).Tag.Get(dbTag)
		if tag == "" {
			continue
		}

		addr := v.Elem().Field(i).Addr().Interface()
		fieldsAddrs = append(fieldsAddrs, addr)
	}
	return fieldsAddrs
}

func (r *repository[T]) setFields(e *T, addrs []any) {
	v := reflect.ValueOf(e)

	for i := 0; i < v.Elem().NumField(); i++ {
		tag := v.Elem().Type().Field(i).Tag.Get(dbTag)
		if tag == "" {
			continue
		}

		elem := reflect.ValueOf(addrs[i]).Elem()
		sugar.Debugf("setting field %s with value %s", reflect.TypeOf(*e).Field(i).Name, elem.Interface())
		v.Elem().Field(i).Set(elem)
	}
}
