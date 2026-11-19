// handlers/admin.go
package handlers

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"Proyecto_web/aplicacion"
	"Proyecto_web/dominio"
	"Proyecto_web/views"
)

type AdminHandler struct {
	admin *aplicacion.AdminService
}

func NewAdminHandler(s *aplicacion.AdminService) *AdminHandler {
	return &AdminHandler{admin: s}
}

// RequiereAdmin: exige sesión y rol de admin, consultado a través del servicio.
func (h *AdminHandler) RequiereAdmin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		usuario, ok := usuarioActual(r)
		if !ok {
			http.Redirect(w, r, "/ingresar", http.StatusSeeOther)
			return
		}
		if !h.admin.EsAdmin(r.Context(), usuario) {
			http.Error(w, "No autorizado", http.StatusForbidden)
			return
		}
		next(w, r)
	}
}

// ---------------------------------------------------------------
// GET /admin
// ---------------------------------------------------------------

func (h *AdminHandler) PaginaAdmin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}
	aviso := ""
	switch r.URL.Query().Get("ok") {
	case "carrera":
		aviso = "Carrera programada. Ya se puede apostar."
	case "resultado":
		aviso = "Resultado cargado y puntos sumados."
	}
	h.renderAdmin(w, r, aviso, "", http.StatusOK)
}

func (h *AdminHandler) renderAdmin(w http.ResponseWriter, r *http.Request, aviso, mensaje string, status int) {
	ctx := r.Context()
	data := views.AdminData{Aviso: aviso, Mensaje: mensaje}

	circuitos, err := h.admin.Circuitos(ctx)
	if err != nil {
		log.Println("Listar circuitos:", err)
	}
	data.Circuitos = circuitos

	pendientes, err := h.admin.Pendientes(ctx)
	if err != nil {
		log.Println("Listar carreras pendientes:", err)
	}
	for _, p := range pendientes {
		data.Pendientes = append(data.Pendientes, views.OpcionGP{
			Valor: p.ID + "|" + p.Fecha.Format("2006-01-02"),
			Texto: p.ID + " — " + p.Fecha.Format("02/01/2006"),
		})
	}

	pilotos, err := h.admin.Pilotos(ctx)
	if err != nil {
		log.Println("Listar pilotos:", err)
	}
	for _, p := range pilotos {
		data.Pilotos = append(data.Pilotos, views.OpcionPiloto{ID: p.ID, Nombre: p.Nombre})
	}

	w.WriteHeader(status)
	views.AdminPage(data).Render(ctx, w)
}

// ---------------------------------------------------------------
// POST /admin/carrera: programar una carrera futura
// ---------------------------------------------------------------

func (h *AdminHandler) ProgramarCarrera(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseForm(); err != nil {
		h.renderAdmin(w, r, "", "Formulario inválido", http.StatusBadRequest)
		return
	}

	circuito := r.FormValue("circuito")
	fecha, err := time.Parse("2006-01-02", r.FormValue("fecha"))
	if err != nil {
		h.renderAdmin(w, r, "", "Elegí un circuito y una fecha válida", http.StatusBadRequest)
		return
	}
	if err := h.admin.ProgramarCarrera(r.Context(), circuito, fecha); err != nil {
		mensaje, status := "No se pudo programar la carrera", http.StatusInternalServerError
		var ev dominio.ErrValidacion
		switch {
		case errors.As(err, &ev):
			mensaje, status = ev.Error(), http.StatusBadRequest
		case errors.Is(err, dominio.ErrEdicionDuplicada):
			mensaje, status = "Esa carrera ya está programada", http.StatusConflict
		}
		h.renderAdmin(w, r, "", mensaje, status)
		return
	}
	http.Redirect(w, r, "/admin?ok=carrera", http.StatusSeeOther)
}

// ---------------------------------------------------------------
// POST /admin/resultado: cargar los 22 resultados y repartir puntos
// ---------------------------------------------------------------

func (h *AdminHandler) CargarResultado(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseForm(); err != nil {
		h.renderAdmin(w, r, "", "Formulario inválido", http.StatusBadRequest)
		return
	}
	ctx := r.Context()

	idGP, fechaStr, ok := strings.Cut(r.FormValue("edicion"), "|")
	fecha, err := time.Parse("2006-01-02", fechaStr)
	if !ok || err != nil {
		h.renderAdmin(w, r, "", "Elegí una carrera", http.StatusBadRequest)
		return
	}
	orden, msg := parsearResultado(r)
	if msg != "" {
		h.renderAdmin(w, r, "", msg, http.StatusBadRequest)
		return
	}

	if err := h.admin.CargarResultado(ctx, idGP, fecha, orden); err != nil {
		mensaje, status := "No se pudo cargar el resultado", http.StatusInternalServerError
		var ev dominio.ErrValidacion
		switch {
		case errors.As(err, &ev):
			mensaje, status = ev.Error(), http.StatusBadRequest
		case errors.Is(err, dominio.ErrCarreraNoEncontrada):
			mensaje, status = "Esa carrera no existe o ya tiene resultado", http.StatusConflict
		}
		h.renderAdmin(w, r, "", mensaje, status)
		return
	}
	http.Redirect(w, r, "/admin?ok=resultado", http.StatusSeeOther)
}

// parsearResultado lee el orden de llegada del formulario; las reglas las valida el dominio.
func parsearResultado(r *http.Request) ([]int, string) {
	valores := r.Form["piloto"]
	if len(valores) != 22 {
		return nil, "Tenés que ordenar los 22 pilotos"
	}
	orden := make([]int, 0, 22)
	for _, v := range valores {
		n, err := strconv.Atoi(v)
		if err != nil {
			return nil, "Piloto inválido"
		}
		orden = append(orden, n)
	}
	return orden, ""
}
