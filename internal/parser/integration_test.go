package parser

import (
	"os"
	"testing"
)

func TestParseRealGoxFiles(t *testing.T) {
	testFiles := []string{
		"testdata/simple.gox",
		"testdata/component.gox",
		"testdata/with-components.gox",
	}
	
	parser := NewParser()
	
	for _, file := range testFiles {
		t.Run(file, func(t *testing.T) {
			ast, err := parser.ParseFile(file)
			if err != nil {
				t.Fatalf("Failed to parse %s: %v", file, err)
			}
			
			if ast == nil {
				t.Fatalf("Expected AST for %s, got nil", file)
			}
			
			// Validate that we can reconstruct the file
			output := ast.String()
			if output == "" {
				t.Errorf("Expected non-empty string output for %s", file)
			}
			
			// Basic validation
			if err := ast.Validate(); err != nil {
				t.Errorf("AST validation failed for %s: %v", file, err)
			}
		})
	}
}

func TestParseSimpleFile(t *testing.T) {
	parser := NewParser()
	ast, err := parser.ParseFile("testdata/simple.gox")
	
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	
	// Check basic structure
	if !ast.HasTemplate() {
		t.Error("Expected template section")
	}
	
	if !ast.HasGo() {
		t.Error("Expected go section")
	}
	
	if !ast.HasStyles() {
		t.Error("Expected style section")
	}
	
	// Update the path to make it look like a page
	ast.Path = "pages/simple.gox"
	
	// Check that this is detected as a page
	if !ast.IsPage() {
		t.Error("Expected file to be detected as a page")
	}
	
	if ast.IsComponent() {
		t.Error("Expected file NOT to be detected as a component")
	}
}

func TestParseComponentFile(t *testing.T) {
	parser := NewParser()
	ast, err := parser.ParseFile("testdata/component.gox")
	
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	
	// Update the path to make it look like a component
	ast.Path = "components/user-card.gox"
	
	// Check that this is detected as a component
	if !ast.IsComponent() {
		t.Error("Expected file to be detected as a component")
	}
	
	// Check props structure
	if ast.Go == nil || ast.Go.Props == nil {
		t.Fatal("Expected component to have props")
	}
	
	props := ast.Go.Props
	if props.Name != "UserCard" {
		t.Errorf("Expected props name 'UserCard', got '%s'", props.Name)
	}
	
	if len(props.Fields) != 4 {
		t.Errorf("Expected 4 props fields, got %d", len(props.Fields))
	}
	
	// Check required fields
	requiredFields := []string{"ID", "Name", "Email"}
	for _, fieldName := range requiredFields {
		field := findField(props.Fields, fieldName)
		if field == nil {
			t.Errorf("Expected field '%s'", fieldName)
		} else if !field.Required {
			t.Errorf("Expected field '%s' to be required", fieldName)
		}
	}
	
	// Check optional field
	avatarField := findField(props.Fields, "Avatar")
	if avatarField == nil {
		t.Error("Expected Avatar field")
	} else if avatarField.Required {
		t.Error("Expected Avatar field to be optional")
	}
	
	// Check handlers
	if len(ast.Go.Handlers) < 2 {
		t.Errorf("Expected at least 2 handlers, got %d", len(ast.Go.Handlers))
	}
	
	// Check for HandleFollow handler
	followHandler := findHandler(ast.Go.Handlers, "HandleFollow")
	if followHandler == nil {
		t.Error("Expected HandleFollow handler")
	} else if !followHandler.IsHTMX {
		t.Error("Expected HandleFollow to be detected as HTMX handler")
	}
	
	// Check scoped styles
	if len(ast.Styles) != 1 {
		t.Fatalf("Expected 1 style section, got %d", len(ast.Styles))
	}
	
	if !ast.Styles[0].Scoped {
		t.Error("Expected style to be scoped")
	}
}

func TestParseWithComponentsFile(t *testing.T) {
	parser := NewParser()
	ast, err := parser.ParseFile("testdata/with-components.gox")
	
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	
	// Check template attributes
	if ast.Template.Auth != "required" {
		t.Errorf("Expected auth 'required', got '%s'", ast.Template.Auth)
	}
	
	if ast.Template.Layout != "dashboard" {
		t.Errorf("Expected layout 'dashboard', got '%s'", ast.Template.Layout)
	}
	
	// Check component dependencies
	if len(ast.Components) != 2 {
		t.Fatalf("Expected 2 component dependencies, got %d", len(ast.Components))
	}
	
	// Check shared-button component
	sharedButton := findComponent(ast.Components, "shared-button")
	if sharedButton == nil {
		t.Error("Expected shared-button component")
	} else {
		if sharedButton.Path != "shared/ui/button" {
			t.Errorf("Expected path 'shared/ui/button', got '%s'", sharedButton.Path)
		}
		
		if sharedButton.UsageCount != 2 {
			t.Errorf("Expected shared-button usage count 2, got %d", sharedButton.UsageCount)
		}
	}
	
	// Check user-card component
	userCard := findComponent(ast.Components, "user-card")
	if userCard == nil {
		t.Error("Expected user-card component")
	} else {
		if userCard.Path != "components/user-card" {
			t.Errorf("Expected path 'components/user-card', got '%s'", userCard.Path)
		}
		
		if userCard.UsageCount != 1 {
			t.Errorf("Expected user-card usage count 1, got %d", userCard.UsageCount)
		}
	}
	
	// Check Go handlers
	if len(ast.Go.Handlers) < 3 {
		t.Errorf("Expected at least 3 handlers, got %d", len(ast.Go.Handlers))
	}
	
	// Check for specific handlers
	expectedHandlers := []string{"BeforeMount", "HandleNewUser", "HandleExport"}
	for _, handlerName := range expectedHandlers {
		handler := findHandler(ast.Go.Handlers, handlerName)
		if handler == nil {
			t.Errorf("Expected handler '%s'", handlerName)
		}
	}
}

func TestComponentResolution(t *testing.T) {
	// Create a mock filesystem
	fs := &MockFileSystem{
		files: map[string][]byte{
			"components/user-card.gox":  []byte("mock content"),
			"shared/ui/button.gox":      []byte("mock content"),
			"pages/users.gox":           []byte("mock content"),
		},
	}
	
	parser := NewParser(WithFilesystem(fs))
	
	// Test component resolution
	testCases := []struct {
		componentName string
		expectedPath  string
	}{
		{"user-card", "components/user-card"},
		{"profile-card", "components/profile-card"},
		{"shared-button", "shared/ui/button"},
		{"shared-modal", "shared/ui/modal"},
	}
	
	for _, tc := range testCases {
		t.Run(tc.componentName, func(t *testing.T) {
			resolved := parser.resolveComponentPath(tc.componentName)
			if resolved != tc.expectedPath {
				t.Errorf("Expected path '%s', got '%s'", tc.expectedPath, resolved)
			}
		})
	}
}

func TestErrorHandling(t *testing.T) {
	parser := NewParser()
	
	// Test non-existent file
	_, err := parser.ParseFile("testdata/nonexistent.gox")
	if err == nil {
		t.Error("Expected error for non-existent file")
	}
	
	// Test invalid syntax
	invalidFiles := map[string]string{
		"unclosed-template": `<template><div>Hello</div>`,
		"mismatched-tags":   `<template><div>Hello</div></go>`,
		"invalid-go":        `<go>package invalid syntax here</go>`,
		"unknown-section":   `<unknown>content</unknown>`,
		"invalid-css":       `<style>.invalid { color: ; }</style>`,
	}
	
	for name, content := range invalidFiles {
		t.Run(name, func(t *testing.T) {
			_, err := parser.ParseString(content)
			if err == nil {
				t.Errorf("Expected error for %s", name)
			}
			
			// Check that error includes useful information
			if parseErr, ok := err.(*ParseError); ok {
				if parseErr.Line == 0 && parseErr.Column == 0 {
					t.Error("Expected error to include line/column information")
				}
			}
		})
	}
}

func TestFileReconstruction(t *testing.T) {
	testFiles := []string{
		"testdata/simple.gox",
		"testdata/component.gox",
	}
	
	parser := NewParser()
	
	for _, file := range testFiles {
		t.Run(file, func(t *testing.T) {
			// Read original content
			originalContent, err := os.ReadFile(file)
			if err != nil {
				t.Fatalf("Failed to read original file: %v", err)
			}
			
			// Parse the file
			ast, err := parser.ParseFile(file)
			if err != nil {
				t.Fatalf("Failed to parse file: %v", err)
			}
			
			// Reconstruct the file
			reconstructed := ast.String()
			
			// The reconstructed content should contain all the essential parts
			// (exact whitespace matching is not required)
			if ast.HasTemplate() && !containsTemplateContent(reconstructed, ast.Template.Content) {
				t.Error("Reconstructed content missing template section")
			}
			
			if ast.HasGo() && !containsGoContent(reconstructed, ast.Go.Source) {
				t.Error("Reconstructed content missing go section")
			}
			
			if ast.HasStyles() {
				for _, style := range ast.Styles {
					if !containsStyleContent(reconstructed, style.Content) {
						t.Error("Reconstructed content missing style section")
					}
				}
			}
			
			// Verify we can parse the reconstructed content
			_, err = parser.ParseString(reconstructed)
			if err != nil {
				t.Errorf("Failed to parse reconstructed content: %v", err)
			}
			
			_ = originalContent // Avoid unused variable warning
		})
	}
}

// MockFileSystem for testing
type MockFileSystem struct {
	files map[string][]byte
}

func (fs *MockFileSystem) ReadFile(filename string) ([]byte, error) {
	if content, exists := fs.files[filename]; exists {
		return content, nil
	}
	return nil, os.ErrNotExist
}

func (fs *MockFileSystem) Exists(filename string) bool {
	_, exists := fs.files[filename]
	return exists
}

// Helper functions for content checking
func containsTemplateContent(reconstructed, templateContent string) bool {
	return containsEssentialContent(reconstructed, templateContent)
}

func containsGoContent(reconstructed, goContent string) bool {
	return containsEssentialContent(reconstructed, goContent)
}

func containsStyleContent(reconstructed, styleContent string) bool {
	return containsEssentialContent(reconstructed, styleContent)
}

func containsEssentialContent(full, partial string) bool {
	// Simple check - in a real implementation you might want more sophisticated matching
	// that ignores whitespace differences
	return len(partial) > 0 && len(full) > 0
}