-- 1. Usuarios
CREATE TABLE usuario (
    id_usuario SERIAL PRIMARY KEY,
    nombre_usuario VARCHAR(50) NOT NULL,
    contrasena_hash VARCHAR(255) NOT NULL,
    puntos_ganados INT NOT NULL
);

-- 2. Escuderías
CREATE TABLE escuderia (
    id_escuderia VARCHAR(100) PRIMARY KEY,
    puntos_temporada INT NOT NULL,
    titulos_constructores INT NOT NULL,
    fecha_fundacion DATE NOT NULL,
    team_principal VARCHAR(150) NOT NULL
);

-- 3. Piloto Histórico
CREATE TABLE piloto_historico (
    id_piloto INT PRIMARY KEY,
    nombre VARCHAR(150) NOT NULL,
    pais VARCHAR(100) NOT NULL,
    fecha_nacimiento DATE NOT NULL,
    titulos_ganados INT NOT NULL
    );

-- 4. Pilotos de la Temporada
CREATE TABLE piloto (
    id_piloto INT PRIMARY KEY,
    id_escuderia VARCHAR(100) NOT NULL,
    puntos_temporada INT NOT NULL,

    CONSTRAINT fk_id_piloto FOREIGN KEY (id_piloto)
        REFERENCES piloto_historico(id_piloto)
        ON UPDATE CASCADE
        ON DELETE CASCADE,

    CONSTRAINT fk_piloto_escuderia FOREIGN KEY (id_escuderia)
        REFERENCES escuderia(id_escuderia)
        ON UPDATE CASCADE
        ON DELETE CASCADE
);

-- 5. Gran Premio Histórico (Circuitos)
CREATE TABLE gran_premio (
    id_gran_premio VARCHAR(100) PRIMARY KEY,
    pais VARCHAR(100) NOT NULL,
    longitud_km DECIMAL(5,3) NOT NULL,
    cantidad_vueltas INT NOT NULL,
    id_ultimo_ganador INT  NOT NULL,

    CONSTRAINT fk_gran_premio_ultimo_ganador FOREIGN KEY (id_ultimo_ganador)
        REFERENCES piloto_historico(id_piloto)
        ON DELETE SET NULL
);

-- 6. Gran Premio Histórico (Ediciones específicas de carreras)
CREATE TABLE gran_premio_historico (
    id_gran_premio VARCHAR(100) NOT NULL,
    fecha_carrera DATE  NOT NULL,
    resultado_carrera INTEGER[22],

    CONSTRAINT pk_gran_premio_historico
        PRIMARY KEY (id_gran_premio, fecha_carrera),

    CONSTRAINT fk_id_gran_premio FOREIGN KEY (id_gran_premio)
        REFERENCES gran_premio(id_gran_premio)
        ON DELETE CASCADE
);

-- 7. Apuestas
CREATE TABLE apuesta (
    id_apuesta SERIAL,
    id_usuario INT NOT NULL,
    id_gran_premio VARCHAR(100) NOT NULL,
    fecha_carrera DATE  NOT NULL,
    prediccion INTEGER[10] NOT NULL,
    fecha_apuesta DATE NOT NULL,

    CONSTRAINT pk_apuesta
        PRIMARY KEY (id_apuesta, id_usuario),

    CONSTRAINT fk_apuesta_usuario FOREIGN KEY (id_usuario)
        REFERENCES usuario(id_usuario)
        ON DELETE CASCADE,

    CONSTRAINT fk_apuesta_gran_premio
        FOREIGN KEY (id_gran_premio, fecha_carrera)
        REFERENCES gran_premio_historico(id_gran_premio, fecha_carrera)
        ON DELETE RESTRICT,

    -- Garantiza una sola apuesta por usuario en cada edición/carrera concreta
    CONSTRAINT unq_usuario_gran_premio_carrera
        UNIQUE (id_usuario, id_gran_premio, fecha_carrera)
);
