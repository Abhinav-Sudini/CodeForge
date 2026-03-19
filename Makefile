.PHONY: judge
judge:
	cd judge; \
	make build; \
	cd ..; \
	docker compose up --build; 

.PHONY: restart_master_worker
restart_master:
	cd judge; \
	make build; \
	cd ..; \
	docker compose up judge.master --build -d; 

.PHONY: judge_prod
judge_prod:
	cd judge; \
	make build; \
	cd ..; \
	docker compose -f prod.compose.yaml up --build ; 
