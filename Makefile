.PHONY: proto-gen run run-zero-downtime run-migrate down-migrate run-dump run-rollback

proto-gen:
	protoc --go_out=app/pkg/proto --go-grpc_out=app/pkg/proto app/pkg/proto/*.proto

run:
	docker-compose up -d --build

run-zero-downtime:
	docker-compose up -d --no-deps --build app

run-migrate:
	./migrator.sh up

down-migrate:
	./migrator.sh down -n $(num)

run-dump:
	./dumper.sh

run-rollback:
	docker-compose down && docker-compose up -d
