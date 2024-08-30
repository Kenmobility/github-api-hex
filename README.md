# GitHub API Service## Description

This service fetches data from GitHub's public APIs to retrieve repository commits, saves the data in a persistent store, and continuously monitors the repository for changes at a set interval.

## Requirements- Golang 1.20+
- PostgreSQL

## 1. Clone the repository, cd into the project folder and download required go dependencies
```bash
git clone https://github.com/kenmobility/github-api-service.git
```
```bash
cd github-api-hex
```
```bash
go mod tidy
```
## 2. Run below command to Duplicate .env.example file and rename to .env
```bash
cp .env.example .env
```

## 3. Set environmental variables:
- the duplicated .env.example file already has default variables that the program needs to run except for GIT_HUB_TOKEN env variable (feel free to change the values).
- go to [https://github.com/](GitHub) to set up a GitHub API token (i.e Personal access token) and set value for the GIT_HUB_TOKEN env variable.

- (OPTIONAL) if the default values of the DATABASE_HOST, DATABASE_PORT, DATABASE_USER, DATABASE_PASSWORD or DATABASE_NAME in .env file were altered, ensure that the 'make postgres' command in the makefile matches the new set values.
  ```bash
  postgres: 
	  docker run --name github-api-hex-db-con -p {DATABASE_PORT}:5432 -e POSTGRES_USER={DATABASE_USER} -e POSTGRES_PASSWORD={DATABASE_PASSWORD} -d postgres:14-alpine
  ```

## 4 Open Docker desktop application
- Ensure that docker desktop is started and running on your machine 

## 5. Attempt Removal of postgres container by name
- run 'docker rm github-api-hex-db-con' to ensure that the container name for postgreSQL does not exist already
```bash
docker rm github-api-hex-db-con
``` 

## 6. Run makefile commands 
- run 'make postgres' to pull and run PostgreSQL instance as docker container
```bash
make postgres
```
- run 'make createdb' to create a database
```bash
make createdb
```
- run 'make migrate' to perform database migration
```bash
make migrate
```

## 7. Unit Testing

Run 'make test' to run the unit tests:
```bash
make test
```
## 8. Start web server
- run 'make server' to start the service
```bash
make server
```

## 9. Endpoint requests
- POST application/json Request to add a new repository
``` 
curl -d '{"name": "GoogleChrome/chromium-dashboard"}'\
  -H "Content-Type: application/json" \
  -X POST http://localhost:5000/repository \
```

- GET Request to fetch all the repositories on the database
```
curl -L \
  -X GET http://localhost:5000/repositories \
```

- GET Request to fetch all the commits fetched from github API for any repo using repository Id 
```
curl \
  -X GET http://localhost:5000/commits/5846c0f0-81f5-45e3-9d4a-cfc6fe4f176a \
```

- GET Request to get repository metadata using repository id. 
``` 
curl -L \
  -X GET http://localhost:5000/repository/5846c0f0-81f5-45e3-9d4a-cfc6fe4f176a \
```

- GET Request to fetch N (as limit) top commit authors of the any added repository using its repository id with limit as query param
```
curl -L \
  -X GET http://localhost:5000/top-authors/5846c0f0-81f5-45e3-9d4a-cfc6fe4f176a?limit=5 \
```
  
