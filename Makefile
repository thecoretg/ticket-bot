install-cli:
	go build -o ~/bin/ticketbot ./cmd/cli/main.go

gensql:
	sqlc generate -f internal/db/sqlc.yaml

run:
	go run main.go

up:
	goose up

down:
	goose down