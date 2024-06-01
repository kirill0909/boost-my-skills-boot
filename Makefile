proto-gen:
	protoc --go_out=app/pkg/proto --go-grpc_out=app/pkg/proto app/pkg/proto/*.proto

run:
	docker-compose up --build
