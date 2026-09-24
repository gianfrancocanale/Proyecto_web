package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
)

func crearEscuderia(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var escuderia sqlc.CrearEscuderiaParams

	err := json.NewDecoder(r.Body).Decode(&escuderia)

	if err != nil {
		http.Error(w, "Error al decodificar el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}

	if escuderia.IDEscuderia == "" || escuderia.TeamPrincipal == "" || escuderia.FechaFundacion.IsZero() {
		http.Error(w, "Todos los campos son obligatorios", http.StatusBadRequest)
		return
	}

	escuderia.PuntosTemporada = 0      //setea 0 puntos al crear una nueva escuderia
	escuderia.TitulosConstructores = 0 //setea 0 titulos al crear una nueva escuderia

	queries := sqlc.New(db)

	devolver, err := queries.CrearEscuderia(r.Context(), escuderia)
	if err != nil {
		http.Error(w, "Error al crear la escudería", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(devolver)
}

func listarEscuderias(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	queries := sqlc.New(db)

	escuderias, err := queries.ListarEscuderias(r.Context())
	if err != nil {
		http.Error(w, "Error al obtener las escuderías", http.StatusInternalServerError)
		return
	}
	if len(escuderias) == 0 {
		http.Error(w, "No se encontraron escuderías", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(escuderias)
}

func obtenerEscuderia(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	idEscuderia := r.PathValue("id")

	queries := sqlc.New(db)

	escuderia, err := queries.RecuperarEscuderia(r.Context(), idEscuderia)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Escudería no encontrada", http.StatusNotFound)
		} else {
			http.Error(w, "Error al obtener la escudería", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(escuderia)
}

func actualizarEscuderia(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	idEscuderia := r.PathValue("id")

	var escuderia sqlc.ModificarEscuderiaParams

	err := json.NewDecoder(r.Body).Decode(&escuderia)
	if err != nil {
		http.Error(w, "Error al decodificar el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}

	if escuderia.TeamPrincipal == "" || escuderia.FechaFundacion.IsZero() {
		http.Error(w, "Todos los campos son obligatorios", http.StatusBadRequest)
		return
	}

	escuderia.IDEscuderia = idEscuderia

	queries := sqlc.New(db)

	devolver, err := queries.ModificarEscuderia(r.Context(), escuderia)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Escudería no encontrada", http.StatusNotFound)
		} else {
			http.Error(w, "Error al actualizar la escudería", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(devolver)
}

func eliminarEscuderia(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	idEscuderia := r.PathValue("id")

	queries := sqlc.New(db)

	err := queries.EliminarEscuderia(r.Context(), idEscuderia)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Escudería no encontrada", http.StatusNotFound)
		} else {
			http.Error(w, "Error al eliminar la escudería", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
