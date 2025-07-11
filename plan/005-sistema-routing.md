# Task 005: Sistema de Routing

## Descripción
Implementar el sistema de routing dual que soporte tanto routing automático basado en archivos como routing manual configurable, con soporte para middleware, grupos, y parámetros dinámicos.

## Prioridad
Alta

## Estimación
4-5 días

## Dependencias
- Task 003: Parser de archivos .gox
- Task 004: Compilador .gox a Go

## Subtasks

### 5.1 Router core
- [ ] Implementar router base con árbol radix
- [ ] Soportar métodos HTTP (GET, POST, PUT, DELETE, etc.)
- [ ] Implementar sistema de middleware
- [ ] Crear context personalizado para GOX
- [ ] Manejar parámetros de ruta y query

### 5.2 Routing automático (file-based)
- [ ] Escanear directorio `pages/` recursivamente
- [ ] Mapear estructura de archivos a rutas
- [ ] Soportar rutas dinámicas ([id], [slug])
- [ ] Implementar catch-all routes ([...])
- [ ] Generar mapa de rutas automáticamente

### 5.3 Routing manual
- [ ] Crear API para definir rutas manualmente
- [ ] Implementar grupos de rutas
- [ ] Soportar middleware por ruta/grupo
- [ ] Permitir mount de servicios
- [ ] Integrar con archivo routes.go

### 5.4 Parámetros y validación
- [ ] Extraer parámetros de ruta
- [ ] Parsear query parameters
- [ ] Validación automática de tipos
- [ ] Binding de JSON/Form data
- [ ] Manejo de uploads de archivos

### 5.5 Middleware system
- [ ] Implementar cadena de middleware
- [ ] Middleware globales
- [ ] Middleware por grupo
- [ ] Middleware por ruta
- [ ] Orden de ejecución correcto

### 5.6 Integración con componentes
- [ ] Conectar rutas con páginas .gox
- [ ] Resolver handlers HTMX
- [ ] Manejar layouts
- [ ] Implementar error pages
- [ ] Soportar redirects

## Criterios de Aceptación

1. **Router funcional**
   ```go
   router := gox.NewRouter()
   
   // Routing manual
   router.GET("/", handlers.Home)
   router.POST("/api/users", handlers.CreateUser)
   
   // Grupos
   api := router.Group("/api", middleware.Auth())
   api.GET("/profile", handlers.Profile)
   
   // Parámetros
   router.GET("/users/:id", handlers.GetUser)
   router.GET("/posts/*slug", handlers.GetPost)
   ```

2. **Routing automático**
   ```
   pages/
   ├── index.gox          → GET /
   ├── about.gox          → GET /about
   ├── users/
   │   ├── index.gox      → GET /users
   │   └── [id].gox       → GET /users/:id
   └── blog/
       └── [...slug].gox  → GET /blog/*slug
   ```

3. **Context GOX**
   ```go
   func Handler(ctx *gox.Context) error {
       // Params
       id := ctx.Param("id")
       page := ctx.QueryInt("page", 1)
       
       // Binding
       var user User
       ctx.Bind(&user)
       
       // Response
       return ctx.JSON(user)
       return ctx.HTML("<div>Hello</div>")
       return ctx.Redirect("/login")
   }
   ```

4. **Middleware**
   ```go
   func Logger() gox.Middleware {
       return func(ctx *gox.Context) error {
           start := time.Now()
           err := ctx.Next()
           log.Printf("%s %s %v", ctx.Method(), ctx.Path(), time.Since(start))
           return err
       }
   }
   ```

## Tests Necesarios

### Tests Unitarios

1. **Test router básico**
```go
func TestBasicRouting(t *testing.T) {
    router := gox.NewRouter()
    
    router.GET("/hello", func(ctx *gox.Context) error {
        return ctx.Text("Hello World")
    })
    
    req := httptest.NewRequest("GET", "/hello", nil)
    w := httptest.NewRecorder()
    
    router.ServeHTTP(w, req)
    
    assert.Equal(t, 200, w.Code)
    assert.Equal(t, "Hello World", w.Body.String())
}
```

2. **Test parámetros**
```go
func TestRouteParams(t *testing.T) {
    router := gox.NewRouter()
    
    router.GET("/users/:id", func(ctx *gox.Context) error {
        id := ctx.Param("id")
        return ctx.JSON(map[string]string{"id": id})
    })
    
    req := httptest.NewRequest("GET", "/users/123", nil)
    w := httptest.NewRecorder()
    
    router.ServeHTTP(w, req)
    
    assert.Equal(t, 200, w.Code)
    assert.Contains(t, w.Body.String(), `"id":"123"`)
}
```

3. **Test routing automático**
```go
func TestFileBasedRouting(t *testing.T) {
    router := gox.NewRouter()
    router.EnableAutoRouting("testdata/pages")
    
    routes := router.GetRoutes()
    
    expected := map[string]string{
        "GET /":           "index.gox",
        "GET /about":      "about.gox",
        "GET /users/:id":  "users/[id].gox",
    }
    
    for route, file := range expected {
        assert.Contains(t, routes, route)
        assert.Equal(t, file, routes[route].File)
    }
}
```

### Tests de Integración

1. **Test middleware chain**
```go
func TestMiddlewareExecution(t *testing.T) {
    var order []string
    
    middleware1 := func(ctx *gox.Context) error {
        order = append(order, "m1-before")
        err := ctx.Next()
        order = append(order, "m1-after")
        return err
    }
    
    middleware2 := func(ctx *gox.Context) error {
        order = append(order, "m2-before")
        err := ctx.Next()
        order = append(order, "m2-after")
        return err
    }
    
    handler := func(ctx *gox.Context) error {
        order = append(order, "handler")
        return nil
    }
    
    router := gox.NewRouter()
    router.Use(middleware1)
    router.GET("/test", handler, middleware2)
    
    req := httptest.NewRequest("GET", "/test", nil)
    router.ServeHTTP(httptest.NewRecorder(), req)
    
    expected := []string{"m1-before", "m2-before", "handler", "m2-after", "m1-after"}
    assert.Equal(t, expected, order)
}
```

2. **Test error handling**
```go
func TestErrorHandling(t *testing.T) {
    router := gox.NewRouter()
    
    router.GET("/error", func(ctx *gox.Context) error {
        return ctx.Error(500, "Internal Server Error")
    })
    
    router.OnError(func(ctx *gox.Context, err error) {
        ctx.JSON(map[string]string{"error": err.Error()})
    })
    
    req := httptest.NewRequest("GET", "/error", nil)
    w := httptest.NewRecorder()
    
    router.ServeHTTP(w, req)
    
    assert.Equal(t, 500, w.Code)
    assert.Contains(t, w.Body.String(), "Internal Server Error")
}
```

### Benchmarks

```go
func BenchmarkRouter(b *testing.B) {
    router := gox.NewRouter()
    
    for i := 0; i < 100; i++ {
        path := fmt.Sprintf("/route%d", i)
        router.GET(path, func(ctx *gox.Context) error {
            return ctx.Text("OK")
        })
    }
    
    req := httptest.NewRequest("GET", "/route50", nil)
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        w := httptest.NewRecorder()
        router.ServeHTTP(w, req)
    }
}
```

## Definición de Done

- [ ] Router con árbol radix eficiente
- [ ] Routing automático funcionando
- [ ] Routing manual completo
- [ ] Sistema de middleware robusto
- [ ] Context con helpers útiles
- [ ] Tests con cobertura > 90%
- [ ] Benchmarks mostrando alta performance
- [ ] Documentación con ejemplos

## Notas Adicionales

- Considerar usar chi o httprouter como base
- El router debe ser compatible con net/http
- Mantener API similar a frameworks populares
- Optimizar para rutas estáticas comunes
- Pensar en debugging y logging de rutas
- Soportar websockets en el futuro