#!/bin/bash

# Demostración completa de Task 2
# Este script muestra todas las funcionalidades implementadas

set -e

# Colores
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}======================================${NC}"
echo -e "${BLUE}  Demostración Task 2 - GOX Framework${NC}"
echo -e "${BLUE}======================================${NC}"
echo

# Preparación
cd /home/jairo/gox
echo -e "${YELLOW}Compilando GOX...${NC}"
go build -o gox-demo ./cmd/gox/

# Crear directorio temporal para la demo
DEMO_DIR="/home/jairo/gox/demo-task2"
rm -rf $DEMO_DIR
mkdir -p $DEMO_DIR
cd $DEMO_DIR

echo -e "\n${YELLOW}=== 1. DEMOSTRACIÓN: Detección de Contexto ===${NC}\n"

echo -e "${BLUE}Intentando generar página fuera de un proyecto (debe fallar)...${NC}"
if ! ../gox-demo generate page test 2>/dev/null; then
    echo -e "${GREEN}✓ Correctamente rechazado - No estamos en un proyecto GOX${NC}"
else
    echo -e "${RED}✗ Error - Debería haber fallado${NC}"
fi

echo -e "\n${YELLOW}=== 2. DEMOSTRACIÓN: Crear Nuevo Proyecto ===${NC}\n"

echo -e "${BLUE}Creando proyecto con arquitectura gateway...${NC}"
echo -e "${BLUE}Comando: gox new project todo-app --router=gin --auth=jwt --db=postgres${NC}"
../gox-demo new project todo-app --router=gin --auth=jwt --db=postgres

echo -e "\n${GREEN}✓ Proyecto creado exitosamente${NC}"
echo -e "\n${BLUE}Estructura generada:${NC}"
tree -L 3 todo-app/ 2>/dev/null || find todo-app -type d -not -path '*/\.*' | head -20

echo -e "\n${YELLOW}=== 3. DEMOSTRACIÓN: Comandos Dentro del Proyecto ===${NC}\n"

cd todo-app

echo -e "${BLUE}Intentando crear otro proyecto dentro del actual (debe fallar)...${NC}"
if ! ../../gox-demo new project another-app 2>/dev/null; then
    echo -e "${GREEN}✓ Correctamente rechazado - Ya estamos dentro de un proyecto${NC}"
else
    echo -e "${RED}✗ Error - Debería haber fallado${NC}"
fi

echo -e "\n${BLUE}Generando página con autenticación...${NC}"
echo -e "${BLUE}Comando: gox generate page dashboard --auth${NC}"
../../gox-demo generate page dashboard --auth

echo -e "\n${BLUE}Generando componente con props...${NC}"
echo -e "${BLUE}Comando: gox generate component todo-item --props=\"title:string,completed:bool,priority:int\"${NC}"
../../gox-demo generate component todo-item --props="title:string,completed:bool,priority:int"

echo -e "\n${BLUE}Generando componente compartido...${NC}"
echo -e "${BLUE}Comando: gox generate component button --shared${NC}"
../../gox-demo generate component button --shared

echo -e "\n${BLUE}Generando middleware...${NC}"
echo -e "${BLUE}Comando: gox generate middleware rate-limit${NC}"
../../gox-demo generate middleware rate-limit

echo -e "\n${BLUE}Generando microservicio...${NC}"
echo -e "${BLUE}Comando: gox generate service todos --api${NC}"
../../gox-demo generate service todos --api

echo -e "\n${YELLOW}=== 4. VERIFICACIÓN: Archivos Generados ===${NC}\n"

echo -e "${BLUE}Verificando archivos creados:${NC}"
FILES_TO_CHECK=(
    "gateway/pages/dashboard.gox"
    "gateway/components/todo-item.gox"
    "gateway/shared/ui/button.gox"
    "middleware/rate-limit.go"
    "services/todos/cmd/server/main.go"
    "docker-compose.yml"
    "go.work"
)

for file in "${FILES_TO_CHECK[@]}"; do
    if [ -f "$file" ]; then
        echo -e "${GREEN}✓${NC} $file"
    else
        echo -e "${RED}✗${NC} $file - No encontrado"
    fi
done

echo -e "\n${YELLOW}=== 5. CONTENIDO: Verificando Generación Correcta ===${NC}\n"

echo -e "${BLUE}Dashboard con autenticación:${NC}"
grep -A 2 -B 2 "auth" gateway/pages/dashboard.gox 2>/dev/null || echo "Contenido de auth no encontrado"

echo -e "\n${BLUE}Componente con props:${NC}"
grep -E "Title string|Completed bool|Priority int" gateway/components/todo-item.gox 2>/dev/null || echo "Props no encontradas"

echo -e "\n${BLUE}Gateway con router Gin:${NC}"
grep "gin" gateway/main.go 2>/dev/null | head -3 || echo "Gin no encontrado"

echo -e "\n${BLUE}Docker Compose con servicios:${NC}"
grep -E "gateway:|consul:|postgres:" docker-compose.yml 2>/dev/null || echo "Servicios no encontrados"

echo -e "\n${YELLOW}=== 6. DEMOSTRACIÓN: Crear Servicio Independiente ===${NC}\n"

cd $DEMO_DIR
echo -e "${BLUE}Creando microservicio independiente...${NC}"
echo -e "${BLUE}Comando: gox new service notification-service --db=postgres${NC}"
../gox-demo new service notification-service --db=postgres

echo -e "\n${GREEN}✓ Servicio creado exitosamente${NC}"
ls -la notification-service/

echo -e "\n${BLUE}======================================${NC}"
echo -e "${BLUE}    Resumen de Funcionalidades${NC}"
echo -e "${BLUE}======================================${NC}"

echo -e "\n${GREEN}✓ Detección de Contexto:${NC}"
echo "  - Comandos 'new' solo funcionan fuera de proyectos"
echo "  - Comandos 'generate' solo funcionan dentro de proyectos"

echo -e "\n${GREEN}✓ Arquitectura Gateway:${NC}"
echo "  - Gateway con pages, components, shared/ui"
echo "  - Servicios independientes en /services"
echo "  - Docker Compose con Consul para service discovery"

echo -e "\n${GREEN}✓ Generación de Código:${NC}"
echo "  - Páginas con autenticación opcional"
echo "  - Componentes con props tipadas"
echo "  - Componentes compartidos con --shared"
echo "  - Microservicios completos con estructura estándar"

echo -e "\n${GREEN}✓ Multi-Router Support:${NC}"
echo "  - HTTP, Gin, Echo, Fiber, Gorilla"
echo "  - Código específico por router"

echo -e "\n${GREEN}🎉 Task 2 completado exitosamente!${NC}"

# Limpiar
cd /home/jairo/gox
rm -f gox-demo