# 🧪 Guía de Pruebas y Validación

## 📋 Resumen Ejecutivo

Este documento contiene toda la información necesaria para probar y validar la API Andrei Mes Manur, incluyendo guías manuales, scripts automatizados y mejores prácticas para asegurar el correcto funcionamiento del sistema.

## 🎯 Objetivos de Testing

### Funcionales
- ✅ Verificar que todos los endpoints respondan correctamente
- ✅ Validar el sistema de autenticación JWT
- ✅ Confirmar el control de acceso basado en roles
- ✅ Probar la integridad de datos y relaciones

### No Funcionales
- ✅ Verificar tiempos de respuesta aceptables
- ✅ Validar el manejo de errores
- ✅ Confirmar la seguridad de endpoints
- ✅ Probar la resistencia a ataques básicos

## 🚀 Preparación del Ambiente de Pruebas

### 1. Prerequisitos
```bash
# Instalar dependencias del sistema
sudo apt-get update
sudo apt-get install curl jq git

# Docker para PostgreSQL
docker --version

# Go para ejecutar la aplicación
go version
```

### 2. Setup de Base de Datos
```bash
# Iniciar PostgreSQL con Docker
docker run --name postgres_test \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=mydb \
  -p 5432:5432 \
  -d postgres:15

# Verificar que esté corriendo
docker ps | grep postgres_test
```

### 3. Configuración de la Aplicación
```bash
# Clonar y configurar proyecto
cd /tu/directorio/proyecto

# Verificar variables de entorno
cat .env
# Debe contener:
# PORT=8085
# DB_HOST=localhost
# DB_PORT=5432
# DB_USER=postgres
# DB_PASSWORD=postgres
# DB_NAME=mydb
# JWT_SECRET=tu_secret_jwt

# Construir aplicación
go build -o andrei-api main.go

# Crear usuario Andrei
go run cmd/seed/main.go

# Poblar con datos de prueba
go run cmd/populate/main.go
```

### 4. Iniciar Servidor
```bash
# Terminal 1 - Iniciar API
./andrei-api

# Terminal 2 - Verificar que esté corriendo
curl http://localhost:8085/api/v1/resistance
```

## 🔧 Herramientas de Testing

### Script Automatizado
**Archivo:** `test_api.sh`
```bash
# Ejecutar todas las pruebas automatizadas
chmod +x test_api.sh
./test_api.sh
```

**Salida esperada:**
```
🚀 Iniciando pruebas de la API Andrei Mes Manur

📋 FASE 1: Pruebas de Autenticación
✅ PASSED: Login Andrei exitoso (Status: 200)
✅ PASSED: Login Demonio exitoso (Status: 200)
✅ PASSED: Login Network Admin exitoso (Status: 200)
✅ PASSED: Registro de nuevo demonio (Status: 201)
✅ PASSED: Registro de nuevo network admin (Status: 201)
✅ PASSED: Registro como Andrei debe fallar (Status: 400)

📋 FASE 2: Endpoints Públicos
✅ PASSED: Acceso público a página de resistencia (Status: 200)

📋 FASE 3: Funcionalidades de Andrei
✅ PASSED: Ver todos los usuarios (Status: 200)
✅ PASSED: Ver usuario específico (Status: 200)
✅ PASSED: Ver estadísticas de plataforma (Status: 200)
...

📊 RESUMEN DE PRUEBAS
Total de pruebas: 45
Exitosas: 45
Fallidas: 0

🎉 ¡Todas las pruebas pasaron! La API está funcionando correctamente.
```

### Colección de Postman
**Archivo:** `Andrei_API_Postman_Collection.json`

1. **Importar en Postman:**
   - Abrir Postman
   - File > Import > Upload Files
   - Seleccionar `Andrei_API_Postman_Collection.json`

2. **Configurar variables:**
   - Variables de colección ya configuradas
   - `base_url`: http://localhost:8085/api/v1
   - Tokens se guardan automáticamente

3. **Ejecutar colección:**
   - Click en "Run Collection"
   - Seleccionar todos los tests
   - Click "Run Andrei Mes Manur API"

## 📋 Casos de Prueba Detallados

### 1. Pruebas de Autenticación

#### Caso 1.1: Login Exitoso - Andrei
```bash
# Comando
curl -X POST http://localhost:8085/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "andrei@evil.com",
    "password": "password123"
  }'

# Respuesta esperada (200 OK)
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "username": "AndreiMesManur",
    "email": "andrei@evil.com",
    "role": "andrei"
  }
}

# Validaciones
✅ Status code: 200
✅ Token presente y válido
✅ Usuario con rol 'andrei'
✅ Password no incluido en respuesta
```

#### Caso 1.2: Login Fallido - Credenciales Incorrectas
```bash
# Comando
curl -X POST http://localhost:8085/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "andrei@evil.com",
    "password": "wrong_password"
  }'

# Respuesta esperada (401 Unauthorized)
{
  "error": "Invalid credentials"
}

# Validaciones
✅ Status code: 401
✅ Mensaje de error genérico (no revela qué está mal)
✅ No se devuelve información sensible
```

#### Caso 1.3: Registro Exitoso - Demonio
```bash
# Comando
curl -X POST http://localhost:8085/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "TestDemon",
    "email": "testdemon@evil.com",
    "password": "password123",
    "role": "demon"
  }'

# Respuesta esperada (201 Created)
{
  "message": "User registered successfully",
  "user": {
    "id": 15,
    "username": "TestDemon",
    "email": "testdemon@evil.com",
    "role": "demon"
  }
}

# Validaciones
✅ Status code: 201
✅ Usuario creado con datos correctos
✅ Password no devuelto en respuesta
✅ ID asignado automáticamente
```

### 2. Pruebas de Autorización

#### Caso 2.1: Acceso Autorizado - Andrei
```bash
# Obtener token primero
TOKEN=$(curl -s -X POST http://localhost:8085/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"email":"andrei@evil.com","password":"password123"}' | \
  jq -r '.token')

# Comando autorizado
curl -X GET http://localhost:8085/api/v1/admin/users \
  -H "Authorization: Bearer $TOKEN"

# Respuesta esperada (200 OK)
{
  "users": [
    {
      "id": 1,
      "username": "AndreiMesManur",
      "email": "andrei@evil.com",
      "role": "andrei",
      "created_at": "2025-09-13T10:00:00Z"
    },
    // ... más usuarios
  ]
}

# Validaciones
✅ Status code: 200
✅ Lista de usuarios devuelta
✅ Passwords excluidos de respuesta
✅ Todos los campos necesarios presentes
```

#### Caso 2.2: Acceso Denegado - Rol Insuficiente
```bash
# Obtener token de demonio
DEMON_TOKEN=$(curl -s -X POST http://localhost:8085/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"email":"shadow@evil.com","password":"demon123"}' | \
  jq -r '.token')

# Intentar acceso de admin con token de demonio
curl -X GET http://localhost:8085/api/v1/admin/users \
  -H "Authorization: Bearer $DEMON_TOKEN"

# Respuesta esperada (403 Forbidden)
{
  "error": "Insufficient permissions"
}

# Validaciones
✅ Status code: 403
✅ Acceso correctamente denegado
✅ Mensaje de error apropiado
✅ No se filtran datos sensibles
```

### 3. Pruebas de Funcionalidad por Rol

#### Caso 3.1: Demonio - Registrar Víctima
```bash
# Login como demonio
DEMON_TOKEN=$(curl -s -X POST http://localhost:8085/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"email":"shadow@evil.com","password":"demon123"}' | \
  jq -r '.token')

# Registrar nueva víctima
curl -X POST http://localhost:8085/api/v1/demons/victims \
  -H "Authorization: Bearer $DEMON_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "NewVictim",
    "email": "newvictim@company.com",
    "password": "victim123",
    "role": "network_admin"
  }'

# Respuesta esperada (201 Created)
{
  "message": "Victim registered successfully",
  "victim": {
    "id": 20,
    "username": "NewVictim",
    "email": "newvictim@company.com",
    "role": "network_admin"
  }
}

# Validaciones
✅ Status code: 201
✅ Víctima creada con rol network_admin
✅ Datos correctos almacenados
✅ Password hasheado en BD (no visible)
```

#### Caso 3.2: Network Admin - Post Anónimo
```bash
# Login como Network Admin
ADMIN_TOKEN=$(curl -s -X POST http://localhost:8085/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"email":"john.admin@company.com","password":"admin123"}' | \
  jq -r '.token')

# Crear post anónimo
curl -X POST http://localhost:8085/api/v1/network-admins/posts/anonymous \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Mensaje de Resistencia",
    "body": "¡No se rindan compañeros administradores!",
    "media": ""
  }'

# Respuesta esperada (201 Created)
{
  "post": {
    "id": 25,
    "title": "Mensaje de Resistencia",
    "body": "¡No se rindan compañeros administradores!",
    "media": "",
    "author_id": null,
    "anonymous": true,
    "created_at": "2025-09-13T15:30:00Z"
  }
}

# Validaciones
✅ Status code: 201
✅ Post creado como anónimo
✅ author_id es null
✅ anonymous flag es true
```

### 4. Pruebas de Integridad de Datos

#### Caso 4.1: Estadísticas de Demonio
```bash
# Login como demonio y obtener estadísticas
DEMON_TOKEN=$(curl -s -X POST http://localhost:8085/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"email":"shadow@evil.com","password":"demon123"}' | \
  jq -r '.token')

curl -X GET http://localhost:8085/api/v1/demons/stats \
  -H "Authorization: Bearer $DEMON_TOKEN"

# Respuesta esperada (200 OK)
{
  "stats": {
    "demon_id": 2,
    "victims_count": 3,
    "rewards_count": 2,
    "punishments_count": 1,
    "total_points": 150,
    "reports_count": 5
  }
}

# Validaciones
✅ Status code: 200
✅ Estadísticas calculadas correctamente
✅ Números coherentes con datos en BD
✅ Demon_id corresponde al usuario actual
```

#### Caso 4.2: Página de Resistencia Pública
```bash
# Acceso sin autenticación
curl -X GET http://localhost:8085/api/v1/resistance

# Respuesta esperada (200 OK)
{
  "posts": [
    {
      "id": 1,
      "title": "Welcome to the New Order",
      "body": "My loyal demons, the time has come...",
      "media": "",
      "anonymous": false,
      "author": "AndreiMesManur",
      "created_at": "2025-09-13T10:00:00Z"
    },
    {
      "id": 15,
      "title": "Mensaje de Resistencia",
      "body": "¡No se rindan compañeros!",
      "media": "",
      "anonymous": true,
      "author": "Anonymous",
      "created_at": "2025-09-13T15:30:00Z"
    }
    // ... más posts ordenados por fecha desc
  ]
}

# Validaciones
✅ Status code: 200
✅ Posts ordenados por fecha (más recientes primero)
✅ Posts anónimos muestran "Anonymous"
✅ Posts con autor muestran username
✅ Accesible sin autenticación
```

## 🛡️ Pruebas de Seguridad

### 1. Pruebas de Tokens JWT

#### Caso S1.1: Token Expirado
```bash
# Usar token expirado (simular)
EXPIRED_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MzAwMDAwMDB9.invalid"

curl -X GET http://localhost:8085/api/v1/admin/users \
  -H "Authorization: Bearer $EXPIRED_TOKEN"

# Respuesta esperada (401 Unauthorized)
{
  "error": "Invalid token"
}

# Validaciones
✅ Status code: 401
✅ Token expirado rechazado
✅ No se devuelve información sensible
```

#### Caso S1.2: Token Malformado
```bash
# Token inválido
curl -X GET http://localhost:8085/api/v1/admin/users \
  -H "Authorization: Bearer token_malformado"

# Respuesta esperada (401 Unauthorized)
{
  "error": "Invalid token"
}

# Validaciones
✅ Status code: 401
✅ Token malformado rechazado
✅ Error genérico (no detalles del problema)
```

#### Caso S1.3: Sin Token
```bash
# Request sin header Authorization
curl -X GET http://localhost:8085/api/v1/admin/users

# Respuesta esperada (401 Unauthorized)
{
  "error": "Authorization header required"
}

# Validaciones
✅ Status code: 401
✅ Request sin token rechazado
✅ Mensaje claro pero no sensible
```

### 2. Pruebas de SQL Injection

#### Caso S2.1: Injection en Login
```bash
# Intentar SQL injection en login
curl -X POST http://localhost:8085/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@evil.com'\'' OR 1=1--",
    "password": "anything"
  }'

# Respuesta esperada (401 Unauthorized)
{
  "error": "Invalid credentials"
}

# Validaciones
✅ Status code: 401
✅ SQL injection no funciona
✅ GORM protege automáticamente
✅ No se devuelve información de BD
```

### 3. Pruebas de Enumeración de Usuarios

#### Caso S3.1: Usuario Inexistente
```bash
# Login con usuario que no existe
curl -X POST http://localhost:8085/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "noexiste@evil.com",
    "password": "password123"
  }'

# Respuesta esperada (401 Unauthorized)
{
  "error": "Invalid credentials"
}

# Validaciones
✅ Status code: 401
✅ Mismo mensaje que password incorrecto
✅ No revela si usuario existe o no
✅ Previene enumeración de usuarios
```

## 📊 Pruebas de Performance

### 1. Tiempo de Respuesta

#### Script de Performance Básica
```bash
#!/bin/bash
# performance_test.sh

echo "🚀 Pruebas de Performance Básica"

# Función para medir tiempo de respuesta
measure_response_time() {
    local url=$1
    local headers=$2
    local method=${3:-GET}
    
    time_total=$(curl -o /dev/null -s -w '%{time_total}' -X $method "$url" $headers)
    echo "Endpoint: $url - Tiempo: ${time_total}s"
}

# Obtener token
TOKEN=$(curl -s -X POST http://localhost:8085/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"email":"andrei@evil.com","password":"password123"}' | \
  jq -r '.token')

echo "📊 Midiendo tiempos de respuesta..."

# Endpoints públicos
measure_response_time "http://localhost:8085/api/v1/resistance"

# Endpoints autenticados
measure_response_time "http://localhost:8085/api/v1/admin/users" "-H 'Authorization: Bearer $TOKEN'"
measure_response_time "http://localhost:8085/api/v1/admin/stats" "-H 'Authorization: Bearer $TOKEN'"

echo "✅ Pruebas de performance completadas"
```

#### Resultados Esperados
```
📊 Midiendo tiempos de respuesta...
Endpoint: /api/v1/resistance - Tiempo: 0.045s
Endpoint: /api/v1/admin/users - Tiempo: 0.067s
Endpoint: /api/v1/admin/stats - Tiempo: 0.089s

✅ Todos los endpoints responden en menos de 100ms
```

### 2. Prueba de Carga Básica

#### Script con Apache Bench
```bash
# Instalar apache bench si no está disponible
sudo apt-get install apache2-utils

# Prueba de carga en endpoint público
ab -n 1000 -c 10 http://localhost:8085/api/v1/resistance

# Resultados esperados
# Requests per second: > 500 req/s
# Time per request: < 20ms (mean)
# Failed requests: 0
```

## 🔍 Debugging y Troubleshooting

### Problemas Comunes

#### 1. Error de Conexión a Base de Datos
```bash
# Síntoma
curl http://localhost:8085/api/v1/resistance
# Response: Connection refused

# Diagnóstico
docker ps | grep postgres
# Si no aparece PostgreSQL

# Solución
docker run --name postgres_local -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=mydb -p 5432:5432 -d postgres:15
```

#### 2. Error 500 en Endpoints
```bash
# Verificar logs del servidor
# Los logs deben mostrar el error específico

# Errores comunes:
# - "Failed to connect to database" → Verificar .env
# - "Migration failed" → Verificar estructura de BD
# - "Invalid JWT secret" → Verificar JWT_SECRET en .env
```

#### 3. Tests Fallan Inesperadamente
```bash
# Limpiar estado de BD
docker stop postgres_local
docker rm postgres_local

# Recrear desde cero
docker run --name postgres_local -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=mydb -p 5432:5432 -d postgres:15

# Esperar que inicie completamente
sleep 10

# Re-poblar datos
go run cmd/seed/main.go
go run cmd/populate/main.go

# Reiniciar API
./andrei-api
```

## 📋 Checklist de Testing

### Pre-Deployment Checklist
- [ ] ✅ Todos los tests automatizados pasan
- [ ] ✅ Colección de Postman ejecuta sin errores
- [ ] ✅ Pruebas de seguridad básicas pasan
- [ ] ✅ Performance dentro de límites aceptables
- [ ] ✅ Endpoints públicos accesibles sin auth
- [ ] ✅ Endpoints protegidos requieren auth correcta
- [ ] ✅ Roles y permisos funcionan correctamente
- [ ] ✅ Datos de prueba consistentes
- [ ] ✅ Error handling apropiado
- [ ] ✅ Logs no contienen información sensible

### Post-Deployment Verification
- [ ] ✅ Health check endpoint responde
- [ ] ✅ Base de datos accesible
- [ ] ✅ Autenticación funcional
- [ ] ✅ Todas las funcionalidades principales operativas
- [ ] ✅ Tiempos de respuesta aceptables
- [ ] ✅ No hay errores en logs

## 📈 Métricas de Calidad

### Criterios de Aceptación
- **Cobertura de Tests**: 100% de endpoints críticos
- **Tiempo de Respuesta**: < 100ms para endpoints simples
- **Disponibilidad**: 99.9% uptime esperado
- **Seguridad**: 0 vulnerabilidades críticas
- **Funcionalidad**: 100% de casos de uso principales

### KPIs de Testing
- Tiempo de ejecución de suite completa: < 5 minutos
- Tests automatizados ejecutados: 45+
- Tests manuales ejecutados: 20+
- Casos de seguridad validados: 10+
- Casos de performance verificados: 5+

---

Esta guía proporciona una cobertura completa de testing para asegurar que la API Andrei Mes Manur funcione correctamente en todos los escenarios críticos antes de ser desplegada en producción.