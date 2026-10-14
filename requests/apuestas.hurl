
# Crear usuario necesario para la apuesta
POST http://localhost:8080/usuarios
Content-Type: application/json

{
  "nombre_usuario": "gian_apuesta_test",
  "contrasena_hash": "1234",
  "puntos_ganados": 0
}

HTTP 201


# Crear piloto histórico necesario para el circuito
POST http://localhost:8080/pilotos_historicos
Content-Type: application/json

{
  "id_piloto": 99,
  "nombre": "Piloto Prueba",
  "pais": "Argentina",
  "fecha_nacimiento": "2000-01-01",
  "titulos_ganados": 0
}

HTTP 201


# Crear circuito necesario para la carrera
POST http://localhost:8080/gran_premios
Content-Type: application/json

{
  "id_gran_premio": "GP_TEST_APUESTA",
  "pais": "Argentina",
  "longitud_km": 5,
  "cantidad_vueltas": 50,
  "id_ultimo_ganador": 99
}

HTTP 201


# Crear edición de carrera
POST http://localhost:8080/gran_premios_historicos
Content-Type: application/json

{
  "id_gran_premio": "GP_TEST_APUESTA",
  "fecha_carrera": "2026-11-22",
  "resultado_carrera": [99, 1, 2]
}

HTTP 201


# Crear apuesta
POST http://localhost:8080/apuestas
Content-Type: application/json

{
  "nombre_usuario": "gian_apuesta_test",
  "id_gran_premio": "GP_TEST_APUESTA",
  "fecha_carrera": "2026-11-22",
  "prediccion": [99, 1, 2, 3, 4, 5, 6, 7, 8, 9],
  "fecha_apuesta": "2026-10-14"
}

HTTP 201
[Asserts]
jsonpath "$.id_apuesta" exists
jsonpath "$.nombre_usuario" == "gian_apuesta_test"
jsonpath "$.id_gran_premio" == "GP_TEST_APUESTA"
jsonpath "$.prediccion[0]" == 99

[Captures]
id_apuesta: jsonpath "$.id_apuesta"


# Listar todas las apuestas
GET http://localhost:8080/apuestas

HTTP 200
[Asserts]
jsonpath "$" isCollection
jsonpath "$[?(@.id_apuesta == {{id_apuesta}})].nombre_usuario" includes "gian_apuesta_test"


# Listar apuestas de un usuario
GET http://localhost:8080/apuestas/gian_apuesta_test

HTTP 200
[Asserts]
jsonpath "$" isCollection
jsonpath "$[?(@.id_apuesta == {{id_apuesta}})].id_gran_premio" includes "GP_TEST_APUESTA"


# Obtener apuesta por ID y usuario
GET http://localhost:8080/apuestas/{{id_apuesta}}/gian_apuesta_test

HTTP 200
[Asserts]
jsonpath "$.id_apuesta" == {{id_apuesta}}
jsonpath "$.nombre_usuario" == "gian_apuesta_test"
jsonpath "$.id_gran_premio" == "GP_TEST_APUESTA"


# Modificar apuesta
PUT http://localhost:8080/apuestas/{{id_apuesta}}/gian_apuesta_test
Content-Type: application/json

{
  "prediccion": [1, 99, 2, 3, 4, 5, 6, 7, 8, 9],
  "fecha_apuesta": "2026-10-14",
  "id_gran_premio": "GP_TEST_APUESTA",
  "fecha_carrera": "2026-11-22"
}

HTTP 200
[Asserts]
jsonpath "$.id_apuesta" == {{id_apuesta}}
jsonpath "$.prediccion[0]" == 1
jsonpath "$.prediccion[1]" == 99


# Verificar modificación
GET http://localhost:8080/apuestas/{{id_apuesta}}/gian_apuesta_test

HTTP 200
[Asserts]
jsonpath "$.prediccion[0]" == 1
jsonpath "$.prediccion[1]" == 99


# Eliminar apuesta
DELETE http://localhost:8080/apuestas/{{id_apuesta}}/gian_apuesta_test

HTTP 204


# Verificar eliminación
GET http://localhost:8080/apuestas/{{id_apuesta}}/gian_apuesta_test

HTTP 404


# Limpiar edición de carrera
DELETE http://localhost:8080/gran_premios_historicos/GP_TEST_APUESTA/2026-11-22

HTTP 204


# Limpiar circuito
DELETE http://localhost:8080/gran_premios/GP_TEST_APUESTA

HTTP 204


# Limpiar piloto histórico
DELETE http://localhost:8080/pilotos_historicos/99

HTTP 204


# Limpiar usuario
DELETE http://localhost:8080/usuarios/gian_apuesta_test

HTTP 204