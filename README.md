# RegisterTime - Sistema de Control de Horario (PWA Premium)

Este proyecto es un sistema integral y avanzado de control de horario y asistencia (Time Tracking), diseñado con una arquitectura robusta y moderna. Permite gestionar el registro de jornada laboral de los empleados, soportando múltiples roles, modalidades de trabajo y generación de reportes automáticos. Todo el ecosistema está completamente "dockerizado" para garantizar un despliegue rápido y consistente.

## 🚀 ¿Qué hace el proyecto?

El sistema está diseñado para facilitar el registro de tiempos de trabajo con las siguientes funcionalidades clave:

### 👥 Gestión de Usuarios y Roles
- **Autenticación Segura:** Sistema de login mediante tokens JWT.
- **Roles de Acceso:** Diferenciación entre usuarios regulares y administradores.
- **Recuperación de Contraseña:** Flujos completos de "olvidé mi contraseña", "restablecer" y "cambiar contraseña" vía email.
- **Administración:** Los administradores pueden crear nuevos usuarios, listar empleados y auditar la plataforma.

### ⏱️ Control de Asistencia (Fichaje)
- **Clock In / Clock Out:** Los usuarios pueden iniciar (fichar entrada) y detener (fichar salida) su jornada laboral.
- **Modalidades de Trabajo:** Selección de modo de trabajo (Remoto o Presencial/On-site) al iniciar la sesión.
- **Gestión Administrativa:** Los administradores pueden añadir registros manuales, detener jornadas en curso de otros usuarios y consultar el historial completo de fichajes de cualquier empleado.

### 📊 Reportes y Notificaciones
- **Exportación a Excel:** Generación automática de reportes mensuales y detallados en formato `.xlsx`.
- **Notificaciones por Correo:** Envíos de correos transaccionales (por ejemplo, para recuperación de contraseñas o finalización de jornadas).

### 📱 Experiencia de Usuario (PWA)
- **Progressive Web App:** Instalable en dispositivos móviles y de escritorio, comportándose como una aplicación nativa.
- **Interfaz Premium:** Diseño moderno, responsivo y estético con componentes de alta calidad.

---

## 🛠️ Herramientas y Tecnologías Utilizadas

El proyecto sigue una arquitectura Cliente-Servidor separada, utilizando las siguientes tecnologías:

### Backend (API REST)
Desarrollado en **Go (Golang 1.24)** por su alto rendimiento y concurrencia.
- **[Gin Framework](https://gin-gonic.com/):** Enrutador HTTP ultrarrápido y robusto.
- **[GORM](https://gorm.io/):** ORM para interactuar de forma segura y estructurada con la base de datos.
- **[JWT (JSON Web Tokens)](https://jwt.io/):** Para la seguridad, autenticación y autorización de rutas protegidas.
- **[Excelize](https://github.com/qax-os/excelize):** Librería para la generación de archivos y reportes en formato Microsoft Excel.
- **[Gomail](https://github.com/go-gomail/gomail):** Para la gestión y envío de notificaciones por correo electrónico SMTP.
- **PostgreSQL:** Motor de base de datos relacional primario (soportando también SQLite para entornos ligeros).

### Frontend (Cliente Web / PWA)
Construido con el framework **[Nuxt 3](https://nuxt.com/)** (basado en Vue.js 3).
- **[Vue 3](https://vuejs.org/):** Core de la interfaz reactiva (Composition API).
- **[PrimeVue](https://primevue.org/) + PrimeFlex + PrimeIcons:** Sistema de componentes UI, utilidades CSS e íconos para lograr una estética premium.
- **[Pinia](https://pinia.vuejs.org/):** Gestor de estado global de la aplicación.
- **[Vite PWA](https://vite-pwa-org.netlify.app/):** Para la configuración e instalación de la aplicación como PWA (Service Workers y Manifest).
- **[Axios](https://axios-http.com/):** Cliente HTTP para el consumo de la API REST del backend.

### Infraestructura y DevOps
- **[Docker](https://www.docker.com/) & Docker Compose:** Contenerización de la base de datos, el backend y el frontend. Orquestación mediante `docker-compose.yml` para levantar todo el ecosistema con un solo comando.
- **Vitest & Go Test:** Cobertura de pruebas unitarias tanto en frontend como en backend.

---

## 🚀 Inicio Rápido

1. Clona el repositorio.
2. Configura las variables de entorno en el archivo `.env` (credenciales de DB, JWT_SECRET, y configuraciones SMTP para correos).
3. Levanta el sistema con Docker:
   ```bash
   docker-compose up -d --build
   ```
4. Accede a la aplicación en tu navegador: **`http://localhost:3000`**

### Scripts y Comandos Útiles

- **Ver logs de los contenedores:**
  ```bash
  docker-compose logs -f
  ```
- **Tests del backend:**
  ```bash
  docker-compose exec control-horario-backend go test ./...
  ```
- **Tests del frontend:**
  ```bash
  docker-compose exec control-horario-frontend npm run test
  ```
