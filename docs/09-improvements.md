# 🚀 Propuestas de Mejoras Futuras

## 📋 Visión General

Este documento presenta propuestas detalladas para evolucionar la API Andrei Mes Manur desde su estado actual hacia una aplicación enterprise-ready, organizadas por prioridad y impacto en el proyecto.

## 🎯 Roadmap de Mejoras

### 🔴 Fase 1: Seguridad y Estabilidad (1-2 semanas)

#### 1.1 Reforzamiento de Seguridad JWT
**Problema Actual:** Token JWT sin validación de algoritmo
**Archivo Afectado:** `middleware/auth.go:29`

**Implementación Propuesta:**
```go
// middleware/auth.go - Versión mejorada
func AuthRequired() gin.HandlerFunc {
    return gin.HandlerFunc(func(c *gin.Context) {
        // ... código existente hasta línea 29
        
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            // NUEVO: Validar algoritmo de firma
            if method, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
            } else if method != jwt.SigningMethodHS256 {
                return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
            }
            
            return []byte(os.Getenv("JWT_SECRET")), nil
        })
        
        // NUEVO: Validar expiración explícitamente
        if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
            if exp, ok := claims["exp"].(float64); ok {
                if time.Now().Unix() > int64(exp) {
                    c.JSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
                    c.Abort()
                    return
                }
            }
        }
        
        // ... resto del código existente
    })
}
```

#### 1.2 Rate Limiting Inteligente
**Problema Actual:** Sin protección contra ataques de fuerza bruta
**Archivos Nuevos:** `middleware/ratelimit.go`

**Implementación Propuesta:**
```go
// middleware/ratelimit.go
package middleware

import (
    "net/http"
    "sync"
    "time"
    "github.com/gin-gonic/gin"
)

type rateLimiter struct {
    requests map[string][]time.Time
    mutex    sync.RWMutex
    limit    int
    window   time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *rateLimiter {
    rl := &rateLimiter{
        requests: make(map[string][]time.Time),
        limit:    limit,
        window:   window,
    }
    
    // Cleanup goroutine
    go rl.cleanupExpiredRequests()
    return rl
}

func (rl *rateLimiter) Middleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        ip := c.ClientIP()
        
        rl.mutex.Lock()
        now := time.Now()
        
        // Limpiar requests antiguos
        if requests, exists := rl.requests[ip]; exists {
            var validRequests []time.Time
            for _, req := range requests {
                if now.Sub(req) < rl.window {
                    validRequests = append(validRequests, req)
                }
            }
            rl.requests[ip] = validRequests
        }
        
        // Verificar límite
        if len(rl.requests[ip]) >= rl.limit {
            rl.mutex.Unlock()
            c.JSON(http.StatusTooManyRequests, gin.H{
                "error": "Rate limit exceeded",
                "retry_after": rl.window.Seconds(),
            })
            c.Abort()
            return
        }
        
        // Agregar request actual
        rl.requests[ip] = append(rl.requests[ip], now)
        rl.mutex.Unlock()
        
        c.Next()
    }
}

// Uso en routes/routes.go
func SetupRoutes(r *gin.Engine) {
    api := r.Group("/api/v1")
    
    // Rate limiting específico para auth
    authLimiter := NewRateLimiter(5, time.Minute) // 5 requests per minute
    
    api.POST("/login", authLimiter.Middleware(), controllers.Login)
    api.POST("/register", authLimiter.Middleware(), controllers.Register)
    
    // ... resto de rutas
}
```

#### 1.3 Validación Robusta de Input
**Problema Actual:** Validaciones básicas insuficientes
**Archivo Nuevo:** `validators/user_validator.go`

**Implementación Propuesta:**
```go
// validators/user_validator.go
package validators

import (
    "regexp"
    "strings"
    "unicode"
)

type UserValidator struct{}

func NewUserValidator() *UserValidator {
    return &UserValidator{}
}

func (v *UserValidator) ValidatePassword(password string) []string {
    var errors []string
    
    if len(password) < 8 {
        errors = append(errors, "Password must be at least 8 characters long")
    }
    
    var hasUpper, hasLower, hasNumber, hasSpecial bool
    
    for _, char := range password {
        switch {
        case unicode.IsUpper(char):
            hasUpper = true
        case unicode.IsLower(char):
            hasLower = true
        case unicode.IsDigit(char):
            hasNumber = true
        case unicode.IsPunct(char) || unicode.IsSymbol(char):
            hasSpecial = true
        }
    }
    
    if !hasUpper {
        errors = append(errors, "Password must contain at least one uppercase letter")
    }
    if !hasLower {
        errors = append(errors, "Password must contain at least one lowercase letter")
    }
    if !hasNumber {
        errors = append(errors, "Password must contain at least one number")
    }
    if !hasSpecial {
        errors = append(errors, "Password must contain at least one special character")
    }
    
    return errors
}

func (v *UserValidator) ValidateEmail(email string) []string {
    var errors []string
    
    emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
    if !emailRegex.MatchString(email) {
        errors = append(errors, "Invalid email format")
    }
    
    return errors
}

func (v *UserValidator) ValidateUsername(username string) []string {
    var errors []string
    
    if len(username) < 3 {
        errors = append(errors, "Username must be at least 3 characters long")
    }
    
    if len(username) > 50 {
        errors = append(errors, "Username must be less than 50 characters")
    }
    
    usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
    if !usernameRegex.MatchString(username) {
        errors = append(errors, "Username can only contain letters, numbers, and underscores")
    }
    
    return errors
}

// Middleware de validación
// middleware/validation.go
func ValidateUserRegistration() gin.HandlerFunc {
    validator := validators.NewUserValidator()
    
    return func(c *gin.Context) {
        var input models.UserRegistration
        if err := c.ShouldBindJSON(&input); err != nil {
            c.JSON(400, gin.H{"error": "Invalid JSON format"})
            c.Abort()
            return
        }
        
        var allErrors []string
        
        if errors := validator.ValidateUsername(input.Username); len(errors) > 0 {
            allErrors = append(allErrors, errors...)
        }
        
        if errors := validator.ValidateEmail(input.Email); len(errors) > 0 {
            allErrors = append(allErrors, errors...)
        }
        
        if errors := validator.ValidatePassword(input.Password); len(errors) > 0 {
            allErrors = append(allErrors, errors...)
        }
        
        if len(allErrors) > 0 {
            c.JSON(400, gin.H{"errors": allErrors})
            c.Abort()
            return
        }
        
        c.Set("validated_input", input)
        c.Next()
    }
}
```

### 🟡 Fase 2: Arquitectura y Performance (3-4 semanas)

#### 2.1 Implementación de Service Layer
**Problema Actual:** Lógica de negocio en controllers
**Archivos Nuevos:** `services/user_service.go`, `services/demon_service.go`, `services/admin_service.go`

**Implementación Propuesta:**
```go
// services/user_service.go
package services

import (
    "errors"
    "andrei-api/models"
    "andrei-api/repositories"
    "golang.org/x/crypto/bcrypt"
)

type UserService struct {
    userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) *UserService {
    return &UserService{userRepo: userRepo}
}

func (s *UserService) CreateUser(input models.UserRegistration) (*models.User, error) {
    // Verificar si email ya existe
    if existingUser, _ := s.userRepo.FindByEmail(input.Email); existingUser != nil {
        return nil, errors.New("email already exists")
    }
    
    // Hash password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
    if err != nil {
        return nil, errors.New("failed to hash password")
    }
    
    user := &models.User{
        Username: input.Username,
        Email:    input.Email,
        Password: string(hashedPassword),
        Role:     input.Role,
    }
    
    if err := s.userRepo.Create(user); err != nil {
        return nil, err
    }
    
    return user, nil
}

func (s *UserService) AuthenticateUser(email, password string) (*models.User, string, error) {
    user, err := s.userRepo.FindByEmail(email)
    if err != nil {
        return nil, "", errors.New("invalid credentials")
    }
    
    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
        return nil, "", errors.New("invalid credentials")
    }
    
    // Generar token JWT
    token, err := s.generateJWTToken(user)
    if err != nil {
        return nil, "", err
    }
    
    return user, token, nil
}

// services/demon_service.go
type DemonService struct {
    userRepo   repositories.UserRepository
    reportRepo repositories.ReportRepository
    rewardRepo repositories.RewardRepository
}

func (s *DemonService) GetPersonalStats(demonID uint) (*models.DemonStats, error) {
    // Implementar lógica de estadísticas con cache
    return s.calculateAndCacheStats(demonID)
}

func (s *DemonService) RegisterVictim(demonID uint, victimData models.UserRegistration) (*models.User, error) {
    // Validar que solo se registren network_admin
    if victimData.Role != models.RoleNetworkAdmin {
        return nil, errors.New("can only register network administrators as victims")
    }
    
    // Usar UserService para crear la víctima
    userService := NewUserService(s.userRepo)
    return userService.CreateUser(victimData)
}
```

#### 2.2 Repository Pattern
**Problema Actual:** Acceso directo a GORM en toda la aplicación
**Archivos Nuevos:** `repositories/interfaces.go`, `repositories/user_repository.go`

**Implementación Propuesta:**
```go
// repositories/interfaces.go
package repositories

import "andrei-api/models"

type UserRepository interface {
    Create(user *models.User) error
    FindByID(id uint) (*models.User, error)
    FindByEmail(email string) (*models.User, error)
    Update(user *models.User) error
    Delete(id uint) error
    FindByRole(role models.UserRole, limit, offset int) ([]models.User, error)
    Count() (int64, error)
    CountByRole(role models.UserRole) (int64, error)
}

type ReportRepository interface {
    Create(report *models.Report) error
    FindByID(id uint) (*models.Report, error)
    FindByDemonID(demonID uint, limit, offset int) ([]models.Report, error)
    FindByVictimID(victimID uint, limit, offset int) ([]models.Report, error)
    Update(report *models.Report) error
    Delete(id uint) error
    CountByDemonID(demonID uint) (int64, error)
}

// repositories/user_repository.go
package repositories

import (
    "andrei-api/models"
    "gorm.io/gorm"
)

type userRepository struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
    return &userRepository{db: db}
}

func (r *userRepository) Create(user *models.User) error {
    return r.db.Create(user).Error
}

func (r *userRepository) FindByEmail(email string) (*models.User, error) {
    var user models.User
    err := r.db.Where("email = ?", email).First(&user).Error
    if err != nil {
        return nil, err
    }
    return &user, nil
}

func (r *userRepository) FindByRole(role models.UserRole, limit, offset int) ([]models.User, error) {
    var users []models.User
    err := r.db.Where("role = ?", role).Limit(limit).Offset(offset).Find(&users).Error
    return users, err
}

// Implementar resto de métodos...
```

#### 2.3 Caching con Redis
**Problema Actual:** Consultas estadísticas repetitivas sin cache
**Archivos Nuevos:** `cache/redis_client.go`, `cache/cache_service.go`

**Implementación Propuesta:**
```go
// cache/redis_client.go
package cache

import (
    "context"
    "os"
    "time"
    
    "github.com/go-redis/redis/v8"
)

type RedisClient struct {
    client *redis.Client
    ctx    context.Context
}

func NewRedisClient() *RedisClient {
    rdb := redis.NewClient(&redis.Options{
        Addr:     os.Getenv("REDIS_ADDR"),
        Password: os.Getenv("REDIS_PASSWORD"),
        DB:       0,
    })
    
    return &RedisClient{
        client: rdb,
        ctx:    context.Background(),
    }
}

// cache/cache_service.go
type CacheService struct {
    redis *RedisClient
}

func NewCacheService(redis *RedisClient) *CacheService {
    return &CacheService{redis: redis}
}

func (c *CacheService) GetDemonStats(demonID uint) (*models.DemonStats, error) {
    key := fmt.Sprintf("demon_stats:%d", demonID)
    
    // Intentar obtener del cache
    cached := c.redis.Get(key)
    if cached != nil {
        var stats models.DemonStats
        if err := json.Unmarshal(cached, &stats); err == nil {
            return &stats, nil
        }
    }
    
    // Si no está en cache, devolver nil para que el servicio calcule
    return nil, nil
}

func (c *CacheService) SetDemonStats(demonID uint, stats *models.DemonStats, ttl time.Duration) error {
    key := fmt.Sprintf("demon_stats:%d", demonID)
    data, err := json.Marshal(stats)
    if err != nil {
        return err
    }
    
    return c.redis.Set(key, data, ttl)
}

// Integración en services/demon_service.go
func (s *DemonService) GetPersonalStats(demonID uint) (*models.DemonStats, error) {
    // Intentar obtener del cache
    if cachedStats, _ := s.cache.GetDemonStats(demonID); cachedStats != nil {
        return cachedStats, nil
    }
    
    // Calcular estadísticas
    stats := &models.DemonStats{DemonID: demonID}
    
    // ... realizar cálculos complejos
    
    // Guardar en cache por 5 minutos
    s.cache.SetDemonStats(demonID, stats, 5*time.Minute)
    
    return stats, nil
}
```

### 🟢 Fase 3: Observabilidad y DevOps (5-6 semanas)

#### 3.1 Sistema de Métricas con Prometheus
**Problema Actual:** Sin visibilidad de métricas de aplicación
**Archivos Nuevos:** `metrics/prometheus.go`, `middleware/metrics.go`

**Implementación Propuesta:**
```go
// metrics/prometheus.go
package metrics

import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
    HTTPRequestsTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "endpoint", "status", "role"},
    )
    
    HTTPRequestDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "http_request_duration_seconds",
            Help:    "Duration of HTTP requests in seconds",
            Buckets: prometheus.DefBuckets,
        },
        []string{"method", "endpoint", "role"},
    )
    
    ActiveUsersByRole = promauto.NewGaugeVec(
        prometheus.GaugeOpts{
            Name: "active_users_by_role_total",
            Help: "Number of active users by role",
        },
        []string{"role"},
    )
    
    DatabaseConnections = promauto.NewGauge(
        prometheus.GaugeOpts{
            Name: "database_connections_active",
            Help: "Number of active database connections",
        },
    )
    
    JWTTokensIssued = promauto.NewCounter(
        prometheus.CounterOpts{
            Name: "jwt_tokens_issued_total",
            Help: "Total number of JWT tokens issued",
        },
    )
    
    FailedAuthAttempts = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "failed_auth_attempts_total",
            Help: "Total number of failed authentication attempts",
        },
        []string{"reason", "client_ip"},
    )
)

// middleware/metrics.go
func PrometheusMiddleware() gin.HandlerFunc {
    return gin.HandlerFunc(func(c *gin.Context) {
        start := time.Now()
        
        c.Next()
        
        duration := time.Since(start).Seconds()
        status := strconv.Itoa(c.Writer.Status())
        
        // Obtener rol del usuario si está autenticado
        role := "anonymous"
        if user, exists := c.Get("user"); exists {
            if u, ok := user.(models.User); ok {
                role = string(u.Role)
            }
        }
        
        metrics.HTTPRequestsTotal.WithLabelValues(
            c.Request.Method,
            c.FullPath(),
            status,
            role,
        ).Inc()
        
        metrics.HTTPRequestDuration.WithLabelValues(
            c.Request.Method,
            c.FullPath(),
            role,
        ).Observe(duration)
        
        // Métricas específicas de seguridad
        if status == "401" {
            reason := "invalid_token"
            if strings.Contains(c.FullPath(), "login") {
                reason = "invalid_credentials"
            }
            
            metrics.FailedAuthAttempts.WithLabelValues(
                reason,
                c.ClientIP(),
            ).Inc()
        }
    })
}
```

#### 3.2 Logging Estructurado
**Problema Actual:** Logging básico con log.Printf
**Archivos Nuevos:** `logger/structured_logger.go`

**Implementación Propuesta:**
```go
// logger/structured_logger.go
package logger

import (
    "os"
    "time"
    
    "github.com/sirupsen/logrus"
)

type StructuredLogger struct {
    logger *logrus.Logger
}

func NewStructuredLogger() *StructuredLogger {
    logger := logrus.New()
    
    // Configuración basada en ambiente
    if os.Getenv("ENV") == "production" {
        logger.SetLevel(logrus.InfoLevel)
        logger.SetFormatter(&logrus.JSONFormatter{})
    } else {
        logger.SetLevel(logrus.DebugLevel)
        logger.SetFormatter(&logrus.TextFormatter{
            FullTimestamp: true,
        })
    }
    
    return &StructuredLogger{logger: logger}
}

func (l *StructuredLogger) LogAuthAttempt(success bool, userEmail, clientIP string) {
    fields := logrus.Fields{
        "event":     "auth_attempt",
        "success":   success,
        "email":     userEmail,
        "client_ip": clientIP,
        "timestamp": time.Now().ISO8601(),
    }
    
    if success {
        l.logger.WithFields(fields).Info("Successful authentication")
    } else {
        l.logger.WithFields(fields).Warn("Failed authentication attempt")
    }
}

func (l *StructuredLogger) LogPrivilegedAction(userID uint, action, resource string) {
    l.logger.WithFields(logrus.Fields{
        "event":     "privileged_action",
        "user_id":   userID,
        "action":    action,
        "resource":  resource,
        "timestamp": time.Now().ISO8601(),
    }).Info("Privileged action performed")
}

func (l *StructuredLogger) LogAPIError(err error, endpoint, method string, statusCode int) {
    l.logger.WithFields(logrus.Fields{
        "event":       "api_error",
        "error":       err.Error(),
        "endpoint":    endpoint,
        "method":      method,
        "status_code": statusCode,
        "timestamp":   time.Now().ISO8601(),
    }).Error("API error occurred")
}

// Middleware de logging
// middleware/logging.go
func StructuredLoggingMiddleware(logger *logger.StructuredLogger) gin.HandlerFunc {
    return gin.HandlerFunc(func(c *gin.Context) {
        start := time.Now()
        
        c.Next()
        
        latency := time.Since(start)
        status := c.Writer.Status()
        
        // Log para requests importantes
        if status >= 400 || c.Request.URL.Path == "/api/v1/login" {
            logrus.WithFields(logrus.Fields{
                "method":    c.Request.Method,
                "path":      c.Request.URL.Path,
                "status":    status,
                "latency":   latency,
                "client_ip": c.ClientIP(),
                "user_agent": c.Request.UserAgent(),
            }).Info("API request processed")
        }
    })
}
```

#### 3.3 Health Checks Avanzados
**Problema Actual:** Sin endpoints de health check
**Archivos Nuevos:** `health/health_checker.go`, `controllers/health.go`

**Implementación Propuesta:**
```go
// health/health_checker.go
package health

import (
    "context"
    "database/sql"
    "time"
    
    "gorm.io/gorm"
)

type HealthChecker struct {
    db    *gorm.DB
    redis *redis.Client
}

func NewHealthChecker(db *gorm.DB, redis *redis.Client) *HealthChecker {
    return &HealthChecker{db: db, redis: redis}
}

type HealthStatus struct {
    Status    string                 `json:"status"`
    Timestamp time.Time             `json:"timestamp"`
    Version   string                 `json:"version"`
    Checks    map[string]CheckResult `json:"checks"`
}

type CheckResult struct {
    Status  string        `json:"status"`
    Message string        `json:"message,omitempty"`
    Latency time.Duration `json:"latency"`
}

func (h *HealthChecker) CheckHealth() HealthStatus {
    checks := make(map[string]CheckResult)
    
    // Check database
    checks["database"] = h.checkDatabase()
    
    // Check Redis
    checks["redis"] = h.checkRedis()
    
    // Check external dependencies
    checks["jwt_secret"] = h.checkJWTSecret()
    
    // Determine overall status
    overallStatus := "healthy"
    for _, check := range checks {
        if check.Status != "ok" {
            overallStatus = "degraded"
            break
        }
    }
    
    return HealthStatus{
        Status:    overallStatus,
        Timestamp: time.Now(),
        Version:   os.Getenv("APP_VERSION"),
        Checks:    checks,
    }
}

func (h *HealthChecker) checkDatabase() CheckResult {
    start := time.Now()
    
    sqlDB, err := h.db.DB()
    if err != nil {
        return CheckResult{
            Status:  "error",
            Message: "Failed to get database connection",
            Latency: time.Since(start),
        }
    }
    
    if err := sqlDB.Ping(); err != nil {
        return CheckResult{
            Status:  "error",
            Message: "Database ping failed",
            Latency: time.Since(start),
        }
    }
    
    return CheckResult{
        Status:  "ok",
        Latency: time.Since(start),
    }
}

// controllers/health.go
func HealthCheck(c *gin.Context) {
    healthChecker := health.NewHealthChecker(config.DB, cache.RedisClient)
    healthStatus := healthChecker.CheckHealth()
    
    statusCode := 200
    if healthStatus.Status != "healthy" {
        statusCode = 503
    }
    
    c.JSON(statusCode, healthStatus)
}

func ReadinessCheck(c *gin.Context) {
    // Check solo si la app está lista para recibir tráfico
    healthChecker := health.NewHealthChecker(config.DB, cache.RedisClient)
    dbCheck := healthChecker.checkDatabase()
    
    if dbCheck.Status != "ok" {
        c.JSON(503, gin.H{
            "status": "not_ready",
            "message": "Database not available",
        })
        return
    }
    
    c.JSON(200, gin.H{
        "status": "ready",
    })
}
```

### 🔵 Fase 4: Testing y Calidad (7-8 semanas)

#### 4.1 Test Suite Comprehensiva
**Problema Actual:** Sin tests automatizados
**Archivos Nuevos:** `tests/` directory structure

**Implementación Propuesta:**
```go
// tests/setup_test.go
package tests

import (
    "database/sql/driver"
    "regexp"
    "testing"
    
    "github.com/DATA-DOG/go-sqlmock"
    "github.com/gin-gonic/gin"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

func SetupTestDB() (*gorm.DB, sqlmock.Sqlmock, error) {
    db, mock, err := sqlmock.New()
    if err != nil {
        return nil, nil, err
    }
    
    gormDB, err := gorm.Open(postgres.New(postgres.Config{
        Conn: db,
    }), &gorm.Config{})
    
    return gormDB, mock, err
}

func SetupTestRouter() *gin.Engine {
    gin.SetMode(gin.TestMode)
    router := gin.New()
    return router
}

// tests/unit/services/user_service_test.go
func TestUserService_CreateUser(t *testing.T) {
    db, mock, err := SetupTestDB()
    assert.NoError(t, err)
    defer db.DB()
    
    userRepo := repositories.NewUserRepository(db)
    userService := services.NewUserService(userRepo)
    
    tests := []struct {
        name          string
        input         models.UserRegistration
        mockBehavior  func(mock sqlmock.Sqlmock)
        expectedError string
    }{
        {
            name: "successful user creation",
            input: models.UserRegistration{
                Username: "testuser",
                Email:    "test@example.com",
                Password: "password123",
                Role:     models.RoleDemon,
            },
            mockBehavior: func(mock sqlmock.Sqlmock) {
                mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users"`)).
                    WithArgs("test@example.com").
                    WillReturnError(gorm.ErrRecordNotFound)
                
                mock.ExpectBegin()
                mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users"`)).
                    WithArgs("testuser", "test@example.com", sqlmock.AnyArg(), models.RoleDemon, sqlmock.AnyArg(), sqlmock.AnyArg()).
                    WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
                mock.ExpectCommit()
            },
            expectedError: "",
        },
        {
            name: "email already exists",
            input: models.UserRegistration{
                Username: "testuser",
                Email:    "existing@example.com",
                Password: "password123",
                Role:     models.RoleDemon,
            },
            mockBehavior: func(mock sqlmock.Sqlmock) {
                mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users"`)).
                    WithArgs("existing@example.com").
                    WillReturnRows(sqlmock.NewRows([]string{"id", "email"}).
                        AddRow(1, "existing@example.com"))
            },
            expectedError: "email already exists",
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            tt.mockBehavior(mock)
            
            user, err := userService.CreateUser(tt.input)
            
            if tt.expectedError != "" {
                assert.Error(t, err)
                assert.Contains(t, err.Error(), tt.expectedError)
                assert.Nil(t, user)
            } else {
                assert.NoError(t, err)
                assert.NotNil(t, user)
                assert.Equal(t, tt.input.Username, user.Username)
                assert.Equal(t, tt.input.Email, user.Email)
            }
            
            assert.NoError(t, mock.ExpectationsWereMet())
        })
    }
}

// tests/integration/api_test.go
func TestDemonWorkflow(t *testing.T) {
    // Setup test database
    testDB := setupTestDatabase()
    defer cleanupTestDatabase(testDB)
    
    // Setup test server
    app := setupTestApp(testDB)
    
    // Test scenario: Complete demon workflow
    t.Run("complete demon workflow", func(t *testing.T) {
        // 1. Create demon user
        demon := createTestUser(testDB, "demon@evil.com", "password123", models.RoleDemon)
        
        // 2. Login as demon
        loginResp := app.POST("/api/v1/login", gin.H{
            "email":    "demon@evil.com",
            "password": "password123",
        })
        assert.Equal(t, 200, loginResp.Code)
        
        var loginData map[string]interface{}
        json.Unmarshal(loginResp.Body.Bytes(), &loginData)
        token := loginData["token"].(string)
        
        // 3. Register victim
        victimResp := app.POST("/api/v1/demons/victims", gin.H{
            "username": "victim1",
            "email":    "victim1@company.com",
            "password": "victim123",
            "role":     "network_admin",
        }, "Bearer "+token)
        assert.Equal(t, 201, victimResp.Code)
        
        // 4. Create report
        reportResp := app.POST("/api/v1/demons/reports", gin.H{
            "victim_id":   getVictimID(victimResp),
            "title":       "Initial Contact",
            "description": "Successfully made contact with target",
        }, "Bearer "+token)
        assert.Equal(t, 201, reportResp.Code)
        
        // 5. Check personal stats
        statsResp := app.GET("/api/v1/demons/stats", "Bearer "+token)
        assert.Equal(t, 200, statsResp.Code)
        
        var statsData map[string]interface{}
        json.Unmarshal(statsResp.Body.Bytes(), &statsData)
        stats := statsData["stats"].(map[string]interface{})
        
        assert.Equal(t, float64(1), stats["victims_count"])
        assert.Equal(t, float64(1), stats["reports_count"])
    })
}
```

#### 4.2 Configuración de CI/CD
**Problema Actual:** Sin pipeline automatizado
**Archivos Nuevos:** `.github/workflows/ci.yml`, `Dockerfile.test`

**Implementación Propuesta:**
```yaml
# .github/workflows/ci.yml
name: CI/CD Pipeline

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]

jobs:
  test:
    runs-on: ubuntu-latest
    
    services:
      postgres:
        image: postgres:15
        env:
          POSTGRES_PASSWORD: postgres
          POSTGRES_DB: testdb
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
        ports:
          - 5432:5432
      
      redis:
        image: redis:7
        options: >-
          --health-cmd "redis-cli ping"
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
        ports:
          - 6379:6379
    
    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Go
      uses: actions/setup-go@v3
      with:
        go-version: 1.21
    
    - name: Cache Go modules
      uses: actions/cache@v3
      with:
        path: ~/go/pkg/mod
        key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
        restore-keys: |
          ${{ runner.os }}-go-
    
    - name: Install dependencies
      run: go mod download
    
    - name: Run linting
      uses: golangci/golangci-lint-action@v3
      with:
        version: latest
        args: --timeout=5m
    
    - name: Run tests
      env:
        DB_HOST: localhost
        DB_PORT: 5432
        DB_USER: postgres
        DB_PASSWORD: postgres
        DB_NAME: testdb
        REDIS_ADDR: localhost:6379
        JWT_SECRET: test-secret-key
      run: |
        go test ./... -v -cover -coverprofile=coverage.out
        go tool cover -html=coverage.out -o coverage.html
    
    - name: Upload coverage reports
      uses: codecov/codecov-action@v3
      with:
        file: ./coverage.out
    
    - name: Build application
      run: go build -v ./...

  security-scan:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    - name: Run security scan
      uses: securecodewarrior/github-action-gosec@master
      with:
        args: ./...

  docker-build:
    needs: [test, security-scan]
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main'
    
    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Docker Buildx
      uses: docker/setup-buildx-action@v2
    
    - name: Login to DockerHub
      uses: docker/login-action@v2
      with:
        username: ${{ secrets.DOCKER_USERNAME }}
        password: ${{ secrets.DOCKER_PASSWORD }}
    
    - name: Build and push Docker image
      uses: docker/build-push-action@v4
      with:
        context: .
        push: true
        tags: |
          your-registry/andrei-api:latest
          your-registry/andrei-api:${{ github.sha }}
        cache-from: type=gha
        cache-to: type=gha,mode=max

  deploy:
    needs: docker-build
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main'
    
    steps:
    - name: Deploy to production
      run: |
        echo "Deployment steps would go here"
        # kubectl apply -f k8s/
        # or docker-compose up -d
        # or your deployment method
```

## 📊 Cronograma de Implementación

### Semana 1-2: Seguridad Crítica
- ✅ Validación de algoritmo JWT
- ✅ Rate limiting para endpoints de auth
- ✅ Logging de seguridad
- ✅ Validación robusta de input

### Semana 3-4: Refactoring Arquitectónico
- 🔄 Service layer implementation
- 🔄 Repository pattern
- 🔄 Manejo centralizado de errores

### Semana 5-6: Performance y Caching
- ⏳ Redis integration
- ⏳ Database connection pooling
- ⏳ Query optimization
- ⏳ Paginación

### Semana 7-8: Observabilidad
- 📊 Prometheus metrics
- 📝 Structured logging
- 🏥 Advanced health checks
- 📈 Monitoring dashboard

### Semana 9-10: Testing y CI/CD
- 🧪 Unit test suite
- 🔗 Integration tests
- 🚀 CI/CD pipeline
- 🔒 Security scanning

## 💰 Estimación de Esfuerzo

| Fase | Complejidad | Tiempo Estimado | Prioridad |
|------|-------------|-----------------|-----------|
| Fase 1: Seguridad | Media | 2 semanas | 🔴 Crítica |
| Fase 2: Arquitectura | Alta | 4 semanas | 🟡 Alta |
| Fase 3: Observabilidad | Media | 2 semanas | 🟢 Media |
| Fase 4: Testing | Alta | 3 semanas | 🟢 Media |

**Total estimado: 11 semanas de desarrollo**

## 🎯 Beneficios Esperados

### Inmediatos (Fase 1)
- 🛡️ Seguridad robusta contra ataques comunes
- 📋 Validación consistente de datos
- 🔍 Visibilidad de eventos de seguridad

### Mediano Plazo (Fase 2-3)
- 🏗️ Código más mantenible y testeable
- ⚡ Performance mejorada con caching
- 📊 Visibilidad completa del sistema
- 🔧 Debugging más efectivo

### Largo Plazo (Fase 4)
- 🧪 Confianza en cambios con test suite
- 🚀 Deployments automatizados y seguros
- 📈 Capacidad de escalar con confianza
- 👥 Onboarding más rápido de developers

## 🚀 Próximos Pasos

1. **Revisar y aprobar roadmap** - Validar prioridades con stakeholders
2. **Configurar ambiente de desarrollo** - Docker, Redis, herramientas
3. **Implementar Fase 1** - Comenzar con mejoras críticas de seguridad
4. **Setup de testing** - Configurar framework de testing desde el inicio
5. **Documentar decisiones** - Mantener ADRs (Architecture Decision Records)

Este roadmap transformará la API de un MVP funcional a una aplicación production-ready con estándares enterprise, manteniendo la funcionalidad existente mientras se mejora dramáticamente la calidad, seguridad y mantenibilidad del código.