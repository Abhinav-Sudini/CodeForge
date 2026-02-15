# Runners
worker: 
	go run ./cmd/main/start_server.go




# Builders 
.PHONY: worker_server
worker_server: 
	go build -o ./test_environment/worker_server.bin ./cmd/main/start_server.go


.PHONY: runner
runner:
	go build -o ./test_environment/exec_worker.bin ./cmd/code_runner/code_runner.go





# Docker runtimes
.PHONY: build_docker_dev
build_docker_dev:
	docker build -f Dockerfile.dev -t worker_runtime:dev .

.PHONY: run_docker_dev
run_docker_dev:
	docker run \
		--mount type=bind,source=./test_environment,target=/codeForge \
		--network isolated_net \
		-p 8000:8000 \
		--env-file worker.env \
		worker_runtime:dev

.PHONY: it_run_docker_dev
it_run_docker_dev:
	docker run \
		--mount type=bind,source=./test_environment,target=/codeForge \
		-p 8000:8000 \
		--network isolated_net \
		-it --rm --env-file worker.env \
		worker_runtime:dev /bin/bash





# alpine builds
.PHONY: alpine_build_docker_dev
alpine_build_docker_dev:
	docker build -f Dockerfile.dev.alpine -t worker_runtime:dev.alpine .

.PHONY: alpine_run_docker_dev
alpine_run_docker_dev:
	docker run \
		--mount type=bind,source=./test_environment,target=/codeForge \
		--env-file worker.env \
		worker_runtime:dev.alpine

.PHONY: alpine_it_run_docker_dev
alpine_it_run_docker_dev:
	docker run \
		--mount type=bind,source=./test_environment,target=/codeForge \
		-it --rm --env-file worker.env \
		worker_runtime:dev.alpine /bin/sh
