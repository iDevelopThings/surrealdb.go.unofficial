package surrealdb

import (
	"fmt"
	"time"
)

type QueryBuilderSelectField struct {
	key string
	as  *string
}

type OrderDirection string

const (
	OrderDirectionAsc  OrderDirection = "ASC"
	OrderDirectionDesc OrderDirection = "DESC"
)

type QueryBuilderOrderClause struct {
	field     string
	direction OrderDirection
}

type QueryBuilderParam struct {
	field     string
	paramName string
	value     any
}

func (p *QueryBuilderParam) ForQuery() string {
	return "$" + p.paramName
}

type ConditionType string

const (
	RawCondition   ConditionType = "raw"
	BasicCondition ConditionType = "basic"
)

type QueryBuilderBasicCondition struct {
	conditionType ConditionType
	key           string
	param         *QueryBuilderParam
	query         string
	exprOperator  Operator
	queryOperator Operator
}

type QueryBuilder[T any] struct {
	table        []string
	selections   []*QueryBuilderSelectField
	conditions   []*QueryBuilderBasicCondition
	orderClauses []*QueryBuilderOrderClause
	params       map[string]*QueryBuilderParam
	fetch        []string
	limit        int
	start        int

	grammarBuilder *QueryGrammarBuilder[T]

	resolver *ResolvedQuery[T]
}

func NewBuilder[T any](table ...string) *QueryBuilder[T] {
	builder := &QueryBuilder[T]{
		params: make(map[string]*QueryBuilderParam),
		limit:  -1,
		start:  -1,
	}
	if len(table) > 0 {
		builder.table = table
	}
	return builder
}

// From sets the table to query from(when only using one table)
func (qb *QueryBuilder[T]) From(table string) *QueryBuilder[T] {
	qb.table = []string{table}
	return qb
}

// FromMultiple sets the table to query from(when using multiple tables)
func (qb *QueryBuilder[T]) FromMultiple(tables ...string) *QueryBuilder[T] {
	qb.table = tables
	return qb
}

// Select adds a field to the selection
func (qb *QueryBuilder[T]) Select(field string, as ...string) *QueryBuilder[T] {
	selectField := &QueryBuilderSelectField{key: field}
	if len(as) > 0 {
		selectField.as = &as[0]
	}
	qb.selections = append(qb.selections, selectField)
	return qb
}

// SelectMany adds a field to the selection
// This works like: .SelectMany([][]string{"FIELD NAME", "SELECT AS"}, [][]string{"USERNAME", "name"})
func (qb *QueryBuilder[T]) SelectMany(fields ...[]string) *QueryBuilder[T] {
	for _, _field := range fields {
		field := _field
		selectField := &QueryBuilderSelectField{key: field[0]}
		if len(field) > 1 {
			selectField.as = &field[1]
		}
		qb.selections = append(qb.selections, selectField)
	}
	return qb
}

// OrderBy adds an order by clause
func (qb *QueryBuilder[T]) OrderBy(field string, direction ...OrderDirection) *QueryBuilder[T] {
	if len(direction) == 0 {
		direction = []OrderDirection{OrderDirectionAsc}
	}
	if direction[0] != OrderDirectionDesc && direction[0] != OrderDirectionAsc {
		panic("invalid order direction")
	}
	qb.orderClauses = append(qb.orderClauses, &QueryBuilderOrderClause{
		field:     field,
		direction: direction[0],
	})
	return qb
}

func (qb *QueryBuilder[T]) addParam(field string, paramName string, value any) *QueryBuilderParam {
	qb.params[paramName] = &QueryBuilderParam{
		paramName: paramName,
		value:     value,
		field:     field,
	}
	return qb.params[paramName]
}

// Where adds a basic where x = y clause
func (qb *QueryBuilder[T]) Where(field string, value any) *QueryBuilder[T] {
	paramName := fmt.Sprintf("whereVar_%s_%v", field, len(qb.params))

	qb.conditions = append(qb.conditions, &QueryBuilderBasicCondition{
		conditionType: BasicCondition,
		key:           field,
		exprOperator:  Operators.Equal,
		queryOperator: Operators.And,
		param:         qb.addParam(field, paramName, value),
	})

	return qb
}

func (qb *QueryBuilder[T]) Limit(value int) *QueryBuilder[T] {
	qb.limit = value
	return qb
}

func (qb *QueryBuilder[T]) Start(value int) *QueryBuilder[T] {
	qb.start = value
	return qb
}

func (qb *QueryBuilder[T]) Fetch(fields ...string) *QueryBuilder[T] {
	for _, f := range fields {
		field := f
		qb.fetch = append(qb.fetch, field)
	}
	return qb
}

// GetQuery returns the query string
func (qb *QueryBuilder[T]) GetQuery() string {
	if qb.grammarBuilder == nil {
		qb.grammarBuilder = BuildQueryGrammar[T](qb)
	}
	return qb.grammarBuilder.Build()
}

func (qb *QueryBuilder[T]) GetParams() map[string]any {
	params := make(map[string]any)
	for _, _param := range qb.params {
		param := _param
		params[param.paramName] = param.value
	}
	return params
}

func (qb *QueryBuilder[T]) Execute() *ResolvedQuery[T] {
	resolved := Query[T](qb.GetQuery(), qb.GetParams())

	qb.resolver = resolved

	return qb.resolver
}

func (qb *QueryBuilder[T]) First() *T {
	qb.Limit(1)

	return qb.Execute().First()
}

func (qb *QueryBuilder[T]) Get() []T {
	return qb.Execute().All()
}

func (qb *QueryBuilder[T]) HasError() bool {
	return qb.resolver.HasError()
}
func (qb *QueryBuilder[T]) Error() error {
	return qb.resolver.Error()
}
func (qb *QueryBuilder[T]) TotalTimeTaken() time.Duration {
	return qb.resolver.TotalTimeTaken()
}
func (qb *QueryBuilder[T]) FirstQueryResult() *ResultQuery[T] {
	return qb.resolver.FirstQueryResult()
}
func (qb *QueryBuilder[T]) Results() []ResultQuery[T] {
	return qb.resolver.Results()
}
func (qb *QueryBuilder[T]) IsEmpty() bool {
	return qb.resolver.IsEmpty()
}
