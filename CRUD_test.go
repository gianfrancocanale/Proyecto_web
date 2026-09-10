package main

import (
	"context"
	"database/sql"
	"testing"
	"time"

	sqlc "Proyecto_web/db/sqlc"

	_ "github.com/lib/pq"
)

// setupTestDB establece la conexión con la base de datos PostgreSQL de prueba.
func setupTestDB(t *testing.T) *sql.DB {
	connStr := "user=ChiaraGian password=ChiaraGian dbname=DB_PredictOne port=5432 sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatalf("Error al conectar con la base de datos: %v", err)
	}

	err = db.Ping()
	if err != nil {
		t.Fatalf("Error al hacer ping a la base de datos: %v", err)
	}

	return db
}

// cleanupTestDB elimina todos los registros respetando la integridad referencial (Foreign Keys).
func cleanupTestDB(t *testing.T, db *sql.DB) {
	tables := []string{
		"apuestas",
		"gran_premio",
		"pilotos",
		"gran_premio_historico",
		"piloto_historico",
		"escuderias",
		"usuarios",
	}

	for _, table := range tables {
		if _, err := db.Exec("DELETE FROM " + table); err != nil {
			t.Fatalf("Error al limpiar la tabla %s: %v", table, err)
		}
	}
}

// ==========================================
// TESTS DE USUARIOS
// ==========================================

func TestUsuarioRepository_CRUD(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	defer cleanupTestDB(t, db)

	queries := sqlc.New(db)
	ctx := context.Background()

	t.Run("CrearUsuario", func(t *testing.T) {
		cleanupTestDB(t, db)

		usuario, err := queries.CrearUsuario(ctx, sqlc.CrearUsuarioParams{
			NombreUsuario:  "usuario_test_1",
			ContrasenaHash: "hash_seguro_123",
			Column3:        int32(0),
		})

		if err != nil {
			t.Fatalf("CrearUsuario fallo: %v", err)
		}

		if usuario.IDUsuario == 0 {
			t.Error("Se esperaba un IDUsuario autogenerado")
		}

		if usuario.NombreUsuario != "usuario_test_1" {
			t.Errorf("Se esperaba 'usuario_test_1', se obtuvo '%s'", usuario.NombreUsuario)
		}
	})

	t.Run("EliminarUsuario", func(t *testing.T) {
		cleanupTestDB(t, db)

		usuario, err := queries.CrearUsuario(ctx, sqlc.CrearUsuarioParams{
			NombreUsuario:  "usuario_a_eliminar",
			ContrasenaHash: "hash_456",
			Column3:        int32(0),
		})
		if err != nil {
			t.Fatalf("Setup fallo: %v", err)
		}

		err = queries.EliminarUsuario(ctx, usuario.IDUsuario)
		if err != nil {
			t.Fatalf("EliminarUsuario fallo: %v", err)
		}
	})
}

// ==========================================
// TESTS DE ESCUDERÍAS
// ==========================================

func TestEscuderiaRepository_CRUD(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	defer cleanupTestDB(t, db)

	queries := sqlc.New(db)
	ctx := context.Background()

	t.Run("CrearYListarEscuderias", func(t *testing.T) {
		cleanupTestDB(t, db)

		_, err := queries.CrearEscuderia(ctx, sqlc.CrearEscuderiaParams{
			IDEscuderia:    "FERRARI",
			Column2:        int32(100),
			Column3:        int32(16),
			FechaFundacion: sql.NullTime{Time: time.Date(1929, 11, 16, 0, 0, 0, 0, time.UTC), Valid: true},
			TeamPrincipal:  "Fred Vasseur",
		})
		if err != nil {
			t.Fatalf("CrearEscuderia fallo: %v", err)
		}

		escuderias, err := queries.ListarEscuderias(ctx)
		if err != nil {
			t.Fatalf("ListarEscuderias fallo: %v", err)
		}

		if len(escuderias) != 1 {
			t.Errorf("Se esperaba 1 escuderia, se obtuvieron %d", len(escuderias))
		}

		if escuderias[0].IDEscuderia != "FERRARI" {
			t.Errorf("Se esperaba IDEscuderia 'FERRARI', se obtuvo '%s'", escuderias[0].IDEscuderia)
		}
	})
}

// ==========================================
// TESTS DE PILOTOS (HISTÓRICOS Y TEMPORADA)
// ==========================================

func TestPilotoRepository_CRUD(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	defer cleanupTestDB(t, db)

	queries := sqlc.New(db)
	ctx := context.Background()

	t.Run("CrearPilotoHistoricoYTemporada", func(t *testing.T) {
		cleanupTestDB(t, db)

		// 1. Crear Escudería (necesaria para la FK de Piloto)
		_, err := queries.CrearEscuderia(ctx, sqlc.CrearEscuderiaParams{
			IDEscuderia:   "RED_BULL",
			Column2:       int32(200),
			Column3:       int32(6),
			TeamPrincipal: "Christian Horner",
		})
		if err != nil {
			t.Fatalf("Setup Escuderia fallo: %v", err)
		}

		// 2. Crear Piloto Histórico
		pilotoHist, err := queries.CrearPilotoHistorico(ctx, sqlc.CrearPilotoHistoricoParams{
			IDPiloto:        1,
			Nombre:          "Max Verstappen",
			Pais:            "Países Bajos",
			FechaNacimiento: time.Date(1997, 9, 30, 0, 0, 0, 0, time.UTC),
			Column5:         int32(3),
		})
		if err != nil {
			t.Fatalf("CrearPilotoHistorico fallo: %v", err)
		}

		// 3. Crear Piloto en Temporada
		_, err = queries.CrearPilotoTemporada(ctx, sqlc.CrearPilotoTemporadaParams{
			IDPiloto:    pilotoHist.IDPiloto,
			IDEscuderia: sql.NullString{String: "RED_BULL", Valid: true},
			Column3:     int32(50),
		})
		if err != nil {
			t.Fatalf("CrearPilotoTemporada fallo: %v", err)
		}

		// 4. Listar Pilotos Activos
		pilotos, err := queries.ListarPilotos(ctx)
		if err != nil {
			t.Fatalf("ListarPilotos fallo: %v", err)
		}

		if len(pilotos) != 1 {
			t.Errorf("Se esperaba 1 piloto activo, se obtuvieron %d", len(pilotos))
		}

		if pilotos[0].Nombre != "Max Verstappen" {
			t.Errorf("Se esperaba 'Max Verstappen', se obtuvo '%s'", pilotos[0].Nombre)
		}

		// 5. Listar Pilotos Históricos
		historicos, err := queries.ListarPilotosHistoricos(ctx)
		if err != nil {
			t.Fatalf("ListarPilotosHistoricos fallo: %v", err)
		}

		if len(historicos) != 1 {
			t.Errorf("Se esperaba 1 piloto historico, se obtuvieron %d", len(historicos))
		}
	})
}

// ==========================================
// TESTS DE GRANDES PREMIOS
// ==========================================

func TestGranPremioRepository_CRUD(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	defer cleanupTestDB(t, db)

	queries := sqlc.New(db)
	ctx := context.Background()

	t.Run("CrearYListarGrandesPremios", func(t *testing.T) {
		cleanupTestDB(t, db)

		// 1. Crear Gran Premio Histórico (Circuito)
		gpHistorico, err := queries.CrearGranPremioHistorico(ctx, sqlc.CrearGranPremioHistoricoParams{
			IDGranPremio:    "MONACO",
			Pais:            "Mónaco",
			LongitudKm:      sql.NullString{String: "3.337", Valid: true},
			CantidadVueltas: sql.NullInt32{Int32: 78, Valid: true},
		})
		if err != nil {
			t.Fatalf("CrearGranPremioHistorico fallo: %v", err)
		}

		// 2. Programar Evento de Carrera
		fechaCarrera := time.Date(2026, 5, 24, 14, 0, 0, 0, time.UTC)
		_, err = queries.CrearGranPremioEvento(ctx, sqlc.CrearGranPremioEventoParams{
			IDGranPremio:     gpHistorico.IDGranPremio,
			FechaCarrera:     fechaCarrera,
			ResultadoCarrera: []int32{1, 2, 3},
		})
		if err != nil {
			t.Fatalf("CrearGranPremioEvento fallo: %v", err)
		}

		// 3. Listar Eventos
		eventos, err := queries.ListarGrandesPremiosEventos(ctx)
		if err != nil {
			t.Fatalf("ListarGrandesPremiosEventos fallo: %v", err)
		}

		if len(eventos) != 1 {
			t.Errorf("Se esperaba 1 evento, se obtuvieron %d", len(eventos))
		}

		if eventos[0].Pais != "Mónaco" {
			t.Errorf("Se esperaba 'Mónaco', se obtuvo '%s'", eventos[0].Pais)
		}

		// 4. Listar Históricos
		historicos, err := queries.ListarGrandesPremiosHistoricos(ctx)
		if err != nil {
			t.Fatalf("ListarGrandesPremiosHistoricos fallo: %v", err)
		}

		if len(historicos) != 1 {
			t.Errorf("Se esperaba 1 registro historico, se obtuvieron %d", len(historicos))
		}
	})
}

// ==========================================
// TESTS DE APUESTAS
// ==========================================

func TestApuestaRepository_CRUD(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	defer cleanupTestDB(t, db)

	queries := sqlc.New(db)
	ctx := context.Background()

	t.Run("FlujoCompletoApuesta", func(t *testing.T) {
		cleanupTestDB(t, db)

		// Setup de dependencias (Usuario y GranPremio)
		usr, err := queries.CrearUsuario(ctx, sqlc.CrearUsuarioParams{
			NombreUsuario:  "apostador_1",
			ContrasenaHash: "pass_123",
		})
		if err != nil {
			t.Fatalf("Setup Usuario fallo: %v", err)
		}

		gpH, err := queries.CrearGranPremioHistorico(ctx, sqlc.CrearGranPremioHistoricoParams{
			IDGranPremio: "SILVERSTONE",
			Pais:         "Reino Unido",
		})
		if err != nil {
			t.Fatalf("Setup GPHistorico fallo: %v", err)
		}

		fechaCarrera := time.Date(2026, 7, 5, 14, 0, 0, 0, time.UTC)
		_, err = queries.CrearGranPremioEvento(ctx, sqlc.CrearGranPremioEventoParams{
			IDGranPremio:     gpH.IDGranPremio,
			FechaCarrera:     fechaCarrera,
			ResultadoCarrera: []int32{},
		})
		if err != nil {
			t.Fatalf("Setup GPEvento fallo: %v", err)
		}

		// 1. Crear Apuesta
		apuesta, err := queries.CrearApuesta(ctx, sqlc.CrearApuestaParams{
			IDUsuario:    usr.IDUsuario,
			IDGranPremio: gpH.IDGranPremio,
			FechaCarrera: fechaCarrera,
			Prediccion:   []int32{1, 44, 16},
		})
		if err != nil {
			t.Fatalf("CrearApuesta fallo: %v", err)
		}

		if apuesta.IDApuesta == 0 {
			t.Error("Se esperaba un IDApuesta generado")
		}

		// 2. Listar Apuestas por Usuario
		apuestas, err := queries.ListarApuestasPorUsuario(ctx, usr.IDUsuario)
		if err != nil {
			t.Fatalf("ListarApuestasPorUsuario fallo: %v", err)
		}

		if len(apuestas) != 1 {
			t.Errorf("Se esperaba 1 apuesta, se obtuvieron %d", len(apuestas))
		}

		// 3. Borrar Apuesta
		err = queries.BorrarApuesta(ctx, apuesta.IDApuesta)
		if err != nil {
			t.Fatalf("BorrarApuesta fallo: %v", err)
		}

		// Verificar que fue eliminada
		apuestasDespues, err := queries.ListarApuestasPorUsuario(ctx, usr.IDUsuario)
		if err != nil {
			t.Fatalf("ListarApuestasPorUsuario fallo post borrado: %v", err)
		}

		if len(apuestasDespues) != 0 {
			t.Errorf("Se esperaban 0 apuestas despues de borrar, se obtuvieron %d", len(apuestasDespues))
		}
	})
}
