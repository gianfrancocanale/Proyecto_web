-- ==========================================
-- USUARIOS
-- ==========================================

-- name: CrearUsuario :one
-- Crear un nuevo usuario
INSERT INTO usuario (
    nombre_usuario,
    contrasena_hash,
    puntos_ganados
) VALUES (
    $1, $2, $3
)
RETURNING id_usuario, nombre_usuario, puntos_ganados;


-- name: RecuperarUsuarioPorId :one
-- Obtener usuario por ID
SELECT
    id_usuario,
    nombre_usuario,
    puntos_ganados
FROM usuario
WHERE id_usuario = $1;

-- name: RecuperarUsuarios :many
-- Obtener todos los usuarios
SELECT
    id_usuario,
    nombre_usuario,
    puntos_ganados
FROM usuario;


-- name: ModificarUsuario :one
-- Modificar datos de usuario
UPDATE usuario
SET
    nombre_usuario = COALESCE($2, nombre_usuario),
    contrasena_hash = COALESCE($3, contrasena_hash),
    puntos_ganados = COALESCE($4, puntos_ganados)
WHERE id_usuario = $1
RETURNING id_usuario, nombre_usuario, puntos_ganados;


-- name: EliminarUsuario :exec
-- Eliminar un usuario
DELETE FROM usuario
WHERE id_usuario = $1;

-- ==========================================
-- ESCUDERÍAS
-- ==========================================  

-- name: CrearEscuderia :one
-- Registrar una nueva escudería
INSERT INTO escuderia (
    id_escuderia,
    puntos_temporada,
    titulos_constructores,
    fecha_fundacion,
    team_principal
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING *;


-- name: ListarEscuderias :many
-- Mostrar todas las escuderías ordenadas por puntos
SELECT
    id_escuderia,
    puntos_temporada,
    titulos_constructores,
    fecha_fundacion,
    team_principal
FROM escuderia
ORDER BY puntos_temporada DESC;

-- name: RecuperarEscuderia :one
-- Obtener datos de una escudería por ID
SELECT
    id_escuderia,
    puntos_temporada,
    titulos_constructores,
    fecha_fundacion,
    team_principal
FROM escuderia
WHERE id_escuderia = $1;

-- name: ModificarEscuderia :one
-- Actualizar datos de una escudería
UPDATE escuderia
SET
    puntos_temporada = COALESCE($2, puntos_temporada),
    titulos_constructores = COALESCE($3, titulos_constructores),
    fecha_fundacion = COALESCE($4, fecha_fundacion),
    team_principal = COALESCE($5, team_principal)
WHERE id_escuderia = $1
RETURNING *;

-- name: EliminarEscuderia :exec
-- Eliminar una escudería
DELETE FROM escuderia
WHERE id_escuderia = $1;


-- ==========================================
-- PILOTOS HISTÓRICOS
-- ==========================================

-- name: CrearPilotoHistorico :one
-- Registrar datos biográficos/históricos de un piloto
INSERT INTO piloto_historico (
    id_piloto,
    nombre,
    pais,
    fecha_nacimiento,
    titulos_ganados
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING *;


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

-- name: RecuperarPilotoHistorico :one
-- Obtener datos de un piloto histórico por ID
SELECT
    id_piloto,
    nombre,
    pais,
    fecha_nacimiento,
    titulos_ganados
FROM piloto_historico
WHERE id_piloto = $1;

-- name: ModificarPilotoHistorico :one
-- Actualizar datos de un piloto histórico
UPDATE piloto_historico
SET
    nombre = COALESCE($2, nombre),
    pais = COALESCE($3, pais),
    fecha_nacimiento = COALESCE($4, fecha_nacimiento),
    titulos_ganados = COALESCE($5, titulos_ganados)
WHERE id_piloto = $1
RETURNING *;

-- name: EliminarPilotoHistorico :exec
-- Eliminar un piloto histórico
DELETE FROM piloto_historico
WHERE id_piloto = $1;


-- ==========================================
-- PILOTOS DE LA TEMPORADA
-- ==========================================

-- name: CrearPilotoTemporada :one
-- Inscribir a un piloto histórico en la temporada actual
INSERT INTO piloto (
    id_piloto,
    id_escuderia,
    puntos_temporada
) VALUES (
    $1,
    $2,
    $3
)
RETURNING *;


-- name: ListarPilotos :many
-- Listar pilotos activos en la temporada
-- unificando sus datos históricos
SELECT
    p.id_piloto,
    ph.nombre,
    ph.pais,
    ph.fecha_nacimiento,
    ph.titulos_ganados,
    p.id_escuderia,
    p.puntos_temporada
FROM piloto p
INNER JOIN piloto_historico ph
    ON p.id_piloto = ph.id_piloto
ORDER BY p.puntos_temporada DESC, ph.nombre ASC;

-- name: EliminarPilotoTemporada :exec
-- Eliminar un piloto de la temporada
DELETE FROM piloto
WHERE id_piloto = $1;


-- ==========================================
-- GRANDES PREMIOS / CIRCUITOS
-- ==========================================

-- name: CrearGranPremio :one
-- Registrar un circuito en el catálogo histórico
INSERT INTO gran_premio (
    id_gran_premio,
    pais,
    longitud_km,
    cantidad_vueltas,
    id_ultimo_ganador
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING *;


-- name: ListarGrandesPremios :many
-- Listar circuitos con la información de su último ganador
SELECT
    gp.id_gran_premio,
    gp.pais,
    gp.longitud_km,
    gp.cantidad_vueltas,
    gp.id_ultimo_ganador,
    ph.nombre AS nombre_ultimo_ganador
FROM gran_premio gp
LEFT JOIN piloto_historico ph
    ON gp.id_ultimo_ganador = ph.id_piloto
ORDER BY gp.id_gran_premio ASC;

-- name: EliminarGranPremio :exec
-- Eliminar un gran premio
DELETE FROM gran_premio
WHERE id_gran_premio = $1;


-- ==========================================
-- EDICIONES / EVENTOS DE GRANDES PREMIOS
-- ==========================================

-- name: CrearGranPremioHistorico :one
-- Programar una fecha/carrera concreta de un Gran Premio
INSERT INTO gran_premio_historico (
    id_gran_premio,
    fecha_carrera,
    resultado_carrera
) VALUES (
    $1,
    $2,
    $3
)
RETURNING *;


-- name: ListarGrandesPremiosHistoricos :many
-- Obtener las carreras programadas unificando
-- la información del circuito
SELECT
    gph.id_gran_premio,
    gph.fecha_carrera,
    gph.resultado_carrera,
    gp.pais,
    gp.longitud_km,
    gp.cantidad_vueltas
FROM gran_premio_historico gph
INNER JOIN gran_premio gp
    ON gph.id_gran_premio = gp.id_gran_premio
ORDER BY gph.fecha_carrera ASC;


-- name: RegistrarResultadoCarrera :one
-- Cargar el resultado final de una carrera
UPDATE gran_premio_historico
SET resultado_carrera = $3
WHERE id_gran_premio = $1
  AND fecha_carrera = $2
RETURNING *;

-- name: EliminarGranPremioHistorico :exec
-- Eliminar un gran premio histórico
DELETE FROM gran_premio_historico
WHERE id_gran_premio = $1 AND fecha_carrera = $2;

-- ==========================================
-- apuestaS
-- ==========================================

-- name: CrearApuesta :one
-- Crear una apuesta vinculada a un evento de carrera específico
INSERT INTO apuesta (
    id_usuario,
    id_gran_premio,
    fecha_carrera,
    prediccion,
    fecha_apuesta
) VALUES (
    $1,
    $2,
    $3,
    $4,
    CURRENT_TIMESTAMP
)
RETURNING
    id_apuesta,
    id_usuario,
    id_gran_premio,
    fecha_carrera,
    prediccion,
    fecha_apuesta;


-- name: RecuperarApuesta :many
-- Obtener apuestas por ID
SELECT
    id_apuesta,
    id_usuario,
    id_gran_premio,
    fecha_carrera,
    prediccion,
    fecha_apuesta
FROM apuesta
WHERE id_apuesta = $1;


-- name: ListarApuestasPorUsuario :many
-- Listar apuestas registradas por un usuario
SELECT
    id_apuesta,
    id_gran_premio,
    fecha_carrera,
    prediccion,
    fecha_apuesta
FROM apuesta
WHERE id_usuario = $1
ORDER BY fecha_apuesta DESC;


-- name: ModificArapuesta :one
-- Actualizar la predicción de una apuesta
UPDATE apuesta
SET
    prediccion = $2,
    fecha_apuesta = CURRENT_TIMESTAMP
WHERE id_apuesta = $1
RETURNING
    id_apuesta,
    id_usuario,
    id_gran_premio,
    fecha_carrera,
    prediccion,
    fecha_apuesta;


-- name: EliminarApuesta :exec
-- Eliminar una apuesta
DELETE FROM apuesta
WHERE id_apuesta = $1;
