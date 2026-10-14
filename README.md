PredictOne

Descripción del proyecto
Este proyecto consiste en el desarrollo de una API HTTP para gestionar información relacionada con la Fórmula 1 y para realizar predicciones acerca de los resultados de las carreras. Permite administrar escuderías, pilotos históricos, Grandes Premios y sus ediciones históricas.

Tecnologías utilizadas
Go: lenguaje utilizado para desarrollar el servidor HTTP.
PostgreSQL: base de datos relacional para almacenar la información.
sqlc: generación de código Go a partir de consultas SQL.
Hurl: ejecución de pruebas automatizadas de los endpoints HTTP.
Docker y Docker Compose: administración del entorno de base de datos, si están configurados en el proyecto.

Requisitos

Antes de ejecutar el proyecto, es necesario contar con:

Go instalado.
PostgreSQL o Docker con Docker Compose.
sqlc, si no está automatizada su instalación.
Hurl para ejecutar las pruebas HTTP.

Cómo ejecutar el proyecto

Clonar el repositorio:

git clone https://github.com/gianfrancocanale/Proyecto_web.git
cd Proyecto_web


Cambiar al branch del trabajo práctico 3:

git checkout tp3


Para ejecutar los tests y realizar automáticamente todos los pasos necesarios:

```bash
make test
```


Modelo de datos

El sistema administra información de la Fórmula 1 mediante las siguientes entidades:

Entidad	Descripción
Escudería	Información de los equipos, puntos de temporada, títulos de constructores y jefe de equipo.
Piloto histórico	Datos de los pilotos, país, fecha de nacimiento y títulos ganados.
Gran Premio	Información del Gran Premio, país, longitud del circuito, cantidad de vueltas y último ganador.
Gran Premio histórico	Ediciones anteriores de los Grandes Premios, fecha de carrera y resultados.

Las relaciones entre entidades se implementan mediante claves primarias y foráneas en PostgreSQL.

API HTTP

La API permite realizar operaciones sobre las distintas entidades del sistema.

Método	Ruta	Operación
POST	/escuderias	Crear una escudería
GET	/escuderias	Listar escuderías
GET	/escuderias/{id}	Obtener una escudería
PUT	/escuderias/{id}	Actualizar una escudería
DELETE	/escuderias/{id}	Eliminar una escudería
POST	/pilotos_historicos	Crear un piloto histórico
POST	/gran_premios	Crear un Gran Premio
GET	/gran_premios	Listar Grandes Premios
POST	/gran_premios_historicos	Crear una edición histórica
GET	/gran_premios_historicos	Listar ediciones históricas
GET	/gran_premios_historicos/{id}/{fecha}	Obtener una edición por Gran Premio y fecha
PUT	/gran_premios_historicos/{id}/{fecha}	Actualizar el resultado de una carrera
DELETE	/gran_premios_historicos/{id}/{fecha}	Eliminar una edición histórica

Las rutas y operaciones restantes deben verificarse en el registro de rutas de la aplicación.


Pruebas automatizadas

El proyecto utiliza Hurl para verificar el comportamiento de la API mediante solicitudes HTTP.

Las pruebas se encuentran en el directorio requests/ y comprueban:

Creación de registros.
Listado y consulta por identificador.
Actualización de información.
Eliminación de registros.
Códigos de estado HTTP.
Datos devueltos en las respuestas JSON.
Relaciones entre entidades y claves foráneas.

Para ejecutar las pruebas individualmente:

hurl --test requests/escuderias.hurl

hurl --test requests/granpremio.hurl

hurl --test requests/granpremiohistorico.hurl


Para ejecutar las pruebas mediante el Makefile:

make hurl


La API y la base de datos deben estar disponibles antes de ejecutar las pruebas.

Estructura del proyecto
.
├── db/
│   ├── queries/       # Consultas SQL
│   ├── schema/        # Esquema de la base de datos
│   └── sqlc/          # Código generado por sqlc
├── handlers/          # Controladores HTTP
├── requests/          # Pruebas automatizadas con Hurl
├── main.go            # Punto de entrada de la aplicación
├── sqlc.yaml          # Configuración de sqlc
├── docker-compose.yml # Configuración de Docker
├── Makefile           # Comandos automatizados
└── README.md          # Documentación del proyecto


La estructura puede variar según la organización del repositorio.

Autores
Nombre y apellido del integrante 1
Nombre y apellido del integrante 2
