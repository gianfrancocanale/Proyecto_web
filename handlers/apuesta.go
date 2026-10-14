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

// ApuestaHandler agrupa los handlers de apuesta y guarda
// la conexión a la base de datos y las queries de sqlc.
type ApuestaHandler struct {
	DB      *sql.DB
	queries *sqlc.Queries
}

// NewApuestaHandler es el constructor del handler.
func NewApuestaHandler(db *sql.DB) *ApuestaHandler {
	return &ApuestaHandler{
		DB:      db,
		queries: sqlc.New(db),
	}
}

// HandlerApuesta despacha según el método HTTP para la ruta /apuestas.
func (h *ApuestaHandler) HandlerApuesta(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.CrearApuesta(w, r)
	case http.MethodGet:
		h.ListarApuestas(w, r)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

// HandlerApuestaIdUsuario despacha según el método HTTP para la ruta de una apuesta puntual de un usuario.
func (h *ApuestaHandler) HandlerApuestaIdUsuario(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.ObtenerApuestaPorIdUsuario(w, r)
	case http.MethodPut:
		h.ActualizarApuestaIdUsuario(w, r)
	case http.MethodDelete:
		h.EliminarApuesta(w, r)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

func (h *ApuestaHandler) HandlerApuestaUsuario(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.ObtenerApuestaUsuario(w, r)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

func (h *ApuestaHandler) CrearApuesta(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var apuesta sqlc.CrearApuestaParams

	err := json.NewDecoder(r.Body).Decode(&apuesta)
	if err != nil {
		http.Error(w, "Error al decodificar el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}

	// fecha_apuesta la asigna el servidor, por eso no se valida acá
	if apuesta.NombreUsuario == "" || apuesta.IDGranPremio == "" || len(apuesta.Prediccion) == 0 || apuesta.FechaCarrera.IsZero() {
		http.Error(w, "Todos los campos son obligatorios", http.StatusBadRequest)
		return
	}

	// Asignar la fecha actual como fecha de la apuesta
	apuesta.FechaApuesta = time.Now()

	// Verificar que el usuario existe
	_, err = h.queries.RecuperarUsuarioPorId(r.Context(), apuesta.NombreUsuario)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "No existe el usuario", http.StatusBadRequest)
		} else {
			http.Error(w, "Error al verificar el usuario", http.StatusInternalServerError)
		}
		return
	}

	var granPremio sqlc.RecuperarGranPremioHistoricoParams
	granPremio.IDGranPremio = apuesta.IDGranPremio
	granPremio.FechaCarrera = apuesta.FechaCarrera

	// Verificar que el gran premio existe
	_, err = h.queries.RecuperarGranPremioHistorico(r.Context(), granPremio)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "No existe el gran premio", http.StatusBadRequest)
		} else {
			http.Error(w, "Error al verificar el gran premio", http.StatusInternalServerError)
		}
		return
	}

	// Verificar que todos los pilotos en la predicción existen
	for _, pilotoID := range apuesta.Prediccion {
		_, err = h.queries.RecuperarPiloto(r.Context(), pilotoID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.Error(w, "No existe el piloto en la predicción", http.StatusBadRequest)
			} else {
				http.Error(w, "Error al verificar el piloto en la predicción", http.StatusInternalServerError)
			}
			return
		}
	}

	devolver, err := h.queries.CrearApuesta(r.Context(), apuesta)
	if err != nil {
		http.Error(w, "Error al crear la apuesta", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(devolver)
}

func (h *ApuestaHandler) ListarApuestas(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	apuestas, err := h.queries.ListarApuestas(r.Context())
	if err != nil {
		http.Error(w, "Error al obtener las apuestas", http.StatusInternalServerError)
		return
	}
	if len(apuestas) == 0 {
		http.Error(w, "No se encontraron apuestas", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(apuestas)
}

func (h *ApuestaHandler) ObtenerApuestaUsuario(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	nombreUsuario := r.PathValue("usuario")
	if nombreUsuario == "" {
		http.Error(w, "Nombre de usuario inválido", http.StatusBadRequest)
		return
	}

	_, err := h.queries.RecuperarUsuarioPorId(r.Context(), nombreUsuario)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "No existe el usuario", http.StatusBadRequest)
		} else {
			http.Error(w, "Error al verificar el usuario", http.StatusInternalServerError)
		}
		return
	}

	apuestas, err := h.queries.ListarApuestasPorUsuario(r.Context(), nombreUsuario)
	if err != nil {
		http.Error(w, "Error al obtener las apuestas", http.StatusInternalServerError)
		return
	}
	if len(apuestas) == 0 {
		http.Error(w, "No se encontraron apuestas", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(apuestas)
}

func (h *ApuestaHandler) ObtenerApuestaPorIdUsuario(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	idApuesta, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
	if err != nil {
		http.Error(w, "ID de apuesta inválido", http.StatusBadRequest)
		return
	}
	nombreUsuario := r.PathValue("usuario")
	if nombreUsuario == "" {
		http.Error(w, "Nombre de usuario inválido", http.StatusBadRequest)
		return
	}

	_, err = h.queries.RecuperarUsuarioPorId(r.Context(), nombreUsuario)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "No existe el usuario", http.StatusBadRequest)
		} else {
			http.Error(w, "Error al verificar el usuario", http.StatusInternalServerError)
		}
		return
	}

	var PK_apuesta sqlc.RecuperarApuestaPorIdUsuarioParams
	PK_apuesta.IDApuesta = int32(idApuesta)
	PK_apuesta.NombreUsuario = nombreUsuario

	apuesta, err := h.queries.RecuperarApuestaPorIdUsuario(r.Context(), PK_apuesta)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Apuesta no encontrada", http.StatusNotFound)
		} else {
			http.Error(w, "Error al obtener la apuesta", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(apuesta)
}

func (h *ApuestaHandler) ActualizarApuestaIdUsuario(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	idApuesta, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
	if err != nil {
		http.Error(w, "ID de apuesta inválido", http.StatusBadRequest)
		return
	}
	nombreUsuario := r.PathValue("usuario")
	if nombreUsuario == "" {
		http.Error(w, "Nombre de usuario inválido", http.StatusBadRequest)
		return
	}

	var apuesta sqlc.ModificarApuestaParams

	err = json.NewDecoder(r.Body).Decode(&apuesta)
	if err != nil {
		http.Error(w, "Error al decodificar el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}

	// La identidad de la apuesta sale siempre de la URL, no del body
	apuesta.IDApuesta = int32(idApuesta)
	apuesta.NombreUsuario = nombreUsuario

	if apuesta.IDGranPremio == "" || len(apuesta.Prediccion) == 0 || apuesta.FechaCarrera.IsZero() || apuesta.FechaApuesta.IsZero() {
		http.Error(w, "Todos los campos son obligatorios", http.StatusBadRequest)
		return
	}

	devolver, err := h.queries.ModificarApuesta(r.Context(), apuesta)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Apuesta no encontrada", http.StatusNotFound)
		} else {
			http.Error(w, "Error al actualizar la apuesta", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(devolver)
}

func (h *ApuestaHandler) EliminarApuesta(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	idApuesta, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
	if err != nil {
		http.Error(w, "ID de apuesta inválido", http.StatusBadRequest)
		return
	}
	nombreUsuario := r.PathValue("usuario")
	if nombreUsuario == "" {
		http.Error(w, "Nombre de usuario inválido", http.StatusBadRequest)
		return
	}

	var apuestaEliminada sqlc.EliminarApuestaParams
	apuestaEliminada.IDApuesta = int32(idApuesta)
	apuestaEliminada.NombreUsuario = nombreUsuario

	err = h.queries.EliminarApuesta(r.Context(), apuestaEliminada)
	if err != nil {
		http.Error(w, "Error al eliminar la apuesta", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
