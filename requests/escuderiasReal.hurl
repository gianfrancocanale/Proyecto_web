# ============ CREAR ============

# Crear escudería
POST http://localhost:8080/escuderias
Content-Type: application/json

{
  "id_escuderia": "mercedes",
  "puntos_temporada": 0,
  "titulos_constructores": 0,
  "fecha_fundacion": "1954-01-01T00:00:00Z",
  "team_principal": "Toto Wolff"
}

HTTP 201

# Crear escudería
POST http://localhost:8080/escuderias
Content-Type: application/json

{
  "id_escuderia": "ferrari",
  "puntos_temporada": 0,
  "titulos_constructores": 0,
  "fecha_fundacion": "1929-11-16T00:00:00Z",
  "team_principal": "Frédéric Vasseur"
}

HTTP 201

# Crear escudería
POST http://localhost:8080/escuderias
Content-Type: application/json

{
  "id_escuderia": "red_bull",
  "puntos_temporada": 0,
  "titulos_constructores": 0,
  "fecha_fundacion": "2005-01-01T00:00:00Z",
  "team_principal": "Laurent Mekies"
}

HTTP 201

# Crear escudería
POST http://localhost:8080/escuderias
Content-Type: application/json

{
  "id_escuderia": "mclaren",
  "puntos_temporada": 0,
  "titulos_constructores": 0,
  "fecha_fundacion": "1963-01-01T00:00:00Z",
  "team_principal": "Andrea Stella"
}

HTTP 201

# Crear escudería
POST http://localhost:8080/escuderias
Content-Type: application/json

{
  "id_escuderia": "racing_bulls",
  "puntos_temporada": 0,
  "titulos_constructores": 0,
  "fecha_fundacion": "1985-01-01T00:00:00Z",
  "team_principal": "Alan Permane"
}

HTTP 201

# Crear escudería
POST http://localhost:8080/escuderias
Content-Type: application/json

{
  "id_escuderia": "alpine",
  "puntos_temporada": 0,
  "titulos_constructores": 0,
  "fecha_fundacion": "1981-01-01T00:00:00Z",
  "team_principal": "Flavio Briatore"
}

HTTP 201

# Crear escudería
POST http://localhost:8080/escuderias
Content-Type: application/json

{
  "id_escuderia": "haas",
  "puntos_temporada": 0,
  "titulos_constructores": 0,
  "fecha_fundacion": "2014-01-01T00:00:00Z",
  "team_principal": "Ayao Komatsu"
}

HTTP 201

# Crear escudería
POST http://localhost:8080/escuderias
Content-Type: application/json

{
  "id_escuderia": "audi",
  "puntos_temporada": 0,
  "titulos_constructores": 0,
  "fecha_fundacion": "1993-01-01T00:00:00Z",
  "team_principal": "Jonathan Wheatley"
}

HTTP 201

# Crear escudería
POST http://localhost:8080/escuderias
Content-Type: application/json

{
  "id_escuderia": "williams",
  "puntos_temporada": 0,
  "titulos_constructores": 0,
  "fecha_fundacion": "1977-12-16T00:00:00Z",
  "team_principal": "James Vowles"
}

HTTP 201

# Crear escudería
POST http://localhost:8080/escuderias
Content-Type: application/json

{
  "id_escuderia": "aston_martin",
  "puntos_temporada": 0,
  "titulos_constructores": 0,
  "fecha_fundacion": "2021-01-01T00:00:00Z",
  "team_principal": "Adrian Newey"
}

HTTP 201

# Crear escudería
POST http://localhost:8080/escuderias
Content-Type: application/json

{
  "id_escuderia": "cadillac",
  "puntos_temporada": 0,
  "titulos_constructores": 0,
  "fecha_fundacion": "2025-01-01T00:00:00Z",
  "team_principal": "Graeme Lowdon"
}

HTTP 201


# ============ ACTUALIZAR PUNTOS Y TÍTULOS ============

# Actualizar escudería
PUT http://localhost:8080/escuderias/mercedes
Content-Type: application/json

{
  "puntos_temporada": 607,
  "titulos_constructores": 8,
  "fecha_fundacion": "1954-01-01T00:00:00Z",
  "team_principal": "Toto Wolff"
}

HTTP 200

# Actualizar escudería
PUT http://localhost:8080/escuderias/ferrari
Content-Type: application/json

{
  "puntos_temporada": 474,
  "titulos_constructores": 16,
  "fecha_fundacion": "1929-11-16T00:00:00Z",
  "team_principal": "Frédéric Vasseur"
}

HTTP 200

# Actualizar escudería
PUT http://localhost:8080/escuderias/red_bull
Content-Type: application/json

{
  "puntos_temporada": 355,
  "titulos_constructores": 6,
  "fecha_fundacion": "2005-01-01T00:00:00Z",
  "team_principal": "Laurent Mekies"
}

HTTP 200

# Actualizar escudería
PUT http://localhost:8080/escuderias/mclaren
Content-Type: application/json

{
  "puntos_temporada": 355,
  "titulos_constructores": 10,
  "fecha_fundacion": "1963-01-01T00:00:00Z",
  "team_principal": "Andrea Stella"
}

HTTP 200

# Actualizar escudería
PUT http://localhost:8080/escuderias/racing_bulls
Content-Type: application/json

{
  "puntos_temporada": 102,
  "titulos_constructores": 0,
  "fecha_fundacion": "1985-01-01T00:00:00Z",
  "team_principal": "Alan Permane"
}

HTTP 200

# Actualizar escudería
PUT http://localhost:8080/escuderias/alpine
Content-Type: application/json

{
  "puntos_temporada": 69,
  "titulos_constructores": 0,
  "fecha_fundacion": "1981-01-01T00:00:00Z",
  "team_principal": "Flavio Briatore"
}

HTTP 200

# Actualizar escudería
PUT http://localhost:8080/escuderias/haas
Content-Type: application/json

{
  "puntos_temporada": 30,
  "titulos_constructores": 0,
  "fecha_fundacion": "2014-01-01T00:00:00Z",
  "team_principal": "Ayao Komatsu"
}

HTTP 200

# Actualizar escudería
PUT http://localhost:8080/escuderias/audi
Content-Type: application/json

{
  "puntos_temporada": 23,
  "titulos_constructores": 0,
  "fecha_fundacion": "1993-01-01T00:00:00Z",
  "team_principal": "Jonathan Wheatley"
}

HTTP 200

# Actualizar escudería
PUT http://localhost:8080/escuderias/williams
Content-Type: application/json

{
  "puntos_temporada": 12,
  "titulos_constructores": 9,
  "fecha_fundacion": "1977-12-16T00:00:00Z",
  "team_principal": "James Vowles"
}

HTTP 200

# Actualizar escudería
PUT http://localhost:8080/escuderias/aston_martin
Content-Type: application/json

{
  "puntos_temporada": 7,
  "titulos_constructores": 0,
  "fecha_fundacion": "2021-01-01T00:00:00Z",
  "team_principal": "Adrian Newey"
}

HTTP 200

# Actualizar escudería
PUT http://localhost:8080/escuderias/cadillac
Content-Type: application/json

{
  "puntos_temporada": 0,
  "titulos_constructores": 0,
  "fecha_fundacion": "2025-01-01T00:00:00Z",
  "team_principal": "Graeme Lowdon"
}

HTTP 200