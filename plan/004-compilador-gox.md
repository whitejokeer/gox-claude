# Task 004: Compilador .gox a Go

## Descripción
Implementar el compilador que transforma archivos .gox parseados (AST) en código Go ejecutable, generando handlers HTTP, templates compilados, y assets estáticos.

## Prioridad
Alta

## Estimación
6-7 días

## Dependencias
- Task 003: Parser de archivos .gox

## Subtasks

### 4.1 Diseñar arquitectura del compilador
- [ ] Definir estructura de salida del código generado
- [ ] Diseñar sistema de templates para generación
- [ ] Crear pipeline de compilación modular
- [ ] Implementar sistema de caché de compilación

### 4.2 Generador de handlers HTTP
- [ ] Generar struct del componente desde AST
- [ ] Generar métodos lifecycle (Mount, BeforeMount, etc.)
- [ ] Generar handlers para endpoints HTMX
- [ ] Implementar inyección de dependencias
- [ ] Generar routing automático

### 4.3 Compilador de templates
- [ ] Convertir template HTML a Go templates
- [ ] Procesar directivas especiales (auth, layout)
- [ ] Integrar sintaxis Go template
- [ ] Optimizar templates para performance
- [ ] Generar funciones de renderizado

### 4.4 Procesador de estilos
- [ ] Integrar con Tailwind CSS
- [ ] Procesar CSS scoped
- [ ] Generar archivo CSS final
- [ ] Implementar tree-shaking de CSS
- [ ] Soportar PostCSS plugins

### 4.5 Sistema de imports y dependencias
- [ ] Resolver imports automáticamente
- [ ] Detectar dependencias entre componentes
- [ ] Generar go.mod para el proyecto
- [ ] Manejar componentes compartidos
- [ ] Optimizar imports no utilizados

### 4.6 Generación de código auxiliar
- [ ] Generar main.go para el servidor
- [ ] Crear sistema de registro de componentes
- [ ] Generar middleware desde configuración
- [ ] Crear helpers para HTMX
- [ ] Generar tipos TypeScript (opcional)

## Criterios de Aceptación

1. **Código generado válido**
   - El código Go generado debe compilar sin errores
   - Debe seguir las convenciones de Go
   - Debe ser legible y debuggeable

2. **Estructura de salida**
   ```
   .gox/
   ├── generated/
   │   ├── pages/
   │   │   └── index.gen.go
   │   ├── components/
   │   │   └── user_card.gen.go
   │   ├── routes.gen.go
   │   ├── templates.gen.go
   │   └── assets/
   │       └── styles.css
   └── cache/
       └── compile.cache
   ```

3. **Ejemplo de código generado**
   ```go
   // .gox/generated/pages/index.gen.go
   package pages
   
   import (
       "github.com/gox-framework/gox"
       "my-app/common/types"
   )
   
   type IndexPage struct {
       Title string
       Users []types.User
   }
   
   func (p *IndexPage) Mount(ctx *gox.Context) error {
       p.Title = "Welcome to GOX"
       // Original Mount code from .gox file
       return nil
   }
   
   func (p *IndexPage) Render(ctx *gox.Context) error {
       return ctx.RenderTemplate("pages/index", p)
   }
   
   func init() {
       gox.RegisterPage("/", &IndexPage{})
   }
   ```

4. **Templates compilados**
   ```go
   // .gox/generated/templates.gen.go
   var templates = map[string]string{
       "pages/index": `
       <div class="container">
           <h1>{{.Title}}</h1>
           {{range .Users}}
               <div>{{.Name}}</div>
           {{end}}
       </div>
       `,
   }
   ```

## Tests Necesarios

### Tests Unitarios

1. **Test de generación de handlers**
```go
func TestGenerateHandler(t *testing.T) {
    ast := &GoxAST{
        Go: &GoNode{
            StructName: "HomePage",
            Methods: []Method{
                {Name: "Mount", Body: "p.Title = \"Test\""},
            },
        },
    }
    
    code, err := compiler.GenerateHandler(ast)
    assert.NoError(t, err)
    assert.Contains(t, code, "type HomePage struct")
    assert.Contains(t, code, "func (p *HomePage) Mount")
}
```

2. **Test de compilación de templates**
```go
func TestCompileTemplate(t *testing.T) {
    template := &TemplateNode{
        Content: `<div>{{.Title}}</div>`,
        Auth: "required",
    }
    
    compiled, err := compiler.CompileTemplate(template)
    assert.NoError(t, err)
    assert.Contains(t, compiled, "{{.Title}}")
    assert.Contains(t, compiled, "gox.RequireAuth")
}
```

3. **Test de procesamiento de estilos**
```go
func TestProcessStyles(t *testing.T) {
    styles := []*StyleNode{
        {Content: ".test { color: red; }", Scoped: true},
        {Content: "@apply bg-blue-500;", Type: "tailwind"},
    }
    
    css, err := compiler.ProcessStyles(styles)
    assert.NoError(t, err)
    assert.Contains(t, css, ".test")
    assert.Contains(t, css, "background-color")
}
```

### Tests de Integración

1. **Test E2E de compilación**
```go
func TestCompileFullProject(t *testing.T) {
    // Compilar proyecto de prueba
    err := compiler.CompileProject("testdata/sample-project")
    assert.NoError(t, err)
    
    // Verificar archivos generados
    assert.FileExists(t, ".gox/generated/routes.gen.go")
    assert.FileExists(t, ".gox/generated/templates.gen.go")
    
    // Verificar que compila
    cmd := exec.Command("go", "build", "./...")
    err = cmd.Run()
    assert.NoError(t, err)
}
```

2. **Test de hot reload**
```go
func TestIncrementalCompilation(t *testing.T) {
    compiler := NewCompiler(WithCache(true))
    
    // Primera compilación
    start := time.Now()
    err := compiler.CompileFile("pages/index.gox")
    firstTime := time.Since(start)
    
    // Segunda compilación (debe usar caché)
    start = time.Now()
    err = compiler.CompileFile("pages/index.gox")
    secondTime := time.Since(start)
    
    assert.Less(t, secondTime, firstTime/2)
}
```

### Benchmarks

```go
func BenchmarkCompileGoxFile(b *testing.B) {
    ast, _ := parser.ParseFile("testdata/complex.gox")
    compiler := NewCompiler()
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        compiler.Compile(ast)
    }
}
```

## Definición de Done

- [ ] Compilador completo generando código Go válido
- [ ] Sistema de caché funcionando
- [ ] Templates compilados y optimizados
- [ ] Estilos procesados correctamente
- [ ] Tests con cobertura > 85%
- [ ] Documentación del código generado
- [ ] Ejemplos de proyectos compilados
- [ ] Performance < 100ms por archivo

## Notas Adicionales

- El código generado debe ser idempotente
- Considerar generar source maps para debugging
- Implementar modo verbose para ver código generado
- El compilador debe ser incremental para hot reload
- Pensar en compatibilidad con herramientas de Go existentes
- Considerar generar documentación automática del código