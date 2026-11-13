package handlers

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	sqlc "Proyecto_web/db/sqlc"
	vista "Proyecto_web/views"
)

// PilotoHandler agrupa los handlers de piloto y guarda
// la conexión a la base de datos y las queries de sqlc.
type PilotoHandler struct {
	DB      *sql.DB
	queries *sqlc.Queries
}

// NewPilotoHandler es el constructor del handler.
func NewPilotoHandler(db *sql.DB) *PilotoHandler {
	return &PilotoHandler{
		DB:      db,
		queries: sqlc.New(db),
	}
}

// Admin
func (h *PilotoHandler) esAdmin(ctx context.Context, usuario string) bool {
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

func (h *PilotoHandler) Pagina(w http.ResponseWriter, r *http.Request) {
	usuario, _ := usuarioActual(r) // listar es público: no redirige si no hay sesión

	switch r.Method {
	case http.MethodGet:
		orden := r.URL.Query().Get("orden")
		if orden != "asc" {
			orden = "desc"
		}
		data := h.datosPilotos(r.Context(), usuario, orden, "")
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
	ctx := r.Context()

	// 1) Número de piloto
	n, err := strconv.Atoi(strings.TrimSpace(r.FormValue("id_piloto")))
	if err != nil || n <= 0 {
		h.mostrarError(w, r, usuario, "Número de piloto inválido")
		return
	}
	idPiloto := int32(n)

	// 2) Escudería: tiene que existir
	idEscuderia := r.FormValue("id_escuderia")
	if _, err := h.queries.RecuperarEscuderia(ctx, idEscuderia); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.mostrarError(w, r, usuario, "No existe la escudería")
		} else {
			log.Println("RecuperarEscuderia:", err)
			h.mostrarError(w, r, usuario, "Error al verificar la escudería")
		}
		return
	}

	// 3) No puede estar ya en la temporada
	if h.pilotoEnTemporada(ctx, idPiloto) {
		h.mostrarError(w, r, usuario, "Ese piloto ya está en la temporada")
		return
	}

	// Todo o nada: histórico (si hace falta) + temporada
	tx, err := h.DB.BeginTx(ctx, nil)
	if err != nil {
		log.Println("BeginTx:", err)
		h.mostrarError(w, r, usuario, "Error interno")
		return
	}
	defer tx.Rollback()
	qtx := h.queries.WithTx(tx)

	// 4) Si el piloto no existe en el histórico, se crea con los datos del formulario
	_, err = qtx.RecuperarPilotoHistorico(ctx, idPiloto)
	if errors.Is(err, sql.ErrNoRows) {
		nombre := strings.TrimSpace(r.FormValue("nombre"))
		pais := strings.TrimSpace(r.FormValue("pais"))
		nacimiento, errFecha := time.Parse("2006-01-02", r.FormValue("fecha_nacimiento"))
		titulos := 0
		if v := strings.TrimSpace(r.FormValue("titulos_ganados")); v != "" {
			titulos, err = strconv.Atoi(v)
			if err != nil || titulos < 0 {
				h.mostrarError(w, r, usuario, "Títulos ganados inválidos")
				return
			}
		}

		if nombre == "" || pais == "" || errFecha != nil {
			h.mostrarError(w, r, usuario, "Piloto nuevo: completá nombre, país y fecha de nacimiento")
			return
		}

		_, err = qtx.CrearPilotoHistorico(ctx, sqlc.CrearPilotoHistoricoParams{
			IDPiloto:        idPiloto,
			Nombre:          nombre,
			Pais:            pais,
			FechaNacimiento: nacimiento,
			TitulosGanados:  int32(titulos),
		})
		if err != nil {
			log.Println("CrearPilotoHistorico:", err)
			h.mostrarError(w, r, usuario, "No se pudo crear el piloto")
			return
		}
	} else if err != nil {
		log.Println("RecuperarPilotoHistorico:", err)
		h.mostrarError(w, r, usuario, "Error al verificar el piloto")
		return
	}

	// 5) Inscribirlo en la temporada (arranca con 0 puntos)
	_, err = qtx.CrearPilotoTemporada(ctx, sqlc.CrearPilotoTemporadaParams{
		IDPiloto:        idPiloto,
		IDEscuderia:     idEscuderia,
		PuntosTemporada: 0,
	})
	if err != nil {
		log.Println("CrearPilotoTemporada:", err)
		h.mostrarError(w, r, usuario, "No se pudo agregar a la temporada")
		return
	}

	if err := tx.Commit(); err != nil {
		log.Println("Commit:", err)
		h.mostrarError(w, r, usuario, "No se pudo guardar")
		return
	}

	h.renderLista(w, r, usuario, "")
}

// ---------------------------------------------------------------
// Datos y render
// ---------------------------------------------------------------

// listaOrdenada trae los pilotos de la temporada ordenados por puntos.
// Si no hay ninguno devuelve un slice vacío (no es un error).
func (h *PilotoHandler) listaOrdenada(ctx context.Context, orden string) []sqlc.ListarPilotosRow {
	pilotos, err := h.ListarPilotos(ctx)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Println("ListarPilotos:", err)
		}
		return nil
	}
	sort.Slice(pilotos, func(i, j int) bool {
		if orden == "asc" {
			return pilotos[i].PuntosTemporada < pilotos[j].PuntosTemporada
		}
		return pilotos[i].PuntosTemporada > pilotos[j].PuntosTemporada
	})
	return pilotos
}

// datosPilotos arma todo lo que necesitan la página y la lista.
func (h *PilotoHandler) datosPilotos(ctx context.Context, usuario, orden, mensaje string) vista.PilotosData {
	esAdm := h.esAdmin(ctx, usuario)
	data := vista.PilotosData{
		Usuario: usuario,
		EsAdmin: esAdm,
		Orden:   orden,
		Mensaje: mensaje,
		Pilotos: h.listaOrdenada(ctx, orden),
	}

	if esAdm {
		data.Escuderias = h.opcionesEscuderias(ctx)
	}
	return data
}

// renderLista devuelve solo el fragmento de la lista, para HTMX.
func (h *PilotoHandler) renderLista(w http.ResponseWriter, r *http.Request, usuario, mensaje string) {
	data := h.datosPilotos(r.Context(), usuario, "desc", mensaje)
	vista.PilotosLista(data.Pilotos, data.EsAdmin, data.Mensaje).Render(r.Context(), w)
}

// opcionesEscuderias devuelve las escuderías para el select.
func (h *PilotoHandler) opcionesEscuderias(ctx context.Context) []vista.OpcionEscuderia {
	es, err := h.queries.ListarEscuderias(ctx)
	if err != nil {
		log.Println("ListarEscuderias:", err)
		return nil
	}
	out := make([]vista.OpcionEscuderia, 0, len(es))
	for _, e := range es {
		out = append(out, vista.OpcionEscuderia{ID: e.IDEscuderia})
	}
	return out
}

func (h *PilotoHandler) pilotoEnTemporada(ctx context.Context, id int32) bool {
	pilotos, err := h.ListarPilotos(ctx)
	if err != nil {
		return false
	}
	for _, p := range pilotos {
		if p.IDPiloto == id {
			return true
		}
	}
	return false
}

func (h *PilotoHandler) ListarPilotos(ctx context.Context) ([]sqlc.ListarPilotosRow, error) {
	pilotos, err := h.queries.ListarPilotos(ctx)
	if err != nil {
		return nil, err
	}
	if len(pilotos) == 0 {
		return nil, sql.ErrNoRows
	}
	return pilotos, nil
}

// ---------------------------------------------------------------
// Modificar un piloto: /pilotos/modificar  (admin)
// ---------------------------------------------------------------

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
		h.renderModificar(w, r, usuario, int32(id), "", http.StatusOK)

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

		idEscuderia := r.FormValue("id_escuderia")
		if _, err := h.queries.RecuperarEscuderia(r.Context(), idEscuderia); err != nil {
			h.renderModificar(w, r, usuario, int32(id), "No existe la escudería", http.StatusBadRequest)
			return
		}

		puntos, err := strconv.Atoi(r.FormValue("puntos_temporada"))
		if err != nil || puntos < 0 {
			h.renderModificar(w, r, usuario, int32(id), "Puntos inválidos", http.StatusBadRequest)
			return
		}

		_, err = h.queries.ModificarPiloto(r.Context(), sqlc.ModificarPilotoParams{
			IDPiloto:        int32(id),
			IDEscuderia:     idEscuderia,
			PuntosTemporada: int32(puntos),
		})
		if errors.Is(err, sql.ErrNoRows) {
			http.Redirect(w, r, "/pilotos", http.StatusSeeOther)
			return
		}
		if err != nil {
			log.Println("ModificarPiloto:", err)
			h.renderModificar(w, r, usuario, int32(id), "No se pudo guardar el cambio", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/pilotos", http.StatusSeeOther)

	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

func (h *PilotoHandler) renderModificar(w http.ResponseWriter, r *http.Request, usuario string, id int32, mensaje string, status int) {
	ctx := r.Context()

	// Busca el piloto en la temporada
	var actual *sqlc.ListarPilotosRow
	pilotos, _ := h.ListarPilotos(ctx)
	for i := range pilotos {
		if pilotos[i].IDPiloto == id {
			actual = &pilotos[i]
			break
		}
	}
	if actual == nil {
		http.Redirect(w, r, "/pilotos", http.StatusSeeOther)
		return
	}

	data := vista.ModificarPilotoData{
		Usuario:    usuario,
		Mensaje:    mensaje,
		ID:         int(actual.IDPiloto),
		Nombre:     actual.Nombre,
		Escuderia:  actual.IDEscuderia,
		Puntos:     int(actual.PuntosTemporada),
		Escuderias: h.opcionesEscuderias(ctx),
	}

	w.WriteHeader(status)
	vista.ModificarPilotoPage(data).Render(ctx, w)
}

// ---------------------------------------------------------------
// Eliminar un piloto (admin)
// ---------------------------------------------------------------

// DELETE /pilotos/{id}
// Responde 200 vacío para que HTMX reemplace la tarjeta por nada.
// Si no se puede borrar responde 4xx/5xx: HTMX no toca el DOM y la tarjeta queda.
func (h *PilotoHandler) EliminarPiloto(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
	if err != nil {
		http.Error(w, "ID de piloto inválido", http.StatusBadRequest)
		return
	}

	if err := h.queries.EliminarPilotoTemporada(r.Context(), int32(id)); err != nil {
		log.Println("EliminarPilotoTemporada:", err)
		http.Error(w, "No se pudo borrar", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
