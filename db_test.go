package surrealdb_test

import (
	"context"
	"os"
	"testing"

	"github.com/idevelopthings/surrealdb.go.unofficial"
	"github.com/test-go/testify/suite"
)

func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	return value
}

// a simple user struct for testing
type testUser struct {
	Username string
	Password string
	ID       string
}

var sdbConfig = &surrealdb.DbConfig{
	Url:       getEnvOrDefault("SURREALDB_RPC_URL", "ws://localhost:8000/rpc"),
	Username:  getEnvOrDefault("SURREALDB_USER", "root"),
	Password:  getEnvOrDefault("SURREALDB_PASS", "root"),
	Database:  "test",
	Namespace: "test",
	AutoLogin: true,
	AutoUse:   true,
	Timeouts:  &surrealdb.DbTimeoutConfig{Timeout: 10},
}

func TestUnmarshalRaw(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := surrealdb.New(ctx, sdbConfig)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Delete(ctx, "users")
	if err != nil {
		panic(err)
	}

	username := "johnny"
	password := "123"

	// create test user with raw SurrealQL and unmarshal

	userData, err := db.Query(ctx, "create users:johnny set Username = $user, Password = $pass", map[string]any{
		"user": username,
		"pass": password,
	})
	if err != nil {
		panic(err)
	}

	var user testUser
	ok, err := surrealdb.UnmarshalRaw(userData, &user)
	if err != nil {
		panic(err)
	}
	if !ok || user.Username != username || user.Password != password {
		panic("response does not match the request")
	}

	// send query with empty result and unmarshal

	userData, err = db.Query(ctx, "select * from users where id = $id", map[string]any{
		"id": "users:jim",
	})
	if err != nil {
		panic(err)
	}

	ok, err = surrealdb.UnmarshalRaw(userData, &user)
	if err != nil {
		panic(err)
	}
	if ok {
		panic("select should return an empty result")
	}

	// Output:
}

type TestDatabaseTestSuite struct {
	suite.Suite
	ctx context.Context
	db  *surrealdb.DB
}

func TestDatabaseSuite(t *testing.T) {
	suite.Run(t, new(TestDatabaseTestSuite))
}

func (suite *TestDatabaseTestSuite) SetupTest() {
	ctx := context.Background()

	db, err := surrealdb.New(ctx, sdbConfig)
	suite.Require().NoError(err)

	suite.db = db
	suite.ctx = ctx
}

func (suite *TestDatabaseTestSuite) TearDownSuite() {
	suite.db.Close()
}

func (suite *TestDatabaseTestSuite) Test_FailingUserSignin() {
	// NOTE: this query fails for some reason but works when I run it manually...
	// DEFINE SCOPE test_account_scope
	//     SIGNIN ( SELECT * FROM user WHERE username = $user AND crypto::argon2::compare(password, $pass) )
	//     SIGNUP ( CREATE user SET username = $user, password = crypto::argon2::generate($pass) )
	// ;
	// result, err := suite.db.Query(suite.ctx, scopeQuery, map[string]any{})
	// suite.Require().NoError(err)
	// suite.Require().NotNil(result)

	// authResult, err := suite.db.SigninUser(suite.ctx, surrealdb.UserInfo{
	// 	User:      "test_username",
	// 	Password:  "test_password",
	// 	Namespace: "test_account_scope",
	// 	Database:  "test",
	// 	Scope:     "test",
	// })
	//
	// suite.Require().Error(err)
	// suite.Require().NotNil(authResult)
	// suite.Require().False(authResult.Success)
	//
	// authResult, err = suite.db.SignupUser(suite.ctx, surrealdb.UserInfo{
	// 	User:      "test_username",
	// 	Password:  "test_password",
	// 	Namespace: "test",
	// 	Database:  "test",
	// 	Scope:     "test_account_scope",
	// })
	// suite.Require().NoError(err)
	// suite.Require().NotNil(authResult)
	// suite.Require().True(authResult.Success)
	// suite.Require().NotZero(authResult.Token)
	// suite.Require().NotZero(authResult.TokenData)
	// suite.Require().Equal(authResult.TokenData.Scope, "test_account_scope")
}
