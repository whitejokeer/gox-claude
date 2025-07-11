# Task 003: Parser de Archivos .gox

## Descripción
Implementar un parser robusto para archivos .gox que pueda extraer y procesar las secciones template, go, y style, manteniendo la estructura y metadata necesaria para la compilación.

## Prioridad
Alta

## Estimación
5-6 días

## Dependencias
- Task 001: Setup inicial del proyecto
- Task 002: CLI básico (parcial)

## Subtasks

### 3.1 Diseñar AST (Abstract Syntax Tree)
- [ ] Definir estructura del AST para archivos .gox
- [ ] Crear tipos para cada sección (Template, Go, Style, Story)
- [ ] Diseñar sistema de metadata y atributos
- [ ] Implementar visitor pattern para el AST

### 3.2 Implementar lexer/tokenizer
- [ ] Crear tokenizer para identificar tags de apertura/cierre
- [ ] Manejar atributos en tags (auth, layout, scoped)
- [ ] Tokenizar contenido interno respetando sintaxis
- [ ] Implementar manejo de errores con línea/columna

### 3.3 Implementar parser principal
- [ ] Parser para sección `<template>`
- [ ] Parser para sección `<go>`
- [ ] Parser para sección `<style>`
- [ ] Parser para sección `<story>` (opcional)
- [ ] Validar estructura y orden de secciones

### 3.4 Procesamiento de template
- [ ] Extraer directivas especiales (auth, layout)
- [ ] Identificar sintaxis Go template (`{{.Variable}}`)
- [ ] Detectar atributos HTMX (hx-*)
- [ ] Preservar formato y espaciado original

### 3.5 Procesamiento de código Go
- [ ] Validar sintaxis Go usando go/parser
- [ ] Extraer imports automáticamente
- [ ] Identificar struct principal del componente
- [ ] Detectar métodos lifecycle (Mount, BeforeMount, etc.)
- [ ] Extraer handlers HTMX

### 3.6 Procesamiento de estilos
- [ ] Detectar tipo de estilo (CSS, SCSS, Tailwind)
- [ ] Manejar atributo `scoped`
- [ ] Procesar directivas Tailwind (@apply)
- [ ] Validar sintaxis CSS básica

## Criterios de Aceptación

1. **Parser robusto**
   - Debe parsear archivos .gox válidos sin errores
   - Debe dar mensajes de error claros para archivos inválidos
   - Debe preservar formato y comentarios

2. **Estructura AST**
   ```go
   type GoxFile struct {
       Path     string
       Template *TemplateNode
       Go       *GoNode
       Styles   []*StyleNode
       Story    *StoryNode
       Metadata map[string]interface{}
   }
   
   type TemplateNode struct {
       Content    string
       Auth       string // "required", "role:admin", etc.
       Layout     string
       Directives []Directive
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
type HelloPage struct {
    Name string
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

2. **Test de atributos especiales**
```go
func TestParseTemplateAttributes(t *testing.T) {
    input := `<template auth="required" layout="admin">
  <div>Admin Panel</div>
</template>`
    
    ast, err := ParseString(input)
    assert.NoError(t, err)
    assert.Equal(t, "required", ast.Template.Auth)
    assert.Equal(t, "admin", ast.Template.Layout)
}
```

3. **Test de errores**
```go
func TestParseErrors(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        wantErr string
    }{
        {
            name:    "missing closing tag",
            input:   `<template><div>Hello`,
            wantErr: "unclosed tag at line 1",
        },
        {
            name:    "invalid go code",
            input:   `<go>func invalid syntax</go>`,
            wantErr: "invalid Go syntax at line 1",
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            _, err := ParseString(tt.input)
            assert.Contains(t, err.Error(), tt.wantErr)
        })
    }
}
```

### Tests de Integración

1. **Test con archivos reales**
```go
func TestParseRealGoxFiles(t *testing.T) {
    files := []string{
        "testdata/basic.gox",
        "testdata/complex.gox",
        "testdata/with-htmx.gox",
        "testdata/with-tailwind.gox",
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

2. **Test de casos edge**
```go
func TestEdgeCases(t *testing.T) {
    // Archivo solo con template
    ast1, err := ParseString(`<template><div>Hello</div></template>`)
    assert.NoError(t, err)
    assert.Nil(t, ast1.Go)
    
    // Múltiples secciones style
    input := `
<style>/* Global */</style>
<style scoped>/* Scoped */</style>
`
    ast2, err := ParseString(input)
    assert.NoError(t, err)
    assert.Len(t, ast2.Styles, 2)
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
```

## Definición de Done

- [ ] Parser completo con todas las secciones
- [ ] AST bien diseñado y documentado
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
- Considerar implementar un modo de "recuperación" para parsear archivos parcialmente válidos