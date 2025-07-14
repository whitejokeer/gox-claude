# Parser GOX - Guía de Uso

El parser de GOX es capaz de analizar archivos `.gox` y extraer toda la información estructural, componentes, handlers HTTP, y dependencias.

## 🚀 Inicio Rápido

### 1. Parsear un archivo

```go
package main

import (
    "fmt"
    "github.com/gox-framework/gox/internal/parser"
)

func main() {
    // Parsear desde archivo
    ast, err := parser.ParseFile("mi-componente.gox")
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Archivo parseado: %s\n", ast.Path)
    fmt.Printf("Es componente: %v\n", ast.IsComponent())
    fmt.Printf("Componentes usados: %d\n", len(ast.Components))
}
```

### 2. Parsear desde string

```go
content := `<template>
  <user-card name="{{.Name}}" />
</template>

<go>
package pages
type UserPage struct {
    Name string
}
</go>`

ast, err := parser.ParseString(content)
```

## 📋 Qué Detecta el Parser

### ✅ Secciones del archivo
- **Template**: HTML con sintaxis Go template y atributos HTMX
- **Go**: Código Go con validación de sintaxis
- **Style**: CSS con soporte para scoped styles

### ✅ Análisis de Template
- Variables Go template (`{{.Variable}}`)
- Componentes custom (`<user-card>`, `<shared-button>`)
- Atributos HTMX (`hx-post`, `hx-target`, etc.)
- Atributos especiales (`auth`, `layout`)

### ✅ Análisis de Código Go
- Imports automáticos
- Struct principal del componente/página
- Props del componente con validación de tags
- Handlers HTTP con detección automática
- Métodos HTMX

### ✅ Análisis de Estilos
- Detección de scoped styles
- Validación básica de CSS
- Tipo de estilo (CSS, Tailwind, SCSS)

### ✅ Resolución de Componentes
- `<user-card>` → `components/user-card.gox`
- `<shared-button>` → `shared/ui/button.gox`
- Props pasados a cada componente
- Conteo de uso de componentes

## 🔍 Casos de Uso Principales

### 1. **Herramientas de Build**
Analizar dependencias entre componentes para determinar orden de compilación.

### 2. **IDEs y Editores**
Proporcionar autocompletado, validación en tiempo real, y navegación.

### 3. **Generadores de Código**
Crear handlers, rutas, y tipos automáticamente basado en la estructura.

### 4. **Análisis Estático**
Detectar componentes no utilizados, props requeridos faltantes, etc.

### 5. **Hot Reload**
Determinar qué archivos recargar cuando cambian las dependencias.

## 📊 Rendimiento

- **Archivos simples**: ~28μs
- **Archivos complejos**: ~263μs  
- **Detección de componentes**: ~31μs
- **Tokenización**: ~2.4μs

## 🎯 Estructura del AST

```go
type GoxFile struct {
    Path       string                 // Ruta del archivo
    Template   *TemplateNode          // Sección template
    Go         *GoNode                // Sección go
    Styles     []*StyleNode           // Secciones style
    Components []ComponentDependency  // Componentes usados
    Metadata   map[string]interface{} // Metadata adicional
}
```

### TemplateNode
```go
type TemplateNode struct {
    Content  string    // HTML content
    Auth     string    // "required", "role:admin"
    Layout   string    // Layout file
    Elements []Element // Parsed HTML tree
}
```

### GoNode
```go
type GoNode struct {
    Source   string         // Go source code
    Imports  []string       // Import statements
    MainType string         // Main struct name
    Handlers []Handler      // HTTP handlers
    Props    *PropsStruct   // Component props
}
```

### ComponentDependency
```go
type ComponentDependency struct {
    Name       string            // "user-card"
    Path       string            // "components/user-card"
    Props      map[string]string // Props passed
    UsageCount int               // Times used
}
```

## 🛠 Funciones de Conveniencia

```go
// Parsear archivo directamente
ast, err := parser.ParseFile("file.gox")

// Validar archivo
err := parser.ValidateGoxFile("file.gox")

// Obtener solo dependencias
deps, err := parser.GetComponentDependencies("file.gox")

// Parsear directorio completo
files, err := parser.ParseDirectory("./pages")

// Verificar si es archivo .gox
if parser.IsGoxFile("file.gox") {
    // ...
}
```

## 🔧 Detección Avanzada

### Detector de Componentes
```go
detector := parser.NewComponentDetector()

// Detectar en HTML
components, err := detector.DetectComponents(htmlContent)

// Variables de template
vars := detector.DetectGoTemplateVariables(templateContent)

// Atributos HTMX
htmx := detector.DetectHTMXAttributes(templateContent)
```

## 🚨 Manejo de Errores

El parser proporciona errores detallados con información de línea y columna:

```go
ast, err := parser.ParseFile("file.gox")
if err != nil {
    if parseErr, ok := err.(*parser.ParseError); ok {
        fmt.Printf("Error en %s:%d:%d: %s\n", 
            parseErr.File, parseErr.Line, parseErr.Column, parseErr.Message)
    }
}
```

## 📝 Ejemplo Completo

Ver `main.go` y `api-test.go` en este directorio para ejemplos completos de uso.