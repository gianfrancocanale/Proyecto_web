# Crear usuario necesario para la apuesta
POST http://localhost:8080/usuarios
Content-Type: application/json

{
  "nombre_usuario": "gian_apuesta_test",
  "contrasena_hash": "1234",
  "puntos_ganados": 0
}
HTTP 201

# Crear piloto histórico
POST http://localhost:8080/pilotos_historicos
Content-Type: application/json

{
  "id_piloto": 99,
  "nombre": "Piloto Prueba",
  "pais": "Argentina",
  "fecha_nacimiento": "2000-01-01T00:00:00Z",
  "titulos_ganados": 0
}
HTTP 201

# Crear escudería
POST http://localhost:8080/escuderias
Content-Type: application/json

{
  "id_escuderia": "FERRARI_TEST",
  "puntos_temporada": 0,
  "titulos_constructores": 0,
  "fecha_fundacion": "2000-01-01T00:00:00Z",
  "team_principal": "Jefe de Prueba"
}
HTTP 201

# Crear piloto activo
POST http://localhost:8080/pilotos
Content-Type: application/json

{
  "id_piloto": 99,
  "id_escuderia": "FERRARI_TEST",
  "puntos_temporada": 0
}
HTTP 201

# Crear Gran Premio
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

# Crear edición histórica de la carrera
POST http://localhost:8080/gran_premios_historicos
Content-Type: application/json

{
  "id_gran_premio": "GP_TEST_APUESTA",
  "fecha_carrera": "2026-11-22T00:00:00Z",
  "resultado_carrera": [99]
}
HTTP 201

# Crear apuesta
POST http://localhost:8080/apuestas
Content-Type: application/json

{
  "nombre_usuario": "gian_apuesta_test",
  "id_gran_premio": "GP_TEST_APUESTA",
  "fecha_carrera": "2026-11-22T00:00:00Z",
  "prediccion": [99],
  "fecha_apuesta": "2026-10-14T00:00:00Z"
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
jsonpath "$[?(@.id_apuesta == {{id_apuesta}})].nombre_usuario" contains "gian_apuesta_test"

# Listar apuestas de un usuario
GET http://localhost:8080/apuestas/gian_apuesta_test
HTTP 200
[Asserts]
jsonpath "$" isCollection
jsonpath "$[?(@.id_apuesta == {{id_apuesta}})].id_gran_premio" contains "GP_TEST_APUESTA"

# Obtener apuesta por ID y usuario
GET http://localhost:8080/apuestas/{{id_apuesta}}/gian_apuesta_test

HTTP 200
[Asserts]
jsonpath "$.id_apuesta" == {{id_apuesta}}
jsonpath "$.id_gran_premio" == "GP_TEST_APUESTA"
jsonpath "$.prediccion" exists

# Modificar apuesta
PUT http://localhost:8080/apuestas/{{id_apuesta}}/gian_apuesta_test
Content-Type: application/json

{
  "prediccion": [1, 99, 2, 3, 4, 5, 6, 7, 8, 9],
  "fecha_apuesta": "2026-10-14T00:00:00Z",
  "id_gran_premio": "GP_TEST_APUESTA",
  "fecha_carrera": "2026-11-22T00:00:00Z"
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
HTTP 200

# Verificar eliminación
GET http://localhost:8080/apuestas/{{id_apuesta}}/gian_apuesta_test
HTTP 404

# Eliminar edición histórica
DELETE http://localhost:8080/gran_premios_historicos/GP_TEST_APUESTA/22-11-2026
HTTP 200

# Eliminar Gran Premio
DELETE http://localhost:8080/gran_premios/GP_TEST_APUESTA
HTTP 200

# Eliminar piloto activo
DELETE http://localhost:8080/pilotos/99
HTTP 200

# Eliminar piloto histórico
DELETE http://localhost:8080/pilotos_historicos/99
HTTP 200

# Eliminar escudería
DELETE http://localhost:8080/escuderias/FERRARI_TEST
HTTP 200

# Eliminar usuario
DELETE http://localhost:8080/usuarios/gian_apuesta_test
HTTP 200