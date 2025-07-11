# Task 008: Sistema de Componentes

## Descripción
Implementar un sistema de componentes reutilizables que permita crear, compartir y componer componentes .gox, con soporte para props, slots, eventos, y ciclo de vida.

## Prioridad
Alta

## Estimación
5-6 días

## Dependencias
- Task 003: Parser de archivos .gox
- Task 004: Compilador .gox a Go
- Task 007: Integración HTMX

## Subtasks

### 8.1 Definición de componentes
- [ ] Estructura base de componentes
- [ ] Sistema de props con tipos
- [ ] Valores por defecto de props
- [ ] Props requeridos vs opcionales
- [ ] Validación de props en runtime

### 8.2 Sistema de slots
- [ ] Implementar slots con nombre
- [ ] Slot por defecto
- [ ] Slots múltiples
- [ ] Fallback content para slots
- [ ] Scoped slots con datos

### 8.3 Composición de componentes
- [ ] Importar componentes en otros componentes
- [ ] Componentes anidados
- [ ] Herencia de componentes
- [ ] Mixins/traits para compartir lógica
- [ ] Componentes dinámicos

### 8.4 Ciclo de vida
- [ ] BeforeMount hook
- [ ] Mount hook
- [ ] AfterMount hook
- [ ] BeforeUpdate hook
- [ ] AfterUpdate hook
- [ ] Unmount hook

### 8.5 Sistema de eventos
- [ ] Emitir eventos desde componentes
- [ ] Escuchar eventos en padres
- [ ] Bubbling de eventos
- [ ] Eventos personalizados HTMX
- [ ] Event modifiers

### 8.6 Registro y resolución
- [ ] Registro global de componentes
- [ ] Resolución de dependencias
- [ ] Lazy loading de componentes
- [ ] Componentes asíncronos
- [ ] Tree-shaking de componentes no usados

## Criterios de Aceptación

1. **Definición de componente**
   ```gox
   <!-- components/user-card.gox -->
   <template>
     <div class="user-card" {{if .Featured}}featured{{end}}>
       <img src="{{.User.Avatar}}" alt="{{.User.Name}}">
       <h3>{{.User.Name}}</h3>
       <p>{{.User.Bio}}</p>
       
       <slot name="actions">
         <button>Default Action</button>
       </slot>
     </div>
   </template>
   
   <go>
   type UserCard struct {
       User     User   `props:"required"`
       Featured bool   `props:"default:false"`
       OnClick  func() `props:"event"`
   }
   
   func (c *UserCard) BeforeMount(ctx *gox.Context) error {
       // Validar props
       if c.User.ID == 0 {
           return errors.New("User ID is required")
       }
       return nil
   }
   </go>
   ```

2. **Uso de componentes**
   ```gox
   <!-- pages/users.gox -->
   <template>
     <div class="users-page">
       {{range .Users}}
         <gox-component 
           src="user-card" 
           :user="." 
           :featured="{{eq .ID $.FeaturedID}}"
           @click="handleUserClick">
           
           <template slot="actions">
             <button hx-get="/users/{{.ID}}/edit">Edit</button>
             <button hx-delete="/users/{{.ID}}">Delete</button>
           </template>
         </gox-component>
       {{end}}
     </div>
   </template>
   ```

3. **Componentes dinámicos**
   ```gox
   <template>
     <gox-component 
       :src=".ComponentName"
       :props=".ComponentProps">
     </gox-component>
   </template>
   ```

4. **Sistema de eventos**
   ```go
   // En el componente hijo
   func (c *ChildComponent) HandleClick(ctx *gox.Context) error {
       c.Emit("user-selected", map[string]interface{}{
           "userId": c.User.ID,
       })
       return nil
   }
   
   // En el componente padre
   func (p *ParentComponent) OnUserSelected(data map[string]interface{}) {
       userId := data["userId"].(int)
       // Manejar evento
   }
   ```

## Tests Necesarios

### Tests Unitarios

1. **Test props validation**
```go
func TestComponentPropsValidation(t *testing.T) {
    comp := &UserCard{}
    
    // Sin user (requerido)
    err := comp.ValidateProps()
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "User is required")
    
    // Con user válido
    comp.User = User{ID: 1, Name: "John"}
    err = comp.ValidateProps()
    assert.NoError(t, err)
    
    // Valor por defecto
    assert.False(t, comp.Featured)
}
```

2. **Test slots**
```go
func TestComponentSlots(t *testing.T) {
    template := `
    <gox-component src="card">
      <template slot="header">
        <h1>Title</h1>
      </template>
      <template slot="footer">
        <p>Footer</p>
      </template>
      Default content
    </gox-component>
    `
    
    comp := ParseComponent(template)
    slots := comp.GetSlots()
    
    assert.Equal(t, "<h1>Title</h1>", slots["header"])
    assert.Equal(t, "<p>Footer</p>", slots["footer"])
    assert.Equal(t, "Default content", slots["default"])
}
```

3. **Test ciclo de vida**
```go
func TestComponentLifecycle(t *testing.T) {
    var calls []string
    
    comp := &TestComponent{
        BeforeMount: func() error {
            calls = append(calls, "before-mount")
            return nil
        },
        Mount: func() error {
            calls = append(calls, "mount")
            return nil
        },
        AfterMount: func() error {
            calls = append(calls, "after-mount")
            return nil
        },
    }
    
    err := comp.Render(context.Background())
    assert.NoError(t, err)
    
    expected := []string{"before-mount", "mount", "after-mount"}
    assert.Equal(t, expected, calls)
}
```

### Tests de Integración

1. **Test componentes anidados**
```go
func TestNestedComponents(t *testing.T) {
    // Crear app con componentes
    app := NewApp()
    app.RegisterComponent("user-card", &UserCard{})
    app.RegisterComponent("user-list", &UserList{})
    
    // Renderizar página con componentes anidados
    page := `
    <gox-component src="user-list" :users="{{.Users}}">
      <template slot="item" slot-scope="user">
        <gox-component src="user-card" :user="user" />
      </template>
    </gox-component>
    `
    
    result, err := app.RenderPage(page, map[string]interface{}{
        "Users": []User{{ID: 1}, {ID: 2}},
    })
    
    assert.NoError(t, err)
    assert.Contains(t, result, "user-card")
    assert.Equal(t, 2, strings.Count(result, "user-card"))
}
```

2. **Test eventos entre componentes**
```go
func TestComponentEvents(t *testing.T) {
    var eventData map[string]interface{}
    
    parent := &ParentComponent{
        OnChildEvent: func(data map[string]interface{}) {
            eventData = data
        },
    }
    
    child := &ChildComponent{
        Parent: parent,
    }
    
    // Emitir evento desde hijo
    child.Emit("test-event", map[string]interface{}{
        "value": "test",
    })
    
    assert.Equal(t, "test", eventData["value"])
}
```

### Tests de Performance

```go
func BenchmarkComponentRendering(b *testing.B) {
    app := NewApp()
    app.RegisterComponent("item", &Item{})
    
    // Template con muchos componentes
    template := `
    {{range .Items}}
      <gox-component src="item" :data="." />
    {{end}}
    `
    
    items := make([]interface{}, 100)
    for i := range items {
        items[i] = map[string]interface{}{"id": i}
    }
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        app.RenderTemplate(template, map[string]interface{}{
            "Items": items,
        })
    }
}
```

## Definición de Done

- [ ] Sistema de componentes completo
- [ ] Props con validación funcionando
- [ ] Slots implementados
- [ ] Ciclo de vida completo
- [ ] Sistema de eventos funcionando
- [ ] Componentes anidados soportados
- [ ] Tests con cobertura > 85%
- [ ] Documentación y ejemplos

## Notas Adicionales

- Los componentes deben ser eficientes en memoria
- Considerar lazy loading para componentes grandes
- El sistema debe ser extensible
- Mantener compatibilidad con HTMX
- Pensar en componentes del lado del servidor vs cliente
- Documentar mejores prácticas de componentes