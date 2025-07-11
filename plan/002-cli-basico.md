# Task 002: CLI Básico

## Descripción
Implementar la estructura base del CLI de GOX usando cobra/cli framework, incluyendo comandos principales, sistema de flags, y manejo de configuración.

## Prioridad
Alta

## Estimación
3-4 días

## Dependencias
- Task 001: Setup inicial del proyecto

## Subtasks

### 2.1 Configurar framework CLI
- [ ] Evaluar y elegir framework (cobra, urfave/cli, etc.)
- [ ] Agregar dependencias necesarias
- [ ] Crear estructura base del CLI en `cmd/gox/`
- [ ] Implementar sistema de versioning

### 2.2 Implementar comando root
- [ ] Crear comando root con descripción del framework
- [ ] Agregar flags globales (--verbose, --config, --help)
- [ ] Implementar sistema de logging configurable
- [ ] Agregar manejo de errores global

### 2.3 Implementar comando `new`
- [ ] Crear subcomando `gox new project`
- [ ] Crear subcomando `gox new module`
- [ ] Crear subcomando `gox new service`
- [ ] Implementar flags para cada tipo (--template, --db, --auth)
- [ ] Crear sistema de templates base

### 2.4 Implementar comando `generate`
- [ ] Crear subcomando `gox generate page`
- [ ] Crear subcomando `gox generate component`
- [ ] Crear subcomando `gox generate service`
- [ ] Crear subcomando `gox generate middleware`
- [ ] Implementar flags específicos (--path, --props, --auth)

### 2.5 Implementar comando `dev`
- [ ] Crear comando básico que imprime "Starting development server..."
- [ ] Agregar flags (--port, --host, --no-hot-reload)
- [ ] Implementar detección de proyecto GOX válido
- [ ] Agregar validación de configuración

### 2.6 Sistema de ayuda y documentación
- [ ] Implementar `gox help` con información detallada
- [ ] Agregar ejemplos en cada comando
- [ ] Crear sistema de autocompletado para bash/zsh
- [ ] Documentar todos los flags disponibles

## Criterios de Aceptación

1. **Estructura del CLI**
   - Debe usar un framework robusto (preferiblemente cobra)
   - Los comandos deben seguir convenciones Unix
   - Debe soportar subcomandos anidados

2. **Comando `new`**
   ```bash
   # Debe funcionar:
   gox new project my-app
   gox new project my-app --db=postgres --auth=jwt
   gox new module user-management --protocol=grpc
   gox new service notifications --template=minimal
   ```

3. **Comando `generate`**
   ```bash
   # Debe funcionar:
   gox generate page dashboard --auth
   gox generate component user-card --props="name,email,avatar"
   gox generate service users --api
   ```

4. **Comando `dev`**
   ```bash
   # Debe funcionar:
   gox dev
   gox dev --port=4000 --no-hot-reload
   ```

5. **Sistema de ayuda**
   - `gox --help` debe mostrar todos los comandos
   - `gox new --help` debe mostrar ayuda específica
   - Los mensajes de error deben ser claros y útiles

## Tests Necesarios

### Tests Unitarios

1. **Test de parsing de comandos**
```go
func TestNewProjectCommand(t *testing.T) {
    tests := []struct {
        name     string
        args     []string
        wantErr  bool
        expected Config
    }{
        {
            name: "basic project",
            args: []string{"new", "project", "my-app"},
            wantErr: false,
            expected: Config{
                Type: "project",
                Name: "my-app",
            },
        },
        {
            name: "project with flags",
            args: []string{"new", "project", "my-app", "--db=postgres", "--auth=jwt"},
            wantErr: false,
            expected: Config{
                Type: "project",
                Name: "my-app",
                DB: "postgres",
                Auth: "jwt",
            },
        },
    }
    // ... ejecutar tests
}
```

2. **Test de validación**
```go
func TestValidateProjectName(t *testing.T) {
    valid := []string{"my-app", "app123", "test_app"}
    invalid := []string{"", "123app", "app with spaces", "app/slash"}
    
    for _, name := range valid {
        if err := validateProjectName(name); err != nil {
            t.Errorf("Expected %s to be valid", name)
        }
    }
    
    for _, name := range invalid {
        if err := validateProjectName(name); err == nil {
            t.Errorf("Expected %s to be invalid", name)
        }
    }
}
```

### Tests de Integración

1. **Test E2E del CLI**
```bash
#!/bin/bash
# Test comando new
gox new project test-app --db=postgres
if [ ! -d "test-app" ]; then
    echo "Project directory not created"
    exit 1
fi

# Test comando generate
cd test-app
gox generate page home
if [ ! -f "pages/home.gox" ]; then
    echo "Page not generated"
    exit 1
fi
```

2. **Test de help**
```bash
# Verificar que help funciona
gox --help | grep -q "GOX Framework CLI"
gox new --help | grep -q "Create a new GOX project"
```

## Definición de Done

- [ ] Todos los comandos básicos implementados
- [ ] Tests unitarios con cobertura > 80%
- [ ] Tests de integración pasando
- [ ] Documentación de cada comando
- [ ] Sistema de autocompletado funcionando
- [ ] Manejo de errores robusto
- [ ] Logging configurable implementado

## Notas Adicionales

- Considerar usar cobra-cli para generar la estructura inicial
- Implementar un sistema de plugins para futura extensibilidad
- Los comandos deben ser intuitivos y seguir convenciones de otros CLIs populares
- Pensar en la experiencia del desarrollador desde el principio
- Agregar telemetría opcional para entender el uso