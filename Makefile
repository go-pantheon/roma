GOCMD=GO111MODULE=on CGO_ENABLED=0 go
GOBUILD=${GOCMD} build

.PHONY: init
# Initialize environment
init:
	pre-commit install
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.1.6
	go install github.com/google/wire/cmd/wire@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install github.com/envoyproxy/protoc-gen-validate@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-errors/v2@latest
	go install github.com/favadi/protoc-go-inject-tag@latest

.PHONY: version
# Show the generated version
version:
	@find app -type d -depth 1 -print | xargs -L 1 bash -c 'cd "$$0" && pwd && $(MAKE) version'


.PHONY: wire
# Generate wire
wire:
	@find app -type d -depth 1 -print | xargs -L 1 bash -c 'echo "wire: $$0" && cd "$$0" && $(MAKE) wire'
	cd mercury/cmd && wire

.PHONY: proto
# Generate internal proto struct
proto:
	buf generate --template=buf/buf.gen.app.yaml
	buf generate --template=buf/buf.gen.mercury.yaml

.PHONY: api
# Generate API
api:
	buf generate --template=buf/buf.gen.client.yaml
	buf generate --template=buf/buf.gen.server.yaml

.PHONY: db
# Generate API
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

.PHONY: build
# build execute file
build:
	@if [ -z "$(app)" ]; then \
		find app -type d -depth 1 -print | xargs -L 1 bash -c 'cd "$$0" && pwd && $(MAKE) build'; \
	else \
		cd app/$(app) && pwd && $(MAKE) build; \
	fi

.PHONY: run
# start all project services
run: start log

.PHONY: log
# tail -f app/gate/bin/debug.log
log:
	@if [ -z "$(app)" ]; then \
  	    echo "error: app must exist. ex: app=player"; \
	else \
		cd app/$(app) && pwd && $(MAKE) log; \
	fi

.PHONY: start
# start all project services
start: stop build
	@if [ -z "$(app)" ]; then \
		find app -type d -depth 1 -print | xargs -L 1 bash -c 'cd "$$0" && pwd && $(MAKE) start'; \
	else \
		cd app/$(app) && pwd && $(MAKE) start; \
	fi

.PHONY: stop
# stop all project services
stop:
	@if [ -z "$(app)" ]; then \
		find app -type d -depth 1 -print | xargs -L 1 bash -c 'cd "$$0" && pwd && $(MAKE) stop'; \
	else \
		cd app/$(app) && pwd && $(MAKE) stop; \
	fi

.PHONY: docker
# build docker image
docker:
	@if [ -z "$(app)" ]; then \
		find app -type d -depth 1 -print | xargs -L 1 bash -c 'cd "$$0" && pwd && $(MAKE) docker'; \
	else \
		cd app/$(app) && pwd && $(MAKE) docker; \
	fi

.PHONY: test
# test
test:
	go test -v ./... -cover

.PHONY: tools
# build all tools to bin/tools directory
tools:
	rm -rf bin/tools
	mkdir -p bin/tools
	$(GOBUILD) -o bin/tools/gen-api-db -ldflags "-X main.Version=0.0.1"  ./vulcan/app/api/db/cmd
	$(GOBUILD) -o bin/tools/gen-api-client -ldflags "-X main.Version=0.0.1"  ./vulcan/app/api/client/cmd
	$(GOBUILD) -o bin/tools/gen-mercury -ldflags "-X main.Version=0.0.1"  ./vulcan/app/mercury/cmd
	$(GOBUILD) -o bin/tools/gen-data-json -ldflags "-X main.Version=0.0.1"  ./vulcan/app/gamedata/cmd/json
	$(GOBUILD) -o bin/tools/gen-data-base -ldflags "-X main.Version=0.0.1"  ./vulcan/app/gamedata/cmd/base
	$(GOBUILD) -o bin/tools/gen-datas -ldflags "-X main.Version=0.0.1"  ./vulcan/app/gamedata/cmd/data

.PHONY: mercury
# start mercury client
mercury:
	-pkill -f mercury/bin/mercury
	cd mercury \
	  && rm -rf bin/mercury && mkdir -p bin/mercury && $(GOBUILD) -ldflags "-X main.Version=0.0.1" -o bin/mercury/client ./cmd
	cd mercury \
	  && nohup bin/mercury/client -conf=configs -gamedata gen/gamedata/json > bin/mercury/debug.log &
	tail -f mercury/bin/mercury/debug.log

.PHONY: mercury-log
mercury-log:
	tail -f mercury/bin/mercury/debug.log

.PHONY: mercury-stop
mercury-stop:
	-pkill -f mercury/bin/mercury

.PHONY: gen-api-client
# generate api code
gen-api-client:
	bin/tools/gen-api-client

.PHONY: gen-api-db
# generate api code
gen-api-db:
	bin/tools/gen-api-db

.PHONY: gen-mercury
# generate mercury client code
gen-mercury:
	bin/tools/gen-mercury

.PHONY: gen-data-json
# generate data json file
gen-data-json:
	bin/tools/gen-data-json

.PHONY: gen-data-base
# generate data base code
gen-data-base:
	bin/tools/gen-data-base

.PHONY: gen-datas
# generate data code
gen-datas: 
	bin/tools/gen-datas

.PHONY: gen-all-data
# generate all data code
gen-all-data: gen-data-json gen-data-base gen-datas

.PHONY: vet
vet:
	go vet ./...

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
