
# Crear piloto histórico necesario para la clave foránea
POST http://localhost:8080/pilotos_historicos
Content-Type: application/json

{
  "id_piloto": 44,
  "nombre": "Max Verstappen",
  "pais": "Países Bajos",
  "fecha_nacimiento": "1997-09-30T00:00:00Z",
  "titulos_ganados": 0
}


HTTP/1.1 201


# Crear escudería necesaria para la clave foránea
POST http://localhost:8080/escuderias
Content-Type: application/json

{
  "id_escuderia": "red_bull",
  "puntos_temporada": 0,
  "titulos_constructores": 0,
  "fecha_fundacion": "2005-11-15T00:00:00Z",
  "team_principal": "Christian Horner"
}

HTTP/1.1 201


# Crear piloto de la temporada
POST http://localhost:8080/pilotos
Content-Type: application/json

{
  "id_piloto": 44,
  "id_escuderia": "red_bull",
  "puntos_temporada": 0
}

HTTP/1.1 201
[Asserts]
jsonpath "$.id_piloto" == 44
jsonpath "$.id_escuderia" == "red_bull"
jsonpath "$.puntos_temporada" == 0

[Captures]
id_piloto: jsonpath "$.id_piloto"


# Listar pilotos de la temporada
GET http://localhost:8080/pilotos

HTTP/1.1 200
[Asserts]
jsonpath "$" isCollection
jsonpath "$[?(@.id_piloto == {{id_piloto}})].id_escuderia" includes "red_bull"


# Obtener piloto de la temporada
GET http://localhost:8080/pilotos/{{id_piloto}}

HTTP/1.1 200
[Asserts]
jsonpath "$.id_piloto" == {{id_piloto}}
jsonpath "$.id_escuderia" == "red_bull"
jsonpath "$.puntos_temporada" == 0


# Actualizar piloto de la temporada
PUT http://localhost:8080/pilotos/{{id_piloto}}
Content-Type: application/json

{
  "id_escuderia": "red_bull",
  "puntos_temporada": 250
}

HTTP/1.1 200
[Asserts]
jsonpath "$.id_piloto" == {{id_piloto}}
jsonpath "$.id_escuderia" == "red_bull"
jsonpath "$.puntos_temporada" == 250


# Verificar actualización
GET http://localhost:8080/pilotos/{{id_piloto}}

HTTP/1.1 200
[Asserts]
jsonpath "$.puntos_temporada" == 250


# Eliminar piloto de la temporada
DELETE http://localhost:8080/pilotos/{{id_piloto}}

HTTP/1.1 200


# Verificar eliminación del piloto de temporada
GET http://localhost:8080/pilotos/{{id_piloto}}

HTTP/1.1 404


# Limpiar escudería de prueba
DELETE http://localhost:8080/escuderias/red_bull

HTTP/1.1 200


# Limpiar piloto histórico de prueba
DELETE http://localhost:8080/pilotos_historicos/{{id_piloto}}

HTTP/1.1 200