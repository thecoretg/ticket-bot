create-bin-dir:
	mkdir -p bin

build-cli: create-bin-dir
	go build -o bin/cli ./cmd/cli && cp bin/cli ~/go/bin/tbot

gensql:
	sqlc generate

runserver:
	go run ./cmd/server

test-db-up:
	docker compose -f ./docker/docker-compose-db.yml up -d

test-db-down:
	docker compose -f ./docker/docker-compose-db.yml down -v

docker-build:
	docker buildx build --platform=linux/amd64 -t ticketbot:v1.2 --load -f ./docker/DockerfileMain .

deploy-lightsail: docker-build
	aws lightsail push-container-image \
	--region us-west-2 \
	--service-name ticketbot \
	--label ticketbot-server \
	--image ticketbot:latest

lightsail-logs:
	aws lightsail get-container-log \
	--service-name ticketbot \
	--container-name ticketbot \
	--output text
