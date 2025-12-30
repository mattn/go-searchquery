module github.com/mattn/go-searchquery/example/meilisearch

go 1.25.0

require (
	github.com/mattn/go-searchquery v0.0.0-00010101000000-000000000000
	github.com/meilisearch/meilisearch-go v0.35.0
)

require (
	github.com/andybalholm/brotli v1.1.1 // indirect
	github.com/golang-jwt/jwt/v4 v4.5.2 // indirect
)

replace github.com/mattn/go-searchquery => ../..
