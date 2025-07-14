# Task 003: Parser de Archivos .gox

## Descripción
Implementar un parser robusto para archivos .gox que pueda extraer y procesar las secciones template, go, y style, manteniendo la estructura y metadata necesaria para la compilación.

## Prioridad
Alta

## Estimación
5-6 días

## Dependencias
- Task 001: Setup inicial del proyecto
- Task 002: NO BLOQUEANTE - Se puede desarrollar en paralelo
- Sincronización: Definir estructura .gox antes de empezar

## Subtasks

### 3.1 Diseñar AST (Abstract Syntax Tree)
- [ ] Definir estructura del AST para archivos .gox
- [ ] Crear tipos para cada sección (Template, Go, Style)
- [ ] Diseñar sistema de metadata y atributos
- [ ] Implementar visitor pattern para el AST
- [ ] Definir estructura para componentes embebidos

### 3.2 Implementar lexer/tokenizer
- [ ] Crear tokenizer para identificar tags de apertura/cierre
- [ ] Manejar atributos en tags (auth, layout, scoped)
- [ ] Tokenizar contenido interno respetando sintaxis
- [ ] Implementar manejo de errores con línea/columna
- [ ] Detectar componentes custom (`<user-card>`, `<shared-button>`)

### 3.3 Implementar parser principal
- [ ] Parser para sección `<template>`
- [ ] Parser para sección `<go>`
- [ ] Parser para sección `<style>`
- [ ] Validar estructura y orden de secciones
- [ ] Manejar secciones opcionales

### 3.4 Procesamiento de template
- [ ] Extraer directivas especiales (auth, layout)
- [ ] Identificar sintaxis Go template (`{{.Variable}}`)
- [ ] Detectar atributos HTMX (hx-*)
- [ ] Parsear componentes custom y sus props
- [ ] Resolver rutas de componentes (components/ vs shared/)
- [ ] Preservar formato y espaciado original

### 3.5 Procesamiento de código Go
- [ ] Validar sintaxis Go usando go/parser
- [ ] Extraer imports automáticamente
- [ ] Identificar struct principal del componente
- [ ] Detectar handlers HTTP (HandleFunc pattern)
- [ ] Identificar props struct para componentes
- [ ] Extraer métodos HTMX (acciones/validaciones)

### 3.6 Procesamiento de estilos
- [ ] Detectar tipo de estilo (CSS o Tailwind classes)
- [ ] Manejar atributo `scoped`
- [ ] Procesar directivas Tailwind (@apply)
- [ ] Validar sintaxis CSS básica
- [ ] Preparar para futuro CSS-in-Go si necesario

## Criterios de Aceptación

1. **Parser robusto**
   - Debe parsear archivos .gox válidos sin errores
   - Debe dar mensajes de error claros para archivos inválidos
   - Debe preservar formato y comentarios
   - Debe detectar componentes y sus dependencias

2. **Estructura AST**
   ```go
   type GoxFile struct {
       Path       string
       Template   *TemplateNode
       Go         *GoNode
       Styles     []*StyleNode
       Components []ComponentDependency // Componentes que usa
       Metadata   map[string]interface{}
   }
   
   type TemplateNode struct {
       Content    string
       Auth       string // "required", "role:admin", etc.
       Layout     string
       Elements   []Element // Árbol de elementos HTML y componentes
   }
   
   type Element struct {
       Type       string // "div", "user-card", etc.
       IsComponent bool
       Props      map[string]string // Incluye hx-* attributes
       Children   []Element
   }
   
   type GoNode struct {
       Source    string
       Imports   []string
       MainType  string // "UserCard" struct name
       Handlers  []Handler
       Props     *PropsStruct // Si es un componente
   }
   ```

3. **Manejo de errores**
   ```go
   type ParseError struct {
       File    string
       Line    int
       Column  int
       Message string
   }
   ```

4. **API del parser**
   ```go
   // Debe funcionar:
   ast, err := parser.ParseFile("pages/index.gox")
   if err != nil {
       // Error con información de línea/columna
   }
   
   // Acceso fácil a secciones
   template := ast.Template.Content
   goCode := ast.Go.Source
   
   // Detección de dependencias
   for _, comp := range ast.Components {
       fmt.Printf("Uses: %s from %s\n", comp.Name, comp.Path)
   }
   ```

## Tests Necesarios

### Tests Unitarios

1. **Test de parser básico**
```go
func TestParseBasicGoxFile(t *testing.T) {
    input := `
<template>
  <div>Hello {{.Name}}</div>
</template>

<go>
package pages

type HelloPage struct {
    Name string
}

func (p *HelloPage) HandleRequest(w http.ResponseWriter, r *http.Request) {
    // Handler logic
}
</go>

<style>
.hello { color: blue; }
</style>
`
    ast, err := ParseString(input)
    assert.NoError(t, err)
    assert.NotNil(t, ast.Template)
    assert.NotNil(t, ast.Go)
    assert.Len(t, ast.Styles, 1)
}
```

2. **Test de componentes embebidos**
```go
func TestParseWithComponents(t *testing.T) {
    input := `
<template>
  <div>
    <user-card name="{{.User.Name}}" email="{{.User.Email}}" />
    <shared-button text="Save" hx-post="/save" />
  </div>
</template>
`
    ast, err := ParseString(input)
    assert.NoError(t, err)
    assert.Len(t, ast.Components, 2)
    assert.Equal(t, "user-card", ast.Components[0].Name)
    assert.Equal(t, "components/user-card", ast.Components[0].Path)
    assert.Equal(t, "shared-button", ast.Components[1].Name) 
    assert.Equal(t, "shared/ui/button", ast.Components[1].Path)
}
```

3. **Test de handlers HTTP**
```go
func TestParseHTTPHandlers(t *testing.T) {
    input := `
<go>
package components

type UserForm struct {
    User *User
}

func (f *UserForm) HandleSubmit(w http.ResponseWriter, r *http.Request) {
    // Process form submission
}

func (f *UserForm) ValidateEmail(email string) error {
    // HTMX validation endpoint
    return nil
}
</go>
`
    ast, err := ParseString(input)
    assert.NoError(t, err)
    assert.Len(t, ast.Go.Handlers, 2)
    assert.Equal(t, "HandleSubmit", ast.Go.Handlers[0].Name)
    assert.Equal(t, "ValidateEmail", ast.Go.Handlers[1].Name)
}
```

### Tests de Integración

1. **Test con archivos reales**
```go
func TestParseRealGoxFiles(t *testing.T) {
    files := []string{
        "testdata/pages/home.gox",
        "testdata/components/user-card.gox",
        "testdata/shared/ui/button.gox",
    }
    
    for _, file := range files {
        t.Run(file, func(t *testing.T) {
            ast, err := ParseFile(file)
            assert.NoError(t, err)
            assert.NotNil(t, ast)
            
            // Validar que se puede reconstruir
            output := ast.String()
            assert.NotEmpty(t, output)
        })
    }
}
```

2. **Test de resolución de componentes**
```go
func TestComponentResolution(t *testing.T) {
    // Setup filesystem mock
    fs := mockfs.New()
    fs.AddFile("components/user-card.gox", userCardContent)
    fs.AddFile("shared/ui/button.gox", buttonContent)
    
    parser := NewParser(WithFilesystem(fs))
    ast, err := parser.ParseFile("pages/users.gox")
    
    assert.NoError(t, err)
    // Verificar que resolvió correctamente las rutas
    for _, comp := range ast.Components {
        assert.True(t, fs.Exists(comp.Path + ".gox"))
    }
}
```

### Benchmarks

```go
func BenchmarkParseGoxFile(b *testing.B) {
    content, _ := os.ReadFile("testdata/complex.gox")
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, _ = ParseBytes(content)
    }
}

func BenchmarkParseWithComponents(b *testing.B) {
    // Benchmark específico para archivos con muchos componentes
    content, _ := os.ReadFile("testdata/many-components.gox")
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, _ = ParseBytes(content)
    }
}
```

## Definición de Done

- [ ] Parser completo con todas las secciones
- [ ] AST bien diseñado y documentado
- [ ] Detección automática de componentes
- [ ] Resolución de rutas de componentes
- [ ] Manejo de errores con información útil
- [ ] Tests unitarios con cobertura > 90%
- [ ] Tests de integración con archivos reales
- [ ] Benchmarks mostrando performance aceptable
- [ ] Documentación de la API del parser
- [ ] Ejemplos de uso en la documentación

## Notas Adicionales

- Considerar usar una librería de parsing como participle o escribir parser manual
- El parser debe ser lo suficientemente flexible para futuras extensiones
- Mantener compatibilidad con herramientas de Go (gofmt, gopls)
- Pensar en mensajes de error amigables para desarrolladores
- La resolución de componentes debe seguir convenciones:
  - `<user-card>` → `components/user-card.gox`
  - `<shared-button>` → `shared/ui/button.gox`
  - Permitir configuración de estas convenciones más adelante