package parser

import (
	"go/ast"
	"go/parser"
	"go/token"
	"regexp"
	"strings"
)

// parseGo parses the Go code section
func (p *Parser) parseGo(section *Section) (*GoNode, error) {
	node := &GoNode{
		Source:   strings.TrimSpace(section.Content),
		Imports:  []string{},
		Handlers: []Handler{},
	}
	
	// Parse Go code using go/parser
	fset := token.NewFileSet()
	
	// Add a dummy package declaration if not present for parsing
	source := node.Source
	if !strings.HasPrefix(strings.TrimSpace(source), "package ") {
		source = "package main\n" + source
	}
	
	file, err := parser.ParseFile(fset, "", source, parser.ParseComments)
	if err != nil {
		return nil, &ParseError{
			File:    "",
			Line:    section.StartLine,
			Column:  section.StartCol,
			Message: "invalid Go syntax: " + err.Error(),
		}
	}
	
	// Extract imports
	for _, imp := range file.Imports {
		if imp.Path != nil {
			importPath := strings.Trim(imp.Path.Value, "\"")
			node.Imports = append(node.Imports, importPath)
		}
	}
	
	// Find the main struct type - prefer Page/Component structs
	type candidate struct {
		name       string
		structType *ast.StructType
		spec       *ast.TypeSpec
	}
	
	var candidates []candidate
	
	for _, decl := range file.Decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.TYPE {
			for _, spec := range genDecl.Specs {
				if typeSpec, ok := spec.(*ast.TypeSpec); ok {
					if structType, ok := typeSpec.Type.(*ast.StructType); ok {
						candidates = append(candidates, candidate{
							name:       typeSpec.Name.Name,
							structType: structType,
							spec:       typeSpec,
						})
					}
				}
			}
		}
	}
	
	// Select the best candidate - prefer Page/Component suffix
	var selected *candidate
	
	for i := range candidates {
		c := &candidates[i]
		if strings.HasSuffix(c.name, "Page") || 
		   strings.HasSuffix(c.name, "Component") ||
		   strings.HasSuffix(c.name, "Card") {
			selected = c
			break
		}
	}
	
	// If no Page/Component found, use the first struct
	if selected == nil && len(candidates) > 0 {
		selected = &candidates[0]
	}
	
	if selected != nil {
		node.MainType = selected.name
		// Parse struct fields as props for both components and pages
		// In GOX, all public fields are props
		node.Props = p.parsePropsStruct(selected.name, selected.structType, fset)
	}
	
	// Extract HTTP handlers and methods
	for _, decl := range file.Decls {
		if funcDecl, ok := decl.(*ast.FuncDecl); ok {
			handler := p.parseHandler(funcDecl, fset)
			if handler != nil {
				node.Handlers = append(node.Handlers, *handler)
			}
		}
	}
	
	return node, nil
}

// isComponentStruct checks if a struct is a component by looking for gox tags
func (p *Parser) isComponentStruct(structType *ast.StructType) bool {
	if structType.Fields == nil {
		return false
	}
	
	for _, field := range structType.Fields.List {
		if field.Tag != nil {
			tag := strings.Trim(field.Tag.Value, "`")
			if strings.Contains(tag, "gox:") {
				return true
			}
		}
	}
	
	return false
}

// parsePropsStruct parses a props struct for a component
func (p *Parser) parsePropsStruct(name string, structType *ast.StructType, fset *token.FileSet) *PropsStruct {
	props := &PropsStruct{
		Name:   name,
		Fields: []Field{},
	}
	
	if structType.Fields == nil {
		return props
	}
	
	for _, field := range structType.Fields.List {
		for _, name := range field.Names {
			f := Field{
				Name: name.Name,
				Type: p.typeToString(field.Type),
			}
			
			// Parse struct tags
			if field.Tag != nil {
				tag := strings.Trim(field.Tag.Value, "`")
				f.Tags = tag
				
				// Check if field is required
				// Look for gox:"required" or json:",required"
				if strings.Contains(tag, "required") || strings.Contains(tag, `gox:"required"`) {
					f.Required = true
				}
			}
			
			// Get position info
			pos := fset.Position(field.Pos())
			f.Line = pos.Line
			f.Column = pos.Column
			
			props.Fields = append(props.Fields, f)
		}
	}
	
	return props
}

// parseHandler parses a function declaration to extract handler information
func (p *Parser) parseHandler(funcDecl *ast.FuncDecl, fset *token.FileSet) *Handler {
	if funcDecl.Name == nil {
		return nil
	}
	
	handler := &Handler{
		Name:       funcDecl.Name.Name,
		Parameters: []Parameter{},
	}
	
	// Get position info
	pos := fset.Position(funcDecl.Pos())
	handler.Line = pos.Line
	handler.Column = pos.Column
	
	// Check if this is an HTTP handler
	if p.isHTTPHandler(funcDecl) {
		handler.IsHTMX = true
		
		// Try to extract HTTP method and path from function name or comments
		handler.Method, handler.Path = p.extractHTTPInfo(funcDecl)
	}
	
	// Parse parameters
	if funcDecl.Type.Params != nil {
		for _, param := range funcDecl.Type.Params.List {
			paramType := p.typeToString(param.Type)
			
			if len(param.Names) > 0 {
				for _, name := range param.Names {
					handler.Parameters = append(handler.Parameters, Parameter{
						Name: name.Name,
						Type: paramType,
					})
				}
			} else {
				// Anonymous parameter
				handler.Parameters = append(handler.Parameters, Parameter{
					Name: "",
					Type: paramType,
				})
			}
		}
	}
	
	// Parse return type
	if funcDecl.Type.Results != nil && len(funcDecl.Type.Results.List) > 0 {
		// For simplicity, just take the first return type
		handler.ReturnType = p.typeToString(funcDecl.Type.Results.List[0].Type)
	}
	
	return handler
}

// isHTTPHandler checks if a function is an HTTP handler
func (p *Parser) isHTTPHandler(funcDecl *ast.FuncDecl) bool {
	if funcDecl.Type.Params == nil || len(funcDecl.Type.Params.List) < 2 {
		return false
	}
	
	params := funcDecl.Type.Params.List
	
	// Check for http.ResponseWriter and *http.Request parameters
	if len(params) >= 2 {
		firstType := p.typeToString(params[0].Type)
		secondType := p.typeToString(params[1].Type)
		
		return (strings.Contains(firstType, "ResponseWriter") || firstType == "http.ResponseWriter") &&
			(strings.Contains(secondType, "Request") || secondType == "*http.Request")
	}
	
	return false
}

// extractHTTPInfo extracts HTTP method and path from function name or comments
func (p *Parser) extractHTTPInfo(funcDecl *ast.FuncDecl) (method, path string) {
	name := funcDecl.Name.Name
	
	// Try to extract from function name patterns
	methodPatterns := map[string]string{
		"HandleGet":    "GET",
		"HandlePost":   "POST",
		"HandlePut":    "PUT",
		"HandleDelete": "DELETE",
		"HandlePatch":  "PATCH",
		"Get":          "GET",
		"Post":         "POST",
		"Put":          "PUT",
		"Delete":       "DELETE",
		"Patch":        "PATCH",
	}
	
	for pattern, httpMethod := range methodPatterns {
		if strings.HasPrefix(name, pattern) {
			method = httpMethod
			break
		}
	}
	
	// Try to extract path from function name
	// HandleGetUsers -> /users
	// HandlePostUser -> /user
	if strings.HasPrefix(name, "Handle") {
		// Remove "Handle" and HTTP method
		remaining := name[6:] // Remove "Handle"
		for httpMethod := range methodPatterns {
			if strings.HasPrefix(remaining, httpMethod) {
				remaining = remaining[len(httpMethod):]
				break
			}
		}
		
		if remaining != "" {
			// Convert PascalCase to kebab-case
			path = "/" + p.pascalToKebab(remaining)
		}
	}
	
	return method, path
}

// typeToString converts an AST type to string representation
func (p *Parser) typeToString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return "*" + p.typeToString(t.X)
	case *ast.SelectorExpr:
		return p.typeToString(t.X) + "." + t.Sel.Name
	case *ast.ArrayType:
		return "[]" + p.typeToString(t.Elt)
	case *ast.MapType:
		return "map[" + p.typeToString(t.Key) + "]" + p.typeToString(t.Value)
	case *ast.InterfaceType:
		return "interface{}"
	case *ast.StructType:
		return "struct{}"
	case *ast.FuncType:
		return "func"
	default:
		return "unknown"
	}
}

// pascalToKebab converts PascalCase to kebab-case
func (p *Parser) pascalToKebab(s string) string {
	// Use regex to insert hyphens before uppercase letters
	re := regexp.MustCompile("([a-z])([A-Z])")
	return strings.ToLower(re.ReplaceAllString(s, "$1-$2"))
}