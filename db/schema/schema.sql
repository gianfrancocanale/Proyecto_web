-- 1. Usuarios
CREATE TABLE usuarios (
    id_usuario SERIAL PRIMARY KEY,
    nombre_usuario VARCHAR(50) UNIQUE NOT NULL,
    contrasena_hash VARCHAR(255) NOT NULL,
    puntos_ganados INT DEFAULT 0 CHECK (puntos_ganados >= 0)
);

-- 2. Escuderías
CREATE TABLE escuderias (
    id_escuderia VARCHAR(100) PRIMARY KEY,
    puntos_temporada INT DEFAULT 0 CHECK (puntos_temporada >= 0),
    titulos_constructores INT DEFAULT 0 CHECK (titulos_constructores >= 0),
    fecha_fundacion DATE,
    team_principal VARCHAR(150) NOT NULL
);

-- 3. Piloto Histórico
CREATE TABLE piloto_historico (
    id_piloto INT PRIMARY KEY, 
    nombre VARCHAR(150) NOT NULL,
    pais VARCHAR(100) NOT NULL,
    fecha_nacimiento DATE NOT NULL,
    titulos_ganados INT DEFAULT 0 CHECK (titulos_ganados >= 0)
);

-- 4. Pilotos de la Temporada
CREATE TABLE pilotos (
    id_piloto INT PRIMARY KEY, 
    id_escuderia VARCHAR(100),
    puntos_temporada INT DEFAULT 0 CHECK (puntos_temporada >= 0),
    
    CONSTRAINT fk_id_piloto FOREIGN KEY (id_piloto) 
        REFERENCES piloto_historico(id_piloto) ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT fk_piloto_escuderia FOREIGN KEY (id_escuderia) 
        REFERENCES escuderias(id_escuderia) ON UPDATE CASCADE ON DELETE SET NULL
);

-- 5. Gran Premio Histórico (Circuitos)
CREATE TABLE gran_premio_historico (
    id_gran_premio VARCHAR(100) PRIMARY KEY,
    pais VARCHAR(100) NOT NULL,
    id_ultimo_ganador INT,
    longitud_km DECIMAL(5,3) CHECK (longitud_km > 0),
    cantidad_vueltas INT CHECK (cantidad_vueltas > 0),
    
    CONSTRAINT fk_gran_premio_ultimo_ganador FOREIGN KEY (id_ultimo_ganador) 
        REFERENCES pilotos(id_piloto) ON DELETE SET NULL
);

-- 6. Gran Premio (Ediciones específicas de carreras)
CREATE TABLE gran_premio (
    id_gran_premio VARCHAR(100),
    resultado_carrera INTEGER[22],
    fecha_carrera TIMESTAMP WITH TIME ZONE,
    
    CONSTRAINT pk_gran_premio PRIMARY KEY(id_gran_premio, fecha_carrera),
    CONSTRAINT fk_id_gran_premio FOREIGN KEY (id_gran_premio) 
        REFERENCES gran_premio_historico(id_gran_premio) ON DELETE CASCADE
);

-- 7. Apuestas
CREATE TABLE apuestas (
    id_apuesta SERIAL PRIMARY KEY,
    id_usuario INT NOT NULL,
    id_gran_premio VARCHAR(100) NOT NULL,
    fecha_carrera TIMESTAMP WITH TIME ZONE NOT NULL,
    prediccion INTEGER[10] NOT NULL,
    fecha_apuesta TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT fk_apuesta_usuario FOREIGN KEY (id_usuario) 
        REFERENCES usuarios(id_usuario) ON DELETE CASCADE,
    CONSTRAINT fk_apuesta_gran_premio FOREIGN KEY (id_gran_premio, fecha_carrera) 
        REFERENCES gran_premio(id_gran_premio, fecha_carrera) ON DELETE RESTRICT,
        
    -- Garantiza una sola apuesta por usuario en cada edición/carrera concreta
    CONSTRAINT unq_usuario_gran_premio_carrera UNIQUE (id_usuario, id_gran_premio, fecha_carrera)
);
