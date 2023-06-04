AGENT_VERSION := 1.0.0
SERVER_VERSION := 1.0.0
BUILD_DATE := $(shell date +'%d.%m.%Y %H:%M:%S')
BUILD_COMMIT := $(shell git rev-parse --short HEAD)

.PHONY: all
all: clean mock_gen build test

.PHONY: prepare_env
prepare_env: build_mylinter
	@echo "\n### $@"
	@go install github.com/golang/mock/mockgen@v1.6.0
	@wget https://github.com/Yandex-Practicum/go-autotests/releases/download/v0.9.6/devopstest -O ./devopstest
	@chmod u+x ./devopstest
	@wget https://github.com/Yandex-Practicum/go-autotests/releases/download/v0.9.6/statictest -O ./statictest
	@chmod u+x ./statictest
	@wget https://github.com/Yandex-Practicum/go-autotests/releases/download/v0.9.6/random -O ./random
	@chmod u+x ./random

.PHONY: mock_gen
mock_gen:
	@echo "\n### $@"
	@mockgen -destination=internal/server/mocks/mock_storage.go -package=mocks github.com/devldavydov/promytheus/internal/server/storage Storage

.PHONY: proto_gen
proto_gen:
	@echo "\n### $@"
	@protoc --go_out=. --go_opt=paths=import --go-grpc_out=. --go-grpc_opt=paths=import internal/grpc/proto/metric.proto

.PHONY: build
build: build_agent build_server

.PHONY: build_agent
build_agent:
	@echo "\n### $@"
	@mkdir -p ./bin
	@cd cmd/agent && \
	 go build \
	 -ldflags "-X main.buildVersion=$(AGENT_VERSION) -X 'main.buildDate=$(BUILD_DATE)' -X main.buildCommit=$(BUILD_COMMIT)" \
	 -o ../../bin/agent .

.PHONY: build_server
build_server:
	@echo "\n### $@"
	@mkdir -p ./bin
	@cd cmd/server && \
	 go build \
	 -ldflags "-X main.buildVersion=$(SERVER_VERSION) -X 'main.buildDate=$(BUILD_DATE)' -X main.buildCommit=$(BUILD_COMMIT)" \
	 -o ../../bin/server .

.PHONY: build_staticlint
build_staticlint:
	@echo "\n### $@"
	@mkdir -p ./bin
	@cd cmd/staticlint && \
	 go build -o ../../bin/staticlint .

.PHONY: test
test: test_units test_static test_devops

.PHONY: test_units
test_units: 
	@echo "\n### $@"
	@echo "DON'T FORGET TO START postgres.sh\n"
	@export TEST_DATABASE_DSN=postgres://postgres:postgres@127.0.0.1:5432/praktikum?sslmode=disable && \
	 go test ./... -v --count 1

.PHONY: test_static
test_static: clean build_staticlint
	@echo "\n### $@"
	@go vet -vettool=./statictest ./...
	@bin/staticlint -json -test=false ./...

.PHONY: test_devops
test_devops: build
	@echo "\n### $@"
	@echo "DON'T FORGET TO START postgres.sh\n"
	@./devopstest -test.v -test.run=^TestIteration1$$ -agent-binary-path=bin/agent
	@./devopstest -test.v -test.run=^TestIteration2[b]*$$ -source-path=. -binary-path=bin/server
	@./devopstest -test.v -test.run=^TestIteration3[b]*$$ -source-path=. -binary-path=bin/server
	@./devopstest -test.v -test.run=^TestIteration4$$ -source-path=. -binary-path=bin/server -agent-binary-path=bin/agent
	@export SERVER_PORT=$$(./random unused-port) && \
	 export ADDRESS="localhost:$${SERVER_PORT}" && \
	 ./devopstest -test.v -test.run=^TestIteration5$$ \
	 -source-path=. \
	 -agent-binary-path=bin/agent \
	 -binary-path=bin/server \
	 -server-port=$${SERVER_PORT} 
	@export SERVER_PORT=$$(./random unused-port) && \
	 export ADDRESS="localhost:$${SERVER_PORT}" && \
	 export TEMP_FILE=$$(./random tempfile) && \
     ./devopstest -test.v -test.run=^TestIteration6$$ \
	 -source-path=. \
	 -agent-binary-path=bin/agent \
	 -binary-path=bin/server \
	 -server-port=$${SERVER_PORT} \
	 -database-dsn='postgres://postgres:postgres@127.0.0.1:5432/praktikum?sslmode=disable' \
	 -file-storage-path=$${TEMP_FILE}
	@export SERVER_PORT=$$(./random unused-port) && \
	 export ADDRESS="localhost:$${SERVER_PORT}" && \
	 export TEMP_FILE=$$(./random tempfile) && \
	 ./devopstest -test.v -test.run=^TestIteration7$$ \
	 -source-path=. \
	 -agent-binary-path=bin/agent \
	 -binary-path=bin/server \
	 -server-port=$${SERVER_PORT} \
     -database-dsn='postgres://postgres:postgres@127.0.0.1:5432/praktikum?sslmode=disable' \
	 -file-storage-path=$${TEMP_FILE}
	@export SERVER_PORT=$$(./random unused-port) && \
	 export ADDRESS="localhost:$${SERVER_PORT}" && \
	 export TEMP_FILE=$$(./random tempfile) && \
	 ./devopstest -test.v -test.run=^TestIteration8 \
	 -source-path=. \
	 -agent-binary-path=bin/agent \
	 -binary-path=bin/server \
	 -server-port=$${SERVER_PORT} \
	 -database-dsn='postgres://postgres:postgres@127.0.0.1:5432/praktikum?sslmode=disable' \
	 -file-storage-path=$${TEMP_FILE}
	@export SERVER_PORT=$$(./random unused-port) && \
	 export ADDRESS="localhost:$${SERVER_PORT}" && \
	 export TEMP_FILE=$$(./random tempfile) && \
	 ./devopstest -test.v -test.run=^TestIteration9$$ \
	 -source-path=. \
	 -agent-binary-path=bin/agent \
	 -binary-path=bin/server \
	 -server-port=$${SERVER_PORT} \
	 -file-storage-path=$${TEMP_FILE} \
	 -database-dsn='postgres://postgres:postgres@127.0.0.1:5432/praktikum?sslmode=disable' \
	 -key="$${TEMP_FILE}"
	@export SERVER_PORT=$$(./random unused-port) && \
	 export ADDRESS="localhost:$${SERVER_PORT}" && \
	 export TEMP_FILE=$$(./random tempfile) && \
	 ./devopstest -test.v -test.run=^TestIteration10[b]*$$ \
	 -source-path=. \
	 -agent-binary-path=bin/agent \
	 -binary-path=bin/server \
	 -server-port=$${SERVER_PORT} \
	 -database-dsn='postgres://postgres:postgres@127.0.0.1:5432/praktikum?sslmode=disable' \
	 -key="$${TEMP_FILE}"
	@export SERVER_PORT=$$(./random unused-port) && \
	 export ADDRESS="localhost:$${SERVER_PORT}" && \
	 export TEMP_FILE=$$(./random tempfile) && \
	 ./devopstest -test.v -test.run=^TestIteration11$$ \
	 -source-path=. \
	 -agent-binary-path=bin/agent \
	 -binary-path=bin/server \
	 -server-port=$${SERVER_PORT} \
	 -database-dsn='postgres://postgres:postgres@127.0.0.1:5432/praktikum?sslmode=disable' \
	 -key="$${TEMP_FILE}"
	@export SERVER_PORT=$$(./random unused-port) && \
	 export ADDRESS="localhost:$${SERVER_PORT}" && \
	 export TEMP_FILE=$$(./random tempfile) && \
	 ./devopstest -test.v -test.run=^TestIteration12$$ \
	 -source-path=. \
	 -agent-binary-path=bin/agent \
	 -binary-path=bin/server \
	 -server-port=$${SERVER_PORT} \
	 -database-dsn='postgres://postgres:postgres@127.0.0.1:5432/praktikum?sslmode=disable' \
	 -key="$${TEMP_FILE}"
	@./devopstest -test.v -test.run=^TestIteration13$$ \
	 -source-path=.
	@export SERVER_PORT=$$(./random unused-port) && \
	 export ADDRESS="localhost:$${SERVER_PORT}" && \
	 export TEMP_FILE=$$(./random tempfile) && \
	 ./devopstest -test.v -test.run=^TestIteration14$$ \
	 -source-path=. \
	 -agent-binary-path=bin/agent \
	 -binary-path=bin/server \
	 -server-port=$${SERVER_PORT} \
	 -file-storage-path=$${TEMP_FILE} \
	 -database-dsn='postgres://postgres:postgres@127.0.0.1:5432/praktikum?sslmode=disable' \
	 -key="$${TEMP_FILE}"
	@go test -v -race ./...

.PHONY: test_bench
test_bench: build
	@echo "\n### $@"
	@mkdir -p profiles
	@cd internal/server/handler/metric && go test . -run=$^ -bench=. -memprofile=mem.pprof
	@mv internal/server/handler/metric/mem.pprof profiles/

.PHONY: test_cover
test_cover:
	@echo "\n### $@"
	@echo "DON'T FORGET TO START postgres.sh\n"
	@export TEST_DATABASE_DSN=postgres://postgres:postgres@127.0.0.1:5432/praktikum?sslmode=disable && \
	 go test ./... -coverprofile cover.html -v --count 1
	@go tool cover -html=cover.html

.PHONY: diff_bench
diff_bench:
	@echo "\n### $@"
	@go tool pprof -top -diff_base=profiles/base.pprof profiles/result.pprof

.PHONY: run_docs
run_docs:
	@echo "See docs in http://localhost:8080/pkg/github.com/devldavydov/promytheus/?m=all"
	@godoc -http=:8080

.PHONY: swagger_gen
swagger_gen:
	@echo "\n### $@"
	@swag init -g handler.go -d internal/server/handler/metric --parseDependency --output ./swagger/
	@swag fmt

.PHONY: tls_gen
tls_gen:
	@echo "\n### $@"
	@mkdir -p tls
	@rm -rf tls/*
	@echo "subjectAltName=IP:127.0.0.1" > tls/server-ext.cnf
	@echo "> Generate CA's private key and self-signed certificate"
	@openssl req -x509 -newkey rsa:4096 -days 365 -nodes -keyout tls/ca-key.pem -out tls/ca-cert.pem -subj "/C=RU/ST=Moscow/L=Moscow/O=Yandex/OU=Praktikum/CN=promytheus"
	@echo "> CA's self-signed certificate"
	@openssl x509 -in tls/ca-cert.pem -noout -text
	@echo "> Generate web server's private key and certificate signing request (CSR)"
	@openssl req -newkey rsa:4096 -nodes -keyout tls/server-key.pem -out tls/server-req.pem -subj "/C=RU/ST=Moscow/L=Moscow/O=Yandex/OU=Praktikum/CN=promytheus"
	@echo "> Use CA's private key to sign web server's CSR and get back the signed certificate"
	@openssl x509 -req -in tls/server-req.pem -days 365 -CA tls/ca-cert.pem -CAkey tls/ca-key.pem -CAcreateserial -out tls/server-cert.pem -extfile tls/server-ext.cnf
	@echo "> Server's signed certificate"
	@openssl x509 -in tls/server-cert.pem -noout -text

.PHONY: clean
clean:
	@echo "\n### $@"
	@rm -rf ./bin	
