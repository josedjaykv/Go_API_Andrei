# 📄 Documentación Detallada del Código

## 📋 Índice por Archivos

- [🚀 main.go](#-maingo---punto-de-entrada)
- [🗄️ config/database.go](#-configdatabasego---configuración-de-bd)
- [🛣️ routes/routes.go](#-routesroutesgo---configuración-de-rutas)
- [🔐 middleware/auth.go](#-middlewareauthgo---autenticación)
- [🛡️ middleware/rbac.go](#-middlewarerbacgo---autorización)
- [🎯 controllers/auth.go](#-controllersauthgo---autenticación)
- [👑 controllers/andrei.go](#-controllersandreigo---funcionalidades-de-andrei)
- [👹 controllers/demon.go](#-controllersdemongodemonio-funcionalidades)
- [👨‍💻 controllers/network_admin.go](#-controllersnetwork_admingo---network-admin)
- [🌍 controllers/public.go](#-controllerspublicgo---endpoints-públicos)
- [👤 models/user.go](#-modelsusergo---modelo-de-usuario)
- [📝 models/post.go](#-modelspostgo---modelo-de-posts)
- [📊 models/report.go](#-modelsreportgo---modelo-de-reportes)
- [🏆 models/reward.go](#-modelsrewardgo---modelo-de-recompensas)
- [📈 models/statistics.go](#-modelsstatisticsgo---modelos-de-estadísticas)
- [🧪 cmd/seed/main.go](#-cmdseedmaingo---seeder-de-andrei)
- [🌱 cmd/populate/main.go](#-cmdpopulatemaingo---poblador-de-datos)

---

## 🚀 main.go - Punto de Entrada

**Ubicación:** `main.go`
**Propósito:** Inicialización y configuración del servidor

### Función Principal
```go
func main() {
    // Load environment variables
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found")
    }

    // Connect to database
    config.ConnectDatabase()

    // Create Gin router
    r := gin.Default()

    // Add CORS middleware
    r.Use(func(c *gin.Context) {
        c.Header("Access-Control-Allow-Origin", "*")
        c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }

        c.Next()
    })

    // Setup routes
    routes.SetupRoutes(r)

    // Get port from environment or use default
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    log.Printf("Server starting on port %s", port)
    if err := r.Run(":" + port); err != nil {
        log.Fatal("Failed to start server:", err)
    }
}
```

### Responsabilidades:
- **Líneas 15-17:** Carga de variables de entorno desde `.env`
- **Línea 20:** Inicialización de conexión a base de datos
- **Línea 23:** Creación del router Gin
- **Líneas 26-35:** Configuración de CORS para permitir requests desde frontend
- **Línea 38:** Configuración de todas las rutas
- **Líneas 41-48:** Inicio del servidor en puerto configurado

---

## 🗄️ config/database.go - Configuración de BD

**Ubicación:** `config/database.go`
**Propósito:** Gestión de conexión y migración de base de datos

### Variable Global
```go
var DB *gorm.DB
```

### Función de Conexión
```go
func ConnectDatabase() {
    dbHost := os.Getenv("DB_HOST")
    dbPort := os.Getenv("DB_PORT")
    dbUser := os.Getenv("DB_USER")
    dbPassword := os.Getenv("DB_PASSWORD")
    dbName := os.Getenv("DB_NAME")

    dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        dbHost, dbPort, dbUser, dbPassword, dbName)

    database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }

    err = database.AutoMigrate(
        &models.User{},
        &models.Post{},
        &models.Report{},
        &models.Reward{},
    )
    if err != nil {
        log.Fatal("Failed to migrate database:", err)
    }

    DB = database
    log.Println("Database connected and migrated successfully")
}
```

### Responsabilidades:
- **Líneas 12-16:** Lectura de variables de entorno de BD
- **Línea 18:** Construcción de DSN de PostgreSQL
- **Línea 20:** Conexión GORM a PostgreSQL
- **Líneas 25-31:** Migración automática de todas las entidades
- **Línea 33:** Asignación de conexión a variable global

---

## 🛣️ routes/routes.go - Configuración de Rutas

**Ubicación:** `routes/routes.go`
**Propósito:** Definición y organización de endpoints

### Función Principal
```go
func SetupRoutes(r *gin.Engine) {
    api := r.Group("/api/v1")

    // Public routes
    api.POST("/register", controllers.Register)
    api.POST("/login", controllers.Login)
    api.GET("/resistance", controllers.GetResistancePage)

    // Protected routes
    auth := api.Group("/")
    auth.Use(middleware.AuthRequired())

    // Andrei routes (admin only)
    andrei := auth.Group("/admin")
    andrei.Use(middleware.RequireAndrei())
    {
        andrei.GET("/users", controllers.GetAllUsers)
        andrei.GET("/users/:id", controllers.GetUserByID)
        andrei.DELETE("/users/:id", controllers.DeleteUser)
        andrei.POST("/rewards", controllers.CreateReward)
        andrei.GET("/stats", controllers.GetPlatformStats)
        andrei.GET("/demons/ranking", controllers.GetDemonRanking)
        andrei.GET("/posts", controllers.GetAllPosts)
        andrei.DELETE("/posts/:id", controllers.DeletePost)
        andrei.POST("/posts", controllers.CreateAndreiPost)
    }

    // Demon routes
    demons := auth.Group("/demons")
    demons.Use(middleware.RequireDemon())
    {
        demons.POST("/victims", controllers.RegisterVictim)
        demons.POST("/reports", controllers.CreateReport)
        demons.GET("/stats", controllers.GetMyStats)
        demons.GET("/victims", controllers.GetMyVictims)
        demons.GET("/reports", controllers.GetMyReports)
        demons.PUT("/reports/:id", controllers.UpdateReportStatus)
        demons.POST("/posts", controllers.CreateDemonPost)
    }

    // Network Admin routes
    networkAdmins := auth.Group("/network-admins")
    networkAdmins.Use(middleware.RequireNetworkAdmin())
    {
        networkAdmins.POST("/posts/anonymous", controllers.CreateAnonymousPost)
    }
}
```

### Organización Jerárquica:
- **Líneas 9-12:** Rutas públicas sin autenticación
- **Líneas 15-16:** Grupo de rutas protegidas con autenticación
- **Líneas 19-29:** Rutas exclusivas de Andrei con middleware de rol
- **Líneas 32-41:** Rutas exclusivas de demonios con middleware de rol
- **Líneas 44-48:** Rutas exclusivas de Network Admins con middleware de rol

---

## 🔐 middleware/auth.go - Autenticación

**Ubicación:** `middleware/auth.go`
**Propósito:** Validación de tokens JWT

### Middleware Principal
```go
func AuthRequired() gin.HandlerFunc {
    return gin.HandlerFunc(func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
            c.Abort()
            return
        }

        bearerToken := strings.Split(authHeader, " ")
        if len(bearerToken) != 2 || strings.ToLower(bearerToken[0]) != "bearer" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
            c.Abort()
            return
        }

        tokenString := bearerToken[1]
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            return []byte(os.Getenv("JWT_SECRET")), nil
        })

        if err != nil || !token.Valid {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }

        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
            c.Abort()
            return
        }

        userID := uint(claims["user_id"].(float64))
        var user models.User
        if err := config.DB.First(&user, userID).Error; err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
            c.Abort()
            return
        }

        c.Set("user", user)
        c.Next()
    })
}
```

### Proceso de Validación:
- **Líneas 15-20:** Verificación de header Authorization
- **Líneas 22-27:** Validación de formato Bearer token
- **Líneas 29-35:** Parsing y validación de JWT
- **Líneas 37-42:** Extracción de claims del token
- **Líneas 44-49:** Verificación de que el usuario existe en BD
- **Línea 51:** Almacenamiento del usuario en contexto

---

## 🛡️ middleware/rbac.go - Autorización

**Ubicación:** `middleware/rbac.go`
**Propósito:** Control de acceso basado en roles

### Middleware Base
```go
func RequireRole(allowedRoles ...models.UserRole) gin.HandlerFunc {
    return gin.HandlerFunc(func(c *gin.Context) {
        user, exists := c.Get("user")
        if !exists {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
            c.Abort()
            return
        }

        currentUser := user.(models.User)
        
        for _, role := range allowedRoles {
            if currentUser.Role == role {
                c.Next()
                return
            }
        }

        c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
        c.Abort()
    })
}
```

### Middleware Específicos
```go
func RequireAndrei() gin.HandlerFunc {
    return RequireRole(models.RoleAndrei)
}

func RequireDemon() gin.HandlerFunc {
    return RequireRole(models.RoleDemon)
}

func RequireNetworkAdmin() gin.HandlerFunc {
    return RequireRole(models.RoleNetworkAdmin)
}

func RequireAndreiOrDemon() gin.HandlerFunc {
    return RequireRole(models.RoleAndrei, models.RoleDemon)
}
```

### Funcionalidades:
- **RequireRole:** Middleware genérico que acepta múltiples roles
- **RequireAndrei:** Solo usuarios con rol 'andrei'
- **RequireDemon:** Solo usuarios con rol 'demon'
- **RequireNetworkAdmin:** Solo usuarios con rol 'network_admin'
- **RequireAndreiOrDemon:** Usuarios con rol 'andrei' o 'demon'

---

## 🎯 controllers/auth.go - Autenticación

**Ubicación:** `controllers/auth.go`
**Propósito:** Funcionalidades de registro y login

### Función Register
```go
func Register(c *gin.Context) {
    var input models.UserRegistration

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if input.Role != models.RoleDemon && input.Role != models.RoleNetworkAdmin {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role. Only demon and network_admin roles can register"})
        return
    }

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
        return
    }

    user := models.User{
        Username: input.Username,
        Email:    input.Email,
        Password: string(hashedPassword),
        Role:     input.Role,
    }

    if err := config.DB.Create(&user).Error; err != nil {
        c.JSON(http.StatusConflict, gin.H{"error": "Username or email already exists"})
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "message": "User registered successfully",
        "user": gin.H{
            "id":       user.ID,
            "username": user.Username,
            "email":    user.Email,
            "role":     user.Role,
        },
    })
}
```

### Función Login
```go
func Login(c *gin.Context) {
    var input models.UserLogin

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var user models.User
    if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": user.ID,
        "email":   user.Email,
        "role":    user.Role,
        "exp":     time.Now().Add(time.Hour * 24).Unix(),
    })

    tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
        return
    }

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

### Características Clave:
- **Registro:** Solo permite roles 'demon' y 'network_admin'
- **Hash de Contraseñas:** Uso de bcrypt para seguridad
- **Login:** Validación de credenciales y generación de JWT
- **Token JWT:** Expira en 24 horas, contiene datos del usuario

---

## 👑 controllers/andrei.go - Funcionalidades de Andrei

**Ubicación:** `controllers/andrei.go`
**Propósito:** Endpoints exclusivos del rol supremo

### Funciones Principales:

#### GetAllUsers - Ver todos los usuarios
```go
func GetAllUsers(c *gin.Context) {
    var users []models.User
    if err := config.DB.Find(&users).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"users": users})
}
```

#### CreateReward - Crear recompensas/castigos
```go
func CreateReward(c *gin.Context) {
    var input models.RewardCreate

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var demon models.User
    if err := config.DB.Where("id = ? AND role = ?", input.DemonID, models.RoleDemon).First(&demon).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Demon not found"})
        return
    }

    reward := models.Reward{
        DemonID:     input.DemonID,
        Type:        input.Type,
        Title:       input.Title,
        Description: input.Description,
        Points:      input.Points,
    }

    if err := config.DB.Create(&reward).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create reward"})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"reward": reward})
}
```

#### GetDemonRanking - Estadísticas de demonios
```go
func GetDemonRanking(c *gin.Context) {
    var demons []models.User
    if err := config.DB.Where("role = ?", models.RoleDemon).Find(&demons).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch demons"})
        return
    }

    var demonStats []models.DemonStats
    for _, demon := range demons {
        var stats models.DemonStats
        stats.DemonID = demon.ID

        config.DB.Model(&models.User{}).Where("role = ? AND id IN (SELECT victim_id FROM reports WHERE demon_id = ?)", 
            models.RoleNetworkAdmin, demon.ID).Count(&stats.VictimsCount)
        
        config.DB.Model(&models.Reward{}).Where("demon_id = ? AND type = ?", demon.ID, models.RewardTypeReward).Count(&stats.RewardsCount)
        config.DB.Model(&models.Reward{}).Where("demon_id = ? AND type = ?", demon.ID, models.RewardTypePunishment).Count(&stats.PunishmentsCount)
        config.DB.Model(&models.Report{}).Where("demon_id = ?", demon.ID).Count(&stats.ReportsCount)
        
        config.DB.Model(&models.Reward{}).Where("demon_id = ?", demon.ID).Select("COALESCE(SUM(points), 0)").Scan(&stats.TotalPoints)

        demonStats = append(demonStats, stats)
    }

    c.JSON(http.StatusOK, gin.H{"demon_rankings": demonStats})
}
```

### Capacidades de Andrei:
- **CRUD Usuarios:** Ver, obtener por ID, eliminar usuarios
- **Gestión de Recompensas:** Crear recompensas y castigos para demonios
- **Estadísticas:** Ver estadísticas de plataforma y ranking de demonios
- **Gestión de Posts:** CRUD completo sobre todas las publicaciones
- **Publicaciones:** Crear posts como Andrei identificada

---

## 👹 controllers/demon.go - Demonio Funcionalidades

**Ubicación:** `controllers/demon.go`
**Propósito:** Endpoints para el rol demon

### Funciones Principales:

#### RegisterVictim - Registrar víctimas
```go
func RegisterVictim(c *gin.Context) {
    var input models.UserRegistration

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if input.Role != models.RoleNetworkAdmin {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Only network_admin role can be registered as victim"})
        return
    }

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
        return
    }

    user := models.User{
        Username: input.Username,
        Email:    input.Email,
        Password: string(hashedPassword),
        Role:     input.Role,
    }

    if err := config.DB.Create(&user).Error; err != nil {
        c.JSON(http.StatusConflict, gin.H{"error": "Username or email already exists"})
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "message": "Victim registered successfully",
        "victim": gin.H{
            "id":       user.ID,
            "username": user.Username,
            "email":    user.Email,
            "role":     user.Role,
        },
    })
}
```

#### CreateReport - Crear reportes
```go
func CreateReport(c *gin.Context) {
    var input models.ReportCreate

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var victim models.User
    if err := config.DB.Where("id = ? AND role = ?", input.VictimID, models.RoleNetworkAdmin).First(&victim).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Victim not found"})
        return
    }

    user := c.MustGet("user").(models.User)

    report := models.Report{
        DemonID:     user.ID,
        VictimID:    input.VictimID,
        Title:       input.Title,
        Description: input.Description,
        Status:      "pending",
    }

    if err := config.DB.Create(&report).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create report"})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"report": report})
}
```

#### GetMyStats - Estadísticas personales
```go
func GetMyStats(c *gin.Context) {
    user := c.MustGet("user").(models.User)

    var stats models.DemonStats
    stats.DemonID = user.ID

    config.DB.Model(&models.User{}).Where("role = ? AND id IN (SELECT victim_id FROM reports WHERE demon_id = ?)", 
        models.RoleNetworkAdmin, user.ID).Count(&stats.VictimsCount)
    
    config.DB.Model(&models.Reward{}).Where("demon_id = ? AND type = ?", user.ID, models.RewardTypeReward).Count(&stats.RewardsCount)
    config.DB.Model(&models.Reward{}).Where("demon_id = ? AND type = ?", user.ID, models.RewardTypePunishment).Count(&stats.PunishmentsCount)
    config.DB.Model(&models.Report{}).Where("demon_id = ?", user.ID).Count(&stats.ReportsCount)
    
    config.DB.Model(&models.Reward{}).Where("demon_id = ?", user.ID).Select("COALESCE(SUM(points), 0)").Scan(&stats.TotalPoints)

    c.JSON(http.StatusOK, gin.H{"stats": stats})
}
```

### Capacidades de Demonios:
- **Gestión de Víctimas:** Registrar y ver víctimas capturadas
- **Reportes:** Crear, ver y actualizar reportes sobre víctimas
- **Estadísticas:** Ver estadísticas personales de desempeño
- **Publicaciones:** Crear posts identificados como demonio

---

## 👨‍💻 controllers/network_admin.go - Network Admin

**Ubicación:** `controllers/network_admin.go`
**Propósito:** Funcionalidades limitadas para víctimas

### Función Principal:

#### CreateAnonymousPost - Posts anónimos
```go
func CreateAnonymousPost(c *gin.Context) {
    var input models.PostCreate

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    post := models.Post{
        Title:     input.Title,
        Body:      input.Body,
        Media:     input.Media,
        AuthorID:  nil,
        Anonymous: true,
    }

    if err := config.DB.Create(&post).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"post": post})
}
```

### Capacidades de Network Admins:
- **Posts Anónimos:** Crear publicaciones sin identificación de autor
- **Acceso Limitado:** Solo funcionalidades básicas de resistencia

---

## 🌍 controllers/public.go - Endpoints Públicos

**Ubicación:** `controllers/public.go`
**Propósito:** Funcionalidades sin autenticación

### Función Principal:

#### GetResistancePage - Página de resistencia
```go
func GetResistancePage(c *gin.Context) {
    var posts []models.Post
    if err := config.DB.Preload("Author").Order("created_at DESC").Find(&posts).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch posts"})
        return
    }

    var publicPosts []gin.H
    for _, post := range posts {
        postData := gin.H{
            "id":         post.ID,
            "title":      post.Title,
            "body":       post.Body,
            "media":      post.Media,
            "anonymous":  post.Anonymous,
            "created_at": post.CreatedAt,
        }

        if post.Anonymous {
            postData["author"] = "Anonymous"
        } else if post.Author != nil {
            postData["author"] = post.Author.Username
        } else {
            postData["author"] = "Unknown"
        }

        publicPosts = append(publicPosts, postData)
    }

    c.JSON(http.StatusOK, gin.H{"posts": publicPosts})
}
```

### Características:
- **Acceso Público:** No requiere autenticación
- **Preload de Autor:** Carga información de autor si existe
- **Manejo de Anonimato:** Muestra "Anonymous" para posts sin autor

---

## 👤 models/user.go - Modelo de Usuario

**Ubicación:** `models/user.go`
**Propósito:** Definición de entidad User y relacionados

### Estructura Principal:
```go
type User struct {
    ID        uint           `json:"id" gorm:"primaryKey"`
    Username  string         `json:"username" gorm:"unique;not null"`
    Email     string         `json:"email" gorm:"unique;not null"`
    Password  string         `json:"-" gorm:"not null"`
    Role      UserRole       `json:"role" gorm:"not null"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
    
    Posts     []Post     `json:"posts,omitempty" gorm:"foreignKey:AuthorID"`
    Reports   []Report   `json:"reports,omitempty" gorm:"foreignKey:DemonID"`
    Rewards   []Reward   `json:"rewards,omitempty" gorm:"foreignKey:DemonID"`
}
```

### DTOs de Usuario:
```go
type UserRegistration struct {
    Username string   `json:"username" binding:"required"`
    Email    string   `json:"email" binding:"required,email"`
    Password string   `json:"password" binding:"required,min=6"`
    Role     UserRole `json:"role" binding:"required"`
}

type UserLogin struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required"`
}
```

### Características:
- **Soft Delete:** Usando `DeletedAt` con GORM
- **Relaciones:** Posts, Reports, Rewards como foreign keys
- **Validaciones:** Campos requeridos, formato email, longitud mínima
- **Seguridad:** Password excluido de JSON con `json:"-"`

---

## 📝 models/post.go - Modelo de Posts

**Ubicación:** `models/post.go`
**Propósito:** Definición de publicaciones

### Estructura:
```go
type Post struct {
    ID        uint           `json:"id" gorm:"primaryKey"`
    Title     string         `json:"title" gorm:"not null"`
    Body      string         `json:"body" gorm:"not null"`
    Media     string         `json:"media,omitempty"`
    AuthorID  *uint          `json:"author_id,omitempty"`
    Author    *User          `json:"author,omitempty" gorm:"foreignKey:AuthorID"`
    Anonymous bool           `json:"anonymous" gorm:"default:false"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
```

### Características Especiales:
- **Autor Opcional:** `AuthorID` puede ser `nil` para posts anónimos
- **Flag de Anonimato:** `Anonymous` indica si el post es anónimo
- **Media:** Campo opcional para URLs de multimedia

---

## 📊 models/report.go - Modelo de Reportes

**Ubicación:** `models/report.go`
**Propósito:** Reportes de demonios sobre víctimas

### Estructura:
```go
type Report struct {
    ID          uint           `json:"id" gorm:"primaryKey"`
    DemonID     uint           `json:"demon_id" gorm:"not null"`
    Demon       User           `json:"demon" gorm:"foreignKey:DemonID"`
    VictimID    uint           `json:"victim_id" gorm:"not null"`
    Victim      User           `json:"victim" gorm:"foreignKey:VictimID"`
    Title       string         `json:"title" gorm:"not null"`
    Description string         `json:"description" gorm:"not null"`
    Status      string         `json:"status" gorm:"default:'pending'"`
    CreatedAt   time.Time      `json:"created_at"`
    UpdatedAt   time.Time      `json:"updated_at"`
    DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}
```

### Estados Posibles:
- `pending`: Recién creado
- `in_progress`: En progreso  
- `completed`: Completado

---

## 🏆 models/reward.go - Modelo de Recompensas

**Ubicación:** `models/reward.go`
**Propósito:** Sistema de recompensas/castigos

### Estructura:
```go
type Reward struct {
    ID          uint           `json:"id" gorm:"primaryKey"`
    DemonID     uint           `json:"demon_id" gorm:"not null"`
    Demon       User           `json:"demon" gorm:"foreignKey:DemonID"`
    Type        RewardType     `json:"type" gorm:"not null"`
    Title       string         `json:"title" gorm:"not null"`
    Description string         `json:"description" gorm:"not null"`
    Points      int            `json:"points" gorm:"default:0"`
    CreatedAt   time.Time      `json:"created_at"`
    UpdatedAt   time.Time      `json:"updated_at"`
    DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}
```

### Tipos:
```go
const (
    RewardTypeReward     RewardType = "reward"     // Puntos positivos
    RewardTypePunishment RewardType = "punishment" // Puntos negativos
)
```

---

## 📈 models/statistics.go - Modelos de Estadísticas

**Ubicación:** `models/statistics.go`
**Propósito:** DTOs para estadísticas calculadas

### Estructuras:
```go
type DemonStats struct {
    DemonID        uint  `json:"demon_id"`
    VictimsCount   int64 `json:"victims_count"`
    RewardsCount   int64 `json:"rewards_count"`
    PunishmentsCount int64  `json:"punishments_count"`
    TotalPoints    int64  `json:"total_points"`
    ReportsCount   int64  `json:"reports_count"`
}

type PlatformStats struct {
    TotalUsers       int64 `json:"total_users"`
    TotalDemons      int64 `json:"total_demons"`
    TotalNetworkAdmins int64 `json:"total_network_admins"`
    TotalPosts       int64 `json:"total_posts"`
    TotalReports     int64 `json:"total_reports"`
}
```

### Uso:
- **DemonStats:** Estadísticas individuales de cada demonio
- **PlatformStats:** Estadísticas generales de la plataforma

---

## 🧪 cmd/seed/main.go - Seeder de Andrei

**Ubicación:** `cmd/seed/main.go`
**Propósito:** Crear usuario inicial Andrei

### Funcionalidad Principal:
```go
func main() {
    // Load environment variables desde la raíz
    if err := godotenv.Load(".env"); err != nil {
        log.Println("⚠️ No .env file found (usando variables del sistema)")
    }

    // Connect to database
    config.ConnectDatabase()

    // Check if Andrei user exists
    var existingUser models.User
    err := config.DB.Where("email = ?", "andrei@evil.com").First(&existingUser).Error
    if err == nil {
        log.Println("✅ Andrei user already exists")
        return
    }
    if err != nil && err != gorm.ErrRecordNotFound {
        log.Fatal("❌ Error checking for Andrei user:", err)
    }

    // Hash password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
    if err != nil {
        log.Fatal("❌ Failed to hash password:", err)
    }

    // Create Andrei user
    andrei := models.User{
        Username: "AndreiMesManur",
        Email:    "andrei@evil.com",
        Password: string(hashedPassword),
        Role:     models.RoleAndrei,
    }

    if err := config.DB.Create(&andrei).Error; err != nil {
        log.Fatal("❌ Failed to create Andrei user:", err)
    }

    log.Println("🎉 Andrei user created successfully!")
    log.Println("   Email:    andrei@evil.com")
    log.Println("   Password: password123")
}
```

### Características:
- **Verificación de Existencia:** No crea duplicados
- **Credenciales Fijas:** Email y password predefinidos
- **Logging Detallado:** Emojis y mensajes claros

---

## 🌱 cmd/populate/main.go - Poblador de Datos

**Ubicación:** `cmd/populate/main.go`
**Propósito:** Crear datos de prueba realistas

### Funciones Principales:

#### createDemons - Crear demonios
```go
func createDemons() []models.User {
    demons := []models.User{
        {Username: "ShadowMaster", Email: "shadow@evil.com", Role: models.RoleDemon},
        {Username: "DarkLord666", Email: "darklord@evil.com", Role: models.RoleDemon},
        {Username: "ChaosReaper", Email: "chaos@evil.com", Role: models.RoleDemon},
        {Username: "VoidWhisperer", Email: "void@evil.com", Role: models.RoleDemon},
        {Username: "NightmareKing", Email: "nightmare@evil.com", Role: models.RoleDemon},
    }

    hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("demon123"), bcrypt.DefaultCost)

    for i := range demons {
        demons[i].Password = string(hashedPassword)
        if err := config.DB.Create(&demons[i]).Error; err != nil {
            log.Printf("⚠️ Error creando demonio %s: %v", demons[i].Username, err)
        }
    }

    return demons
}
```

#### createPosts - Crear publicaciones variadas
```go
func createPosts(demons, networkAdmins []models.User) []models.Post {
    // Get Andrei user
    var andrei models.User
    config.DB.Where("role = ?", models.RoleAndrei).First(&andrei)

    posts := []struct {
        authorID  *uint
        title     string
        body      string
        media     string
        anonymous bool
    }{
        // Andrei's posts
        {&andrei.ID, "Welcome to the New Order", "My loyal demons, the time has come to expand our dominion over the digital realm. Every network administrator captured brings us closer to total control!", "", false},
        
        // Demon posts
        {&demons[0].ID, "Infiltration Techniques", "Brothers, I've discovered that posing as IT support is incredibly effective. Humans trust anyone who claims they can fix their computer problems.", "", false},
        
        // Anonymous posts from Network Admins
        {nil, "RESIST THE DARKNESS", "Fellow administrators, do not give in to their influence! We must fight back against these digital demons!", "", true},
        // ... más posts
    }
    
    // Crear todos los posts...
}
```

### Datos Creados:
- **5 Demonios** con nombres temáticos
- **8 Network Admins** representando víctimas reales
- **10 Reportes** con estados variados
- **10 Recompensas/Castigos** con puntos balanceados
- **12 Posts** mezclando autores identificados y anónimos

---

*Esta documentación proporciona una referencia completa del código con ubicaciones exactas de archivos y líneas específicas para facilitar el mantenimiento y desarrollo futuro.*