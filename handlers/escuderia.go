package handlers

import (
	sqlc "Proyecto_web/db/sqlc"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
)

// EscuderiaHandler agrupa los handlers de escudería y guarda
// la conexión a la base de datos y las queries de sqlc.
type EscuderiaHandler struct {
	DB      *sql.DB
	queries *sqlc.Queries
}

// NewEscuderiaHandler es el constructor del handler.
func NewEscuderiaHandler(db *sql.DB) *EscuderiaHandler {
	return &EscuderiaHandler{
		DB:      db,
		queries: sqlc.New(db),
	}
}

// HandlerEscuderia despacha según el método HTTP para la ruta /escuderias.
func (h *EscuderiaHandler) HandlerEscuderia(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.CrearEscuderia(w, r)
	case http.MethodGet:
		h.ListarEscuderias(w, r)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

// HandlerEscuderiaId despacha según el método HTTP para la ruta /escuderias/{id}.
func (h *EscuderiaHandler) HandlerEscuderiaId(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.ObtenerEscuderia(w, r)
	case http.MethodPut:
		h.ActualizarEscuderia(w, r)
	case http.MethodDelete:
		h.EliminarEscuderia(w, r)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

func (h *EscuderiaHandler) CrearEscuderia(w http.ResponseWriter, r *http.Request) {
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

	devolver, err := h.queries.CrearEscuderia(r.Context(), escuderia)
	if err != nil {
		http.Error(w, "Error al crear la escudería", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(devolver)
}

func (h *EscuderiaHandler) ListarEscuderias(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	escuderias, err := h.queries.ListarEscuderias(r.Context())
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

func (h *EscuderiaHandler) ObtenerEscuderia(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	idEscuderia := r.PathValue("id")

	escuderia, err := h.queries.RecuperarEscuderia(r.Context(), idEscuderia)
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

func (h *EscuderiaHandler) ActualizarEscuderia(w http.ResponseWriter, r *http.Request) {
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

	devolver, err := h.queries.ModificarEscuderia(r.Context(), escuderia)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Escudería no encontrada", http.StatusNotFound)
		} else {
			http.Error(w, "Error al actualizar la escudería", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(devolver)
}

func (h *EscuderiaHandler) EliminarEscuderia(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	idEscuderia := r.PathValue("id")

	err := h.queries.EliminarEscuderia(r.Context(), idEscuderia)
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
