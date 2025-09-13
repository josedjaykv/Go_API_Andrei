# 📁 Estructura de Archivos del Proyecto

## 🌳 Árbol de Directorios

```
andrei-api/
├── 📂 cmd/                          # Comandos y utilidades
│   ├── 📂 populate/                 # Script de población de datos
│   │   └── 📄 main.go              # Poblador de base de datos
│   └── 📂 seed/                    # Script de seeding
│       └── 📄 main.go              # Creador de usuario Andrei
├── 📂 config/                       # Configuración de la aplicación  
│   └── 📄 database.go              # Configuración de base de datos
├── 📂 controllers/                  # Lógica de negocio (Controladores)
│   ├── 📄 andrei.go                # Funcionalidades de Andrei
│   ├── 📄 auth.go                  # Autenticación y registro
│   ├── 📄 demon.go                 # Funcionalidades de demonios
│   ├── 📄 network_admin.go         # Funcionalidades de Network Admins
│   └── 📄 public.go                # Endpoints públicos
├── 📂 docs/                        # Documentación del proyecto
│   ├── 📄 README.md                # Índice de documentación
│   ├── 📄 01-project-overview.md   # Visión general
│   ├── 📄 02-architecture-flow.md  # Arquitectura y flujo
│   ├── 📄 03-database-design.md    # Diseño de BD
│   ├── 📄 04-jwt-authentication.md # Autenticación JWT
│   ├── 📄 05-role-based-access.md  # Control de acceso por roles
│   └── 📄 [más documentación...]   # Otros documentos
├── 📂 middleware/                   # Middleware de la aplicación
│   ├── 📄 auth.go                  # Middleware de autenticación JWT
│   └── 📄 rbac.go                  # Middleware de control de roles
├── 📂 models/                      # Modelos de datos (Entidades)
│   ├── 📄 post.go                  # Modelo de publicaciones
│   ├── 📄 report.go                # Modelo de reportes
│   ├── 📄 reward.go                # Modelo de recompensas
│   ├── 📄 statistics.go            # Modelos de estadísticas
│   └── 📄 user.go                  # Modelo de usuarios
├── 📂 routes/                      # Definición de rutas
│   └── 📄 routes.go                # Configuración de endpoints
├── 📄 .env                         # Variables de entorno
├── 📄 go.mod                       # Dependencias de Go
├── 📄 go.sum                       # Hash de dependencias
├── 📄 main.go                      # Punto de entrada de la aplicación
├── 📄 README.md                    # Documentación principal del proyecto
├── 📄 API_TESTING_GUIDE.md         # Guía de pruebas de la API
├── 📄 test_api.sh                  # Script de pruebas automatizadas
└── 📄 Andrei_API_Postman_Collection.json  # Colección de Postman
```

## 📋 Descripción de Carpetas

### 📂 `/cmd` - Comandos y Utilidades

Contiene ejecutables auxiliares del proyecto que no son parte del servidor principal.

**Propósito:** Herramientas para desarrollo, testing y administración
**Patrón:** Cada subcarpeta representa un comando ejecutable independiente

#### 📂 `/cmd/populate`
- **Archivo:** `main.go`
- **Función:** Poblar la base de datos con datos de prueba realistas
- **Uso:** `go run cmd/populate/main.go`

#### 📂 `/cmd/seed`
- **Archivo:** `main.go` 
- **Función:** Crear el usuario inicial Andrei
- **Uso:** `go run cmd/seed/main.go`

### 📂 `/config` - Configuración

Centraliza toda la configuración de la aplicación.

**Propósito:** Gestión de configuración y conexiones externas
**Patrón:** Un archivo por servicio/configuración

#### 📄 `database.go`
- **Función:** Configuración y conexión a PostgreSQL
- **Responsabilidades:**
  - Conexión GORM
  - Migración automática
  - Variable global `DB`

### 📂 `/controllers` - Controladores

Contiene la lógica de negocio de la aplicación, organizados por dominio/rol.

**Propósito:** Procesamiento de requests HTTP y lógica de negocio
**Patrón:** Un archivo por dominio funcional

#### 📄 `auth.go`
- **Función:** Autenticación y registro de usuarios
- **Endpoints:** `/register`, `/login`
- **Responsabilidades:**
  - Validación de credenciales
  - Generación de JWT tokens
  - Hash de contraseñas

#### 📄 `andrei.go`
- **Función:** Funcionalidades exclusivas del rol Andrei
- **Endpoints:** `/admin/*`
- **Responsabilidades:**
  - CRUD de usuarios
  - Gestión de recompensas/castigos
  - Estadísticas globales
  - Administración de contenido

#### 📄 `demon.go`
- **Función:** Funcionalidades del rol Demon
- **Endpoints:** `/demons/*`
- **Responsabilidades:**
  - Registro de víctimas
  - Gestión de reportes
  - Estadísticas personales

#### 📄 `network_admin.go`
- **Función:** Funcionalidades del rol Network Admin
- **Endpoints:** `/network-admins/*`
- **Responsabilidades:**
  - Creación de posts anónimos

#### 📄 `public.go`
- **Función:** Endpoints públicos sin autenticación
- **Endpoints:** `/resistance`
- **Responsabilidades:**
  - Página de resistencia pública

### 📂 `/middleware` - Middleware

Funciones que se ejecutan entre el request y el controller.

**Propósito:** Procesamiento transversal de requests
**Patrón:** Un archivo por tipo de middleware

#### 📄 `auth.go`
- **Función:** Validación de tokens JWT
- **Middleware:** `AuthRequired()`
- **Responsabilidades:**
  - Validar formato de token
  - Verificar firma JWT
  - Cargar usuario en contexto

#### 📄 `rbac.go`
- **Función:** Control de acceso basado en roles
- **Middleware:** `RequireRole()`, `RequireAndrei()`, etc.
- **Responsabilidades:**
  - Verificar rol de usuario
  - Autorizar acceso a endpoints

### 📂 `/models` - Modelos

Definición de estructuras de datos y entidades de la base de datos.

**Propósito:** Representación de datos y validaciones
**Patrón:** Un archivo por entidad principal

#### 📄 `user.go`
- **Entidades:** `User`, `UserLogin`, `UserRegistration`
- **Relaciones:** Posts, Reports, Rewards
- **Validaciones:** Campos requeridos, formatos

#### 📄 `post.go`
- **Entidades:** `Post`, `PostCreate`
- **Relaciones:** Author (User)
- **Características:** Soporte para posts anónimos

#### 📄 `report.go`
- **Entidades:** `Report`, `ReportCreate`
- **Relaciones:** Demon (User), Victim (User)
- **Estados:** pending, in_progress, completed

#### 📄 `reward.go`
- **Entidades:** `Reward`, `RewardCreate`
- **Relaciones:** Demon (User)
- **Tipos:** reward, punishment

#### 📄 `statistics.go`
- **Entidades:** `DemonStats`, `PlatformStats`
- **Propósito:** DTOs para estadísticas calculadas

### 📂 `/routes` - Rutas

Configuración y mapeo de endpoints HTTP.

**Propósito:** Definición de API y aplicación de middleware
**Patrón:** Agrupación jerárquica de rutas

#### 📄 `routes.go`
- **Función:** Configuración completa de rutas
- **Estructura:**
  - Rutas públicas
  - Rutas autenticadas por rol
  - Aplicación de middleware

## 📄 Archivos Raíz

### 📄 `main.go`
**Punto de entrada de la aplicación**
- Carga de variables de entorno
- Inicialización de base de datos
- Configuración del servidor Gin
- Configuración de CORS
- Inicio del servidor

### 📄 `.env`
**Variables de entorno**
```env
PORT=8085
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=mydb
JWT_SECRET=clave_secreta_jwt
```

### 📄 `go.mod`
**Gestión de dependencias**
- Definición del módulo
- Lista de dependencias directas
- Versión de Go requerida

### 📄 `go.sum`
**Hash de dependencias**
- Checksums de todas las dependencias
- Garantiza integridad de dependencias

## 🧪 Archivos de Testing

### 📄 `API_TESTING_GUIDE.md`
- Guía manual de pruebas
- Comandos curl para cada endpoint
- Verificación de permisos

### 📄 `test_api.sh`
- Script automatizado de pruebas
- Validación de todos los endpoints
- Verificación de casos de seguridad

### 📄 `Andrei_API_Postman_Collection.json`
- Colección de Postman
- Tests automatizados
- Variables de entorno

## 🏗️ Patrones de Organización

### 1. **Separación por Dominio**
Cada rol (Andrei, Demon, Network Admin) tiene su propio controller con funcionalidades específicas.

### 2. **Middleware Modular**
Autenticación y autorización separados en middleware independientes y reutilizables.

### 3. **Modelos Ricos**
Estructuras con validaciones, relaciones y métodos específicos de dominio.

### 4. **Configuración Centralizada**
Toda la configuración en la carpeta `/config` con variables de entorno.

### 5. **Comandos Auxiliares**
Scripts de utilidad organizados en `/cmd` como ejecutables independientes.

## 📊 Métricas de Código

```
Total archivos .go: 15
├── Controllers: 5 archivos
├── Models: 5 archivos  
├── Middleware: 2 archivos
├── Config: 1 archivo
├── Routes: 1 archivo
└── Main: 1 archivo

Líneas de código aproximadas: ~1,500 líneas
├── Controllers: ~600 líneas (40%)
├── Models: ~300 líneas (20%)
├── Middleware: ~150 líneas (10%)
├── Routes: ~100 líneas (7%)
├── Config: ~50 líneas (3%)
└── Otros: ~300 líneas (20%)
```

## 🔄 Flujo de Dependencias

```
main.go
├── config/database.go
├── routes/routes.go
│   ├── middleware/auth.go
│   ├── middleware/rbac.go
│   └── controllers/*.go
│       └── models/*.go
└── cmd/*/main.go (independientes)
```

---
*Esta estructura sigue las mejores prácticas de Go y facilita el mantenimiento y escalabilidad del proyecto.*