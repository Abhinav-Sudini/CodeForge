.PHONY: judge
judge:
	cd judge; \
	make build; \
	cd ..; \
	docker compose up --build; 

.PHONY: judge_prod
judge_prod:
	cd judge; \
	make build; \
	cd ..; \
	docker compose -f prod.compose.yaml up --build ; 
