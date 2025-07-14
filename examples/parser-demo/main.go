package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/gox-framework/gox/internal/parser"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Uso: go run main.go <archivo.gox>")
		fmt.Println("\nEjemplos:")
		fmt.Println("  go run main.go ../../examples/my-awesome-app/components/user-card.gox")
		fmt.Println("  go run main.go ../../examples/my-awesome-app/pages/home.gox")
		os.Exit(1)
	}

	filename := os.Args[1]
	
	fmt.Printf("🔍 Parseando archivo: %s\n\n", filename)
	
	// Parsear el archivo
	ast, err := parser.ParseFile(filename)
	if err != nil {
		log.Fatalf("❌ Error parseando archivo: %v", err)
	}
	
	// 1. Información básica
	fmt.Println("📋 INFORMACIÓN BÁSICA:")
	fmt.Printf("  - Path: %s\n", ast.Path)
	fmt.Printf("  - Es componente: %v\n", ast.IsComponent())
	fmt.Printf("  - Es página: %v\n", ast.IsPage())
	fmt.Printf("  - Tiene template: %v\n", ast.HasTemplate())
	fmt.Printf("  - Tiene código Go: %v\n", ast.HasGo())
	fmt.Printf("  - Tiene estilos: %v\n", ast.HasStyles())
	
	// 2. Template info
	if ast.Template != nil {
		fmt.Println("\n📝 TEMPLATE:")
		fmt.Printf("  - Auth: %s\n", ast.Template.Auth)
		fmt.Printf("  - Layout: %s\n", ast.Template.Layout)
		fmt.Printf("  - Contenido (primeros 100 chars): %.100s...\n", ast.Template.Content)
		
		// Detectar variables de template
		detector := parser.NewComponentDetector()
		vars := detector.DetectGoTemplateVariables(ast.Template.Content)
		if len(vars) > 0 {
			fmt.Printf("  - Variables Go Template detectadas: %v\n", vars)
		}
		
		// Detectar atributos HTMX
		htmxAttrs := detector.DetectHTMXAttributes(ast.Template.Content)
		if len(htmxAttrs) > 0 {
			fmt.Printf("  - Atributos HTMX detectados: %v\n", htmxAttrs)
		}
	}
	
	// 3. Go code info
	if ast.Go != nil {
		fmt.Println("\n🐹 CÓDIGO GO:")
		fmt.Printf("  - Tipo principal: %s\n", ast.Go.MainType)
		fmt.Printf("  - Imports: %v\n", ast.Go.Imports)
		
		if ast.Go.Props != nil {
			fmt.Printf("\n  📌 Props del componente '%s':\n", ast.Go.Props.Name)
			for _, field := range ast.Go.Props.Fields {
				required := ""
				if field.Required {
					required = " (required)"
				}
				fmt.Printf("    - %s: %s%s\n", field.Name, field.Type, required)
			}
		}
		
		if len(ast.Go.Handlers) > 0 {
			fmt.Println("\n  🔌 Handlers HTTP detectados:")
			for _, handler := range ast.Go.Handlers {
				fmt.Printf("    - %s", handler.Name)
				if handler.IsHTMX {
					fmt.Print(" [HTMX]")
				}
				if handler.Method != "" {
					fmt.Printf(" %s %s", handler.Method, handler.Path)
				}
				fmt.Println()
			}
		}
	}
	
	// 4. Styles info
	if len(ast.Styles) > 0 {
		fmt.Println("\n🎨 ESTILOS:")
		for i, style := range ast.Styles {
			fmt.Printf("  - Estilo %d:\n", i+1)
			fmt.Printf("    - Scoped: %v\n", style.Scoped)
			fmt.Printf("    - Tipo: %s\n", style.Type)
			fmt.Printf("    - Contenido (primeros 50 chars): %.50s...\n", style.Content)
		}
	}
	
	// 5. Component dependencies
	if len(ast.Components) > 0 {
		fmt.Println("\n🧩 COMPONENTES USADOS:")
		for _, comp := range ast.Components {
			fmt.Printf("  - %s (path: %s, usado %d veces)\n", comp.Name, comp.Path, comp.UsageCount)
			if len(comp.Props) > 0 {
				fmt.Printf("    Props: %v\n", comp.Props)
			}
		}
	}
	
	// 6. Export JSON (opcional)
	if len(os.Args) > 2 && os.Args[2] == "--json" {
		fmt.Println("\n📄 AST EN JSON:")
		jsonData, _ := json.MarshalIndent(ast, "", "  ")
		fmt.Println(string(jsonData))
	}
	
	// 7. Reconstruir archivo
	fmt.Println("\n🔄 ARCHIVO RECONSTRUIDO:")
	fmt.Println("```gox")
	fmt.Println(ast.String())
	fmt.Println("```")
}