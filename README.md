# surrealdb.go

Unofficial fork of the [surrealdb](https://github.com/surrealdb/surrealdb.go) package

# Credit:

A lot of the underlying package was written by people in the community on
the [original package](https://github.com/surrealdb/surrealdb.go)

I just wanted to have my modified version easily accessible to my self, and maybe others. So I take no credits for that(
hopefully no one sees me as ripping there work and understands)

# Features:

- Contains my system for a "Query Resolver" using generics
- Has options for auto context setup/timeouts(it's a pain in the ass to pass these around in go)
- A little easier to set up for server usage only
- Has a global DB instance

# Installation:
```shell
go get github.com/idevelopthings/surrealdb.go.unofficial
```

# Examples:

## Setup:

```go
import (
	"github.com/idevelopthings/surrealdb.go.unofficial"
    "github.com/idevelopthings/surrealdb.go.unofficial/config"
)
db, err := surrealdb.New(ctx, &&Config.DbConfig{
   Url:       "ws://localhost:8000/rpc",
   Username:  "root",
   Password:  "root",
   Database:  "test",
   Namespace: "test",
   // This will call db.signin() with the supplied credentials
   AutoLogin: true,
   // This will call db.use() with the supplied credentials
   AutoUse:   true,
   // When using surrealdb.Query[MyModel](), timeouts will be configured
   Timeouts:  &Config.DbTimeoutConfig{Timeout: time.Duration(10) * time.Second},
})
```

# Query Resolver/Generics

## Quick Overview:
A lot of the query resolver use a similar interface/resolver setup

These methods exist on pretty much all resolver responses:
```go
result.HasError()  // Check if there was an error
result.Error()     // Get the go error if there is one
result.First()     // If you expect only 1 user object for example, it will take the first out of the response and return it(so you don't have to play with arrays, it will be a "User" instance for example)
result.All()       // Get all the results(an array of Users)
result.Results()   // Get the raw surreal response/results
result.IsEmpty()   // Check if there is any items in the result
```

**In some cases, for example, create, since only one record is ever expected(as we can only create one)**

It will use a ``.Item()`` method instead of ``.First()``, and will not contain ``.All()``


## Query
Run a query with parameters(parameters should be used to prevent injection)
```go
result := surrealdb.Query[User]("select * from users where name = $name", map[string]any{
   "name": "bob",
})

result.HasError()           // Check if there was an error
result.Error()              // Get the go error if there is one
result.AllAreSuccessful()   // Check if all queries were successful
result.TotalTimeTaken()     // Calculate the total time of all queries
result.FirstQueryResult()   // Get the first query result(surreal returns multiple results for a query, so if you run more than 1, this is useful)
result.First()              // If you expect only 1 user object for example, it will take the first out of the response and return it(so you don't have to play with arrays, it will be a "User" instance for example)
result.All()                // Get all the results(an array of Users)
result.Results()            // Get the raw surreal response/results
```


## Select
Select one or many records

If an id is supplied for the first arg, it will select the one record with that id
Otherwise, it expects a table name, where it will then select all records from that table.

```go
// Select one:
result := surrealdb.Select[User]("user:12345")
// Select all:
result := surrealdb.Select[User]("user")

// Refer to the above overview for the methods available on the result
```

## Create
Create one record

```go
// We can pass a map[string]any as the value
surrealdb.Create[User]("user", map[string]any {"username": "Bob"})
// Or a struct:
surrealdb.Create[User]("user", User{Username: "Bob"})
// Refer to the above overview for the methods available on the result
```

## Update
Update one or many records

If an id is supplied for the first arg, it will update the one record with that id
Otherwise, it expects a table name, where it will then update all records from that table.

This method will overwrite the entire record, if you want to update only a few fields, use the ``.Change()`` method

```go
// We can pass a map[string]any as the value
surrealdb.Update[User]("user:12345", map[string]any {"username": "Bob"})
// Or a struct:
surrealdb.Update[User]("user:12345", User{Username: "Bob"})
// Refer to the above overview for the methods available on the result
```

## Change
Update one or many records, but will use a MERGE operation, so it will only update the fields you specify

If an id is supplied for the first arg, it will update the one record with that id
Otherwise, it expects a table name, where it will then update all records from that table.

```go
// We can pass a map[string]any as the value
surrealdb.Change[User]("user:12345", map[string]any {"username": "Bob"})
// Or a struct:
surrealdb.Change[User]("user:12345", User{Username: "Bob"})
// Refer to the above overview for the methods available on the result
```

## Delete
Delete one record, or all records in a table

If an id is supplied for the first arg, it will delete the one record with that id
Otherwise, it expects a table name, where it will then delete all records from that table.

```go
// We can pass a map[string]any as the value
surrealdb.Delete[User]("user:12345")
// Or a struct:
surrealdb.Delete[User]("user:12345")
// Refer to the above overview for the methods available on the result
```



# WIP Query Builder

Example:

```go 
query := surrealdb.NewBuilder[User]("user").
    Where("username", "bob")

queryResult := query.First() // = *User{Username: "bob"}
```

`query` contains some "proxy" methods that the QueryResolvers above also have

So if you keep a reference to it, it can be used to check if there was an error and such. 

For example, with the above

```go
query := surrealdb.NewBuilder[User]("user").
    Where("username", "bob")
user := query.First()

if query.HasError() {
    log.Fatal(query.Error())
}
if query.IsEmpty() {
    log.Fatal("No user found")
}
// Do something with user
```