.PHONY: cpp_container
cpp_container:
	docker run -it --rm \
		-p 8000:8000 \
		--env-file Dockerfiles/cpp/cpp.env \
		--mount type=bind,source=./scraped_questions,target=/questionsDir \
		codeforge.cpp:latest

.PHONY: it_cpp_dev
it_cpp_dev:
	cd worker; \
	make prod_bin; \
	cd ..; \
	docker build -t codeforge.cpp -f Dockerfiles/cpp/Dockerfile.cpp.prod . ; \
	docker run -it --rm \
		-p 8000:8000 \
		--env-file Dockerfiles/cpp/cpp.env \
		--mount type=bind,source=./scraped_questions,target=/questionsDir \
		codeforge.cpp:latest; 

