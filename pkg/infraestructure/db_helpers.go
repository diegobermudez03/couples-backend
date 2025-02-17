package infraestructure

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)


type Scanable interface{
	Scan(dest ...any) error
}

func getDBContext(ctx context.Context, fallback Executor) Executor {
	tx := ctx.Value(dbKey{})
	if tx == nil {
		return fallback
	}
	return tx.(*sql.Tx)
}

func ExecSQL(ctx context.Context, db *sql.DB, sql func(ex Executor) (sql.Result, error)) (int, error) {
	executor := getDBContext(ctx, db)
	result, err := sql(executor)
	if err != nil {
		return 0, err
	}
	num, _ := result.RowsAffected()
	return int(num), nil
}

func GetFilteredQuery(baseQuery string, filters map[string]any) (string, []any) {
	builder := strings.Builder{}
	builder.WriteString(baseQuery)
	counter := 1
	args := []any{}
	for k, v := range filters {
		if reflect.TypeOf(v).Kind() == reflect.Ptr && !reflect.ValueOf(v).IsNil(){
			builder.WriteString(fmt.Sprintf(" AND $%d = %s ", counter, k))
			counter++
			args = append(args, v)
		}
	}
	return builder.String(), args
}


func GetFetchingQuery(baseQuery string, counter int, limit int, page *int) (string, []any){
	builder := strings.Builder{}
	builder.WriteString(baseQuery)
	counter++
	args := []any{}
	builder.WriteString(fmt.Sprintf(" LIMIT $%d ", counter))
	args = append(args, limit)
	counter++
	if page != nil{
		offset := (*page) * limit - (*page)
		builder.WriteString(fmt.Sprintf(" OFFSET $%d", counter))
		counter++
		args = append(args, offset)
	}
	return builder.String(), args
}