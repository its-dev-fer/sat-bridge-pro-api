# PHP SAT CFDI Scraper Microservice

Microservicio PHP para descarga de CFDIs desde el portal del SAT usando PhpCfdi/CfdiSatScraper.

## Instalación

1. Instalar dependencias:
```bash
cd php-sat-scraper
composer install
```

2. Configurar variables de entorno:
```bash
cp env.example .env
# Editar .env con tus credenciales
```

3. Crear directorios necesarios:
```bash
mkdir -p storage/cfdis logs
chmod 755 storage logs
```

## Ejecución

### Servidor de desarrollo:
```bash
composer start
# O manualmente:
php -S localhost:8080 -t public
```

### Producción (con nginx/apache):
Configurar el document root a la carpeta `public/`

## Endpoints

### POST /api/download
Descarga CFDIs con filtros

**Request:**
```json
{
  "auth_type": "ciec",
  "rfc": "XAXX010101000",
  "ciec": "password123",
  "tipo_cfdi": "emitidos",
  "fecha_inicio": "2024-01-01",
  "fecha_fin": "2024-12-31",
  "rfc_emisor": "",
  "rfc_receptor": "",
  "estado": "activo",
  "complemento": ""
}
```

### POST /api/query-metadata
Consulta solo metadatos sin descargar XMLs

### POST /api/download-by-uuid
Descarga CFDI específico por UUID

### GET /health
Health check del servicio

## Autenticación

Todas las solicitudes requieren header `X-API-Key` con el valor configurado en `.env`

