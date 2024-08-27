# GitHub API Service## Description

This service fetches data from GitHub's public APIs to retrieve repository commits, saves the data in a persistent store, and continuously monitors the repository for changes at a set interval.

## Requirements- Golang 1.20+
- PostgreSQL

## 1. Clone the repository and download go dependencies
```bash
git clone https://github.com/kenmobility/github-api-hex.git
cd github-api-hex
go mod tidy
```
## 2. Duplicate .env.example file and rename to .env
```bash
cp .env.example .env
```

## 3. environment variables:
- the .env file already has default variables that the program needs to run but for security reason, you can change the variables to your custom details.
- to setup GitHub API token in order to change the GIT_HUB_TOKEN env variable, go to  [https://github.com/](GitHub) to set up a GitHub API token.

## 4 Open Docker desktop application
- Ensure that docker desktop is started and running 

## 5. Run makefile commands 
- run 'make postgres' to pull and run PostgreSQL instance as docker container
```bash
make postgres
```
- run 'make createdb' to create a database
```bash
make createdb
```
## 6. Testing (optional)

Run 'make test' to run the unit tests:
```bash
make test
```
## 7. Start web server
- run 'make server' to start the service
```bash
make server
```

## 8. Endpoint requests
- POST application/json Request to add a new repository
``` 
curl -d '{"name": "GoogleChrome/chromium-dashboard","description": "","url": "https://github.com/GoogleChrome/chromium-dashboard"}'\
  -H "Content-Type: application/json" \
  -X POST http://127.0.0.1:5000/repository \
```

- GET Request to fetch all the repositories on the database
```
curl -L \
  -X GET http://127.0.0.1:5000/repositories \
```

- GET Request to fetch all the commits fetched from github API for any repo using repository public Id 
```
curl \
  -X GET http://127.0.0.1:5000/commits/5846c0f0-81f5-45e3-9d4a-cfc6fe4f176a \
```

- POST application/json Request to set a new repository to track an added repository using its repository public id. 
- payload only requires 'repo-public_id' field, others are optional.
- if a valid start_date is passed, date value is used as 'since' while the service is fetching commits for that particular repository else the default start date in the config is used.
- if a valid end_date field value is passed, date value is used as 'until' while the service is fetching commits for that particular repository else the default end date in the config is used.
``` 
curl -d '{"repo_public_id": "5846c0f0-81f5-45e3-9d4a-cfc6fe4f176a","start_date": "","end_date": ""}'\
  -H "Content-Type: application/json" \
  -X POST http://127.0.0.1:5000/repository/track \
```

- GET Request to fetch N (as limit) top commit authors of the any added repository using its repository id with limit as query param
```
curl -L \
  -X GET http://127.0.0.1:5000/top-authors/5846c0f0-81f5-45e3-9d4a-cfc6fe4f176a?limit=5 \
```
  
