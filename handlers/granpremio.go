package handlers

import (
	sqlc "Proyecto_web/db/sqlc"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
)

func CrearGranPremio(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var granPremio sqlc.CrearGranPremioParams

	err := json.NewDecoder(r.Body).Decode(&granPremio)
	if err != nil {
		http.Error(w, "Error al decodificar el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}

	if granPremio.IDGranPremio == "" || granPremio.Pais == "" || granPremio.CantidadVueltas == 0 || granPremio.CantidadVueltas < 0 || granPremio.IDUltimoGanador == 0 {
		http.Error(w, "Todos los campos son obligatorios", http.StatusBadRequest)
		return
	}

	queries := sqlc.New(db)

	_, err = queries.RecuperarPiloto(r.Context(), granPremio.IDUltimoGanador)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "El último ganador no existe", http.StatusBadRequest)
		} else {
			http.Error(w, "Error al verificar el último ganador", http.StatusInternalServerError)
		}
		return
	}

	devolver, err := queries.CrearGranPremio(r.Context(), granPremio)
	if err != nil {
		http.Error(w, "Error al crear el Gran Premio", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(devolver)
}

func ListarGranPremios(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	queries := sqlc.New(db)

	granPremios, err := queries.ListarGrandesPremios(r.Context())
	if err != nil {
		http.Error(w, "Error al obtener los Gran Premios", http.StatusInternalServerError)
		return
	}
	if len(granPremios) == 0 {
		http.Error(w, "No se encontraron Gran Premios", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(granPremios)
}

func ObtenerGranPremio(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	idGranPremio := r.PathValue("id")

	queries := sqlc.New(db)

	granPremio, err := queries.RecuperarGranPremio(r.Context(), idGranPremio)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Gran Premio no encontrado", http.StatusNotFound)
		} else {
			http.Error(w, "Error al obtener el Gran Premio", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(granPremio)
}

func ActualizarGranPremio(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	idGranPremio := r.PathValue("id")

	var granPremio sqlc.ModificarGranPremioParams

	err := json.NewDecoder(r.Body).Decode(&granPremio)
	if err != nil {
		http.Error(w, "Error al decodificar el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}

	if granPremio.Pais == "" || granPremio.LongitudKm == 0 || granPremio.CantidadVueltas < 0 || granPremio.IDUltimoGanador == 0 {
		http.Error(w, "Todos los campos son obligatorios", http.StatusBadRequest)
		return
	}

	granPremio.IDGranPremio = idGranPremio

	queries := sqlc.New(db)

	_, err = queries.RecuperarPiloto(r.Context(), granPremio.IDUltimoGanador)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "El último ganador no existe", http.StatusBadRequest)
		} else {
			http.Error(w, "Error al verificar el último ganador", http.StatusInternalServerError)
		}
		return
	}

	_, err = queries.ModificarGranPremio(r.Context(), granPremio)
	if err != nil {
		http.Error(w, "Error al actualizar el Gran Premio", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(granPremio)
}

func EliminarGranPremio(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	idGranPremio := r.PathValue("id")

	queries := sqlc.New(db)

	err := queries.EliminarGranPremio(r.Context(), idGranPremio)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Gran Premio no encontrado", http.StatusNotFound)
		} else {
			http.Error(w, "Error al eliminar el Gran Premio", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}
