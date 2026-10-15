# PredictOne: API de Fórmula 1 y Predicción de Resultados

## Descripción del Proyecto

**PredictOne** es un proyecto enfocado en el desarrollo de una aplicación web para gestionar información sobre la Fórmula 1 y realizar predicciones acerca de los resultados de las carreras.

La aplicación permite administrar:
* **Usuarios** registrados en el sistema.
* **Apuestas y predicciones** de carreras realizadas por los usuarios.
* **Escuderías** y sus estadísticas de temporada.
* **Pilotos** actuales e históricos.
* **Grandes Premios** y sus circuitos.
* **Ediciones históricas de Grandes Premios**, incluyendo fechas y resultados detallados.

---

## Tecnologías Utilizadas

| Tecnología | Descripción |
| :--- | :--- |
| **Go** | Desarrollo del servidor HTTP y lógica de negocio. |
| **PostgreSQL** | Base de datos relacional para el almacenamiento persistente. |
| **sqlc** | Generación de código Go seguro a partir de consultas SQL. |
| **Hurl** | Ejecución de pruebas automatizadas de los endpoints HTTP. |
| **Docker** | Containerización y administración de servicios. |
| **Docker Compose** | Configuración del entorno de base de datos y servicios. |

---

## Cómo Ejecutar el Proyecto

### 1. Clonar el repositorio
```bash
git clone https://github.com/gianfrancocanale/Proyecto_web.git
cd Proyecto_web
```

### 2. Cambiar al branch del tp3
```bash
git checkout tp3
```

### 3. Ejecutar la aplicación y las pruebas

En la terminal principal, inicia el servidor y los servicios requeridos:
```bash
make run
```

En otra terminal, ejecuta la parte de pruebas:
```bash
make hurl
```

---

## Modelo de Datos

El sistema administra la información mediante las siguientes entidades principales:

| Entidad | Descripción |
| :--- | :--- |
| **Usuario** | Gestión de los usuarios registrados en el sistema. |
| **Apuesta** | Predicciones realizadas por los usuarios para cada Gran Premio. |
| **Escudería** | Información de los equipos, puntos de temporada, títulos de constructores y jefe de equipo. |
| **Piloto** | Información de los pilotos actuales y activos en la competición. |
| **Piloto histórico** | Datos de los pilotos históricos, país, fecha de nacimiento y títulos ganados. |
| **Gran Premio** | Información del Gran Premio, país, longitud del circuito, cantidad de vueltas y último ganador. |
| **Gran Premio histórico** | Ediciones anteriores de los Grandes Premios, fecha de carrera y resultados obtenidos. |
---

## API 
### Usuarios
| Método | Ruta | Operación |
| :---: | :--- | :--- |
| `POST` | `/usuarios` | Crear un usuario |
| `GET` | `/usuarios` | Listar todos los usuarios |
| `GET` | `/usuarios/{usuario}` | Obtener un usuario por su nombre de usuario |
| `PUT` | `/usuarios/{usuario}` | Actualizar un usuario |
| `DELETE` | `/usuarios/{usuario}` | Eliminar un usuario |

### Apuestas
| Método | Ruta | Operación |
| :---: | :--- | :--- |
| `POST` | `/apuestas` | Crear una apuesta |
| `GET` | `/apuestas` | Listar todas las apuestas |
| `GET` | `/apuestas/{usuario}` | Obtener todas las apuestas de un usuario |
| `GET` | `/apuestas/{usuario}/{id}` | Obtener una apuesta específica de un usuario |
| `PUT` | `/apuestas/{usuario}/{id}` | Actualizar una apuesta de un usuario |
| `DELETE` | `/apuestas/{usuario}/{id}` | Eliminar una apuesta |

### Escuderías
| Método | Ruta | Operación |
| :---: | :--- | :--- |
| `POST` | `/escuderias` | Crear una escudería |
| `GET` | `/escuderias` | Listar todas las escuderías |
| `GET` | `/escuderias/{id}` | Obtener una escudería por identificador |
| `PUT` | `/escuderias/{id}` | Actualizar una escudería |
| `DELETE` | `/escuderias/{id}` | Eliminar una escudería |

### Pilotos
| Método | Ruta | Operación |
| :---: | :--- | :--- |
| `POST` | `/pilotos` | Crear un piloto |
| `GET` | `/pilotos` | Listar todos los pilotos |
| `GET` | `/pilotos/{id}` | Obtener un piloto por identificador |
| `PUT` | `/pilotos/{id}` | Actualizar un piloto |
| `DELETE` | `/pilotos/{id}` | Eliminar un piloto |

### Pilotos Históricos
| Método | Ruta | Operación |
| :---: | :--- | :--- |
| `POST` | `/pilotos_historicos` | Crear un piloto histórico |
| `GET` | `/pilotos_historicos` | Listar todos los pilotos históricos |
| `GET` | `/pilotos_historicos/{id}` | Obtener un piloto histórico por identificador |
| `PUT` | `/pilotos_historicos/{id}` | Actualizar un piloto histórico |
| `DELETE` | `/pilotos_historicos/{id}` | Eliminar un piloto histórico |

### Grandes Premios
| Método | Ruta | Operación |
| :---: | :--- | :--- |
| `POST` | `/gran_premios` | Crear un Gran Premio |
| `GET` | `/gran_premios` | Listar todos los Grandes Premios |
| `GET` | `/gran_premios/{id}` | Obtener un Gran Premio por identificador |
| `PUT` | `/gran_premios/{id}` | Actualizar un Gran Premio |
| `DELETE` | `/gran_premios/{id}` | Eliminar un Gran Premio |

### Grandes Premios Históricos
| Método | Ruta | Operación |
| :---: | :--- | :--- |
| `POST` | `/gran_premios_historicos` | Crear una edición histórica |
| `GET` | `/gran_premios_historicos` | Listar las ediciones históricas |
| `GET` | `/gran_premios_historicos/{id}/{fecha}` | Obtener una edición por Gran Premio y fecha |
| `PUT` | `/gran_premios_historicos/{id}/{fecha}` | Actualizar el resultado de una carrera |
| `DELETE` | `/gran_premios_historicos/{id}/{fecha}` | Eliminar una edición histórica |
