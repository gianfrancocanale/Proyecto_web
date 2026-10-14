
# =========================
# USUARIOS
# =========================

# Crear usuario
POST http://localhost:8080/usuarios
Content-Type: application/json

{
  "nombre_usuario": "gian",
  "contrasena_hash": "1234",
  "puntos_ganados": 0
}

HTTP/1.1 201
[Asserts]
jsonpath "$.nombre_usuario" == "gian"
jsonpath "$.puntos_ganados" == 0

[Captures]
nombre_usuario: jsonpath "$.nombre_usuario"


# Listar usuarios
GET http://localhost:8080/usuarios

HTTP/1.1 200
[Asserts]
jsonpath "$" isCollection
jsonpath "$[?(@.nombre_usuario == 'gian')].nombre_usuario" includes "gian"


# Obtener usuario
GET http://localhost:8080/usuarios/{{nombre_usuario}}

HTTP/1.1 200
[Asserts]
jsonpath "$.nombre_usuario" == "gian"
jsonpath "$.puntos_ganados" == 0


# Actualizar usuario
PUT http://localhost:8080/usuarios/{{nombre_usuario}}
Content-Type: application/json

{
  "nombre_usuario": "gian",
  "contrasena_hash": "1234",
  "puntos_ganados": 13
}

HTTP/1.1 200
[Asserts]
jsonpath "$.nombre_usuario" == "gian"
jsonpath "$.puntos_ganados" == 13


# Verificar actualización
GET http://localhost:8080/usuarios/{{nombre_usuario}}

HTTP/1.1 200
[Asserts]
jsonpath "$.nombre_usuario" == "gian"
jsonpath "$.puntos_ganados" == 13


# Eliminar usuario
DELETE http://localhost:8080/usuarios/{{nombre_usuario}}

HTTP/1.1 204


# Verificar eliminación
GET http://localhost:8080/usuarios/{{nombre_usuario}}

HTTP/1.1 404