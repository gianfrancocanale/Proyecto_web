// handlers/escuderia.go
package handlers

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"Proyecto_web/aplicacion"
	"Proyecto_web/dominio"
	view "Proyecto_web/views"
)

type EscuderiaHandler struct {
	escuderias *aplicacion.EscuderiaService
	usuarios   *aplicacion.UsuarioService
}

func NewEscuderiaHandler(e *aplicacion.EscuderiaService, u *aplicacion.UsuarioService) *EscuderiaHandler {
	return &EscuderiaHandler{escuderias: e, usuarios: u}
}

// esAdmin: la cookie es cosa de HTTP; el rol lo consulta el servicio.
func (h *EscuderiaHandler) esAdmin(r *http.Request) bool {
	usuario, ok := usuarioActual(r)
	if !ok {
		return false
	}
	esAdmin, err := h.usuarios.EsAdmin(r.Context(), usuario)
	return err == nil && esAdmin
}

// ---------------------------------------------------------------
// GET /escuderias
// ---------------------------------------------------------------

func (h *EscuderiaHandler) Pagina(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	orden := normalizarOrden(r.URL.Query().Get("orden"))

	escuderias, err := h.escuderias.Listar(r.Context(), orden)
	if err != nil {
		log.Println("Listar escuderías:", err)
		http.Error(w, "Error al obtener las escuderías", http.StatusInternalServerError)
		return
	}
	esAdmin := h.esAdmin(r)

	w.Header().Add("Vary", "HX-Request")
	if r.Header.Get("HX-Request") == "true" {
		view.EscuderiasLista(escuderias, orden, esAdmin, "").Render(r.Context(), w)
		return
	}

	view.EscuderiasPage(view.EscuderiasData{
		EsAdmin:    esAdmin,
		Orden:      orden,
		Escuderias: escuderias,
	}).Render(r.Context(), w)
}

// ---------------------------------------------------------------
// POST /escuderias (htmx)
// ---------------------------------------------------------------

func (h *EscuderiaHandler) CrearEscuderia(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.renderLista(w, r, "", "Formulario inválido")
		return
	}
	orden := r.FormValue("orden")

	e, msg := parsearEscuderiaForm(r)
	if msg != "" {
		h.renderLista(w, r, orden, msg)
		return
	}

	if err := h.escuderias.Crear(r.Context(), e); err != nil {
		h.renderLista(w, r, orden, mensajeDeError(err))
		return
	}
	h.renderLista(w, r, orden, "")
}

// renderLista devuelve solo el fragmento #lista-escuderias.
func (h *EscuderiaHandler) renderLista(w http.ResponseWriter, r *http.Request, orden, mensaje string) {
	orden = normalizarOrden(orden)

	es, err := h.escuderias.Listar(r.Context(), orden)
	if err != nil {
		log.Println("Listar escuderías:", err)
		mensaje = "Error interno"
	}

	view.EscuderiasLista(es, orden, true, mensaje).Render(r.Context(), w)
}

// ---------------------------------------------------------------
// DELETE /escuderias/{id} (htmx)
// ---------------------------------------------------------------

func (h *EscuderiaHandler) EliminarEscuderia(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	err := h.escuderias.Eliminar(r.Context(), r.PathValue("id"))
	if err != nil {
		var ev dominio.ErrValidacion
		if errors.As(err, &ev) {
			http.Error(w, ev.Error(), http.StatusBadRequest)
			return
		}
		log.Println("Eliminar escudería:", err)
		http.Error(w, "Error al eliminar la escudería", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK) // cuerpo vacío: htmx saca la tarjeta
}

// ---------------------------------------------------------------
// GET /escuderias/modificar?id=...
// ---------------------------------------------------------------

func (h *EscuderiaHandler) PaginaModificar(w http.ResponseWriter, r *http.Request) {
	e, err := h.escuderias.Recuperar(r.Context(), r.URL.Query().Get("id"))
	if err != nil {
		http.Redirect(w, r, "/escuderias", http.StatusSeeOther)
		return
	}

	view.ModificarEscuderiaPage(view.ModificarEscuderiaData{
		ID:        e.ID,
		Principal: e.Principal,
		Titulos:   e.Titulos,
		Fundacion: e.Fundacion.Format("2006-01-02"),
		Puntos:    e.Puntos,
	}).Render(r.Context(), w)
}

// ---------------------------------------------------------------
// POST /escuderias/modificar
// ---------------------------------------------------------------

func (h *EscuderiaHandler) ModificarEscuderia(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Redirect(w, r, "/escuderias", http.StatusSeeOther)
		return
	}

	e, msg := parsearEscuderiaForm(r)
	if msg == "" {
		err := h.escuderias.Modificar(r.Context(), e)
		if err == nil || errors.Is(err, dominio.ErrNoEncontrada) {
			http.Redirect(w, r, "/escuderias", http.StatusSeeOther)
			return
		}
		msg = mensajeDeError(err)
		var ev dominio.ErrValidacion
		if errors.As(err, &ev) {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}

	// Error: volvemos a mostrar el formulario con lo que escribió el usuario.
	view.ModificarEscuderiaPage(view.ModificarEscuderiaData{
		Mensaje:   msg,
		ID:        r.FormValue("id_escuderia"),
		Principal: r.FormValue("team_principal"),
		Titulos:   atoiOCero(r.FormValue("titulos_constructores")),
		Fundacion: r.FormValue("fecha_fundacion"),
		Puntos:    atoiOCero(r.FormValue("puntos_temporada")),
	}).Render(r.Context(), w)
}

// ---------------------------------------------------------------
// Helpers (solo cosas de HTTP: parseo y traducción de errores)
// ---------------------------------------------------------------

func normalizarOrden(orden string) string {
	if orden != "desc" {
		return "asc"
	}
	return orden
}

// parsearEscuderiaForm convierte strings a tipos. Las reglas de negocio las valida el dominio.
func parsearEscuderiaForm(r *http.Request) (dominio.Escuderia, string) {
	e := dominio.Escuderia{
		ID:        r.FormValue("id_escuderia"),
		Principal: r.FormValue("team_principal"),
	}

	var err error
	if e.Titulos, err = strconv.Atoi(r.FormValue("titulos_constructores")); err != nil {
		return e, "Los títulos tienen que ser un número"
	}
	if e.Puntos, err = strconv.Atoi(r.FormValue("puntos_temporada")); err != nil {
		return e, "Los puntos tienen que ser un número"
	}
	if e.Fundacion, err = time.Parse("2006-01-02", r.FormValue("fecha_fundacion")); err != nil {
		return e, "La fecha de fundación no es válida"
	}
	return e, ""
}

// mensajeDeError traduce errores de dominio a texto para la vista.
func mensajeDeError(err error) string {
	var ev dominio.ErrValidacion
	switch {
	case errors.As(err, &ev):
		return ev.Error()
	case errors.Is(err, dominio.ErrDuplicada):
		return "Ya existe una escudería con ese nombre"
	default:
		log.Println("escudería:", err)
		return "No se pudo completar la operación"
	}
}

func atoiOCero(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}
