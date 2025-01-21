# Yangder (Yang Finder)
Yang (/jeɪŋ/) means couple in Javanese slang. Yangder is a find-a-partner application that matches you with your potential mate.

## Documentation
Please refer to this [pdf](docs/doc.pdf)

## Structure

```
├── auth // Authentication handler with JWT
├── config // Configuration for database
├── controllers // Business logic
├── docs // Project Documentation
├── go.mod
├── gorm.db
├── go.sum
├── main.go
├── migrate // Database migration tool
├── models // Entities
├── README.md
└── scheduler // To execute scheduled action
```

## Run
1. Make sure `go` is installed.
2. Initialize the database using SQLite
```
go run migrate/migrate.go
```
3. Run with
```
go run main.go
```

## Testing
Testing can be done via automatic unit test or manual testing by Postman.

1. Unit Test
```
go test ./...
```
2. Download this [Postman collection](docs/postman/yangnder.postman_collection.json)

## Note
This project to showcase my proficiency in Go, Gin, GORM. This is my first time actually creating something with Go.