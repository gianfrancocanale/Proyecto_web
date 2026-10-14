package handlers

import (
	sqlc "Proyecto_web/db/sqlc"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
)

// GranPremioHandler agrupa los handlers de gran premio y guarda
// la conexión a la base de datos y las queries de sqlc.
type GranPremioHandler struct {
	DB      *sql.DB
	queries *sqlc.Queries
}

// NewGranPremioHandler es el constructor del handler.
func NewGranPremioHandler(db *sql.DB) *GranPremioHandler {
	return &GranPremioHandler{
		DB:      db,
		queries: sqlc.New(db),
	}
}

// HandlerGranPremio despacha según el método HTTP para la ruta /grandes-premios.
func (h *GranPremioHandler) HandlerGranPremio(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.CrearGranPremio(w, r)
	case http.MethodGet:
		h.ListarGranPremios(w, r)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

// HandlerGranPremioId despacha según el método HTTP para la ruta /grandes-premios/{id}.
func (h *GranPremioHandler) HandlerGranPremioId(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.ObtenerGranPremio(w, r)
	case http.MethodPut:
		h.ActualizarGranPremio(w, r)
	case http.MethodDelete:
		h.EliminarGranPremio(w, r)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

func (h *GranPremioHandler) CrearGranPremio(w http.ResponseWriter, r *http.Request) {
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

	if granPremio.IDGranPremio == "" || granPremio.Pais == "" || granPremio.CantidadVueltas <= 0 || granPremio.IDUltimoGanador == 0 {
		http.Error(w, "Todos los campos son obligatorios", http.StatusBadRequest)
		return
	}

	// Verificar que el último ganador existe
	_, err = h.queries.RecuperarPilotoHistorico(r.Context(), granPremio.IDUltimoGanador)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "El último ganador no existe", http.StatusBadRequest)
		} else {
			http.Error(w, "Error al verificar el último ganador", http.StatusInternalServerError)
		}
		return
	}

	devolver, err := h.queries.CrearGranPremio(r.Context(), granPremio)
	if err != nil {
		http.Error(w, "Error al crear el Gran Premio", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(devolver)
}

func (h *GranPremioHandler) ListarGranPremios(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	granPremios, err := h.queries.ListarGrandesPremios(r.Context())
	if err != nil {
		http.Error(w, "Error al obtener los Gran Premios", http.StatusInternalServerError)
		return
	}
	if len(granPremios) == 0 {
		http.Error(w, "No se encontraron Gran Premios", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(granPremios)
}

func (h *GranPremioHandler) ObtenerGranPremio(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	idGranPremio := r.PathValue("id")
	if idGranPremio == "" {
		http.Error(w, "ID de Gran Premio inválido", http.StatusBadRequest)
		return
	}

	granPremio, err := h.queries.RecuperarGranPremio(r.Context(), idGranPremio)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Gran Premio no encontrado", http.StatusNotFound)
		} else {
			http.Error(w, "Error al obtener el Gran Premio", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(granPremio)
}

func (h *GranPremioHandler) ActualizarGranPremio(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	idGranPremio := r.PathValue("id")
	if idGranPremio == "" {
		http.Error(w, "ID de Gran Premio inválido", http.StatusBadRequest)
		return
	}

	var granPremio sqlc.ModificarGranPremioParams

	err := json.NewDecoder(r.Body).Decode(&granPremio)
	if err != nil {
		http.Error(w, "Error al decodificar el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}

	// La identidad del Gran Premio sale siempre de la URL, no del body
	granPremio.IDGranPremio = idGranPremio

	if granPremio.Pais == "" || granPremio.LongitudKm == 0 || granPremio.CantidadVueltas < 0 || granPremio.IDUltimoGanador == 0 {
		http.Error(w, "Todos los campos son obligatorios", http.StatusBadRequest)
		return
	}

	// Verificar que el último ganador existe
	_, err = h.queries.RecuperarPilotoHistorico(r.Context(), granPremio.IDUltimoGanador)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "El último ganador no existe", http.StatusBadRequest)
		} else {
			http.Error(w, "Error al verificar el último ganador", http.StatusInternalServerError)
		}
		return
	}

	devolver, err := h.queries.ModificarGranPremio(r.Context(), granPremio)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Gran Premio no encontrado", http.StatusNotFound)
		} else {
			http.Error(w, "Error al actualizar el Gran Premio", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(devolver)
}

func (h *GranPremioHandler) EliminarGranPremio(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	idGranPremio := r.PathValue("id")
	if idGranPremio == "" {
		http.Error(w, "ID de Gran Premio inválido", http.StatusBadRequest)
		return
	}

	err := h.queries.EliminarGranPremio(r.Context(), idGranPremio)
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
