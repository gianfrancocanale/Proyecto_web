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

func CrearGranPremioHistorico(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var granPremioHistorico sqlc.CrearGranPremioHistoricoParams

	err := json.NewDecoder(r.Body).Decode(&granPremioHistorico)
	if err != nil {
		http.Error(w, "Error al decodificar el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}

	if granPremioHistorico.FechaCarrera.IsZero() {
		http.Error(w, "Todos los campos son obligatorios", http.StatusBadRequest)
		return
	}

	queries := sqlc.New(db)

	for _, piloto := range granPremioHistorico.ResultadoCarrera {
		_, err = queries.RecuperarPiloto(r.Context(), piloto)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.Error(w, "El piloto no existe en la grilla", http.StatusBadRequest)
			} else {
				http.Error(w, "Error al verificar resultado", http.StatusInternalServerError)
			}
			return
		}
	}

	_, err = queries.RecuperarGranPremio(r.Context(), granPremioHistorico.IDGranPremio)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "El gran premio no existe", http.StatusBadRequest)
		} else {
			http.Error(w, "Error al verificar el gran premio", http.StatusInternalServerError)
		}
		return
	}

	devolver, err := queries.CrearGranPremioHistorico(r.Context(), granPremioHistorico)
	if err != nil {
		http.Error(w, "Error al crear el Gran Premio Histórico", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(devolver)
}

func ListarGrandesPremiosHistoricos(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	queries := sqlc.New(db)

	granPremioHistorico, err := queries.ListarGrandesPremiosHistoricos(r.Context())
	if err != nil {
		http.Error(w, "Error al obtener los Gran Premios Históricos", http.StatusInternalServerError)
		return
	}
	if len(granPremioHistorico) == 0 {
		http.Error(w, "No se encontraron Gran Premios Histórico", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(granPremioHistorico)
}

func ObtenerGranPremioHistorico(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	idGranPremioHistorico := r.PathValue("id")

	queries := sqlc.New(db)

	granPremioHistorico, err := queries.RecuperarGranPremioHistorico(r.Context(), idGranPremioHistorico)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Gran Premio Histórico no encontrado", http.StatusNotFound)
		} else {
			http.Error(w, "Error al obtener el Gran Premio Histórico", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(granPremioHistorico)
}

func ActualizarGranPremioHistorico(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	idGranPremioHistorico := r.PathValue("id")
	fechaCarrera, err := time.Parse("02-01-2006", r.PathValue("fecha_carrera"))

	var granPremioHistorico sqlc.RegistrarResultadoCarreraParams

	err = json.NewDecoder(r.Body).Decode(&granPremioHistorico)
	if err != nil {
		http.Error(w, "Error al decodificar el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}

	if granPremioHistorico.FechaCarrera.IsZero() {
		http.Error(w, "Todos los campos son obligatorios", http.StatusBadRequest)
		return
	}

	granPremioHistorico.IDGranPremio = idGranPremioHistorico
	granPremioHistorico.FechaCarrera = fechaCarrera

	queries := sqlc.New(db)

	for _, piloto := range granPremioHistorico.ResultadoCarrera {
		_, err = queries.RecuperarPiloto(r.Context(), piloto)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.Error(w, "El piloto %s no existe en la grilla", http.StatusBadRequest)
			} else {
				http.Error(w, "Error al verificar resultado", http.StatusInternalServerError)
			}
			return
		}
	}

	_, err = queries.RecuperarGranPremio(r.Context(), granPremioHistorico.IDGranPremio)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "El gran premio no existe", http.StatusBadRequest)
		} else {
			http.Error(w, "Error al verificar el gran premio", http.StatusInternalServerError)
		}
		return
	}

	_, err = queries.RegistrarResultadoCarrera(r.Context(), granPremioHistorico)
	if err != nil {
		http.Error(w, "Error al actualizar el Gran Premio Histórico", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(granPremioHistorico)
}

func EliminarGranPremioHistorico(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	idGranPremioHistorico := r.PathValue("id")

	idGranPremioHistorico, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
	if err != nil {
		http.Error(w, "ID de piloto histórico inválido", http.StatusBadRequest)
		return
	}

	queries := sqlc.New(db)

	err := queries.EliminarGranPremioHistorico(r.Context(), int32(idGranPremioHistorico))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Gran Premio Histórico no encontrado", http.StatusNotFound)
		} else {
			http.Error(w, "Error al eliminar el Gran Premio Histórico", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}
