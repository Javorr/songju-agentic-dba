up:
	docker compose up -d

down:
	docker compose down

run: up
	go run ./cmd/web

load-test:
	k6 run tests/load_test.js
