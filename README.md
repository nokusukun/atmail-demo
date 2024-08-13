
# Atmail API Demo

## Setup
### Database
Create a MySQL instance and the database using Docker:

```go
$ docker run --name mysql-container -e MYSQL_ROOT_PASSWORD=my-secret-pw -p 3306:3306 -d mysql:latest
$ docker exec -it mysql-container mysql -u root -p
mysql> CREATE DATABASE demo;
mysql> EXIT;
```

## Run
```bash
$ go run ./cmd/server
```

## Docs
```
http://localhost:8080/docs.html
```
## OpenAPI Documentation
The JSON and YAML documentation is available at:
```go
./docs/openapi.json
./docs/openapi.yaml
```
