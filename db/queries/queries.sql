-- name: CrearApuesta :one
-- Crear una nueva apuesta en la base de datos
INSERT INTO apuestas (
    id_usuario,
    id_gran_premio,
    prediccion
) VALUES (
    $1, $2, $3
) RETURNING id_apuesta, id_usuario, id_gran_premio, prediccion, fecha_apuesta;

-- name: RecuperarApuesta :one
-- Obtener una apuesta por su ID
SELECT 
    id_apuesta,
    id_usuario,
    id_gran_premio,
    prediccion,
    fecha_apuesta
FROM apuestas
WHERE id_apuesta = $1;

-- name: ListarApuestasPorUsuario :many
-- Listar todas las apuestas registradas por un usuario
SELECT 
    id_apuesta,
    id_gran_premio,
    prediccion,
    fecha_apuesta
FROM apuestas
WHERE id_usuario = $1
ORDER BY fecha_apuesta DESC;

-- name: ModificarApuesta :exec
-- Actualizar la predicción de una apuesta por su ID
UPDATE apuestas
SET 
    prediccion = $2,
    fecha_apuesta = CURRENT_TIMESTAMP
WHERE id_apuesta = $1
RETURNING id_apuesta, id_usuario, id_gran_premio, prediccion, fecha_apuesta;

-- name: BorrarApuesta :exec
-- Eliminar una apuesta por su ID
DELETE FROM apuestas
WHERE id_apuesta = $1;


-- name: ListarPilotos :many
-- Mostrar todos los pilotos registrados
SELECT 
    id_piloto,
    nombre,
    pais,
    fecha_nacimiento,
    id_escuderia,
    titulos_ganados,
    puntos_temporada
FROM pilotos
ORDER BY puntos_temporada DESC, nombre ASC;

-- name: ListarEscuderias :many
-- Mostrar todas las escuderías registradas
SELECT 
    id_escuderia,
    puntos_temporada,
    titulos_constructores,
    fecha_fundacion,
    team_principal
FROM escuderias
ORDER BY puntos_temporada DESC;

-- name: ListarGrandesPremios :many
-- Mostrar todos los Grandes Premios registrados
SELECT 
    id_gran_premio,
    pais,
    id_ultimo_ganador,
    longitud_km,
    cantidad_vueltas
FROM gran_premios
ORDER BY id_gran_premio ASC;


-- name: CrearUsuario :one
-- Crear un nuevo usuario
INSERT INTO usuarios (
    nombre_usuario,
    contrasena_hash,
    puntos_ganados
) VALUES (
    $1, $2, COALESCE($3, 0)
) RETURNING id_usuario, nombre_usuario, puntos_ganados;

-- name: ModificarUsuario :exec
-- Modificar un usuario existente
UPDATE usuarios
SET 
    nombre_usuario = COALESCE($2, nombre_usuario),
    contrasena_hash = COALESCE($3, contrasena_hash),
    puntos_ganados = COALESCE($4, puntos_ganados)
WHERE id_usuario = $1
RETURNING id_usuario, nombre_usuario, puntos_ganados;
