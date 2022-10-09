package surrealdb

type QueryResolver[T any] struct {
	err    error
	query  string
	params any
}

func createResolver[T any](query string, params any) *QueryResolver[T] {
	resolver := &QueryResolver[T]{
		query:  query,
		params: params,
	}

	return resolver
}

type QueryConfig struct {
	Db     *DB
	Query  string
	Params any
}

// Query creates a new query resolver
// Automatically uses the global db instance, ctx and uses ctx timeouts if configured
func Query[T any](query string, params ...map[string]any) *ResolvedQuery[T] {
	if len(params) == 0 {
		// Ensure there's always a default, surreal doesn't like it missing
		params = append(params, map[string]any{})
	}

	var config = QueryConfig{
		Db:     Connection,
		Params: params[0],
		Query:  query,
	}

	// if dbConfig.Timeouts != nil && dbConfig.Timeouts.Timeout > 0 {
	// 	ctx, cancel := context.WithTimeout(dbConfig.Ctx, dbConfig.Timeouts.Timeout*time.Millisecond)
	// 	defer cancel()
	// 	config.Ctx = ctx
	// }

	return QueryWithConfig[T](config)
}

// QueryWithConfig creates a new query resolver
// Uses a specific db instance and ctx, does not use auto ctx timeouts
func QueryWithConfig[T any](config QueryConfig) *ResolvedQuery[T] {
	return createResolver[T](config.Query, config.Params).runQuery(config.Db)
}

func (resolver *QueryResolver[T]) runQuery(db *DB) *ResolvedQuery[T] {
	result, err := db.send("query", resolver.query, resolver.params)
	if err != nil {
		panic(err)
	}

	return NewResolvedQuery[T](result)
}

// Select this will select one or many documents
// It is the same as: https://surrealdb.com/docs/integration/http#select-all
func Select[T any](what string) *ResolvedCrudResult[T] {
	return createResolver[T](what, map[string]any{}).runCrud(Connection, "select")
}

// Create This will create a new document
// It is the same as: https://surrealdb.com/docs/integration/http#create-all
func Create[T any, DType any | map[string]any](what string, data DType) ResolvedCreateResult[T] {
	return createResolver[T](what, data).runCrud(Connection, "create")
}

// Update This will apply a "replace" change to the document
// It is the same as: https://surrealdb.com/docs/integration/http#update-one
func Update[T any, DType any | map[string]any](what string, data DType) ResolvedUpdateResult[T] {
	return createResolver[T](what, data).runCrud(Connection, "update")
}

// Change This will apply a "merge" change to the document
// It is the same as: https://surrealdb.com/docs/integration/http#modify-one
func Change[T any, DType any | map[string]any](what string, data DType) ResolvedUpdateResult[T] {
	return createResolver[T](what, data).runCrud(Connection, "change")
}

// Modify applies a JSONPatch to the document
func Modify(what string, data []Patch) *ResolvedModifyResult {
	return createResolver[any](what, data).runModify(Connection)
}

// Delete deletes a document or all documents
func Delete[T any](what string) ResolvedUpdateResult[T] {
	return createResolver[T](what, map[string]any{}).runCrud(Connection, "delete")
}

func (resolver *QueryResolver[T]) runCrud(db *DB, method string) *ResolvedCrudResult[T] {
	result, err := db.send(method, resolver.query, resolver.params)
	if err != nil {
		panic(err)
	}
	return NewResolvedCrudResult[T](result)
}

func (resolver *QueryResolver[T]) runModify(db *DB) *ResolvedModifyResult {
	result, err := db.send("modify", resolver.query, resolver.params)
	if err != nil {
		panic(err)
	}
	return NewResolvedModifyResult(result)
}
