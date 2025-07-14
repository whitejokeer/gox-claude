package parser

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func BenchmarkParseSimpleGoxFile(b *testing.B) {
	content := `<template>
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
</style>`

	parser := NewParser()
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		_, err := parser.ParseString(content)
		if err != nil {
			b.Fatalf("Parse error: %v", err)
		}
	}
}

func BenchmarkParseComplexGoxFile(b *testing.B) {
	content, err := os.ReadFile("testdata/with-components.gox")
	if err != nil {
		b.Skipf("Could not read test file: %v", err)
	}
	
	parser := NewParser()
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		_, err := parser.ParseBytes(content, "testdata/with-components.gox")
		if err != nil {
			b.Fatalf("Parse error: %v", err)
		}
	}
}

func BenchmarkParseWithManyComponents(b *testing.B) {
	// Generate a file with many components
	var sb strings.Builder
	sb.WriteString("<template>\n<div>\n")
	
	// Add 100 different components
	for i := 0; i < 100; i++ {
		sb.WriteString(fmt.Sprintf(`  <component-%d name="test" id="%d" />%s`, i, i, "\n"))
	}
	
	sb.WriteString("</div>\n</template>\n")
	sb.WriteString("<go>\npackage test\ntype TestPage struct{}\n</go>\n")
	
	content := sb.String()
	parser := NewParser()
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		_, err := parser.ParseString(content)
		if err != nil {
			b.Fatalf("Parse error: %v", err)
		}
	}
}

func BenchmarkComponentDetection(b *testing.B) {
	content := `<div>
		<user-card name="John" email="john@example.com" />
		<shared-button text="Save" hx-post="/save" />
		<user-card name="Jane" email="jane@example.com" />
		<shared-modal title="Confirm" />
		<data-table columns="name,email,status" />
	</div>`
	
	detector := NewComponentDetector()
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		_, err := detector.DetectComponents(content)
		if err != nil {
			b.Fatalf("Detection error: %v", err)
		}
	}
}

func BenchmarkLexerTokenization(b *testing.B) {
	content := []byte(`<template auth="required">
<div class="container">
  <h1>{{.Title}}</h1>
  <user-card name="{{.User.Name}}" />
</div>
</template>

<go>
package pages
type Page struct {
    Title string
    User  User
}
</go>

<style scoped>
.container { padding: 20px; }
</style>`)
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		lexer := NewLexer(content, "test.gox")
		_, err := lexer.Tokenize()
		if err != nil {
			b.Fatalf("Tokenization error: %v", err)
		}
	}
}

func BenchmarkGoCodeParsing(b *testing.B) {
	section := &Section{
		Type: "go",
		Content: `package components

import (
	"net/http"
	"errors"
)

type UserCard struct {
	ID     int    ` + "`gox:\"required\"`" + `
	Name   string ` + "`gox:\"required\"`" + `
	Email  string ` + "`gox:\"required\"`" + `
	Avatar string ` + "`gox:\"optional\"`" + `
}

func (c *UserCard) HandleClick(w http.ResponseWriter, r *http.Request) {
	// Handle click
}

func (c *UserCard) HandleEdit(w http.ResponseWriter, r *http.Request) {
	// Handle edit
}

func (c *UserCard) Validate() error {
	if c.Name == "" {
		return errors.New("name required")
	}
	return nil
}`,
		StartLine: 1,
		StartCol:  1,
	}
	
	parser := NewParser()
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		_, err := parser.parseGo(section)
		if err != nil {
			b.Fatalf("Go parsing error: %v", err)
		}
	}
}

func BenchmarkStyleParsing(b *testing.B) {
	section := &Section{
		Type: "style",
		Content: `.container {
	max-width: 1200px;
	margin: 0 auto;
	padding: 20px;
}

.header {
	background: #f8f9fa;
	border-bottom: 1px solid #dee2e6;
	padding: 1rem 0;
}

.button {
	display: inline-block;
	padding: 0.375rem 0.75rem;
	margin-bottom: 0;
	font-size: 1rem;
	font-weight: 400;
	line-height: 1.5;
	color: #212529;
	text-align: center;
	text-decoration: none;
	vertical-align: middle;
	cursor: pointer;
	user-select: none;
	background-color: transparent;
	border: 1px solid transparent;
	border-radius: 0.25rem;
}

.button:hover {
	color: #212529;
	text-decoration: none;
}`,
		Attributes: map[string]string{"scoped": ""},
		StartLine:  1,
		StartCol:   1,
	}
	
	parser := NewParser()
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		_, err := parser.parseStyle(section)
		if err != nil {
			b.Fatalf("Style parsing error: %v", err)
		}
	}
}

func BenchmarkHTMXAttributeDetection(b *testing.B) {
	content := `<div>
		<button hx-post="/submit" hx-target="#result" hx-swap="innerHTML">Submit</button>
		<div hx-get="/data" hx-trigger="click" hx-indicator="#spinner">Load Data</div>
		<form hx-boost="true" hx-confirm="Are you sure?">
			<input hx-post="/validate" hx-trigger="blur" hx-target="#error" />
		</form>
		<table hx-get="/table" hx-trigger="load" hx-swap="outerHTML">
			<tr hx-delete="/items/1" hx-confirm="Delete?">
				<td>Item 1</td>
			</tr>
		</table>
	</div>`
	
	detector := NewComponentDetector()
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		_ = detector.DetectHTMXAttributes(content)
	}
}

func BenchmarkGoTemplateVariableDetection(b *testing.B) {
	content := `<div>
		<h1>{{.Title}}</h1>
		<p>Welcome {{.User.Name}}</p>
		<span>{{.User.Email}}</span>
		<div>Score: {{.User.Score}}</div>
		<ul>
		{{range .Items}}
			<li>{{.Name}} - {{.Price}}</li>
		{{end}}
		</ul>
		<footer>{{.Copyright}} {{.Year}}</footer>
	</div>`
	
	detector := NewComponentDetector()
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		_ = detector.DetectGoTemplateVariables(content)
	}
}

func BenchmarkFileReconstruction(b *testing.B) {
	goxFile := &GoxFile{
		Template: &TemplateNode{
			Content: `<div class="container">
  <h1>{{.Title}}</h1>
  <user-card name="{{.User.Name}}" email="{{.User.Email}}" />
</div>`,
			Auth:   "required",
			Layout: "dashboard",
		},
		Go: &GoNode{
			Source: `package pages

import "net/http"

type UsersPage struct {
	Title string
	User  User
}

func (p *UsersPage) BeforeMount(w http.ResponseWriter, r *http.Request) error {
	p.Title = "Users"
	return nil
}`,
			MainType: "UsersPage",
		},
		Styles: []*StyleNode{
			{
				Content: `.container {
	padding: 20px;
	max-width: 800px;
	margin: 0 auto;
}`,
				Scoped: true,
			},
		},
		Components: []ComponentDependency{
			{
				Name: "user-card",
				Path: "components/user-card",
				Props: map[string]string{
					"name":  "{{.User.Name}}",
					"email": "{{.User.Email}}",
				},
			},
		},
	}
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		_ = goxFile.String()
	}
}

func BenchmarkParseFromFile(b *testing.B) {
	parser := NewParser()
	
	b.Run("simple", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := parser.ParseFile("testdata/simple.gox")
			if err != nil {
				b.Fatalf("Parse error: %v", err)
			}
		}
	})
	
	b.Run("component", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := parser.ParseFile("testdata/component.gox")
			if err != nil {
				b.Fatalf("Parse error: %v", err)
			}
		}
	})
	
	b.Run("with-components", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := parser.ParseFile("testdata/with-components.gox")
			if err != nil {
				b.Fatalf("Parse error: %v", err)
			}
		}
	})
}