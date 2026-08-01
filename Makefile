include .env
export

export PROJECT_ROOT=${shell pwd}


resender-run:
	@go run ${PROJECT_ROOT}/main.go

ps:
	@docker compose ps
	
