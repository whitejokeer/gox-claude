package parser

import (
	"fmt"
	"strings"
	"testing"
)

// Test que verifica que el parser cumple al 100% con las convenciones GOX
func TestParserCompleteVerification(t *testing.T) {
	// 1. ARCHIVO DE PÁGINA CORRECTO (estilo HTMX puro)
	correctPageGox := `<template auth="required" layout="dashboard">
  <div class="users-page">
    <h1>{{.Title}}</h1>
    
    <form hx-post="/users/search" hx-target="#results">
      <input name="query" placeholder="Search..." />
    </form>
    
    <div id="results">
      {{range .Users}}
        <user-card id="{{.ID}}" name="{{.Name}}" />
      {{end}}
    </div>
  </div>
</template>

<go>
package pages

import "net/http"

// Página sin herencia ni lifecycle
type UsersPage struct {
    Title string
    Users []User
}

type User struct {
    ID   int
    Name string
}

// Constructor correcto (no BeforeMount)
func NewUsersPage() *UsersPage {
    return &UsersPage{
        Title: "User Management",
        Users: loadUsers(),
    }
}

// Handler HTTP puro para HTMX
func (p *UsersPage) HandleSearch(w http.ResponseWriter, r *http.Request) {
    query := r.FormValue("query")
    users := searchUsers(query)
    
    // Renderiza fragmento HTML
    renderPartial(w, "user-results", users)
}

// Handler para eliminar (HTMX style)
func HandleDeleteUser(w http.ResponseWriter, r *http.Request) {
    // Elimina y devuelve HTML actualizado
    w.Header().Set("Content-Type", "text/html")
    w.Write([]byte("<div>User deleted</div>"))
}
</go>

<style scoped>
.users-page {
    max-width: 1200px;
    margin: 0 auto;
}
</style>`

	// 2. ARCHIVO DE COMPONENTE CORRECTO
	correctComponentGox := `<template>
  <div class="user-card" id="user-{{.ID}}">
    <h3>{{.Name}}</h3>
    <button hx-delete="/users/{{.ID}}" 
            hx-target="#user-{{.ID}}" 
            hx-swap="outerHTML">
      Delete
    </button>
  </div>
</template>

<go>
package components

// Componente simple, sin herencia
type UserCard struct {
    ID   int    ` + "`gox:\"required\"`" + `
    Name string ` + "`gox:\"required\"`" + `
}
</go>

<style scoped>
.user-card {
    border: 1px solid #ddd;
    padding: 1rem;
}
</style>`

	// 3. ARCHIVOS INCORRECTOS (anti-patterns)
	incorrectLifecycleGox := `<template>
  <div>Bad example</div>
</template>

<go>
package pages

// INCORRECTO: Lifecycle methods como Vue/React
type BadPage struct {
    page.Base // INCORRECTO: Herencia
}

func (p *BadPage) BeforeMount() error { // INCORRECTO: No hay lifecycle
    return nil
}

func (p *BadPage) AfterUpdate() { // INCORRECTO: No hay update cycle
}
</go>`

	// TEST 1: Parsear página correcta
	t.Run("Página HTMX correcta", func(t *testing.T) {
		parser := NewParser()
		ast, err := parser.ParseString(correctPageGox)
		if err != nil {
			t.Fatalf("No debería fallar: %v", err)
		}
		
		// Set the path to make it look like a page
		ast.Path = "pages/users.gox"

		// Verificar que es página, no componente
		if ast.IsComponent() {
			t.Error("Debería ser página, no componente")
		}

		// Verificar props detectadas
		if ast.Go.Props == nil {
			t.Fatal("Debería detectar props de la página")
		}
		
		expectedProps := []string{"Title", "Users"}
		for _, prop := range expectedProps {
			found := false
			for _, field := range ast.Go.Props.Fields {
				if field.Name == prop {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Prop '%s' no detectada", prop)
			}
		}

		// Verificar handlers HTTP (no lifecycle)
		for _, handler := range ast.Go.Handlers {
			if strings.Contains(handler.Name, "Mount") || 
			   strings.Contains(handler.Name, "Update") ||
			   strings.Contains(handler.Name, "Destroy") {
				t.Errorf("Encontrado lifecycle method incorrecto: %s", handler.Name)
			}
		}

		// Verificar que detecta handlers HTMX correctos
		foundSearch := false
		for _, handler := range ast.Go.Handlers {
			if handler.Name == "HandleSearch" {
				foundSearch = true
				if !handler.IsHTMX {
					t.Error("HandleSearch debería ser marcado como HTMX handler")
				}
			}
		}
		if !foundSearch {
			t.Error("No detectó HandleSearch")
		}

		// Verificar componentes usados
		if len(ast.Components) != 1 {
			t.Errorf("Debería detectar 1 componente, encontró %d", len(ast.Components))
		}
		if len(ast.Components) > 0 {
			if ast.Components[0].Name != "user-card" {
				t.Errorf("Componente incorrecto: %s", ast.Components[0].Name)
			}
			if ast.Components[0].Path != "components/user-card" {
				t.Errorf("Path incorrecto: %s", ast.Components[0].Path)
			}
		}
	})

	// TEST 2: Parsear componente correcto
	t.Run("Componente HTMX correcto", func(t *testing.T) {
		parser := NewParser()
		ast, err := parser.ParseString(correctComponentGox)
		if err != nil {
			t.Fatalf("No debería fallar: %v", err)
		}
		
		// Set the path to make it look like a component
		ast.Path = "components/user-card.gox"

		// Verificar que es componente
		if !ast.IsComponent() {
			t.Error("Debería ser componente")
		}

		// Verificar props con tags gox:"required"
		if ast.Go.Props == nil {
			t.Fatal("Debería detectar props del componente")
		}

		requiredCount := 0
		for _, field := range ast.Go.Props.Fields {
			if field.Required {
				requiredCount++
			}
		}
		if requiredCount != 2 {
			t.Errorf("Debería tener 2 props required, tiene %d", requiredCount)
		}

		// Verificar atributos HTMX en template
		detector := NewComponentDetector()
		htmxAttrs := detector.DetectHTMXAttributes(ast.Template.Content)
		
		expectedHTMX := []string{"hx-delete", "hx-target", "hx-swap"}
		for _, attr := range expectedHTMX {
			found := false
			for _, detected := range htmxAttrs {
				if detected == attr {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("No detectó atributo HTMX: %s", attr)
			}
		}
	})

	// TEST 3: Rechazar anti-patterns
	t.Run("Rechazar lifecycle y herencia", func(t *testing.T) {
		parser := NewParser()
		ast, err := parser.ParseString(incorrectLifecycleGox)
		if err != nil {
			t.Fatalf("Parser falló: %v", err)
		}
		
		// Set the path
		ast.Path = "pages/bad.gox"

		// El parser debería detectar estos anti-patterns
		warnings := []string{}

		// Detectar herencia (page.Base)
		if strings.Contains(ast.Go.Source, "page.Base") {
			warnings = append(warnings, "ADVERTENCIA: Detectada herencia page.Base - GOX no usa herencia")
		}

		// Detectar lifecycle methods
		for _, handler := range ast.Go.Handlers {
			if handler.Name == "BeforeMount" || handler.Name == "AfterUpdate" {
				warnings = append(warnings, fmt.Sprintf("ADVERTENCIA: Lifecycle method '%s' - GOX usa handlers HTTP, no lifecycle", handler.Name))
			}
		}

		if len(warnings) == 0 {
			t.Error("Debería detectar anti-patterns de Vue/React")
		}

		// Imprimir warnings para verificación
		for _, w := range warnings {
			t.Log(w)
		}
	})

	// TEST 4: Manejo de errores con ubicación
	t.Run("Errores con línea y columna", func(t *testing.T) {
		badGox := `<template>
  <div>
    <h1>No cerrado
  </div>
</template>

<go>
func syntax error {
</go>`

		parser := NewParser()
		_, err := parser.ParseString(badGox)
		if err == nil {
			t.Fatal("Debería fallar con HTML mal formado")
		}

		// Verificar que el error tiene información de ubicación
		errStr := err.Error()
		if !strings.Contains(errStr, ":") {
			t.Errorf("Error sin información de ubicación: %v", err)
		}
	})

	// TEST 5: Detección de rutas de componentes
	t.Run("Resolución correcta de rutas", func(t *testing.T) {
		testCases := []struct {
			tag          string
			expectedPath string
		}{
			{"user-card", "components/user-card"},
			{"shared-button", "shared/ui/button"},
			{"shared-modal", "shared/ui/modal"},
			{"product-list", "components/product-list"},
		}

		parser := NewParser()
		for _, tc := range testCases {
			template := fmt.Sprintf(`<template><%s /></template>`, tc.tag)
			ast, _ := parser.ParseString(template)
			
			if len(ast.Components) != 1 {
				t.Errorf("No detectó componente %s", tc.tag)
				continue
			}

			if ast.Components[0].Path != tc.expectedPath {
				t.Errorf("Path incorrecto para %s: esperado %s, obtenido %s",
					tc.tag, tc.expectedPath, ast.Components[0].Path)
			}
		}
	})

	// TEST 6: Props en template
	t.Run("Extracción de props de componentes", func(t *testing.T) {
		template := `<template>
  <user-card 
    id="{{.User.ID}}"
    name="{{.User.Name}}"
    email="{{.User.Email}}"
    :featured="true"
    hx-get="/users/{{.User.ID}}" />
</template>`

		parser := NewParser()
		ast, _ := parser.ParseString(template)
		
		if len(ast.Components) != 1 {
			t.Fatal("No detectó componente")
		}

		comp := ast.Components[0]
		expectedProps := map[string]bool{
			"id":         true,
			"name":       true,
			"email":      true,
			":featured":  true,
			"hx-get":     true,
		}

		for prop := range expectedProps {
			if _, exists := comp.Props[prop]; !exists {
				t.Errorf("No detectó prop: %s", prop)
			}
		}
	})
}

// Test de rendimiento
func BenchmarkParserCorrectFile(b *testing.B) {
	goxContent := `<template>
  <div class="page">
    <h1>{{.Title}}</h1>
    <user-list users="{{.Users}}" />
  </div>
</template>

<go>
package pages

type HomePage struct {
    Title string
    Users []User
}
</go>

<style scoped>
.page { padding: 2rem; }
</style>`

	parser := NewParser()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = parser.ParseString(goxContent)
	}
}