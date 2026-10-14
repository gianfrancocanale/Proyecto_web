package handlers

import (
	sqlc "Proyecto_web/db/sqlc"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

// PilotoHistoricoHandler agrupa los handlers de piloto histórico y guarda
// la conexión a la base de datos y las queries de sqlc.
type PilotoHistoricoHandler struct {
	DB      *sql.DB
	queries *sqlc.Queries
}

// NewPilotoHistoricoHandler es el constructor del handler.
func NewPilotoHistoricoHandler(db *sql.DB) *PilotoHistoricoHandler {
	return &PilotoHistoricoHandler{
		DB:      db,
		queries: sqlc.New(db),
	}
}

// HandlerPilotoHistorico despacha según el método HTTP para la ruta /pilotos-historicos.
func (h *PilotoHistoricoHandler) HandlerPilotoHistorico(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.CrearPilotoHistorico(w, r)
	case http.MethodGet:
		h.ListarPilotosHistoricos(w, r)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

// HandlerPilotoHistoricoId despacha según el método HTTP para la ruta /pilotos-historicos/{id}.
func (h *PilotoHistoricoHandler) HandlerPilotoHistoricoId(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.ObtenerPilotoHistorico(w, r)
	case http.MethodPut:
		h.ActualizarPilotoHistorico(w, r)
	case http.MethodDelete:
		h.EliminarPilotoHistorico(w, r)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

func (h *PilotoHistoricoHandler) CrearPilotoHistorico(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var pilotoHistorico sqlc.CrearPilotoHistoricoParams

	err := json.NewDecoder(r.Body).Decode(&pilotoHistorico)
	if err != nil {
		http.Error(w, "Error al decodificar el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}

	pilotoHistorico.TitulosGanados = 0 // setea 0 títulos ganados al crear un nuevo piloto histórico

	if pilotoHistorico.Nombre == "" || pilotoHistorico.Pais == "" || pilotoHistorico.FechaNacimiento.IsZero() {
		http.Error(w, "Todos los campos son obligatorios", http.StatusBadRequest)
		return
	}

	devolver, err := h.queries.CrearPilotoHistorico(r.Context(), pilotoHistorico)
	if err != nil {
		http.Error(w, "Error al crear el piloto histórico", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(devolver)
}

func (h *PilotoHistoricoHandler) ListarPilotosHistoricos(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	pilotosHistoricos, err := h.queries.ListarPilotosHistoricos(r.Context())
	if err != nil {
		http.Error(w, "Error al obtener los pilotos históricos", http.StatusInternalServerError)
		return
	}
	if len(pilotosHistoricos) == 0 {
		http.Error(w, "No se encontraron pilotos históricos", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pilotosHistoricos)
}

func (h *PilotoHistoricoHandler) ObtenerPilotoHistorico(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	idPiloto, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
	if err != nil {
		http.Error(w, "ID de piloto histórico inválido", http.StatusBadRequest)
		return
	}

	pilotoHistorico, err := h.queries.RecuperarPilotoHistorico(r.Context(), int32(idPiloto))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Piloto histórico no encontrado", http.StatusNotFound)
		} else {
			http.Error(w, "Error al obtener el piloto histórico", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pilotoHistorico)
}

func (h *PilotoHistoricoHandler) ActualizarPilotoHistorico(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	idPiloto, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
	if err != nil {
		http.Error(w, "ID de piloto histórico inválido", http.StatusBadRequest)
		return
	}

	var pilotoHistorico sqlc.ModificarPilotoHistoricoParams

	err = json.NewDecoder(r.Body).Decode(&pilotoHistorico)
	if err != nil {
		http.Error(w, "Error al decodificar el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}

	// La identidad del piloto histórico sale siempre de la URL, no del body
	pilotoHistorico.IDPiloto = int32(idPiloto)

	if pilotoHistorico.Nombre == "" || pilotoHistorico.Pais == "" || pilotoHistorico.FechaNacimiento.IsZero() {
		http.Error(w, "Todos los campos son obligatorios", http.StatusBadRequest)
		return
	}

	devolver, err := h.queries.ModificarPilotoHistorico(r.Context(), pilotoHistorico)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Piloto histórico no encontrado", http.StatusNotFound)
		} else {
			http.Error(w, "Error al actualizar el piloto histórico", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(devolver)
}

func (h *PilotoHistoricoHandler) EliminarPilotoHistorico(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	idPiloto, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
	if err != nil {
		http.Error(w, "ID de piloto histórico inválido", http.StatusBadRequest)
		return
	}

	err = h.queries.EliminarPilotoHistorico(r.Context(), int32(idPiloto))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Piloto histórico no encontrado", http.StatusNotFound)
		} else {
			http.Error(w, "Error al eliminar el piloto histórico", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}
