# GOX Framework - Product Requirements Document

## 🎯 Visión del Producto

**GOX** es un framework de desarrollo web moderno que combina el poder de Go con la simplicidad de HTMX para crear aplicaciones web ultrarrápidas mediante Server-Side Rendering. Inspirado en Vue.js SFC, React y Svelte, pero diseñado para ser "fácil de aprender, difícil de masterizar".

### Filosofía Core
- **Performance First**: SSR nativo + Go + HTMX = velocidad extrema
- **Developer Experience**: Herramientas integradas que facilitan el desarrollo
- **Distributed by Default**: Arquitectura distribuida desde el día uno
- **Zero Client Bundle**: Sin JavaScript cliente pesado

---

## 🏗️ Arquitectura del Framework

### Estructura de Proyecto
```
my-app/
├── pages/              # Frontend (.gox pages) - UI Layer
│   ├── index.gox
│   ├── dashboard/
│   │   └── index.gox
│   └── _layout.gox    # Layout compartido
│
├── services/           # Backend services/modules - Business Logic
│   ├── auth/          # Servicio puro (Go)
│   │   ├── cmd/
│   │   ├── internal/
│   │   ├── api/
│   │   └── service.yaml
│   │
│   └── users/         # Módulo con componentes
│       ├── cmd/
│       ├── internal/
│       ├── components/   # Solo si es módulo
│       │   └── user-card.gox
│       └── module.yaml
│
├── routing/            # Configuración de rutas - Presentation Layer
│   ├── routes.go      # Si prefiere manual
│   └── middleware.go  # Middleware global
│
├── common/            # Código compartido - Shared Layer
│   ├── types/
│   ├── utils/
│   └── contracts/     # Interfaces entre servicios
│
├── infra/             # Infraestructura - Infrastructure Layer
│   ├── database/
│   ├── cache/
│   └── messaging/     # Event bus, queues
│
└── gox.config.yaml    # Configuración principal
```

---

## 📋 Especificaciones Técnicas

### 1. Formato de Archivo .gox

Los archivos .gox son componentes de archivo único inspirados en Vue SFC:

```gox
<template>
  <div class="bg-blue-500 p-4" hx-get="/api/stats" hx-target="#stats">
    <h1>{{.Title}}</h1>
    <div id="stats">{{.Stats}}</div>
  </div>
</template>

<go>
type DashboardData struct {
  Title string
  Stats interface{}
}

func (d *Dashboard) Mount(ctx *gox.Context) error {
  d.Title = "Mi Dashboard"
  return nil
}

func (d *Dashboard) GetStats(ctx *gox.Context) error {
  // Lógica del handler HTMX
  return ctx.JSON(map[string]interface{}{
    "users": 150,
    "sales": 1200,
  })
}
</go>

<style scoped>
.dashboard {
  padding: 2rem;
  background: var(--bg-primary);
}
</style>
```

**Características:**
- `<template>`: HTML con sintaxis Go template + HTMX
- `<go>`: Lógica del componente en Go puro
- `<style>`: CSS con scoped opcional + soporte Tailwind
- Todas las secciones son opcionales

### 2. Sistema de Routing

#### Automático (File-based)
```
pages/
├── index.gox          → /
├── about.gox          → /about
├── users/
│   ├── index.gox      → /users
│   └── [id].gox       → /users/:id
└── api/
    └── [...].gox      → /api/* (catch-all)
```

#### Manual (Configurable)
```go
// routing/routes.go
func SetupRoutes(r *gox.Router) {
    r.Use(middleware.Auth())
    
    r.Group("/api", func(g *gox.Group) {
        g.Use(middleware.RateLimit())
        g.Mount("/users", services.Users)
        g.Mount("/auth", services.Auth)
    })
    
    r.Page("/dashboard", pages.Dashboard, middleware.RequireAuth())
}
```

### 3. Configuración Principal

```yaml
# gox.config.yaml
name: "my-app"
version: "1.0.0"

# Arquitectura
architecture:
  pattern: "distributed"  # distributed, monolithic
  routing: "auto"         # auto, manual
  
# Protocolo de comunicación entre servicios
communication:
  protocol: "grpc"        # grpc, rest, graphql
  service_discovery: true
  
# Base de datos
database:
  strategy: "per-service" # shared, per-service
  default: "postgres"
  migrations: "auto"      # Sincronización automática
  
# Autenticación
auth:
  provider: "jwt"         # jwt, oauth2, session
  service: "auth"         # Servicio que maneja auth
  
# Desarrollo
dev:
  hot_reload: true
  storybook: true
  port: 3000
  
# Estilos
styles:
  framework: "tailwind"   # tailwind, css, scss
  css_modules: true
```

---

## 🛠️ CLI Commands

### Comandos Principales

#### Creación de Proyectos
```bash
# Crear nuevo proyecto completo
gox new project my-app

# Crear módulo independiente
gox new module user-management

# Crear servicio puro
gox new service auth-service
```

#### Generación de Componentes
```bash
# Generar página
gox generate page dashboard

# Generar componente
gox generate component user-card

# Generar servicio
gox generate service users
```

#### Desarrollo
```bash
# Servidor de desarrollo
gox dev

# Servidor de desarrollo con storybook
gox dev --storybook

# Build para producción
gox build

# Deploy
gox deploy --target=docker
```

#### Gestión de Dependencias
```bash
# Instalar dependencias
gox install

# Agregar servicio
gox add service auth

# Remover servicio
gox remove service auth
```

#### Utilidades
```bash
# Generar tipos TypeScript para frontend
gox generate types

# Ejecutar migraciones
gox migrate

# Ejecutar tests
gox test

# Linting
gox lint
```

---

## 🎨 Storybook Integrado

### Auto-Discovery
El storybook se genera automáticamente basado en la estructura del proyecto:

```
src/
├── pages/
│   └── dashboard.gox     # → Storybook: Pages/Dashboard
├── components/
│   └── user-card.gox     # → Storybook: Components/UserCard
└── services/
    └── users/
        └── components/
            └── profile.gox # → Storybook: Services/Users/Profile
```

### Características
- **Hot Reload**: Cambios en tiempo real
- **Props Playground**: Generado automáticamente de los structs Go
- **Casos de Uso**: Basados en los handlers definidos
- **Responsive Testing**: Diferentes breakpoints
- **Accessibility Testing**: Validación automática

---

## 🔐 Middleware y Autenticación

### Middleware Global
```go
// routing/middleware.go
func Auth() gox.Middleware {
    return func(ctx *gox.Context) error {
        token := ctx.Header("Authorization")
        user, err := auth.Verify(token)
        if err != nil {
            return ctx.Unauthorized()
        }
        ctx.Set("user", user)
        return ctx.Next()
    }
}
```

### Autenticación en Componentes
```gox
<template auth="required">
  <div>
    <h1>Welcome {{.User.Name}}</h1>
  </div>
</template>

<go>
func (p *Dashboard) BeforeMount(ctx *gox.Context) error {
    p.User = ctx.Get("user").(*User)
    return nil
}
</go>
```

---

## 🗄️ Estrategias de Base de Datos

### Opción 1: Base de Datos Compartida
```yaml
database:
  strategy: "shared"
  migrations:
    tool: "migrate"
    auto_sync: true
```

### Opción 2: Base de Datos por Servicio
```yaml
database:
  strategy: "per-service"
  sync:
    method: "event-sourcing"  # event-sourcing, cdc, saga
    bus: "nats"               # nats, kafka, rabbitmq
```

---

## 📡 Comunicación Entre Servicios

### Contracts
```go
// common/contracts/users.go
type UserService interface {
    GetUser(ctx context.Context, id string) (*User, error)
    ListUsers(ctx context.Context, filters Filters) ([]*User, error)
}
```

### Cliente Auto-generado
```go
// En cualquier servicio
users := gox.Service("users").Client()
user, err := users.GetUser(ctx, "123")
```

---

## 🎯 MVP Features

### Core (Esencial)
- ✅ Parser .gox básico
- ✅ Compilador Go handlers
- ✅ Router file-based y manual
- ✅ Servidor de desarrollo con hot reload
- ✅ HTMX integration
- ✅ Storybook integrado
- ✅ CLI completo
- ✅ Middleware y Auth básico
- ✅ Soporte Tailwind CSS

### Arquitectura
- ✅ Estructura distribuida
- ✅ Comunicación gRPC/REST
- ✅ Base de datos por servicio
- ✅ Event sourcing básico

---

## 🚀 Roadmap Post-MVP

### v1.1 - Mejoras de DX
- 🔄 Agente IA integrado (GenKit)
- 🔄 Testing framework integrado
- 🔄 Observabilidad built-in
- 🔄 Performance monitoring

### v1.2 - Integraciones
- 🔄 Embeddings de terceros (React/Vue)
- 🔄 Theming avanzado
- 🔄 Edge computing support
- 🔄 Multi-database support

### v1.3 - Enterprise
- 🔄 Kubernetes deployment
- 🔄 Microservices orchestration
- 🔄 Advanced security features
- 🔄 Multi-tenant support

---

## 📊 Métricas de Éxito

### Performance
- **Time to First Byte**: < 50ms
- **First Contentful Paint**: < 100ms
- **Bundle Size**: 0KB cliente (Solo HTML/CSS)
- **Memory Usage**: < 50MB por servicio

### Developer Experience
- **Setup Time**: < 5 minutos
- **Hot Reload**: < 500ms
- **Build Time**: < 30 segundos
- **Learning Curve**: 1 semana para productividad

### Adoption
- **Community**: 1000+ GitHub stars en 6 meses
- **Ecosystem**: 50+ componentes community
- **Enterprise**: 5+ empresas usando en producción

---

## 🎯 Propuesta de Valor

### Para Desarrolladores
- **Simplicidad**: Un solo lenguaje (Go) para frontend y backend
- **Performance**: Aplicaciones ultra-rápidas sin optimizaciones complejas
- **Productividad**: Herramientas integradas que aceleran el desarrollo
- **Escalabilidad**: Arquitectura distribuida desde el día uno

### Para Empresas
- **Costo**: Menor infrastructure footprint
- **Velocidad**: Time to market reducido
- **Mantenibilidad**: Código más simple y predecible
- **Seguridad**: Beneficios inherentes de Go

---

## 🏁 Conclusión

GOX Framework representa una evolución natural del desarrollo web, combinando las mejores prácticas de frameworks modernos con la simplicidad y performance de Go. Su arquitectura distribuida y herramientas integradas lo posicionan como la opción ideal para equipos que buscan velocidad, simplicidad y escalabilidad.

El framework está diseñado para ser **fácil de aprender pero difícil de masterizar**, siguiendo la filosofía de Go, y ofrece una alternativa compelling a los frameworks JavaScript tradicionales para aplicaciones web de alta performance.