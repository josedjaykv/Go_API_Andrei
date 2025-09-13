# 📚 Documentación del Proyecto - API Andrei Mes Manur

Esta carpeta contiene toda la documentación técnica del proyecto API Andrei Mes Manur, una aplicación backend con sistema de roles para el manejo de demonios, administradores de red y la suprema líder Andrei.

## 📋 Índice de Documentación

### 🏗️ Documentación de Arquitectura
- **[01-project-overview.md](./01-project-overview.md)** - Visión general del proyecto y contexto
- **[02-architecture-flow.md](./02-architecture-flow.md)** - Arquitectura y flujo de datos del sistema
- **[03-database-design.md](./03-database-design.md)** - Diseño de la base de datos y modelos

### 🔐 Documentación de Seguridad
- **[04-jwt-authentication.md](./04-jwt-authentication.md)** - Implementación de JWT y autenticación
- **[05-role-based-access.md](./05-role-based-access.md)** - Sistema de roles y autorización

### 📁 Documentación de Código
- **[06-file-structure.md](./06-file-structure.md)** - Estructura de archivos y responsabilidades
- **[07-code-documentation.md](./07-code-documentation.md)** - Documentación detallada de cada archivo

### 📈 Mejoras y Recomendaciones
- **[08-best-practices.md](./08-best-practices.md)** - Mejores prácticas y recomendaciones
- **[09-improvements.md](./09-improvements.md)** - Propuestas de mejoras futuras

### 🧪 Documentación de Pruebas
- **[10-testing-guide.md](./10-testing-guide.md)** - Guía de pruebas y validación

## 🎯 Objetivo del Proyecto

El proyecto simula un sistema de gestión para el "imperio digital" de Andrei Mes Manur, donde:
- **Andrei** (rol admin) controla todo el sistema
- **Demonios** capturan y gestionan administradores de red
- **Administradores de Red** (víctimas) tienen acceso limitado

## 🛠️ Tecnologías Utilizadas

- **Backend**: Go con framework Gin
- **Base de datos**: PostgreSQL con GORM
- **Autenticación**: JWT (JSON Web Tokens)
- **Containerización**: Docker
- **Testing**: Scripts bash y colección Postman

## 📦 Estructura del Proyecto

```
andrei-api/
├── cmd/                    # Comandos y utilidades
│   ├── seed/               # Seeder para usuario Andrei
│   └── populate/           # Poblador de datos de prueba
├── config/                 # Configuración de la aplicación
├── controllers/            # Lógica de negocio
├── middleware/            # Middleware de autenticación y autorización
├── models/                # Modelos de datos
├── routes/                # Definición de rutas
├── docs/                  # Documentación (esta carpeta)
├── main.go                # Punto de entrada de la aplicación
└── .env                   # Variables de entorno
```

## 🚀 Inicio Rápido

Para comenzar a trabajar con la documentación:

1. **Leer la visión general**: `01-project-overview.md`
2. **Entender la arquitectura**: `02-architecture-flow.md`
3. **Revisar la seguridad**: `04-jwt-authentication.md` y `05-role-based-access.md`
4. **Explorar el código**: `07-code-documentation.md`

## 📝 Notas de Versión

- **v1.0.0**: Versión inicial con funcionalidades básicas
- Sistema de roles implementado
- Autenticación JWT funcional
- CRUD completo para todas las entidades
- API REST completamente documentada

---
*Última actualización: 2025-09-13*