module github.com/mattn/go-searchquery/example/mysql

go 1.25.0

require (
	github.com/go-sql-driver/mysql v1.9.3
	github.com/mattn/go-searchquery v0.0.0-20251228165728-790fffe5df31
)

require filippo.io/edwards25519 v1.1.0 // indirect

replace github.com/mattn/go-searchquery => ../..
