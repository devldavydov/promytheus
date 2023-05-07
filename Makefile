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

.PHONY: build
build: build_agent build_server

.PHONY: build_agent
build_agent:
	@echo "\n### $@"
	@cd cmd/agent && \
	 go build \
	 -ldflags "-X main.buildVersion=$(AGENT_VERSION) -X 'main.buildDate=$(BUILD_DATE)' -X main.buildCommit=$(BUILD_COMMIT)" \
	 .

.PHONY: build_server
build_server:
	@echo "\n### $@"
	@cd cmd/server && \
	 go build \
	 -ldflags "-X main.buildVersion=$(SERVER_VERSION) -X 'main.buildDate=$(BUILD_DATE)' -X main.buildCommit=$(BUILD_COMMIT)" \
	 .


.PHONY: build_staticlint
build_staticlint:
	@echo "\n### $@"
	@cd cmd/staticlint && go build .

.PHONY: test
test: test_units test_static test_devops

.PHONY: test_units
test_units: 
	@echo "\n### $@"
	@echo "DON'T FORGET TO START postgres.sh\n"
	@export TEST_DATABASE_DSN=postgres://postgres:postgres@127.0.0.1:5432/praktikum?sslmode=disable && \
	 go test ./... -v --count 1

.PHONY: test_static
test_static: build_staticlint
	@echo "\n### $@"
	@go vet -vettool=./statictest ./...
	@cmd/staticlint/staticlint -json -test=false ./...

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
	@go test ./... -coverprofile cover.html -v --count 1
	@go tool cover -html=cover.html

.PHONY: diff_bench
diff_bench:
	@echo "\n### $@"
	@go tool pprof -top -diff_base=profiles/base.pprof profiles/result.pprof

.PHONY: run_docs
run_docs:
	@echo "See docs in http://localhost:8080/pkg/github.com/devldavydov/promytheus/?m=all"
	@godoc -http=:8080

.PHONY: gen_swagger
gen_swagger:
	@echo "\n### $@"
	@swag init -g handler.go -d internal/server/handler/metric --parseDependency --output ./swagger/
	@swag fmt

.PHONY: clean
clean:
	@echo "\n### $@"
	@rm -rf cmd/agent/agent
	@rm -rf cmd/server/server
