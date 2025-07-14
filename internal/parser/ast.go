package parser

import (
	"fmt"
	"strings"
)

// GoxFile represents a parsed .gox file
type GoxFile struct {
	Path       string                 `json:"path"`
	Template   *TemplateNode          `json:"template,omitempty"`
	Go         *GoNode                `json:"go,omitempty"`
	Styles     []*StyleNode           `json:"styles,omitempty"`
	Components []ComponentDependency  `json:"components,omitempty"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
}

// String returns a string representation of the GoxFile
func (gf *GoxFile) String() string {
	var sb strings.Builder
	
	if gf.Template != nil {
		sb.WriteString("<template")
		if gf.Template.Auth != "" {
			sb.WriteString(fmt.Sprintf(" auth=\"%s\"", gf.Template.Auth))
		}
		if gf.Template.Layout != "" {
			sb.WriteString(fmt.Sprintf(" layout=\"%s\"", gf.Template.Layout))
		}
		sb.WriteString(">\n")
		sb.WriteString(gf.Template.Content)
		sb.WriteString("</template>\n\n")
	}
	
	if gf.Go != nil {
		sb.WriteString("<go>\n")
		sb.WriteString(gf.Go.Source)
		sb.WriteString("</go>\n\n")
	}
	
	for _, style := range gf.Styles {
		sb.WriteString("<style")
		if style.Scoped {
			sb.WriteString(" scoped")
		}
		sb.WriteString(">\n")
		sb.WriteString(style.Content)
		sb.WriteString("</style>\n\n")
	}
	
	return strings.TrimSpace(sb.String())
}

// TemplateNode represents the template section of a .gox file
type TemplateNode struct {
	Content  string    `json:"content"`
	Auth     string    `json:"auth,omitempty"`     // "required", "role:admin", etc.
	Layout   string    `json:"layout,omitempty"`   // layout file to use
	Elements []Element `json:"elements,omitempty"` // Parsed HTML elements
}

// Element represents an HTML element or component in the template
type Element struct {
	Type        string            `json:"type"`        // "div", "user-card", etc.
	IsComponent bool              `json:"is_component"`
	Props       map[string]string `json:"props,omitempty"`       // Includes hx-* attributes
	Children    []Element         `json:"children,omitempty"`
	Content     string            `json:"content,omitempty"`     // Text content
	Line        int               `json:"line"`
	Column      int               `json:"column"`
}

// GoNode represents the Go code section of a .gox file
type GoNode struct {
	Source   string         `json:"source"`
	Imports  []string       `json:"imports,omitempty"`
	MainType string         `json:"main_type,omitempty"` // "UserCard" struct name
	Handlers []Handler      `json:"handlers,omitempty"`
	Props    *PropsStruct   `json:"props,omitempty"` // If this is a component
}

// Handler represents an HTTP handler method in the Go code
type Handler struct {
	Name        string            `json:"name"`
	Method      string            `json:"method,omitempty"`      // HTTP method if auto-detected
	Path        string            `json:"path,omitempty"`        // Path if auto-detected
	IsHTMX      bool              `json:"is_htmx"`               // Whether this is an HTMX handler
	Parameters  []Parameter       `json:"parameters,omitempty"`
	ReturnType  string            `json:"return_type,omitempty"`
	Line        int               `json:"line"`
	Column      int               `json:"column"`
}

// Parameter represents a function parameter
type Parameter struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// PropsStruct represents the props structure for a component
type PropsStruct struct {
	Name   string  `json:"name"`
	Fields []Field `json:"fields"`
}

// Field represents a struct field
type Field struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Tags     string `json:"tags,omitempty"`
	Required bool   `json:"required"`
	Line     int    `json:"line"`
	Column   int    `json:"column"`
}

// StyleNode represents a style section of a .gox file
type StyleNode struct {
	Content string `json:"content"`
	Scoped  bool   `json:"scoped"`
	Type    string `json:"type,omitempty"` // "css", "tailwind", etc.
	Line    int    `json:"line"`
	Column  int    `json:"column"`
}

// ComponentDependency represents a component used by this .gox file
type ComponentDependency struct {
	Name      string            `json:"name"`       // "user-card"
	Path      string            `json:"path"`       // "components/user-card"
	Props     map[string]string `json:"props,omitempty"`      // Props passed to component
	Line      int               `json:"line"`
	Column    int               `json:"column"`
	UsageCount int              `json:"usage_count"` // How many times this component is used
}

// ParseError represents a parsing error with location information
type ParseError struct {
	File    string `json:"file"`
	Line    int    `json:"line"`
	Column  int    `json:"column"`
	Message string `json:"message"`
}

// Error implements the error interface
func (pe *ParseError) Error() string {
	if pe.File != "" {
		return fmt.Sprintf("%s:%d:%d: %s", pe.File, pe.Line, pe.Column, pe.Message)
	}
	return fmt.Sprintf("%d:%d: %s", pe.Line, pe.Column, pe.Message)
}

// HasTemplate returns true if the file has a template section
func (gf *GoxFile) HasTemplate() bool {
	return gf.Template != nil && strings.TrimSpace(gf.Template.Content) != ""
}

// HasGo returns true if the file has a Go section
func (gf *GoxFile) HasGo() bool {
	return gf.Go != nil && strings.TrimSpace(gf.Go.Source) != ""
}

// HasStyles returns true if the file has style sections
func (gf *GoxFile) HasStyles() bool {
	return len(gf.Styles) > 0
}

// IsComponent returns true if this file represents a component
func (gf *GoxFile) IsComponent() bool {
	// A component is determined by its path, not by having props
	// All structs in GOX have props (public fields)
	return strings.Contains(gf.Path, "/components/") || 
	       strings.Contains(gf.Path, "/shared/ui/") ||
	       strings.Contains(gf.Path, "components/") ||
	       strings.Contains(gf.Path, "shared/ui/")
}

// IsPage returns true if this file represents a page
func (gf *GoxFile) IsPage() bool {
	return strings.Contains(gf.Path, "/pages/") || strings.HasPrefix(gf.Path, "pages/")
}

// GetComponentName returns the component name from the file path
func (gf *GoxFile) GetComponentName() string {
	if !gf.IsComponent() {
		return ""
	}
	
	// Extract component name from path
	// components/user-card.gox -> user-card
	// shared/ui/button.gox -> shared-button
	
	base := strings.TrimSuffix(gf.Path, ".gox")
	
	if strings.Contains(base, "shared/ui/") {
		parts := strings.Split(base, "/")
		if len(parts) > 0 {
			return "shared-" + parts[len(parts)-1]
		}
	}
	
	if strings.Contains(base, "components/") {
		parts := strings.Split(base, "/")
		if len(parts) > 0 {
			return parts[len(parts)-1]
		}
	}
	
	return ""
}

// Validate performs basic validation on the AST
func (gf *GoxFile) Validate() error {
	// A .gox file must have at least one section
	if !gf.HasTemplate() && !gf.HasGo() && !gf.HasStyles() {
		return &ParseError{
			File:    gf.Path,
			Message: "file must contain at least one section (template, go, or style)",
		}
	}
	
	// Component files must have a Go section with props
	if gf.IsComponent() && gf.Go.Props == nil {
		return &ParseError{
			File:    gf.Path,
			Message: "component files must define a props struct",
		}
	}
	
	return nil
}