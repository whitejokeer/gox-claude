package parser

import (
	"io"
	"os"
)

// ParseFile is a convenience function that creates a parser and parses a file
func ParseFile(filename string) (*GoxFile, error) {
	parser := NewParser()
	return parser.ParseFile(filename)
}

// ParseString is a convenience function that creates a parser and parses a string
func ParseString(content string) (*GoxFile, error) {
	parser := NewParser()
	return parser.ParseString(content)
}

// ParseBytes is a convenience function that creates a parser and parses bytes
func ParseBytes(content []byte, filename string) (*GoxFile, error) {
	parser := NewParser()
	return parser.ParseBytes(content, filename)
}

// ParseReader is a convenience function that creates a parser and parses from a reader
func ParseReader(r io.Reader, filename string) (*GoxFile, error) {
	parser := NewParser()
	return parser.ParseReader(r, filename)
}

// MustParseFile parses a file and panics on error (useful for tests)
func MustParseFile(filename string) *GoxFile {
	ast, err := ParseFile(filename)
	if err != nil {
		panic(err)
	}
	return ast
}

// MustParseString parses a string and panics on error (useful for tests)
func MustParseString(content string) *GoxFile {
	ast, err := ParseString(content)
	if err != nil {
		panic(err)
	}
	return ast
}

// ValidateGoxFile validates a .gox file and returns any errors
func ValidateGoxFile(filename string) error {
	ast, err := ParseFile(filename)
	if err != nil {
		return err
	}
	return ast.Validate()
}

// GetComponentDependencies extracts component dependencies from a .gox file
func GetComponentDependencies(filename string) ([]ComponentDependency, error) {
	ast, err := ParseFile(filename)
	if err != nil {
		return nil, err
	}
	return ast.Components, nil
}

// IsGoxFile checks if a file has a .gox extension
func IsGoxFile(filename string) bool {
	return len(filename) > 4 && filename[len(filename)-4:] == ".gox"
}

// ParseDirectory parses all .gox files in a directory
func ParseDirectory(dirPath string) (map[string]*GoxFile, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}
	
	results := make(map[string]*GoxFile)
	parser := NewParser()
	
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		
		filename := entry.Name()
		if !IsGoxFile(filename) {
			continue
		}
		
		fullPath := dirPath + "/" + filename
		ast, err := parser.ParseFile(fullPath)
		if err != nil {
			return nil, err
		}
		
		results[filename] = ast
	}
	
	return results, nil
}