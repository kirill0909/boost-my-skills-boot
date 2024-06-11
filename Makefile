proto-gen:
	protoc --go_out=app/pkg/proto --go-grpc_out=app/pkg/proto app/pkg/proto/*.proto

run:
	docker-compose up -d --build

run-zero-downtime:
	docker-compose up -d --no-deps --build app

run-migrate:
	./migrator.sh up

down-migrate:
	./migrator.sh down

run-dump:
	./dumper.sh
