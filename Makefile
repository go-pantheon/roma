GOCMD=GO111MODULE=on CGO_ENABLED=0 go
GOBUILD=${GOCMD} build

.PHONY: init
# Initialize environment
init:
	pre-commit install
	go install github.com/google/go-licenses@latest
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.1.6
	go install github.com/google/wire/cmd/wire@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install github.com/envoyproxy/protoc-gen-validate@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@latest
	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-errors/v2@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
	go install github.com/favadi/protoc-go-inject-tag@latest

.PHONY: generate
# Generate all
generate: proto api wire

.PHONY: version
# Show the generated version
version:
	@find app -type d -depth 1 -print | xargs -L 1 bash -c 'echo "version: $$0" && cd "$$0" && $(MAKE) version'

.PHONY: wire
# Generate wire
wire:
	@find app -type d -depth 1 -print | xargs -L 1 bash -c 'echo "wire: $$0" && cd "$$0" && $(MAKE) wire'

.PHONY: proto
# Generate internal proto pb.go files.
proto:
	buf generate --template=buf/buf.gen.app.yaml
	buf generate --template=buf/buf.gen.mercury.yaml

.PHONY: api
# Generate api/client and api/server pb.go files.
api:
	buf generate --template=buf/buf.gen.client.yaml
	buf generate --template=buf/buf.gen.server.yaml

.PHONY: db
# Generate api/db pb.go files and inject proto tag
db:
	@buf generate --template=buf/buf.gen.db.yaml && \
	cd gen/api/db/ && find . -type d -print0 | while IFS= read -r -d '' dir; do \
		if [[ -z $$(find "$$dir" -mindepth 1 -type d) ]]; then \
			if ls "$$dir"/*.pb.go &>/dev/null; then \
				echo "Running protoc-go-inject-tag in $$dir"; \
				protoc-go-inject-tag -input="$$dir/*.pb.go"; \
			fi; \
		fi; \
	done
	bin/tools/gen-api-db

.PHONY: run
# Run app
run:
	@if [ -z "$(app)" ]; then \
  	echo "error: app must exist. ex: app=player"; \
	else \
		echo "run: app/$(app)" && cd app/$(app) && $(MAKE) run; \
	fi

.PHONY: build
# Build app execute file. Use app=app_name to build specific service.
build:
	@if [ -z "$(app)" ]; then \
		find app -type d -depth 1 -print | xargs -L 1 bash -c 'echo "build: $$0" && cd "$$0" && $(MAKE) build'; \
	else \
		echo "build: app/$(app)" && cd app/$(app) && $(MAKE) build; \
	fi

.PHONY: start
# Start all app services. Use app=app_name to start specific service.
start: stop build
	@if [ -z "$(app)" ]; then \
		find app -type d -depth 1 -print | xargs -L 1 bash -c 'echo "start: $$0" && cd "$$0" && $(MAKE) start'; \
	else \
		echo "start: app/$(app)" && cd app/$(app) && $(MAKE) start; \
	fi

.PHONY: stop
# Stop all app services. Use app=app_name to stop specific service.
stop:
	@if [ -z "$(app)" ]; then \
		find app -type d -depth 1 -print | xargs -L 1 bash -c 'echo "stop: $$0" && cd "$$0" && $(MAKE) stop'; \
	else \
		echo "stop: app/$(app)" && cd app/$(app) && $(MAKE) stop; \
	fi

.PHONY: log
# Tail app service log file. Must use app=app_name to tail specific service.
log:
	@if [ -z "$(app)" ]; then \
  	echo "error: app must exist. ex: app=player"; \
	else \
		echo "log: app/$(app)" && cd app/$(app) && $(MAKE) log; \
	fi

.PHONY: mercury
# Run mercury client
mercury: mercury-build mercury-start mercury-log

.PHONY: mercury-build
mercury-build:
	cd mercury && echo "build: mercury" \
	  && rm -rf bin/mercury && mkdir -p bin/mercury && $(GOBUILD) -ldflags "-X main.Version=0.0.1" -o bin/mercury/client ./cmd

.PHONY: mercury-start
# Start mercury client
mercury-start:
	cd mercury && echo "start: mercury" \
	  && nohup bin/mercury/client -conf=configs -gamedata ../gen/gamedata/json > bin/mercury/debug.log &

.PHONY: mercury-stop
# Kill mercury process
mercury-stop:
	-pkill -f mercury/bin/mercury

.PHONY: mercury-log
# tail -f mercury/bin/mercury/debug.log
mercury-log:
	cd mercury && echo "log: mercury" && tail -f bin/mercury/debug.log

.PHONY: tools
# Build all tools to directory:bin/tools
tools:
	@rm -rf bin/tools
	@mkdir -p bin/tools
	@echo "build: gen-api-db" && $(GOBUILD) -o bin/tools/gen-api-db -ldflags "-X main.Version=0.0.1"  ./vulcan/app/api/db/cmd
	@echo "build: gen-api-client" && $(GOBUILD) -o bin/tools/gen-api-client -ldflags "-X main.Version=0.0.1"  ./vulcan/app/api/client/cmd
	@echo "build: gen-mercury" && $(GOBUILD) -o bin/tools/gen-mercury -ldflags "-X main.Version=0.0.1"  ./vulcan/app/mercury/cmd
	@echo "build: gen-data-json" && $(GOBUILD) -o bin/tools/gen-data-json -ldflags "-X main.Version=0.0.1"  ./vulcan/app/gamedata/cmd/json
	@echo "build: gen-data-base" && $(GOBUILD) -o bin/tools/gen-data-base -ldflags "-X main.Version=0.0.1"  ./vulcan/app/gamedata/cmd/base
	@echo "build: gen-datas" && $(GOBUILD) -o bin/tools/gen-datas -ldflags "-X main.Version=0.0.1"  ./vulcan/app/gamedata/cmd/data

.PHONY: gen-api-client
# Generate api/client code
gen-api-client:
	@echo "run: gen-api-client" && bin/tools/gen-api-client

.PHONY: gen-api-db
# Generate api/db code
gen-api-db:
	@echo "run: gen-api-db" && bin/tools/gen-api-db

.PHONY: gen-mercury
# Generate mercury client code
gen-mercury:
	@echo "run: gen-mercury" && bin/tools/gen-mercury

.PHONY: gen-data-json
# Generate data json file
gen-data-json:
	@echo "run: gen-data-json" && bin/tools/gen-data-json

.PHONY: gen-data-base
# Generate data base code
gen-data-base:
	@echo "run: gen-data-base" && bin/tools/gen-data-base

.PHONY: gen-datas
# Generate data code
gen-datas:
	@echo "run: gen-datas" && bin/tools/gen-datas

.PHONY: gen-all-data
# Generate all data code
gen-all-data: gen-data-json gen-data-base gen-datas

.PHONY: test
# Run test and show coverage
test:
	go test -race ./...

.PHONY: vet
# Run vet
vet:
	go vet ./...

.PHONY: license-check
# Run license check
license-check:
	go-licenses check ./...

.PHONY: lint
# Run lint
lint:
	golangci-lint run ./...

# show help
help:
	@echo ''
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help
