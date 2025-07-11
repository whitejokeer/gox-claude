# Task 011: Storybook Integrado

## Descripción
Implementar un sistema de Storybook nativo integrado que auto-descubra componentes .gox, genere documentación automática, y proporcione un entorno de desarrollo interactivo para componentes.

## Prioridad
Alta

## Estimación
5-6 días

## Dependencias
- Task 008: Sistema de componentes
- Task 006: Servidor de desarrollo

## Subtasks

### 11.1 Auto-discovery de componentes
- [ ] Escanear estructura de proyecto
- [ ] Detectar archivos .gox automáticamente
- [ ] Organizar por estructura de carpetas
- [ ] Filtrar componentes vs páginas
- [ ] Detectar stories definidas manualmente

### 11.2 Generación automática de stories
- [ ] Extraer props de componentes Go
- [ ] Generar controles automáticamente
- [ ] Crear variantes por defecto
- [ ] Detectar tipos de datos (string, int, bool, etc.)
- [ ] Generar documentación desde comentarios

### 11.3 Interface de usuario
- [ ] Sidebar con navegación jerárquica
- [ ] Área de preview de componentes
- [ ] Panel de controles dinámicos
- [ ] Visor de código fuente
- [ ] Documentación automática

### 11.4 Sistema de controles
- [ ] Controls automáticos basados en tipos
- [ ] Selectores para enums
- [ ] Sliders para números
- [ ] Text inputs para strings
- [ ] Checkboxes para booleans
- [ ] Object editors para structs

### 11.5 Addons y extensiones
- [ ] Viewport addon (responsive testing)
- [ ] Backgrounds addon
- [ ] Actions addon (eventos)
- [ ] Accessibility addon
- [ ] Performance addon

### 11.6 Integración con hot reload
- [ ] Sincronizar cambios en tiempo real
- [ ] Recargar components modificados
- [ ] Preservar estado de controles
- [ ] Mostrar errores de compilación
- [ ] Live documentation updates

## Criterios de Aceptación

1. **Auto-discovery funcionando**
   ```
   src/
   ├── components/
   │   ├── button.gox        → Components/Button
   │   ├── card.gox          → Components/Card
   │   └── forms/
   │       └── input.gox     → Components/Forms/Input
   ├── pages/
   │   └── dashboard.gox     → Pages/Dashboard
   └── services/
       └── users/
           └── components/
               └── profile.gox → Services/Users/Profile
   ```

2. **Stories automáticas**
   ```gox
   <!-- components/button.gox -->
   <template>
     <button class="btn {{.Variant}}" {{if .Disabled}}disabled{{end}}>
       {{.Label}}
     </button>
   </template>
   
   <go>
   // Button component for user interactions
   type Button struct {
       Label    string `story:"Button text" default:"Click me"`
       Variant  string `story:"Button variant" options:"primary,secondary,danger" default:"primary"`
       Disabled bool   `story:"Disabled state" default:"false"`
       Size     string `story:"Button size" options:"sm,md,lg" default:"md"`
   }
   </go>
   
   <!-- Genera automáticamente controles para cada prop -->
   ```

3. **Interface completa**
   ```
   Storybook UI:
   ├── Sidebar
   │   ├── 📁 Components
   │   │   ├── 📄 Button
   │   │   ├── 📄 Card
   │   │   └── 📁 Forms
   │   │       └── 📄 Input
   │   ├── 📁 Pages
   │   └── 📁 Services
   ├── Preview Area
   │   └── [Component rendering]
   ├── Controls Panel
   │   ├── Label: [text input]
   │   ├── Variant: [select: primary/secondary/danger]
   │   ├── Disabled: [checkbox]
   │   └── Size: [select: sm/md/lg]
   └── Docs/Code Panel
   ```

4. **Addons funcionando**
   ```javascript
   // Viewport addon
   const viewports = [
       { name: 'Mobile', width: 375, height: 667 },
       { name: 'Tablet', width: 768, height: 1024 },
       { name: 'Desktop', width: 1440, height: 900 }
   ];
   
   // Actions addon
   button.addEventListener('click', action('button-clicked'));
   ```

## Tests Necesarios

### Tests Unitarios

1. **Test auto-discovery**
```go
func TestComponentDiscovery(t *testing.T) {
    discoverer := NewStorybook()
    
    components, err := discoverer.DiscoverComponents("testdata/project")
    assert.NoError(t, err)
    
    expected := []Component{
        {Name: "Button", Path: "components/button.gox", Category: "Components"},
        {Name: "Card", Path: "components/card.gox", Category: "Components"},
        {Name: "Input", Path: "components/forms/input.gox", Category: "Components/Forms"},
    }
    
    assert.Equal(t, expected, components)
}
```

2. **Test generación de controles**
```go
func TestControlGeneration(t *testing.T) {
    component := &Button{
        Label:    "Click me",
        Variant:  "primary",
        Disabled: false,
    }
    
    controls := GenerateControls(component)
    
    assert.Contains(t, controls, Control{
        Name: "Label",
        Type: "text",
        Default: "Click me",
    })
    
    assert.Contains(t, controls, Control{
        Name: "Variant",
        Type: "select",
        Options: []string{"primary", "secondary", "danger"},
        Default: "primary",
    })
}
```

3. **Test props extraction**
```go
func TestPropsExtraction(t *testing.T) {
    sourceCode := `
    type Button struct {
        Label string ` + "`story:\"Button text\" default:\"Click me\"`" + `
        Size  string ` + "`story:\"Size\" options:\"sm,md,lg\"`" + `
    }`
    
    props, err := ExtractProps(sourceCode)
    assert.NoError(t, err)
    
    assert.Equal(t, "Button text", props["Label"].Description)
    assert.Equal(t, "Click me", props["Label"].Default)
    assert.Equal(t, []string{"sm", "md", "lg"}, props["Size"].Options)
}
```

### Tests de Integración

1. **Test servidor Storybook**
```go
func TestStorybookServer(t *testing.T) {
    storybook := NewStorybook(Config{
        Port: 6007,
        Root: "testdata/project",
    })
    
    go storybook.Start()
    defer storybook.Stop()
    
    // Esperar que inicie
    waitForServer(t, "http://localhost:6007")
    
    // Verificar página principal
    resp, err := http.Get("http://localhost:6007")
    assert.NoError(t, err)
    assert.Equal(t, 200, resp.StatusCode)
    
    // Verificar API de componentes
    resp, err = http.Get("http://localhost:6007/api/components")
    assert.NoError(t, err)
    
    var components []Component
    json.NewDecoder(resp.Body).Decode(&components)
    assert.NotEmpty(t, components)
}
```

2. **Test hot reload**
```go
func TestStorybookHotReload(t *testing.T) {
    storybook := NewStorybook()
    reloaded := make(chan string, 1)
    
    storybook.OnReload(func(component string) {
        reloaded <- component
    })
    
    go storybook.Start()
    defer storybook.Stop()
    
    // Modificar componente
    modifyFile("components/button.gox")
    
    select {
    case component := <-reloaded:
        assert.Equal(t, "button", component)
    case <-time.After(2 * time.Second):
        t.Fatal("Hot reload not triggered")
    }
}
```

3. **Test addons**
```go
func TestStorybookAddons(t *testing.T) {
    storybook := NewStorybook()
    
    // Registrar addon
    storybook.RegisterAddon("viewport", ViewportAddon{
        Viewports: []Viewport{
            {Name: "Mobile", Width: 375},
            {Name: "Desktop", Width: 1440},
        },
    })
    
    // Verificar addon cargado
    addons := storybook.GetAddons()
    assert.Contains(t, addons, "viewport")
    
    // Verificar configuración
    config := storybook.GetAddonConfig("viewport")
    assert.Contains(t, config["viewports"], Viewport{Name: "Mobile"})
}
```

### Tests de Performance

```go
func BenchmarkComponentDiscovery(b *testing.B) {
    discoverer := NewStorybook()
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        discoverer.DiscoverComponents("testdata/large-project")
    }
}

func BenchmarkStoryGeneration(b *testing.B) {
    component := createLargeComponent()
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        GenerateStory(component)
    }
}
```

## Definición de Done

- [ ] Auto-discovery de componentes funcionando
- [ ] Stories automáticas generadas
- [ ] Interface de usuario completa
- [ ] Sistema de controles dinámicos
- [ ] Addons básicos implementados
- [ ] Hot reload integrado
- [ ] Tests con cobertura > 85%
- [ ] Documentación del sistema

## Notas Adicionales

- La UI debe ser responsive y moderna
- Considerar usar Web Components para la interfaz
- El sistema debe ser extensible con plugins
- Mantener compatibilidad con Storybook estándar cuando posible
- Optimizar para proyectos grandes con muchos componentes
- Documentar cómo crear addons personalizados