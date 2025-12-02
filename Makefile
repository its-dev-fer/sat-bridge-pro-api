.PHONY: help start stop restart logs build clean test docker-up docker-down docker-logs docker-build

# Variables
DOCKER_COMPOSE = docker-compose
GO = go

# Colors para output
GREEN  = \033[0;32m
YELLOW = \033[0;33m
RED    = \033[0;31m
NC     = \033[0m # No Color

help: ## Mostrar esta ayuda
	@echo "$(GREEN)SAT Bridge Pro API - Comandos Disponibles:$(NC)"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(YELLOW)%-20s$(NC) %s\n", $$1, $$2}'
	@echo ""

# ============================================
# Docker Commands
# ============================================

start: ## Iniciar todos los servicios con Docker
	@echo "$(GREEN) Iniciando SAT Bridge Pro API...$(NC)"
	@./start.sh

docker-up: ## Iniciar servicios Docker (detached)
	@echo "$(GREEN) Levantando servicios Docker...$(NC)"
	$(DOCKER_COMPOSE) up -d
	@echo "$(GREEN) Servicios iniciados$(NC)"
	@echo "📖 Swagger: http://localhost:3000/v1/docs"

docker-build: ## Construir imágenes Docker
	@echo "$(GREEN)  Construyendo imágenes Docker...$(NC)"
	$(DOCKER_COMPOSE) build --no-cache

docker-down: ## Detener y eliminar servicios Docker
	@echo "$(YELLOW) Deteniendo servicios Docker...$(NC)"
	$(DOCKER_COMPOSE) down

stop: docker-down ## Alias para docker-down

docker-restart: ## Reiniciar servicios Docker
	@echo "$(YELLOW) Reiniciando servicios...$(NC)"
	$(DOCKER_COMPOSE) restart

restart: docker-restart ## Alias para docker-restart

docker-logs: ## Ver logs de todos los servicios
	$(DOCKER_COMPOSE) logs -f

logs: docker-logs ## Alias para docker-logs

docker-logs-backend: ## Ver logs del backend Go
	$(DOCKER_COMPOSE) logs -f backend

docker-logs-php: ## Ver logs del PHP scraper
	$(DOCKER_COMPOSE) logs -f php-scraper

docker-logs-db: ## Ver logs de PostgreSQL
	$(DOCKER_COMPOSE) logs -f postgresdb

docker-ps: ## Ver estado de los containers
	@echo "$(GREEN)Estado de los servicios:$(NC)"
	$(DOCKER_COMPOSE) ps

status: docker-ps ## Alias para docker-ps

# ============================================
# Desarrollo Local (Sin Docker)
# ============================================

dev: ## Iniciar backend en modo desarrollo (con Air si está instalado)
	@if command -v air >/dev/null 2>&1; then \
		echo "$(GREEN) Iniciando con hot reload (Air)...$(NC)"; \
		air; \
	else \
		echo "$(YELLOW)  Air no instalado, usando go run...$(NC)"; \
		$(GO) run src/main.go; \
	fi

run: ## Ejecutar backend sin Docker
	@echo "$(GREEN) Iniciando backend Go...$(NC)"
	$(GO) run src/main.go

php: ## Iniciar PHP scraper sin Docker
	@echo "$(GREEN) Iniciando PHP scraper...$(NC)"
	cd php-sat-scraper && composer start

# ============================================
# Testing
# ============================================

test: ## Ejecutar todos los tests
	@echo "$(GREEN) Ejecutando tests...$(NC)"
	$(GO) test ./... -v

test-coverage: ## Ejecutar tests con cobertura
	@echo "$(GREEN) Generando reporte de cobertura...$(NC)"
	$(GO) test ./... -coverprofile=coverage.out
	$(GO) tool cover -html=coverage.out

# ============================================
# Database
# ============================================

db-create: ## Crear base de datos
	@echo "$(GREEN) Creando base de datos...$(NC)"
	createdb fiberdb || echo "$(YELLOW)Base de datos ya existe$(NC)"

db-drop: ## Eliminar base de datos
	@echo "$(RED)  Eliminando base de datos...$(NC)"
	dropdb fiberdb || echo "$(YELLOW)Base de datos no existe$(NC)"

db-reset: db-drop db-create ## Resetear base de datos (eliminar y crear)
	@echo "$(GREEN) Base de datos reseteada$(NC)"

db-connect: ## Conectar a PostgreSQL
	@echo "$(GREEN) Conectando a PostgreSQL...$(NC)"
	psql -d fiberdb

db-backup: ## Crear backup de la base de datos
	@echo "$(GREEN) Creando backup...$(NC)"
	@mkdir -p backups
	pg_dump -U postgres fiberdb > backups/backup_$(shell date +%Y%m%d_%H%M%S).sql
	@echo "$(GREEN) Backup creado en backups/$(NC)"

# ============================================
# Code Quality
# ============================================

lint: ## Ejecutar linter
	@echo "$(GREEN) Ejecutando linter...$(NC)"
	golangci-lint run

fmt: ## Formatear código Go
	@echo "$(GREEN) Formateando código...$(NC)"
	$(GO) fmt ./...

vet: ## Ejecutar go vet
	@echo "$(GREEN) Ejecutando go vet...$(NC)"
	$(GO) vet ./...

# ============================================
# Dependencies
# ============================================

deps: ## Instalar/actualizar dependencias Go
	@echo "$(GREEN) Instalando dependencias Go...$(NC)"
	$(GO) mod download
	$(GO) mod tidy

deps-php: ## Instalar dependencias PHP
	@echo "$(GREEN) Instalando dependencias PHP...$(NC)"
	cd php-sat-scraper && composer install

deps-all: deps deps-php ## Instalar todas las dependencias

# ============================================
# Swagger
# ============================================

swagger: ## Generar documentación Swagger
	@echo "$(GREEN) Generando documentación Swagger...$(NC)"
	swag init -g src/main.go -o src/docs

# ============================================
# Cleanup
# ============================================

clean: ## Limpiar archivos generados
	@echo "$(YELLOW) Limpiando archivos generados...$(NC)"
	rm -f coverage.out
	rm -rf bin/

clean-docker: ## Limpiar containers, imágenes y volúmenes Docker
	@echo "$(RED)  Limpiando Docker (esto eliminará datos)...$(NC)"
	$(DOCKER_COMPOSE) down -v
	docker system prune -af

# ============================================
# Production
# ============================================

build-prod: ## Construir binario para producción
	@echo "$(GREEN) Construyendo binario de producción...$(NC)"
	CGO_ENABLED=0 GOOS=linux $(GO) build -a -installsuffix cgo -ldflags="-w -s" -o bin/sat-bridge-pro src/main.go
	@echo "$(GREEN) Binario creado en bin/sat-bridge-pro$(NC)"

deploy: docker-build docker-up ## Desplegar en producción (build + up)
	@echo "$(GREEN) Desplegado en producción$(NC)"

# ============================================
# Utilities
# ============================================

health: ## Verificar health de los servicios
	@echo "$(GREEN) Verificando salud de los servicios...$(NC)"
	@echo -n "Backend Go: "
	@curl -s http://localhost:3000/v1/health-check > /dev/null && echo "$(GREEN)✅$(NC)" || echo "$(RED)❌$(NC)"
	@echo -n "PHP Scraper: "
	@curl -s http://localhost:8081/health > /dev/null && echo "$(GREEN)✅$(NC)" || echo "$(RED)❌$(NC)"
	@echo -n "PostgreSQL: "
	@docker exec postgresdb pg_isready -U postgres > /dev/null 2>&1 && echo "$(GREEN)✅$(NC)" || echo "$(RED)❌$(NC)"

urls: ## Mostrar URLs importantes
	@echo "$(GREEN) URLs del sistema:$(NC)"
	@echo "   Swagger UI:  http://localhost:3000/v1/docs"
	@echo "   Backend API: http://localhost:3000"
	@echo "   Adminer:     http://localhost:8080"
	@echo "   PHP Scraper: http://localhost:8081"

setup: deps-all db-create ## Setup inicial del proyecto
	@echo "$(GREEN) Setup completado$(NC)"
	@echo "$(YELLOW)Recuerda configurar tu archivo .env$(NC)"

.DEFAULT_GOAL := help
