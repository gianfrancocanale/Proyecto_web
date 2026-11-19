// handlers/granpremiohistorico.go
package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"Proyecto_web/aplicacion"
	"Proyecto_web/dominio"
)

type GranPremioHistoricoHandler struct {
	carreras *aplicacion.GranPremioHistoricoService
}

func NewGranPremioHistoricoHandler(s *aplicacion.GranPremioHistoricoService) *GranPremioHistoricoHandler {
	return &GranPremioHistoricoHandler{carreras: s}
}

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

	var g dominio.GranPremioHistorico
	if err := json.NewDecoder(r.Body).Decode(&g); err != nil {
		http.Error(w, "Error al decodificar el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}
	if err := h.carreras.Crear(r.Context(), g); err != nil {
		escribirError(w, err, "Error al crear el Gran Premio Histórico")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(g)
}

func (h *GranPremioHistoricoHandler) ListarGrandesPremiosHistoricos(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	carreras, err := h.carreras.Listar(r.Context())
	if err != nil {
		http.Error(w, "Error al obtener los Gran Premios Históricos", http.StatusInternalServerError)
		return
	}
	if len(carreras) == 0 {
		http.Error(w, "No se encontraron Gran Premios Históricos", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(carreras)
}

func (h *GranPremioHistoricoHandler) ObtenerGranPremioHistorico(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	id := r.PathValue("id")
	fecha, err := time.Parse("02-01-2006", r.PathValue("fecha_carrera"))
	if id == "" || err != nil {
		http.Error(w, "Son obligatorios el ID del Gran Premio y su fecha", http.StatusBadRequest)
		return
	}
	g, err := h.carreras.Recuperar(r.Context(), id, fecha)
	if err != nil {
		escribirError(w, err, "Error al obtener el Gran Premio Histórico")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(g)
}

func (h *GranPremioHistoricoHandler) ActualizarGranPremioHistorico(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	id := r.PathValue("id")
	fecha, err := time.Parse("02-01-2006", r.PathValue("fecha_carrera"))
	if id == "" || err != nil {
		http.Error(w, "El ID del Gran Premio y la fecha de carrera son obligatorios", http.StatusBadRequest)
		return
	}
	var g dominio.GranPremioHistorico
	if err := json.NewDecoder(r.Body).Decode(&g); err != nil {
		http.Error(w, "Error al decodificar el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}
	// La identidad sale de la URL, nunca del body.
	g.ID = id
	g.Fecha = fecha

	if err := h.carreras.RegistrarResultado(r.Context(), g); err != nil {
		escribirError(w, err, "Error al actualizar el Gran Premio Histórico")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(g)
}

func (h *GranPremioHistoricoHandler) EliminarGranPremioHistorico(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	id := r.PathValue("id")
	fecha, err := time.Parse("02-01-2006", r.PathValue("fecha_carrera"))
	if id == "" || err != nil {
		http.Error(w, "El ID del Gran Premio y la fecha de carrera son obligatorios", http.StatusBadRequest)
		return
	}
	if err := h.carreras.Eliminar(r.Context(), id, fecha); err != nil {
		escribirError(w, err, "Error al eliminar el Gran Premio Histórico")
		return
	}
	w.WriteHeader(http.StatusOK)
}
