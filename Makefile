.PHONY: all app1 app2  

all: app1 app2  

app1:  
	@echo "Запуск Оркестратора"  
	@go run cmd/Orchestrator/main.go &  

app2:  
	@echo "Запуск Агента"  
	@go run cmd/Agent/main.go &  

wait:  
	@wait  

run: all wait