up:
	docker compose up -d

down:
	docker compose down

run: up
	go run main.go

load-test:
	k6 run tests/load_test.js
