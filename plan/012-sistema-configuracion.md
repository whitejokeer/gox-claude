# Task 012: Sistema de Configuración

## Descripción
Implementar un sistema robusto de configuración que soporte múltiples formatos, validación, variables de entorno, configuración por ambiente, y hot reload de configuración.

## Prioridad
Alta

## Estimación
3-4 días

## Dependencias
- Task 001: Setup inicial del proyecto

## Subtasks

### 12.1 Formatos de configuración soportados
- [ ] YAML (gox.config.yaml) - principal
- [ ] JSON (gox.config.json) - alternativo
- [ ] TOML (gox.config.toml) - alternativo
- [ ] Variables de entorno
- [ ] Flags de línea de comandos

### 12.2 Estructura de configuración
- [ ] Definir schema completo de configuración
- [ ] Configuración jerárquica anidada
- [ ] Herencia de configuración
- [ ] Overrides por ambiente
- [ ] Configuración por servicio

### 12.3 Validación y tipos
- [ ] Validación de schema con JSON Schema
- [ ] Tipos fuertemente tipados en Go
- [ ] Validación de valores en runtime
- [ ] Mensajes de error descriptivos
- [ ] Configuración requerida vs opcional

### 12.4 Variables de entorno
- [ ] Interpolación de variables ${VAR}
- [ ] Variables requeridas vs opcionales
- [ ] Valores por defecto
- [ ] Validación de formato
- [ ] Documentación automática de env vars

### 12.5 Configuración por ambiente
- [ ] Archivos por ambiente (dev, staging, prod)
- [ ] Override automático por NODE_ENV/GOX_ENV
- [ ] Configuración específica por servicio
- [ ] Secrets management
- [ ] Configuración distribuida

### 12.6 Hot reload y observación
- [ ] Detectar cambios en archivos de config
- [ ] Reload automático sin reiniciar
- [ ] Notificar componentes de cambios
- [ ] Validar nueva configuración
- [ ] Rollback en caso de error

## Criterios de Aceptación

1. **Configuración base funcionando**
   ```yaml
   # gox.config.yaml
   name: "my-app"
   version: "1.0.0"
   
   dev:
     port: ${PORT:-3000}
     host: "localhost"
     debug: true
     
   database:
     driver: "postgres"
     host: ${DB_HOST}
     port: ${DB_PORT:-5432}
     name: ${DB_NAME}
     user: ${DB_USER}
     password: ${DB_PASSWORD}
     
   auth:
     jwt:
       secret: ${JWT_SECRET}
       expires: "24h"
   ```

2. **API de configuración**
   ```go
   // Cargar configuración
   config, err := gox.LoadConfig("gox.config.yaml")
   if err != nil {
       log.Fatal(err)
   }
   
   // Acceso tipado
   port := config.Dev.Port
   dbHost := config.Database.Host
   
   // Validación automática
   if err := config.Validate(); err != nil {
       log.Fatal("Invalid configuration:", err)
   }
   
   // Watch para cambios
   config.OnChange(func(key string, oldVal, newVal interface{}) {
       log.Printf("Config changed: %s = %v", key, newVal)
   })
   ```

3. **Configuración por ambiente**
   ```
   config/
   ├── gox.config.yaml         # Base
   ├── gox.config.dev.yaml     # Development
   ├── gox.config.staging.yaml # Staging
   └── gox.config.prod.yaml    # Production
   ```

4. **Variables de entorno**
   ```bash
   # .env
   GOX_ENV=development
   PORT=3000
   DB_HOST=localhost
   DB_PASSWORD=secret
   JWT_SECRET=super-secret-key
   ```

## Tests Necesarios

### Tests Unitarios

1. **Test carga de configuración**
```go
func TestConfigLoad(t *testing.T) {
    configYAML := `
name: test-app
version: 1.0.0
dev:
  port: 3000
  debug: true
database:
  driver: postgres
  host: localhost
`
    
    config, err := gox.ParseConfig([]byte(configYAML))
    assert.NoError(t, err)
    
    assert.Equal(t, "test-app", config.Name)
    assert.Equal(t, "1.0.0", config.Version)
    assert.Equal(t, 3000, config.Dev.Port)
    assert.True(t, config.Dev.Debug)
}
```

2. **Test interpolación de variables**
```go
func TestEnvironmentInterpolation(t *testing.T) {
    os.Setenv("TEST_PORT", "4000")
    os.Setenv("TEST_HOST", "example.com")
    
    configYAML := `
dev:
  port: ${TEST_PORT}
  host: ${TEST_HOST}
  fallback: ${MISSING_VAR:-default_value}
`
    
    config, err := gox.ParseConfig([]byte(configYAML))
    assert.NoError(t, err)
    
    assert.Equal(t, 4000, config.Dev.Port)
    assert.Equal(t, "example.com", config.Dev.Host)
    assert.Equal(t, "default_value", config.Dev.Fallback)
}
```

3. **Test validación**
```go
func TestConfigValidation(t *testing.T) {
    tests := []struct {
        name    string
        config  string
        wantErr string
    }{
        {
            name: "missing required field",
            config: `
name: test
# version is required
`,
            wantErr: "version is required",
        },
        {
            name: "invalid port",
            config: `
name: test
version: 1.0.0
dev:
  port: -1
`,
            wantErr: "port must be between 1 and 65535",
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            config, err := gox.ParseConfig([]byte(tt.config))
            if err == nil {
                err = config.Validate()
            }
            
            if tt.wantErr != "" {
                assert.Error(t, err)
                assert.Contains(t, err.Error(), tt.wantErr)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

### Tests de Integración

1. **Test configuración por ambiente**
```go
func TestEnvironmentConfiguration(t *testing.T) {
    // Crear archivos de config
    createConfigFile(t, "gox.config.yaml", `
name: test-app
dev:
  port: 3000
  debug: true
`)
    
    createConfigFile(t, "gox.config.prod.yaml", `
dev:
  port: 8080
  debug: false
`)
    
    // Test environment development
    os.Setenv("GOX_ENV", "development")
    config, err := gox.LoadConfig(".")
    assert.NoError(t, err)
    assert.Equal(t, 3000, config.Dev.Port)
    assert.True(t, config.Dev.Debug)
    
    // Test environment production
    os.Setenv("GOX_ENV", "production")
    config, err = gox.LoadConfig(".")
    assert.NoError(t, err)
    assert.Equal(t, 8080, config.Dev.Port)
    assert.False(t, config.Dev.Debug)
}
```

2. **Test hot reload**
```go
func TestConfigHotReload(t *testing.T) {
    configFile := "test-config.yaml"
    createConfigFile(t, configFile, `
name: test
dev:
  port: 3000
`)
    
    config, err := gox.LoadConfig(configFile)
    assert.NoError(t, err)
    
    // Watch for changes
    changed := make(chan bool, 1)
    config.OnChange(func(key string, old, new interface{}) {
        if key == "dev.port" {
            changed <- true
        }
    })
    
    // Modify config
    time.Sleep(100 * time.Millisecond)
    updateConfigFile(t, configFile, `
name: test
dev:
  port: 4000
`)
    
    // Wait for change
    select {
    case <-changed:
        assert.Equal(t, 4000, config.Dev.Port)
    case <-time.After(2 * time.Second):
        t.Fatal("Config change not detected")
    }
}
```

3. **Test múltiples formatos**
```go
func TestMultipleConfigFormats(t *testing.T) {
    tests := []struct {
        format   string
        filename string
        content  string
    }{
        {
            format:   "yaml",
            filename: "config.yaml",
            content:  "name: test\nversion: 1.0.0",
        },
        {
            format:   "json",
            filename: "config.json",
            content:  `{"name": "test", "version": "1.0.0"}`,
        },
        {
            format:   "toml",
            filename: "config.toml",
            content:  "name = \"test\"\nversion = \"1.0.0\"",
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.format, func(t *testing.T) {
            createFile(t, tt.filename, tt.content)
            
            config, err := gox.LoadConfig(tt.filename)
            assert.NoError(t, err)
            assert.Equal(t, "test", config.Name)
            assert.Equal(t, "1.0.0", config.Version)
        })
    }
}
```

### Tests de Performance

```go
func BenchmarkConfigLoad(b *testing.B) {
    configContent := createLargeConfig()
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        gox.ParseConfig(configContent)
    }
}

func BenchmarkConfigValidation(b *testing.B) {
    config := createValidConfig()
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        config.Validate()
    }
}
```

## Definición de Done

- [ ] Múltiples formatos soportados (YAML, JSON, TOML)
- [ ] Interpolación de variables de entorno
- [ ] Configuración por ambiente
- [ ] Validación robusta con mensajes claros
- [ ] Hot reload funcionando
- [ ] API tipada y fácil de usar
- [ ] Tests con cobertura > 90%
- [ ] Documentación completa de configuración

## Notas Adicionales

- La configuración debe ser fácil de usar y entender
- Los errores deben ser descriptivos y útiles
- Considerar encriptación para secretos sensibles
- Documentar todas las opciones disponibles
- El sistema debe ser extensible para nuevas opciones
- Pensar en configuración distribuida para el futuro