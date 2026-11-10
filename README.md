# PredictOne

**PredictOne** es una aplicación web para consultar información sobre la Fórmula 1 y realizar predicciones sobre los resultados de las carreras.

Permite:

* Registrarse, ingresar y salir del sitio.
* Consultar **escuderías**, **pilotos** y **Grandes Premios**.
* Cargar, modificar y eliminar **apuestas** (predicciones de los 10 primeros puestos de un Gran Premio), solo con la sesión iniciada.

---

## Tecnologías utilizadas

| Tecnología | Descripción |
| :--- | :--- |
| **Go** | Servidor HTTP (`net/http`) y lógica de negocio. |
| **templ** | Generación de las vistas HTML desde Go (paquete `views`). |
| **PostgreSQL** | Base de datos relacional para el almacenamiento persistente. |
| **sqlc** | Generación de código Go seguro a partir de consultas SQL. |
| **Docker / Docker Compose** | Containerización de la API y de la base de datos. |
| **Make** | Automatización de compilación y ejecución. |

---

## Cómo ejecutar el proyecto

### Requisitos previos

Tener instalados: `git`, `make`, `go`, `sqlc`, `templ`, `docker` (con Docker Compose).

### 1. Clonar el repositorio

```bash
git clone https://github.com/gianfrancocanale/Proyecto_web.git
cd Proyecto_web
```

### 2. Cambiar al branch del tp4

```bash
git checkout tp4
```

### 3. Configurar las variables de entorno

Las credenciales de la base de datos se leen desde un archivo `.env`. El `Makefile` lo genera a partir de `.env.example` si no existe. Variables utilizadas:

| Variable | Descripción |
| :--- | :--- |
| `DB_HOST` | Host de la base de datos (`database` dentro de Docker). |
| `DB_PORT` | Puerto de PostgreSQL. |
| `DB_USER` | Usuario de la base. |
| `DB_PASSWORD` | Contraseña de la base. |
| `DB_NAME` | Nombre de la base. |

### 4. Ejecutar la aplicación

```bash
make run
```

Esto genera el código de `sqlc` y `templ`, compila el proyecto y levanta PostgreSQL y la API en Docker. Cuando termina, el sitio queda disponible en:

**http://localhost:8080**

---

## Páginas y rutas

El sitio se renderiza del lado del servidor. Los archivos estáticos (CSS, imágenes) se sirven bajo `/front/`.

| Ruta | Acceso | Descripción |
| :--- | :---: | :--- |
| `/` | Público | Página de inicio. |
| `/escuderias` | Público | Listado de escuderías, con un botón para invertir el orden. |
| `/pilotos` | Público | Listado de pilotos. |
| `/gran_premios` | Público | Listado de Grandes Premios y circuitos. |
| `/ingresar` | Público | Formulario de inicio de sesión. |
| `/registrar` | Público | Formulario de registro de usuario. |
| `/salir` | Público | Cierra la sesión. |
| `/apuestas` | Requiere sesión | Ver las apuestas propias y cargar una nueva eligiendo un Gran Premio y los 10 pilotos. |
| `/apuestas/modificar` | Requiere sesión | Modificar una apuesta. |
| `/apuestas/eliminar` | Requiere sesión | Eliminar una apuesta. |

### Sesiones

Al ingresar, el servidor crea una cookie `sesion` con un token asociado al nombre de usuario (guardado en memoria). Un middleware (`ConSesion`) procesa la cookie en cada request y `RequiereLogin` protege las rutas de apuestas, redirigiendo al login si no hay sesión activa. El usuario de cada apuesta se obtiene siempre de la cookie.

> Como las sesiones viven en memoria, se pierden al reiniciar el servidor.

---

## Modelo de datos

| Entidad | Descripción |
| :--- | :--- |
| **Usuario** | Usuarios registrados. Clave primaria: `nombre_usuario`. Guarda los puntos ganados. |
| **Apuesta** | Predicción de un usuario para un Gran Premio y una fecha de carrera: 10 posiciones y fecha de apuesta. Única por usuario, Gran Premio y fecha. |
| **Escudería** | Equipos, puntos de temporada, títulos de constructores y jefe de equipo. |
| **Piloto** | Pilotos actuales (por temporada), asociados a una escudería y a su ficha histórica. |
| **Piloto histórico** | Datos de pilotos: país, fecha de nacimiento y títulos ganados. |
| **Gran Premio** | País, circuito, longitud, cantidad de vueltas y último ganador. |
| **Gran Premio histórico** | Ediciones pasadas: fecha de carrera y resultado (22 posiciones). |

---

## Estructura del repositorio

```
.
├── db/                  # Esquema SQL, queries y código generado por sqlc
├── handlers/            # Handlers HTTP y middleware de sesión
├── views/               # Vistas templ
├── front/               # Archivos estáticos (CSS, imágenes)
├── main.go              # Rutas y arranque del servidor
├── Makefile
├── Dockerfile
├── docker-compose.yml
├── sqlc.yaml
├── .env.example
└── README.md
```