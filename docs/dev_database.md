# Add new database type

Write a new package to implement the interface [database.Connection:](https://github.com/FreifunkBremen/yanic/blob/master/database/database.go)

```go
type Connection interface {
	InsertNode(node *runtime.Node)

	InsertLink(*runtime.Link, time.Time)

	InsertGlobals(*runtime.GlobalStats, time.Time, string)

	PruneNodes(deleteAfter time.Duration)

	Close()
}
```

**InsertNode** is stores statistics per node

**InsertLink** is stores statistics per link

**InsertGlobals** is stores global statistics (by `site_code`, and "global" like in `runtime.GLOBAL_SITE` overall sites).

**PruneNodes** is prunes historical per-node data

**Close** is called during shutdown of Yanic.



For startup, you need to bind your database type by calling `database.RegisterAdapter("typeofdatabase",ConnectFunction)`

it should be in the `func init() {}` of your package.



The _typeofdatabase_ is used as mapping in the configuration `[[database.connection.typeofdatabase]]` the `map[string]interface{}` of the content are parsed to the _ConnectFunction_ and on of your implemented `Connection` or a `error` is needed as result.



Short: the function signature of _ConnectFunction_ should be `func Connect(configuration interface{}) (Connection, error)`



At last add you import string to compile the your database as well in this [all](https://github.com/FreifunkBremen/yanic/blob/master/database/all/main.go) package.



TIP: take a look in the easy database type [logging](https://github.com/FreifunkBremen/yanic/blob/master/database/logging/file.go).
