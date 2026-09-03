CREATE TABLE apuestas (
    id_apuesta SERIAL PRIMARY KEY,
    id_usuario INT NOT NULL,
    id_gran_premio VARCHAR(100) NOT NULL,
    prediccion INTEGER[10] NOT NULL,
    fecha_apuesta TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    -- Restricciones de Claves Foráneas
    CONSTRAINT fk_apuesta_usuario FOREIGN KEY (id_usuario) 
        REFERENCES usuarios(id_usuario) ON DELETE CASCADE,
    CONSTRAINT fk_apuesta_gran_premio FOREIGN KEY (id_gran_premio) 
        REFERENCES gran_premios(id_gran_premio) ON DELETE RESTRICT,
        
    -- Garantiza que un usuario solo pueda hacer una apuesta por Gran Premio
    CONSTRAINT unq_usuario_gran_premio UNIQUE (id_usuario, id_gran_premio)
);

-- Tabla: escuderias
CREATE TABLE escuderias (
    id_escuderia VARCHAR(100) PRIMARY KEY,
    puntos_temporada INT DEFAULT 0 CHECK (puntos_temporada >= 0),
    titulos_constructores INT DEFAULT 0 CHECK (titulos_constructores >= 0),
    fecha_fundacion DATE,
    team_principal VARCHAR(150) NOT NULL
);

-- Tabla: pilotos
CREATE TABLE pilotos (
    id_piloto INT PRIMARY KEY, 
    nombre VARCHAR(150) NOT NULL,
    pais VARCHAR(100) NOT NULL,
    fecha_nacimiento DATE NOT NULL,
    id_escuderia VARCHAR(100) NOT NULL,
    titulos_ganados INT DEFAULT 0 CHECK (titulos_ganados >= 0),
    puntos_temporada INT DEFAULT 0 CHECK (puntos_temporada >= 0),
    
    CONSTRAINT fk_piloto_escuderia FOREIGN KEY (id_escuderia) 
        REFERENCES escuderias(id_escuderia) ON UPDATE CASCADE
);

-- Tabla: usuarios
CREATE TABLE usuarios (
    id_usuario SERIAL PRIMARY KEY,
    nombre_usuario VARCHAR(50) UNIQUE NOT NULL,
    contrasena_hash VARCHAR(255) NOT NULL,
    puntos_ganados INT DEFAULT 0 CHECK (puntos_ganados >= 0)
);

-- Tabla: gran_premios
CREATE TABLE gran_premios (
    id_gran_premio VARCHAR(100) PRIMARY KEY,
    pais VARCHAR(100) NOT NULL,
    id_ultimo_ganador INT,
    longitud_km DECIMAL(5,3) CHECK (longitud_km > 0),
    cantidad_vueltas INT CHECK (cantidad_vueltas > 0),
    
    CONSTRAINT fk_gran_premio_ultimo_ganador FOREIGN KEY (id_ultimo_ganador) 
        REFERENCES pilotos(id_piloto) ON DELETE SET NULL
);
