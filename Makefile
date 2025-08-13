install-cli:
	go build -o ~/bin/tbot ./cmd/cli/main.go

update-lambda:
	scripts/deploy_lambda.sh

tail-lambda-logs:
	aws logs tail /aws/lambda/ticketbot --follow

gensql:
	sqlc generate -f db/sqlc.yaml

init-hooks:
	go run cmd/cli/main.go init-hooks

up:
	goose up

down:
	goose down