postgres: 
	docker run --name github-api-db-con -p 5439:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:14-alpine

createdb:
	docker exec -it github-api-db-con createdb --username=root --owner=root github_api_db

dropdb:
	docker exec -it github-api-db-con dropdb github_api_db

test:
	go test -v -cover ./...

server: 
	go run cmd/main.go

.PHONY: postgres createdb dropdb test server