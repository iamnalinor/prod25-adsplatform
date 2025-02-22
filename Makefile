lint:
	cd backend && go vet ./... && go fmt ./...
	cd tg_bot && ruff check --fix && ruff format

test: up
	cd backend && go test ./...
	cd backend/tests && py.test

up:
	docker-compose up -d --build

docs:
	cd backend && swag init

data:
	cd backend/tests && pytest test_07_stress_ads.py -k "not cleanup"

views:
	python3 tg_bot/generate_views.py

.PHONY: lint test up docs data views