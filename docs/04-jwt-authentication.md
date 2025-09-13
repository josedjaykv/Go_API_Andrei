# 🔐 Autenticación JWT

## 📋 Visión General

El sistema implementa autenticación stateless usando **JSON Web Tokens (JWT)** para identificar y validar usuarios. Cada token contiene información del usuario encriptada y firmada que permite verificar la identidad sin necesidad de consultar la base de datos en cada request.

## 🛠️ Implementación Técnica

### Dependencias Utilizadas
```go
// go.mod:7 - JWT library
github.com/golang-jwt/jwt/v5 v5.3.0

// go.mod:9 - Password hashing
golang.org/x/crypto v0.42.0
```

### Configuración del Secret
**Archivo:** `.env:7`
```env
JWT_SECRET=sadasfasdlkjLKAJSDFHLKjhlkaj098734LKJASDH09LKAS1lmmc
```

## 🔑 Proceso de Autenticación

### 1. Login y Generación de Token
**Archivo:** `controllers/auth.go:38-67`

```go
func Login(c *gin.Context) {
    var input models.UserLogin
    
    // 1. Validar entrada
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // 2. Buscar usuario en BD
    var user models.User
    if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }

    // 3. Verificar contraseña con bcrypt
    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }

    // 4. Crear JWT token
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": user.ID,
        "email":   user.Email,
        "role":    user.Role,
        "exp":     time.Now().Add(time.Hour * 24).Unix(),
    })

    // 5. Firmar token con secret
    tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
        return
    }

    // 6. Devolver token y datos del usuario
    c.JSON(http.StatusOK, gin.H{
        "token": tokenString,
        "user": gin.H{
            "id":       user.ID,
            "username": user.Username,
            "email":    user.Email,
            "role":     user.Role,
        },
    })
}
```

### 2. Estructura del Token JWT

**Claims incluidos en el token:**
```json
{
  "user_id": 123,
  "email": "user@example.com", 
  "role": "demon",
  "exp": 1726246800
}
```

**Formato del Token:**
```
Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InNoYWRvd0BldmlsLmNvbSIsImV4cCI6MTcyNjI0NjgwMCwicm9sZSI6ImRlbW9uIiwidXNlcl9pZCI6Mn0.example-signature
```

## 🛡️ Middleware de Autenticación

### AuthRequired Middleware
**Archivo:** `middleware/auth.go:15-57`

```go
func AuthRequired() gin.HandlerFunc {
    return gin.HandlerFunc(func(c *gin.Context) {
        // 1. Obtener header de autorización
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
            c.Abort()
            return
        }

        // 2. Validar formato Bearer token
        bearerToken := strings.Split(authHeader, " ")
        if len(bearerToken) != 2 || strings.ToLower(bearerToken[0]) != "bearer" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
            c.Abort()
            return
        }

        // 3. Parsear y validar JWT
        tokenString := bearerToken[1]
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            return []byte(os.Getenv("JWT_SECRET")), nil
        })

        if err != nil || !token.Valid {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }

        // 4. Extraer claims
        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
            c.Abort()
            return
        }

        // 5. Verificar que el usuario aún existe
        userID := uint(claims["user_id"].(float64))
        var user models.User
        if err := config.DB.First(&user, userID).Error; err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
            c.Abort()
            return
        }

        // 6. Almacenar usuario en contexto
        c.Set("user", user)
        c.Next()
    })
}
```

### Flujo de Validación del Token

```mermaid
graph TD
    A[Request con Header] --> B{Header 'Authorization'?}
    B -->|No| C[401: Header requerido]
    B -->|Sí| D{Formato 'Bearer token'?}
    D -->|No| E[401: Formato inválido]
    D -->|Sí| F[Parsear JWT]
    F --> G{Token válido?}
    G -->|No| H[401: Token inválido]
    G -->|Sí| I[Extraer Claims]
    I --> J{Claims válidos?}
    J -->|No| K[401: Claims inválidos]
    J -->|Sí| L[Verificar usuario en BD]
    L --> M{Usuario existe?}
    M -->|No| N[401: Usuario no encontrado]
    M -->|Sí| O[Almacenar en contexto]
    O --> P[Continuar al controller]
```

## 🔒 Hash de Contraseñas

### Registro de Usuario con bcrypt
**Archivo:** `controllers/auth.go:15-37`

```go
func Register(c *gin.Context) {
    var input models.UserRegistration
    
    // Validación de entrada...
    
    // Hash de contraseña con bcrypt
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
        return
    }

    user := models.User{
        Username: input.Username,
        Email:    input.Email,
        Password: string(hashedPassword), // Almacenar hash, no contraseña original
        Role:     input.Role,
    }
    
    // Guardar en base de datos...
}
```

### Verificación de Contraseña
**Archivo:** `controllers/auth.go:53-57`

```go
// Durante el login, comparar hash almacenado con contraseña ingresada
if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
    c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
    return
}
```

## 📊 Modelos de Autenticación

### UserLogin (Entrada de Login)
**Archivo:** `models/user.go:35-38`

```go
type UserLogin struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required"`
}
```

### UserRegistration (Entrada de Registro)
**Archivo:** `models/user.go:29-34`

```go
type UserRegistration struct {
    Username string   `json:"username" binding:"required"`
    Email    string   `json:"email" binding:"required,email"`
    Password string   `json:"password" binding:"required,min=6"`
    Role     UserRole `json:"role" binding:"required"`
}
```

### User Model (Almacenamiento)
**Archivo:** `models/user.go:18-19`

```go
type User struct {
    // ... otros campos
    Email     string         `json:"email" gorm:"unique;not null"`
    Password  string         `json:"-" gorm:"not null"` // Excluido del JSON
    // ... otros campos
}
```

**Nota**: El campo `password` tiene la etiqueta `json:"-"` para excluirlo automáticamente de las respuestas JSON.

## 🔧 Configuración y Uso

### Variables de Entorno
```env
JWT_SECRET=tu_clave_secreta_muy_segura_aqui
```

### Headers de Request
```http
POST /api/v1/login
Content-Type: application/json

{
  "email": "usuario@example.com",
  "password": "password123"
}
```

### Headers de Response
```http
200 OK
Content-Type: application/json

{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "username": "usuario",
    "email": "usuario@example.com", 
    "role": "demon"
  }
}
```

### Uso del Token en Requests Posteriores
```http
GET /api/v1/demons/stats
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

## ⚡ Características de Seguridad

### 1. **Expiración de Token**
- Tokens válidos por 24 horas
- Expiración incluida en claims (`exp`)
- Verificación automática de expiración

### 2. **Validación de Usuario**
- Verificación en BD en cada request autenticado
- Usuario eliminado → token inválido
- Soporte para soft delete

### 3. **Hash Seguro de Contraseñas**
- bcrypt con DefaultCost (10 rounds)
- Salt automático por contraseña
- Comparación segura sin timing attacks

### 4. **Formato Estándar**
- Bearer token standard
- Algoritmo HS256 (HMAC SHA-256)
- Claims estructurados y consistentes

## 🚨 Manejo de Errores

### Códigos de Error HTTP
- **401 Unauthorized**: Token inválido, expirado o faltante
- **400 Bad Request**: Datos de login inválidos
- **409 Conflict**: Email/username ya existe (registro)
- **500 Internal Server Error**: Error interno del servidor

### Respuestas de Error Típicas
```json
// Token faltante
{
  "error": "Authorization header required"
}

// Credenciales incorrectas
{
  "error": "Invalid credentials"
}

// Token expirado
{
  "error": "Invalid token"
}
```

---
*Esta implementación de JWT proporciona autenticación stateless segura y escalable para toda la aplicación.*