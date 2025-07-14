package parser

import (
	"regexp"
	"strings"
)

// parseStyle parses the style section
func (p *Parser) parseStyle(section *Section) (*StyleNode, error) {
	node := &StyleNode{
		Content: strings.TrimSpace(section.Content),
		Scoped:  false,
		Type:    "css", // default to CSS
		Line:    section.StartLine,
		Column:  section.StartCol,
	}
	
	// Check for scoped attribute
	if _, exists := section.Attributes["scoped"]; exists {
		node.Scoped = true
	}
	
	// Detect style type
	node.Type = p.detectStyleType(node.Content)
	
	// Basic CSS syntax validation
	if err := p.validateCSS(node.Content); err != nil {
		return nil, &ParseError{
			File:    "",
			Line:    section.StartLine,
			Column:  section.StartCol,
			Message: "invalid CSS syntax: " + err.Error(),
		}
	}
	
	return node, nil
}

// detectStyleType detects whether the style content is CSS, Tailwind, or other
func (p *Parser) detectStyleType(content string) string {
	content = strings.TrimSpace(content)
	
	// Check for Tailwind @apply directives
	if strings.Contains(content, "@apply") {
		return "tailwind"
	}
	
	// Check for SCSS/SASS features
	if strings.Contains(content, "&") || strings.Contains(content, "$") || strings.Contains(content, "@mixin") {
		return "scss"
	}
	
	// Check for CSS custom properties
	if strings.Contains(content, "--") {
		return "css-custom-properties"
	}
	
	// Default to CSS
	return "css"
}

// validateCSS performs basic CSS syntax validation
func (p *Parser) validateCSS(content string) error {
	// Remove comments first
	content = p.removeComments(content)
	
	// Check for balanced braces
	if err := p.checkBalancedBraces(content); err != nil {
		return err
	}
	
	// Check for basic CSS structure
	if err := p.validateCSSStructure(content); err != nil {
		return err
	}
	
	return nil
}

// removeComments removes CSS comments from content
func (p *Parser) removeComments(content string) string {
	// Remove /* ... */ comments
	commentRegex := regexp.MustCompile(`/\*.*?\*/`)
	return commentRegex.ReplaceAllString(content, "")
}

// checkBalancedBraces checks if braces are balanced
func (p *Parser) checkBalancedBraces(content string) error {
	count := 0
	for _, char := range content {
		switch char {
		case '{':
			count++
		case '}':
			count--
			if count < 0 {
				return &ParseError{
					Message: "unexpected closing brace '}'",
				}
			}
		}
	}
	
	if count > 0 {
		return &ParseError{
			Message: "unclosed opening brace '{'",
		}
	}
	
	return nil
}

// validateCSSStructure validates basic CSS structure
func (p *Parser) validateCSSStructure(content string) error {
	content = strings.TrimSpace(content)
	if content == "" {
		return nil // Empty styles are valid
	}
	
	// Split into rules
	rules := p.splitCSSRules(content)
	
	for i, rule := range rules {
		rule = strings.TrimSpace(rule)
		if rule == "" {
			continue
		}
		
		// Check if rule has proper structure
		if err := p.validateCSSRule(rule, i+1); err != nil {
			return err
		}
	}
	
	return nil
}

// splitCSSRules splits CSS content into individual rules
func (p *Parser) splitCSSRules(content string) []string {
	var rules []string
	var currentRule strings.Builder
	braceCount := 0
	
	for _, char := range content {
		currentRule.WriteRune(char)
		
		switch char {
		case '{':
			braceCount++
		case '}':
			braceCount--
			if braceCount == 0 {
				rules = append(rules, currentRule.String())
				currentRule.Reset()
			}
		}
	}
	
	// Add any remaining content
	if currentRule.Len() > 0 {
		rules = append(rules, currentRule.String())
	}
	
	return rules
}

// validateCSSRule validates a single CSS rule
func (p *Parser) validateCSSRule(rule string, ruleNumber int) error {
	rule = strings.TrimSpace(rule)
	
	// Handle at-rules (@apply, @media, etc.)
	if strings.HasPrefix(rule, "@") {
		return p.validateAtRule(rule, ruleNumber)
	}
	
	// Find the opening brace
	braceIndex := strings.Index(rule, "{")
	if braceIndex == -1 {
		return &ParseError{
			Message: "CSS rule missing opening brace",
		}
	}
	
	// Extract selector and declarations
	selector := strings.TrimSpace(rule[:braceIndex])
	declarations := strings.TrimSpace(rule[braceIndex+1:])
	
	// Remove closing brace
	if strings.HasSuffix(declarations, "}") {
		declarations = strings.TrimSpace(declarations[:len(declarations)-1])
	}
	
	// Validate selector
	if err := p.validateCSSSelector(selector, ruleNumber); err != nil {
		return err
	}
	
	// Validate declarations
	if err := p.validateCSSDeclarations(declarations, ruleNumber); err != nil {
		return err
	}
	
	return nil
}

// validateAtRule validates CSS at-rules
func (p *Parser) validateAtRule(rule string, ruleNumber int) error {
	// Basic validation for common at-rules
	validAtRules := []string{"@apply", "@media", "@import", "@charset", "@keyframes", "@font-face"}
	
	for _, validRule := range validAtRules {
		if strings.HasPrefix(rule, validRule) {
			return nil // Valid at-rule
		}
	}
	
	// Unknown at-rule - allow it but could warn
	return nil
}

// validateCSSSelector validates a CSS selector
func (p *Parser) validateCSSSelector(selector string, ruleNumber int) error {
	if selector == "" {
		return &ParseError{
			Message: "empty CSS selector",
		}
	}
	
	// Basic selector validation (could be more comprehensive)
	// Allow common selector patterns
	selectorRegex := regexp.MustCompile(`^[a-zA-Z0-9\s\-_#.:\[\]='"(),>+~*]+$`)
	if !selectorRegex.MatchString(selector) {
		return &ParseError{
			Message: "invalid CSS selector: " + selector,
		}
	}
	
	return nil
}

// validateCSSDeclarations validates CSS declarations
func (p *Parser) validateCSSDeclarations(declarations string, ruleNumber int) error {
	if declarations == "" {
		return nil // Empty declarations are valid
	}
	
	// Split declarations by semicolon
	decls := strings.Split(declarations, ";")
	
	for _, decl := range decls {
		decl = strings.TrimSpace(decl)
		if decl == "" {
			continue
		}
		
		// Validate individual declaration
		if err := p.validateCSSDeclaration(decl, ruleNumber); err != nil {
			return err
		}
	}
	
	return nil
}

// validateCSSDeclaration validates a single CSS declaration
func (p *Parser) validateCSSDeclaration(declaration string, ruleNumber int) error {
	// Declaration should have format: property: value
	colonIndex := strings.Index(declaration, ":")
	if colonIndex == -1 {
		return &ParseError{
			Message: "CSS declaration missing colon: " + declaration,
		}
	}
	
	property := strings.TrimSpace(declaration[:colonIndex])
	value := strings.TrimSpace(declaration[colonIndex+1:])
	
	if property == "" {
		return &ParseError{
			Message: "empty CSS property",
		}
	}
	
	if value == "" {
		return &ParseError{
			Message: "empty CSS value for property: " + property,
		}
	}
	
	// Basic property name validation
	propertyRegex := regexp.MustCompile(`^[a-zA-Z-]+$`)
	if !propertyRegex.MatchString(property) {
		return &ParseError{
			Message: "invalid CSS property name: " + property,
		}
	}
	
	return nil
}

// extractCSSSelectors extracts all CSS selectors from the content
func (p *Parser) extractCSSSelectors(content string) []string {
	var selectors []string
	
	// Remove comments
	content = p.removeComments(content)
	
	// Find all selectors before opening braces
	selectorRegex := regexp.MustCompile(`([^{}]+)\s*\{`)
	matches := selectorRegex.FindAllStringSubmatch(content, -1)
	
	for _, match := range matches {
		if len(match) >= 2 {
			selector := strings.TrimSpace(match[1])
			if selector != "" && !strings.HasPrefix(selector, "@") {
				selectors = append(selectors, selector)
			}
		}
	}
	
	return selectors
}

// extractCSSProperties extracts all CSS properties from the content
func (p *Parser) extractCSSProperties(content string) []string {
	var properties []string
	seen := make(map[string]bool)
	
	// Remove comments
	content = p.removeComments(content)
	
	// Find all property declarations
	propertyRegex := regexp.MustCompile(`([a-zA-Z-]+)\s*:`)
	matches := propertyRegex.FindAllStringSubmatch(content, -1)
	
	for _, match := range matches {
		if len(match) >= 2 {
			property := strings.TrimSpace(match[1])
			if property != "" && !seen[property] {
				properties = append(properties, property)
				seen[property] = true
			}
		}
	}
	
	return properties
}