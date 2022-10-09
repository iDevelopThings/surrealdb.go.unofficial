package surrealdb_test

import (
	"os"
	"testing"
	"time"

	"github.com/idevelopthings/surrealdb.go.unofficial"
	Config "github.com/idevelopthings/surrealdb.go.unofficial/config"
)

func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	return value
}

var sdbConfig = &Config.DbConfig{
	Url:       getEnvOrDefault("SURREALDB_RPC_URL", "ws://localhost:8000/rpc"),
	Username:  getEnvOrDefault("SURREALDB_USER", "root"),
	Password:  getEnvOrDefault("SURREALDB_PASS", "root"),
	Database:  "test",
	Namespace: "test",
	AutoLogin: true,
	AutoUse:   true,
	Timeouts:  &Config.DbTimeoutConfig{Timeout: time.Duration(10) * time.Second},
}

func setupTests(t *testing.T) *surrealdb.DB {
	db, err := surrealdb.New(sdbConfig)
	if err != nil {
		if t != nil {
			t.Errorf("Error creating db: %s", err)
		}
		return nil
	}
	// insert testing data

	if _, err := db.Query("DELETE user:bob; UPDATE user:bob SET username = $username;", map[string]any{"username": "bob"}); err != nil {
		t.Errorf("Update user:bob errored: %d", err)
	}
	if _, err := db.Query("DELETE user:bob_two; UPDATE user:bob_two SET username = $username;", map[string]any{"username": "bob"}); err != nil {
		t.Errorf("Update user:bob_two errored: %d", err)
	}
	if _, err := db.Query("DELETE user:bob_three; UPDATE user:bob_three SET username = $username;", map[string]any{"username": "bob"}); err != nil {
		t.Errorf("Update user:bob_three errored: %d", err)
	}
	if _, err := db.Query("DELETE user:bob_four; UPDATE user:bob_three SET username = $username;", map[string]any{"username": "bob"}); err != nil {
		t.Errorf("Update user:bob_three errored: %d", err)
	}

	return db
}

type testUserInformation struct {
	Username string `json:"username,omitempty"`
	NewValue string `json:"newValue,omitempty"`
	Nickname string `json:"nickname,omitempty"`
	Age      int    `json:"age,omitempty"`
}

func Test_QueryResolver_Query(t *testing.T) {
	_ = setupTests(t)

	result := surrealdb.Query[testUserInformation]("SELECT * FROM user:bob WHERE username = $username;", map[string]any{
		"username": "bob",
	})

	if result.HasError() {
		t.Errorf("Query errored: %d", result.Error())
		return
	}

	bob := result.First()

	if bob == nil {
		t.Errorf("Expected object for bob, got nil")
		return
	}
	if bob.Username != "bob" {
		t.Errorf("Expected bob, got %s", bob.Username)
		return
	}

}

func Test_QueryResolver_Select(t *testing.T) {
	_ = setupTests(t)

	result := surrealdb.Select[testUserInformation]("user")
	if result.HasError() {
		t.Errorf("Query errored: %d", result.Error())
		return
	}
	allBobs := result.All()
	if len(allBobs) == 0 {
		t.Errorf("Expected array of bobs, got empty array")
		return
	}

	bob := result.First()
	if bob == nil {
		t.Errorf("Expected object for bob, got nil")
		return
	}
	if bob.Username != "bob" {
		t.Errorf("Expected bob, got %s", bob.Username)
		return
	}
}

func Test_QueryResolver_Create(t *testing.T) {
	_ = setupTests(t)

	result := surrealdb.Create[testUserInformation]("user", testUserInformation{Username: "bob"})
	if result.HasError() {
		t.Errorf("Query errored: %d", result.Error())
		return
	}
	bob := result.Item()
	if bob == nil {
		t.Errorf("Expected object for bob, got nil")
		return
	}
	if bob.Username != "bob" {
		t.Errorf("Expected bob, got %s", bob.Username)
		return
	}

}

func Test_QueryResolver_Update(t *testing.T) {
	_ = setupTests(t)

	cResult := surrealdb.Create[testUserInformation]("user:bob_four", testUserInformation{
		Username: "bob_four",
		NewValue: "empty",
	})
	if cResult.HasError() {
		t.Errorf("Query errored: %d", cResult.Error())
		return
	}
	if cResult.Item().Username != "bob_four" {
		t.Errorf("Expected bob_four, got %s", cResult.Item().Username)
		return
	}
	if cResult.Item().NewValue != "empty" {
		t.Errorf("Expected empty, got %s", cResult.Item().NewValue)
		return
	}

	result := surrealdb.Change[testUserInformation]("user:bob_four", testUserInformation{NewValue: "hello world"})
	if result.HasError() {
		t.Errorf("Query errored: %d", result.Error())
		return
	}
	bob := result.First()
	if bob == nil {
		t.Errorf("Expected object for bob, got nil")
		return
	}
	if bob.Username != "bob_four" {
		t.Errorf("Expected bob, got %s", bob.Username)
		return
	}
	if bob.NewValue != "hello world" {
		t.Errorf("Expected 'hello world', got %s", bob.NewValue)
		return
	}

}

func Test_QueryResolver_Change(t *testing.T) {
	_ = setupTests(t)

	result := surrealdb.Change[testUserInformation]("user:bob_two", testUserInformation{NewValue: "changed value"})
	if result.HasError() {
		t.Errorf("Query errored: %d", result.Error())
		return
	}
	bob := result.First()
	if bob == nil {
		t.Errorf("Expected object for bob, got nil")
		return
	}
	if bob.Username != "bob" {
		t.Errorf("Expected bob, got %s", bob.Username)
		return
	}
	if bob.NewValue != "changed value" {
		t.Errorf("Expected changed value, got %s", bob.NewValue)
		return
	}

}

func Test_QueryResolver_Modify(t *testing.T) {
	_ = setupTests(t)

	patches := []surrealdb.Patch{
		{Op: "add", Path: "nickname", Value: "Bobs nickname"},
		{Op: "add", Path: "age", Value: 44},
	}

	result := surrealdb.Modify("user:bob_three", patches)

	if result.HasError() {
		t.Errorf("Query errored: %d", result.Error())
		return
	}
	ops := result.First()
	if ops == nil || len(ops) == 0 {
		t.Errorf("Expected array of ops, got nil or empty")
		return
	}
	if len(ops) != 2 {
		t.Errorf("Expected 2 ops, got %d", len(ops))
		return
	}

	if ops[0].Op != "add" {
		t.Errorf("Expected add, got %s", ops[0].Op)
		return
	}
}

func Test_QueryResolver_Delete(t *testing.T) {
	_ = setupTests(t)

	result := surrealdb.Delete[testUserInformation]("user:bob_three")

	if result.HasError() {
		t.Errorf("Query errored: %d", result.Error())
		return
	}

}
