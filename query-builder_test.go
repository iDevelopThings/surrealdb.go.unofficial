package surrealdb_test

import (
	"testing"

	"github.com/idevelopthings/surrealdb.go.unofficial"
)

func TestQueryBuilder_Basic(t *testing.T) {
	builder := surrealdb.NewBuilder[any]("user").
		Select("id").
		Select("name").
		Select("something", "something_else").
		Where("name", "Bob").
		OrderBy("id", surrealdb.OrderDirectionDesc).
		Limit(1)

	query := builder.GetQuery()

	if query != "SELECT id, name, something AS something_else FROM user WHERE name = $whereVar_name_0 ORDER BY id DESC LIMIT 1" {
		t.Error("query is not correct")
	}
}

func TestQueryBuilder_Resolving(t *testing.T) {
	_ = setupTests(t)

	query := surrealdb.NewBuilder[testUserInformation]("user").
		Where("username", "bob")

	queryResult := query.First()

	if query.HasError() {
		t.Errorf("Query errored: %d", query.Error())
		return
	}

	if queryResult == nil {
		t.Errorf("Expected object for bob, got nil")
		return
	}

	if queryResult.Username != "bob" {
		t.Errorf("Expected bob, got %s", queryResult.Username)
		return
	}
}
