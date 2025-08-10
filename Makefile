up:
	docker compose up --build -d
down:
	docker compose down -v
producer:
	go run ./cmd/producer/main.go
consumer:
	go run ./cmd/consumer/main.go