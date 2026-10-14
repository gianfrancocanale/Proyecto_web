
# Crear piloto histórico
POST http://localhost:8080/pilotos_historicos
Content-Type: application/json

{
  "id_piloto": 44,
  "nombre": "Max Verstappen",
  "pais": "Países Bajos",
  "fecha_nacimiento": "1997-09-30",
  "titulos_ganados": 0
}

HTTP/1.1 201
[Asserts]
jsonpath "$.id_piloto" == 44
jsonpath "$.nombre" == "Max Verstappen"
jsonpath "$.pais" == "Países Bajos"
jsonpath "$.titulos_ganados" == 0

[Captures]
id_piloto: jsonpath "$.id_piloto"


# Listar pilotos históricos
GET http://localhost:8080/pilotos_historicos

HTTP/1.1 200
[Asserts]
jsonpath "$" isCollection
jsonpath "$[?(@.id_piloto == {{id_piloto}})].nombre" includes "Max Verstappen"


# Obtener piloto histórico
GET http://localhost:8080/pilotos_historicos/{{id_piloto}}

HTTP/1.1 200
[Asserts]
jsonpath "$.id_piloto" == {{id_piloto}}
jsonpath "$.nombre" == "Max Verstappen"
jsonpath "$.pais" == "Países Bajos"
jsonpath "$.titulos_ganados" == 0


# Actualizar piloto histórico
PUT http://localhost:8080/pilotos_historicos/{{id_piloto}}
Content-Type: application/json

{
  "nombre": "Max Verstappen",
  "pais": "Países Bajos",
  "fecha_nacimiento": "1997-09-30",
  "titulos_ganados": 4
}

HTTP/1.1 200
[Asserts]
jsonpath "$.id_piloto" == {{id_piloto}}
jsonpath "$.nombre" == "Max Verstappen"
jsonpath "$.pais" == "Países Bajos"
jsonpath "$.titulos_ganados" == 4


# Verificar actualización
GET http://localhost:8080/pilotos_historicos/{{id_piloto}}

HTTP/1.1 200
[Asserts]
jsonpath "$.titulos_ganados" == 4


# Eliminar piloto histórico
DELETE http://localhost:8080/pilotos_historicos/{{id_piloto}}

HTTP/1.1 200


# Verificar eliminación
GET http://localhost:8080/pilotos_historicos/{{id_piloto}}

HTTP/1.1 404