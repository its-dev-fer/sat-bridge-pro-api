#!/bin/bash

# Script de inicio simplificado para SAT Bridge Pro API
# Uso: ./start.sh

set -e

echo " Iniciando SAT Bridge Pro API..."
echo ""

# Verificar que existe .env
if [ ! -f .env ]; then
    echo " Advertencia: No se encontró archivo .env"
    echo "Copiando env.example a .env..."
    cp env.example .env
    echo " Archivo .env creado. Por favor, edita las variables antes de continuar."
    echo ""
    read -p "¿Deseas continuar con la configuración por defecto? (y/n) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo " Operación cancelada. Edita el archivo .env y ejecuta este script nuevamente."
        exit 1
    fi
fi

# Verificar Docker
if ! command -v docker &> /dev/null; then
    echo " Error: Docker no está instalado"
    echo "Por favor instala Docker: https://docs.docker.com/get-docker/"
    exit 1
fi

if ! command -v docker-compose &> /dev/null && ! docker compose version &> /dev/null; then
    echo " Error: Docker Compose no está instalado"
    echo "Por favor instala Docker Compose: https://docs.docker.com/compose/install/"
    exit 1
fi

# Detener containers existentes
echo "Deteniendo containers existentes (si los hay)..."
docker-compose down 2>/dev/null || true

echo ""
echo "  Construyendo imágenes..."
docker-compose build

echo ""
echo " Levantando servicios..."
docker-compose up -d

echo ""
echo " Esperando que los servicios estén listos..."
sleep 10

# Verificar servicios
echo ""
echo " Verificando servicios..."

# Backend Go
if curl -s http://localhost:3000/v1/health-check > /dev/null 2>&1; then
    echo " Backend Go: http://localhost:3000"
    echo "    Swagger: http://localhost:3000/v1/docs"
else
    echo "  Backend Go: iniciando... (puede tardar unos segundos)"
fi

# PHP Scraper
if curl -s http://localhost:8081/health > /dev/null 2>&1; then
    echo " PHP Scraper: http://localhost:8081"
else
    echo "  PHP Scraper: iniciando..."
fi

# PostgreSQL
if docker exec postgresdb pg_isready -U postgres > /dev/null 2>&1; then
    echo " PostgreSQL: localhost:5432"
else
    echo "  PostgreSQL: iniciando..."
fi

# Adminer
echo " Adminer (DB Manager): http://localhost:8080"

echo ""
echo " Estado de los containers:"
docker-compose ps

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo  ¡SAT Bridge Pro API está corriendo!"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo " URLs importantes:"
echo "   🔗 Backend API: http://localhost:3000"
echo "   📖 Swagger UI: http://localhost:3000/v1/docs"
echo "   🐘 Adminer: http://localhost:8080"
echo "   🔧 PHP Scraper: http://localhost:8081"
echo ""
echo " Comandos útiles:"
echo "   Ver logs:     docker-compose logs -f"
echo "   Ver logs Go:  docker-compose logs -f backend"
echo "   Ver logs PHP: docker-compose logs -f php-scraper"
echo "   Detener:      docker-compose down"
echo "   Reiniciar:    docker-compose restart"
echo ""
echo " Documentación: ./API_DOCUMENTATION.md"
echo " Quick Start: ./QUICK_START.md"
echo ""

# Preguntar si desea ver los logs
read -p "¿Deseas ver los logs en tiempo real? (y/n) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    docker-compose logs -f
fi

