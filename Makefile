GOCMD=GO111MODULE=on CGO_ENABLED=0 go
GOBUILD=${GOCMD} build

.PHONY: version
# Show the generated version
version:
	@find app -type d -depth 1 -print | xargs -L 1 bash -c 'cd "$$0" && pwd && $(MAKE) version'


.PHONY: wire
# Generate wire
wire:
	@find app -type d -depth 1 -print | xargs -L 1 bash -c 'echo "wire: $$0" && cd "$$0" && $(MAKE) wire'
	cd mock/cmd && wire

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
	$(GOBUILD) -o bin/tools/gen-api -ldflags "-X main.Version=0.0.1"  ./generator/app/api/cmd
	$(GOBUILD) -o bin/tools/gen-mock -ldflags "-X main.Version=0.0.1"  ./generator/app/mock/cmd
	$(GOBUILD) -o bin/tools/gen-data-json -ldflags "-X main.Version=0.0.1"  ./generator/app/gamedata/cmd/json
	$(GOBUILD) -o bin/tools/gen-data-base -ldflags "-X main.Version=0.0.1"  ./generator/app/gamedata/cmd/base
	$(GOBUILD) -o bin/tools/gen-datas -ldflags "-X main.Version=0.0.1"  ./generator/app/gamedata/cmd/data

.PHONY: mock
# start mock server
mock:
	-pkill -f mock/bin/mock
	cd mock \
	  && rm -rf bin/mock && mkdir -p bin/mock && $(GOBUILD) -ldflags "-X main.Version=0.0.1" -o bin/mock/client ./cmd
	cd mock \
	  && nohup bin/mock/client -conf=configs -gamedata gen/gamedata/json > bin/mock/debug.log &
	tail -f mock/bin/mock/debug.log

.PHONY: mock-log
mock-log:
	tail -f mock/bin/mock/debug.log

.PHONY: mock-stop
mock-stop:
	-pkill -f mock/bin/mock

.PHONY: gen-api
# generate api code
gen-api:
	bin/tools/gen-api

.PHONY: gen-mock	
# generate mock code
gen-mock:
	bin/tools/gen-mock

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
