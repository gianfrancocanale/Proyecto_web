package handlers

import (
	sqlc "Proyecto_web/db/sqlc"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

// GranPremioHistoricoHandler agrupa los handlers de gran premio histórico y guarda
// la conexión a la base de datos y las queries de sqlc.
type GranPremioHistoricoHandler struct {
	DB      *sql.DB
	queries *sqlc.Queries
}

// NewGranPremioHistoricoHandler es el constructor del handler.
func NewGranPremioHistoricoHandler(db *sql.DB) *GranPremioHistoricoHandler {
	return &GranPremioHistoricoHandler{
		DB:      db,
		queries: sqlc.New(db),
	}
}

// HandlerGranPremioHistorico despacha según el método HTTP para la ruta /grandes-premios-historicos.
func (h *GranPremioHistoricoHandler) HandlerGranPremioHistorico(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.CrearGranPremioHistorico(w, r)
	case http.MethodGet:
		h.ListarGrandesPremiosHistoricos(w, r)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

// HandlerGranPremioHistoricoId despacha según el método HTTP para la ruta
// /grandes-premios-historicos/{id}/{fecha_carrera}.
func (h *GranPremioHistoricoHandler) HandlerGranPremioHistoricoId(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.ObtenerGranPremioHistorico(w, r)
	case http.MethodPut:
		h.ActualizarGranPremioHistorico(w, r)
	case http.MethodDelete:
		h.EliminarGranPremioHistorico(w, r)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

func (h *GranPremioHistoricoHandler) CrearGranPremioHistorico(w http.ResponseWriter, r *http.Request) {
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

	// Verificar que el gran premio existe
	_, err = h.queries.RecuperarGranPremio(r.Context(), granPremioHistorico.IDGranPremio)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "El gran premio no existe", http.StatusBadRequest)
		} else {
			http.Error(w, "Error al verificar el gran premio", http.StatusInternalServerError)
		}
		return
	}

	// Verificar que todos los pilotos del resultado existen
	for _, piloto := range granPremioHistorico.ResultadoCarrera {
		_, err = h.queries.RecuperarPilotoHistorico(r.Context(), piloto)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.Error(w, "El piloto no existe", http.StatusBadRequest)
			} else {
				http.Error(w, "Error al verificar resultado", http.StatusInternalServerError)
			}
			return
		}
	}

	devolver, err := h.queries.CrearGranPremioHistorico(r.Context(), granPremioHistorico)
	if err != nil {
		http.Error(w, "Error al crear el Gran Premio Histórico", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(devolver)
}

func (h *GranPremioHistoricoHandler) ListarGrandesPremiosHistoricos(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	grandesPremiosHistoricos, err := h.queries.ListarGrandesPremiosHistoricos(r.Context())
	if err != nil {
		http.Error(w, "Error al obtener los Gran Premios Históricos", http.StatusInternalServerError)
		return
	}
	if len(grandesPremiosHistoricos) == 0 {
		http.Error(w, "No se encontraron Gran Premios Históricos", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(grandesPremiosHistoricos)
}

func (h *GranPremioHistoricoHandler) ObtenerGranPremioHistorico(w http.ResponseWriter, r *http.Request) {
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

	granPremioHistorico, err := h.queries.RecuperarGranPremioHistorico(r.Context(), sqlc.RecuperarGranPremioHistoricoParams{
		IDGranPremio: idGranPremioHistorico,
		FechaCarrera: fechaCarrera,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Gran Premio Histórico no encontrado", http.StatusNotFound)
		} else {
			http.Error(w, "Error al obtener el Gran Premio Histórico", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(granPremioHistorico)
}

func (h *GranPremioHistoricoHandler) ActualizarGranPremioHistorico(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, "El ID del Gran Premio y la fecha de carrera son obligatorios", http.StatusBadRequest)
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

	// La identidad del gran premio histórico sale siempre de la URL, no del body
	granPremioHistorico.IDGranPremio = idGranPremioHistorico
	granPremioHistorico.FechaCarrera = fechaCarrera

	// Verificar que el gran premio existe
	_, err = h.queries.RecuperarGranPremio(r.Context(), idGranPremioHistorico)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "El gran premio no existe", http.StatusBadRequest)
		} else {
			http.Error(w, "Error al verificar el gran premio", http.StatusInternalServerError)
		}
		return
	}

	// Verificar que todos los pilotos del resultado existen
	for _, piloto := range granPremioHistorico.ResultadoCarrera {
		_, err = h.queries.RecuperarPilotoHistorico(r.Context(), piloto)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.Error(w, "Uno de los pilotos no existe en la grilla", http.StatusBadRequest)
			} else {
				http.Error(w, "Error al verificar resultado", http.StatusInternalServerError)
			}
			return
		}
	}

	devolver, err := h.queries.RegistrarResultadoCarrera(r.Context(), granPremioHistorico)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Gran Premio Histórico no encontrado", http.StatusNotFound)
		} else {
			http.Error(w, "Error al actualizar el Gran Premio Histórico", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(devolver)
}

func (h *GranPremioHistoricoHandler) EliminarGranPremioHistorico(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	idGranPremioHistorico := r.PathValue("id")
	if idGranPremioHistorico == "" {
		http.Error(w, "El ID del Gran Premio es obligatorio", http.StatusBadRequest)
		return
	}

	fechaCarrera, err := time.Parse("02-01-2006", r.PathValue("fecha_carrera"))
	if err != nil {
		http.Error(w, "Fecha de carrera inválida", http.StatusBadRequest)
		return
	}

	err = h.queries.EliminarGranPremioHistorico(r.Context(), sqlc.EliminarGranPremioHistoricoParams{
		IDGranPremio: idGranPremioHistorico,
		FechaCarrera: fechaCarrera,
	})
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
