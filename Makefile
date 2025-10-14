build:
	go build -o bin/tbot main.go && sudo cp bin/tbot /usr/local/bin/tbot

gensql:
	sqlc generate -f internal/db/sqlc.yaml

up:
	goose up

down:
	goose down