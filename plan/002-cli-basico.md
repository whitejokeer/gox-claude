# Task 002: CLI Básico

## Descripción
Implementar la estructura base del CLI de GOX usando cobra/cli framework, incluyendo comandos principales, sistema de flags, manejo de configuración y soporte para arquitectura distribuida con service discovery.

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
- [ ] Implementar detección de contexto (dentro/fuera de proyecto GOX)

### 2.2 Implementar comando root
- [ ] Crear comando root con descripción del framework
- [ ] Agregar flags globales (--verbose, --config, --help)
- [ ] Implementar sistema de logging configurable
- [ ] Agregar manejo de errores global
- [ ] Implementar detección de `gox.config.yaml` para contexto

### 2.3 Implementar comando `new`
- [ ] Crear subcomando `gox new project` (solo fuera de proyecto)
- [ ] Crear subcomando `gox new service` (servicio independiente)
- [ ] Implementar flags para cada tipo (--template, --db, --auth, --router)
- [ ] Crear sistema de templates base con estructura distribuida

### 2.4 Implementar comando `generate`
- [ ] Crear subcomando `gox generate service` (dentro de proyecto)
- [ ] Crear subcomando `gox generate page`
- [ ] Crear subcomando `gox generate component` (con flag --shared)
- [ ] Crear subcomando `gox generate middleware`
- [ ] Implementar flags específicos (--path, --props, --auth, --api, --db)

### 2.5 Implementar comando `dev`
- [ ] Crear comando básico que imprime "Starting development server..."
- [ ] Agregar flags (--port, --host, --no-hot-reload, --services)
- [ ] Implementar detección de proyecto GOX válido
- [ ] Agregar validación de configuración
- [ ] Placeholder para integración con docker-compose

### 2.6 Sistema de ayuda y documentación
- [ ] Implementar `gox help` con información detallada
- [ ] Agregar ejemplos en cada comando
- [ ] Crear sistema de autocompletado para bash/zsh
- [ ] Documentar todos los flags disponibles
- [ ] Incluir ejemplos de arquitectura distribuida

## Criterios de Aceptación

1. **Estructura del CLI**
   - Debe usar cobra como framework
   - Los comandos deben seguir convenciones Unix
   - Debe soportar subcomandos anidados
   - Debe detectar contexto (dentro/fuera de proyecto)

2. **Comando `new` (fuera de proyecto)**
   ```bash
   # Debe funcionar:
   gox new project my-app --router=gin --auth=jwt --db=postgres
   gox new service payment-service --db=postgres --queue=rabbitmq
   
   # Debe fallar dentro de un proyecto:
   cd my-app && gox new project another-app  # Error: Already inside a GOX project
   ```

3. **Comando `generate` (dentro de proyecto)**
   ```bash
   # Debe funcionar dentro de un proyecto:
   gox generate service users --api --db=postgres
   gox generate page dashboard --auth
   gox generate component user-card --props="name:string,email:string"
   gox generate component button --shared  # Crea en gateway/shared/ui/
   gox generate middleware rate-limit
   
   # Debe fallar fuera de proyecto:
   gox generate service users  # Error: Not inside a GOX project
   ```

4. **Comando `dev`**
   ```bash
   # Debe funcionar:
   gox dev
   gox dev --port=4000 --no-hot-reload
   gox dev --services=users,products  # Solo levantar ciertos servicios
   ```

5. **Sistema de ayuda**
   - `gox --help` debe mostrar todos los comandos
   - `gox new --help` debe mostrar ayuda específica
   - Los mensajes de error deben indicar si el comando requiere estar dentro/fuera de proyecto

## Tests Necesarios

### Tests Unitarios

1. **Test de detección de contexto**
```go
func TestIsInsideGoxProject(t *testing.T) {
    tests := []struct {
        name     string
        setup    func() error
        cleanup  func() error
        expected bool
    }{
        {
            name: "inside project",
            setup: func() error {
                return os.WriteFile("gox.config.yaml", []byte("version: 1"), 0644)
            },
            cleanup: func() error {
                return os.Remove("gox.config.yaml")
            },
            expected: true,
        },
        {
            name: "outside project",
            setup: func() error { return nil },
            cleanup: func() error { return nil },
            expected: false,
        },
    }
    // ... ejecutar tests
}
```

2. **Test de generación de componentes**
```go
func TestGenerateComponentCommand(t *testing.T) {
    tests := []struct {
        name           string
        args           []string
        expectedPath   string
    }{
        {
            name:         "regular component",
            args:         []string{"generate", "component", "user-card"},
            expectedPath: "components/user-card.gox",
        },
        {
            name:         "shared component",
            args:         []string{"generate", "component", "button", "--shared"},
            expectedPath: "shared/ui/button.gox",
        },
    }
    // ... ejecutar tests
}
```

3. **Test de estructura generada**
```go
func TestProjectStructure(t *testing.T) {
    // Test que verifica que se genera la estructura completa
    cmd := NewProjectCmd()
    cmd.Run(cmd, []string{"test-app", "--router=gin"})
    
    expectedFiles := []string{
        "test-app/gateway/main.go",
        "test-app/gateway/routing/routes.go",
        "test-app/gateway/shared/ui/.gitkeep",
        "test-app/gateway/shared/layouts/.gitkeep",
        "test-app/gateway/.env",
        "test-app/gateway/.env.example",
        "test-app/common/middleware/service_router.go",
        "test-app/common/discovery/interface.go",
        "test-app/docker-compose.yml",
        "test-app/go.work",
    }
    
    for _, file := range expectedFiles {
        if _, err := os.Stat(file); os.IsNotExist(err) {
            t.Errorf("Expected file %s not created", file)
        }
    }
}
```

### Tests de Integración

1. **Test E2E del flujo completo**
```bash
#!/bin/bash
# Test creación de proyecto
gox new project test-app --router=gin --db=postgres
cd test-app

# Verificar estructura base
[ -f "gateway/main.go" ] || exit 1
[ -d "gateway/shared/ui" ] || exit 1
[ -d "services" ] || exit 1

# Test generación de servicios
gox generate service users --api --db=postgres
[ -f "services/users/cmd/server/main.go" ] || exit 1
[ -f "services/users/go.mod" ] || exit 1

# Test generación de componentes
cd gateway
gox generate component user-card
[ -f "components/user-card.gox" ] || exit 1

gox generate component button --shared
[ -f "shared/ui/button.gox" ] || exit 1

# Test error fuera de proyecto
cd ../../..
gox generate service products 2>&1 | grep -q "Not inside a GOX project" || exit 1
```

## Definición de Done

- [ ] Todos los comandos básicos implementados con detección de contexto
- [ ] Estructura distribuida generada correctamente
- [ ] Tests unitarios con cobertura > 80%
- [ ] Tests de integración pasando
- [ ] Documentación de cada comando actualizada
- [ ] Sistema de autocompletado funcionando
- [ ] Manejo de errores robusto con mensajes claros
- [ ] Logging configurable implementado
- [ ] Service discovery integrado en templates

## Notas Adicionales

- Usar cobra-cli para generar la estructura inicial
- Los comandos `new` solo funcionan fuera de un proyecto GOX
- Los comandos `generate` solo funcionan dentro de un proyecto GOX
- Implementar templates que incluyan service discovery desde el inicio
- Asegurar que cada servicio generado sea completamente independiente
- Los templates deben incluir health checks y métricas básicas

## Estructura de Templates a Implementar

### Template: Project Gateway
```
gateway/
├── main.go                 # Con service discovery configurado
├── config/
│   └── services.yaml      # Configuración de servicios
├── routing/
│   └── routes.go          # Rutas principales
├── middleware/            # Middleware local del gateway
├── pages/                 # Páginas de la app
├── components/            # Componentes específicos
├── shared/                # Componentes compartidos
│   ├── ui/               # Design system base
│   └── layouts/          # Layouts reutilizables
├── .env
├── .env.example
├── Dockerfile
└── go.mod
```

### Template: Service
```
services/{name}/
├── cmd/server/main.go     # Con auto-registro en Consul
├── internal/
│   ├── config/
│   ├── handlers/
│   │   ├── health.go      # /health y /ready
│   │   └── routes.go
│   ├── models/
│   ├── repository/
│   └── service/
├── k8s/                   # Manifiestos Kubernetes
├── migrations/            # Si --db flag
├── api/                   # Si --api flag
├── Dockerfile
├── .env
├── .env.example
├── go.mod
└── README.md
```

## Documentación de Decisiones Arquitecturales

### 1. **¿Por qué no hay comando `new module` o `generate module`?**

En una arquitectura verdaderamente distribuida, los "módulos" como servicios independientes con UI crean problemas:

- **Problema de imports**: No podemos importar componentes de otros servicios
- **Acoplamiento**: Los módulos con UI violan el principio de servicios independientes
- **Complejidad**: Añade otra capa de abstracción innecesaria

**Decisión**: Los componentes reutilizables viven en el gateway bajo `shared/`.

### 2. **¿Por qué `gateway/shared/` en lugar de módulos separados?**

Los frontend developers esperan poder compartir componentes, pero en GOX con SSR:

- Los componentes se renderizan en el servidor
- No hay "build step" que combine módulos
- HTMX espera HTML del servidor, no componentes JS

**Decisión**: Carpeta `shared/` dentro del gateway para componentes reutilizables localmente.

### 3. **¿Por qué servicios no tienen UI (páginas/componentes)?**

Principio de responsabilidad única:

- **Services**: Solo exponen APIs (REST/gRPC)
- **Gateway**: Maneja toda la UI y composición

Esto permite:
- Escalar servicios independientemente
- Cambiar UI sin tocar servicios
- Diferentes equipos en backend/frontend

### 4. **¿Por qué Service Discovery desde el inicio?**

En lugar de hardcodear URLs de servicios:

```go
// ❌ Malo
userService := "http://localhost:8081"

// ✅ Bueno
userService := discovery.Find("users")
```

**Beneficios**:
- Desarrollo local igual que producción
- Fácil agregar/quitar servicios
- Preparado para Kubernetes/Cloud

### 5. **¿Por qué detección de contexto (dentro/fuera de proyecto)?**

Evita errores comunes:
- Crear proyectos dentro de proyectos
- Generar recursos fuera de un proyecto
- Confusión sobre dónde ejecutar comandos

**Regla simple**:
- `new` = crear algo nuevo (fuera de proyecto)
- `generate` = agregar a proyecto existente (dentro de proyecto)

### 6. **¿Qué pasa si realmente necesito módulos UI compartidos entre proyectos?**

Para equipos enterprise con múltiples gateways, se puede:

1. Crear un repositorio de componentes
2. Copiar componentes entre proyectos
3. Usar git submodules (no recomendado)

Pero esto es un caso edge. La mayoría de proyectos tienen un solo gateway.

### 7. **¿Por qué cada servicio tiene su propio go.mod?**

Independencia total:
- Diferentes versiones de dependencias
- Deploy independiente
- Equipos autónomos
- Sin conflictos de versiones

El `go.work` es solo para desarrollo local.

### 8. **¿Por qué no hay plantillas para diferentes tipos de servicios?**

Mantener simplicidad:
- Un solo blueprint probado
- Fácil de entender y mantener
- Flexibilidad para evolucionar

Los flags (--db, --api) agregan capacidades, no cambian la estructura base.

### 9. **Decisión sobre HTMX y componentes**

Los componentes .gox se renderizan en el servidor y devuelven HTML:

```vue
<!-- Componente -->
<template>
  <button hx-post="{{ .Action }}">{{ .Label }}</button>
</template>

<!-- Se renderiza como -->
<button hx-post="/api/users/create">Create User</button>
```

No hay JavaScript del lado del cliente, todo es HTML + HTMX.

### 10. **¿Por qué Consul para Service Discovery?**

- Maduro y probado en producción
- Funciona igual en local y cloud
- UI web para debugging
- Health checks integrados
- Compatible con Kubernetes

Alternativas consideradas:
- etcd: Más complejo
- Kubernetes DNS: Solo en K8s
- Hardcoded: No escala