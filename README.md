# SAT Bridge Pro API 🚀

![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)
![PHP Version](https://img.shields.io/badge/PHP-8.0+-777BB4?style=flat&logo=php)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-14+-316192?style=flat&logo=postgresql)

> Sistema completo de gestión y descarga automatizada de CFDIs desde el portal del SAT, con control de suscripciones y generación de reportes contables.

## 🌟 Características Principales

- ✅ **Descarga Automatizada de CFDIs** desde el portal del SAT
- ✅ **Autenticación Dual**: CIEC y FIEL (Firma Electrónica)
- ✅ **Gestión de Suscripciones** con planes Básico, Pro y Empresarial
- ✅ **Reportes Contables** detallados (mensuales, anuales, por proveedor/cliente)
- ✅ **Control de Límites** de descarga por plan
- ✅ **Almacenamiento Seguro** de certificados FIEL con encriptación AES-256-GCM
- ✅ **Resolución Automática de Captchas** con BoxFactura AI
- ✅ **API RESTful** completa con documentación Swagger
- ✅ **Arquitectura de Microservicios** (Go + PHP)
- ✅ **Docker Ready** para despliegue fácil

## 📋 Tabla de Contenidos

- [Inicio Rápido](#-inicio-rápido)
- [Características Detalladas](#-características-detalladas)
- [Arquitectura](#-arquitectura)
- [Instalación](#-instalación)
- [Documentación](#-documentación)
- [Endpoints API](#-endpoints-api)
- [Ejemplos de Uso](#-ejemplos-de-uso)
- [Despliegue](#-despliegue)
- [Testing](#-testing)
- [Contribuir](#-contribuir)

## 🚀 Inicio Rápido

### Opción 1: Con Docker (Recomendado - Un solo comando)

```bash
# 1. Clonar repositorio
git clone https://github.com/tu-usuario/sat-bridge-pro-api.git
cd sat-bridge-pro-api

# 2. Configurar variables de entorno (opcional, hay valores por defecto)
cp env.example .env
nano .env  # Editar si necesitas cambiar algo

# 3. ¡Iniciar todo!
./start.sh

# O directamente:
docker-compose up -d
```

**¡Eso es todo!** El sistema completo se levanta automáticamente:
- ✅ Backend Go en http://localhost:3000
- ✅ Swagger UI en http://localhost:3000/v1/docs
- ✅ PHP Scraper en http://localhost:8081
- ✅ PostgreSQL en localhost:5432
- ✅ Adminer en http://localhost:8080

### Opción 2: Sin Docker (Desarrollo)

```bash
# 1. Prerequisitos
# - Go 1.22.5+, PHP 8.0+, PostgreSQL 14+, Composer

# 2. Configurar
cp env.example .env
createdb fiberdb

# 3. Instalar dependencias
go mod download
cd php-sat-scraper && composer install && cd ..

# 4. Iniciar servicios (2 terminales)
# Terminal 1: Backend Go
go run src/main.go

# Terminal 2: PHP Scraper
cd php-sat-scraper && composer start
```

## ✨ Características Detalladas

### 🔐 Autenticación y Seguridad

- **JWT Authentication** con tokens de acceso y renovación
- **OAuth2 con Google**
- **Encriptación AES-256-GCM** para credenciales SAT
- **Reset de contraseña** por email
- **Verificación de email**
- **Sistema de roles** (user/admin)
- **Rate limiting** en endpoints críticos

### 📝 Gestión de FIEL y CIEC

- **Almacenamiento seguro** de certificados .cer y .key en base64
- **Encriptación** de contraseñas de FIEL
- **Validación de vigencia** de certificados
- **Múltiples FIELs** por usuario
- **CIEC encriptado** como alternativa a FIEL

### ⬇️ Descarga de CFDIs

- **Autenticación CIEC o FIEL** con el SAT
- **Filtros avanzados**:
  - Por tipo (emitidos/recibidos)
  - Por rango de fechas
  - Por RFC emisor/receptor
  - Por estado (activo/cancelado)
  - Por complemento
- **Descarga por UUID** específico
- **Consulta de metadatos** sin consumir cuota
- **Almacenamiento automático** en base de datos
- **Gestión de XMLs** en base64

### 💳 Planes de Suscripción

| Plan | Descargas/Mes | Precio | Características |
|------|---------------|--------|-----------------|
| **Básico** | 10 | $299 | Ideal para freelancers |
| **Pro** | 30 | $599 | Para pequeñas empresas |
| **Empresarial** | Ilimitadas* | $1,999 | Para grandes empresas |

*Ilimitadas con más de 5 usuarios

### 📊 Reportes Contables

- **Reporte mensual** con top proveedores y clientes
- **Reporte anual** con desglose mensual
- **Resumen de gastos** por proveedor
- **Resumen de ingresos** por cliente
- **CFDIs por rango de fechas**
- **Totales de ingresos y egresos**
- **Análisis por tipo y estado**

## 🏗️ Arquitectura

```
┌─────────────────┐
│   Frontend      │
│  (React/Vue)    │
└────────┬────────┘
         │ HTTPS/REST
         ▼
┌─────────────────┐
│   Backend Go    │◄─── JWT Auth
│   (Fiber)       │◄─── Business Logic
└────────┬────────┘
         │ HTTP
         ▼
┌─────────────────┐
│  PHP Scraper    │◄─── SAT Portal
│  (PhpCfdi)      │◄─── Captcha Resolver
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  PostgreSQL     │◄─── Data Storage
│  Database       │
└─────────────────┘
```

### Tecnologías

**Backend Go:**
- Fiber v2 - Framework web
- GORM - ORM
- JWT - Autenticación
- Viper - Configuración
- Validator - Validación

**Microservicio PHP:**
- PhpCfdi/CfdiSatScraper - Scraping del SAT
- PhpCfdi/ImageCaptchaResolver - Captchas
- Guzzle - Cliente HTTP
- Monolog - Logging

**Base de Datos:**
- PostgreSQL 14+
- Migraciones con GORM AutoMigrate

## 📦 Instalación

Ver [SETUP_GUIDE.md](./SETUP_GUIDE.md) para instrucciones detalladas de instalación.

### Configuración Rápida

1. **Variables de Entorno (.env)**

```env
APP_ENV=development
APP_PORT=3000

DB_HOST=localhost:5432
DB_NAME=sat_bridge_pro
DB_PASSWORD=tu_password

JWT_SECRET=tu_secret_minimo_32_caracteres

PHP_SAT_SCRAPER_URL=http://localhost:8080
PHP_SAT_SCRAPER_API_KEY=clave_secreta_compartida

CAPTCHA_RESOLVER_TOKEN=tu_token_boxfactura
```

2. **Base de Datos**

```bash
# Crear base de datos
createdb sat_bridge_pro

# Las migraciones se ejecutan automáticamente al iniciar
```

3. **Iniciar Servicios**

```bash
# Backend Go (con hot reload)
air

# PHP Microservice
cd php-sat-scraper && composer start
```

## 📚 Documentación

- **[API_DOCUMENTATION.md](./API_DOCUMENTATION.md)** - Documentación completa de API con ejemplos de Insomnia
- **[SETUP_GUIDE.md](./SETUP_GUIDE.md)** - Guía de instalación y configuración
- **[ENDPOINTS_REFERENCE.md](./ENDPOINTS_REFERENCE.md)** - Referencia rápida de endpoints
- **[IMPLEMENTATION_SUMMARY.md](./IMPLEMENTATION_SUMMARY.md)** - Resumen de implementación

### Swagger UI

Accede a la documentación interactiva en: `http://localhost:3000/v1/docs`

## 🔌 Endpoints API

### Autenticación
```
POST   /v1/auth/register        - Registrar usuario
POST   /v1/auth/login           - Iniciar sesión
POST   /v1/auth/refresh-tokens  - Renovar token
POST   /v1/auth/logout          - Cerrar sesión
```

### FIEL (Firma Electrónica)
```
POST   /v1/fiel       - Crear FIEL
GET    /v1/fiel       - Obtener FIEL activa
GET    /v1/fiel/all   - Listar todas
PATCH  /v1/fiel/:id   - Actualizar
DELETE /v1/fiel/:id   - Eliminar
```

### CIEC
```
POST   /v1/ciec        - Guardar CIEC
PATCH  /v1/ciec        - Actualizar
DELETE /v1/ciec        - Eliminar
GET    /v1/ciec/status - Verificar estado
```

### Descarga de CFDIs
```
POST /v1/sat/download         - Descargar CFDIs
POST /v1/sat/query-metadata   - Consultar metadatos
POST /v1/sat/download-uuid    - Descargar por UUID
GET  /v1/sat/stats            - Estadísticas
```

### Reportes
```
GET /v1/reports/monthly     - Reporte mensual
GET /v1/reports/yearly      - Reporte anual
GET /v1/reports/date-range  - Por rango de fechas
GET /v1/reports/expenses    - Resumen de gastos
GET /v1/reports/income      - Resumen de ingresos
```

Ver [ENDPOINTS_REFERENCE.md](./ENDPOINTS_REFERENCE.md) para lista completa.

## 💡 Ejemplos de Uso

### 1. Registrar Usuario

```bash
curl -X POST http://localhost:3000/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Juan Pérez",
    "email": "juan@ejemplo.com",
    "password": "Password123!",
    "role": "user"
  }'
```

### 2. Guardar FIEL

```bash
curl -X POST http://localhost:3000/v1/fiel \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "rfc": "XAXX010101000",
    "certificado_cer": "BASE64_CER_CONTENT",
    "clave_privada_key": "BASE64_KEY_CONTENT",
    "password_key": "password123",
    "nombre_certificado": "FIEL Principal",
    "fecha_vigencia_inicio": "2024-01-01T00:00:00Z",
    "fecha_vigencia_fin": "2028-01-01T00:00:00Z"
  }'
```

### 3. Descargar CFDIs

```bash
curl -X POST http://localhost:3000/v1/sat/download \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "auth_type": "fiel",
    "tipo_cfdi": "emitidos",
    "fecha_inicio": "2024-01-01",
    "fecha_fin": "2024-01-31",
    "save_to_database": true
  }'
```

### 4. Obtener Reporte Mensual

```bash
curl -X GET "http://localhost:3000/v1/reports/monthly?year=2024&month=11" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

Ver [API_DOCUMENTATION.md](./API_DOCUMENTATION.md) para más ejemplos.

## 🐳 Despliegue

### Con Docker Compose

```bash
# Producción
docker-compose -f docker-compose.prod.yml up -d

# Ver logs
docker-compose logs -f backend
docker-compose logs -f php-scraper

# Detener
docker-compose down
```

### Manual

Ver [SETUP_GUIDE.md](./SETUP_GUIDE.md#despliegue-en-producción) para instrucciones detalladas.

## 🧪 Testing

```bash
# Ejecutar todos los tests
go test ./... -v

# Tests específicos
go test ./src/service/... -v

# Con cobertura
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## 📊 Estado del Proyecto

### ✅ Completado

- [x] Sistema de autenticación completo
- [x] Gestión de FIEL y CIEC
- [x] Integración con SAT vía PHP
- [x] Control de límites por suscripción
- [x] Reportes contables detallados
- [x] Documentación completa
- [x] Docker setup
- [x] Migrations

### 🔜 Próximas Funcionalidades

- [ ] Frontend web (React/Vue)
- [ ] WebSockets para notificaciones en tiempo real
- [ ] Cola de trabajos con Redis
- [ ] Exportación de reportes a PDF/Excel
- [ ] Multi-tenancy
- [ ] Integración con sistemas contables externos

## 🤝 Contribuir

Las contribuciones son bienvenidas! Por favor:

1. Fork el proyecto
2. Crea tu Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit tus cambios (`git commit -m 'Add some AmazingFeature'`)
4. Push al Branch (`git push origin feature/AmazingFeature`)
5. Abre un Pull Request

Ver [CONTRIBUTING.md](./CONTRIBUTING.md) para más detalles.

## 📝 Licencia

Este proyecto está bajo la Licencia MIT - ver el archivo [LICENSE](LICENSE) para más detalles.

## 👥 Autores

- **Tu Nombre** - *Trabajo Inicial* - [tu-usuario](https://github.com/tu-usuario)

## 🙏 Agradecimientos

- [PhpCfdi/CfdiSatScraper](https://github.com/phpcfdi/cfdisat-scraper) - Librería para scraping del SAT
- [Fiber](https://gofiber.io/) - Framework web para Go
- [BoxFactura](https://boxfactura.com/) - Servicio de resolución de captchas

## 📞 Soporte

- 📧 Email: soporte@tudominio.com
- 📖 Docs: [Ver documentación completa](./API_DOCUMENTATION.md)
- 🐛 Issues: [GitHub Issues](https://github.com/tu-usuario/sat-bridge-pro-api/issues)

---

**Desarrollado con ❤️ para la comunidad contable mexicana**

![Made in Mexico](https://img.shields.io/badge/Made%20in-Mexico-green?style=flat&labelColor=red)
