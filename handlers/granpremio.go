package handlers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"

	sqlc "Proyecto_web/db/sqlc"
	vista "Proyecto_web/views"
)

// GranPremioHandler agrupa los handlers de gran premio y guarda
// la conexión a la base de datos y las queries de sqlc.
type GranPremioHandler struct {
	DB      *sql.DB
	queries *sqlc.Queries
}

// NewGranPremioHandler es el constructor del handler.
func NewGranPremioHandler(db *sql.DB) *GranPremioHandler {
	return &GranPremioHandler{
		DB:      db,
		queries: sqlc.New(db),
	}
}

// esAdmin usa la misma regla que AdminHandler.RequiereAdmin: el rol sale de la base.
func (h *GranPremioHandler) esAdmin(ctx context.Context, usuario string) bool {
	if usuario == "" {
		return false
	}
	ok, err := h.queries.EsAdmin(ctx, usuario)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Println("EsAdmin:", err)
		}
		return false
	}
	return ok
}

// ---------------------------------------------------------------
// Página principal: listar (todos) y crear (admin)
// ---------------------------------------------------------------

func (h *GranPremioHandler) Pagina(w http.ResponseWriter, r *http.Request) {
	usuario, _ := usuarioActual(r) // listar es público

	switch r.Method {
	case http.MethodGet:
		orden := r.URL.Query().Get("orden")
		if orden != "desc" {
			orden = "asc"
		}
		data := h.datosGranPremios(r.Context(), usuario, orden, "")
		vista.GranPremiosPage(data).Render(r.Context(), w)

	case http.MethodPost:
		// RequiereAdmin ya validó el rol en main.go
		h.crearDesdeForm(w, r, usuario)

	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

// mostrarError devuelve la lista con un mensaje arriba.
// Responde 200 a propósito: así HTMX siempre reemplaza la lista.
func (h *GranPremioHandler) mostrarError(w http.ResponseWriter, r *http.Request, usuario, mensaje string) {
	h.renderLista(w, r, usuario, mensaje)
}

func (h *GranPremioHandler) crearDesdeForm(w http.ResponseWriter, r *http.Request, usuario string) {
	if err := r.ParseForm(); err != nil {
		h.mostrarError(w, r, usuario, "Formulario inválido")
		return
	}
	ctx := r.Context()

	id := strings.TrimSpace(r.FormValue("id_gran_premio"))
	pais := strings.TrimSpace(r.FormValue("pais"))
	if id == "" || pais == "" {
		h.mostrarError(w, r, usuario, "Completá el nombre del gran premio y el país")
		return
	}

	vueltas, err := strconv.Atoi(strings.TrimSpace(r.FormValue("cantidad_vueltas")))
	if err != nil || vueltas <= 0 {
		h.mostrarError(w, r, usuario, "Cantidad de vueltas inválida")
		return
	}

	ganador, err := strconv.Atoi(r.FormValue("id_ultimo_ganador"))
	if err != nil {
		h.mostrarError(w, r, usuario, "Elegí el último ganador")
		return
	}
	if _, err := h.queries.RecuperarPilotoHistorico(ctx, int32(ganador)); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.mostrarError(w, r, usuario, "El último ganador no existe")
		} else {
			log.Println("RecuperarPilotoHistorico:", err)
			h.mostrarError(w, r, usuario, "Error al verificar el último ganador")
		}
		return
	}

	if h.existeGranPremio(ctx, id) {
		h.mostrarError(w, r, usuario, "Ya existe un gran premio con ese nombre")
		return
	}

	params := sqlc.CrearGranPremioParams{
		IDGranPremio:    id,
		Pais:            pais,
		CantidadVueltas: int32(vueltas),
		IDUltimoGanador: int32(ganador),
	}
	if !asignarLongitud(r.FormValue("longitud_km"), &params.LongitudKm) {
		h.mostrarError(w, r, usuario, "Longitud inválida")
		return
	}

	if _, err := h.queries.CrearGranPremio(ctx, params); err != nil {
		log.Println("CrearGranPremio:", err)
		h.mostrarError(w, r, usuario, "No se pudo crear el gran premio")
		return
	}

	h.renderLista(w, r, usuario, "")
}

// ---------------------------------------------------------------
// Datos y render
// ---------------------------------------------------------------

// listaOrdenada trae los grandes premios ordenados por nombre.
// Si no hay ninguno devuelve un slice vacío (no es un error).
func (h *GranPremioHandler) listaOrdenada(ctx context.Context, orden string) []sqlc.ListarGrandesPremiosRow {
	gps, err := h.ListarGranPremios(ctx)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Println("ListarGrandesPremios:", err)
		}
		return nil
	}
	sort.Slice(gps, func(i, j int) bool {
		if orden == "desc" {
			return gps[i].IDGranPremio > gps[j].IDGranPremio
		}
		return gps[i].IDGranPremio < gps[j].IDGranPremio
	})
	return gps
}

func (h *GranPremioHandler) datosGranPremios(ctx context.Context, usuario, orden, mensaje string) vista.GranPremiosData {
	adm := h.esAdmin(ctx, usuario)
	data := vista.GranPremiosData{
		EsAdmin:     adm,
		Orden:       orden,
		Mensaje:     mensaje,
		GranPremios: h.listaOrdenada(ctx, orden),
	}
	if adm {
		data.Pilotos = h.opcionesPilotos(ctx)
	}
	return data
}

// renderLista devuelve solo el fragmento de la lista, para HTMX.
func (h *GranPremioHandler) renderLista(w http.ResponseWriter, r *http.Request, usuario, mensaje string) {
	data := h.datosGranPremios(r.Context(), usuario, "asc", mensaje)
	vista.GranPremiosLista(data.GranPremios, data.EsAdmin, data.Mensaje).Render(r.Context(), w)
}

// opcionesPilotos: pilotos históricos para elegir el último ganador.
func (h *GranPremioHandler) opcionesPilotos(ctx context.Context) []vista.OpcionPiloto {
	hist, err := h.queries.ListarPilotosHistoricos(ctx)
	if err != nil {
		log.Println("ListarPilotosHistoricos:", err)
		return nil
	}
	out := make([]vista.OpcionPiloto, 0, len(hist))
	for _, p := range hist {
		out = append(out, vista.OpcionPiloto{ID: int(p.IDPiloto), Nombre: p.Nombre})
	}
	return out
}

func (h *GranPremioHandler) existeGranPremio(ctx context.Context, id string) bool {
	gps, err := h.ListarGranPremios(ctx)
	if err != nil {
		return false
	}
	for _, g := range gps {
		if g.IDGranPremio == id {
			return true
		}
	}
	return false
}

func (h *GranPremioHandler) ListarGranPremios(ctx context.Context) ([]sqlc.ListarGrandesPremiosRow, error) {
	granPremios, err := h.queries.ListarGrandesPremios(ctx)
	if err != nil {
		return nil, err
	}
	if len(granPremios) == 0 {
		return nil, sql.ErrNoRows
	}
	return granPremios, nil
}

// asignarLongitud valida que el texto sea un número > 0 y lo escribe en el
// campo de destino (sirve sea int, float o string, según lo que genere sqlc).
func asignarLongitud(texto string, destino any) bool {
	texto = strings.Replace(strings.TrimSpace(texto), ",", ".", 1)
	f, err := strconv.ParseFloat(texto, 64)
	if err != nil || f <= 0 {
		return false
	}
	_, err = fmt.Sscan(texto, destino)
	return err == nil
}

// ---------------------------------------------------------------
// Modificar un gran premio: /gran_premios/modificar  (admin)
// ---------------------------------------------------------------

func (h *GranPremioHandler) ModificarGranPremio(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Redirect(w, r, "/gran_premios", http.StatusSeeOther)
			return
		}
		h.renderModificar(w, r, id, "", http.StatusOK)

	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			http.Redirect(w, r, "/gran_premios", http.StatusSeeOther)
			return
		}
		ctx := r.Context()

		// La identidad del GP viaja en un campo oculto; no se puede cambiar.
		id := strings.TrimSpace(r.FormValue("id_gran_premio"))
		if id == "" {
			http.Redirect(w, r, "/gran_premios", http.StatusSeeOther)
			return
		}

		pais := strings.TrimSpace(r.FormValue("pais"))
		if pais == "" {
			h.renderModificar(w, r, id, "El país es obligatorio", http.StatusBadRequest)
			return
		}
		vueltas, err := strconv.Atoi(strings.TrimSpace(r.FormValue("cantidad_vueltas")))
		if err != nil || vueltas <= 0 {
			h.renderModificar(w, r, id, "Cantidad de vueltas inválida", http.StatusBadRequest)
			return
		}
		ganador, err := strconv.Atoi(r.FormValue("id_ultimo_ganador"))
		if err != nil {
			h.renderModificar(w, r, id, "Elegí el último ganador", http.StatusBadRequest)
			return
		}
		if _, err := h.queries.RecuperarPilotoHistorico(ctx, int32(ganador)); err != nil {
			h.renderModificar(w, r, id, "El último ganador no existe", http.StatusBadRequest)
			return
		}

		params := sqlc.ModificarGranPremioParams{
			IDGranPremio:    id,
			Pais:            pais,
			CantidadVueltas: int32(vueltas),
			IDUltimoGanador: int32(ganador),
		}
		if !asignarLongitud(r.FormValue("longitud_km"), &params.LongitudKm) {
			h.renderModificar(w, r, id, "Longitud inválida", http.StatusBadRequest)
			return
		}

		_, err = h.queries.ModificarGranPremio(ctx, params)
		if errors.Is(err, sql.ErrNoRows) {
			http.Redirect(w, r, "/gran_premios", http.StatusSeeOther)
			return
		}
		if err != nil {
			log.Println("ModificarGranPremio:", err)
			h.renderModificar(w, r, id, "No se pudo guardar el cambio", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/gran_premios", http.StatusSeeOther)

	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

func (h *GranPremioHandler) renderModificar(w http.ResponseWriter, r *http.Request, id, mensaje string, status int) {
	ctx := r.Context()

	var actual *sqlc.ListarGrandesPremiosRow
	gps, _ := h.ListarGranPremios(ctx)
	for i := range gps {
		if gps[i].IDGranPremio == id {
			actual = &gps[i]
			break
		}
	}
	if actual == nil {
		http.Redirect(w, r, "/gran_premios", http.StatusSeeOther)
		return
	}

	data := vista.ModificarGranPremioData{
		Mensaje:  mensaje,
		ID:       actual.IDGranPremio,
		Pais:     actual.Pais,
		Longitud: fmt.Sprint(actual.LongitudKm),
		Vueltas:  fmt.Sprint(actual.CantidadVueltas),
		Ganador:  fmt.Sprint(actual.IDUltimoGanador),
		Pilotos:  h.opcionesPilotos(ctx),
	}

	w.WriteHeader(status)
	vista.ModificarGranPremioPage(data).Render(ctx, w)
}

// ---------------------------------------------------------------
// Eliminar un gran premio (admin)
// ---------------------------------------------------------------

// DELETE /gran_premios/{id}
// Siempre devuelve la lista (200): si falla, con un mensaje arriba.
func (h *GranPremioHandler) EliminarGranPremio(w http.ResponseWriter, r *http.Request) {
	usuario, _ := usuarioActual(r)

	id := r.PathValue("id")
	if id == "" {
		h.renderLista(w, r, usuario, "ID de gran premio inválido")
		return
	}

	if err := h.queries.EliminarGranPremio(r.Context(), id); err != nil {
		log.Println("EliminarGranPremio:", err)
		msg := "No se pudo borrar el gran premio"
		if strings.Contains(strings.ToLower(err.Error()), "foreign key") {
			msg = "No se puede borrar: tiene carreras o apuestas asociadas"
		}
		h.renderLista(w, r, usuario, msg)
		return
	}

	h.renderLista(w, r, usuario, "")
}
