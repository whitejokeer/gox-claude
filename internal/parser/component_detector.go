package parser

import (
	"regexp"
	"strings"
)

// ComponentDetector detects component usage in templates
type ComponentDetector struct {
	componentRegex *regexp.Regexp
}

// NewComponentDetector creates a new component detector
func NewComponentDetector() *ComponentDetector {
	// Regex to match custom component tags
	// Matches: <user-card name="..." email="..." />
	// Matches: <shared-button text="..." hx-post="..." />
	componentRegex := regexp.MustCompile(`<([a-z][a-z0-9]*(?:-[a-z0-9]+)+)([^>]*?)(/?)>`)
	
	return &ComponentDetector{
		componentRegex: componentRegex,
	}
}

// DetectComponents finds all component usages in template content
func (cd *ComponentDetector) DetectComponents(content string) ([]ComponentDependency, error) {
	var components []ComponentDependency
	usageCount := make(map[string]int)
	
	// Find all component tags
	matches := cd.componentRegex.FindAllStringSubmatch(content, -1)
	
	for _, match := range matches {
		if len(match) < 4 {
			continue
		}
		
		componentName := match[1]
		attributesStr := match[2]
		// isSelfClosing := match[3] == "/"
		
		// Parse component props/attributes
		props := cd.parseComponentProps(attributesStr)
		
		// Count usage
		usageCount[componentName]++
		
		// Find the component in our list or create new entry
		var existingComponent *ComponentDependency
		for i := range components {
			if components[i].Name == componentName {
				existingComponent = &components[i]
				break
			}
		}
		
		if existingComponent == nil {
			// Create new component dependency
			component := ComponentDependency{
				Name:       componentName,
				Props:      props,
				UsageCount: 1,
			}
			components = append(components, component)
		} else {
			// Update existing component
			existingComponent.UsageCount++
			// Merge props (could have different props in different usages)
			for key, value := range props {
				existingComponent.Props[key] = value
			}
		}
	}
	
	return components, nil
}

// parseComponentProps parses props/attributes from component tag
func (cd *ComponentDetector) parseComponentProps(attributesStr string) map[string]string {
	props := make(map[string]string)
	
	if strings.TrimSpace(attributesStr) == "" {
		return props
	}
	
	// Regex to match attribute="value" or attribute='value' - includes : prefix
	attrRegex := regexp.MustCompile(`([:a-zA-Z][a-zA-Z0-9:_-]*)\s*=\s*["']([^"']*)["']`)
	matches := attrRegex.FindAllStringSubmatch(attributesStr, -1)
	
	for _, match := range matches {
		if len(match) >= 3 {
			attrName := match[1]
			attrValue := match[2]
			props[attrName] = attrValue
		}
	}
	
	// Handle boolean attributes by parsing the string manually
	// We need to find attribute names that are not part of key=value pairs
	remaining := strings.TrimSpace(attributesStr)
	
	for remaining != "" {
		// Skip spaces
		remaining = strings.TrimLeft(remaining, " \t\n\r")
		if remaining == "" {
			break
		}
		
		// Find the next space or end of string
		spaceIndex := strings.IndexAny(remaining, " \t\n\r")
		var word string
		if spaceIndex == -1 {
			word = remaining
			remaining = ""
		} else {
			word = remaining[:spaceIndex]
			remaining = remaining[spaceIndex:]
		}
		
		// Skip if this word contains = (it's a key=value pair)
		if strings.Contains(word, "=") {
			continue
		}
		
		// Check if it's a valid attribute name
		if matched, _ := regexp.MatchString(`^[a-zA-Z:][a-zA-Z0-9:_-]*$`, word); matched {
			if _, exists := props[word]; !exists {
				props[word] = ""
			}
		}
	}
	
	return props
}

// DetectHTMXAttributes finds HTMX attributes in template content
func (cd *ComponentDetector) DetectHTMXAttributes(content string) []string {
	var attributes []string
	seen := make(map[string]bool)
	
	// Regex to find hx-* attributes
	hxRegex := regexp.MustCompile(`hx-([a-zA-Z-]+)`)
	matches := hxRegex.FindAllStringSubmatch(content, -1)
	
	for _, match := range matches {
		if len(match) >= 2 {
			attr := "hx-" + match[1]
			if !seen[attr] {
				attributes = append(attributes, attr)
				seen[attr] = true
			}
		}
	}
	
	return attributes
}

// DetectGoTemplateVariables finds Go template variables in content
func (cd *ComponentDetector) DetectGoTemplateVariables(content string) []string {
	var variables []string
	seen := make(map[string]bool)
	
	// Regex to find {{.Variable}} and {{.Nested.Field}} patterns
	varRegex := regexp.MustCompile(`\{\{\s*\.([a-zA-Z][a-zA-Z0-9]*(?:\.[a-zA-Z][a-zA-Z0-9]*)*)\s*\}\}`)
	matches := varRegex.FindAllStringSubmatch(content, -1)
	
	for _, match := range matches {
		if len(match) >= 2 {
			fullVariable := match[1]
			// Extract the root variable (before the first dot)
			parts := strings.Split(fullVariable, ".")
			variable := parts[0]
			if !seen[variable] {
				variables = append(variables, variable)
				seen[variable] = true
			}
		}
	}
	
	return variables
}

// DetectGoTemplateFunctions finds Go template function calls in content
func (cd *ComponentDetector) DetectGoTemplateFunctions(content string) []string {
	var functions []string
	seen := make(map[string]bool)
	
	// Regex to find function calls like {{functionName .arg}}
	funcRegex := regexp.MustCompile(`\{\{\s*([a-zA-Z][a-zA-Z0-9]*)\s+`)
	matches := funcRegex.FindAllStringSubmatch(content, -1)
	
	for _, match := range matches {
		if len(match) >= 2 {
			function := match[1]
			// Skip common template keywords
			if !cd.isTemplateKeyword(function) && !seen[function] {
				functions = append(functions, function)
				seen[function] = true
			}
		}
	}
	
	return functions
}

// isTemplateKeyword checks if a word is a Go template keyword
func (cd *ComponentDetector) isTemplateKeyword(word string) bool {
	keywords := map[string]bool{
		"if":     true,
		"else":   true,
		"end":    true,
		"range":  true,
		"with":   true,
		"define": true,
		"block":  true,
		"template": true,
	}
	
	return keywords[word]
}

// IsValidComponentName checks if a tag name is a valid component name
func (cd *ComponentDetector) IsValidComponentName(name string) bool {
	// Component names must:
	// 1. Start with lowercase letter
	// 2. Contain at least one hyphen
	// 3. Only contain lowercase letters, numbers, and hyphens
	// 4. Not start or end with hyphen
	
	if len(name) < 3 { // Minimum length for valid component name
		return false
	}
	
	if name[0] < 'a' || name[0] > 'z' {
		return false
	}
	
	if name[len(name)-1] == '-' {
		return false
	}
	
	hasHyphen := false
	for i, char := range name {
		if char == '-' {
			hasHyphen = true
			// No consecutive hyphens
			if i > 0 && name[i-1] == '-' {
				return false
			}
		} else if !((char >= 'a' && char <= 'z') || (char >= '0' && char <= '9')) {
			return false
		}
	}
	
	return hasHyphen
}

// ExtractComponentsFromHTML extracts components from HTML content with more detail
func (cd *ComponentDetector) ExtractComponentsFromHTML(content string) ([]ComponentDependency, error) {
	var components []ComponentDependency
	
	// Find all custom component tags with their positions
	lines := strings.Split(content, "\n")
	
	for lineNum, line := range lines {
		matches := cd.componentRegex.FindAllStringSubmatch(line, -1)
		
		for _, match := range matches {
			if len(match) < 4 {
				continue
			}
			
			componentName := match[1]
			
			// Validate component name
			if !cd.IsValidComponentName(componentName) {
				continue
			}
			
			attributesStr := match[2]
			props := cd.parseComponentProps(attributesStr)
			
			// Find column position
			columnPos := strings.Index(line, match[0]) + 1
			
			component := ComponentDependency{
				Name:       componentName,
				Props:      props,
				Line:       lineNum + 1,
				Column:     columnPos,
				UsageCount: 1,
			}
			
			components = append(components, component)
		}
	}
	
	return components, nil
}