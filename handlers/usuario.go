package handlers

import (
	sqlc "Proyecto_web/db/sqlc"
	views "Proyecto_web/views"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"golang.org/x/crypto/bcrypt" // para hashear contraseñas
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

func (h *UsuarioHandler) PaginaRegistrar(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		views.RegistrarPage("").Render(r.Context(), w)
	case http.MethodPost:
		h.CrearUsuario(w, r)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

func (h *UsuarioHandler) PaginaIngresar(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if _, ok := usuarioActual(r); ok {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		views.IngresarPage("").Render(r.Context(), w)

	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			views.IngresarPage("Formulario inválido").Render(r.Context(), w)
			return
		}

		nombre_usuario := r.FormValue("nombre_usuario")
		contraseña := r.FormValue("contraseña")

		usuario, error := h.autenticarUsuario(r.Context(), nombre_usuario, contraseña)
		if error != nil {
			w.WriteHeader(http.StatusUnauthorized)
			views.IngresarPage("Usuario o contraseña incorrectos").Render(r.Context(), w)
			return
		}

		if err := crearSesion(w, usuario.NombreUsuario); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			views.IngresarPage("Error interno").Render(r.Context(), w)
			return
		}
		if usuario.EsAdmin {
			http.Redirect(w, r, "/admin", http.StatusSeeOther)
		} else {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func (h *UsuarioHandler) Salir(w http.ResponseWriter, r *http.Request) {
	cerrarSesion(w, r)
	http.Redirect(w, r, "/ingresar", http.StatusSeeOther)
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

func (h *UsuarioHandler) autenticarUsuario(ctx context.Context, nombre, password string) (sqlc.Usuario, error) {
	if nombre == "" || password == "" {
		return sqlc.Usuario{}, errors.New("usuario o contraseña incorrectos")
	}

	usuario, err := h.queries.RecuperarUsuarioPorId(ctx, nombre)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sqlc.Usuario{}, errors.New("usuario o contraseña incorrectos")
		}
		return sqlc.Usuario{}, err
	}

	if bcrypt.CompareHashAndPassword([]byte(usuario.ContrasenaHash), []byte(password)) != nil {
		return sqlc.Usuario{}, errors.New("usuario o contraseña incorrectos")
	}

	return usuario, nil
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

	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		views.RegistrarPage("Formulario inválido").Render(r.Context(), w)
		return
	}

	nombre := r.FormValue("nombre_usuario")
	contraseña := r.FormValue("contraseña")

	if nombre == "" || contraseña == "" {
		w.WriteHeader(http.StatusBadRequest)
		views.RegistrarPage("Todos los campos son obligatorios").Render(r.Context(), w)
		return
	}

	//hasheo de contraseña para guardarla de manera segura en la base de datos
	hash, err := bcrypt.GenerateFromPassword([]byte(contraseña), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		views.RegistrarPage("Error interno").Render(r.Context(), w)
		return
	}

	_, err = h.queries.CrearUsuario(r.Context(), sqlc.CrearUsuarioParams{
		NombreUsuario:  nombre,
		ContrasenaHash: string(hash),
		PuntosGanados:  0, // un usuario nuevo arranca en 0 puntos
	})
	if err != nil {
		w.WriteHeader(http.StatusConflict)
		views.RegistrarPage("Error al crear el usuario").Render(r.Context(), w)
		return
	}

	http.Redirect(w, r, "/ingresar", http.StatusSeeOther)
}
