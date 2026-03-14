.PHONY: judge
build:
	cd judge; \
	make build; \
	cd ..; \
	docker compose up ; 
