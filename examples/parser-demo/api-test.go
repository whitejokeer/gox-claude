package main

import (
	"fmt"
	"log"

	"github.com/gox-framework/gox/internal/parser"
)

func testParserAPI() {
	// Test 1: Parser simple desde string
	fmt.Println("🧪 TEST 1: Parser desde string")
	simpleGox := `<template>
  <h1>{{.Title}}</h1>
  <user-card name="John" email="john@example.com" />
</template>

<go>
package pages

type SimplePage struct {
    Title string
}
</go>`

	ast, err := parser.ParseString(simpleGox)
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Printf("✅ Parseado exitosamente!\n")
		fmt.Printf("   - Componentes detectados: %d\n", len(ast.Components))
		if len(ast.Components) > 0 {
			fmt.Printf("   - Primer componente: %s -> %s\n", 
				ast.Components[0].Name, ast.Components[0].Path)
		}
	}

	// Test 2: Validación de archivo
	fmt.Println("\n🧪 TEST 2: Validación")
	err = ast.Validate()
	if err != nil {
		fmt.Printf("❌ Validación falló: %v\n", err)
	} else {
		fmt.Printf("✅ Validación exitosa!\n")
	}

	// Test 3: Obtener dependencias
	fmt.Println("\n🧪 TEST 3: Dependencias de componentes")
	for _, comp := range ast.Components {
		fmt.Printf("   - %s (usado %d veces)\n", comp.Name, comp.UsageCount)
		for prop, value := range comp.Props {
			fmt.Printf("     * %s = %s\n", prop, value)
		}
	}

	// Test 4: Detector de componentes independiente
	fmt.Println("\n🧪 TEST 4: Detector independiente")
	detector := parser.NewComponentDetector()
	
	htmlContent := `<div>
		<user-profile name="Jane" avatar="/jane.jpg" />
		<shared-modal title="Confirm" closable />
		<data-table columns="name,email" sortable />
	</div>`
	
	components, err := detector.DetectComponents(htmlContent)
	if err != nil {
		log.Printf("Error detectando componentes: %v", err)
	} else {
		fmt.Printf("✅ Detectados %d componentes:\n", len(components))
		for _, comp := range components {
			fmt.Printf("   - %s\n", comp.Name)
		}
	}

	// Test 5: Variables de template
	fmt.Println("\n🧪 TEST 5: Variables de template")
	templateContent := `<div>
		<h1>{{.User.Name}}</h1>
		<p>{{.User.Bio}}</p>
		<span>{{.Count}} items</span>
	</div>`
	
	variables := detector.DetectGoTemplateVariables(templateContent)
	fmt.Printf("✅ Variables detectadas: %v\n", variables)

	// Test 6: Atributos HTMX
	fmt.Println("\n🧪 TEST 6: Atributos HTMX")
	htmxContent := `<div>
		<button hx-post="/api/save" hx-target="#result">Save</button>
		<form hx-boost="true" hx-indicator="#spinner">
			<input hx-validate-on="blur" />
		</form>
	</div>`
	
	htmxAttrs := detector.DetectHTMXAttributes(htmxContent)
	fmt.Printf("✅ Atributos HTMX detectados: %v\n", htmxAttrs)
}

func main() {
	fmt.Println("🚀 PRUEBAS DE API DEL PARSER GOX")
	fmt.Println("================================\n")
	
	testParserAPI()
	
	fmt.Println("\n✨ Todas las pruebas completadas!")
}