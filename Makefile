build:
	go build -o ~/bin/tbot main.go

gensql:
	sqlc generate -f internal/db/sqlc.yaml

up:
	goose up

down:
	goose down