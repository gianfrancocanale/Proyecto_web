INSERT INTO escuderia (id_escuderia, puntos_temporada, titulos_constructores, fecha_fundacion, team_principal) VALUES
('Mercedes',     630,  8, '1970-01-01', 'Toto Wolff'),
('Ferrari',      504, 16, '1929-11-16', 'Frédéric Vasseur'),
('Red Bull',     382,  6, '2005-01-01', 'Laurent Mekies'),
('McLaren',      365, 10, '1963-01-01', 'Andrea Stella'),
('Racing Bulls', 108,  0, '1985-01-01', 'Alan Permane'),
('Alpine',        69,  2, '1981-01-01', 'Flavio Briatore'),
('Haas',          31,  0, '2014-01-01', 'Ayao Komatsu'),
('Audi',          27,  0, '1970-01-01', 'Mattia Binotto'),
('Williams',      12,  9, '1977-12-16', 'James Vowles'),
('Aston Martin',   7,  0, '1991-01-01', 'Adrian Newey'),
('Cadillac',       0,  0, '2023-01-01', 'Graeme Lowdon');


INSERT INTO piloto_historico (id_piloto, nombre, pais, fecha_nacimiento, titulos_ganados) VALUES
(44,  'Lewis Hamilton',      'Reino Unido',  '1985-01-07', 7),
(3,  'Max Verstappen',      'Países Bajos', '1997-09-30', 4),
(16,  'Charles Leclerc',     'Mónaco',       '1997-10-16', 0),
(1,  'Lando Norris',        'Reino Unido',  '1999-11-13', 1),
(14,  'Fernando Alonso',     'España',       '1981-07-29', 2),
(63,  'George Russell',      'Reino Unido',  '1998-02-15', 0),
(81,  'Oscar Piastri',       'Australia',    '2001-04-06', 0),
(10,  'Pierre Gasly',        'Francia',      '1996-02-07', 0),
(12, 'Kimi Antonelli',      'Italia',       '2006-08-25', 1),
(6, 'Isack Hadjar',        'Francia',      '2004-09-28', 0),
(30, 'Liam Lawson',         'Nueva Zelanda','2002-02-11', 0),
(41, 'Arvid Lindblad',      'Reino Unido',  '2007-08-08', 0),
(43, 'Franco Colapinto',    'Argentina',    '2003-05-27', 0),
(87, 'Oliver Bearman',      'Reino Unido',  '2005-05-08', 0),
(5, 'Gabriel Bortoleto',   'Brasil',       '2004-10-14', 0),
(27, 'Nico Hülkenberg',     'Alemania',     '1987-08-19', 0),
(51, 'Rafael Camara',        'Brasil',      '2005-05-05', 0),
(55, 'Carlos Sainz',        'España',       '1994-09-01', 0),
(23, 'Alex Albon',          'Tailandia',    '1996-03-23', 0),
(18, 'Lance Stroll',        'Canadá',       '1998-10-29', 0),
(77, 'Valtteri Bottas',     'Finlandia',    '1989-08-28', 0),
(11, 'Sergio Pérez',        'México',       '1990-01-26', 0);

INSERT INTO piloto (id_piloto, id_escuderia, puntos_temporada) VALUES
(12, 'Mercedes',     380),
(3,  'Red Bull',     261),
(44,  'Ferrari',      259),
(63,  'Mercedes',     250),
(16,  'Ferrari',      245),
(1,  'McLaren',      221),
(81,  'McLaren',      144),
(6, 'Red Bull',     107),
(30, 'Racing Bulls',  83),
(10,  'Alpine',        42),
(41, 'Racing Bulls',  38),
(43, 'Alpine',        27),
(87, 'Haas',          23),
(5, 'Audi',          14),
(27, 'Audi',          13),
(14,  'Aston Martin',   7),
(55, 'Williams',       7),
(23, 'Williams',       5),
(18, 'Aston Martin',   0),
(77, 'Cadillac',       0),
(11, 'Cadillac',       0),
(51, 'Haas',           0);

INSERT INTO gran_premio (id_gran_premio, pais, longitud_km, cantidad_vueltas, id_ultimo_ganador) VALUES
('Albert Park',     'Australia',             5, 58, 63), -- George Russell
('Shanghai',        'China',                 5, 56, 12), -- Kimi Antonelli
('Suzuka',          'Japón',                 6, 53, 12), -- Kimi Antonelli
('Miami',           'Estados Unidos',        5, 57, 12), -- Kimi Antonelli
('Gilles Villeneuve','Canadá',               4, 70, 12), -- Kimi Antonelli
('Monaco',          'Mónaco',                3, 78, 12), -- Kimi Antonelli
('Catalunya',       'España (Barcelona)',    5, 66, 44), -- Lewis Hamilton
('Red Bull Ring',   'Austria',               4, 71, 63), -- George Russell
('Silverstone',     'Reino Unido',           6, 52, 16), -- Charles Leclerc
('Spa',             'Bélgica',               7, 44, 12), -- Kimi Antonelli
('Hungaroring',     'Hungría',               4, 70,  1), -- Lando Norris
('Zandvoort',       'Países Bajos',          4, 72,  1), -- Lando Norris
('Monza',           'Italia',                6, 53, 12), -- Kimi Antonelli
('Madring',         'España (Madrid)',       5, 57, 12), -- Kimi Antonelli
('Baku',            'Azerbaiyán',            6, 51, 63), -- George Russell
('Sepang',          'Baréin (Malasia)',      6, 56,  3), -- Max Verstappen
('Marina Bay',      'Singapur',              5, 62,  3), -- Max Verstappen
('Cota',            'Estados Unidos (Austin)',5, 56, 12), -- Kimi Antonelli
('Hermanos Rodriguez','México',              4, 71,  3), -- Max Verstappen
('Interlagos',      'Brasil',                4, 71,  3), -- Pendiente
('Las Vegas',       'Estados Unidos',        6, 50, 3), -- Pendiente por disputar
('Losail',          'Catar',                 5, 57, 3), -- Pendiente por disputar
('Yas Marina',      'Emiratos Árabes Unidos',5, 58, 3); -- Pendiente por disputar