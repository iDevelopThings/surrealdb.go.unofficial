package surrealdb

import (
	"errors"
	"log"

	Config "github.com/idevelopthings/surrealdb.go.unofficial/config"
	"github.com/idevelopthings/surrealdb.go.unofficial/internal"
)

var (
	ErrInvalidLoginResponse = errors.New("invalid login response")
)

// DB is a client for the SurrealDB database that holds are websocket connection.
type DB struct {
	ws *internal.WS
}

var dbConfig *Config.DbConfig
var Connection *DB

// New Creates a new DB instance given a WebSocket URL.
func New(config *Config.DbConfig) (*DB, error) {
	saveGlobalConf := true
	if Connection != nil {
		if dbConfig.Url == config.Url &&
			dbConfig.Username == config.Username &&
			dbConfig.Password == config.Password &&
			dbConfig.Database == config.Database &&
			dbConfig.Namespace == config.Namespace {
			log.Println("SurrealDB: DB already initialized, configurations match the current configuration, reusing the existing connection")
			return Connection, nil
		} else {
			log.Println("SurrealDB: DB already initialized, this new initialization will not override the existing, so make sure you save a reference.")
			saveGlobalConf = false
		}
	}

	conf := &Config.DbConfig{
		Url:       config.Url,
		Username:  config.Username,
		Password:  config.Password,
		Database:  config.Database,
		Namespace: config.Namespace,
		AutoLogin: config.AutoLogin,
		AutoUse:   config.AutoUse,
		Timeouts:  config.Timeouts,
	}
	if saveGlobalConf {
		dbConfig = conf
	}

	ws, err := internal.NewWebsocket(conf.Url, config.Timeouts)
	if err != nil {
		return nil, err
	}

	inst := &DB{ws}

	if conf.AutoLogin {
		_, err = inst.Signin(UserInfo{User: conf.Username, Password: conf.Password})
		if err != nil {
			return nil, err
		}
	}

	if conf.AutoUse {
		_, err = inst.Use(conf.Namespace, conf.Database)
		if err != nil {
			return nil, err
		}
	}

	if saveGlobalConf {
		Connection = inst
	}

	return inst, nil
}

// --------------------------------------------------
// Public methods
// --------------------------------------------------

// Close closes the underlying WebSocket connection.
func (db *DB) Close() error {
	return db.ws.Close()
}

// --------------------------------------------------

// Use is a method to select the namespace and table to use.
func (db *DB) Use(ns string, dbname string) (any, error) {
	return db.send("use", ns, dbname)
}

func (db *DB) Info() (any, error) {
	return db.send("info")
}

// Signup is a helper method for signing up a new user.
func (db *DB) Signup(vars any) (any, error) {
	return db.send("signup", vars)
}

// SignupUser is a helper method for signing in a user and returning a typed response
func (db *DB) SignupUser(vars UserInfo) (*AuthenticationResult, error) {
	authResult := &AuthenticationResult{Success: false}
	result, err := db.send("signup", vars)
	if err != nil {
		return authResult, err
	}

	err = authResult.fromQuery(result)

	return authResult, err
}

// Signin is a helper method for signing in a user.
func (db *DB) Signin(vars UserInfo) (any, error) {
	return db.send("signin", vars)
}

// SigninUser is a helper method for signing in a user and returning a typed response
// Note: This will probably fail when signing in as a root user, but for
// a regular user(via a scope for example) we get a JWT response
func (db *DB) SigninUser(vars UserInfo) (*AuthenticationResult, error) {
	authResult := &AuthenticationResult{Success: false}
	result, err := db.send("signin", vars)
	if err != nil {
		return authResult, err
	}
	if err != nil {
		return authResult, err
	}

	err = authResult.fromQuery(result)

	return authResult, err
}

func (db *DB) Invalidate() (any, error) {
	return db.send("invalidate")
}

func (db *DB) Authenticate(token string) (any, error) {
	return db.send("authenticate", token)
}

// --------------------------------------------------

func (db *DB) Live(table string) (any, error) {
	return db.send("live", table)
}

func (db *DB) Kill(query string) (any, error) {
	return db.send("kill", query)
}

func (db *DB) Let(key string, val any) (any, error) {
	return db.send("let", key, val)
}

// Query is a convenient method for sending a query to the database.
func (db *DB) Query(sql string, vars any) (any, error) {
	return db.send("query", sql, vars)
}

// --------------------------------------------------
// Private methods
// --------------------------------------------------

// send is a helper method for sending a query to the database.
func (db *DB) send(method string, params ...any) (*internal.RPCRawResponse, error) {
	id := xid()

	// response, err := db.ws.Send(id, method, params)

	chn := db.ws.Once(id, method)
	// here we send the args through our websocket connection
	db.ws.Send(id, method, params)

	ctx, cancel := db.ws.NewContext()
	defer cancel()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()

	case r := <-chn:
		if r.Err != nil {
			return nil, r.Err
		}

		return r.Value, nil
	}
}
