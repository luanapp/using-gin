package crud

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/luanapp/gin-example/config/env"

	"go.uber.org/zap"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/luanapp/gin-example/config/database"
	"github.com/luanapp/gin-example/pkg/logger"
	"github.com/luanapp/gin-example/pkg/model"
)

type (
	Repository[T any] interface {
		GetAll() ([]T, error)
		GetById(id string) (*T, error)
		Save(entity *T) (*T, error)
		Update(id string, entity *T) error
		Delete(id string) error
	}

	repository[T any] struct {
		conn *pgxpool.Pool
	}
)

const (
	dbTag        = "db"
	selectAllSQL = `
SELECT %s FROM %s.%s %s;
`
	getByIdSQL = `
SELECT %s FROM %s.%s %s WHERE id = $1;
`
	insertSQL = `
INSERT INTO %s.%s (%s) VALUES (%s) RETURNING id;
`
	updateSQL = `
UPDATE %s.%s SET %s WHERE %s;
`
	updateWhere = `id = $%d`
	deleteSQL   = `
DELETE FROM %s.%s
WHERE id = $1;
`
)

var (
	schema     string
	sugar      *zap.SugaredLogger
	tableNames = map[reflect.Type][]string{
		reflect.TypeOf(model.Species{}): {"species", "sp"},
	}
)

func init() {
	sugar = logger.New()
	schema = env.Instance.Database.Schema
}

func defaultRepository[T any]() Repository[T] {
	return &repository[T]{
		conn: database.GetConnection(),
	}
}

func (r repository[T]) GetAll() ([]T, error) {
	fields, tableName, tablePrefix := r.getSelectAllFields()
	selectAll := fmt.Sprintf(selectAllSQL, fields, schema, tableName, tablePrefix)
	sugar.Debugf("select all: %s", selectAll)

	rows, err := r.conn.Query(context.Background(), selectAll)
	if err != nil {
		sugar.Errorw("failed to get data from database", "error", err.Error(), "table", tableNames)
		return nil, err
	}
	defer rows.Close()

	entities := make([]T, 0)
	for rows.Next() {
		t := new(T)
		fieldsAddrs := r.getVarFields(t, true, false)
		err = rows.Scan(fieldsAddrs...)
		if err != nil {
			sugar.Errorw("failed to get data from database", "error", err.Error(), "table", tableName)
			return nil, err
		}
		entities = append(entities, *t)
	}

	return entities, nil
}

func (r repository[T]) GetById(id string) (*T, error) {
	fields, tableName, tablePrefix := r.getSelectAllFields()
	selectById := fmt.Sprintf(getByIdSQL, fields, schema, tableName, tablePrefix)
	sugar.Debugf("select by id: %s", selectById)

	rows := r.conn.QueryRow(context.Background(), selectById, id)

	t := new(T)
	fieldsAddrs := r.getVarFields(t, true, false)
	err := rows.Scan(fieldsAddrs...)
	if err != nil {
		sugar.Errorw("failed to get data by id from database", "error", err.Error(), "table", tableName, "id", id)
		return nil, err
	}

	return t, nil
}

func (r repository[T]) Save(t *T) (*T, error) {
	table, fieldNames, fieldValues := r.getInsertionFields()
	insert := fmt.Sprintf(insertSQL, schema, table, fieldNames, fieldValues)
	fields := r.getVarFields(t, false, false)

	sugar.Debugf("insert: %s", insert)

	row := r.conn.QueryRow(context.Background(), insert, fields...)

	var id string
	err := row.Scan(&id)
	if err != nil {
		sugar.Errorw(fmt.Sprintf("failed to save %s into database", table), "error", err.Error())
		return nil, err
	}

	t, _ = r.GetById(id)
	return t, nil
}

func (r repository[T]) Update(id string, t *T) error {
	fields := r.getVarFields(t, false, true)
	tableName, fieldNames := r.getUpdateFields()

	where := fmt.Sprintf(updateWhere, len(fields)+1)
	update := fmt.Sprintf(updateSQL, schema, tableName, fieldNames, where)

	sugar.Debugf("update: %s", update)

	fields = append(fields, id)

	_, err := r.conn.Exec(context.Background(), update, fields...)
	if err != nil {
		return err
	}
	return nil
}

func (r repository[T]) Delete(id string) error {
	tableData := r.getTableData()
	_, err := r.conn.Exec(context.Background(), fmt.Sprintf(deleteSQL, schema, tableData[0]), id)
	if err != nil {
		return err
	}
	return nil
}

func (r repository[T]) getDBFieldsNames(ignorePk bool) []string {
	var (
		t           T
		fieldsNames []string
	)

	v := reflect.ValueOf(t)
	for i := 0; i < v.NumField(); i++ {
		tag := v.Type().Field(i).Tag.Get(dbTag)
		if tag == "" {
			continue
		}
		dbField, isPk := getTagParts(tag)
		if ignorePk && isPk {
			continue
		}

		fieldsNames = append(fieldsNames, dbField)
	}
	return fieldsNames
}

func getTagParts(tag string) (string, bool) {
	values := strings.Split(tag, ",")

	var isPk bool
	if len(values) > 1 && values[1] == "pk" {
		isPk = true
	}
	return values[0], isPk
}

func (r repository[T]) getSelectAllFields() (sel, tableName, tableAlias string) {
	fieldsNames := r.getDBFieldsNames(false)

	var selectFields strings.Builder
	tableData := r.getTableData()
	for _, field := range fieldsNames {
		selectFields.WriteString(tableData[1])
		selectFields.WriteString(".\"")
		selectFields.WriteString(field)
		selectFields.WriteString("\", ")
	}

	return strings.TrimRight(selectFields.String(), ", "), tableData[0], tableData[1]
}

func (r repository[T]) getTableData() []string {
	var t T
	return tableNames[reflect.ValueOf(t).Type()]
}

func (r repository[T]) getInsertionFields() (tableName, fields, fieldsValues string) {
	fieldsNames := r.getDBFieldsNames(true)

	var selectFields, valueFields strings.Builder
	tableData := r.getTableData()
	for i, field := range fieldsNames {
		// fields
		selectFields.WriteString("\"")
		selectFields.WriteString(field)
		selectFields.WriteString("\", ")

		//values
		valueFields.WriteString("$")
		valueFields.WriteString(strconv.Itoa(i + 1))
		valueFields.WriteString(", ")
	}

	return tableData[0], strings.TrimRight(selectFields.String(), ", "), strings.TrimRight(valueFields.String(), ", ")
}

func (r repository[T]) getUpdateFields() (tableName, fields string) {
	fieldsNames := r.getDBFieldsNames(true)

	var selectFields strings.Builder
	tableData := r.getTableData()
	for i, field := range fieldsNames {
		// field
		selectFields.WriteString("\"")
		selectFields.WriteString(field)
		selectFields.WriteString("\"")

		selectFields.WriteString(" = ")

		//value
		selectFields.WriteString("$")
		selectFields.WriteString(strconv.Itoa(i + 1))
		selectFields.WriteString(", ")
	}

	return tableData[0], strings.TrimRight(selectFields.String(), ", ")
}

func (r repository[T]) getVarFields(t *T, returnAddr, ignorePk bool) []any {
	var fields []any
	v := reflect.ValueOf(t)

	for i := 0; i < v.Elem().NumField(); i++ {
		tag := v.Elem().Type().Field(i).Tag.Get(dbTag)
		if tag == "" {
			continue
		}
		_, isPk := getTagParts(tag)
		if ignorePk && isPk {
			continue
		}

		var value any
		if returnAddr {
			value = v.Elem().Field(i).Addr().Interface()
		} else {
			value = v.Elem().Field(i).Interface()
		}
		fields = append(fields, value)
	}
	return fields
}
