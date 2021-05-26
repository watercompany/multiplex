deploy-env:
ifndef WORKERS
	$(error WORKERS is not defined)
endif

deploy-workers: deploy-env
	WORKERS=$(WORKERS) ./scripts/workers.sh

build-env:
ifndef NAME
	$(error NAME is not defined)
endif

build: build-env
	go mod tidy
	go mod download
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build ./worker/cmd/server.go
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build ./job/dispatcher/main.go
	docker build -t multiplex/$(NAME) . 
	rm server
	rm main

deploy-docker:build
	docker-compose up

push-docker: build-env build
	docker tag multiplex/$(NAME):latest kenje4090/multiplex
	docker push kenje4090/multiplex:latest
