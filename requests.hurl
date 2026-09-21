# =========================
# USUARIOS
# =========================

# Crear usuario
POST http://localhost:8080/usuarios
Content-Type: application/json

{
  "nombre_usuario": "gian",
  "contrasena_hash": "1234"
}

HTTP/1.1 201
[Asserts]
jsonpath "$.id_usuario" exists
jsonpath "$.nombre_usuario" == "gian"
jsonpath "$.puntos_ganados" == 0

[Captures]
id_usuario: jsonpath "$.id_usuario"


# Listar usuarios
GET http://localhost:8080/usuarios

HTTP/1.1 200
[Asserts]
jsonpath "$[0].nombre_usuario" == "gian"


# Obtener usuario
GET http://localhost:8080/usuarios/{{id_usuario}}

HTTP/1.1 200
[Asserts]
jsonpath "$.id_usuario" == {{id_usuario}}
jsonpath "$.nombre_usuario" == "gian"


# Actualizar usuario
PUT http://localhost:8080/usuarios/{{id_usuario}}
Content-Type: application/json

{
  "nombre_usuario": "gian",
  "contrasena_hash": "1234",
  "puntos_ganados": 13
}

HTTP/1.1 200
[Asserts]
jsonpath "$.id_usuario" == {{id_usuario}}
jsonpath "$.puntos_ganados" == 13


# Eliminar usuario
DELETE http://localhost:8080/usuarios/1

HTTP/1.1 204

# ==========================================
# PILOTOS HISTÓRICOS
# ==========================================

# Crear piloto histórico
POST http://localhost:8080/pilotos_historicos
Content-Type: application/json

{
  "id_piloto": 44,
  "nombre": "Max Verstappen",
  "pais": "Países Bajos",
  "fecha_nacimiento": "1997-09-30T00:00:00Z"
}

HTTP/1.1 201
[Asserts]
jsonpath "$.id_piloto" == 44
jsonpath "$.nombre" == "Max Verstappen"
jsonpath "$.pais" == "Países Bajos"
jsonpath "$.titulos_ganados" == 0


# Listar pilotos históricos
GET http://localhost:8080/pilotos_historicos

HTTP/1.1 200
[Asserts]
jsonpath "$" isCollection
jsonpath "$[0].id_piloto" exists


# Obtener piloto histórico
GET http://localhost:8080/pilotos_historicos/44

HTTP/1.1 200
[Asserts]
jsonpath "$.id_piloto" == 44
jsonpath "$.nombre" == "Max Verstappen"


# Actualizar piloto histórico
PUT http://localhost:8080/pilotos_historicos/44
Content-Type: application/json

{
  "nombre": "Max Verstappen",
  "pais": "Países Bajos",
  "fecha_nacimiento": "1997-09-30T00:00:00Z",
  "titulos_ganados": 4
}

HTTP/1.1 200
[Asserts]
jsonpath "$.id_piloto" == 44
jsonpath "$.nombre" == "Max Verstappen"
jsonpath "$.titulos_ganados" == 4


# Eliminar piloto histórico
DELETE http://localhost:8080/pilotos_historicos/44

HTTP/1.1 204