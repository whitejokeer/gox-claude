package parser

import (
	"regexp"
	"strings"
)

// parseTemplate parses the template section
func (p *Parser) parseTemplate(section *Section) (*TemplateNode, error) {
	node := &TemplateNode{
		Content:  strings.TrimSpace(section.Content),
		Elements: []Element{},
	}
	
	// Extract attributes from section
	if auth, exists := section.Attributes["auth"]; exists {
		node.Auth = auth
	}
	
	if layout, exists := section.Attributes["layout"]; exists {
		node.Layout = layout
	}
	
	// Parse HTML structure (basic implementation)
	elements, err := p.parseHTMLElements(node.Content)
	if err != nil {
		return nil, err
	}
	node.Elements = elements
	
	return node, nil
}

// parseHTMLElements parses HTML elements from template content
func (p *Parser) parseHTMLElements(content string) ([]Element, error) {
	// This is a simplified HTML parser
	// In a production implementation, you might want to use a proper HTML parser
	// or implement a more sophisticated parsing algorithm
	
	var elements []Element
	
	// For now, we'll do a basic regex-based parsing to identify components
	// This will be enhanced in later iterations
	
	// Find all HTML tags and components
	tagRegex := regexp.MustCompile(`<([a-zA-Z][a-zA-Z0-9-]*)[^>]*>`)
	matches := tagRegex.FindAllStringSubmatch(content, -1)
	
	for _, match := range matches {
		if len(match) < 2 {
			continue
		}
		
		tagName := match[1]
		isComponent := p.isComponentTag(tagName)
		
		element := Element{
			Type:        tagName,
			IsComponent: isComponent,
			Props:       make(map[string]string),
			Children:    []Element{},
		}
		
		// Parse attributes/props from the full match
		fullTag := match[0]
		props := p.parseElementProps(fullTag)
		element.Props = props
		
		elements = append(elements, element)
	}
	
	return elements, nil
}

// isComponentTag determines if a tag is a custom component
func (p *Parser) isComponentTag(tagName string) bool {
	// Standard HTML tags are not components
	standardTags := map[string]bool{
		"a": true, "abbr": true, "address": true, "area": true, "article": true,
		"aside": true, "audio": true, "b": true, "base": true, "bdi": true,
		"bdo": true, "blockquote": true, "body": true, "br": true, "button": true,
		"canvas": true, "caption": true, "cite": true, "code": true, "col": true,
		"colgroup": true, "data": true, "datalist": true, "dd": true, "del": true,
		"details": true, "dfn": true, "dialog": true, "div": true, "dl": true,
		"dt": true, "em": true, "embed": true, "fieldset": true, "figcaption": true,
		"figure": true, "footer": true, "form": true, "h1": true, "h2": true,
		"h3": true, "h4": true, "h5": true, "h6": true, "head": true, "header": true,
		"hr": true, "html": true, "i": true, "iframe": true, "img": true,
		"input": true, "ins": true, "kbd": true, "label": true, "legend": true,
		"li": true, "link": true, "main": true, "map": true, "mark": true,
		"meta": true, "meter": true, "nav": true, "noscript": true, "object": true,
		"ol": true, "optgroup": true, "option": true, "output": true, "p": true,
		"param": true, "picture": true, "pre": true, "progress": true, "q": true,
		"rp": true, "rt": true, "ruby": true, "s": true, "samp": true,
		"script": true, "section": true, "select": true, "small": true, "source": true,
		"span": true, "strong": true, "style": true, "sub": true, "summary": true,
		"sup": true, "svg": true, "table": true, "tbody": true, "td": true,
		"template": true, "textarea": true, "tfoot": true, "th": true, "thead": true,
		"time": true, "title": true, "tr": true, "track": true, "u": true,
		"ul": true, "var": true, "video": true, "wbr": true,
	}
	
	return !standardTags[strings.ToLower(tagName)]
}

// parseElementProps parses props/attributes from an HTML tag
func (p *Parser) parseElementProps(tag string) map[string]string {
	props := make(map[string]string)
	
	// Remove < and > from tag
	tag = strings.TrimPrefix(tag, "<")
	tag = strings.TrimSuffix(tag, ">")
	tag = strings.TrimSuffix(tag, "/") // Handle self-closing tags
	
	// Split by spaces to get tag name and attributes
	parts := strings.Fields(tag)
	if len(parts) <= 1 {
		return props
	}
	
	// Parse attributes (simplified) - includes : prefix for Vue-style binding
	attrRegex := regexp.MustCompile(`([:a-zA-Z][a-zA-Z0-9:-]*)\s*=\s*["']([^"']*)["']`)
	matches := attrRegex.FindAllStringSubmatch(tag, -1)
	
	for _, match := range matches {
		if len(match) >= 3 {
			attrName := match[1]
			attrValue := match[2]
			props[attrName] = attrValue
		}
	}
	
	// Also handle boolean attributes (no value)
	boolAttrRegex := regexp.MustCompile(`\s([a-zA-Z:][a-zA-Z0-9:-]*)\s`)
	boolMatches := boolAttrRegex.FindAllStringSubmatch(" "+tag+" ", -1)
	
	for _, match := range boolMatches {
		if len(match) >= 2 {
			attrName := match[1]
			// Don't override attributes that already have values
			if _, exists := props[attrName]; !exists {
				// Check if this is actually a boolean attribute by seeing if it's followed by =
				if !strings.Contains(tag, attrName+"=") {
					props[attrName] = ""
				}
			}
		}
	}
	
	return props
}

// extractGoTemplateVariables finds Go template variables in content
func (p *Parser) extractGoTemplateVariables(content string) []string {
	var variables []string
	
	// Find {{.Variable}} patterns
	varRegex := regexp.MustCompile(`\{\{\s*\.([a-zA-Z][a-zA-Z0-9]*)\s*\}\}`)
	matches := varRegex.FindAllStringSubmatch(content, -1)
	
	seen := make(map[string]bool)
	for _, match := range matches {
		if len(match) >= 2 {
			variable := match[1]
			if !seen[variable] {
				variables = append(variables, variable)
				seen[variable] = true
			}
		}
	}
	
	return variables
}

// extractHTMXAttributes finds HTMX attributes in content
func (p *Parser) extractHTMXAttributes(content string) []string {
	var attributes []string
	
	// Find hx-* attributes
	hxRegex := regexp.MustCompile(`hx-([a-zA-Z-]+)`)
	matches := hxRegex.FindAllStringSubmatch(content, -1)
	
	seen := make(map[string]bool)
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