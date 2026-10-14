package handlers

import (
	sqlc "Proyecto_web/db/sqlc"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// UsuarioHandler agrupa los handlers de usuario y guarda
// la conexión a la base de datos y las queries de sqlc.
type UsuarioHandler struct {
	DB      *sql.DB
	queries *sqlc.Queries
}

// NewUsuarioHandler es el constructor del handler.
func NewUsuarioHandler(db *sql.DB) *UsuarioHandler {
	return &UsuarioHandler{
		DB:      db,
		queries: sqlc.New(db),
	}
}

func (h *UsuarioHandler) EliminarUsuario(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	nombreUsuario := r.PathValue("nombre")
	if nombreUsuario == "" {
		http.Error(w, "Nombre de usuario invalido", http.StatusBadRequest)
		return
	}

	err := h.queries.EliminarUsuario(r.Context(), nombreUsuario)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Usuario no encontrado", http.StatusNotFound)
		} else {
			http.Error(w, "Error al eliminar el usuario", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *UsuarioHandler) ActualizarUsuario(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	nombreUsuario := r.PathValue("nombre")
	if nombreUsuario == "" {
		http.Error(w, "Nombre de usuario invalido", http.StatusBadRequest)
		return
	}

	var usuario sqlc.ModificarUsuarioParams

	err := json.NewDecoder(r.Body).Decode(&usuario)
	if err != nil {
		http.Error(w, "Error al decodificar el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}

	if usuario.ContrasenaHash == "" {
		http.Error(w, "Todos los campos son obligatorios", http.StatusBadRequest)
		return
	}

	if usuario.PuntosGanados < 0 {
		http.Error(w, "Los puntos ganados no pueden ser negativos", http.StatusBadRequest)
		return
	}

	usuario.NombreUsuario = nombreUsuario

	devolver, err := h.queries.ModificarUsuario(r.Context(), usuario)
	if err != nil {
		http.Error(w, "Error al actualizar el usuario", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(devolver)
}

func (h *UsuarioHandler) ObtenerUsuario(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	nombreUsuario := r.PathValue("nombre")
	if nombreUsuario == "" {
		http.Error(w, "Nombre de usuario invalido", http.StatusBadRequest)
		return
	}

	usuario, err := h.queries.RecuperarUsuarioPorId(r.Context(), nombreUsuario)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Usuario no encontrado", http.StatusNotFound)
		} else {
			http.Error(w, "Error al obtener el usuario", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(usuario)
}

func (h *UsuarioHandler) ListarUsuarios(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	usuarios, err := h.queries.RecuperarUsuarios(r.Context())
	if err != nil {
		http.Error(w, "Error al obtener los usuarios", http.StatusInternalServerError)
		return
	}
	if len(usuarios) == 0 {
		http.Error(w, "No se encontraron usuarios", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(usuarios)
}

func (h *UsuarioHandler) CrearUsuario(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var usuario sqlc.CrearUsuarioParams

	err := json.NewDecoder(r.Body).Decode(&usuario)
	if err != nil {
		http.Error(w, "Error al decodificar el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}

	if usuario.NombreUsuario == "" || usuario.ContrasenaHash == "" {
		http.Error(w, "Todos los campos son obligatorios", http.StatusBadRequest)
		return
	}

	usuario.PuntosGanados = 0 //setea 0 puntos ganados al crear un nuevo usuario

	devolver, err := h.queries.CrearUsuario(r.Context(), usuario)
	if err != nil {
		http.Error(w, "Error al crear el usuario", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(devolver)
}

// HandlerUsuario despacha según el método HTTP para la ruta /usuarios.
func (h *UsuarioHandler) HandlerUsuario(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.CrearUsuario(w, r)
	case http.MethodGet:
		h.ListarUsuarios(w, r)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		fmt.Println("entro al default")
	}
}

// HandlerUsuarioId despacha según el método HTTP para la ruta /usuarios/{...}.
func (h *UsuarioHandler) HandlerUsuarioId(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.ObtenerUsuario(w, r)
	case http.MethodPut:
		h.ActualizarUsuario(w, r)
	case http.MethodDelete:
		h.EliminarUsuario(w, r)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}
