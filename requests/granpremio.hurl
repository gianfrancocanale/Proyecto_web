
# Crear piloto histórico necesario para la clave foránea
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


# Crear gran premio
POST http://localhost:8080/gran_premios
Content-Type: application/json

{
  "id_gran_premio": "SILVERSTONE",
  "pais": "Reino Unido",
  "longitud_km": 5,
  "cantidad_vueltas": 52,
  "id_ultimo_ganador": 44
}

HTTP/1.1 201
[Asserts]
jsonpath "$.id_gran_premio" == "SILVERSTONE"
jsonpath "$.pais" == "Reino Unido"
jsonpath "$.longitud_km" == 5
jsonpath "$.cantidad_vueltas" == 52
jsonpath "$.id_ultimo_ganador" == 44

[Captures]
id_gran_premio: jsonpath "$.id_gran_premio"


# Listar grandes premios
GET http://localhost:8080/gran_premios

HTTP/1.1 200
[Asserts]
jsonpath "$" isCollection
jsonpath "$[?(@.id_gran_premio == 'SILVERSTONE')].pais" includes "Reino Unido"


# Obtener gran premio
GET http://localhost:8080/gran_premios/{{id_gran_premio}}

HTTP/1.1 200
[Asserts]
jsonpath "$.id_gran_premio" == "SILVERSTONE"
jsonpath "$.pais" == "Reino Unido"
jsonpath "$.longitud_km" == 5
jsonpath "$.cantidad_vueltas" == 52
jsonpath "$.id_ultimo_ganador" == 44


# Actualizar gran premio
PUT http://localhost:8080/gran_premios/{{id_gran_premio}}
Content-Type: application/json

{
  "pais": "Reino Unido",
  "longitud_km": 6,
  "cantidad_vueltas": 52,
  "id_ultimo_ganador": 44
}

HTTP/1.1 200
[Asserts]
jsonpath "$.id_gran_premio" == "SILVERSTONE"
jsonpath "$.pais" == "Reino Unido"
jsonpath "$.longitud_km" == 6
jsonpath "$.cantidad_vueltas" == 52
jsonpath "$.id_ultimo_ganador" == 44


# Verificar actualización
GET http://localhost:8080/gran_premios/{{id_gran_premio}}

HTTP/1.1 200
[Asserts]
jsonpath "$.longitud_km" == 6


# Eliminar gran premio
DELETE http://localhost:8080/gran_premios/{{id_gran_premio}}

HTTP/1.1 204


# Verificar eliminación
GET http://localhost:8080/gran_premios/{{id_gran_premio}}

HTTP/1.1 404


# Limpiar piloto histórico
DELETE http://localhost:8080/pilotos_historicos/44

HTTP/1.1 204