
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

HTTP 201


# Crear circuito necesario para la clave foránea
POST http://localhost:8080/gran_premios
Content-Type: application/json

{
  "id_gran_premio": "SILVERSTONE",
  "pais": "Reino Unido",
  "longitud_km": 5,
  "cantidad_vueltas": 52,
  "id_ultimo_ganador": 44
}

HTTP 201


# Crear edición de Gran Premio
POST http://localhost:8080/gran_premios_historicos
Content-Type: application/json

{
  "id_gran_premio": "SILVERSTONE",
  "fecha_carrera": "2026-07-05",
  "resultado_carrera": [44, 1, 16]
}

HTTP 201
[Asserts]
jsonpath "$.id_gran_premio" == "SILVERSTONE"
jsonpath "$.fecha_carrera" exists
jsonpath "$.resultado_carrera[0]" == 44
jsonpath "$.resultado_carrera[1]" == 1
jsonpath "$.resultado_carrera[2]" == 16


# Listar ediciones de Grandes Premios
GET http://localhost:8080/gran_premios_historicos

HTTP 200
[Asserts]
jsonpath "$" isCollection
jsonpath "$[?(@.id_gran_premio == 'SILVERSTONE')].fecha_carrera" includes "2026-07-05"


# Obtener edición por circuito y fecha
GET http://localhost:8080/gran_premios_historicos/SILVERSTONE/2026-07-05

HTTP 200
[Asserts]
jsonpath "$.id_gran_premio" == "SILVERSTONE"
jsonpath "$.fecha_carrera" exists
jsonpath "$.resultado_carrera[0]" == 44


# Registrar resultado de la carrera
PUT http://localhost:8080/gran_premios_historicos/SILVERSTONE/2026-07-05
Content-Type: application/json

{
  "resultado_carrera": [1, 44, 16]
}

HTTP 200
[Asserts]
jsonpath "$.id_gran_premio" == "SILVERSTONE"
jsonpath "$.resultado_carrera[0]" == 1
jsonpath "$.resultado_carrera[1]" == 44
jsonpath "$.resultado_carrera[2]" == 16


# Verificar resultado actualizado
GET http://localhost:8080/gran_premios_historicos/SILVERSTONE/2026-07-05

HTTP 200
[Asserts]
jsonpath "$.resultado_carrera[0]" == 1
jsonpath "$.resultado_carrera[1]" == 44


# Eliminar edición del Gran Premio
DELETE http://localhost:8080/gran_premios_historicos/SILVERSTONE/2026-07-05

HTTP 200


# Verificar eliminación de la edición
GET http://localhost:8080/gran_premios_historicos/SILVERSTONE/2026-07-05

HTTP 404


# Limpiar circuito de prueba
DELETE http://localhost:8080/gran_premios/SILVERSTONE

HTTP 200


# Limpiar piloto histórico de prueba
DELETE http://localhost:8080/pilotos_historicos/44

HTTP 200