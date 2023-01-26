.PHONY: default
default: list ;

.PHONY: unit-tests
unit-tests:
	go test -race -cover -v ./...

.PHONY: lint
lint:
	docker run --rm \
			-v $$PWD:/go/src/github.com/rizzza/echoserver \
			-w /go/src/github.com/rizzza/echoserver golangci/golangci-lint golangci-lint -j4 run \
			--deadline=120s \
			--max-same-issues=0

.PHONY: docker-build
docker-build:
	docker build --no-cache -t echoserver .

.PHONY: docker-compose-echoserver
docker-compose-echoserver:
	docker-compose up --build --force-recreate echoserver

.PHONY: list
list:
	@echo =============
	@echo Targets:
	@echo =============
	@$(MAKE) -pRrq -f $(lastword $(MAKEFILE_LIST)) : 2>/dev/null | awk -v RS= -F: '/^# File/,/^# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | sort | egrep -v -e '^[^[:alnum:]]' -e '^$@$$'
