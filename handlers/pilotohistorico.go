package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"Proyecto_web/aplicacion"
	"Proyecto_web/dominio"
)

type PilotoHistoricoHandler struct {
	pilotos *aplicacion.PilotoHistoricoService
}

func NewPilotoHistoricoHandler(s *aplicacion.PilotoHistoricoService) *PilotoHistoricoHandler {
	return &PilotoHistoricoHandler{pilotos: s}
}

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

	var p dominio.PilotoHistorico
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Error al decodificar el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}
	p.Titulos = 0 // los pilotos nuevos arrancan sin títulos
	if err := h.pilotos.Crear(r.Context(), p); err != nil {
		escribirError(w, err, "Error al crear el piloto histórico")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(p)
}

func (h *PilotoHistoricoHandler) ListarPilotosHistoricos(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	pilotos, err := h.pilotos.Listar(r.Context())
	if err != nil {
		http.Error(w, "Error al obtener los pilotos históricos", http.StatusInternalServerError)
		return
	}
	if len(pilotos) == 0 {
		http.Error(w, "No se encontraron pilotos históricos", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pilotos)
}

func (h *PilotoHistoricoHandler) ObtenerPilotoHistorico(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
	if err != nil {
		http.Error(w, "ID de piloto histórico inválido", http.StatusBadRequest)
		return
	}
	p, err := h.pilotos.Recuperar(r.Context(), int(id))
	if err != nil {
		escribirError(w, err, "Error al obtener el piloto histórico")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

func (h *PilotoHistoricoHandler) ActualizarPilotoHistorico(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
	if err != nil {
		http.Error(w, "ID de piloto histórico inválido", http.StatusBadRequest)
		return
	}
	var p dominio.PilotoHistorico
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Error al decodificar el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}
	// La identidad sale de la URL, nunca del body.
	p.ID = int(id)

	if err := h.pilotos.Modificar(r.Context(), p); err != nil {
		escribirError(w, err, "Error al actualizar el piloto histórico")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

func (h *PilotoHistoricoHandler) EliminarPilotoHistorico(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
	if err != nil {
		http.Error(w, "ID de piloto histórico inválido", http.StatusBadRequest)
		return
	}
	if err := h.pilotos.Eliminar(r.Context(), int(id)); err != nil {
		escribirError(w, err, "Error al eliminar el piloto histórico")
		return
	}
	w.WriteHeader(http.StatusOK)
}

// escribirError traduce errores de dominio a respuestas HTTP.
func escribirError(w http.ResponseWriter, err error, mensajeGenerico string) {
	var ev dominio.ErrValidacion
	switch {
	case errors.As(err, &ev):
		http.Error(w, ev.Error(), http.StatusBadRequest)
	case errors.Is(err, dominio.ErrPilotoHistoricoNoEncontrado),
		errors.Is(err, dominio.ErrPilotoNoEncontrado),
		errors.Is(err, dominio.ErrCarreraNoEncontrada),
		errors.Is(err, dominio.ErrApuestaNoEncontrada),
		errors.Is(err, dominio.ErrUsuarioNoEncontrado):
		http.Error(w, err.Error(), http.StatusNotFound)
	case errors.Is(err, dominio.ErrPilotoHistoricoDuplicado),
		errors.Is(err, dominio.ErrPilotoDuplicado),
		errors.Is(err, dominio.ErrApuestaDuplicada),
		errors.Is(err, dominio.ErrEdicionDuplicada):
		http.Error(w, err.Error(), http.StatusConflict)
	default:
		http.Error(w, mensajeGenerico, http.StatusInternalServerError)
	}
}
