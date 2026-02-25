
.PHONY: default
default:
	docker run -it --rm \
		-p 8000:8000 \
		--env-file cpp.env \
		--mount type=bind,source=./test_environment/questionsDir,target=/codeForge/questionsDir \
		codeforge.cpp:latest
