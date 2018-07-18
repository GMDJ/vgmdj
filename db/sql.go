package db

import (
	"fmt"
	"github.com/vgmdj/utils/chars"
	"github.com/vgmdj/utils/logger"
	"strconv"
	"strings"
)

type OP string

const (
	EQ   OP = "="
	LT   OP = "<"
	LE   OP = "<="
	NE   OP = "!="
	GT   OP = ">"
	GE   OP = ">="
	LIKE OP = "like"
)

func AttachOr(sql string, query interface{}, data interface{}, op OP) string {
	return attach(sql, query, data, op, "or")
}

func AttachAnd(sql string, query interface{}, data interface{}, op OP) string {
	return attach(sql, query, data, op, "and")
}

func Attach(sql string, query interface{}, data interface{}, op OP) string {
	return attach(sql, query, data, op, " ")
}

func Limit(sql string, pageCount, pageIndex interface{}) string {
	var count, index int

	switch pageCount.(type) {
	default:
		logger.Error("invalid type ")
		return sql
	case string:
		count, _ = strconv.Atoi(pageCount.(string))
	}

	switch pageIndex.(type) {
	default:
		logger.Error("invalid type ")
		return sql
	case string:
		index, _ = strconv.Atoi(pageIndex.(string))
	}

	limit, offset := LimitQuery(count, index)

	return fmt.Sprintf(" %s limit %d, %d ", sql, limit, offset)

}

func LimitQuery(pageCount, pageIndex int) (limit int, offset int) {
	limit = pageCount
	offset = limit * (pageIndex - 1)

	if limit == 0 {
		return 0, 0
	}

	if offset <= 0 {
		return limit, 0
	}

	return limit, offset

}

func attach(sql string, query interface{}, data interface{}, op OP, relation string) string {
	if data == "" || !checkOp(op) {
		return sql
	}

	sql = strings.Replace(sql, "\t", " ", -1)
	sql = strings.Replace(sql, "\n", " ", -1)
	sql = strings.Replace(sql, "\r", " ", -1)

	if strings.Contains(strings.ToLower(sql), " where ") || sql == "" {
		sql = fmt.Sprintf(" %s %s %v %v '%v' ", sql, relation, query, op, data)
		return sql
	}

	sql = fmt.Sprintf(" %s where %s %v %v '%v' ", sql, relation, query, op, data)
	return sql
}

func checkOp(op OP) bool {
	ops := []interface{}{EQ, LT, LE, NE, GT, GE}

	return chars.IsContain(ops, op)
}
