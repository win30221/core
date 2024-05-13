package mysqldb

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/win30221/core/http/ctx"
	"github.com/win30221/core/http/middleware"
)

func buildSQLLog(ctx *gin.Context, query string, args ...any) {
	res := query
	for _, arg := range args {
		switch v := arg.(type) {
		case string:
			res = strings.Replace(res, "?", fmt.Sprintf("'%s'", v), 1)
		case time.Time:
			res = strings.Replace(res, "?", fmt.Sprintf("'%s'", v.Format("2006-01-02 15:04:05")), 1)
		default:
			res = strings.Replace(res, "?", fmt.Sprintf("%v", arg), 1)
		}
	}
	logs, exists := ctx.Get(middleware.SQLLogs)
	if !exists {
		ctx.Set(middleware.SQLLogs, []string{})
		logs, _ = ctx.Get(middleware.SQLLogs)
	}
	ctx.Set(middleware.SQLLogs, append(logs.([]string), res))
}

func QueryRowContext(ctx *ctx.Context, db *sql.DB, query string, args ...any) (res *sql.Row) {
	buildSQLLog(ctx.GinContext, query, args...)
	res = db.QueryRowContext(ctx.Context, query, args...)
	return
}

func QueryContext(ctx *ctx.Context, db *sql.DB, query string, args ...any) (res *sql.Rows, err error) {
	buildSQLLog(ctx.GinContext, query, args...)
	res, err = db.QueryContext(ctx.Context, query, args...)
	return
}

func ExecContext(ctx *ctx.Context, db *sql.DB, query string, args ...any) (res sql.Result, err error) {
	buildSQLLog(ctx.GinContext, query, args...)
	res, err = db.ExecContext(ctx.Context, query, args...)
	return
}

func BeginTx(ctx *ctx.Context, db *sql.DB) (tx *sql.Tx, err error) {
	buildSQLLog(ctx.GinContext, "BEGIN;")
	tx, err = db.BeginTx(ctx.Context, nil)
	return
}

func Commit(ctx *ctx.Context, tx *sql.Tx) (err error) {
	buildSQLLog(ctx.GinContext, "COMMIT;")
	err = tx.Commit()
	return
}

func Rollback(ctx *ctx.Context, tx *sql.Tx) {
	buildSQLLog(ctx.GinContext, "ROLLBACK;")
	tx.Rollback()
}

func QueryRowContextTx(ctx *ctx.Context, db *sql.Tx, query string, args ...any) (res *sql.Row) {
	buildSQLLog(ctx.GinContext, query, args...)
	res = db.QueryRowContext(ctx.Context, query, args...)
	return
}

func QueryContextTx(ctx *ctx.Context, db *sql.Tx, query string, args ...any) (res *sql.Rows, err error) {
	buildSQLLog(ctx.GinContext, query, args...)
	res, err = db.QueryContext(ctx.Context, query, args...)
	return
}

func ExecContextTx(ctx *ctx.Context, db *sql.Tx, query string, args ...any) (res sql.Result, err error) {
	buildSQLLog(ctx.GinContext, query, args...)
	res, err = db.ExecContext(ctx.Context, query, args...)
	return
}

// ===================================================

// BuildInCondition 函數用於生成 SQL 的 IN 條件子句，可用於多種查詢中。
// 它接收欄位名稱和對應的值陣列，返回構建的條件子句和對應的參數陣列。
func BuildInCondition(field string, values []any, conditions *[]string, args *[]any) {
	if len(values) == 0 {
		return
	}

	placeholders := make([]string, len(values))
	for i, value := range values {
		placeholders[i] = "?"        // 為每個值創建一個佔位符
		*args = append(*args, value) // 將值添加到參數列表中
	}
	*conditions = append(*conditions, "`"+field+"` IN ("+strings.Join(placeholders, ", ")+") ")
}

// BuildOrder 函數根據提供的排序欄位和排序順序陣列生成 SQL 的 ORDER BY 子句。
// 它接收兩個字串陣列：sort 包含排序欄位名，order 包含對應的排序順序（如 "ASC" 或 "DESC"）。
// 返回構建好的 ORDER BY 子句字串。
func BuildOrder(sort []string, order []string) (res string) {
	if len(sort) == 0 || len(sort) != len(order) {
		return "" // 如果沒有排序欄位或排序欄位與順序陣列長度不匹配，返回空字串
	}

	var orderClauses []string
	for i, field := range sort {
		orderClause := field + " " + order[i] // 組合欄位名和排序順序
		orderClauses = append(orderClauses, orderClause)
	}

	res = " ORDER BY " + strings.Join(orderClauses, ", ") // 連接所有排序子句，形成完整的 ORDER BY 子句
	return
}

// GetScanFields 函數用於從任何指定的結構體動態生成用於資料庫查詢結果的掃描接收器。
// 此函數接收任意類型的 data 參數，返回一個包含每個欄位接收器的切片。
func GetScanFields(data any) (res []any) {
	val := reflect.ValueOf(data).Elem() // 確保傳入的是結構體的指針

	// 遍歷結構體中的每個欄位
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		res = append(res, field.Addr().Interface())
	}

	return
}

// GetInsertFields 函數用於從任何給定的結構體動態生成 INSERT 語句所需的欄位名稱、佔位符和參數值。
// 此函數接受任意類型的 data 參數，返回三個切片：欄位名稱、佔位符和參數值。
func GetInsertFields(data any) (columns []string, placeholders []string, args []interface{}) {
	val := reflect.ValueOf(data) // 取得 data 的反射值物件
	typ := reflect.TypeOf(data)  // 取得 data 的類型資訊

	// 遍歷結構體中的每個欄位
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)     // 取得欄位的反射值物件
		fieldType := typ.Field(i) // 取得當前欄位的類型資訊

		// 如果當前欄位是 ID 或者為 nil（對於指針類型），則跳過不處理
		if fieldType.Name == "ID" || (field.Kind() == reflect.Ptr && field.IsNil()) {
			continue
		}

		fieldName := strings.ToLower(fieldType.Name)  // 將欄位名稱轉為小寫
		columns = append(columns, fieldName)          // 加入欄位名稱至 columns 切片
		placeholders = append(placeholders, "?")      // 加入對應的 SQL 佔位符 "?"
		args = append(args, field.Elem().Interface()) // 加入當前欄位的值至 args 切片
	}
	return
}

// GetNonNullFields 函數用於檢查結構體中的非 nil 指針欄位，並為這些欄位生成 SQL 更新語句的 SET 子句部分。
// 此函數接受任意類型的 data 參數，返回兩個切片：一個包含 SET 子句的字符串，另一個包含相應的參數值。
func GetNonNullFields(data any) (res []string, args []any) {
	val := reflect.ValueOf(data) // 取得 data 的反射值物件
	typ := reflect.TypeOf(data)  // 取得 data 的類型資訊

	// 遍歷結構體中的每個欄位
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)     // 取得欄位的反射值物件
		fieldType := typ.Field(i) // 取得當前欄位的類型資訊

		// 如果當前欄位是 ID 或者為 nil（對於指針類型），則跳過不處理
		if fieldType.Name == "ID" {
			continue
		}

		// 如果欄位是指標類型且不為 nil，則將其包括在 SET 子句中
		if field.Kind() == reflect.Ptr && !field.IsNil() {
			fieldName := strings.ToLower(fieldType.Name)        // 將欄位名稱轉換為小寫
			res = append(res, fmt.Sprintf("%s = ?", fieldName)) // 構建 SET 子句，使用欄位名並指定值將由參數提供
			args = append(args, field.Elem().Interface())       // 添加當前欄位的值到參數切片中
		}
	}

	return
}
