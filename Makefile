deploy-env:
ifndef WORKERS
	$(error WORKERS is not defined)
endif

deploy-workers: deploy-env
	WORKERS=$(WORKERS) ./scripts/workers.sh
