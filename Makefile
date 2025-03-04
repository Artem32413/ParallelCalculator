.PHONY: all app1 app2  

all: app1 app2  

app1:  
	@echo "Запуск приложения 1"  
	@go run cmd/Orchestrator/main.go &  

app2:  
	@echo "Запуск приложения 2"  
	@go run cmd/Agent/main.go &  

# Дожидаемся завершения обеих задач  
wait:  
	@wait  

run: all wait