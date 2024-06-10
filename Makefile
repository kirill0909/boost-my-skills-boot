proto-gen:
	protoc --go_out=app/pkg/proto --go-grpc_out=app/pkg/proto app/pkg/proto/*.proto

run:
	docker-compose up --build

zero-downtime-run:
	docker-compose up -d --no-deps --build app
