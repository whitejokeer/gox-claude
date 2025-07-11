# GOX Framework - Documentación Completa

## 📚 Índice

1. [Introducción](#introducción)
2. [Instalación](#instalación)
3. [Primeros Pasos](#primeros-pasos)
4. [Comandos CLI](#comandos-cli)
5. [Arquitectura](#arquitectura)
6. [Formato .gox](#formato-gox)
7. [Sistema de Routing](#sistema-de-routing)
8. [Middleware y Autenticación](#middleware-y-autenticación)
9. [Base de Datos](#base-de-datos)
10. [Storybook](#storybook)
11. [Configuración](#configuración)
12. [Ejemplos](#ejemplos)

---

## 🎯 Introducción

GOX Framework es un framework de desarrollo web moderno que combina el poder de Go con la simplicidad de HTMX para crear aplicaciones web ultrarrápidas mediante Server-Side Rendering.

### ¿Por qué GOX?

- **🚀 Ultra Performance**: SSR nativo + Go + HTMX
- **🎨 Developer Experience**: Herramientas integradas
- **🏗️ Arquitectura Distribuida**: Escalable desde el día uno
- **📦 Zero Bundle**: Sin JavaScript cliente pesado
- **🔧 Fácil de aprender**: Curva de aprendizaje suave
- **💎 Difícil de masterizar**: Opciones avanzadas potentes

---

## 📦 Instalación

### Requisitos
- Go 1.21+
- Node.js 18+ (para herramientas de desarrollo)

### Instalación del CLI
```bash
# Usando Go
go install github.com/gox-framework/gox@latest

# Usando curl
curl -sSL https://get.gox.dev | sh

# Usando Homebrew (macOS)
brew install gox-framework/tap/gox
```

### Verificar instalación
```bash
gox version
# GOX Framework v1.0.0
```

---

## 🚀 Primeros Pasos

### 1. Crear nuevo proyecto
```bash
gox new project my-app
cd my-app
```

### 2. Estructura generada
```
my-app/
├── pages/
│   ├── index.gox
│   └── _layout.gox
├── services/
├── routing/
├── common/
├── infra/
├── gox.config.yaml
├── go.mod
└── README.md
```

### 3. Ejecutar en desarrollo
```bash
gox dev
# 🚀 Server running on http://localhost:3000
# 📖 Storybook running on http://localhost:6006
```

### 4. Tu primera página
```gox
<!-- pages/hello.gox -->
<template>
  <div class="p-4">
    <h1 class="text-2xl font-bold">Hello {{.Name}}!</h1>
    <button hx-post="/api/greet" hx-target="#response">
      Greet
    </button>
    <div id="response"></div>
  </div>
</template>

<go>
type HelloPage struct {
  Name string
}

func (h *HelloPage) Mount(ctx *gox.Context) error {
  h.Name = "World"
  return nil
}

func (h *HelloPage) Greet(ctx *gox.Context) error {
  return ctx.HTML("<p>Hello from GOX!</p>")
}
</go>
```

---

## 🛠️ Comandos CLI

### Comandos de Proyecto

#### `gox new`
Crea nuevos proyectos, módulos o servicios.

```bash
# Crear proyecto completo
gox new project my-app [flags]

# Crear módulo independiente
gox new module user-management [flags]

# Crear servicio puro
gox new service auth-service [flags]
```

**Flags:**
- `--template, -t`: Template a usar (default, minimal, enterprise)
- `--db`: Base de datos (postgres, mysql, sqlite, none)
- `--auth`: Sistema de auth (jwt, oauth2, session, none)
- `--protocol`: Protocolo comunicación (grpc, rest, graphql)
- `--no-storybook`: Deshabilitar storybook
- `--no-tailwind`: Deshabilitar Tailwind CSS

**Ejemplos:**
```bash
# Proyecto con PostgreSQL y JWT
gox new project e-commerce --db=postgres --auth=jwt

# Módulo con gRPC
gox new module payments --protocol=grpc

# Servicio minimal
gox new service notifications --template=minimal
```

#### `gox generate`
Genera componentes, páginas y servicios.

```bash
# Generar página
gox generate page <name> [flags]

# Generar componente
gox generate component <name> [flags]

# Generar servicio
gox generate service <name> [flags]

# Generar middleware
gox generate middleware <name> [flags]
```

**Flags:**
- `--path, -p`: Ruta donde generar el archivo
- `--template, -t`: Template específico
- `--props`: Props del componente (separadas por coma)
- `--auth`: Requiere autenticación
- `--api`: Generar endpoints API

**Ejemplos:**
```bash
# Página con autenticación
gox generate page dashboard --auth

# Componente con props
gox generate component user-card --props="name,email,avatar"

# Servicio con API
gox generate service users --api
```

### Comandos de Desarrollo

#### `gox dev`
Inicia el servidor de desarrollo.

```bash
gox dev [flags]
```

**Flags:**
- `--port, -p`: Puerto (default: 3000)
- `--host`: Host (default: localhost)
- `--no-storybook`: Deshabilitar storybook
- `--no-hot-reload`: Deshabilitar hot reload
- `--verbose, -v`: Modo verbose

#### `gox build`
Construye la aplicación para producción.

```bash
gox build [flags]
```

**Flags:**
- `--output, -o`: Directorio de salida
- `--target`: Target platform (linux, darwin, windows)
- `--compress`: Comprimir binarios
- `--static`: Build estático

#### `gox test`
Ejecuta las pruebas.

```bash
gox test [flags]
```

**Flags:**
- `--coverage`: Generar reporte de cobertura
- `--verbose, -v`: Modo verbose
- `--watch`: Modo watch

### Comandos de Gestión

#### `gox install`
Instala dependencias del proyecto.

```bash
gox install [package]
```

#### `gox add`
Agrega servicios o módulos al proyecto.

```bash
gox add service <name>
gox add module <name>
```

#### `gox remove`
Remueve servicios o módulos del proyecto.

```bash
gox remove service <name>
gox remove module <name>
```

### Comandos de Utilidad

#### `gox migrate`
Maneja migraciones de base de datos.

```bash
gox migrate up      # Ejecutar migraciones
gox migrate down    # Revertir migraciones
gox migrate create  # Crear nueva migración
gox migrate status  # Ver estado
```

#### `gox deploy`
Despliega la aplicación.

```bash
gox deploy [flags]
```

**Flags:**
- `--target`: Target (docker, k8s, cloud-run, lambda)
- `--env`: Environment (dev, staging, prod)
- `--config`: Archivo de configuración

---

## 🏗️ Arquitectura

### Estructura Detallada

```
my-app/
├── pages/                  # 🎨 Páginas frontend (.gox)
│   ├── index.gox          # Página principal
│   ├── about.gox          # Página estática
│   ├── auth/              # Páginas de autenticación
│   │   ├── login.gox
│   │   └── register.gox
│   ├── dashboard/         # Páginas dinámicas
│   │   ├── index.gox
│   │   └── [id].gox       # Página dinámica
│   └── _layout.gox        # Layout global
│
├── services/              # 🔧 Servicios y módulos
│   ├── auth/             # Servicio de autenticación
│   │   ├── cmd/          # Comando principal
│   │   │   └── main.go
│   │   ├── internal/     # Lógica interna
│   │   │   ├── handlers/
│   │   │   ├── models/
│   │   │   └── service/
│   │   ├── api/          # Definición API
│   │   │   └── auth.proto
│   │   └── service.yaml  # Configuración
│   │
│   └── users/            # Módulo con componentes
│       ├── cmd/
│       ├── internal/
│       ├── api/
│       ├── components/   # Componentes del módulo
│       │   ├── user-card.gox
│       │   └── user-list.gox
│       └── module.yaml
│
├── routing/              # 🛣️ Configuración de rutas
│   ├── routes.go         # Rutas manuales (opcional)
│   ├── middleware.go     # Middleware global
│   └── handlers.go       # Handlers compartidos
│
├── common/               # 📦 Código compartido
│   ├── types/            # Tipos compartidos
│   │   ├── user.go
│   │   └── response.go
│   ├── utils/            # Utilidades
│   │   ├── crypto.go
│   │   └── validation.go
│   └── contracts/        # Interfaces entre servicios
│       ├── auth.go
│       └── users.go
│
├── infra/                # 🏗️ Infraestructura
│   ├── database/         # Configuración DB
│   │   ├── migrations/
│   │   ├── models/
│   │   └── connection.go
│   ├── cache/            # Cache (Redis, etc.)
│   │   └── redis.go
│   └── messaging/        # Message bus
│       ├── nats.go
│       └── events.go
│
├── static/               # 📁 Archivos estáticos
│   ├── css/
│   ├── js/
│   └── images/
│
├── docs/                 # 📚 Documentación
│   └── api/              # Documentación API
│
├── scripts/              # 📝 Scripts de desarrollo
│   ├── setup.sh
│   └── deploy.sh
│
├── .gox/                 # 🔧 Archivos del framework
│   ├── generated/        # Código generado
│   ├── cache/            # Cache de compilación
│   └── storybook/        # Configuración storybook
│
├── gox.config.yaml       # ⚙️ Configuración principal
├── go.mod                # 📦 Dependencias Go
├── package.json          # 📦 Dependencias Node (dev)
├── tailwind.config.js    # 🎨 Configuración Tailwind
└── README.md             # 📖 Documentación
```

### Tipos de Archivos

#### Páginas (.gox)
- **Ubicación**: `pages/`
- **Propósito**: Interfaz de usuario
- **Routing**: Automático basado en estructura
- **Características**: Template + Go + Styles

#### Servicios
- **Ubicación**: `services/*/`
- **Propósito**: Lógica de negocio
- **Protocolo**: gRPC, REST, GraphQL
- **Estructura**: Go estándar (cmd, internal, api)

#### Módulos
- **Ubicación**: `services/*/`
- **Propósito**: Servicio + Componentes UI
- **Características**: Combina backend y frontend
- **Reutilización**: Componentes exportables

---

## 📄 Formato .gox

### Estructura Básica

```gox
<template>
  <!-- HTML con Go templates + HTMX -->
</template>

<go>
// Código Go puro
</go>

<style>
/* CSS con soporte Tailwind */
</style>
```

### Sección Template

```gox
<template>
  <div class="container mx-auto p-4">
    <h1 class="text-3xl font-bold">{{.Title}}</h1>
    
    <!-- Iteración -->
    {{range .Users}}
      <div class="user-card">
        <h2>{{.Name}}</h2>
        <p>{{.Email}}</p>
      </div>
    {{end}}
    
    <!-- Condicionales -->
    {{if .IsAuthenticated}}
      <button hx-post="/api/logout">Logout</button>
    {{else}}
      <a href="/auth/login">Login</a>
    {{end}}
    
    <!-- HTMX -->
    <div hx-get="/api/stats" 
         hx-trigger="load"
         hx-target="#stats">
      <div id="stats">Loading...</div>
    </div>
  </div>
</template>
```

### Sección Go

```gox
<go>
// Importaciones
import (
  "context"
  "fmt"
  "github.com/gox-framework/gox"
)

// Estructura de datos
type DashboardPage struct {
  Title           string
  Users           []User
  IsAuthenticated bool
}

type User struct {
  ID    int    `json:"id"`
  Name  string `json:"name"`
  Email string `json:"email"`
}

// Lifecycle hooks
func (d *DashboardPage) Mount(ctx *gox.Context) error {
  d.Title = "Dashboard"
  d.IsAuthenticated = ctx.IsAuthenticated()
  
  // Cargar datos
  users, err := d.loadUsers(ctx)
  if err != nil {
    return err
  }
  d.Users = users
  
  return nil
}

func (d *DashboardPage) BeforeMount(ctx *gox.Context) error {
  // Validaciones antes del mount
  return nil
}

func (d *DashboardPage) AfterMount(ctx *gox.Context) error {
  // Lógica después del mount
  return nil
}

// Handlers HTMX
func (d *DashboardPage) GetStats(ctx *gox.Context) error {
  stats := map[string]interface{}{
    "users":  len(d.Users),
    "active": d.countActiveUsers(),
  }
  
  return ctx.JSON(stats)
}

func (d *DashboardPage) CreateUser(ctx *gox.Context) error {
  var user User
  if err := ctx.Bind(&user); err != nil {
    return ctx.BadRequest("Invalid data")
  }
  
  // Crear usuario
  if err := d.saveUser(&user); err != nil {
    return ctx.InternalError("Failed to create user")
  }
  
  // Devolver HTML actualizado
  return ctx.HTML(`<div class="user-card">
    <h2>` + user.Name + `</h2>
    <p>` + user.Email + `</p>
  </div>`)
}

// Métodos privados
func (d *DashboardPage) loadUsers(ctx *gox.Context) ([]User, error) {
  // Lógica para cargar usuarios
  return []User{}, nil
}

func (d *DashboardPage) countActiveUsers() int {
  count := 0
  for _, user := range d.Users {
    if user.IsActive {
      count++
    }
  }
  return count
}

func (d *DashboardPage) saveUser(user *User) error {
  // Lógica para guardar usuario
  return nil
}
</go>
```

### Sección Style

```gox
<style>
/* CSS global */
.container {
  max-width: 1200px;
  margin: 0 auto;
}

.user-card {
  @apply bg-white shadow-lg rounded-lg p-4 mb-4;
  border-left: 4px solid #3b82f6;
}

.user-card h2 {
  @apply text-xl font-semibold text-gray-800;
}

.user-card p {
  @apply text-gray-600;
}
</style>

<!-- O con scope -->
<style scoped>
/* CSS solo para este componente */
.dashboard-header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 2rem;
  border-radius: 0.5rem;
}
</style>
```

### Propiedades Especiales

#### Template con Autenticación
```gox
<template auth="required">
  <!-- Solo usuarios autenticados -->
</template>

<template auth="role:admin">
  <!-- Solo administradores -->
</template>
```

#### Template con Layout
```gox
<template layout="admin">
  <!-- Usa layout específico -->
</template>
```

#### Metadata
```gox
<template>
  <head>
    <title>{{.Title}}</title>
    <meta name="description" content="{{.Description}}">
    <meta property="og:title" content="{{.Title}}">
  </head>
  <body>
    <!-- Contenido -->
  </body>
</template>
```

---

## 🛣️ Sistema de Routing

### Routing Automático (File-based)

El sistema de routing automático mapea archivos a rutas:

```
pages/
├── index.gox              → /
├── about.gox              → /about
├── contact.gox            → /contact
├── blog/
│   ├── index.gox          → /blog
│   ├── [slug].gox         → /blog/:slug
│   └── categories/
│       ├── index.gox      → /blog/categories
│       └── [category].gox → /blog/categories/:category
├── users/
│   ├── index.gox          → /users
│   ├── [id].gox           → /users/:id
│   └── [id]/
│       ├── edit.gox       → /users/:id/edit
│       └── profile.gox    → /users/:id/profile
├── admin/
│   ├── _layout.gox        → Layout para /admin/*
│   ├── index.gox          → /admin
│   └── users.gox          → /admin/users
└── api/
    ├── auth.gox           → /api/auth
    ├── users/
    │   ├── index.gox      → /api/users
    │   └── [id].gox       → /api/users/:id
    └── [...].gox          → /api/* (catch-all)
```

### Parámetros Dinámicos

#### Parámetros Simples
```gox
<!-- pages/users/[id].gox -->
<template>
  <div>
    <h1>User {{.UserID}}</h1>
    <p>Name: {{.User.Name}}</p>
  </div>
</template>

<go>
type UserPage struct {
  UserID string
  User   User
}

func (u *UserPage) Mount(ctx *gox.Context) error {
  u.UserID = ctx.Param("id")
  
  // Cargar usuario
  user, err := u.loadUser(u.UserID)
  if err != nil {
    return ctx.NotFound("User not found")
  }
  u.User = user
  
  return nil
}
</go>
```

#### Parámetros Múltiples
```gox
<!-- pages/blog/[year]/[month]/[slug].gox -->
<go>
func (b *BlogPostPage) Mount(ctx *gox.Context) error {
  b.Year = ctx.Param("year")
  b.Month = ctx.Param("month")
  b.Slug = ctx.Param("slug")
  return nil
}
</go>
```

#### Catch-all
```gox
<!-- pages/docs/[...].gox -->
<go>
func (d *DocsPage) Mount(ctx *gox.Context) error {
  // ctx.Param("...") contiene toda la ruta
  d.Path = ctx.Param("...")
  return nil
}
</go>
```

### Routing Manual

```go
// routing/routes.go
package routing

import (
  "github.com/gox-framework/gox"
  "my-app/pages"
  "my-app/services"
)

func SetupRoutes(r *gox.Router) {
  // Middleware global
  r.Use(middleware.Logger())
  r.Use(middleware.CORS())
  
  // Rutas públicas
  r.Page("/", pages.Home)
  r.Page("/about", pages.About)
  r.Page("/contact", pages.Contact)
  
  // Grupo con autenticación
  auth := r.Group("/", middleware.Auth())
  {
    auth.Page("/dashboard", pages.Dashboard)
    auth.Page("/profile", pages.Profile)
  }
  
  // Grupo admin
  admin := r.Group("/admin", middleware.RequireRole("admin"))
  {
    admin.Page("/", pages.AdminDashboard)
    admin.Page("/users", pages.AdminUsers)
    admin.Page("/settings", pages.AdminSettings)
  }
  
  // API Routes
  api := r.Group("/api")
  {
    api.Use(middleware.JSON())
    
    // Servicios
    api.Mount("/auth", services.Auth)
    api.Mount("/users", services.Users)
    api.Mount("/posts", services.Posts)
    
    // Rutas específicas
    api.GET("/health", handlers.Health)
    api.POST("/upload", handlers.Upload)
  }
  
  // Archivos estáticos
  r.Static("/static", "./static")
  
  // Catch-all para SPA
  r.NoRoute(pages.NotFound)
}
```

### Middleware

#### Middleware Global
```go
// routing/middleware.go
func Logger() gox.Middleware {
  return func(ctx *gox.Context) error {
    start := time.Now()
    
    err := ctx.Next()
    
    log.Printf("%s %s %v", 
      ctx.Method(), 
      ctx.Path(), 
      time.Since(start))
    
    return err
  }
}

func CORS() gox.Middleware {
  return func(ctx *gox.Context) error {
    ctx.Header("Access-Control-Allow-Origin", "*")
    ctx.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
    ctx.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
    
    if ctx.Method() == "OPTIONS" {
      return ctx.NoContent(204)
    }
    
    return ctx.Next()
  }
}
```

#### Middleware de Autenticación
```go
func Auth() gox.Middleware {
  return func(ctx *gox.Context) error {
    token := ctx.Header("Authorization")
    if token == "" {
      return ctx.Unauthorized("Missing token")
    }
    
    user, err := auth.ValidateToken(token)
    if err != nil {
      return ctx.Unauthorized("Invalid token")
    }
    
    ctx.Set("user", user)
    return ctx.Next()
  }
}

func RequireRole(role string) gox.Middleware {
  return func(ctx *gox.Context) error {
    user := ctx.Get("user").(*User)
    if user.Role != role {
      return ctx.Forbidden("Insufficient permissions")
    }
    return ctx.Next()
  }
}
```

---

## 🔐 Middleware y Autenticación

### Sistema de Autenticación

#### Configuración
```yaml
# gox.config.yaml
auth:
  provider: "jwt"           # jwt, oauth2, session
  secret: "${JWT_SECRET}"   # Variable de entorno
  expires: "24h"            # Tiempo de expiración
  refresh: true             # Tokens de refresh
  service: "auth"           # Servicio que maneja auth
```

#### Servicio de Autenticación
```go
// services/auth/internal/service/auth.go
package service

type AuthService struct {
  db     *gorm.DB
  config *Config
}

func (s *AuthService) Login(email, password string) (*TokenResponse, error) {
  user, err := s.validateCredentials(email, password)
  if err != nil {
    return nil, err
  }
  
  token, err := s.generateToken(user)
  if err != nil {
    return nil, err
  }
  
  return &TokenResponse{
    Token:     token,
    ExpiresAt: time.Now().Add(24 * time.Hour),
    User:      user,
  }, nil
}

func (s *AuthService) ValidateToken(token string) (*User, error) {
  claims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
    return []byte(s.config.Secret), nil
  })
  
  if err != nil {
    return nil, err
  }
  
  if !claims.Valid {
    return nil, errors.New("invalid token")
  }
  
  userID := claims.Claims.(*Claims).UserID
  return s.getUserByID(userID)
}
```

### Middleware de Autenticación

#### Middleware JWT
```go
// routing/middleware/auth.go
func JWT() gox.Middleware {
  return func(ctx *gox.Context) error {
    token := extractToken(ctx)
    if token == "" {
      return ctx.Unauthorized("Missing token")
    }
    
    user, err := auth.ValidateToken(token)
    if err != nil {
      return ctx.Unauthorized("Invalid token")
    }
    
    ctx.Set("user", user)
    ctx.Set("authenticated", true)
    
    return ctx.Next()
  }
}

func extractToken(ctx *gox.Context) string {
  // Header Authorization: Bearer <token>
  auth := ctx.Header("Authorization")
  if strings.HasPrefix(auth, "Bearer ") {
    return strings.TrimPrefix(auth, "Bearer ")
  }
  
  // Query parameter
  if token := ctx.Query("token"); token != "" {
    return token
  }
  
  // Cookie
  if cookie, err := ctx.Cookie("auth_token"); err == nil {
    return cookie.Value
  }
  
  return ""
}
```

#### Middleware de Roles
```go
func RequireRole(roles ...string) gox.Middleware {
  return func(ctx *gox.Context) error {
    user := ctx.Get("user").(*User)
    
    for _, role := range roles {
      if user.HasRole(role) {
        return ctx.Next()
      }
    }
    
    return ctx.Forbidden("Insufficient permissions")
  }
}

func RequirePermission(permission string) gox.Middleware {
  return func(ctx *gox.Context) error {
    user := ctx.Get("user").(*User)
    
    if user.HasPermission(permission) {
      return ctx.Next()
    }
    
    return ctx.Forbidden("Permission denied")
  }
}
```

### Uso en Páginas

#### Autenticación Requerida
```gox
<template auth="required">
  <div>
    <h1>Welcome {{.User.Name}}</h1>
    <p>Email: {{.User.Email}}</p>
  </div>
</template>

<go>
type ProfilePage struct {
  User User
}

func (p *ProfilePage) Mount(ctx *gox.Context) error {
  p.User = ctx.Get("user").(User)
  return nil
}
</go>
```

#### Roles Específicos
```gox
<template auth="role:admin">
  <div>
    <h1>Admin Dashboard</h1>
    <!-- Contenido solo para administradores -->
  </div>
</template>
```

#### Permisos Específicos
```gox
<template auth="permission:users.manage">
  <div>
    <h1>User Management</h1>
    <!-- Contenido solo para usuarios con permiso -->
  </div>
</template>
```

### Middleware Personalizados

#### Rate Limiting
```go
func RateLimit(requests int, window time.Duration) gox.Middleware {
  limiter := rate.NewLimiter(rate.Every(window/time.Duration(requests)), requests)
  
  return func(ctx *gox.Context) error {
    if !limiter.Allow() {
      return ctx.TooManyRequests("Rate limit exceeded")
    }
    return ctx.Next()
  }
}
```

#### Validación de Input
```go
func ValidateJSON(schema interface{}) gox.Middleware {
  return func(ctx *gox.Context) error {
    if err := ctx.Bind(schema); err != nil {
      return ctx.BadRequest("Invalid JSON")
    }
    
    if err := validator.Validate(schema); err != nil {
      return ctx.BadRequest(err.Error())
    }
    
    ctx.Set("validated", schema)
    return ctx.Next()
  }
}
```

---

## 🗄️ Base de Datos

### Configuración

#### Estrategia Compartida
```yaml
# gox.config.yaml
database:
  strategy: "shared"
  driver: "postgres"
  host: "localhost"
  port: 5432
  name: "myapp"
  user: "postgres"
  password: "${DB_PASSWORD}"
  ssl_mode: "disable"
  
  migrations:
    auto: true
    path: "./infra/database/migrations"
```

#### Estrategia por Servicio
```yaml
database:
  strategy: "per-service"
  default_driver: "postgres"
  
  services:
    auth:
      driver: "postgres"
      name: "auth_db"
    users:
      driver: "postgres"
      name: "users_db"
    posts:
      driver: "mysql"
      name: "posts_db"
  
  sync:
    method: "event-sourcing"
    bus: "nats"
```

### Modelos

#### Modelo Base
```go
// common/types/base.go
type BaseModel struct {
  ID        uint      `json:"id" gorm:"primaryKey"`
  CreatedAt time.Time `json:"created_at"`
  UpdatedAt time.Time `json:"updated_at"`
  DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
```

#### Modelos Específicos
```go
// common/types/user.go
type User struct {
  BaseModel
  Email     string    `json:"email" gorm:"uniqueIndex"`
  Name      string    `json:"name"`
  Password  string    `json:"-"`
  Role      string    `json:"role" gorm:"default:user"`
  IsActive  bool      `json:"is_active" gorm:"default:true"`
  Profile   Profile   `json:"profile"`
  Posts     []Post    `json:"posts"`
}

type Profile struct {
  BaseModel
  UserID    uint   `json:"user_id" gorm:"uniqueIndex"`
  Bio       string `json:"bio"`
  Avatar    string `json:"avatar"`
  Website   string `json:"website"`
}

type Post struct {
  BaseModel
  UserID    uint     `json:"user_id"`
  Title     string   `json:"title"`
  Content   string   `json:"content"`
  Status    string   `json:"status" gorm:"default:draft"`
  Tags      []Tag    `json:"tags" gorm:"many2many:post_tags;"`
  User      User     `json:"user"`
}
```

### Repositorios

#### Repositorio Base
```go
// common/repository/base.go
type BaseRepository[T any] struct {
  db *gorm.DB
}

func NewBaseRepository[T any](db *gorm.DB) *BaseRepository[T] {
  return &BaseRepository[T]{db: db}
}

func (r *BaseRepository[T]) Create(entity *T) error {
  return r.db.Create(entity).Error
}

func (r *BaseRepository[T]) GetByID(id uint) (*T, error) {
  var entity T
  err := r.db.First(&entity, id).Error
  if err != nil {
    return nil, err
  }
  return &entity, nil
}

func (r *BaseRepository[T]) Update(entity *T) error {
  return r.db.Save(entity).Error
}

func (r *BaseRepository[T]) Delete(id uint) error {
  var entity T
  return r.db.Delete(&entity, id).Error
}
```

#### Repositorio Específico
```go
// services/users/internal/repository/user.go
type UserRepository struct {
  *BaseRepository[User]
}

func NewUserRepository(db *gorm.DB) *UserRepository {
  return &UserRepository{
    BaseRepository: NewBaseRepository[User](db),
  }
}

func (r *UserRepository) GetByEmail(email string) (*User, error) {
  var user User
  err := r.db.Where("email = ?", email).First(&user).Error
  if err != nil {
    return nil, err
  }
  return &user, nil
}

func (r *UserRepository) GetActiveUsers() ([]User, error) {
  var users []User
  err := r.db.Where("is_active = ?", true).Find(&users).Error
  return users, err
}
```

### Migraciones

#### Crear Migración
```bash
gox migrate create create_users_table
```

#### Archivo de Migración
```go
// infra/database/migrations/001_create_users_table.go
package migrations

import (
  "gorm.io/gorm"
  "my-app/common/types"
)

func init() {
  migrations = append(migrations, &Migration{
    ID:   "001_create_users_table",
    Up:   up001,
    Down: down001,
  })
}

func up001(db *gorm.DB) error {
  return db.AutoMigrate(
    &types.User{},
    &types.Profile{},
    &types.Post{},
    &types.Tag{},
  )
}

func down001(db *gorm.DB) error {
  return db.Migrator().DropTable(
    &types.User{},
    &types.Profile{},
    &types.Post{},
    &types.Tag{},
  )
}
```

### Uso en Servicios

#### Servicio con Repository
```go
// services/users/internal/service/user.go
type UserService struct {
  repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
  return &UserService{repo: repo}
}

func (s *UserService) CreateUser(req *CreateUserRequest) (*User, error) {
  // Validar datos
  if err := s.validateCreateRequest(req); err != nil {
    return nil, err
  }
  
  // Verificar que no exista el email
  if _, err := s.repo.GetByEmail(req.Email); err == nil {
    return nil, errors.New("email already exists")
  }
  
  // Crear usuario
  user := &User{
    Email:    req.Email,
    Name:     req.Name,
    Password: s.hashPassword(req.Password),
    Role:     "user",
  }
  
  if err := s.repo.Create(user); err != nil {
    return nil, err
  }
  
  return user, nil
}
```

### Sincronización entre Servicios

#### Event Sourcing
```go
// common/events/user.go
type UserCreatedEvent struct {
  UserID uint   `json:"user_id"`
  Email  string `json:"email"`
  Name   string `json:"name"`
}

type UserUpdatedEvent struct {
  UserID uint   `json:"user_id"`
  Email  string `json:"email"`
  Name   string `json:"name"`
}

// Publicar evento
func (s *UserService) CreateUser(req *CreateUserRequest) (*User, error) {
  user, err := s.createUser(req)
  if err != nil {
    return nil, err
  }
  
  // Publicar evento
  event := &UserCreatedEvent{
    UserID: user.ID,
    Email:  user.Email,
    Name:   user.Name,
  }
  
  if err := s.eventBus.Publish("user.created", event); err != nil {
    log.Printf("Failed to publish event: %v", err)
  }
  
  return user, nil
}
```

---

## 📖 Storybook

### Configuración Automática

El Storybook se genera automáticamente cuando ejecutas:

```bash
gox dev --storybook
```

### Auto-discovery

```
src/
├── pages/
│   ├── index.gox           # → Pages/Index
│   ├── about.gox           # → Pages/About
│   └── dashboard/
│       └── index.gox       # → Pages/Dashboard/Index
├── components/
│   ├── button.gox          # → Components/Button
│   ├── card.gox            # → Components/Card
│   └── forms/
│       └── input.gox       # → Components/Forms/Input
└── services/
    └── users/
        └── components/
            └── profile.gox  # → Services/Users/Profile
```

### Configuración por Componente

#### Metadata básica
```gox
<!-- components/button.gox -->
<template>
  <button class="{{.Class}}" {{if .Disabled}}disabled{{end}}>
    {{.Label}}
  </button>
</template>

<go>
type Button struct {
  Label    string `story:"Button Text" default:"Click me"`
  Class    string `story:"CSS Classes" default:"btn-primary"`
  Disabled bool   `story:"Disabled" default:"false"`
}
</go>

<style>
.btn-primary {
  @apply bg-blue-500 text-white px-4 py-2 rounded;
}
</style>

<!-- Configuración del Story -->
<story>
title: "Components/Button"
description: "Basic button component with different styles"

variants:
  - name: "Primary"
    props:
      label: "Primary Button"
      class: "btn-primary"
  
  - name: "Secondary"
    props:
      label: "Secondary Button"
      class: "btn-secondary"
  
  - name: "Disabled"
    props:
      label: "Disabled Button"
      class: "btn-primary"
      disabled: true

controls:
  - name: "label"
    type: "text"
    description: "Button text"
  
  - name: "class"
    type: "select"
    options: ["btn-primary", "btn-secondary", "btn-danger"]
  
  - name: "disabled"
    type: "boolean"
</story>
```

#### Casos de uso complejos
```gox
<!-- components/user-card.gox -->
<template>
  <div class="user-card">
    <img src="{{.User.Avatar}}" alt="{{.User.Name}}">
    <h3>{{.User.Name}}</h3>
    <p>{{.User.Email}}</p>
    <span class="status {{.StatusClass}}">{{.User.Status}}</span>
  </div>
</template>

<go>
type UserCard struct {
  User User
}

type User struct {
  ID     int    `json:"id"`
  Name   string `json:"name"`
  Email  string `json:"email"`
  Avatar string `json:"avatar"`
  Status string `json:"status"`
}

func (u *UserCard) StatusClass() string {
  switch u.User.Status {
  case "active":
    return "status-active"
  case "inactive":
    return "status-inactive"
  default:
    return "status-unknown"
  }
}
</go>

<story>
title: "Components/UserCard"
description: "Display user information in a card format"

fixtures:
  active_user:
    id: 1
    name: "John Doe"
    email: "john@example.com"
    avatar: "https://via.placeholder.com/150"
    status: "active"
  
  inactive_user:
    id: 2
    name: "Jane Smith"
    email: "jane@example.com"
    avatar: "https://via.placeholder.com/150"
    status: "inactive"

variants:
  - name: "Active User"
    props:
      user: "{{.fixtures.active_user}}"
  
  - name: "Inactive User"
    props:
      user: "{{.fixtures.inactive_user}}"
</story>
```

### Funcionalidades del Storybook

#### Navegación
- **Sidebar**: Organizado por carpetas
- **Search**: Búsqueda de componentes
- **Filters**: Filtrar por tipo, tags, etc.

#### Controles
- **Props playground**: Editar props en tiempo real
- **Actions**: Ver eventos disparados
- **Docs**: Documentación auto-generada

#### Testing
- **Visual testing**: Screenshots automáticos
- **Accessibility**: Validaciones A11y
- **Responsive**: Diferentes breakpoints

#### Addons
- **Viewport**: Simular diferentes dispositivos
- **Backgrounds**: Diferentes fondos
- **Themes**: Cambiar entre temas

### Personalización

#### Configuración global
```yaml
# .gox/storybook/config.yaml
title: "My App Design System"
theme: "light"
logo: "./static/logo.png"

addons:
  - viewport
  - backgrounds
  - accessibility
  - docs

backgrounds:
  - name: "Light"
    value: "#ffffff"
  - name: "Dark"
    value: "#1a1a1a"
  - name: "Gray"
    value: "#f5f5f5"

viewports:
  - name: "Mobile"
    width: 375
    height: 667
  - name: "Tablet"
    width: 768
    height: 1024
  - name: "Desktop"
    width: 1440
    height: 900
```

---

## ⚙️ Configuración

### gox.config.yaml

```yaml
# Información del proyecto
name: "my-app"
version: "1.0.0"
description: "Mi aplicación GOX"

# Configuración de desarrollo
dev:
  port: 3000
  host: "localhost"
  hot_reload: true
  storybook: true
  storybook_port: 6006
  verbose: false

# Configuración de producción
prod:
  port: 8080
  host: "0.0.0.0"
  compress: true
  cache_static: true

# Arquitectura
architecture:
  pattern: "distributed"     # distributed, monolithic
  routing: "auto"            # auto, manual

# Protocolos de comunicación
communication:
  protocol: "grpc"           # grpc, rest, graphql
  service_discovery: true
  load_balancer: "round_robin"

# Base de datos
database:
  strategy: "per-service"    # shared, per-service
  default_driver: "postgres"
  
  # Configuración por servicio
  services:
    auth:
      driver: "postgres"
      host: "localhost"
      port: 5432
      name: "auth_db"
      user: "postgres"
      password: "${DB_PASSWORD}"
      ssl_mode: "disable"
    
    users:
      driver: "postgres"
      host: "localhost"
      port: 5432
      name: "users_db"
      user: "postgres"
      password: "${DB_PASSWORD}"
      ssl_mode: "disable"
  
  # Configuración compartida
  shared:
    driver: "postgres"
    host: "localhost"
    port: 5432
    name: "myapp_db"
    user: "postgres"
    password: "${DB_PASSWORD}"
    ssl_mode: "disable"
  
  # Migraciones
  migrations:
    auto: true
    path: "./infra/database/migrations"
    table: "gox_migrations"

# Sincronización entre servicios
sync:
  method: "event-sourcing"   # event-sourcing, cdc, saga
  bus: "nats"               # nats, kafka, rabbitmq
  
  nats:
    url: "nats://localhost:4222"
    subjects:
      - "user.*"
      - "post.*"
      - "notification.*"

# Autenticación
auth:
  provider: "jwt"            # jwt, oauth2, session
  secret: "${JWT_SECRET}"
  expires: "24h"
  refresh: true
  service: "auth"
  
  jwt:
    algorithm: "HS256"
    issuer: "my-app"
    audience: "my-app-users"
  
  oauth2:
    providers:
      google:
        client_id: "${GOOGLE_CLIENT_ID}"
        client_secret: "${GOOGLE_CLIENT_SECRET}"
        scopes: ["email", "profile"]
      
      github:
        client_id: "${GITHUB_CLIENT_ID}"
        client_secret: "${GITHUB_CLIENT_SECRET}"
        scopes: ["user:email"]

# Estilos y frontend
styles:
  framework: "tailwind"      # tailwind, css, scss
  css_modules: true
  postcss: true
  
  tailwind:
    config: "./tailwind.config.js"
    content: ["./pages/**/*.gox", "./components/**/*.gox"]
  
  sass:
    include_paths: ["./styles", "./node_modules"]

# Cache
cache:
  provider: "redis"          # redis, memory, none
  
  redis:
    host: "localhost"
    port: 6379
    password: "${REDIS_PASSWORD}"
    db: 0
  
  memory:
    max_size: "100MB"
    ttl: "1h"

# Logging
logging:
  level: "info"              # debug, info, warn, error
  format: "json"             # json, text
  output: "stdout"           # stdout, file, both
  
  file:
    path: "./logs/app.log"
    max_size: "100MB"
    max_backups: 7
    compress: true

# Monitoring
monitoring:
  enabled: true
  metrics: true
  tracing: true
  
  prometheus:
    port: 9090
    path: "/metrics"
  
  jaeger:
    endpoint: "http://localhost:14268/api/traces"

# Seguridad
security:
  cors:
    enabled: true
    origins: ["http://localhost:3000"]
    methods: ["GET", "POST", "PUT", "DELETE"]
    headers: ["Content-Type", "Authorization"]
  
  rate_limit:
    enabled: true
    requests: 100
    window: "1m"
  
  csrf:
    enabled: true
    secret: "${CSRF_SECRET}"

# Build
build:
  target: "linux/amd64"
  compress: true
  embed_static: true
  
  docker:
    base_image: "alpine:latest"
    port: 8080
    healthcheck: "/health"

# Deploy
deploy:
  targets:
    - name: "staging"
      type: "docker"
      registry: "registry.example.com"
      
    - name: "production"
      type: "kubernetes"
      namespace: "production"
      replicas: 3
      
    - name: "serverless"
      type: "cloud-run"
      region: "us-central1"

# Plugins
plugins:
  - name: "swagger"
    enabled: true
    path: "/docs"
  
  - name: "graphql"
    enabled: false
    path: "/graphql"
  
  - name: "websocket"
    enabled: true
    path: "/ws"

# Variables de entorno
env:
  required:
    - JWT_SECRET
    - DB_PASSWORD
  
  optional:
    - REDIS_PASSWORD
    - GOOGLE_CLIENT_ID
    - GITHUB_CLIENT_ID
```

### Variables de Entorno

```bash
# .env
NODE_ENV=development
PORT=3000

# Database
DB_HOST=localhost
DB_PORT=5432
DB_NAME=myapp
DB_USER=postgres
DB_PASSWORD=secret

# Auth
JWT_SECRET=your-super-secret-key
CSRF_SECRET=another-secret-key

# OAuth
GOOGLE_CLIENT_ID=your-google-client-id
GOOGLE_CLIENT_SECRET=your-google-client-secret
GITHUB_CLIENT_ID=your-github-client-id
GITHUB_CLIENT_SECRET=your-github-client-secret

# Cache
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=redis-password

# Monitoring
JAEGER_ENDPOINT=http://localhost:14268/api/traces
PROMETHEUS_PORT=9090
```

---

## 🎯 Ejemplos

### Ejemplo 1: Blog Simple

#### Estructura
```
blog/
├── pages/
│   ├── index.gox          # Lista de posts
│   ├── post/
│   │   └── [slug].gox     # Post individual
│   └── admin/
│       ├── index.gox      # Admin dashboard
│       └── posts/
│           ├── index.gox  # Lista admin
│           └── edit/
│               └── [id].gox # Editar post
├── services/
│   └── blog/
│       ├── internal/
│       └── api/
└── components/
    ├── post-card.gox
    └── post-form.gox
```

#### Página Principal
```gox
<!-- pages/index.gox -->
<template>
  <div class="container mx-auto px-4 py-8">
    <h1 class="text-4xl font-bold mb-8">Mi Blog</h1>
    
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      {{range .Posts}}
        <gox-component src="post-card" props="{{.}}"/>
      {{end}}
    </div>
    
    <!-- Paginación -->
    <div class="mt-8 flex justify-center">
      {{if .HasPrevious}}
        <a href="?page={{.PreviousPage}}" class="btn btn-secondary mr-2">
          Anterior
        </a>
      {{end}}
      
      {{if .HasNext}}
        <a href="?page={{.NextPage}}" class="btn btn-secondary">
          Siguiente
        </a>
      {{end}}
    </div>
  </div>
</template>

<go>
type IndexPage struct {
  Posts        []Post
  CurrentPage  int
  TotalPages   int
  HasPrevious  bool
  HasNext      bool
  PreviousPage int
  NextPage     int
}

func (p *IndexPage) Mount(ctx *gox.Context) error {
  page := ctx.QueryInt("page", 1)
  
  posts, pagination, err := p.loadPosts(page)
  if err != nil {
    return err
  }
  
  p.Posts = posts
  p.CurrentPage = pagination.Page
  p.TotalPages = pagination.TotalPages
  p.HasPrevious = pagination.HasPrevious
  p.HasNext = pagination.HasNext
  p.PreviousPage = pagination.PreviousPage
  p.NextPage = pagination.NextPage
  
  return nil
}

func (p *IndexPage) loadPosts(page int) ([]Post, *Pagination, error) {
  // Cargar posts desde el servicio
  return blogService.GetPosts(page, 9)
}
</go>
```

#### Componente Post Card
```gox
<!-- components/post-card.gox -->
<template>
  <article class="bg-white rounded-lg shadow-md overflow-hidden">
    <img src="{{.Post.Image}}" alt="{{.Post.Title}}" class="w-full h-48 object-cover">
    
    <div class="p-6">
      <h2 class="text-xl font-bold mb-2">
        <a href="/post/{{.Post.Slug}}" class="hover:text-blue-600">
          {{.Post.Title}}
        </a>
      </h2>
      
      <p class="text-gray-600 mb-4">{{.Post.Excerpt}}</p>
      
      <div class="flex items-center justify-between">
        <div class="flex items-center">
          <img src="{{.Post.Author.Avatar}}" alt="{{.Post.Author.Name}}" 
               class="w-8 h-8 rounded-full mr-2">
          <span class="text-sm text-gray-500">{{.Post.Author.Name}}</span>
        </div>
        
        <time class="text-sm text-gray-500">
          {{.Post.CreatedAt.Format "2006-01-02"}}
        </time>
      </div>
    </div>
  </article>
</template>

<go>
type PostCard struct {
  Post Post
}

type Post struct {
  ID        int       `json:"id"`
  Title     string    `json:"title"`
  Slug      string    `json:"slug"`
  Excerpt   string    `json:"excerpt"`
  Image     string    `json:"image"`
  CreatedAt time.Time `json:"created_at"`
  Author    Author    `json:"author"`
}

type Author struct {
  Name   string `json:"name"`
  Avatar string `json:"avatar"`
}
</go>
```

### Ejemplo 2: Dashboard con HTMX

```gox
<!-- pages/dashboard.gox -->
<template auth="required">
  <div class="min-h-screen bg-gray-100">
    <nav class="bg-white shadow-sm">
      <div class="container mx-auto px-4 py-3">
        <div class="flex justify-between items-center">
          <h1 class="text-xl font-bold">Dashboard</h1>
          <button hx-post="/api/auth/logout" class="btn btn-secondary">
            Logout
          </button>
        </div>
      </div>
    </nav>
    
    <main class="container mx-auto px-4 py-8">
      <!-- Stats Cards -->
      <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
        <div class="bg-white p-6 rounded-lg shadow">
          <h3 class="text-lg font-semibold mb-2">Total Users</h3>
          <p class="text-3xl font-bold text-blue-600" 
             hx-get="/api/stats/users" 
             hx-trigger="load, every 30s">
            {{.Stats.Users}}
          </p>
        </div>
        
        <div class="bg-white p-6 rounded-lg shadow">
          <h3 class="text-lg font-semibold mb-2">Revenue</h3>
          <p class="text-3xl font-bold text-green-600"
             hx-get="/api/stats/revenue" 
             hx-trigger="load, every 30s">
            ${{.Stats.Revenue}}
          </p>
        </div>
        
        <div class="bg-white p-6 rounded-lg shadow">
          <h3 class="text-lg font-semibold mb-2">Orders</h3>
          <p class="text-3xl font-bold text-purple-600"
             hx-get="/api/stats/orders" 
             hx-trigger="load, every 30s">
            {{.Stats.Orders}}
          </p>
        </div>
      </div>
      
      <!-- Chart -->
      <div class="bg-white p-6 rounded-lg shadow mb-8">
        <h3 class="text-lg font-semibold mb-4">Sales Chart</h3>
        <div id="chart" 
             hx-get="/api/chart/sales" 
             hx-trigger="load"
             hx-indicator="#chart-loading">
          <div id="chart-loading" class="text-center">
            Loading chart...
          </div>
        </div>
      </div>
      
      <!-- Recent Activity -->
      <div class="bg-white rounded-lg shadow">
        <div class="p-6 border-b">
          <h3 class="text-lg font-semibold">Recent Activity</h3>
        </div>
        
        <div id="activity-list" 
             hx-get="/api/activity" 
             hx-trigger="load, every 60s">
          <!-- Activity items loaded via HTMX -->
        </div>
      </div>
    </main>
  </div>
</template>

<go>
type DashboardPage struct {
  User  User
  Stats DashboardStats
}

type DashboardStats struct {
  Users   int
  Revenue float64
  Orders  int
}

func (d *DashboardPage) Mount(ctx *gox.Context) error {
  d.User = ctx.Get("user").(User)
  
  // Cargar estadísticas iniciales
  stats, err := d.loadStats()
  if err != nil {
    return err
  }
  d.Stats = stats
  
  return nil
}

func (d *DashboardPage) GetStatsUsers(ctx *gox.Context) error {
  count, err := d.getUserCount()
  if err != nil {
    return err
  }
  
  return ctx.HTML(fmt.Sprintf("%d", count))
}

func (d *DashboardPage) GetStatsRevenue(ctx *gox.Context) error {
  revenue, err := d.getRevenue()
  if err != nil {
    return err
  }
  
  return ctx.HTML(fmt.Sprintf("$%.2f", revenue))
}

func (d *DashboardPage) GetChartSales(ctx *gox.Context) error {
  data, err := d.getSalesData()
  if err != nil {
    return err
  }
  
  // Generar HTML del chart
  html := d.generateChartHTML(data)
  return ctx.HTML(html)
}

func (d *DashboardPage) GetActivity(ctx *gox.Context) error {
  activities, err := d.getRecentActivity()
  if err != nil {
    return err
  }
  
  var html strings.Builder
  for _, activity := range activities {
    html.WriteString(fmt.Sprintf(`
      <div class="p-4 border-b">
        <p class="font-medium">%s</p>
        <p class="text-sm text-gray-500">%s</p>
      </div>
    `, activity.Description, activity.CreatedAt.Format("2006-01-02 15:04")))
  }
  
  return ctx.HTML(html.String())
}
</go>
```

### Ejemplo 3: API REST con Servicios

```go
// services/users/internal/handlers/user.go
package handlers

type UserHandler struct {
  service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
  return &UserHandler{service: service}
}

func (h *UserHandler) GetUsers(ctx *gox.Context) error {
  page := ctx.QueryInt("page", 1)
  limit := ctx.QueryInt("limit", 10)
  
  users, pagination, err := h.service.GetUsers(page, limit)
  if err != nil {
    return ctx.InternalError("Failed to get users")
  }
  
  return ctx.JSON(gox.Response{
    Data:       users,
    Pagination: pagination,
  })
}

func (h *UserHandler) GetUser(ctx *gox.Context) error {
  id := ctx.ParamInt("id")
  
  user, err := h.service.GetUser(id)
  if err != nil {
    return ctx.NotFound("User not found")
  }
  
  return ctx.JSON(user)
}

func (h *UserHandler) CreateUser(ctx *gox.Context) error {
  var req CreateUserRequest
  if err := ctx.Bind(&req); err != nil {
    return ctx.BadRequest("Invalid request")
  }
  
  user, err := h.service.CreateUser(&req)
  if err != nil {
    return ctx.InternalError("Failed to create user")
  }
  
  return ctx.Created(user)
}

func (h *UserHandler) UpdateUser(ctx *gox.Context) error {
  id := ctx.ParamInt("id")
  
  var req UpdateUserRequest
  if err := ctx.Bind(&req); err != nil {
    return ctx.BadRequest("Invalid request")
  }
  
  user, err := h.service.UpdateUser(id, &req)
  if err != nil {
    return ctx.InternalError("Failed to update user")
  }
  
  return ctx.JSON(user)
}

func (h *UserHandler) DeleteUser(ctx *gox.Context) error {
  id := ctx.ParamInt("id")
  
  if err := h.service.DeleteUser(id); err != nil {
    return ctx.InternalError("Failed to delete user")
  }
  
  return ctx.NoContent()
}
```

Estos ejemplos muestran cómo usar GOX Framework para crear aplicaciones web modernas con una arquitectura limpia y escalable.