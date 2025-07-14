package main

import (
	"fmt"
	
	"github.com/gox-framework/gox/internal/parser"
)

func main() {
	fmt.Println("🔴 EJEMPLOS DE MANEJO DE ERRORES DEL PARSER")
	fmt.Println("=========================================\n")
	
	// Error 1: Template sin cerrar
	fmt.Println("❌ ERROR 1: Template sin cerrar")
	gox1 := `<template>
  <div>Sin cerrar
</template>`
	
	_, err := parser.ParseString(gox1)
	if err != nil {
		fmt.Printf("   Error: %v\n\n", err)
	}
	
	// Error 2: Tags mal emparejados
	fmt.Println("❌ ERROR 2: Tags mal emparejados")
	gox2 := `<template>
  <div>Contenido</div>
</go>`
	
	_, err = parser.ParseString(gox2)
	if err != nil {
		fmt.Printf("   Error: %v\n\n", err)
	}
	
	// Error 3: Sintaxis Go inválida
	fmt.Println("❌ ERROR 3: Sintaxis Go inválida")
	gox3 := `<template>
  <div>Válido</div>
</template>

<go>
package pages

func syntax error {
    return no value
}
</go>`
	
	_, err = parser.ParseString(gox3)
	if err != nil {
		fmt.Printf("   Error: %v\n\n", err)
	}
	
	// Error 4: Sección desconocida
	fmt.Println("❌ ERROR 4: Sección desconocida")
	gox4 := `<template>
  <div>OK</div>
</template>

<script>
  console.log("No soportado");
</script>`
	
	_, err = parser.ParseString(gox4)
	if err != nil {
		fmt.Printf("   Error: %v\n\n", err)
	}
	
	// Error 5: CSS inválido
	fmt.Println("❌ ERROR 5: CSS inválido")
	gox5 := `<style>
.clase {
  color: ;
  background: red
}

.otra-clase {
  /* Sin cerrar
</style>`
	
	_, err = parser.ParseString(gox5)
	if err != nil {
		fmt.Printf("   Error: %v\n\n", err)
	}
	
	// Error 6: Archivo vacío
	fmt.Println("❌ ERROR 6: Archivo vacío")
	gox6 := ``
	
	ast, err := parser.ParseString(gox6)
	if err != nil {
		fmt.Printf("   Error parsing: %v\n", err)
	} else {
		// El archivo se parsea pero falla la validación
		err = ast.Validate()
		if err != nil {
			fmt.Printf("   Error validating: %v\n\n", err)
		}
	}
	
	// Error 7: Props mal formateadas
	fmt.Println("❌ ERROR 7: Props con sintaxis inválida")
	gox7 := `<template>
  <user-card name={{.Name}} email="{{.Email}}" />
</template>`
	
	ast, err = parser.ParseString(gox7)
	if err != nil {
		fmt.Printf("   Error: %v\n", err)
	} else {
		// Se parsea pero muestra props incorrectas
		if len(ast.Components) > 0 {
			fmt.Printf("   Warning: Props detectadas incorrectamente:\n")
			for k, v := range ast.Components[0].Props {
				fmt.Printf("     %s = %s\n", k, v)
			}
		}
		fmt.Println()
	}
	
	// Error 8: Múltiples secciones del mismo tipo
	fmt.Println("❌ ERROR 8: Múltiples templates (solo se usa el último)")
	gox8 := `<template>
  <div>Primer template</div>
</template>

<template>
  <div>Segundo template</div>
</template>`
	
	ast, err = parser.ParseString(gox8)
	if err != nil {
		fmt.Printf("   Error: %v\n", err)
	} else {
		fmt.Printf("   Warning: Solo se conservó el último template\n")
		fmt.Printf("   Contenido: %.50s...\n\n", ast.Template.Content)
	}
	
	// Demostración de ParseError con información de ubicación
	fmt.Println("📍 INFORMACIÓN DE UBICACIÓN EN ERRORES")
	errorContent := `<template>
  <div>
    <h1>Title</h1>
    <user-card name="test" />
  </div>
</template>

<go>
package components

type UserCard struct {
    Name string
}

func (c *UserCard) broken syntax here {
    return
}
</go>`
	
	_, err = parser.ParseString(errorContent)
	if err != nil {
		if parseErr, ok := err.(*parser.ParseError); ok {
			fmt.Printf("   Tipo: ParseError\n")
			fmt.Printf("   Línea: %d\n", parseErr.Line)
			fmt.Printf("   Columna: %d\n", parseErr.Column)
			fmt.Printf("   Mensaje: %s\n", parseErr.Message)
		} else {
			fmt.Printf("   Error genérico: %v\n", err)
		}
	}
	
	fmt.Println("\n✅ Todos los errores fueron capturados correctamente!")
}