// handlers/usuario.go
package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"Proyecto_web/aplicacion"
	"Proyecto_web/dominio"
	views "Proyecto_web/views"
)

// usuarioPublico mantiene la forma JSON de la API (sin exponer el hash).
type usuarioPublico struct {
	Nombre string `json:"nombre_usuario"`
	Puntos int    `json:"puntos_ganados"`
}

type UsuarioHandler struct {
	usuarios *aplicacion.UsuarioService
}

func NewUsuarioHandler(u *aplicacion.UsuarioService) *UsuarioHandler {
	return &UsuarioHandler{usuarios: u}
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

		usuario, err := h.usuarios.Autenticar(r.Context(), r.FormValue("nombre_usuario"), r.FormValue("contraseña"))
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			views.IngresarPage("Usuario o contraseña incorrectos").Render(r.Context(), w)
			return
		}

		if err := crearSesion(w, usuario.Nombre); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			views.IngresarPage("Error interno").Render(r.Context(), w)
			return
		}
		if usuario.EsAdmin {
			http.Redirect(w, r, "/admin", http.StatusSeeOther)
		} else {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
	}
}

func (h *UsuarioHandler) Salir(w http.ResponseWriter, r *http.Request) {
	cerrarSesion(w, r)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *UsuarioHandler) EliminarUsuario(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	nombre := r.PathValue("nombre")
	if nombre == "" {
		http.Error(w, "Nombre de usuario invalido", http.StatusBadRequest)
		return
	}
	if err := h.usuarios.Eliminar(r.Context(), nombre); err != nil {
		if errors.Is(err, dominio.ErrUsuarioNoEncontrado) {
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

	nombre := r.PathValue("nombre")
	if nombre == "" {
		http.Error(w, "Nombre de usuario invalido", http.StatusBadRequest)
		return
	}

	var u dominio.Usuario
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "Error al decodificar el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}
	if u.Puntos < 0 {
		http.Error(w, "Los puntos ganados no pueden ser negativos", http.StatusBadRequest)
		return
	}
	u.Nombre = nombre

	if err := h.usuarios.Modificar(r.Context(), u); err != nil {
		var ev dominio.ErrValidacion
		if errors.As(err, &ev) {
			http.Error(w, ev.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "Error al actualizar el usuario", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(usuarioPublico{Nombre: u.Nombre, Puntos: u.Puntos})
}

func (h *UsuarioHandler) ListarUsuarios(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	usuarios, err := h.usuarios.Listar(r.Context())
	if err != nil {
		http.Error(w, "Error al obtener los usuarios", http.StatusInternalServerError)
		return
	}
	if len(usuarios) == 0 {
		http.Error(w, "No se encontraron usuarios", http.StatusNotFound)
		return
	}

	out := make([]usuarioPublico, 0, len(usuarios))
	for _, u := range usuarios {
		out = append(out, usuarioPublico{Nombre: u.Nombre, Puntos: u.Puntos})
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(out)
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

	u := dominio.Usuario{Nombre: r.FormValue("nombre_usuario"), Contrasena: r.FormValue("contraseña")}
	if err := h.usuarios.Crear(r.Context(), u); err != nil {
		var ev dominio.ErrValidacion
		switch {
		case errors.As(err, &ev):
			w.WriteHeader(http.StatusBadRequest)
			views.RegistrarPage(ev.Error()).Render(r.Context(), w)
		case errors.Is(err, dominio.ErrUsuarioDuplicado):
			w.WriteHeader(http.StatusConflict)
			views.RegistrarPage("Ya existe ese nombre de usuario").Render(r.Context(), w)
		default:
			w.WriteHeader(http.StatusInternalServerError)
			views.RegistrarPage("Error al crear el usuario").Render(r.Context(), w)
		}
		return
	}
	http.Redirect(w, r, "/ingresar", http.StatusSeeOther)
}
