package parser

import (
	"testing"
)

func TestDetectComponents(t *testing.T) {
	detector := NewComponentDetector()
	
	content := `<div>
		<user-card name="John" email="john@example.com" />
		<shared-button text="Save" hx-post="/save" />
		<user-card name="Jane" email="jane@example.com" />
	</div>`
	
	components, err := detector.DetectComponents(content)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	
	if len(components) != 2 {
		t.Fatalf("Expected 2 unique components, got %d", len(components))
	}
	
	// Check user-card component
	userCard := findComponentByName(components, "user-card")
	if userCard == nil {
		t.Error("Expected user-card component")
	} else {
		if userCard.UsageCount != 2 {
			t.Errorf("Expected user-card usage count 2, got %d", userCard.UsageCount)
		}
		
		if userCard.Props["name"] == "" {
			t.Error("Expected user-card to have name prop")
		}
		
		if userCard.Props["email"] == "" {
			t.Error("Expected user-card to have email prop")
		}
	}
	
	// Check shared-button component
	sharedButton := findComponentByName(components, "shared-button")
	if sharedButton == nil {
		t.Error("Expected shared-button component")
	} else {
		if sharedButton.UsageCount != 1 {
			t.Errorf("Expected shared-button usage count 1, got %d", sharedButton.UsageCount)
		}
		
		if sharedButton.Props["text"] != "Save" {
			t.Errorf("Expected text prop 'Save', got '%s'", sharedButton.Props["text"])
		}
		
		if sharedButton.Props["hx-post"] != "/save" {
			t.Errorf("Expected hx-post prop '/save', got '%s'", sharedButton.Props["hx-post"])
		}
	}
}

func TestDetectHTMXAttributes(t *testing.T) {
	detector := NewComponentDetector()
	
	content := `<div>
		<button hx-post="/submit" hx-target="#result">Submit</button>
		<div hx-get="/data" hx-trigger="click">Load Data</div>
		<form hx-boost="true">Form</form>
	</div>`
	
	attributes := detector.DetectHTMXAttributes(content)
	
	expectedAttrs := []string{"hx-post", "hx-target", "hx-get", "hx-trigger", "hx-boost"}
	
	if len(attributes) != len(expectedAttrs) {
		t.Fatalf("Expected %d HTMX attributes, got %d", len(expectedAttrs), len(attributes))
	}
	
	for _, expected := range expectedAttrs {
		found := false
		for _, attr := range attributes {
			if attr == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected HTMX attribute '%s' not found", expected)
		}
	}
}

func TestDetectGoTemplateVariables(t *testing.T) {
	detector := NewComponentDetector()
	
	content := `<div>
		<h1>{{.Title}}</h1>
		<p>Welcome {{.User.Name}}</p>
		<span>{{.Count}} items</span>
		<div>{{.User.Name}}</div> <!-- duplicate should not be counted twice -->
	</div>`
	
	variables := detector.DetectGoTemplateVariables(content)
	
	expectedVars := []string{"Title", "User", "Count"}
	
	if len(variables) != len(expectedVars) {
		t.Fatalf("Expected %d variables, got %d: %v", len(expectedVars), len(variables), variables)
	}
	
	for _, expected := range expectedVars {
		found := false
		for _, variable := range variables {
			if variable == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected variable '%s' not found", expected)
		}
	}
}

func TestIsValidComponentName(t *testing.T) {
	detector := NewComponentDetector()
	
	tests := []struct {
		name     string
		expected bool
	}{
		{"user-card", true},
		{"shared-button", true},
		{"my-awesome-component", true},
		{"user-profile-card", true},
		{"component-1", true},
		{"user", false},        // no hyphen
		{"User-Card", false},   // uppercase
		{"user_card", false},   // underscore
		{"-user-card", false},  // starts with hyphen
		{"user-card-", false},  // ends with hyphen
		{"user--card", false},  // double hyphen
		{"1user-card", false},  // starts with number
		{"", false},            // empty
		{"a", false},           // too short
		{"ab", false},          // too short, no hyphen
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := detector.IsValidComponentName(tt.name)
			if result != tt.expected {
				t.Errorf("IsValidComponentName(%s) = %v, expected %v", tt.name, result, tt.expected)
			}
		})
	}
}

func TestParseComponentPropsDetector(t *testing.T) {
	detector := NewComponentDetector()
	
	tests := []struct {
		name       string
		attributes string
		expected   map[string]string
	}{
		{
			name:       "simple props",
			attributes: ` name="John" email="john@example.com"`,
			expected: map[string]string{
				"name":  "John",
				"email": "john@example.com",
			},
		},
		{
			name:       "htmx attributes",
			attributes: ` text="Save" hx-post="/save" hx-target="#result"`,
			expected: map[string]string{
				"text":      "Save",
				"hx-post":   "/save",
				"hx-target": "#result",
			},
		},
		{
			name:       "single quotes",
			attributes: ` name='John' email='john@example.com'`,
			expected: map[string]string{
				"name":  "John",
				"email": "john@example.com",
			},
		},
		{
			name:       "boolean attributes",
			attributes: ` required disabled checked`,
			expected: map[string]string{
				"required": "",
				"disabled": "",
				"checked":  "",
			},
		},
		{
			name:       "mixed attributes",
			attributes: ` name="John" required class="btn btn-primary"`,
			expected: map[string]string{
				"name":     "John",
				"required": "",
				"class":    "btn btn-primary",
			},
		},
		{
			name:       "empty",
			attributes: "",
			expected:   map[string]string{},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := detector.parseComponentProps(tt.attributes)
			
			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d props, got %d. Got: %+v", len(tt.expected), len(result), result)
			}
			
			for key, expectedValue := range tt.expected {
				if actualValue, exists := result[key]; !exists {
					t.Errorf("Expected prop '%s' not found", key)
				} else if actualValue != expectedValue {
					t.Errorf("Prop '%s': expected '%s', got '%s'", key, expectedValue, actualValue)
				}
			}
		})
	}
}

func TestExtractComponentsFromHTML(t *testing.T) {
	detector := NewComponentDetector()
	
	content := `<div class="container">
	<user-card name="John" email="john@example.com" />
	<div>
		<shared-button text="Save" hx-post="/save" />
	</div>
	<user-profile name="Jane" avatar="/avatars/jane.jpg" />
</div>`
	
	components, err := detector.ExtractComponentsFromHTML(content)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	
	if len(components) != 3 {
		t.Fatalf("Expected 3 components, got %d", len(components))
	}
	
	// Check line numbers are set
	for i, comp := range components {
		if comp.Line == 0 {
			t.Errorf("Component %d should have line number set", i)
		}
		
		if comp.Column == 0 {
			t.Errorf("Component %d should have column number set", i)
		}
	}
	
	// Check first component is on line 2
	firstComp := components[0]
	if firstComp.Line != 2 {
		t.Errorf("Expected first component on line 2, got line %d", firstComp.Line)
	}
	
	if firstComp.Name != "user-card" {
		t.Errorf("Expected first component name 'user-card', got '%s'", firstComp.Name)
	}
}

func TestDetectGoTemplateFunctions(t *testing.T) {
	detector := NewComponentDetector()
	
	content := `<div>
		{{formatDate .CreatedAt}}
		{{if .IsLoggedIn}}
			{{capitalize .User.Name}}
		{{end}}
		{{range .Items}}
			{{truncate .Description 100}}
		{{end}}
	</div>`
	
	functions := detector.DetectGoTemplateFunctions(content)
	
	expectedFuncs := []string{"formatDate", "capitalize", "truncate"}
	
	if len(functions) != len(expectedFuncs) {
		t.Fatalf("Expected %d functions, got %d: %v", len(expectedFuncs), len(functions), functions)
	}
	
	for _, expected := range expectedFuncs {
		found := false
		for _, function := range functions {
			if function == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected function '%s' not found", expected)
		}
	}
}

// Helper function
func findComponentByName(components []ComponentDependency, name string) *ComponentDependency {
	for i := range components {
		if components[i].Name == name {
			return &components[i]
		}
	}
	return nil
}