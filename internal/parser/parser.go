package parser

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Parser is the main parser struct for .gox files
type Parser struct {
	fs FileSystem
}

// FileSystem interface to allow for testing with mock filesystems
type FileSystem interface {
	ReadFile(filename string) ([]byte, error)
	Exists(filename string) bool
}

// DefaultFileSystem implements FileSystem using the standard library
type DefaultFileSystem struct{}

func (fs DefaultFileSystem) ReadFile(filename string) ([]byte, error) {
	return os.ReadFile(filename)
}

func (fs DefaultFileSystem) Exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

// NewParser creates a new parser instance
func NewParser(options ...ParserOption) *Parser {
	p := &Parser{
		fs: DefaultFileSystem{},
	}
	
	for _, opt := range options {
		opt(p)
	}
	
	return p
}

// ParserOption configures the parser
type ParserOption func(*Parser)

// WithFilesystem sets a custom filesystem for the parser
func WithFilesystem(fs FileSystem) ParserOption {
	return func(p *Parser) {
		p.fs = fs
	}
}

// ParseFile parses a .gox file from the filesystem
func (p *Parser) ParseFile(filename string) (*GoxFile, error) {
	content, err := p.fs.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filename, err)
	}
	
	return p.ParseBytes(content, filename)
}

// ParseString parses a .gox file from a string
func (p *Parser) ParseString(content string) (*GoxFile, error) {
	return p.ParseBytes([]byte(content), "")
}

// ParseBytes parses a .gox file from byte content
func (p *Parser) ParseBytes(content []byte, filename string) (*GoxFile, error) {
	lexer := NewLexer(content, filename)
	tokens, err := lexer.Tokenize()
	if err != nil {
		return nil, err
	}
	
	return p.parseTokens(tokens, filename)
}

// ParseReader parses a .gox file from an io.Reader
func (p *Parser) ParseReader(r io.Reader, filename string) (*GoxFile, error) {
	content, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read content: %w", err)
	}
	
	return p.ParseBytes(content, filename)
}

// parseTokens converts tokens into AST
func (p *Parser) parseTokens(tokens []Token, filename string) (*GoxFile, error) {
	goxFile := &GoxFile{
		Path:       filename,
		Components: []ComponentDependency{},
		Metadata:   make(map[string]interface{}),
	}
	
	var currentSection *Section
	
	for _, token := range tokens {
		switch token.Type {
		case TokenSectionStart:
			currentSection = &Section{
				Type:       token.Value,
				Attributes: token.Attributes,
				StartLine:  token.Line,
				StartCol:   token.Column,
			}
			
		case TokenSectionEnd:
			if currentSection == nil {
				return nil, &ParseError{
					File:    filename,
					Line:    token.Line,
					Column:  token.Column,
					Message: fmt.Sprintf("unexpected closing tag </%s>", token.Value),
				}
			}
			
			if currentSection.Type != token.Value {
				return nil, &ParseError{
					File:    filename,
					Line:    token.Line,
					Column:  token.Column,
					Message: fmt.Sprintf("mismatched tags: expected </%s>, got </%s>", currentSection.Type, token.Value),
				}
			}
			
			if err := p.processSection(goxFile, currentSection); err != nil {
				return nil, err
			}
			currentSection = nil
			
		case TokenContent:
			if currentSection == nil {
				return nil, &ParseError{
					File:    filename,
					Line:    token.Line,
					Column:  token.Column,
					Message: "content outside of section",
				}
			}
			currentSection.Content += token.Value
		}
	}
	
	if currentSection != nil {
		return nil, &ParseError{
			File:    filename,
			Line:    currentSection.StartLine,
			Column:  currentSection.StartCol,
			Message: fmt.Sprintf("unclosed section <%s>", currentSection.Type),
		}
	}
	
	// Detect component dependencies
	if err := p.detectComponents(goxFile); err != nil {
		return nil, err
	}
	
	return goxFile, nil
}

// processSection processes a completed section
func (p *Parser) processSection(goxFile *GoxFile, section *Section) error {
	switch section.Type {
	case "template":
		templateNode, err := p.parseTemplate(section)
		if err != nil {
			return err
		}
		goxFile.Template = templateNode
		
	case "go":
		goNode, err := p.parseGo(section)
		if err != nil {
			return err
		}
		goxFile.Go = goNode
		
	case "style":
		styleNode, err := p.parseStyle(section)
		if err != nil {
			return err
		}
		goxFile.Styles = append(goxFile.Styles, styleNode)
		
	default:
		return &ParseError{
			File:    goxFile.Path,
			Line:    section.StartLine,
			Column:  section.StartCol,
			Message: fmt.Sprintf("unknown section type: %s", section.Type),
		}
	}
	
	return nil
}

// detectComponents finds component dependencies in the template
func (p *Parser) detectComponents(goxFile *GoxFile) error {
	if goxFile.Template == nil {
		return nil
	}
	
	detector := NewComponentDetector()
	components, err := detector.DetectComponents(goxFile.Template.Content)
	if err != nil {
		return err
	}
	
	// Resolve component paths
	for _, comp := range components {
		resolved := p.resolveComponentPath(comp.Name)
		comp.Path = resolved
		goxFile.Components = append(goxFile.Components, comp)
	}
	
	return nil
}

// resolveComponentPath resolves component names to file paths
func (p *Parser) resolveComponentPath(componentName string) string {
	// Convert kebab-case to path
	// user-card -> components/user-card
	// shared-button -> shared/ui/button
	
	if len(componentName) > 7 && componentName[:7] == "shared-" {
		// shared components go to shared/ui/
		name := componentName[7:] // remove "shared-" prefix
		return filepath.Join("shared", "ui", name)
	}
	
	// regular components go to components/
	return filepath.Join("components", componentName)
}

// Section represents a parsed section from the .gox file
type Section struct {
	Type       string
	Content    string
	Attributes map[string]string
	StartLine  int
	StartCol   int
}