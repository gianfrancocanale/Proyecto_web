package handlers

import (
	sqlc "Proyecto_web/db/sqlc"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
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

	if granPremioHistorico.FechaCarrera.IsZero() || granPremioHistorico.IDGranPremio == "" || len(granPremioHistorico.ResultadoCarrera) == 0 {
		http.Error(w, "Todos los campos son obligatorios", http.StatusBadRequest)
		return
	}

	queries := sqlc.New(db)

	_, err = queries.RecuperarGranPremio(r.Context(), granPremioHistorico.IDGranPremio)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "El gran premio no existe", http.StatusBadRequest)
		} else {
			http.Error(w, "Error al verificar el gran premio", http.StatusInternalServerError)
		}
		return
	}

	for _, piloto := range granPremioHistorico.ResultadoCarrera {
		_, err = queries.RecuperarPiloto(r.Context(), piloto)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.Error(w, "Uno de los pilotos no existe en la grilla", http.StatusBadRequest)
			} else {
				http.Error(w, "Error al verificar resultado", http.StatusInternalServerError)
			}
			return
		}
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
	fechaCarreraStr := r.PathValue("fecha_carrera")

	if idGranPremioHistorico == "" || fechaCarreraStr == "" {
		http.Error(w, "Son obligatorios el ID del Gran Premio y su fecha", http.StatusBadRequest)
		return
	}

	fechaCarrera, err := time.Parse("02-01-2006", fechaCarreraStr)
	if err != nil {
		http.Error(w, "Fecha de carrera inválida", http.StatusBadRequest)
		return
	}

	queries := sqlc.New(db)

	granPremioHistorico, err := queries.RecuperarGranPremioHistorico(r.Context(), sqlc.RecuperarGranPremioHistoricoParams{
		IDGranPremio: idGranPremioHistorico,
		FechaCarrera: fechaCarrera})

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
	if err != nil {
		http.Error(w, "Fecha de carrera inválida", http.StatusBadRequest)
		return
	}

	if idGranPremioHistorico == "" || fechaCarrera.IsZero() {
		http.Error(
			w,
			"El ID del Gran Premio y la fecha de carrera son obligatorios",
			http.StatusBadRequest,
		)
		return
	}

	var granPremioHistorico sqlc.RegistrarResultadoCarreraParams

	err = json.NewDecoder(r.Body).Decode(&granPremioHistorico)
	if err != nil {
		http.Error(w, "Error al decodificar el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}

	if len(granPremioHistorico.ResultadoCarrera) == 0 {
		http.Error(w, "El resultado de la carrera es obligatorio", http.StatusBadRequest)
		return
	}

	granPremioHistorico.IDGranPremio = idGranPremioHistorico
	granPremioHistorico.FechaCarrera = fechaCarrera

	queries := sqlc.New(db)

	_, err = queries.RecuperarGranPremio(r.Context(), idGranPremioHistorico)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "El gran premio no existe", http.StatusBadRequest)
		} else {
			http.Error(w, "Error al verificar el gran premio", http.StatusInternalServerError)
		}
		return
	}

	for _, piloto := range granPremioHistorico.ResultadoCarrera {
		_, err = queries.RecuperarPiloto(r.Context(), piloto)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.Error(w, "Uno de los pilotos no existe en la grilla", http.StatusBadRequest)
			} else {
				http.Error(w, "Error al verificar resultado", http.StatusInternalServerError)
			}
			return
		}
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
	fechaCarreraStr := r.PathValue("fecha_carrera")

	fechaCarrera, err := time.Parse("02-01-2006", fechaCarreraStr)
	if err != nil {
		http.Error(w, "Fecha de carrera inválida", http.StatusBadRequest)
		return
	}

	queries := sqlc.New(db)

	err = queries.EliminarGranPremioHistorico(r.Context(), sqlc.EliminarGranPremioHistoricoParams{
		IDGranPremio: idGranPremioHistorico,
		FechaCarrera: fechaCarrera})
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
