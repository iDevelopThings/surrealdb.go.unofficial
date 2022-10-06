package surrealdb

import (
	"context"
	"time"
)

type QueryResolver[T any] struct {
	err    error
	query  string
	params any
	ctx    context.Context
}

func createResolver[T any](ctx context.Context, query string, params any) *QueryResolver[T] {
	resolver := &QueryResolver[T]{
		query:  query,
		params: params,
		ctx:    ctx,
	}

	return resolver
}

type QueryConfig struct {
	Ctx    context.Context
	Db     *DB
	Query  string
	Params map[string]any
}

// Query creates a new query resolver
// Automatically uses the global db instance, ctx and uses ctx timeouts if configured
func Query[T any](query string, params ...map[string]any) *ResolvedQuery[T] {
	if len(params) == 0 {
		// Ensure there's always a default, surreal doesn't like it missing
		params = append(params, map[string]any{})
	}

	var config = QueryConfig{
		Ctx:    dbConfig.Ctx,
		Db:     db,
		Params: params[0],
		Query:  query,
	}

	if dbConfig.Timeouts != nil && dbConfig.Timeouts.Timeout > 0 {
		ctx, cancel := context.WithTimeout(dbConfig.Ctx, dbConfig.Timeouts.Timeout*time.Millisecond)
		defer cancel()
		config.Ctx = ctx
	}

	return QueryWithConfig[T](config)
}

// QueryWithConfig creates a new query resolver
// Uses a specific db instance and ctx, does not use auto ctx timeouts
func QueryWithConfig[T any](config QueryConfig) *ResolvedQuery[T] {
	return createResolver[T](config.Ctx, config.Query, config.Params).runQuery(config.Db)
}

func (resolver *QueryResolver[T]) runQuery(db *DB) *ResolvedQuery[T] {
	result, err := db.send(resolver.ctx, "query", resolver.query, resolver.params)
	if err != nil {
		panic(err)
	}
	if _, ok := result.(*RPCRawResponse); !ok {
		panic("Invalid response")
	}
	return NewResolvedQuery[T](result.(*RPCRawResponse))
}

func Create[T any](what string, params ...map[string]any) ResolvedCreateResult[T] {
	if len(params) == 0 {
		// Ensure there's always a default, surreal doesn't like it missing
		params = append(params, map[string]any{})
	}
	var config = QueryConfig{
		Ctx:    dbConfig.Ctx,
		Db:     db,
		Params: params[0],
		Query:  what,
	}
	if dbConfig.Timeouts != nil && dbConfig.Timeouts.Timeout > 0 {
		ctx, cancel := context.WithTimeout(dbConfig.Ctx, dbConfig.Timeouts.Timeout*time.Millisecond)
		defer cancel()
		config.Ctx = ctx
	}
	return createResolver[T](config.Ctx, config.Query, config.Params).runCrud(config.Db, "create")
}

func Update[T any](what string, params ...map[string]any) ResolvedUpdateResult[T] {
	if len(params) == 0 {
		// Ensure there's always a default, surreal doesn't like it missing
		params = append(params, map[string]any{})
	}
	var config = QueryConfig{
		Ctx:    dbConfig.Ctx,
		Db:     db,
		Params: params[0],
		Query:  what,
	}
	if dbConfig.Timeouts != nil && dbConfig.Timeouts.Timeout > 0 {
		ctx, cancel := context.WithTimeout(dbConfig.Ctx, dbConfig.Timeouts.Timeout*time.Millisecond)
		defer cancel()
		config.Ctx = ctx
	}
	return createResolver[T](config.Ctx, config.Query, config.Params).runCrud(config.Db, "update")
}

func Change[T any](what string, params ...map[string]any) ResolvedUpdateResult[T] {
	if len(params) == 0 {
		// Ensure there's always a default, surreal doesn't like it missing
		params = append(params, map[string]any{})
	}
	var config = QueryConfig{
		Ctx:    dbConfig.Ctx,
		Db:     db,
		Params: params[0],
		Query:  what,
	}
	if dbConfig.Timeouts != nil && dbConfig.Timeouts.Timeout > 0 {
		ctx, cancel := context.WithTimeout(dbConfig.Ctx, dbConfig.Timeouts.Timeout*time.Millisecond)
		defer cancel()
		config.Ctx = ctx
	}
	return createResolver[T](config.Ctx, config.Query, config.Params).runCrud(config.Db, "change")
}

func Modify(what string, data []Patch) *ResolvedModifyResult {

	var config = QueryConfig{
		Ctx:   dbConfig.Ctx,
		Db:    db,
		Query: what,
	}
	if dbConfig.Timeouts != nil && dbConfig.Timeouts.Timeout > 0 {
		ctx, cancel := context.WithTimeout(dbConfig.Ctx, dbConfig.Timeouts.Timeout*time.Millisecond)
		defer cancel()
		config.Ctx = ctx
	}
	return createResolver[any](config.Ctx, config.Query, data).runModify(config.Db)
}

func Delete[T any](what string, params ...map[string]any) ResolvedUpdateResult[T] {
	if len(params) == 0 {
		// Ensure there's always a default, surreal doesn't like it missing
		params = append(params, map[string]any{})
	}
	var config = QueryConfig{
		Ctx:    dbConfig.Ctx,
		Db:     db,
		Params: params[0],
		Query:  what,
	}
	if dbConfig.Timeouts != nil && dbConfig.Timeouts.Timeout > 0 {
		ctx, cancel := context.WithTimeout(dbConfig.Ctx, dbConfig.Timeouts.Timeout*time.Millisecond)
		defer cancel()
		config.Ctx = ctx
	}
	return createResolver[T](config.Ctx, config.Query, config.Params).runCrud(config.Db, "delete")
}

func (resolver *QueryResolver[T]) runCrud(db *DB, method string) *ResolvedCrudResult[T] {
	result, err := db.send(resolver.ctx, method, resolver.query, resolver.params)
	if err != nil {
		panic(err)
	}
	if _, ok := result.(*RPCRawResponse); !ok {
		panic("Invalid response")
	}
	return NewResolvedCrudResult[T](result.(*RPCRawResponse))
}

func (resolver *QueryResolver[T]) runModify(db *DB) *ResolvedModifyResult {
	result, err := db.send(resolver.ctx, "modify", resolver.query, resolver.params)
	if err != nil {
		panic(err)
	}
	if _, ok := result.(*RPCRawResponse); !ok {
		panic("Invalid response")
	}
	return NewResolvedModifyResult(result.(*RPCRawResponse))
}
