lint:
	cd backend; golangci-lint run
	cd frontend; npm run lint

test:
	cd backend; go test ./...
	cd frontend; npm run test

coverage:
	cd backend; go test ./... --cover
	cd frontend; npm run coverage

build:
	docker-compose build

run:
	docker-compose up
