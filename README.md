PredictOne
API de Fórmula 1 y predicción de resultados de carreras

Descripción del proyecto

PredictOne es un proyecto que consiste en el desarrollo de una página web para gestionar información relacionada con la Fórmula 1 y realizar predicciones acerca de los resultados de las carreras.

La aplicación permite administrar escuderías, pilotos históricos, Grandes Premios y sus ediciones históricas, incluyendo las fechas y los resultados de las carreras.


Tecnologías utilizadas
Tecnología	Descripción
Go	Desarrollo del servidor HTTP y la lógica de la aplicación.
PostgreSQL	Base de datos relacional para almacenar la información.
sqlc	Generación de código Go a partir de consultas SQL.
Hurl	Ejecución de pruebas automatizadas de los endpoints HTTP.
Docker	Ejecución y administración de los servicios del proyecto.
Docker Compose	Configuración y administración del entorno de base de datos.


Cómo ejecutar el proyecto
1. Clonar el repositorio
git clone https://github.com/gianfrancocanale/Proyecto_web.git
cd Proyecto_web

2. Cambiar a la rama del trabajo práctico
git checkout tp3

3. Ejecutar los tests

Para ejecutar las pruebas automatizadas y realizar los pasos configurados en el Makefile:

make run

Y en otra terminal

make hurl



Modelo de datos

El sistema administra información de la Fórmula 1 mediante las siguientes entidades:

Entidad	Descripción
Escudería	Información de los equipos, puntos de temporada, títulos de constructores y jefe de equipo.
Piloto histórico	Datos de los pilotos, país, fecha de nacimiento y títulos ganados.
Gran Premio	Información del Gran Premio, país, longitud del circuito, cantidad de vueltas y último ganador.
Gran Premio histórico	Ediciones anteriores de los Grandes Premios, fecha de carrera y resultados obtenidos.

Las relaciones entre las entidades se implementan mediante claves primarias y foráneas en PostgreSQL.

API HTTP

La API permite realizar operaciones sobre las distintas entidades del sistema mediante solicitudes HTTP.

Escuderías
Método	Ruta	Operación
POST	/escuderias	Crear una escudería.
GET	/escuderias	Listar todas las escuderías.
GET	/escuderias/{id}	Obtener una escudería por su identificador.
PUT	/escuderias/{id}	Actualizar una escudería.
DELETE	/escuderias/{id}	Eliminar una escudería.
Pilotos históricos
Método	Ruta	Operación
POST	/pilotos_historicos	Crear un piloto histórico.
Grandes Premios
Método	Ruta	Operación
POST	/gran_premios	Crear un Gran Premio.
GET	/gran_premios	Listar todos los Grandes Premios.
Grandes Premios históricos
Método	Ruta	Operación
POST	/gran_premios_historicos	Crear una edición histórica.
GET	/gran_premios_historicos	Listar las ediciones históricas.
GET	/gran_premios_historicos/{id}/{fecha}	Obtener una edición por Gran Premio y fecha.
PUT	/gran_premios_historicos/{id}/{fecha}	Actualizar el resultado de una carrera.
DELETE	/gran_premios_historicos/{id}/{fecha}	Eliminar una edición histórica.
Nota: las tablas documentan las rutas principales identificadas. Para reflejar la API completa, deben agregarse las operaciones restantes que estén implementadas en el código.

Pruebas automatizadas

El proyecto utiliza Hurl para verificar el comportamiento de la API mediante solicitudes HTTP.

Las pruebas se encuentran en el directorio requests/ y permiten comprobar:

Creación de registros.
Listado y consulta por identificador.
Actualización de información.
Eliminación de registros.
Códigos de estado HTTP.
Datos devueltos en las respuestas JSON.
Relaciones entre entidades y claves foráneas.

Ejecutar las pruebas individualmente

Escuderías

hurl --test requests/escuderias.hurl


Grandes Premios

hurl --test requests/granpremio.hurl


Grandes Premios históricos

hurl --test requests/granpremiohistorico.hurl

Ejecutar las pruebas mediante Make
make hurl


La API y la base de datos deben estar disponibles antes de ejecutar las pruebas HTTP, salvo que el Makefile se encargue de iniciarlas automáticamente.

Estructura del proyecto
Proyecto_web/
├── db/
│   ├── queries/          # Consultas SQL
│   ├── schema/           # Esquema de la base de datos
│   └── sqlc/             # Código generado por sqlc
├── handlers/             # Controladores HTTP
├── requests/             # Pruebas automatizadas con Hurl
├── main.go               # Punto de entrada de la aplicación
├── sqlc.yaml             # Configuración de sqlc
├── docker-compose.yml    # Configuración de Docker Compose
├── Makefile              # Comandos automatizados
└── README.md              # Documentación del proyecto
