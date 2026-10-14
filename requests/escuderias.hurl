
# Crear escudería
POST http://localhost:8080/escuderias
Content-Type: application/json

{
  "id_escuderia": "ferrari",
  "puntos_temporada": 0,
  "titulos_constructores": 0,
  "fecha_fundacion": "1929-12-16",
  "team_principal": "Frédéric Vasseur"
}

HTTP/1.1 201
[Asserts]
jsonpath "$.id_escuderia" == "ferrari"
jsonpath "$.puntos_temporada" == 0
jsonpath "$.titulos_constructores" == 0
jsonpath "$.fecha_fundacion" exists
jsonpath "$.team_principal" == "Frédéric Vasseur"

[Captures]
id_escuderia: jsonpath "$.id_escuderia"


# Listar escuderías
GET http://localhost:8080/escuderias

HTTP/1.1 200
[Asserts]
jsonpath "$" isCollection
jsonpath "$[?(@.id_escuderia == 'ferrari')].team_principal" includes "Frédéric Vasseur"


# Obtener escudería
GET http://localhost:8080/escuderias/{{id_escuderia}}

HTTP/1.1 200
[Asserts]
jsonpath "$.id_escuderia" == "ferrari"
jsonpath "$.puntos_temporada" == 0
jsonpath "$.titulos_constructores" == 0
jsonpath "$.team_principal" == "Frédéric Vasseur"


# Actualizar escudería
PUT http://localhost:8080/escuderias/{{id_escuderia}}
Content-Type: application/json

{
  "puntos_temporada": 120,
  "titulos_constructores": 16,
  "fecha_fundacion": "1929-12-16",
  "team_principal": "Frédéric Vasseur"
}

HTTP/1.1 200
[Asserts]
jsonpath "$.id_escuderia" == "ferrari"
jsonpath "$.puntos_temporada" == 120
jsonpath "$.titulos_constructores" == 16
jsonpath "$.team_principal" == "Frédéric Vasseur"


# Verificar actualización
GET http://localhost:8080/escuderias/{{id_escuderia}}

HTTP/1.1 200
[Asserts]
jsonpath "$.puntos_temporada" == 120
jsonpath "$.titulos_constructores" == 16


# Eliminar escudería
DELETE http://localhost:8080/escuderias/{{id_escuderia}}

HTTP/1.1 204


# Verificar eliminación
GET http://localhost:8080/escuderias/{{id_escuderia}}

HTTP/1.1 404