# surrealdb.go

Unofficial fork of the [surrealdb](https://github.com/surrealdb/surrealdb.go) package

# Credit:

A lot of the underlying package was written by people in the community on
the [original package](https://github.com/surrealdb/surrealdb.go)

I just wanted to have my modified version easily accessible to my self, and maybe others. So I take no credits for that(
hopefully no one sees me as ripping there work and understands)

### Features:

- Contains my system for a "Query Resolver" using generics
- Has options for auto context setup/timeouts(it's a pain in the ass to pass these around in go)
- A little easier to set up for server usage only
- Has a global DB instance

### Installation:
```shell
go get github.com/idevelopthings/surrealdb.go.unofficial
```

### Examples:

#### Setup:

```go
db, err := surrealdb.New(ctx, &surrealdb.DbConfig{
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
   Timeouts:  &surrealdb.DbTimeoutConfig{Timeout: 10},
})
```

#### Query Resolver/Generics

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

