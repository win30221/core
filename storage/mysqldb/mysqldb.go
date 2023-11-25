package mysqldb

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/win30221/core/http/ctx"
	"github.com/win30221/core/http/middleware"
)

func restoreSQLLog(ctx *gin.Context, query string, args ...any) {
	res := query
	for _, arg := range args {
		switch arg.(type) {
		case string:
			res = strings.Replace(res, "?", fmt.Sprintf("'%v'", arg), 1)
		default:
			res = strings.Replace(res, "?", fmt.Sprintf("%v", arg), 1)
		}
	}
	if logs, exists := ctx.Get(middleware.SQLLogs); exists {
		ctx.Set(middleware.SQLLogs, append(logs.([]string), res))
	}
}

func QueryRowContext(ctx ctx.Context, db *sql.DB, query string, args ...any) (res *sql.Row) {
	restoreSQLLog(ctx.GinContext, query, args...)
	res = db.QueryRowContext(ctx.Context, query, args...)
	return
}

func QueryContext(ctx ctx.Context, db *sql.DB, query string, args ...any) (res *sql.Rows, err error) {
	restoreSQLLog(ctx.GinContext, query, args...)
	res, err = db.QueryContext(ctx.Context, query, args...)
	return
}

func ExecContext(ctx ctx.Context, db *sql.DB, query string, args ...any) (res sql.Result, err error) {
	restoreSQLLog(ctx.GinContext, query, args...)
	res, err = db.ExecContext(ctx.Context, query, args...)
	return
}

func BeginTx(ctx ctx.Context, db *sql.DB) (tx *sql.Tx, err error) {
	restoreSQLLog(ctx.GinContext, "BEGIN;")
	tx, err = db.BeginTx(ctx.Context, nil)
	return
}

func Commit(ctx ctx.Context, tx *sql.Tx) (err error) {
	restoreSQLLog(ctx.GinContext, "COMMIT;")
	err = tx.Commit()
	return
}

func Rollback(ctx ctx.Context, tx *sql.Tx) {
	restoreSQLLog(ctx.GinContext, "ROLLBACK;")
	tx.Rollback()
}

func QueryRowContextTx(ctx ctx.Context, db *sql.Tx, query string, args ...any) (res *sql.Row) {
	restoreSQLLog(ctx.GinContext, query, args...)
	res = db.QueryRowContext(ctx.Context, query, args...)
	return
}

func QueryContextTx(ctx ctx.Context, db *sql.Tx, query string, args ...any) (res *sql.Rows, err error) {
	restoreSQLLog(ctx.GinContext, query, args...)
	res, err = db.QueryContext(ctx.Context, query, args...)
	return
}

func ExecContextTx(ctx ctx.Context, db *sql.Tx, query string, args ...any) (res sql.Result, err error) {
	restoreSQLLog(ctx.GinContext, query, args...)
	res, err = db.ExecContext(ctx.Context, query, args...)
	return
}
