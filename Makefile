run:
	go run main.go default -c ./configs/config.toml

debug:
	go run main.go debug -c ./configs/config.toml

build:
	go build -a -installsuffix cgo -o ./main .

coverage:
	go test ./... -coverprofile=coverage.out

report:
	go tool cover -html=coverage.out

proto:
	protoc \
	-I ./pkg/proto \
	--go_opt=module=github.com/dathuynh1108/clean-arch-base --go_out=. \
	./pkg/proto/*.proto

proto-grpc:
	protoc \
	-I ./pkg/proto \
	--go_opt=module=github.com/dathuynh1108/clean-arch-base --go_out=. \
	--go-grpc_opt=module=github.com/dathuynh1108/clean-arch-base --go-grpc_out=. \
	./pkg/proto/*.proto

