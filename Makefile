.PHONY: all
all: clean mock_gen build test

.PHONY: prepare_env
prepare_env:
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

.PHONY: build
build: build_agent build_server

.PHONY: build_agent
build_agent:
	@echo "\n### $@"
	@cd cmd/agent && go build .

.PHONY: build_server
build_server:
	@echo "\n### $@"
	@cd cmd/server && go build .

.PHONY: test
test: test_units test_static test_devops

.PHONY: test_units
test_units: 
	@echo "\n### $@"
	@go test ./... -v --count 1

.PHONY: test_static
test_static:
	@echo "\n### $@"
	@go vet -vettool=./statictest ./...

.PHONY: test_devops
test_devops: build
	@echo "\n### $@"
	@echo "DON'T FORGET TO START postgres.sh\n"
	@./devopstest -test.v -test.run=^TestIteration1$$ -agent-binary-path=cmd/agent/agent
	@./devopstest -test.v -test.run=^TestIteration2[b]*$$ -source-path=. -binary-path=cmd/server/server
	@./devopstest -test.v -test.run=^TestIteration3[b]*$$ -source-path=. -binary-path=cmd/server/server
	@./devopstest -test.v -test.run=^TestIteration4$$ -source-path=. -binary-path=cmd/server/server -agent-binary-path=cmd/agent/agent
	@export SERVER_PORT=$$(./random unused-port) && \
	 export ADDRESS="localhost:$${SERVER_PORT}" && \
	 ./devopstest -test.v -test.run=^TestIteration5$$ \
	 -source-path=. \
	 -agent-binary-path=cmd/agent/agent \
	 -binary-path=cmd/server/server \
	 -server-port=$${SERVER_PORT} 
	@export SERVER_PORT=$$(./random unused-port) && \
	 export ADDRESS="localhost:$${SERVER_PORT}" && \
	 export TEMP_FILE=$$(./random tempfile) && \
     ./devopstest -test.v -test.run=^TestIteration6$$ \
	 -source-path=. \
	 -agent-binary-path=cmd/agent/agent \
	 -binary-path=cmd/server/server \
	 -server-port=$${SERVER_PORT} \
	 -database-dsn='postgres://postgres:postgres@127.0.0.1:5432/praktikum?sslmode=disable' \
	 -file-storage-path=$${TEMP_FILE}
	@export SERVER_PORT=$$(./random unused-port) && \
	 export ADDRESS="localhost:$${SERVER_PORT}" && \
	 export TEMP_FILE=$$(./random tempfile) && \
	 ./devopstest -test.v -test.run=^TestIteration7$$ \
	 -source-path=. \
	 -agent-binary-path=cmd/agent/agent \
	 -binary-path=cmd/server/server \
	 -server-port=$${SERVER_PORT} \
     -database-dsn='postgres://postgres:postgres@127.0.0.1:5432/praktikum?sslmode=disable' \
	 -file-storage-path=$${TEMP_FILE}
	@export SERVER_PORT=$$(./random unused-port) && \
	 export ADDRESS="localhost:$${SERVER_PORT}" && \
	 export TEMP_FILE=$$(./random tempfile) && \
	 ./devopstest -test.v -test.run=^TestIteration8 \
	 -source-path=. \
	 -agent-binary-path=cmd/agent/agent \
	 -binary-path=cmd/server/server \
	 -server-port=$${SERVER_PORT} \
	 -database-dsn='postgres://postgres:postgres@127.0.0.1:5432/praktikum?sslmode=disable' \
	 -file-storage-path=$${TEMP_FILE}
	@export SERVER_PORT=$$(./random unused-port) && \
	 export ADDRESS="localhost:$${SERVER_PORT}" && \
	 export TEMP_FILE=$$(./random tempfile) && \
	 ./devopstest -test.v -test.run=^TestIteration9$$ \
	 -source-path=. \
	 -agent-binary-path=cmd/agent/agent \
	 -binary-path=cmd/server/server \
	 -server-port=$${SERVER_PORT} \
	 -file-storage-path=$${TEMP_FILE} \
	 -database-dsn='postgres://postgres:postgres@127.0.0.1:5432/praktikum?sslmode=disable' \
	 -key="$${TEMP_FILE}"
	@export SERVER_PORT=$$(./random unused-port) && \
	 export ADDRESS="localhost:$${SERVER_PORT}" && \
	 export TEMP_FILE=$$(./random tempfile) && \
	 ./devopstest -test.v -test.run=^TestIteration10[b]*$$ \
	 -source-path=. \
	 -agent-binary-path=cmd/agent/agent \
	 -binary-path=cmd/server/server \
	 -server-port=$${SERVER_PORT} \
	 -database-dsn='postgres://postgres:postgres@127.0.0.1:5432/praktikum?sslmode=disable' \
	 -key="$${TEMP_FILE}"
	@export SERVER_PORT=$$(./random unused-port) && \
	 export ADDRESS="localhost:$${SERVER_PORT}" && \
	 export TEMP_FILE=$$(./random tempfile) && \
	 ./devopstest -test.v -test.run=^TestIteration11$$ \
	 -source-path=. \
	 -agent-binary-path=cmd/agent/agent \
	 -binary-path=cmd/server/server \
	 -server-port=$${SERVER_PORT} \
	 -database-dsn='postgres://postgres:postgres@127.0.0.1:5432/praktikum?sslmode=disable' \
	 -key="$${TEMP_FILE}"
	@export SERVER_PORT=$$(./random unused-port) && \
	 export ADDRESS="localhost:$${SERVER_PORT}" && \
	 export TEMP_FILE=$$(./random tempfile) && \
	 ./devopstest -test.v -test.run=^TestIteration12$$ \
	 -source-path=. \
	 -agent-binary-path=cmd/agent/agent \
	 -binary-path=cmd/server/server \
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
	 -agent-binary-path=cmd/agent/agent \
	 -binary-path=cmd/server/server \
	 -server-port=$${SERVER_PORT} \
	 -file-storage-path=$${TEMP_FILE} \
	 -database-dsn='postgres://postgres:postgres@127.0.0.1:5432/praktikum?sslmode=disable' \
	 -key="$${TEMP_FILE}"
	@go test -v -race ./...

.PHONY: clean
clean:
	@echo "\n### $@"
	@rm -rf cmd/agent/agent
	@rm -rf cmd/server/server
