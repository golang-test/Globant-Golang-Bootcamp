include server/configs.env
RUNPATH=server/

help:
	@echo ""
	@echo "	go-get		install & verify dependencies"
	@echo ""
	@echo "	go-test		run docker-compose & run tests"
	@echo ""
	@echo "	go-run		run docker-compose in detach mode"
	@echo ""
	@echo "	clear		downs docker-compose"
	@echo ""
	
go-get: 
	@echo "-----Installing dependencies-----"
	@cd $(RUNPATH) && go mod tidy
	@echo ""
	@echo "-----Verifying dependencies-----"
	@cd $(RUNPATH) && go mod verify

go-test:
	@echo "-----Running docker-compose (detach mode)-----"
	@docker-compose up -d
	@echo ""
	@echo "-----Running tests-----"
	@docker exec -it go_server go test -coverprofile cover.out ./...
	@echo ""
	@echo "-----Removing docker-compose-----"
	@docker-compose down

go-run:
	@echo "-----Running docker-compose (detach mode)-----"
	@docker-compose up

clear:
	@echo "-----Removing docker-compose-----"	
	@docker-compose down