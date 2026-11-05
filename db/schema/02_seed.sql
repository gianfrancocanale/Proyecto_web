INSERT INTO escuderia (id_escuderia, puntos_temporada, titulos_constructores, fecha_fundacion, team_principal) VALUES
('mercedes',     630,  8, '1970-01-01', 'Toto Wolff'),
('ferrari',      504, 16, '1929-11-16', 'Frédéric Vasseur'),
('red_bull',     382,  6, '2005-01-01', 'Laurent Mekies'),
('mclaren',      365, 10, '1963-01-01', 'Andrea Stella'),
('racing_bulls', 108,  0, '1985-01-01', 'Alan Permane'),
('alpine',        69,  2, '1981-01-01', 'Flavio Briatore'),
('haas',          31,  0, '2014-01-01', 'Ayao Komatsu'),
('audi',          27,  0, '1970-01-01', 'Mattia Binotto'),
('williams',      12,  9, '1977-12-16', 'James Vowles'),
('aston_martin',   7,  0, '1991-01-01', 'Adrian Newey'),
('cadillac',       0,  0, '2023-01-01', 'Graeme Lowdon');


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
(12, 'mercedes',     380),
(3,  'red_bull',     261),
(44,  'ferrari',      259),
(63,  'mercedes',     250),
(16,  'ferrari',      245),
(1,  'mclaren',      221),
(81,  'mclaren',      144),
(6, 'red_bull',     107),
(30, 'racing_bulls',  83),
(10,  'alpine',        42),
(41, 'racing_bulls',  38),
(43, 'alpine',        27),
(87, 'haas',          23),
(5, 'audi',          14),
(27, 'audi',          13),
(14,  'aston_martin',   7),
(55, 'williams',       7),
(23, 'williams',       5),
(18, 'aston_martin',   0),
(77, 'cadillac',       0),
(11, 'cadillac',       0),
(51, 'haas',           0);

INSERT INTO gran_premio (id_gran_premio, pais, longitud_km, cantidad_vueltas, id_ultimo_ganador) VALUES
('albert_park',     'Australia',             5, 58, 63), -- George Russell
('shanghai',        'China',                 5, 56, 12), -- Kimi Antonelli
('suzuka',          'Japón',                 6, 53, 12), -- Kimi Antonelli
('jeddah',          'Arabia Saudita',        6, 50, 3), -- 
('miami',           'Estados Unidos',        5, 57, 12), -- Kimi Antonelli
('gilles_villeneuve','Canadá',               4, 70, 12), -- Kimi Antonelli
('monaco',          'Mónaco',                3, 78, 12), -- Kimi Antonelli
('catalunya',       'España (Barcelona)',    5, 66, 44), -- Lewis Hamilton
('red_bull_ring',   'Austria',               4, 71, 63), -- George Russell
('silverstone',     'Reino Unido',           6, 52, 16), -- Charles Leclerc
('spa',             'Bélgica',               7, 44, 12), -- Kimi Antonelli
('hungaroring',     'Hungría',               4, 70,  1), -- Lando Norris
('zandvoort',       'Países Bajos',          4, 72,  1), -- Lando Norris
('monza',           'Italia',                6, 53, 12), -- Kimi Antonelli
('madring',         'España (Madrid)',       5, 57, 12), -- Kimi Antonelli
('baku',            'Azerbaiyán',            6, 51, 63), -- George Russell
('sepang',          'Baréin (Malasia)',      6, 56,  3), -- Max Verstappen
('marina_bay',      'Singapur',              5, 62,  3), -- Max Verstappen
('cota',            'Estados Unidos (Austin)',5, 56, 12), -- Kimi Antonelli
('hermanos_rodriguez','México',              4, 71,  3), -- Max Verstappen
('interlagos',      'Brasil',                4, 71,  3), -- Pendiente
('las_vegas',       'Estados Unidos',        6, 50, 3), -- Pendiente por disputar
('losail',          'Catar',                 5, 57, 3), -- Pendiente por disputar
('yas_marina',      'Emiratos Árabes Unidos',5, 58, 3); -- Pendiente por disputar