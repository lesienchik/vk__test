PROJECT_NAME := vktest
INNER_PORT := 9100
OUTER_PORT := 9100

.PHONY: start
start: build drun

.PHONY: build
build:
	docker build -t ${PROJECT_NAME} .

.PHONY: drun
drun:
	@echo 'starting ${PROJECT_NAME} application...'
	docker run --env-file ./local_files/.env -p ${OUTER_PORT}:${INNER_PORT} -d --restart unless-stopped --name ${PROJECT_NAME} ${PROJECT_NAME}

.PHONY: clean
clean:
	docker stop ${PROJECT_NAME}  && \
	docker rm ${PROJECT_NAME} && \
	docker rmi ${PROJECT_NAME}

# https://github.com/golangci/golangci-lint/releases/latest
.PHONY: lint
lint:
	@echo 'golangci-lint:'
	@golangci-lint run
	@echo ' ok'

.PHONY: swag
swag:
	@echo 'swag creation started'
	swag init -g cmd/vktest/vktest.go

.PHONY: test
test:
	@echo 'tests is started'
	go test ./...
