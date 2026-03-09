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
		--mount type=bind,source=./worker/bin,target=/dev_work \
		codeforge.cpp:latest; 

.PHONY: master_it_dev
master_it_dev:
	cd master; \
	make prod_bin; \
	cd ..; \
	docker build -t codeforge.master -f Dockerfiles/master/Dockerfile . ; \
	docker run -it --rm \
		-p 7000:7000 \
		--env-file Dockerfiles/master/master.env \
		--mount type=bind,source=./master/bin,target=/dev_master \
		codeforge.master:latest; 
