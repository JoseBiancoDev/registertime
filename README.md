# Control de Horas - PWA Premium

Este proyecto es un sistema de control de horario desarrollado con **Go (Gin)** en el backend y **Nuxt 3 (PrimeVue)** en el frontend. Es una PWA (Progressive Web App) totalmente dockerizada.

## Características
- 🚀 **Backend Go**: Rápido y eficiente con Gin y GORM.
- 🎨 **Frontend Premium**: Interfaz moderna y amigable con PrimeVue.
- 📱 **PWA**: Instalable en móvil y escritorio.
- 📧 **Notificaciones**: Avisos por correo al finalizar jornadas.
- 📊 **Reportes**: Generación de reportes detallados en Excel.
- 🐳 **Docker**: Ejecución aislada sin dependencias locales.
- ✅ **Tests & Seeders**: Cobertura de pruebas y generación de datos iniciales.

## Requisitos
- Docker y Docker Compose.

## Inicio Rápido

1. Clona el repositorio.
2. Configura las variables de entorno en el archivo `.env`.
3. Levanta el sistema:
   ```bash
   docker-compose up --build
   ```
4. Accede a `http://localhost:3000`.

## Usuarios y Seeders
El sistema crea automáticamente un administrador al iniciar:
- **Email**: `admin@admin.com`
- **Password**: `admin123`

Para generar **100 usuarios de prueba** adicionales con registros de horas, puedes ejecutar el seeder manualmente:
```bash
docker-compose exec backend ./main -seed
```

## Pruebas
Para ejecutar los tests del backend:
```bash
docker-compose exec backend go test ./...
```

Para los tests del frontend:
```bash
docker-compose exec frontend npm test
```
