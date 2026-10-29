package handlers

import (
	sqlc "Proyecto_web/db/sqlc"
	vista "Proyecto_web/views"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

// PilotoHandler agrupa los handlers de piloto y guarda
// la conexión a la base de datos y las queries de sqlc.
type PilotoHandler struct {
	DB      *sql.DB
	queries *sqlc.Queries
}

// NewPilotoHandler es el constructor del handler.
func NewPilotoHandler(db *sql.DB) *PilotoHandler {
	return &PilotoHandler{
		DB:      db,
		queries: sqlc.New(db),
	}
}

// HandlerPiloto despacha según el método HTTP para la ruta /pilotos.
func (h *PilotoHandler) HandlerPiloto(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.CrearPiloto(w, r)
	case http.MethodGet:
		h.ListarPilotos(w, r)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

// HandlerPilotoId despacha según el método HTTP para la ruta /pilotos/{id}.
func (h *PilotoHandler) HandlerPilotoId(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.ObtenerPiloto(w, r)
	case http.MethodPut:
		h.ActualizarPiloto(w, r)
	case http.MethodDelete:
		h.EliminarPiloto(w, r)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

func (h *PilotoHandler) Pagina(w http.ResponseWriter, r *http.Request) {
	pilotos, err := h.queries.ListarPilotos(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	vista.PilotosPage(pilotos).Render(r.Context(), w)
}

func (h *PilotoHandler) CrearPiloto(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var piloto sqlc.CrearPilotoTemporadaParams

	err := json.NewDecoder(r.Body).Decode(&piloto)
	if err != nil {
		http.Error(w, "Error al decodificar el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}

	piloto.PuntosTemporada = 0 // el piloto arranca con 0 puntos

	// Verificar que el piloto histórico existe
	_, err = h.queries.RecuperarPilotoHistorico(r.Context(), piloto.IDPiloto)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "No existe el piloto", http.StatusBadRequest)
		} else {
			http.Error(w, "Error al verificar el piloto", http.StatusInternalServerError)
		}
		return
	}

	// Verificar que la escudería existe
	_, err = h.queries.RecuperarEscuderia(r.Context(), piloto.IDEscuderia)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "No existe la escudería", http.StatusBadRequest)
		} else {
			http.Error(w, "Error al verificar la escudería", http.StatusInternalServerError)
		}
		return
	}

	devolver, err := h.queries.CrearPilotoTemporada(r.Context(), piloto)
	if err != nil {
		http.Error(w, "Error al crear el piloto", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(devolver)
}

func (h *PilotoHandler) ListarPilotos(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	pilotos, err := h.queries.ListarPilotos(r.Context())
	if err != nil {
		http.Error(w, "Error al obtener los pilotos", http.StatusInternalServerError)
		return
	}
	if len(pilotos) == 0 {
		http.Error(w, "No se encontraron pilotos", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pilotos)
}

func (h *PilotoHandler) ObtenerPiloto(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	idPiloto, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
	if err != nil {
		http.Error(w, "ID de piloto inválido", http.StatusBadRequest)
		return
	}

	piloto, err := h.queries.RecuperarPiloto(r.Context(), int32(idPiloto))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Piloto no encontrado", http.StatusNotFound)
		} else {
			http.Error(w, "Error al obtener el piloto", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(piloto)
}

func (h *PilotoHandler) ActualizarPiloto(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	idPiloto, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
	if err != nil {
		http.Error(w, "ID de piloto inválido", http.StatusBadRequest)
		return
	}

	var piloto sqlc.ModificarPilotoParams

	err = json.NewDecoder(r.Body).Decode(&piloto)
	if err != nil {
		http.Error(w, "Error al decodificar el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}

	// La identidad del piloto sale siempre de la URL, no del body
	piloto.IDPiloto = int32(idPiloto)

	devolver, err := h.queries.ModificarPiloto(r.Context(), piloto)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Piloto no encontrado", http.StatusNotFound)
		} else {
			http.Error(w, "Error al actualizar el piloto", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(devolver)
}

func (h *PilotoHandler) EliminarPiloto(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	idPiloto, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
	if err != nil {
		http.Error(w, "ID de piloto inválido", http.StatusBadRequest)
		return
	}

	err = h.queries.EliminarPilotoTemporada(r.Context(), int32(idPiloto))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Piloto no encontrado", http.StatusNotFound)
		} else {
			http.Error(w, "Error al eliminar el piloto", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}
