update-lambda:
	scripts/deploy_lambda.sh

gensql:
	sqlc generate -f db/sqlc.yaml