-- name: CrearUsuario :one
-- Crear un nuevo usuario
INSERT INTO usuarios (
    nombre_usuario,
    contrasena_hash,
    puntos_ganados
) VALUES (
    $1, $2, COALESCE($3, 0)
) RETURNING id_usuario, nombre_usuario, puntos_ganados;

-- name: RecuperarUsuarioPorId :one
-- Obtener usuario por ID
SELECT id_usuario, nombre_usuario, puntos_ganados 
FROM usuarios 
WHERE id_usuario = $1;

-- name: ModificarUsuario :one
-- Modificar datos de usuario
UPDATE usuarios
SET 
    nombre_usuario = COALESCE($2, nombre_usuario),
    contrasena_hash = COALESCE($3, contrasena_hash),
    puntos_ganados = COALESCE($4, puntos_ganados)
WHERE id_usuario = $1
RETURNING id_usuario, nombre_usuario, puntos_ganados;

-- name: EliminarUsuario :exec
-- Eliminar un usuario
DELETE FROM usuarios 
WHERE id_usuario = $1;

-- name: CrearEscuderia :one
-- Registrar una nueva escudería
INSERT INTO escuderias (
    id_escuderia,
    puntos_temporada,
    titulos_constructores,
    fecha_fundacion,
    team_principal
) VALUES (
    $1, COALESCE($2, 0), COALESCE($3, 0), $4, $5
) RETURNING *;

-- name: ListarEscuderias :many
-- Mostrar todas las escuderías ordenadas por puntos
SELECT 
    id_escuderia,
    puntos_temporada,
    titulos_constructores,
    fecha_fundacion,
    team_principal
FROM escuderias
ORDER BY puntos_temporada DESC;


-- name: CrearPilotoHistorico :one
-- Registrar datos biográficos/históricos de un piloto
INSERT INTO piloto_historico (
    id_piloto,
    nombre,
    pais,
    fecha_nacimiento,
    titulos_ganados
) VALUES (
    $1, $2, $3, $4, COALESCE($5, 0)
) RETURNING *;

-- name: ListarPilotosHistoricos :many
-- Obtener catálogo histórico general de pilotos
SELECT 
    id_piloto,
    nombre,
    pais,
    fecha_nacimiento,
    titulos_ganados
FROM piloto_historico
ORDER BY titulos_ganados DESC, nombre ASC;

-- name: CrearPilotoTemporada :one
-- Inscribir a un piloto histórico en la temporada actual asignándole escudería
INSERT INTO pilotos (
    id_piloto,
    id_escuderia,
    puntos_temporada
) VALUES (
    $1, $2, COALESCE($3, 0)
) RETURNING *;

-- name: ListarPilotos :many
-- Listar pilotos activos en la temporada unificando sus datos históricos
SELECT 
    p.id_piloto,
    ph.nombre,
    ph.pais,
    ph.fecha_nacimiento,
    ph.titulos_ganados,
    p.id_escuderia,
    p.puntos_temporada
FROM pilotos p
INNER JOIN piloto_historico ph ON p.id_piloto = ph.id_piloto
ORDER BY p.puntos_temporada DESC, ph.nombre ASC;


-- name: CrearGranPremioHistorico :one
-- Registrar un circuito en el catálogo histórico
INSERT INTO gran_premio_historico (
    id_gran_premio,
    pais,
    id_ultimo_ganador,
    longitud_km,
    cantidad_vueltas
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: ListarGrandesPremiosHistoricos :many
-- Listar circuitos con la información de su último ganador
SELECT 
    gph.id_gran_premio,
    gph.pais,
    gph.longitud_km,
    gph.cantidad_vueltas,
    gph.id_ultimo_ganador,
    ph.nombre AS nombre_ultimo_ganador
FROM gran_premio_historico gph
LEFT JOIN piloto_historico ph ON gph.id_ultimo_ganador = ph.id_piloto
ORDER BY gph.id_gran_premio ASC;

-- name: CrearGranPremioEvento :one
-- Programar una fecha/carrera concreta de un Gran Premio
INSERT INTO gran_premio (
    id_gran_premio,
    fecha_carrera,
    resultado_carrera
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: ListarGrandesPremiosEventos :many
-- Obtener las carreras programadas unificando la información del circuito
SELECT 
    gp.id_gran_premio,
    gp.fecha_carrera,
    gp.resultado_carrera,
    gph.pais,
    gph.longitud_km,
    gph.cantidad_vueltas
FROM gran_premio gp
INNER JOIN gran_premio_historico gph ON gp.id_gran_premio = gph.id_gran_premio
ORDER BY gp.fecha_carrera ASC;

-- name: RegistrarResultadoCarrera :one
-- Cargar el resultado final de una carrera
UPDATE gran_premio
SET resultado_carrera = $3
WHERE id_gran_premio = $1 AND fecha_carrera = $2
RETURNING *;


-- name: CrearApuesta :one
-- Crear una apuesta vinculada a un evento de carrera específico
INSERT INTO apuestas (
    id_usuario,
    id_gran_premio,
    fecha_carrera,
    prediccion
) VALUES (
    $1, $2, $3, $4
) RETURNING id_apuesta, id_usuario, id_gran_premio, fecha_carrera, prediccion, fecha_apuesta;

-- name: RecuperarApuesta :one
-- Obtener una apuesta por su ID
SELECT 
    id_apuesta,
    id_usuario,
    id_gran_premio,
    fecha_carrera,
    prediccion,
    fecha_apuesta
FROM apuestas
WHERE id_apuesta = $1;

-- name: ListarApuestasPorUsuario :many
-- Listar apuestas registradas por un usuario
SELECT 
    id_apuesta,
    id_gran_premio,
    fecha_carrera,
    prediccion,
    fecha_apuesta
FROM apuestas
WHERE id_usuario = $1
ORDER BY fecha_apuesta DESC;

-- name: ModificarApuesta :one
-- Actualizar la predicción de una apuesta
UPDATE apuestas
SET 
    prediccion = $2,
    fecha_apuesta = CURRENT_TIMESTAMP
WHERE id_apuesta = $1
RETURNING id_apuesta, id_usuario, id_gran_premio, fecha_carrera, prediccion, fecha_apuesta;

-- name: BorrarApuesta :exec
-- Eliminar una apuesta
DELETE FROM apuestas
WHERE id_apuesta = $1;
