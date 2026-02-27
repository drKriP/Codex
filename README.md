# Codex

This repository now contains two simple CRUD microservices and one placeholder for the third service.

## 1) Java microservice (`java-service`)

Spring Boot REST API with in-memory storage.

### Run

```bash
cd java-service
mvn spring-boot:run
```

Service listens on `http://localhost:8080`.

### Endpoints

- `GET /items`
- `POST /items`
- `GET /items/{id}`
- `PUT /items/{id}`
- `DELETE /items/{id}`

Example request body:

```json
{
  "name": "item-1",
  "description": "first item"
}
```

## 2) Go microservice (`go-service`)

Net/http REST API with in-memory storage.

### Run

```bash
cd go-service
go run .
```

Service listens on `http://localhost:8081`.

### Endpoints

- `GET /items`
- `POST /items`
- `GET /items/{id}`
- `PUT /items/{id}`
- `DELETE /items/{id}`

## 3) Ruby/Express-style microservice (`ruby-express-service`)

Reserved for later implementation (as requested).
