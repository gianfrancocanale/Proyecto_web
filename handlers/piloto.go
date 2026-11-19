package handlers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"Proyecto_web/aplicacion"
	"Proyecto_web/dominio"
	vista "Proyecto_web/views"
)

type PilotoHandler struct {
	pilotos  *aplicacion.PilotoService
	usuarios *aplicacion.UsuarioService
}

func NewPilotoHandler(p *aplicacion.PilotoService, u *aplicacion.UsuarioService) *PilotoHandler {
	return &PilotoHandler{pilotos: p, usuarios: u}
}

func (h *PilotoHandler) esAdmin(ctx context.Context, usuario string) bool {
	if usuario == "" {
		return false
	}
	esAdmin, err := h.usuarios.EsAdmin(ctx, usuario)
	return err == nil && esAdmin
}

// Página principal: listar (todos) y crear (admin)

func (h *PilotoHandler) Pagina(w http.ResponseWriter, r *http.Request) {
	usuario, _ := usuarioActual(r) // listar es público: no redirige si no hay sesión

	switch r.Method {
	case http.MethodGet:
		orden := r.URL.Query().Get("orden")
		if orden != "asc" {
			orden = "desc"
		}
		data := h.datosPilotos(r.Context(), usuario, orden, "")

		w.Header().Add("Vary", "HX-Request")
		if r.Header.Get("HX-Request") == "true" {
			vista.PilotosLista(data.Pilotos, data.Orden, data.EsAdmin, data.Mensaje).Render(r.Context(), w)
			return
		}
		vista.PilotosPage(data).Render(r.Context(), w)

	case http.MethodPost:
		h.crearPilotoDesdeForm(w, r, usuario)

	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

// mostrarError devuelve la lista con un mensaje arriba.
// Responde 200 a propósito: HTMX no reemplaza el contenido con 4xx/5xx.
func (h *PilotoHandler) mostrarError(w http.ResponseWriter, r *http.Request, usuario, mensaje string) {
	h.renderLista(w, r, usuario, mensaje)
}

func (h *PilotoHandler) crearPilotoDesdeForm(w http.ResponseWriter, r *http.Request, usuario string) {
	if err := r.ParseForm(); err != nil {
		h.mostrarError(w, r, usuario, "Formulario inválido")
		return
	}

	n, err := strconv.Atoi(strings.TrimSpace(r.FormValue("id_piloto")))
	if err != nil || n <= 0 {
		h.mostrarError(w, r, usuario, "Número de piloto inválido")
		return
	}

	p := dominio.Piloto{
		ID:        n,
		Nombre:    strings.TrimSpace(r.FormValue("nombre")),
		Pais:      strings.TrimSpace(r.FormValue("pais")),
		Escuderia: strings.TrimSpace(r.FormValue("id_escuderia")),
		Puntos:    0,
	}
	if p.Nacimiento, err = time.Parse("2006-01-02", r.FormValue("fecha_nacimiento")); err != nil {
		// Si el piloto ya existe en el histórico, la fecha del formulario viene vacía.
		p.Nacimiento = time.Time{}
	}
	if v := strings.TrimSpace(r.FormValue("titulos_ganados")); v != "" {
		p.Titulos, err = strconv.Atoi(v)
		if err != nil || p.Titulos < 0 {
			h.mostrarError(w, r, usuario, "Títulos ganados inválidos")
			return
		}
	}

	if err := h.pilotos.Inscribir(r.Context(), p); err != nil {
		h.mostrarError(w, r, usuario, mensajeErrorPiloto(err))
		return
	}
	h.renderLista(w, r, usuario, "")
}

// Datos y render

func (h *PilotoHandler) datosPilotos(ctx context.Context, usuario, orden, mensaje string) vista.PilotosData {
	esAdm := h.esAdmin(ctx, usuario)
	data := vista.PilotosData{
		Usuario: usuario,
		EsAdmin: esAdm,
		Orden:   orden,
		Mensaje: mensaje,
	}

	pilotos, err := h.pilotos.Listar(ctx, orden)
	if err != nil {
		log.Println("Listar pilotos:", err)
	} else {
		data.Pilotos = pilotos
	}
	if esAdm {
		data.Escuderias = h.opcionesEscuderias(ctx)
	}
	return data
}

func (h *PilotoHandler) renderLista(w http.ResponseWriter, r *http.Request, usuario, mensaje string) {
	orden := r.FormValue("orden")
	if orden != "asc" {
		orden = "desc"
	}
	data := h.datosPilotos(r.Context(), usuario, orden, mensaje)
	vista.PilotosLista(data.Pilotos, data.Orden, data.EsAdmin, data.Mensaje).Render(r.Context(), w)
}

func (h *PilotoHandler) opcionesEscuderias(ctx context.Context) []vista.OpcionEscuderia {
	ids, err := h.pilotos.Escuderias(ctx)
	if err != nil {
		log.Println("Listar escuderías:", err)
		return nil
	}
	out := make([]vista.OpcionEscuderia, 0, len(ids))
	for _, id := range ids {
		out = append(out, vista.OpcionEscuderia{ID: id})
	}
	return out
}

func mensajeErrorPiloto(err error) string {
	var ev dominio.ErrValidacion
	switch {
	case errors.As(err, &ev):
		return ev.Error()
	case errors.Is(err, dominio.ErrEscuderiaNoExiste):
		return "No existe la escudería"
	case errors.Is(err, dominio.ErrPilotoDuplicado):
		return "Ese piloto ya está en la temporada"
	case errors.Is(err, dominio.ErrHistoricoIncompleto):
		return "Piloto nuevo: completá nombre, país y fecha de nacimiento"
	default:
		log.Println("piloto:", err)
		return "Error interno"
	}
}

// Modificar un piloto: /pilotos/modificar  (admin)

func (h *PilotoHandler) ModificarPiloto(w http.ResponseWriter, r *http.Request) {
	usuario, ok := usuarioActual(r)
	if !ok {
		http.Redirect(w, r, "/ingresar", http.StatusSeeOther)
		return
	}

	switch r.Method {
	case http.MethodGet:
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil {
			http.Redirect(w, r, "/pilotos", http.StatusSeeOther)
			return
		}
		h.renderModificar(w, r, usuario, id, "", http.StatusOK)

	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			http.Redirect(w, r, "/pilotos", http.StatusSeeOther)
			return
		}
		id, err := strconv.Atoi(r.FormValue("id_piloto"))
		if err != nil {
			http.Redirect(w, r, "/pilotos", http.StatusSeeOther)
			return
		}
		puntos, err := strconv.Atoi(r.FormValue("puntos_temporada"))
		if err != nil || puntos < 0 {
			h.renderModificar(w, r, usuario, id, "Puntos inválidos", http.StatusBadRequest)
			return
		}

		p := dominio.Piloto{
			ID:        id,
			Escuderia: r.FormValue("id_escuderia"),
			Puntos:    puntos,
		}
		if err := h.pilotos.Modificar(r.Context(), p); err != nil {
			h.renderModificar(w, r, usuario, id, mensajeErrorPiloto(err), http.StatusBadRequest)
			return
		}
		http.Redirect(w, r, "/pilotos", http.StatusSeeOther)

	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

func (h *PilotoHandler) renderModificar(w http.ResponseWriter, r *http.Request, usuario string, id int, mensaje string, status int) {
	ctx := r.Context()

	p, err := h.pilotos.Recuperar(ctx, id)
	if err != nil {
		http.Redirect(w, r, "/pilotos", http.StatusSeeOther)
		return
	}

	data := vista.ModificarPilotoData{
		Usuario:    usuario,
		Mensaje:    mensaje,
		ID:         p.ID,
		Nombre:     p.Nombre,
		Escuderia:  p.Escuderia,
		Puntos:     p.Puntos,
		Escuderias: h.opcionesEscuderias(ctx),
	}
	w.WriteHeader(status)
	vista.ModificarPilotoPage(data).Render(ctx, w)
}

// Eliminar un piloto (admin)

// DELETE /pilotos/{id}: 200 vacío para que HTMX saque la tarjeta.
func (h *PilotoHandler) EliminarPiloto(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
	if err != nil {
		http.Error(w, "ID de piloto inválido", http.StatusBadRequest)
		return
	}
	if err := h.pilotos.Eliminar(r.Context(), int(id)); err != nil {
		log.Println("Eliminar piloto:", err)
		http.Error(w, "No se pudo borrar", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
