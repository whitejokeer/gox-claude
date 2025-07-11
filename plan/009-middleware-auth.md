# Task 009: Middleware y Autenticación

## Descripción
Implementar un sistema robusto de middleware con autenticación integrada, soportando JWT, OAuth2, y sesiones, con helpers para proteger rutas y manejar permisos.

## Prioridad
Alta

## Estimación
5-6 días

## Dependencias
- Task 005: Sistema de routing
- Task 008: Sistema de componentes

## Subtasks

### 9.1 Sistema de middleware base
- [ ] Implementar cadena de middleware
- [ ] Middleware global vs por ruta
- [ ] Sistema de prioridades
- [ ] Error handling en middleware
- [ ] Context propagation

### 9.2 Autenticación JWT
- [ ] Generación de tokens JWT
- [ ] Validación de tokens
- [ ] Refresh tokens
- [ ] Blacklist de tokens
- [ ] Claims personalizados

### 9.3 Autenticación OAuth2
- [ ] Implementar flujo OAuth2
- [ ] Soportar múltiples providers (Google, GitHub, etc.)
- [ ] Callback handling
- [ ] State management
- [ ] Token storage

### 9.4 Sistema de sesiones
- [ ] Session store (memoria, Redis, DB)
- [ ] Cookie management
- [ ] Session lifecycle
- [ ] CSRF protection
- [ ] Session hijacking prevention

### 9.5 Autorización y permisos
- [ ] Sistema de roles
- [ ] Permisos granulares
- [ ] Políticas de acceso
- [ ] Middleware de autorización
- [ ] Helpers para templates

### 9.6 Middleware comunes
- [ ] CORS middleware
- [ ] Rate limiting
- [ ] Request logging
- [ ] Compression
- [ ] Security headers

## Criterios de Aceptación

1. **Middleware chain funcionando**
   ```go
   router := gox.NewRouter()
   
   // Global middleware
   router.Use(middleware.Logger())
   router.Use(middleware.CORS())
   router.Use(middleware.RateLimit(100, time.Minute))
   
   // Grupo con auth
   protected := router.Group("/api", middleware.JWT())
   protected.GET("/profile", handlers.Profile)
   
   // Ruta con múltiples middleware
   router.POST("/admin", 
       middleware.JWT(),
       middleware.RequireRole("admin"),
       handlers.AdminPanel,
   )
   ```

2. **Autenticación JWT**
   ```go
   // Generar token
   token, err := auth.GenerateJWT(user, auth.JWTOptions{
       ExpiresIn: 24 * time.Hour,
       Claims: map[string]interface{}{
           "role": user.Role,
           "permissions": user.Permissions,
       },
   })
   
   // Middleware JWT
   func JWT() gox.Middleware {
       return func(ctx *gox.Context) error {
           token := ctx.GetHeader("Authorization")
           
           user, err := auth.ValidateJWT(token)
           if err != nil {
               return ctx.Unauthorized("Invalid token")
           }
           
           ctx.Set("user", user)
           return ctx.Next()
       }
   }
   ```

3. **OAuth2 flow**
   ```go
   // Configurar OAuth2
   auth.ConfigureOAuth2(auth.OAuth2Config{
       Providers: map[string]auth.Provider{
           "google": {
               ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
               ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
               Scopes:       []string{"email", "profile"},
           },
       },
   })
   
   // Routes
   router.GET("/auth/google", auth.OAuth2Login("google"))
   router.GET("/auth/google/callback", auth.OAuth2Callback("google"))
   ```

4. **Protección en templates**
   ```gox
   <template auth="required">
     <!-- Solo usuarios autenticados -->
   </template>
   
   <template auth="role:admin">
     <!-- Solo administradores -->
   </template>
   
   <template>
     {{if .User.HasPermission "posts.edit"}}
       <button>Edit Post</button>
     {{end}}
   </template>
   ```

## Tests Necesarios

### Tests Unitarios

1. **Test middleware chain**
```go
func TestMiddlewareChain(t *testing.T) {
    var executed []string
    
    m1 := func(ctx *gox.Context) error {
        executed = append(executed, "m1-before")
        err := ctx.Next()
        executed = append(executed, "m1-after")
        return err
    }
    
    m2 := func(ctx *gox.Context) error {
        executed = append(executed, "m2")
        return ctx.Next()
    }
    
    handler := func(ctx *gox.Context) error {
        executed = append(executed, "handler")
        return nil
    }
    
    chain := gox.Chain(m1, m2, handler)
    chain(&gox.Context{})
    
    expected := []string{"m1-before", "m2", "handler", "m1-after"}
    assert.Equal(t, expected, executed)
}
```

2. **Test JWT**
```go
func TestJWTAuthentication(t *testing.T) {
    user := &User{ID: 1, Email: "test@example.com"}
    
    // Generar token
    token, err := auth.GenerateJWT(user, auth.JWTOptions{
        Secret: "test-secret",
        ExpiresIn: time.Hour,
    })
    assert.NoError(t, err)
    
    // Validar token
    decoded, err := auth.ValidateJWT(token, "test-secret")
    assert.NoError(t, err)
    assert.Equal(t, user.ID, decoded.ID)
    
    // Token expirado
    expiredToken := generateExpiredToken()
    _, err = auth.ValidateJWT(expiredToken, "test-secret")
    assert.Error(t, err)
}
```

3. **Test autorización**
```go
func TestAuthorization(t *testing.T) {
    tests := []struct {
        user     *User
        required string
        allowed  bool
    }{
        {
            user:     &User{Role: "admin"},
            required: "admin",
            allowed:  true,
        },
        {
            user:     &User{Role: "user"},
            required: "admin",
            allowed:  false,
        },
        {
            user:     &User{Permissions: []string{"posts.edit"}},
            required: "permission:posts.edit",
            allowed:  true,
        },
    }
    
    for _, tt := range tests {
        middleware := RequireAuth(tt.required)
        ctx := &gox.Context{}
        ctx.Set("user", tt.user)
        
        err := middleware(ctx)
        if tt.allowed {
            assert.NoError(t, err)
        } else {
            assert.Error(t, err)
        }
    }
}
```

### Tests de Integración

1. **Test flujo completo de auth**
```go
func TestAuthenticationFlow(t *testing.T) {
    app := createTestApp()
    
    // Login
    loginReq := map[string]string{
        "email": "user@example.com",
        "password": "password123",
    }
    
    resp := httptest.NewRecorder()
    req := httptest.NewRequest("POST", "/auth/login", toJSON(loginReq))
    app.ServeHTTP(resp, req)
    
    assert.Equal(t, 200, resp.Code)
    
    var loginResp map[string]interface{}
    json.Unmarshal(resp.Body.Bytes(), &loginResp)
    token := loginResp["token"].(string)
    
    // Acceder a ruta protegida
    resp2 := httptest.NewRecorder()
    req2 := httptest.NewRequest("GET", "/api/profile", nil)
    req2.Header.Set("Authorization", "Bearer "+token)
    app.ServeHTTP(resp2, req2)
    
    assert.Equal(t, 200, resp2.Code)
}
```

2. **Test OAuth2**
```go
func TestOAuth2Flow(t *testing.T) {
    app := createTestApp()
    
    // Iniciar OAuth2
    resp := httptest.NewRecorder()
    req := httptest.NewRequest("GET", "/auth/google", nil)
    app.ServeHTTP(resp, req)
    
    // Debe redirigir a Google
    assert.Equal(t, 302, resp.Code)
    location := resp.Header().Get("Location")
    assert.Contains(t, location, "accounts.google.com")
    
    // Simular callback
    state := extractState(location)
    code := "test-code"
    
    resp2 := httptest.NewRecorder()
    req2 := httptest.NewRequest("GET", "/auth/google/callback?state="+state+"&code="+code, nil)
    app.ServeHTTP(resp2, req2)
    
    // Debe crear sesión
    assert.Equal(t, 302, resp2.Code)
    assert.Contains(t, resp2.Header().Get("Set-Cookie"), "session=")
}
```

### Tests de Performance

```go
func BenchmarkJWTValidation(b *testing.B) {
    token, _ := auth.GenerateJWT(&User{ID: 1}, auth.JWTOptions{
        Secret: "test-secret",
    })
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        auth.ValidateJWT(token, "test-secret")
    }
}

func BenchmarkMiddlewareChain(b *testing.B) {
    middleware := []gox.Middleware{
        Logger(),
        CORS(),
        RateLimit(1000, time.Minute),
        JWT(),
    }
    
    handler := func(ctx *gox.Context) error {
        return ctx.Text("OK")
    }
    
    chain := gox.Chain(append(middleware, handler)...)
    ctx := &gox.Context{}
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        chain(ctx)
    }
}
```

## Definición de Done

- [ ] Sistema de middleware completo
- [ ] JWT authentication funcionando
- [ ] OAuth2 con múltiples providers
- [ ] Sistema de sesiones robusto
- [ ] Autorización flexible
- [ ] Middleware comunes implementados
- [ ] Tests con cobertura > 85%
- [ ] Documentación completa

## Notas Adicionales

- La seguridad debe ser prioridad
- Seguir mejores prácticas de OWASP
- Los tokens deben ser seguros y eficientes
- Considerar rate limiting por usuario
- Implementar audit logging
- Documentar configuración de seguridad