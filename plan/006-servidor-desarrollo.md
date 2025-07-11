# Task 006: Servidor de Desarrollo

## Descripción
Implementar un servidor de desarrollo con hot reload, detección automática de cambios, compilación incremental, y herramientas de debugging integradas.

## Prioridad
Alta

## Estimación
4-5 días

## Dependencias
- Task 004: Compilador .gox a Go
- Task 005: Sistema de routing

## Subtasks

### 6.1 Servidor HTTP base
- [ ] Crear servidor HTTP configurable
- [ ] Implementar graceful shutdown
- [ ] Manejar señales del sistema (SIGINT, SIGTERM)
- [ ] Configurar timeouts apropiados
- [ ] Implementar health checks

### 6.2 File watcher
- [ ] Implementar watcher para archivos .gox
- [ ] Detectar cambios en archivos Go
- [ ] Monitorear archivos CSS/assets
- [ ] Detectar nuevos archivos/directorios
- [ ] Implementar debouncing de eventos

### 6.3 Hot reload
- [ ] Implementar WebSocket para comunicación
- [ ] Inyectar script de hot reload en HTML
- [ ] Recargar página automáticamente
- [ ] Preservar estado cuando sea posible
- [ ] Mostrar errores en el navegador

### 6.4 Compilación incremental
- [ ] Compilar solo archivos modificados
- [ ] Mantener caché de compilación
- [ ] Detectar dependencias afectadas
- [ ] Recompilar en paralelo
- [ ] Mostrar progreso de compilación

### 6.5 Dev tools
- [ ] Implementar request logger
- [ ] Panel de debugging en browser
- [ ] Visualizador de rutas activas
- [ ] Monitor de performance
- [ ] Inspector de componentes

### 6.6 Integración con Storybook
- [ ] Lanzar Storybook en paralelo
- [ ] Sincronizar cambios con Storybook
- [ ] Compartir configuración
- [ ] Proxy requests cuando necesario
- [ ] Unificar logs

## Criterios de Aceptación

1. **Servidor funcional**
   ```bash
   gox dev
   # Output:
   # 🚀 Server running on http://localhost:3000
   # 📖 Storybook running on http://localhost:6006
   # 👀 Watching for changes...
   ```

2. **Hot reload rápido**
   - Cambios en .gox < 500ms
   - Cambios en CSS < 100ms
   - Sin pérdida de estado cuando posible
   - Errores mostrados en browser

3. **File watching robusto**
   ```
   [watch] pages/index.gox changed
   [compile] Compiling pages/index.gox...
   [compile] ✓ Compiled in 243ms
   [reload] Reloading browsers...
   ```

4. **Dev tools en browser**
   ```html
   <!-- Inyectado automáticamente -->
   <div id="gox-devtools">
     <div class="gox-errors"></div>
     <div class="gox-performance"></div>
     <div class="gox-routes"></div>
   </div>
   <script src="/__gox__/devtools.js"></script>
   ```

## Tests Necesarios

### Tests Unitarios

1. **Test file watcher**
```go
func TestFileWatcher(t *testing.T) {
    watcher := NewWatcher()
    changes := make(chan FileChange, 10)
    
    watcher.OnChange(func(change FileChange) {
        changes <- change
    })
    
    watcher.Watch("testdata/watch")
    
    // Crear archivo
    os.WriteFile("testdata/watch/test.gox", []byte("test"), 0644)
    
    select {
    case change := <-changes:
        assert.Equal(t, "testdata/watch/test.gox", change.Path)
        assert.Equal(t, Created, change.Type)
    case <-time.After(1 * time.Second):
        t.Fatal("No change detected")
    }
}
```

2. **Test compilación incremental**
```go
func TestIncrementalCompilation(t *testing.T) {
    compiler := NewIncrementalCompiler()
    
    // Primera compilación
    result1, err := compiler.Compile("pages/index.gox")
    assert.NoError(t, err)
    assert.True(t, result1.FullCompile)
    
    // Sin cambios - debe usar caché
    result2, err := compiler.Compile("pages/index.gox")
    assert.NoError(t, err)
    assert.True(t, result2.FromCache)
    
    // Con cambios
    os.Chtimes("pages/index.gox", time.Now(), time.Now())
    result3, err := compiler.Compile("pages/index.gox")
    assert.NoError(t, err)
    assert.False(t, result3.FromCache)
}
```

3. **Test hot reload websocket**
```go
func TestHotReloadWebSocket(t *testing.T) {
    server := NewDevServer()
    go server.Start()
    defer server.Stop()
    
    // Conectar cliente WebSocket
    ws, _, err := websocket.DefaultDialer.Dial("ws://localhost:3000/__gox__/ws", nil)
    assert.NoError(t, err)
    defer ws.Close()
    
    // Simular cambio
    server.NotifyChange(FileChange{
        Path: "pages/index.gox",
        Type: Modified,
    })
    
    // Recibir mensaje de reload
    var msg ReloadMessage
    err = ws.ReadJSON(&msg)
    assert.NoError(t, err)
    assert.Equal(t, "reload", msg.Type)
}
```

### Tests de Integración

1. **Test servidor completo**
```go
func TestDevServerIntegration(t *testing.T) {
    // Crear proyecto de prueba
    createTestProject(t, "testdata/project")
    
    // Iniciar servidor
    server := NewDevServer(Config{
        Port: 3001,
        Root: "testdata/project",
    })
    
    go server.Start()
    defer server.Stop()
    
    // Esperar que inicie
    waitForServer(t, "http://localhost:3001")
    
    // Verificar página principal
    resp, err := http.Get("http://localhost:3001")
    assert.NoError(t, err)
    assert.Equal(t, 200, resp.StatusCode)
    
    // Verificar hot reload script
    body, _ := io.ReadAll(resp.Body)
    assert.Contains(t, string(body), "__gox__/client.js")
}
```

2. **Test cambio de archivos**
```go
func TestFileChangeReload(t *testing.T) {
    server := NewDevServer()
    reloaded := make(chan bool, 1)
    
    server.OnReload(func() {
        reloaded <- true
    })
    
    go server.Start()
    defer server.Stop()
    
    // Modificar archivo
    content, _ := os.ReadFile("pages/index.gox")
    os.WriteFile("pages/index.gox", append(content, []byte("\n<!-- test -->")...), 0644)
    
    select {
    case <-reloaded:
        // Success
    case <-time.After(2 * time.Second):
        t.Fatal("Reload not triggered")
    }
}
```

### Tests de Performance

```go
func TestCompilationSpeed(t *testing.T) {
    files := []string{
        "pages/index.gox",
        "pages/about.gox",
        "pages/contact.gox",
    }
    
    compiler := NewIncrementalCompiler()
    
    start := time.Now()
    for _, file := range files {
        _, err := compiler.Compile(file)
        assert.NoError(t, err)
    }
    duration := time.Since(start)
    
    // Debe compilar 3 archivos en menos de 1 segundo
    assert.Less(t, duration, 1*time.Second)
}
```

## Definición de Done

- [ ] Servidor con hot reload funcionando
- [ ] File watcher robusto y eficiente
- [ ] Compilación incremental < 500ms
- [ ] WebSocket para hot reload
- [ ] Dev tools en browser
- [ ] Integración con Storybook
- [ ] Tests con cobertura > 85%
- [ ] Documentación de uso

## Notas Adicionales

- Usar fsnotify para file watching
- Considerar usar esbuild para assets
- El hot reload debe ser opcional
- Mostrar errores de forma amigable
- Considerar modo offline/PWA
- Pensar en debugging remoto
- Logs deben ser claros y útiles