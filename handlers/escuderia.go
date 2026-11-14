package handlers

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	sqlc "Proyecto_web/db/sqlc"
	view "Proyecto_web/views"
)

// EscuderiaHandler agrupa los handlers de escudería y guarda
// la conexión a la base de datos y las queries de sqlc.
type EscuderiaHandler struct {
	DB      *sql.DB
	queries *sqlc.Queries
}

// NewEscuderiaHandler es el constructor del handler.
func NewEscuderiaHandler(db *sql.DB) *EscuderiaHandler {
	return &EscuderiaHandler{
		DB:      db,
		queries: sqlc.New(db),
	}
}

// ---------------------------------------------------------------
// GET /escuderias
// ---------------------------------------------------------------

func (h *EscuderiaHandler) Pagina(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	escuderias, err := h.queries.ListarEscuderias(r.Context())
	if err != nil {
		log.Println("ListarEscuderias:", err)
		http.Error(w, "Error al obtener las escuderías", http.StatusInternalServerError)
		return
	}

	orden := r.URL.Query().Get("orden")
	if orden != "desc" {
		orden = "asc" // valor por defecto
	}
	ordenarEscuderias(escuderias, orden)
	esAdmin := usuarioEsAdmin(r, h.queries)

	// La misma URL devuelve página completa o fragmento según quién pregunte.
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
// POST /escuderias (htmx): crear y devolver la lista actualizada
// ---------------------------------------------------------------

func (h *EscuderiaHandler) CrearEscuderia(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.renderLista(w, r, "", "Formulario inválido")
		return
	}
	orden := r.FormValue("orden")

	f, msg := leerEscuderiaForm(r)
	if msg != "" {
		h.renderLista(w, r, orden, msg)
		return
	}

	_, err := h.queries.CrearEscuderia(r.Context(), sqlc.CrearEscuderiaParams{
		IDEscuderia:          f.ID,
		PuntosTemporada:      f.Puntos,
		TitulosConstructores: f.Titulos,
		FechaFundacion:       f.Fundacion,
		TeamPrincipal:        f.Principal,
	})
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			h.renderLista(w, r, orden, "Ya existe una escudería con ese nombre")
			return
		}
		log.Println("CrearEscuderia:", err)
		h.renderLista(w, r, orden, "No se pudo crear la escudería")
		return
	}

	h.renderLista(w, r, orden, "")
}

// renderLista devuelve solo el fragmento #lista-escuderias.
func (h *EscuderiaHandler) renderLista(w http.ResponseWriter, r *http.Request, orden, mensaje string) {
	es, err := h.queries.ListarEscuderias(r.Context())
	if err != nil {
		log.Println("ListarEscuderias:", err)
		mensaje = "Error interno"
	}
	if orden != "desc" {
		orden = "asc"
	}
	ordenarEscuderias(es, orden)

	view.EscuderiasLista(es, orden, true, mensaje).Render(r.Context(), w)
}

// ---------------------------------------------------------------
// DELETE /escuderias/{id} (htmx): la tarjeta se reemplaza por nada
// ---------------------------------------------------------------

func (h *EscuderiaHandler) EliminarEscuderia(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	idEscuderia := r.PathValue("id")
	if idEscuderia == "" {
		http.Error(w, "id inválido", http.StatusBadRequest)
		return
	}

	if err := h.queries.EliminarEscuderia(r.Context(), idEscuderia); err != nil {
		log.Println("EliminarEscuderia:", err)
		http.Error(w, "Error al eliminar la escudería", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK) // cuerpo vacío: htmx saca la tarjeta
}

// ---------------------------------------------------------------
// GET /escuderias/modificar?id=...
// ---------------------------------------------------------------

func (h *EscuderiaHandler) PaginaModificar(w http.ResponseWriter, r *http.Request) {
	e, err := h.queries.RecuperarEscuderia(r.Context(), r.URL.Query().Get("id"))
	if err != nil {
		http.Redirect(w, r, "/escuderias", http.StatusSeeOther)
		return
	}

	view.ModificarEscuderiaPage(view.ModificarEscuderiaData{
		ID:        e.IDEscuderia,
		Principal: e.TeamPrincipal,
		Titulos:   int(e.TitulosConstructores),
		Fundacion: e.FechaFundacion.Format("2006-01-02"),
		Puntos:    int(e.PuntosTemporada),
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

	f, msg := leerEscuderiaForm(r)
	if msg != "" {
		w.WriteHeader(http.StatusBadRequest)
		view.ModificarEscuderiaPage(view.ModificarEscuderiaData{
			Mensaje:   msg,
			ID:        f.ID,
			Principal: r.FormValue("team_principal"),
			Titulos:   atoiOCero(r.FormValue("titulos_constructores")),
			Fundacion: r.FormValue("fecha_fundacion"),
			Puntos:    atoiOCero(r.FormValue("puntos_temporada")),
		}).Render(r.Context(), w)
		return
	}

	_, err := h.queries.ModificarEscuderia(r.Context(), sqlc.ModificarEscuderiaParams{
		IDEscuderia:          f.ID,
		TeamPrincipal:        f.Principal,
		TitulosConstructores: f.Titulos,
		FechaFundacion:       f.Fundacion,
		PuntosTemporada:      f.Puntos,
	})
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Println("ModificarEscuderia:", err)
		w.WriteHeader(http.StatusInternalServerError)
		view.ModificarEscuderiaPage(view.ModificarEscuderiaData{
			Mensaje:   "No se pudo guardar el cambio",
			ID:        f.ID,
			Principal: f.Principal,
			Titulos:   int(f.Titulos),
			Fundacion: f.Fundacion.Format("2006-01-02"),
			Puntos:    int(f.Puntos),
		}).Render(r.Context(), w)
		return
	}

	http.Redirect(w, r, "/escuderias", http.StatusSeeOther)
}

// ---------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------

type escuderiaForm struct {
	ID        string
	Principal string
	Titulos   int32
	Puntos    int32
	Fundacion time.Time
}

// leerEscuderiaForm valida el formulario. Devuelve un mensaje si algo está mal.
func leerEscuderiaForm(r *http.Request) (escuderiaForm, string) {
	f := escuderiaForm{
		ID:        strings.TrimSpace(r.FormValue("id_escuderia")),
		Principal: strings.TrimSpace(r.FormValue("team_principal")),
	}
	if f.ID == "" || f.Principal == "" {
		return f, "El nombre y el Team Principal son obligatorios"
	}

	titulos, err := strconv.Atoi(r.FormValue("titulos_constructores"))
	if err != nil || titulos < 0 {
		return f, "Los títulos tienen que ser un número mayor o igual a 0"
	}
	puntos, err := strconv.Atoi(r.FormValue("puntos_temporada"))
	if err != nil || puntos < 0 {
		return f, "Los puntos tienen que ser un número mayor o igual a 0"
	}
	fecha, err := time.Parse("2006-01-02", r.FormValue("fecha_fundacion"))
	if err != nil {
		return f, "La fecha de fundación no es válida"
	}

	f.Titulos = int32(titulos)
	f.Puntos = int32(puntos)
	f.Fundacion = fecha
	return f, ""
}

func ordenarEscuderias(es []sqlc.Escuderium, orden string) {
	sort.Slice(es, func(i, j int) bool {
		if orden == "desc" {
			return es[i].IDEscuderia > es[j].IDEscuderia
		}
		return es[i].IDEscuderia < es[j].IDEscuderia
	})
}

func atoiOCero(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}

// usuarioEsAdmin dice si la sesión actual es de un administrador (el rol sale de la base).
// Si ya tenés un helper así para pilotos, usá ese y borrá este.
func usuarioEsAdmin(r *http.Request, q *sqlc.Queries) bool {
	usuario, ok := usuarioActual(r)
	if !ok {
		return false
	}
	esAdmin, err := q.EsAdmin(r.Context(), usuario)
	return err == nil && esAdmin
}
