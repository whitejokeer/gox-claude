# Convenciones GOX Framework

## 🎯 Principios Fundamentales

GOX sigue la filosofía de Go: **simplicidad y claridad sobre magia**.

### ❌ NO hay:
- Lifecycle methods (BeforeMount, AfterMount, etc.)
- Herencia forzada (no page.Base, component.Base)
- Hooks estilo React/Vue
- Virtual DOM
- Estado global mágico

### ✅ SÍ hay:
- Structs simples de Go
- Constructores explícitos
- Handlers HTTP estándar
- HTMX para interactividad
- Server-Side Rendering puro

## 📝 Estructura de Archivos .gox

### 1. **Páginas** (`app/pages/*.gox`)

```gox
<template auth="optional|required" layout="nombre-layout">
  <!-- HTML + Go templates + HTMX -->
</template>

<go>
package pages

// Struct simple - sin herencia
type HomePage struct {
    Title   string
    Posts   []Post
    IsAdmin bool
}

// Constructor - se llama al renderizar la página
func NewHomePage(r *http.Request) (*HomePage, error) {
    // Cargar datos, verificar permisos, etc.
    return &HomePage{
        Title: "Welcome",
        Posts: loadPosts(),
    }, nil
}

// Handlers HTMX - métodos del struct
func (p *HomePage) HandleSearch(w http.ResponseWriter, r *http.Request) {
    // Retorna HTML parcial para HTMX
}
</go>

<style scoped>
/* CSS con scope automático */
</style>
```

### 2. **Componentes** (`app/components/*.gox`)

```gox
<template>
  <!-- HTML del componente -->
</template>

<go>
package components

// Props son campos públicos del struct
type UserCard struct {
    ID     int    `gox:"required"`
    Name   string `gox:"required"`  
    Email  string
    Avatar string `gox:"default:/default-avatar.png"`
}

// Validación opcional
func (c *UserCard) Validate() error {
    if c.Name == "" {
        return errors.New("name is required")
    }
    return nil
}

// Handlers para interacciones
func (c *UserCard) HandleClick(w http.ResponseWriter, r *http.Request) {
    // Acción del componente
}
</go>
```

## 🔧 Convenciones de Código

### Tipos Principales

1. **Páginas**: El struct principal representa el estado de la página
2. **Componentes**: El struct principal define las props
3. **Campos públicos**: Son las variables disponibles en el template

### Constructores

- Páginas: `func NewXxxPage(r *http.Request) (*XxxPage, error)`
- Componentes: Se crean automáticamente con las props del template padre

### Handlers HTTP

- Deben seguir la firma estándar: `func (receiver) HandlerName(w http.ResponseWriter, r *http.Request)`
- Nombres descriptivos: `HandleSubmit`, `HandleDelete`, `HandleSearch`
- Retornan HTML parcial para requests HTMX

### Tags de Struct

```go
type Component struct {
    // Validación GOX
    Name string `gox:"required"`
    Age  int    `gox:"min=0,max=150"`
    
    // Serialización JSON
    Data string `json:"data,omitempty"`
    
    // Valor por defecto
    Icon string `gox:"default=/icons/user.svg"`
}
```

## 🚫 Anti-patrones a Evitar

### ❌ NO hacer:
```go
// Herencia forzada
type MyPage struct {
    framework.BasePage // ❌
}

// Lifecycle methods
func (p *Page) BeforeMount() { } // ❌
func (p *Page) AfterRender() { } // ❌

// Estado global mágico
func (p *Page) SetState(key, value) { } // ❌
```

### ✅ Hacer en su lugar:
```go
// Struct simple
type MyPage struct {
    Title string
    Data  []Item
}

// Constructor explícito
func NewMyPage(r *http.Request) (*MyPage, error) {
    // Inicialización clara
}

// Handlers HTTP estándar
func (p *MyPage) HandleAction(w http.ResponseWriter, r *http.Request) {
    // Lógica explícita
}
```

## 🎨 Templates

### Variables Disponibles

En el template tienes acceso a:
- Todos los campos públicos del struct
- Funciones helper registradas
- Sintaxis estándar de Go templates

### HTMX Integration

```html
<!-- Actualización parcial -->
<button hx-post="/users/{{.ID}}/delete" 
        hx-target="#user-list"
        hx-confirm="¿Estás seguro?">
    Eliminar
</button>

<!-- Polling -->
<div hx-get="/notifications" 
     hx-trigger="every 30s">
</div>

<!-- Validación en tiempo real -->
<input name="email" 
       hx-post="/validate/email" 
       hx-trigger="blur"
       hx-target="#email-error">
```

## 🗂 Organización del Proyecto

```
app/
├── pages/          # Páginas - rutas automáticas
├── components/     # Componentes reutilizables  
├── shared/         
│   ├── ui/        # Componentes compartidos (shared-*)
│   └── layouts/   # Layouts para páginas
└── routing/       # Rutas manuales si se necesitan
```

## 🔍 Resolución de Componentes

- `<user-card>` → `app/components/user-card.gox`
- `<shared-button>` → `app/shared/ui/button.gox`
- Los componentes se resuelven automáticamente por convención

## 💡 Mejores Prácticas

1. **Mantén los componentes simples**: Una sola responsabilidad
2. **Props explícitas**: Usa structs con campos bien definidos
3. **Sin estado oculto**: Todo el estado en el struct principal
4. **Handlers puros**: Sin side effects no documentados
5. **Errores explícitos**: Retorna errores, no los ocultes
6. **Composición sobre herencia**: Usa embebimiento de structs cuando sea necesario

## 🚀 Ejemplo Completo

Ver `/examples/parser-demo/test.gox` para un ejemplo real siguiendo estas convenciones.