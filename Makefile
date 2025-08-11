up:
	docker compose up --build -d
down:
	docker compose down -v
producer:
	./cmd/producer/producer
consumer:
	./cmd/consumer1/consumer1 &
	./cmd/consumer2/consumer2 &
	./cmd/consumer3/consumer3 &
	wait
cbin:
	go build -o ./cmd/producer/producer ./cmd/producer &
	go build -o ./cmd/consumer1/consumer1 ./cmd/consumer1 &
	go build -o ./cmd/consumer2/consumer2 ./cmd/consumer2 &
	go build -o ./cmd/consumer3/consumer3 ./cmd/consumer3 &
	go build -o ./cmd/consumerctl/consumerctl ./cmd/consumerctl &
dbin:
	rm ./cmd/producer/producer &
	rm ./cmd/consumer1/consumer1 &
	rm ./cmd/consumer2/consumer2 &
	rm ./cmd/consumer3/consumer3 &
	rm ./cmd/consumerctl/consumerctl
stop_consumer1:
	./cmd/consumerctl/consumerctl consumer1.pid
stop_consumer2:
	./cmd/consumerctl/consumerctl consumer2.pid
stop_consumer3:
	./cmd/consumerctl/consumerctl consumer3.pid