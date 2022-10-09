package surrealdb

type Model[T any] struct {
}

type SurrealModel interface {
	TableName() string
}
