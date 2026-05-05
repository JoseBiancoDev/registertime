#!/bin/bash

# Colores para los mensajes
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}======================================================${NC}"
echo -e "${BLUE}  Instalador del Sistema de Control de Horas (PWA)    ${NC}"
echo -e "${BLUE}======================================================${NC}\n"

# Función para comprobar si un comando existe
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# 1. Comprobar requisitos previos
echo -e "${GREEN}[1/3] Comprobando requisitos previos...${NC}"

if ! command_exists docker; then
    echo -e "${RED}Error: Docker no está instalado. Por favor, instala Docker e inténtalo de nuevo.${NC}"
    exit 1
fi

if ! command_exists docker-compose; then
    # Comprobar también el plugin 'docker compose' (V2)
    if ! docker compose version >/dev/null 2>&1; then
        echo -e "${RED}Error: Docker Compose no está instalado. Por favor, instala Docker Compose e inténtalo de nuevo.${NC}"
        exit 1
    fi
    DOCKER_COMPOSE_CMD="docker compose"
else
    DOCKER_COMPOSE_CMD="docker-compose"
fi

echo "Docker y Docker Compose están instalados."

# 2. Configuración del entorno
echo -e "\n${GREEN}[2/3] Configurando el entorno...${NC}"

if [ ! -f .env ]; then
    echo -e "${BLUE}No se encontró un archivo .env. Creando uno básico...${NC}"
    cat <<EOF > .env
DB_USER=admin
DB_PASSWORD=admin123
DB_NAME=control_horario
JWT_SECRET=supersecret_cambiar_en_produccion
SMTP_HOST=smtp.ejemplo.com
SMTP_PORT=587
SMTP_USER=tu_correo@ejemplo.com
SMTP_PASS=tu_contraseña
EOF
    echo "Archivo .env creado. ¡Asegúrate de editarlo con tus credenciales reales (especialmente SMTP) más adelante!"
else
    echo "El archivo .env ya existe."
fi

# 3. Levantar los contenedores
echo -e "\n${GREEN}[3/3] Construyendo e iniciando los contenedores...${NC}"
echo "Esto puede tardar unos minutos la primera vez..."

$DOCKER_COMPOSE_CMD up -d --build

if [ $? -eq 0 ]; then
    echo -e "\n${BLUE}======================================================${NC}"
    echo -e "${GREEN}¡Instalación completada con éxito!${NC}"
    echo -e "${BLUE}======================================================${NC}"
    echo -e "La aplicación está ahora en funcionamiento:"
    echo -e "- ${GREEN}Frontend (PWA):${NC} http://localhost:3000"
    echo -e "- ${GREEN}Backend (API):${NC}  http://localhost:8080"
    echo -e "- ${GREEN}Base de Datos:${NC}  localhost:5432"
    echo -e "\nPara ver los logs en tiempo real, ejecuta: ${BLUE}$DOCKER_COMPOSE_CMD logs -f${NC}"
    echo -e "Para detener el sistema, ejecuta: ${BLUE}$DOCKER_COMPOSE_CMD down${NC}"
else
    echo -e "\n${RED}Hubo un error al intentar levantar los contenedores. Revisa la salida de error arriba.${NC}"
    exit 1
fi
