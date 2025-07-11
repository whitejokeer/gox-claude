# Task 013: Generadores de Código

## Descripción
Implementar generadores de código que automaticen la creación de páginas, componentes, servicios, y otros elementos del framework, con templates personalizables y scaffolding inteligente.

## Prioridad
Alta

## Estimación
4-5 días

## Dependencias
- Task 002: CLI básico
- Task 012: Sistema de configuración

## Subtasks

### 13.1 Sistema base de generadores
- [ ] Framework para crear generadores
- [ ] Sistema de templates con variables
- [ ] Engine de templating (text/template)
- [ ] Validación de entrada
- [ ] Post-processing hooks

### 13.2 Generador de páginas
- [ ] Comando `gox generate page`
- [ ] Templates para diferentes tipos de página
- [ ] Routing automático
- [ ] Layout integration
- [ ] Metadata generation

### 13.3 Generador de componentes
- [ ] Comando `gox generate component`
- [ ] Props scaffolding
- [ ] Diferentes variantes (button, card, form, etc.)
- [ ] Story generation automática
- [ ] CSS scaffolding

### 13.4 Generador de servicios
- [ ] Comando `gox generate service`
- [ ] API endpoints scaffolding
- [ ] Repository pattern
- [ ] Model generation
- [ ] Migration generation

### 13.5 Generador de middleware
- [ ] Comando `gox generate middleware`
- [ ] Templates para auth, logging, etc.
- [ ] Integration con routing
- [ ] Error handling patterns
- [ ] Testing scaffolding

### 13.6 Templates personalizables
- [ ] Sistema de templates por proyecto
- [ ] Override de templates default
- [ ] Variables de contexto
- [ ] Conditional generation
- [ ] Custom helpers

## Criterios de Aceptación

1. **Generador de páginas funcionando**
   ```bash
   # Comando básico
   gox generate page dashboard
   # Genera: pages/dashboard.gox
   
   # Con opciones
   gox generate page users/profile --auth --layout=admin
   # Genera: pages/users/profile.gox con auth y layout
   
   # Con props
   gox generate page product/[id] --props="productId:string,tab:string"
   ```

2. **Template generado**
   ```gox
   <!-- pages/dashboard.gox -->
   <template auth="required" layout="admin">
     <div class="dashboard-page">
       <h1>Dashboard</h1>
       <p>Welcome {{.User.Name}}</p>
       
       <!-- Add your content here -->
     </div>
   </template>
   
   <go>
   type DashboardPage struct {
       User User
   }
   
   func (d *DashboardPage) Mount(ctx *gox.Context) error {
       d.User = ctx.Get("user").(User)
       return nil
   }
   </go>
   
   <style scoped>
   .dashboard-page {
       @apply container mx-auto py-8;
   }
   </style>
   ```

3. **Generador de componentes**
   ```bash
   # Componente básico
   gox generate component user-card --props="user:User,featured:bool"
   
   # Con variante
   gox generate component button --variant=form
   
   # Con story
   gox generate component modal --with-story
   ```

4. **Generador de servicios**
   ```bash
   # Servicio completo
   gox generate service users --crud --model=User
   
   # Con endpoints específicos
   gox generate service auth --endpoints="login,logout,refresh"
   ```

## Tests Necesarios

### Tests Unitarios

1. **Test generador de páginas**
```go
func TestPageGenerator(t *testing.T) {
    gen := NewPageGenerator()
    
    options := PageOptions{
        Name:   "dashboard",
        Auth:   true,
        Layout: "admin",
        Props:  []Prop{{Name: "userId", Type: "string"}},
    }
    
    files, err := gen.Generate(options)
    assert.NoError(t, err)
    assert.Len(t, files, 1)
    
    content := files[0].Content
    assert.Contains(t, content, `auth="required"`)
    assert.Contains(t, content, `layout="admin"`)
    assert.Contains(t, content, "UserId string")
}
```

2. **Test generador de componentes**
```go
func TestComponentGenerator(t *testing.T) {
    gen := NewComponentGenerator()
    
    options := ComponentOptions{
        Name: "UserCard",
        Props: []Prop{
            {Name: "user", Type: "User", Required: true},
            {Name: "featured", Type: "bool", Default: "false"},
        },
        WithStory: true,
    }
    
    files, err := gen.Generate(options)
    assert.NoError(t, err)
    assert.Len(t, files, 1)
    
    content := files[0].Content
    assert.Contains(t, content, "type UserCard struct")
    assert.Contains(t, content, "User User")
    assert.Contains(t, content, "Featured bool")
    assert.Contains(t, content, `default:"false"`)
}
```

3. **Test sistema de templates**
```go
func TestTemplateSystem(t *testing.T) {
    templates := map[string]string{
        "component": `
type {{.Name}} struct {
    {{range .Props}}
    {{.Name | title}} {{.Type}} ` + "`" + `{{if .Required}}props:"required"{{end}}` + "`" + `
    {{end}}
}
`,
    }
    
    engine := NewTemplateEngine(templates)
    
    data := map[string]interface{}{
        "Name": "TestComponent",
        "Props": []Prop{
            {Name: "title", Type: "string", Required: true},
        },
    }
    
    result, err := engine.Render("component", data)
    assert.NoError(t, err)
    assert.Contains(t, result, "type TestComponent struct")
    assert.Contains(t, result, `props:"required"`)
}
```

### Tests de Integración

1. **Test generación E2E**
```go
func TestEndToEndGeneration(t *testing.T) {
    // Crear proyecto temporal
    tmpDir := createTempProject(t)
    defer os.RemoveAll(tmpDir)
    
    // Generar página
    err := runCLI([]string{
        "generate", "page", "users",
        "--auth",
        "--props=page:int,search:string",
    }, tmpDir)
    assert.NoError(t, err)
    
    // Verificar archivo generado
    pageFile := filepath.Join(tmpDir, "pages", "users.gox")
    assert.FileExists(t, pageFile)
    
    content, err := os.ReadFile(pageFile)
    assert.NoError(t, err)
    assert.Contains(t, string(content), `auth="required"`)
    assert.Contains(t, string(content), "Page int")
    assert.Contains(t, string(content), "Search string")
    
    // Verificar que compila
    err = compilePage(pageFile)
    assert.NoError(t, err)
}
```

2. **Test templates personalizados**
```go
func TestCustomTemplates(t *testing.T) {
    tmpDir := createTempProject(t)
    defer os.RemoveAll(tmpDir)
    
    // Crear template personalizado
    customTemplate := `
<!-- Custom page template -->
<template>
  <div class="custom-{{.Name}}">
    <h1>{{.Name | title}}</h1>
  </div>
</template>
`
    
    templateDir := filepath.Join(tmpDir, ".gox", "templates")
    os.MkdirAll(templateDir, 0755)
    os.WriteFile(filepath.Join(templateDir, "page.gox"), []byte(customTemplate), 0644)
    
    // Generar con template personalizado
    err := runCLI([]string{"generate", "page", "custom"}, tmpDir)
    assert.NoError(t, err)
    
    // Verificar uso del template
    content, err := os.ReadFile(filepath.Join(tmpDir, "pages", "custom.gox"))
    assert.NoError(t, err)
    assert.Contains(t, string(content), "Custom page template")
    assert.Contains(t, string(content), "custom-custom")
}
```

### Tests de CLI

```go
func TestCLIGeneration(t *testing.T) {
    tests := []struct {
        name     string
        args     []string
        wantFile string
        wantErr  bool
    }{
        {
            name:     "page simple",
            args:     []string{"generate", "page", "about"},
            wantFile: "pages/about.gox",
        },
        {
            name:     "component with props",
            args:     []string{"generate", "component", "button", "--props=label:string,disabled:bool"},
            wantFile: "components/button.gox",
        },
        {
            name:    "invalid component name",
            args:    []string{"generate", "component", "123invalid"},
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            tmpDir := createTempProject(t)
            defer os.RemoveAll(tmpDir)
            
            err := runCLI(tt.args, tmpDir)
            
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
                assert.FileExists(t, filepath.Join(tmpDir, tt.wantFile))
            }
        })
    }
}
```

## Definición de Done

- [ ] Generadores para páginas, componentes, servicios
- [ ] Sistema de templates flexible
- [ ] CLI commands completos
- [ ] Templates personalizables
- [ ] Validación de entrada
- [ ] Tests con cobertura > 85%
- [ ] Documentación de uso

## Notas Adicionales

- Los generadores deben seguir convenciones del framework
- Los templates deben ser fáciles de personalizar
- Considerar generadores para tests automáticamente
- Los nombres generados deben seguir convenciones de Go
- Pensar en generadores para migraciones
- Documentar cómo crear generadores personalizados