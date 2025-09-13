# 🧪 Guía Completa de Pruebas - API Andrei Mes Manur

Esta guía te ayudará a probar todos los endpoints de la API y verificar que los permisos funcionan correctamente.

## 🚀 Preparación

### 1. Iniciar el sistema
```bash
# 1. Iniciar PostgreSQL
docker run --name postgres_local -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=mydb -p 5432:5432 -d postgres:15

# 2. Construir la aplicación
go build -o andrei-api main.go

# 3. Crear usuario Andrei
go run cmd/seed/main.go

# 4. Poblar la base de datos con datos de prueba
go run cmd/populate/main.go

# 5. Iniciar la API
./andrei-api
```

### 2. Credenciales de Prueba
- **Andrei (Admin)**: `andrei@evil.com` / `password123`
- **Demonio**: `shadow@evil.com` / `demon123`
- **Network Admin**: `john.admin@company.com` / `admin123`

### 3. URL Base
```
http://localhost:8080/api/v1
```

## 📋 Pruebas Paso a Paso

### FASE 1: Autenticación

#### 1.1 Login Andrei (Admin)
```bash
curl -X POST http://localhost:8085/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "andrei@evil.com",
    "password": "password123"
  }'
```
**Esperado**: Token JWT + datos del usuario
**Guarda el token**: `ANDREI_TOKEN="Bearer {token}"`

#### 1.2 Login Demonio
```bash
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "shadow@evil.com",
    "password": "demon123"
  }'
```
**Esperado**: Token JWT + datos del demonio
**Guarda el token**: `DEMON_TOKEN="Bearer {token}"`

#### 1.3 Login Network Admin
```bash
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john.admin@company.com",
    "password": "admin123"
  }'
```
**Esperado**: Token JWT + datos del admin
**Guarda el token**: `ADMIN_TOKEN="Bearer {token}"`

#### 1.4 Registro Nuevo Demonio
```bash
curl -X POST http://localhost:8080/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "TestDemon",
    "email": "testdemon@evil.com",
    "password": "password123",
    "role": "demon"
  }'
```
**Esperado**: Usuario creado exitosamente

#### 1.5 Registro Nuevo Network Admin
```bash
curl -X POST http://localhost:8080/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "TestAdmin",
    "email": "testadmin@company.com",
    "password": "password123",
    "role": "network_admin"
  }'
```
**Esperado**: Usuario creado exitosamente

#### 1.6 Intento de Registro como Andrei (Debe Fallar)
```bash
curl -X POST http://localhost:8080/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "FakeAndrei",
    "email": "fake@evil.com",
    "password": "password123",
    "role": "andrei"
  }'
```
**Esperado**: Error - Rol no permitido

### FASE 2: Endpoints Públicos

#### 2.1 Página de Resistencia (Sin Autenticación)
```bash
curl -X GET http://localhost:8080/api/v1/resistance
```
**Esperado**: Lista de todos los posts (con autor o "Anonymous")

### FASE 3: Funcionalidades de Andrei (Admin)

#### 3.1 Ver Todos los Usuarios
```bash
curl -X GET http://localhost:8080/api/v1/admin/users \
  -H "Authorization: $ANDREI_TOKEN"
```
**Esperado**: Lista completa de usuarios

#### 3.2 Ver Usuario Específico
```bash
curl -X GET http://localhost:8080/api/v1/admin/users/2 \
  -H "Authorization: $ANDREI_TOKEN"
```
**Esperado**: Detalles del usuario con ID 2

#### 3.3 Ver Estadísticas de la Plataforma
```bash
curl -X GET http://localhost:8080/api/v1/admin/stats \
  -H "Authorization: $ANDREI_TOKEN"
```
**Esperado**: Estadísticas generales (usuarios, posts, reportes)

#### 3.4 Ver Ranking de Demonios
```bash
curl -X GET http://localhost:8080/api/v1/admin/demons/ranking \
  -H "Authorization: $ANDREI_TOKEN"
```
**Esperado**: Lista de demonios con sus estadísticas

#### 3.5 Crear Recompensa para Demonio
```bash
curl -X POST http://localhost:8080/api/v1/admin/rewards \
  -H "Authorization: $ANDREI_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "demon_id": 2,
    "type": "reward",
    "title": "Excelente Trabajo",
    "description": "Ha cumplido todas sus misiones",
    "points": 100
  }'
```
**Esperado**: Recompensa creada

#### 3.6 Crear Castigo para Demonio
```bash
curl -X POST http://localhost:8080/api/v1/admin/rewards \
  -H "Authorization: $ANDREI_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "demon_id": 3,
    "type": "punishment",
    "title": "Trabajo Descuidado",
    "description": "Dejó escapar a un objetivo",
    "points": -50
  }'
```
**Esperado**: Castigo creado

#### 3.7 Ver Todos los Posts
```bash
curl -X GET http://localhost:8080/api/v1/admin/posts \
  -H "Authorization: $ANDREI_TOKEN"
```
**Esperado**: Lista completa de posts

#### 3.8 Crear Post como Andrei
```bash
curl -X POST http://localhost:8080/api/v1/admin/posts \
  -H "Authorization: $ANDREI_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Nuevo Decreto Imperial",
    "body": "A partir de ahora, todos los demonios deben reportar diariamente.",
    "media": ""
  }'
```
**Esperado**: Post creado con autor Andrei

#### 3.9 Eliminar Post
```bash
curl -X DELETE http://localhost:8080/api/v1/admin/posts/1 \
  -H "Authorization: $ANDREI_TOKEN"
```
**Esperado**: Post eliminado

#### 3.10 Eliminar Usuario
```bash
curl -X DELETE http://localhost:8080/api/v1/admin/users/10 \
  -H "Authorization: $ANDREI_TOKEN"
```
**Esperado**: Usuario eliminado

### FASE 4: Funcionalidades de Demonio

#### 4.1 Registrar Nueva Víctima
```bash
curl -X POST http://localhost:8080/api/v1/demons/victims \
  -H "Authorization: $DEMON_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "NewVictim",
    "email": "victim@company.com",
    "password": "victim123",
    "role": "network_admin"
  }'
```
**Esperado**: Nueva víctima registrada

#### 4.2 Ver Mis Víctimas
```bash
curl -X GET http://localhost:8080/api/v1/demons/victims \
  -H "Authorization: $DEMON_TOKEN"
```
**Esperado**: Lista de víctimas del demonio

#### 4.3 Crear Reporte sobre Víctima
```bash
curl -X POST http://localhost:8080/api/v1/demons/reports \
  -H "Authorization: $DEMON_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "victim_id": 9,
    "title": "Progreso de Hipnosis",
    "description": "El objetivo muestra signos de control mental exitoso"
  }'
```
**Esperado**: Reporte creado

#### 4.4 Ver Mis Reportes
```bash
curl -X GET http://localhost:8080/api/v1/demons/reports \
  -H "Authorization: $DEMON_TOKEN"
```
**Esperado**: Lista de reportes del demonio

#### 4.5 Actualizar Estado de Reporte
```bash
curl -X PUT http://localhost:8080/api/v1/demons/reports/1 \
  -H "Authorization: $DEMON_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "status": "completed"
  }'
```
**Esperado**: Estado del reporte actualizado

#### 4.6 Ver Mis Estadísticas
```bash
curl -X GET http://localhost:8080/api/v1/demons/stats \
  -H "Authorization: $DEMON_TOKEN"
```
**Esperado**: Estadísticas personales del demonio

#### 4.7 Crear Post como Demonio
```bash
curl -X POST http://localhost:8080/api/v1/demons/posts \
  -H "Authorization: $DEMON_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Técnica de Infiltración",
    "body": "He descubierto una nueva forma de acceder a los sistemas",
    "media": ""
  }'
```
**Esperado**: Post creado con autor demonio

### FASE 5: Funcionalidades de Network Admin

#### 5.1 Crear Post Anónimo
```bash
curl -X POST http://localhost:8080/api/v1/network-admins/posts/anonymous \
  -H "Authorization: $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Mensaje de Resistencia",
    "body": "¡No se rindan! Juntos podemos resistir a estos demonios digitales",
    "media": ""
  }'
```
**Esperado**: Post anónimo creado

### FASE 6: Pruebas de Seguridad y Permisos

#### 6.1 Demonio Intentando Acceder a Funciones de Andrei (Debe Fallar)
```bash
curl -X GET http://localhost:8080/api/v1/admin/users \
  -H "Authorization: $DEMON_TOKEN"
```
**Esperado**: Error 403 - Permisos insuficientes

#### 6.2 Network Admin Intentando Crear Reporte (Debe Fallar)
```bash
curl -X POST http://localhost:8080/api/v1/demons/reports \
  -H "Authorization: $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "victim_id": 1,
    "title": "Test",
    "description": "Test"
  }'
```
**Esperado**: Error 403 - Permisos insuficientes

#### 6.3 Acceso Sin Token (Debe Fallar)
```bash
curl -X GET http://localhost:8080/api/v1/admin/users
```
**Esperado**: Error 401 - Token requerido

#### 6.4 Token Inválido (Debe Fallar)
```bash
curl -X GET http://localhost:8080/api/v1/admin/users \
  -H "Authorization: Bearer token_invalido"
```
**Esperado**: Error 401 - Token inválido

#### 6.5 Demonio Intentando Eliminar Usuario (Debe Fallar)
```bash
curl -X DELETE http://localhost:8080/api/v1/admin/users/5 \
  -H "Authorization: $DEMON_TOKEN"
```
**Esperado**: Error 403 - Permisos insuficientes

#### 6.6 Network Admin Intentando Ver Estadísticas (Debe Fallar)
```bash
curl -X GET http://localhost:8080/api/v1/demons/stats \
  -H "Authorization: $ADMIN_TOKEN"
```
**Esperado**: Error 403 - Permisos insuficientes

### FASE 7: Verificación Final

#### 7.1 Verificar Resistencia Actualizada
```bash
curl -X GET http://localhost:8080/api/v1/resistance
```
**Esperado**: Debe incluir todos los posts creados durante las pruebas

#### 7.2 Verificar Estadísticas Finales (Como Andrei)
```bash
curl -X GET http://localhost:8080/api/v1/admin/stats \
  -H "Authorization: $ANDREI_TOKEN"
```
**Esperado**: Números actualizados con los datos creados

## 🎯 Checklist de Verificación

### ✅ Autenticación y Autorización
- [ ] Login exitoso para cada rol
- [ ] Registro exitoso para roles permitidos
- [ ] Registro rechazado para rol Andrei
- [ ] Acceso denegado sin token
- [ ] Acceso denegado con token inválido
- [ ] Acceso denegado con permisos insuficientes

### ✅ Funcionalidades de Andrei
- [ ] Ver todos los usuarios
- [ ] Ver usuario específico
- [ ] Eliminar usuario
- [ ] Crear recompensas/castigos
- [ ] Ver estadísticas de plataforma
- [ ] Ver ranking de demonios
- [ ] CRUD completo de posts
- [ ] Crear posts como Andrei

### ✅ Funcionalidades de Demonio
- [ ] Registrar víctimas
- [ ] Ver mis víctimas
- [ ] Crear reportes
- [ ] Ver mis reportes
- [ ] Actualizar reportes
- [ ] Ver mis estadísticas
- [ ] Crear posts como demonio

### ✅ Funcionalidades de Network Admin
- [ ] Crear posts anónimos
- [ ] Acceder a página de resistencia

### ✅ Endpoint Público
- [ ] Página de resistencia accesible sin autenticación
- [ ] Muestra posts de todos los usuarios correctamente

## 🐛 Posibles Problemas

1. **Error de conexión a BD**: Verificar que PostgreSQL esté corriendo
2. **Token expirado**: Los tokens duran 24 horas, re-autenticarse si es necesario
3. **IDs incorrectos**: Verificar que los IDs existan en la base de datos
4. **Permisos de rol**: Asegurarse de usar el token correcto para cada endpoint

## 📊 Resultados Esperados

Después de completar todas las pruebas deberías tener:
- Varios usuarios de cada rol
- Múltiples reportes entre demonios y víctimas
- Recompensas y castigos asignados
- Posts de diferentes autores (Andrei, demonios, anónimos)
- Verificación completa de permisos y seguridad

¡La API está lista para soportar el imperio digital de Andrei Mes Manur! 👹