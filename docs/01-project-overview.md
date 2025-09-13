# 📖 Visión General del Proyecto

## 🎭 Contexto y Narrativa

El proyecto **API Andrei Mes Manur** simula un sistema de gestión para el "imperio digital" de una bruja/warlock llamada Andrei Mes Manur. En esta narrativa:

- **Andrei Mes Manur** es una bruja fracasada que odia DevOps y las prácticas modernas de software
- A pesar de su odio, necesita una aplicación web moderna para orquestar su caos
- Ha declarado la guerra contra DevOps y administradores de red
- Planea capturar y hipnotizar administradores de red para que sirvan a su imperio

## 🎯 Objetivo del Sistema

Crear una API backend escalable y completa que permita:
1. Gestión de usuarios con diferentes roles y permisos
2. Sistema de autenticación y autorización robusto
3. Comunicación separada entre frontend y backend vía API
4. Funcionalidades específicas según el rol del usuario
5. Página pública de "resistencia" para compartir contenido

## 👥 Roles del Sistema

### 👑 Andrei (Rol: `andrei`)
**Suprema Líder con Control Total**

**Capacidades:**
- Control total de la plataforma
- CRUD completo sobre todas las entidades
- Asignación de recompensas y castigos a demonios
- Visualización de estadísticas completas del sistema
- Visualización de rankings de demonios
- Creación de publicaciones como autora identificada
- Eliminación de cualquier contenido
- Gestión completa de usuarios

**Archivos relacionados:**
- `controllers/andrei.go` - Lógica específica de Andrei
- `middleware/rbac.go:25` - Middleware `RequireAndrei()`

### 👹 Demonios (Rol: `demon`)
**Ejecutores de las Órdenes de Andrei**

**Capacidades:**
- Registro de nuevas víctimas (Network Admins)
- Creación y gestión de reportes sobre sus víctimas
- Visualización de estadísticas personales
- Creación de publicaciones para la página de resistencia
- Gestión de sus propios reportes y víctimas

**Archivos relacionados:**
- `controllers/demon.go` - Lógica específica de demonios
- `middleware/rbac.go:29` - Middleware `RequireDemon()`

### 👨‍💻 Administradores de Red (Rol: `network_admin`)
**Víctimas del Sistema**

**Capacidades:**
- Acceso a la página pública de "Resistencia"
- Creación de publicaciones anónimas
- Acceso limitado como víctimas del sistema

**Archivos relacionados:**
- `controllers/network_admin.go` - Lógica específica de Network Admins
- `middleware/rbac.go:33` - Middleware `RequireNetworkAdmin()`

## 🏗️ Arquitectura Técnica

### Stack Tecnológico

**Backend:**
- **Go 1.25+** - Lenguaje de programación
- **Gin Framework** - Web framework para Go
- **GORM** - ORM para manejo de base de datos
- **PostgreSQL** - Base de datos relacional
- **JWT** - Autenticación con tokens
- **bcrypt** - Hash de contraseñas

**Referencias de archivos:**
- `go.mod:1-45` - Dependencias del proyecto
- `main.go:1-45` - Configuración principal

### Patrones de Diseño

**MVC (Model-View-Controller):**
- **Models** (`models/`) - Representación de datos
- **Controllers** (`controllers/`) - Lógica de negocio
- **Routes** (`routes/`) - Enrutamiento (actúa como Vista en API REST)

**Middleware Pattern:**
- `middleware/auth.go` - Autenticación JWT
- `middleware/rbac.go` - Autorización basada en roles

## 🌐 Funcionalidades Principales

### 🔐 Sistema de Autenticación
- Registro de usuarios (solo demonios y network admins)
- Login con email/password
- Tokens JWT con expiración de 24 horas
- Hash seguro de contraseñas con bcrypt

### 🛡️ Sistema de Autorización
- Middleware de autenticación requerida
- Middleware de autorización por roles
- Protección de endpoints por permisos específicos
- Validación de tokens en cada request

### 📊 Gestión de Datos
- **Usuarios**: Registro, autenticación, CRUD
- **Posts**: Publicaciones con autor o anónimas
- **Reportes**: Informes de demonios sobre víctimas
- **Recompensas**: Sistema de puntos para demonios
- **Estadísticas**: Métricas del sistema y usuarios

### 🌍 Endpoints Públicos
- Página de resistencia accesible sin autenticación
- Visualización de todos los posts (con autor o anónimo)

## 📈 Métricas y Estadísticas

El sistema recopila:
- Total de usuarios por rol
- Cantidad de posts creados
- Reportes generados por demonios
- Puntos de recompensas/castigos
- Ranking de demonios por desempeño

**Archivos relacionados:**
- `models/statistics.go` - Estructuras de estadísticas
- `controllers/andrei.go:83-123` - Lógica de estadísticas

## 🔄 Flujo de Datos Principal

1. **Autenticación**: Usuario envía credenciales → JWT token
2. **Autorización**: Middleware valida token y rol
3. **Procesamiento**: Controller ejecuta lógica de negocio
4. **Persistencia**: GORM interactúa con PostgreSQL
5. **Respuesta**: JSON con datos o error

## 🎯 Casos de Uso Principales

### Para Andrei:
1. Monitorear el estado general del "imperio"
2. Recompensar o castigar a demonios
3. Eliminar usuarios problemáticos
4. Publicar decretos imperiales

### Para Demonios:
1. Registrar nuevas víctimas capturadas
2. Reportar progreso de hipnosis
3. Ver estadísticas personales de desempeño
4. Compartir técnicas de infiltración

### Para Network Admins:
1. Acceder a mensajes de resistencia
2. Publicar mensajes anónimos de esperanza
3. Coordinar resistencia contra demonios

## 📋 Requerimientos del Negocio

- **Escalabilidad**: Soportar crecimiento del "imperio"
- **Seguridad**: Proteger información sensible
- **Separación de Responsabilidades**: Frontend y backend independientes
- **Flexibilidad**: Diferentes capacidades por rol
- **Auditabilidad**: Rastreo de acciones por usuario

---
*Referencia: Este documento describe el contexto general del proyecto implementado en la estructura de archivos del repositorio.*