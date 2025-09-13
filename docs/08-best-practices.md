# 📈 Mejores Prácticas y Recomendaciones

## 🎯 Visión General

Este documento analiza el código actual del proyecto y proporciona recomendaciones específicas para mejorar la calidad, seguridad, mantenibilidad y escalabilidad de la API Andrei Mes Manur.

## 🛡️ Seguridad

### ✅ Implementaciones Actuales Correctas

1. **Hash de Contraseñas**
   ```go
   // controllers/auth.go:19
   hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
   ```
   ✅ Uso correcto de bcrypt con DefaultCost

2. **Exclusión de Contraseñas en JSON**
   ```go
   // models/user.go:21
   Password  string `json:"-" gorm:"not null"`
   ```
   ✅ Contraseñas excluidas de respuestas JSON

3. **Validación de Tokens JWT**
   ```go
   // middleware/auth.go:29-35
   token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
       return []byte(os.Getenv("JWT_SECRET")), nil
   })
   ```
   ✅ Validación correcta de tokens

### 🔄 Mejoras Recomendadas

#### 1. **Validación de Algoritmo JWT**
**Problema Actual:** No se valida el algoritmo del token
```go
// middleware/auth.go:29 - Agregar validación
token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
    // Validar algoritmo
    if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
        return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
    }
    return []byte(os.Getenv("JWT_SECRET")), nil
})
```

#### 2. **Validación de Input Más Robusta**
**Problema Actual:** Validaciones básicas en modelos
```go
// Crear middleware de validación customizada
func ValidateUserRegistration() gin.HandlerFunc {
    return gin.HandlerFunc(func(c *gin.Context) {
        var input models.UserRegistration
        if err := c.ShouldBindJSON(&input); err != nil {
            c.JSON(400, gin.H{"error": err.Error()})
            c.Abort()
            return
        }
        
        // Validaciones adicionales
        if len(input.Password) < 8 {
            c.JSON(400, gin.H{"error": "Password must be at least 8 characters"})
            c.Abort()
            return
        }
        
        // Validar complejidad de contraseña
        if !isValidPassword(input.Password) {
            c.JSON(400, gin.H{"error": "Password must contain letters, numbers and symbols"})
            c.Abort()
            return
        }
        
        c.Next()
    })
}
```

#### 3. **Rate Limiting**
**Problema Actual:** No hay protección contra ataques de fuerza bruta
```go
// Agregar middleware de rate limiting
import "github.com/gin-contrib/limiter"

func setupRateLimiting(r *gin.Engine) {
    // 5 requests por minuto para login
    loginLimiter := limiter.NewRateLimiter(
        redis.NewRedisStore(redisClient),
        limiter.Rate{Period: time.Minute, Limit: 5},
    )
    
    r.POST("/api/v1/login", loginLimiter.Middleware(), controllers.Login)
}
```

#### 4. **Logging de Seguridad**
**Problema Actual:** Falta logging de eventos de seguridad
```go
// Agregar logging de eventos de seguridad
func securityLogger() gin.HandlerFunc {
    return gin.HandlerFunc(func(c *gin.Context) {
        start := time.Now()
        c.Next()
        
        // Log failed auth attempts
        if c.Writer.Status() == 401 {
            log.Printf("SECURITY: Failed auth attempt from %s to %s", 
                c.ClientIP(), c.Request.URL.Path)
        }
        
        // Log privileged operations
        if strings.Contains(c.Request.URL.Path, "/admin/") && c.Writer.Status() < 400 {
            user, _ := c.Get("user")
            log.Printf("AUDIT: Admin operation %s by user %v", 
                c.Request.URL.Path, user)
        }
    })
}
```

## 🏗️ Arquitectura y Estructura

### ✅ Puntos Fuertes Actuales

1. **Separación de Responsabilidades**
   - Controllers, Models, Middleware bien separados
   - Patrón MVC implementado correctamente

2. **Middleware Modular**
   - Autenticación y autorización separados
   - Reutilizable y composable

### 🔄 Mejoras Arquitectónicas

#### 1. **Capa de Servicios**
**Problema Actual:** Lógica de negocio en controllers
```go
// Crear capa de servicios
// services/user_service.go
type UserService struct {
    db *gorm.DB
}

func (s *UserService) CreateUser(input models.UserRegistration) (*models.User, error) {
    // Validaciones de negocio
    if s.EmailExists(input.Email) {
        return nil, errors.New("email already exists")
    }
    
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
    if err != nil {
        return nil, err
    }
    
    user := &models.User{
        Username: input.Username,
        Email:    input.Email,
        Password: string(hashedPassword),
        Role:     input.Role,
    }
    
    if err := s.db.Create(user).Error; err != nil {
        return nil, err
    }
    
    return user, nil
}

// Controller simplificado
func Register(c *gin.Context) {
    var input models.UserRegistration
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    userService := services.NewUserService(config.DB)
    user, err := userService.CreateUser(input)
    if err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(201, gin.H{"user": user})
}
```

#### 2. **Repository Pattern**
**Problema Actual:** Acceso directo a GORM en controllers
```go
// repositories/user_repository.go
type UserRepository interface {
    Create(user *models.User) error
    FindByEmail(email string) (*models.User, error)
    FindByID(id uint) (*models.User, error)
    Update(user *models.User) error
    Delete(id uint) error
}

type userRepository struct {
    db *gorm.DB
}

func (r *userRepository) FindByEmail(email string) (*models.User, error) {
    var user models.User
    err := r.db.Where("email = ?", email).First(&user).Error
    return &user, err
}

// Uso en servicios
type UserService struct {
    repo UserRepository
}

func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
    return s.repo.FindByEmail(email)
}
```

#### 3. **Manejo de Errores Centralizado**
**Problema Actual:** Manejo de errores repetitivo
```go
// errors/handler.go
type AppError struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Err     error  `json:"-"`
}

func (e *AppError) Error() string {
    return e.Message
}

func NewBadRequestError(message string) *AppError {
    return &AppError{
        Code:    400,
        Message: message,
    }
}

func NewNotFoundError(message string) *AppError {
    return &AppError{
        Code:    404,
        Message: message,
    }
}

// Middleware de manejo de errores
func ErrorHandler() gin.HandlerFunc {
    return gin.HandlerFunc(func(c *gin.Context) {
        c.Next()
        
        if len(c.Errors) > 0 {
            err := c.Errors.Last()
            
            if appErr, ok := err.Err.(*AppError); ok {
                c.JSON(appErr.Code, gin.H{"error": appErr.Message})
                return
            }
            
            // Error genérico
            c.JSON(500, gin.H{"error": "Internal server error"})
        }
    })
}
```

## 💾 Base de Datos

### ✅ Implementaciones Correctas

1. **Migración Automática**
   ```go
   // config/database.go:25-31
   err = database.AutoMigrate(&models.User{}, &models.Post{}, &models.Report{}, &models.Reward{})
   ```

2. **Soft Delete**
   ```go
   // models/user.go:25
   DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
   ```

### 🔄 Mejoras de Base de Datos

#### 1. **Pool de Conexiones**
**Problema Actual:** Configuración básica de conexión
```go
// config/database.go - Mejorar configuración
func ConnectDatabase() {
    // ... DSN configuration
    
    database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
    })
    
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    
    // Configurar pool de conexiones
    sqlDB, err := database.DB()
    if err != nil {
        log.Fatal("Failed to get database instance:", err)
    }
    
    sqlDB.SetMaxIdleConns(10)
    sqlDB.SetMaxOpenConns(100)
    sqlDB.SetConnMaxLifetime(time.Hour)
    
    // ... resto de la configuración
}
```

#### 2. **Transacciones**
**Problema Actual:** Operaciones complejas sin transacciones
```go
// services/demon_service.go
func (s *DemonService) RegisterVictimWithReport(victimData models.UserRegistration, reportData models.ReportCreate) error {
    return s.db.Transaction(func(tx *gorm.DB) error {
        // Crear víctima
        victim := &models.User{
            Username: victimData.Username,
            Email:    victimData.Email,
            Password: hashPassword(victimData.Password),
            Role:     models.RoleNetworkAdmin,
        }
        
        if err := tx.Create(victim).Error; err != nil {
            return err
        }
        
        // Crear reporte inicial
        report := &models.Report{
            DemonID:     s.currentUserID,
            VictimID:    victim.ID,
            Title:       reportData.Title,
            Description: reportData.Description,
            Status:      "pending",
        }
        
        if err := tx.Create(report).Error; err != nil {
            return err
        }
        
        return nil
    })
}
```

#### 3. **Índices Optimizados**
**Problema Actual:** Índices básicos automáticos
```go
// models/user.go - Agregar índices específicos
type User struct {
    ID        uint           `json:"id" gorm:"primaryKey"`
    Username  string         `json:"username" gorm:"unique;not null;index"`
    Email     string         `json:"email" gorm:"unique;not null;index"`
    Role      UserRole       `json:"role" gorm:"not null;index"` // Para consultas por rol
    CreatedAt time.Time      `json:"created_at" gorm:"index"`   // Para ordenamiento
    // ...
}

// models/report.go - Índices compuestos
type Report struct {
    // ...
    DemonID     uint   `gorm:"index:idx_demon_status,priority:1"`
    Status      string `gorm:"index:idx_demon_status,priority:2"` // Índice compuesto
    // ...
}
```

## 🚀 Performance

### 🔄 Optimizaciones Recomendadas

#### 1. **Caching con Redis**
**Problema Actual:** Consultas repetitivas sin cache
```go
// services/cache_service.go
type CacheService struct {
    client *redis.Client
}

func (c *CacheService) GetDemonStats(demonID uint) (*models.DemonStats, error) {
    key := fmt.Sprintf("demon_stats:%d", demonID)
    
    // Intentar obtener del cache
    cached, err := c.client.Get(key).Result()
    if err == nil {
        var stats models.DemonStats
        json.Unmarshal([]byte(cached), &stats)
        return &stats, nil
    }
    
    // Si no está en cache, calcular y guardar
    stats := c.calculateDemonStats(demonID)
    
    // Cache por 5 minutos
    data, _ := json.Marshal(stats)
    c.client.Set(key, data, 5*time.Minute)
    
    return stats, nil
}
```

#### 2. **Paginación**
**Problema Actual:** Consultas sin límite
```go
// controllers/andrei.go - Agregar paginación
func GetAllUsers(c *gin.Context) {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
    
    if page < 1 {
        page = 1
    }
    if limit < 1 || limit > 100 {
        limit = 10
    }
    
    offset := (page - 1) * limit
    
    var users []models.User
    var total int64
    
    config.DB.Model(&models.User{}).Count(&total)
    
    if err := config.DB.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
        c.JSON(500, gin.H{"error": "Failed to fetch users"})
        return
    }
    
    c.JSON(200, gin.H{
        "users": users,
        "pagination": gin.H{
            "page":  page,
            "limit": limit,
            "total": total,
            "pages": (total + int64(limit) - 1) / int64(limit),
        },
    })
}
```

#### 3. **Lazy Loading**
**Problema Actual:** Carga automática de relaciones
```go
// Cargar relaciones solo cuando se necesiten
func GetUserWithPosts(c *gin.Context) {
    userID := c.Param("id")
    
    var user models.User
    if err := config.DB.First(&user, userID).Error; err != nil {
        c.JSON(404, gin.H{"error": "User not found"})
        return
    }
    
    // Cargar posts solo si se solicitan
    if c.Query("include_posts") == "true" {
        config.DB.Preload("Posts").First(&user, userID)
    }
    
    c.JSON(200, gin.H{"user": user})
}
```

## 📊 Monitoring y Observabilidad

### 🔄 Implementaciones Recomendadas

#### 1. **Métricas de Aplicación**
```go
// metrics/metrics.go
import "github.com/prometheus/client_golang/prometheus"

var (
    requestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "http_request_duration_seconds",
            Help: "Duration of HTTP requests",
        },
        []string{"method", "endpoint", "status"},
    )
    
    activeUsers = prometheus.NewGaugeVec(
        prometheus.GaugeOpts{
            Name: "active_users_total",
            Help: "Number of active users by role",
        },
        []string{"role"},
    )
)

// Middleware de métricas
func MetricsMiddleware() gin.HandlerFunc {
    return gin.HandlerFunc(func(c *gin.Context) {
        start := time.Now()
        
        c.Next()
        
        duration := time.Since(start).Seconds()
        status := strconv.Itoa(c.Writer.Status())
        
        requestDuration.WithLabelValues(
            c.Request.Method,
            c.FullPath(),
            status,
        ).Observe(duration)
    })
}
```

#### 2. **Health Checks**
```go
// controllers/health.go
func HealthCheck(c *gin.Context) {
    health := gin.H{
        "status":    "ok",
        "timestamp": time.Now().ISO8601(),
        "version":   os.Getenv("APP_VERSION"),
        "checks": gin.H{
            "database": checkDatabase(),
            "redis":    checkRedis(),
        },
    }
    
    allHealthy := true
    for _, check := range health["checks"].(gin.H) {
        if check != "ok" {
            allHealthy = false
            break
        }
    }
    
    status := 200
    if !allHealthy {
        status = 503
        health["status"] = "degraded"
    }
    
    c.JSON(status, health)
}
```

## 🧪 Testing

### 🔄 Framework de Testing Recomendado

#### 1. **Tests Unitarios**
```go
// controllers/auth_test.go
func TestLogin(t *testing.T) {
    // Setup
    gin.SetMode(gin.TestMode)
    db := setupTestDB()
    defer teardownTestDB(db)
    
    router := gin.New()
    router.POST("/login", controllers.Login)
    
    // Crear usuario de prueba
    user := createTestUser(db, "test@example.com", "password123", models.RoleDemon)
    
    tests := []struct {
        name           string
        requestBody    string
        expectedStatus int
        expectedFields []string
    }{
        {
            name:           "Valid login",
            requestBody:    `{"email":"test@example.com","password":"password123"}`,
            expectedStatus: 200,
            expectedFields: []string{"token", "user"},
        },
        {
            name:           "Invalid password",
            requestBody:    `{"email":"test@example.com","password":"wrong"}`,
            expectedStatus: 401,
            expectedFields: []string{"error"},
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            req := httptest.NewRequest("POST", "/login", strings.NewReader(tt.requestBody))
            req.Header.Set("Content-Type", "application/json")
            
            w := httptest.NewRecorder()
            router.ServeHTTP(w, req)
            
            assert.Equal(t, tt.expectedStatus, w.Code)
            
            var response map[string]interface{}
            json.Unmarshal(w.Body.Bytes(), &response)
            
            for _, field := range tt.expectedFields {
                assert.Contains(t, response, field)
            }
        })
    }
}
```

#### 2. **Tests de Integración**
```go
// tests/integration/api_test.go
func TestDemonWorkflow(t *testing.T) {
    // Setup completo del ambiente
    app := setupTestApp()
    defer app.Cleanup()
    
    // 1. Login como demonio
    token := app.LoginAs("demon@evil.com", "password123")
    
    // 2. Registrar víctima
    victim := app.POST("/api/v1/demons/victims", token, gin.H{
        "username": "testvictim",
        "email":    "victim@company.com",
        "password": "victim123",
        "role":     "network_admin",
    })
    assert.Equal(t, 201, victim.Status)
    
    // 3. Crear reporte
    report := app.POST("/api/v1/demons/reports", token, gin.H{
        "victim_id":   victim.JSON()["victim"].(map[string]interface{})["id"],
        "title":       "Test Report",
        "description": "Test description",
    })
    assert.Equal(t, 201, report.Status)
    
    // 4. Verificar estadísticas
    stats := app.GET("/api/v1/demons/stats", token)
    assert.Equal(t, 200, stats.Status)
    assert.True(t, stats.JSON()["stats"].(map[string]interface{})["victims_count"].(float64) > 0)
}
```

## 📝 Documentación

### 🔄 Mejoras de Documentación

#### 1. **Swagger/OpenAPI**
```go
// main.go - Agregar Swagger
import "github.com/swaggo/gin-swagger"
import "github.com/swaggo/files"

//go:generate swag init

// @title Andrei Mes Manur API
// @version 1.0
// @description API for the digital empire of Andrei Mes Manur
// @host localhost:8085
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
    r := gin.Default()
    
    // Swagger endpoint
    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
    
    // ... resto de la configuración
}

// Controllers con anotaciones Swagger
// @Summary Login user
// @Description Authenticate user and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param user body models.UserLogin true "User credentials"
// @Success 200 {object} map[string]interface{} "token and user data"
// @Failure 401 {object} map[string]interface{} "invalid credentials"
// @Router /login [post]
func Login(c *gin.Context) {
    // ... implementación
}
```

## 🔧 Configuración y Deploy

### 🔄 Mejoras de Configuración

#### 1. **Configuración por Ambiente**
```go
// config/config.go
type Config struct {
    Database DatabaseConfig `mapstructure:"database"`
    JWT      JWTConfig      `mapstructure:"jwt"`
    Redis    RedisConfig    `mapstructure:"redis"`
    Server   ServerConfig   `mapstructure:"server"`
}

type DatabaseConfig struct {
    Host     string `mapstructure:"host"`
    Port     int    `mapstructure:"port"`
    User     string `mapstructure:"user"`
    Password string `mapstructure:"password"`
    DBName   string `mapstructure:"dbname"`
    SSLMode  string `mapstructure:"sslmode"`
}

func LoadConfig(path string) (*Config, error) {
    viper.AddConfigPath(path)
    viper.SetConfigName("app")
    viper.SetConfigType("yaml")
    
    viper.AutomaticEnv()
    
    err := viper.ReadInConfig()
    if err != nil {
        return nil, err
    }
    
    var config Config
    err = viper.Unmarshal(&config)
    return &config, err
}
```

#### 2. **Docker Compose para Desarrollo**
```yaml
# docker-compose.yml
version: '3.8'
services:
  postgres:
    image: postgres:15
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: mydb
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
  
  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
  
  api:
    build: .
    ports:
      - "8085:8085"
    depends_on:
      - postgres
      - redis
    environment:
      - DB_HOST=postgres
      - REDIS_HOST=redis
    volumes:
      - .:/app

volumes:
  postgres_data:
```

#### 3. **Dockerfile Optimizado**
```dockerfile
# Dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/.env .

EXPOSE 8085
CMD ["./main"]
```

---

## 📋 Resumen de Prioridades

### 🔴 Críticas (Implementar Inmediatamente)
1. **Validación de algoritmo JWT** - Seguridad crítica
2. **Rate limiting en endpoints de autenticación** - Prevenir ataques
3. **Logging de eventos de seguridad** - Auditabilidad
4. **Manejo de errores centralizado** - Calidad del código

### 🟡 Importantes (Próximas 2-4 semanas)
1. **Capa de servicios** - Mejor arquitectura
2. **Repository pattern** - Testabilidad
3. **Paginación** - Performance
4. **Tests unitarios** - Calidad

### 🟢 Deseables (Mediano plazo)
1. **Caching con Redis** - Performance avanzada
2. **Métricas con Prometheus** - Observabilidad
3. **Swagger documentation** - Developer experience
4. **CI/CD pipeline** - DevOps

Estas mejoras transformarán el proyecto de un MVP funcional a una aplicación enterprise-ready con altos estándares de calidad, seguridad y mantenibilidad.