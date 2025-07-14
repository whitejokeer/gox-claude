package parser

import (
	"strings"
	"testing"
)

func TestParseBasicGoxFile(t *testing.T) {
	input := `<template>
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
	ast, err := parser.ParseString(input)
	
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	
	if ast == nil {
		t.Fatal("Expected AST, got nil")
	}
	
	if ast.Template == nil {
		t.Error("Expected template section")
	}
	
	if ast.Go == nil {
		t.Error("Expected go section")
	}
	
	if len(ast.Styles) != 1 {
		t.Errorf("Expected 1 style section, got %d", len(ast.Styles))
	}
	
	// Check template content
	expectedTemplate := "<div>Hello {{.Name}}</div>"
	if !strings.Contains(ast.Template.Content, expectedTemplate) {
		t.Errorf("Template content mismatch. Expected to contain: %s", expectedTemplate)
	}
	
	// Check Go content
	expectedGo := "type HelloPage struct"
	if !strings.Contains(ast.Go.Source, expectedGo) {
		t.Errorf("Go content mismatch. Expected to contain: %s", expectedGo)
	}
	
	// Check style content
	expectedStyle := ".hello { color: blue; }"
	if !strings.Contains(ast.Styles[0].Content, expectedStyle) {
		t.Errorf("Style content mismatch. Expected to contain: %s", expectedStyle)
	}
}

func TestParseWithComponents(t *testing.T) {
	input := `<template>
  <div>
    <user-card name="{{.User.Name}}" email="{{.User.Email}}" />
    <shared-button text="Save" hx-post="/save" />
  </div>
</template>`

	parser := NewParser()
	ast, err := parser.ParseString(input)
	
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	
	if len(ast.Components) != 2 {
		t.Fatalf("Expected 2 components, got %d", len(ast.Components))
	}
	
	// Check user-card component
	userCard := findComponent(ast.Components, "user-card")
	if userCard == nil {
		t.Error("Expected user-card component")
	} else {
		if userCard.Path != "components/user-card" {
			t.Errorf("Expected path 'components/user-card', got '%s'", userCard.Path)
		}
		
		if userCard.Props["name"] != "{{.User.Name}}" {
			t.Errorf("Expected name prop '{{.User.Name}}', got '%s'", userCard.Props["name"])
		}
	}
	
	// Check shared-button component
	sharedButton := findComponent(ast.Components, "shared-button")
	if sharedButton == nil {
		t.Error("Expected shared-button component")
	} else {
		if sharedButton.Path != "shared/ui/button" {
			t.Errorf("Expected path 'shared/ui/button', got '%s'", sharedButton.Path)
		}
		
		if sharedButton.Props["text"] != "Save" {
			t.Errorf("Expected text prop 'Save', got '%s'", sharedButton.Props["text"])
		}
		
		if sharedButton.Props["hx-post"] != "/save" {
			t.Errorf("Expected hx-post prop '/save', got '%s'", sharedButton.Props["hx-post"])
		}
	}
}

func TestParseHTTPHandlers(t *testing.T) {
	input := `<go>
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
</go>`

	parser := NewParser()
	ast, err := parser.ParseString(input)
	
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	
	if len(ast.Go.Handlers) != 2 {
		t.Fatalf("Expected 2 handlers, got %d", len(ast.Go.Handlers))
	}
	
	// Check HandleSubmit
	handleSubmit := findHandler(ast.Go.Handlers, "HandleSubmit")
	if handleSubmit == nil {
		t.Error("Expected HandleSubmit handler")
	} else {
		if !handleSubmit.IsHTMX {
			t.Error("Expected HandleSubmit to be detected as HTMX handler")
		}
	}
	
	// Check ValidateEmail
	validateEmail := findHandler(ast.Go.Handlers, "ValidateEmail")
	if validateEmail == nil {
		t.Error("Expected ValidateEmail handler")
	}
}

func TestParseComponentProps(t *testing.T) {
	input := `<go>
package components

type UserCard struct {
    Name   string ` + "`gox:\"required\"`" + `
    Email  string ` + "`gox:\"required\"`" + `
    Avatar string ` + "`gox:\"optional\"`" + `
}
</go>`

	parser := NewParser()
	ast, err := parser.ParseString(input)
	
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	
	if ast.Go.Props == nil {
		t.Fatal("Expected props struct")
	}
	
	if ast.Go.Props.Name != "UserCard" {
		t.Errorf("Expected props name 'UserCard', got '%s'", ast.Go.Props.Name)
	}
	
	if len(ast.Go.Props.Fields) != 3 {
		t.Fatalf("Expected 3 fields, got %d", len(ast.Go.Props.Fields))
	}
	
	// Check required fields
	nameField := findField(ast.Go.Props.Fields, "Name")
	if nameField == nil {
		t.Error("Expected Name field")
	} else if !nameField.Required {
		t.Error("Expected Name field to be required")
	}
	
	emailField := findField(ast.Go.Props.Fields, "Email")
	if emailField == nil {
		t.Error("Expected Email field")
	} else if !emailField.Required {
		t.Error("Expected Email field to be required")
	}
	
	avatarField := findField(ast.Go.Props.Fields, "Avatar")
	if avatarField == nil {
		t.Error("Expected Avatar field")
	} else if avatarField.Required {
		t.Error("Expected Avatar field to be optional")
	}
}

func TestParseWithAttributes(t *testing.T) {
	input := `<template auth="required" layout="dashboard">
  <div>Content</div>
</template>

<style scoped>
.content { color: red; }
</style>`

	parser := NewParser()
	ast, err := parser.ParseString(input)
	
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	
	if ast.Template.Auth != "required" {
		t.Errorf("Expected auth 'required', got '%s'", ast.Template.Auth)
	}
	
	if ast.Template.Layout != "dashboard" {
		t.Errorf("Expected layout 'dashboard', got '%s'", ast.Template.Layout)
	}
	
	if !ast.Styles[0].Scoped {
		t.Error("Expected style to be scoped")
	}
}

func TestParseErrors(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "unclosed template",
			input: `<template><div>Hello</div>`,
		},
		{
			name:  "mismatched tags",
			input: `<template><div>Hello</div></go>`,
		},
		{
			name:  "invalid go syntax",
			input: `<go>package invalid syntax here</go>`,
		},
		{
			name:  "unknown section",
			input: `<unknown>content</unknown>`,
		},
	}
	
	parser := NewParser()
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parser.ParseString(tt.input)
			if err == nil {
				t.Error("Expected error but got none")
			}
		})
	}
}

func TestGoxFileValidation(t *testing.T) {
	tests := []struct {
		name    string
		goxFile *GoxFile
		wantErr bool
	}{
		{
			name: "valid file with template",
			goxFile: &GoxFile{
				Template: &TemplateNode{Content: "<div>Hello</div>"},
			},
			wantErr: false,
		},
		{
			name: "valid file with go",
			goxFile: &GoxFile{
				Go: &GoNode{Source: "package main"},
			},
			wantErr: false,
		},
		{
			name:    "empty file",
			goxFile: &GoxFile{},
			wantErr: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.goxFile.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGoxFileStringGeneration(t *testing.T) {
	goxFile := &GoxFile{
		Template: &TemplateNode{
			Content: "<div>Hello {{.Name}}</div>",
			Auth:    "required",
		},
		Go: &GoNode{
			Source: "package pages\n\ntype Page struct {\n    Name string\n}",
		},
		Styles: []*StyleNode{
			{
				Content: ".hello { color: blue; }",
				Scoped:  true,
			},
		},
	}
	
	output := goxFile.String()
	
	if !strings.Contains(output, "<template auth=\"required\">") {
		t.Error("Expected template with auth attribute")
	}
	
	if !strings.Contains(output, "<div>Hello {{.Name}}</div>") {
		t.Error("Expected template content")
	}
	
	if !strings.Contains(output, "<go>") {
		t.Error("Expected go section")
	}
	
	if !strings.Contains(output, "<style scoped>") {
		t.Error("Expected scoped style section")
	}
}

// Helper functions

func findComponent(components []ComponentDependency, name string) *ComponentDependency {
	for i := range components {
		if components[i].Name == name {
			return &components[i]
		}
	}
	return nil
}

func findHandler(handlers []Handler, name string) *Handler {
	for i := range handlers {
		if handlers[i].Name == name {
			return &handlers[i]
		}
	}
	return nil
}

func findField(fields []Field, name string) *Field {
	for i := range fields {
		if fields[i].Name == name {
			return &fields[i]
		}
	}
	return nil
}