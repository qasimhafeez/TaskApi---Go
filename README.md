### Sanity testing and Starting a Server
go mod tidy
go run ./cmd/api

### Health
curl -s localhost:8080/health

### Create
curl -s -X POST localhost:8080/v1/tasks \
  -H "Content-Type: application/json" \
  -d '{"id":"t1","title":"Write blog","status":"todo"}'

### Get / List / Update / Delete
curl -s localhost:8080/v1/tasks/t1

curl -s localhost:8080/v1/tasks

curl -s -X PUT localhost:8080/v1/tasks/t1 -H "Content-Type: application/json" -d '{"title":"Write Go blog","status":"doing"}'

curl -i -X DELETE localhost:8080/v1/tasks/t1
