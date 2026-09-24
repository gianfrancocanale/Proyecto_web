package handlers

import (
	sqlc "Proyecto_web/db/sqlc"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"
)

func CrearApuesta(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var apuesta sqlc.CrearApuestaParams
	err := json.NewDecoder(r.Body).Decode(&apuesta)
	if err != nil {
		http.Error(w, "Error al decodificar el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}
	if apuesta.IDUsuario == 0 || apuesta.IDGranPremio == "" || len(apuesta.Prediccion) == 0 || apuesta.FechaCarrera.IsZero() || apuesta.FechaApuesta.IsZero() {
		http.Error(w, "Todos los campos son obligatorios", http.StatusBadRequest)
		return
	}

	// Asignar la fecha actual como fecha de carrera
	apuesta.FechaApuesta = time.Now()

	queries := sqlc.New(db)

	// Verificar que el usuario existe
	_, err = queries.RecuperarUsuarioPorId(r.Context(), apuesta.IDUsuario)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "No existe el usuario", http.StatusBadRequest)
		} else {
			http.Error(w, "Error al verificar el usuario", http.StatusInternalServerError)
		}
		return
	}

	// Verificar que el gran premio existe
	_, err = queries.RecuperarGranPremioHistorico(r.Context(), apuesta.IDGranPremio, apuesta.FechaCarrera)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "No existe el gran premio", http.StatusBadRequest)
		} else {
			http.Error(w, "Error al verificar el gran premio", http.StatusInternalServerError)
		}
		return
	}

	// Verificar que todos los pilotos en la predicción existen
	for _, pilotoID := range apuesta.Prediccion {
		_, err = queries.RecuperarPiloto(r.Context(), pilotoID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.Error(w, "No existe el piloto en la predicción", http.StatusBadRequest)
			} else {
				http.Error(w, "Error al verificar el piloto en la predicción", http.StatusInternalServerError)
			}
			return
		}
	}

	devolver, err := queries.CrearApuesta(r.Context(), apuesta)
	if err != nil {
		http.Error(w, "Error al crear la apuesta", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(devolver)
}

func ListarApuestas(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	queries := sqlc.New(db)

	apuestas, err := queries.ListarApuestas(r.Context())
	if err != nil {
		http.Error(w, "Error al obtener las apuestas", http.StatusInternalServerError)
		return
	}
	if len(apuestas) == 0 {
		http.Error(w, "No se encontraron apuestas", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(apuestas)
}

func ObtenerApuesta(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	idUsuario, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
	if err != nil {
		http.Error(w, "ID de apuesta invalido", http.StatusBadRequest)
		return
	}

	queries := sqlc.New(db)

	apuesta, err := queries.ListarApuestasPorUsuario(r.Context(), int32(idUsuario))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Apuesta no encontrada", http.StatusNotFound)
		} else {
			http.Error(w, "Error al obtener la apuesta", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(apuesta)
}

func ActualizarApuesta(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}
	idUsuario, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
	if err != nil {
		http.Error(w, "ID de usuario invalido", http.StatusBadRequest)
		return
	}

	var apuesta sqlc.ModificarApuestaParams

	err := json.NewDecoder(r.Body).Decode(&apuesta)
	if err != nil {
		http.Error(w, "Error al decodificar el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}

	if apuesta.IDApuesta == 0 || len(apuesta.Prediccion) == 0 || apuesta.FechaCarrera.IsZero() || apuesta.FechaApuesta.IsZero() {
		http.Error(w, "Todos los campos son obligatorios", http.StatusBadRequest)
		return
	}

	queries := sqlc.New(db)

	devolver, err := queries.ModificarApuesta(r.Context(), apuesta)
	if err != nil {
		http.Error(w, "Error al actualizar la apuesta", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(devolver)
}
