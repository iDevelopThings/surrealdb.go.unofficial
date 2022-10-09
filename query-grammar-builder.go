package surrealdb

import (
	"strconv"
	"strings"
)

type QueryGrammarBuilder[T any] struct {
	builder *QueryBuilder[T]
	query   string
}

func BuildQueryGrammar[T any](builder *QueryBuilder[T]) *QueryGrammarBuilder[T] {
	return &QueryGrammarBuilder[T]{
		builder: builder,
	}
}

func (q *QueryGrammarBuilder[T]) BuildSelects() string {
	selects := ""

	if len(q.builder.selections) == 0 {
		selects = "*"
	} else {
		for _, selection := range q.builder.selections {
			selects += selection.key
			if selection.as != nil {
				selects += " AS " + *selection.as
			}
			selects += ", "
		}

		selects = strings.TrimSuffix(selects, ", ")
	}

	return selects
}

func (q *QueryGrammarBuilder[T]) BuildTables() string {
	tables := ""

	// If there is only one table, we're just doing a normal `select x from table_name`
	if len(q.builder.table) == 1 {
		tables += q.builder.table[0]
		return tables
	}

	tables += "["

	// If there is multiple, we're doing a `select x from [table1, table2]`
	for _, table := range q.builder.table {
		tables += table + ", "
	}

	tables = strings.TrimSuffix(tables, ", ")
	tables += "]"

	return tables
}

func (q *QueryGrammarBuilder[T]) BuildConditions() string {
	conditions := ""

	total := len(q.builder.conditions)
	for idx, condition := range q.builder.conditions {
		conditions += condition.key + " " + condition.exprOperator.String() + " " + condition.param.ForQuery()
		if idx != total-1 {
			conditions += " " + condition.queryOperator.String() + " "
		}
	}

	return conditions
}

func (q *QueryGrammarBuilder[T]) BuildOrderClauses() string {
	clauses := ""

	for _, clause := range q.builder.orderClauses {
		clauses += clause.field + " " + string(clause.direction) + ", "
	}

	clauses = strings.TrimSuffix(clauses, ", ")

	return clauses
}

func (q *QueryGrammarBuilder[T]) BuildFetch() string {
	fetches := ""

	for _, fetch := range q.builder.fetch {
		fetches += fetch + ", "
	}

	fetches = strings.TrimSuffix(fetches, ", ")

	return fetches
}

func (q *QueryGrammarBuilder[T]) Build() string {
	q.query = "SELECT "
	q.query += q.BuildSelects()
	q.query += " FROM "
	q.query += q.BuildTables()

	if q.builder.conditions != nil {
		q.query += " WHERE "
		q.query += q.BuildConditions()
	}

	if len(q.builder.orderClauses) > 0 {
		q.query += " ORDER BY "
		q.query += q.BuildOrderClauses()
	}

	if q.builder.limit != -1 {
		q.query += " LIMIT " + strconv.Itoa(q.builder.limit)
	}

	if q.builder.start != -1 {
		q.query += " START " + strconv.Itoa(q.builder.start)
	}

	if len(q.builder.fetch) > 0 {
		q.query += " FETCH "
		q.query += q.BuildFetch()
	}

	return q.query
}
