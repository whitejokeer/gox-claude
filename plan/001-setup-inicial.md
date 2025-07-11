# Task 001: Setup Inicial del Proyecto

## Descripción
Configurar la estructura base del proyecto GOX Framework, incluyendo la configuración de Go modules, estructura de directorios, y archivos de configuración básicos.

## Prioridad
Alta

## Estimación
2-3 días

## Dependencias
Ninguna

## Subtasks

### 1.1 Crear estructura de directorios base
- [ ] Crear directorio `cmd/gox` para el CLI principal
- [ ] Crear directorio `pkg/` para los paquetes públicos
- [ ] Crear directorio `internal/` para código interno
- [ ] Crear directorio `examples/` para proyectos de ejemplo
- [ ] Crear directorio `scripts/` para scripts de desarrollo
- [ ] Crear directorio `docs/` para documentación

### 1.2 Inicializar Go module
- [ ] Ejecutar `go mod init github.com/gox-framework/gox`
- [ ] Configurar Go 1.21 como versión mínima
- [ ] Agregar archivo `.gitignore` apropiado
- [ ] Configurar `.editorconfig` para consistencia de código

### 1.3 Configurar herramientas de desarrollo
- [ ] Crear `Makefile` con comandos básicos
- [ ] Configurar `golangci-lint` con reglas personalizadas
- [ ] Agregar configuración de VS Code en `.vscode/`
- [ ] Configurar pre-commit hooks

### 1.4 Crear archivos base
- [ ] Crear `README.md` principal
- [ ] Crear `LICENSE` (MIT)
- [ ] Crear `CONTRIBUTING.md`
- [ ] Crear `CHANGELOG.md`
- [ ] Crear `.github/` con templates de issues y PRs

### 1.5 Configurar CI/CD inicial
- [ ] Crear workflow de GitHub Actions para tests
- [ ] Crear workflow para linting
- [ ] Crear workflow para build
- [ ] Configurar dependabot

## Criterios de Aceptación

1. **Estructura de directorios**
   - La estructura debe seguir las convenciones de Go
   - Todos los directorios deben tener un README.md explicativo
   - La separación entre código público e interno debe ser clara

2. **Go Module**
   - El módulo debe estar correctamente inicializado
   - Las dependencias iniciales deben estar documentadas
   - El go.mod debe especificar Go 1.21+

3. **Herramientas de desarrollo**
   - `make help` debe mostrar todos los comandos disponibles
   - `make lint` debe ejecutar golangci-lint sin errores
   - `make test` debe estar configurado aunque no haya tests

4. **Documentación**
   - README debe incluir badges de CI/CD
   - CONTRIBUTING debe explicar el proceso de contribución
   - LICENSE debe ser MIT

5. **CI/CD**
   - Los workflows deben ejecutarse en cada push
   - Los checks deben pasar en verde
   - Debe haber protección de branch en main

## Tests Necesarios

### Tests Manuales
1. Clonar el repositorio en una máquina limpia
2. Ejecutar `make help` y verificar que funciona
3. Ejecutar `make lint` sin errores
4. Verificar que la estructura de directorios es correcta

### Tests Automatizados
1. **Test de estructura**
   ```bash
   #!/bin/bash
   # Verificar que existen los directorios necesarios
   dirs=("cmd/gox" "pkg" "internal" "examples" "scripts" "docs")
   for dir in "${dirs[@]}"; do
     if [ ! -d "$dir" ]; then
       echo "Directory $dir does not exist"
       exit 1
     fi
   done
   ```

2. **Test de Go module**
   ```bash
   # Verificar versión de Go
   go version
   # Verificar que el módulo está bien formado
   go mod verify
   ```

3. **Test de herramientas**
   ```bash
   # Verificar que las herramientas están instaladas
   which golangci-lint
   make lint
   ```

## Definición de Done

- [ ] Todos los subtasks completados
- [ ] Tests manuales pasando
- [ ] Tests automatizados pasando
- [ ] PR revisado y aprobado
- [ ] Documentación actualizada
- [ ] CI/CD configurado y funcionando

## Notas Adicionales

- Usar conventional commits desde el inicio
- Configurar semantic versioning
- Considerar usar goreleaser para futuros releases
- Mantener el CHANGELOG actualizado desde el primer commit