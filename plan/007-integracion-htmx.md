# Task 007: Integración HTMX

## Descripción
Implementar la integración completa con HTMX para permitir interactividad sin JavaScript del lado del cliente, incluyendo helpers, atributos personalizados, y optimizaciones para el servidor.

## Prioridad
Alta

## Estimación
3-4 días

## Dependencias
- Task 004: Compilador .gox a Go
- Task 005: Sistema de routing

## Subtasks

### 7.1 Integración base de HTMX
- [ ] Incluir HTMX en templates automáticamente
- [ ] Configurar versión de HTMX a usar
- [ ] Implementar CDN fallback
- [ ] Soportar modo offline
- [ ] Minificar para producción

### 7.2 Procesamiento de atributos HTMX
- [ ] Detectar atributos hx-* en parser
- [ ] Validar atributos HTMX válidos
- [ ] Generar endpoints automáticamente
- [ ] Mapear a handlers Go
- [ ] Preservar atributos en output

### 7.3 Helpers del lado del servidor
- [ ] Implementar respuestas HTMX (HX-Trigger, HX-Redirect, etc.)
- [ ] Helper para swap strategies
- [ ] Manejo de eventos server-sent
- [ ] Helpers para polling
- [ ] Integración con WebSockets

### 7.4 Optimizaciones de respuesta
- [ ] Detectar requests HTMX (HX-Request header)
- [ ] Enviar solo fragmentos HTML cuando apropiado
- [ ] Comprimir respuestas automáticamente
- [ ] Cache de fragmentos
- [ ] Respuestas streaming

### 7.5 Extensiones y plugins
- [ ] Soportar extensiones HTMX comunes
- [ ] Sistema de plugins personalizado
- [ ] Integración con Alpine.js (opcional)
- [ ] Debug mode para HTMX
- [ ] Métricas de performance

### 7.6 Componentes HTMX nativos
- [ ] Componente de formulario con validación
- [ ] Componente de tabla con paginación
- [ ] Componente de búsqueda en tiempo real
- [ ] Componente de infinite scroll
- [ ] Modal/Dialog helpers

## Criterios de Aceptación

1. **Atributos HTMX funcionando**
   ```gox
   <template>
     <div>
       <button hx-get="/api/users" 
               hx-target="#users-list"
               hx-swap="innerHTML">
         Load Users
       </button>
       <div id="users-list"></div>
     </div>
   </template>
   ```

2. **Helpers del servidor**
   ```go
   func (h *Handler) GetUsers(ctx *gox.Context) error {
       // Detectar si es request HTMX
       if ctx.IsHTMX() {
           // Enviar solo el fragmento
           return ctx.HTMLFragment(`
               <div class="user">John Doe</div>
               <div class="user">Jane Smith</div>
           `)
       }
       
       // Enviar página completa
       return ctx.Render("users", users)
   }
   
   // Triggers y eventos
   func (h *Handler) DeleteUser(ctx *gox.Context) error {
       // Eliminar usuario...
       
       return ctx.HTMXTrigger("user-deleted").
           HTMXRedirect("/users").
           NoContent()
   }
   ```

3. **Componentes nativos**
   ```gox
   <!-- Búsqueda en tiempo real -->
   <gox-search 
     hx-get="/api/search"
     hx-trigger="keyup changed delay:500ms"
     hx-target="#results">
   </gox-search>
   
   <!-- Tabla con paginación -->
   <gox-table 
     hx-get="/api/users"
     hx-include="#filters"
     paginated="true">
   </gox-table>
   ```

4. **Optimizaciones automáticas**
   - Fragmentos HTML cuando HX-Request presente
   - Compresión gzip automática
   - Headers HTMX correctos
   - Cache inteligente

## Tests Necesarios

### Tests Unitarios

1. **Test detección de atributos**
```go
func TestHTMXAttributeDetection(t *testing.T) {
    template := `
    <button hx-get="/test" hx-swap="outerHTML">Test</button>
    `
    
    attrs := parser.ExtractHTMXAttributes(template)
    
    assert.Contains(t, attrs, HTMXAttribute{
        Type: "hx-get",
        Value: "/test",
    })
    assert.Contains(t, attrs, HTMXAttribute{
        Type: "hx-swap",
        Value: "outerHTML",
    })
}
```

2. **Test helpers del servidor**
```go
func TestHTMXHelpers(t *testing.T) {
    ctx := &gox.Context{
        Request: &http.Request{
            Header: http.Header{
                "HX-Request": []string{"true"},
            },
        },
    }
    
    assert.True(t, ctx.IsHTMX())
    
    // Test trigger
    ctx.HTMXTrigger("test-event", map[string]interface{}{
        "id": 123,
    })
    
    assert.Equal(t, `{"test-event":{"id":123}}`, ctx.Header("HX-Trigger"))
}
```

3. **Test optimización de respuesta**
```go
func TestHTMXResponseOptimization(t *testing.T) {
    handler := func(ctx *gox.Context) error {
        users := []User{{Name: "John"}, {Name: "Jane"}}
        
        if ctx.IsHTMX() {
            return ctx.RenderPartial("users-list", users)
        }
        return ctx.Render("users-page", users)
    }
    
    // Request normal
    req1 := httptest.NewRequest("GET", "/users", nil)
    resp1 := executeHandler(handler, req1)
    assert.Contains(t, resp1.Body, "<html>")
    
    // Request HTMX
    req2 := httptest.NewRequest("GET", "/users", nil)
    req2.Header.Set("HX-Request", "true")
    resp2 := executeHandler(handler, req2)
    assert.NotContains(t, resp2.Body, "<html>")
}
```

### Tests de Integración

1. **Test flujo completo HTMX**
```go
func TestHTMXFlow(t *testing.T) {
    app := createTestApp()
    
    // Página inicial
    resp1 := httptest.NewRecorder()
    req1 := httptest.NewRequest("GET", "/", nil)
    app.ServeHTTP(resp1, req1)
    
    // Verificar que incluye HTMX
    assert.Contains(t, resp1.Body.String(), "htmx.org")
    
    // Request HTMX
    resp2 := httptest.NewRecorder()
    req2 := httptest.NewRequest("GET", "/api/data", nil)
    req2.Header.Set("HX-Request", "true")
    req2.Header.Set("HX-Target", "content")
    app.ServeHTTP(resp2, req2)
    
    // Verificar headers HTMX
    assert.NotEmpty(t, resp2.Header().Get("HX-Push"))
}
```

2. **Test componentes HTMX**
```go
func TestHTMXComponents(t *testing.T) {
    // Test búsqueda
    comp := ParseComponent(`<gox-search hx-get="/search" />`)
    html := comp.Render()
    
    assert.Contains(t, html, `hx-trigger="keyup changed delay:500ms"`)
    assert.Contains(t, html, `type="search"`)
    
    // Test tabla paginada
    table := ParseComponent(`
    <gox-table hx-get="/users" paginated="true">
      <columns>
        <column name="name" />
        <column name="email" />
      </columns>
    </gox-table>
    `)
    
    html = table.Render()
    assert.Contains(t, html, "hx-get=\"/users?page=1\"")
}
```

### Tests de Performance

```go
func BenchmarkHTMXResponse(b *testing.B) {
    handler := func(ctx *gox.Context) error {
        data := generateLargeDataset()
        
        if ctx.IsHTMX() {
            return ctx.RenderPartial("list", data)
        }
        return ctx.Render("page", data)
    }
    
    req := httptest.NewRequest("GET", "/", nil)
    req.Header.Set("HX-Request", "true")
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        w := httptest.NewRecorder()
        handler(&gox.Context{Request: req, Writer: w})
    }
}
```

## Definición de Done

- [ ] HTMX integrado automáticamente
- [ ] Atributos hx-* procesados correctamente
- [ ] Helpers del servidor completos
- [ ] Optimizaciones de respuesta funcionando
- [ ] Componentes HTMX nativos
- [ ] Tests con cobertura > 85%
- [ ] Documentación con ejemplos
- [ ] Performance optimizada

## Notas Adicionales

- HTMX debe ser la versión más reciente estable
- Considerar modo de desarrollo vs producción
- Los helpers deben ser intuitivos
- Mantener compatibilidad con HTMX estándar
- Documentar patrones comunes de uso
- Considerar integración con herramientas de debug de HTMX