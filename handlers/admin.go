package handlers

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	sqlc "Proyecto_web/db/sqlc"
	"Proyecto_web/views"
)

// Puntaje de las apuestas: editá estos valores para cambiar la regla.
const (
	puntosPosicionExacta = 3 // el piloto terminó justo en la posición apostada
	puntosEnTop10        = 1 // terminó en el top 10, pero en otra posición
)

type AdminHandler struct {
	DB      *sql.DB
	queries *sqlc.Queries
}

func NewAdminHandler(db *sql.DB) *AdminHandler {
	return &AdminHandler{DB: db, queries: sqlc.New(db)}
}

// RequiereAdmin: exige sesión Y rol de admin, consultado en la base.
func (h *AdminHandler) RequiereAdmin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		usuario, ok := usuarioActual(r)
		if !ok {
			http.Redirect(w, r, "/ingresar", http.StatusSeeOther)
			return
		}
		esAdmin, err := h.queries.EsAdmin(r.Context(), usuario)
		if err != nil || !esAdmin {
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

	circuitos, err := h.queries.ListarCircuitosParaAdmin(ctx)
	if err != nil {
		log.Println("ListarCircuitosParaAdmin:", err)
	}
	data.Circuitos = circuitos

	pend, err := h.queries.ListarEdicionesPendientes(ctx)
	if err != nil {
		log.Println("ListarEdicionesPendientes:", err)
	}
	for _, p := range pend {
		data.Pendientes = append(data.Pendientes, views.OpcionGP{
			Valor: p.IDGranPremio + "|" + p.FechaCarrera.Format("2006-01-02"),
			Texto: p.IDGranPremio + " — " + p.FechaCarrera.Format("02/01/2006"),
		})
	}

	pilotos, err := h.queries.ListarPilotosParaApuesta(ctx)
	if err != nil {
		log.Println("ListarPilotosParaApuesta:", err)
	}
	for _, p := range pilotos {
		data.Pilotos = append(data.Pilotos, views.OpcionPiloto{ID: int(p.IDPiloto), Nombre: p.Nombre})
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
	if circuito == "" || err != nil {
		h.renderAdmin(w, r, "", "Elegí un circuito y una fecha válida", http.StatusBadRequest)
		return
	}

	_, err = h.queries.CrearEdicionGranPremio(r.Context(), sqlc.CrearEdicionGranPremioParams{
		IDGranPremio: circuito,
		FechaCarrera: fecha,
	})
	if err != nil {
		if strings.Contains(err.Error(), "pk_gran_premio_historico") {
			h.renderAdmin(w, r, "", "Esa carrera ya está programada", http.StatusConflict)
			return
		}
		log.Println("CrearEdicionGranPremio:", err)
		h.renderAdmin(w, r, "", "No se pudo programar la carrera", http.StatusInternalServerError)
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
	fecha, errFecha := time.Parse("2006-01-02", fechaStr)
	if !ok || errFecha != nil {
		h.renderAdmin(w, r, "", "Elegí una carrera", http.StatusBadRequest)
		return
	}

	orden, msg := h.validarResultado(r)
	if msg != "" {
		h.renderAdmin(w, r, "", msg, http.StatusBadRequest)
		return
	}

	// Todo o nada: resultado + puntos de todos los usuarios en una transacción.
	tx, err := h.DB.BeginTx(ctx, nil)
	if err != nil {
		log.Println("BeginTx:", err)
		h.renderAdmin(w, r, "", "Error interno", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()
	qtx := h.queries.WithTx(tx)

	_, err = qtx.CargarResultadoCarrera(ctx, sqlc.CargarResultadoCarreraParams{
		IDGranPremio:     idGP,
		FechaCarrera:     fecha,
		ResultadoCarrera: orden,
	})
	if errors.Is(err, sql.ErrNoRows) {
		h.renderAdmin(w, r, "", "Esa carrera no existe o ya tiene resultado", http.StatusConflict)
		return
	}
	if err != nil {
		log.Println("CargarResultadoCarrera:", err)
		h.renderAdmin(w, r, "", "No se pudo cargar el resultado", http.StatusInternalServerError)
		return
	}

	apuestas, err := qtx.ListarApuestasPorGranPremio(ctx, sqlc.ListarApuestasPorGranPremioParams{
		IDGranPremio: idGP,
		FechaCarrera: fecha,
	})
	if err != nil {
		log.Println("ListarApuestasPorGranPremio:", err)
		h.renderAdmin(w, r, "", "No se pudieron leer las apuestas", http.StatusInternalServerError)
		return
	}

	for _, a := range apuestas {
		pts := calcularPuntos(a.Prediccion, orden)
		if pts == 0 {
			continue
		}
		err := qtx.SumarPuntosUsuario(ctx, sqlc.SumarPuntosUsuarioParams{
			NombreUsuario: a.NombreUsuario,
			Puntos:        int32(pts),
		})
		if err != nil {
			log.Println("SumarPuntosUsuario:", err)
			h.renderAdmin(w, r, "", "No se pudieron sumar los puntos", http.StatusInternalServerError)
			return
		}
	}

	if err := tx.Commit(); err != nil {
		log.Println("Commit:", err)
		h.renderAdmin(w, r, "", "No se pudo guardar", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin?ok=resultado", http.StatusSeeOther)
}

// validarResultado: exactamente 22 pilotos de la temporada, sin repetir, en orden de llegada.
func (h *AdminHandler) validarResultado(r *http.Request) ([]int32, string) {
	pilotos, err := h.queries.ListarPilotosParaApuesta(r.Context())
	if err != nil {
		log.Println("ListarPilotosParaApuesta:", err)
		return nil, "Error interno"
	}
	existentes := make(map[int32]bool, len(pilotos))
	for _, p := range pilotos {
		existentes[p.IDPiloto] = true
	}

	valores := r.Form["piloto"]
	if len(valores) != 22 {
		return nil, "Tenés que ordenar los 22 pilotos"
	}

	orden := make([]int32, 0, 22)
	vistos := map[int32]bool{}
	for _, v := range valores {
		n, err := strconv.Atoi(v)
		id := int32(n)
		if err != nil || !existentes[id] {
			return nil, "Piloto inválido"
		}
		if vistos[id] {
			return nil, "No podés repetir pilotos"
		}
		vistos[id] = true
		orden = append(orden, id)
	}
	return orden, ""
}

// calcularPuntos compara el top 10 apostado contra los primeros 10 del resultado.
func calcularPuntos(prediccion, resultado []int32) int {
	top10 := resultado
	if len(top10) > 10 {
		top10 = top10[:10]
	}
	estaEnTop10 := make(map[int32]bool, len(top10))
	for _, id := range top10 {
		estaEnTop10[id] = true
	}

	pts := 0
	for i, id := range prediccion {
		switch {
		case i < len(top10) && top10[i] == id:
			pts += puntosPosicionExacta
		case estaEnTop10[id]:
			pts += puntosEnTop10
		}
	}
	return pts
}
