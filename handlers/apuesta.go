package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"
	"time"

	sqlc "Proyecto_web/db/sqlc"
	"Proyecto_web/views"
)

// ApuestaHandler agrupa los handlers de apuesta y guarda
// la conexión a la base de datos y las queries de sqlc.
type ApuestaHandler struct {
	DB      *sql.DB
	queries *sqlc.Queries
}

// NewApuestaHandler es el constructor del handler.
func NewApuestaHandler(db *sql.DB) *ApuestaHandler {
	return &ApuestaHandler{
		DB:      db,
		queries: sqlc.New(db),
	}
}

func (h *ApuestaHandler) PaginaApuestas(w http.ResponseWriter, r *http.Request) {
	usuario, ok := usuarioActual(r) // el usuario SIEMPRE sale de la cookie
	if !ok {
		http.Redirect(w, r, "/ingresar", http.StatusSeeOther)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.renderApuestas(w, r, usuario, "")

	case http.MethodPost:
		h.crearApuestaDesdeForm(w, r, usuario)

	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

func (h *ApuestaHandler) crearApuestaDesdeForm(w http.ResponseWriter, r *http.Request, usuario string) {
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.renderApuestas(w, r, usuario, "Formulario inválido")
		return
	}

	q := h.queries
	ctx := r.Context()

	// 1) Gran premio: "id|YYYY-MM-DD". Se valida contra los que realmente están disponibles.
	idGP, fechaStr, ok := strings.Cut(r.FormValue("gran_premio"), "|")
	fechaCarrera, errFecha := time.Parse("2006-01-02", fechaStr)
	if !ok || errFecha != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.renderApuestas(w, r, usuario, "Elegí un gran premio")
		return
	}

	gps, err := q.ListarGranPremiosProximos(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.renderApuestas(w, r, usuario, "Error interno")
		return
	}
	disponible := false
	for _, gp := range gps {
		if gp.IDGranPremio == idGP && gp.FechaCarrera.Format("2006-01-02") == fechaStr {
			disponible = true
			break
		}
	}
	if !disponible {
		w.WriteHeader(http.StatusBadRequest)
		h.renderApuestas(w, r, usuario, "Ese gran premio no está disponible")
		return
	}

	// 2) Pilotos: exactamente 10, existentes y sin repetir.
	pilotosDB, err := q.ListarPilotosParaApuesta(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.renderApuestas(w, r, usuario, "Error interno")
		return
	}
	existentes := make(map[int32]bool, len(pilotosDB))
	for _, p := range pilotosDB {
		existentes[p.IDPiloto] = true
	}

	valores := r.Form["piloto"]
	if len(valores) != 10 {
		w.WriteHeader(http.StatusBadRequest)
		h.renderApuestas(w, r, usuario, "Tenés que elegir 10 pilotos")
		return
	}

	prediccion := make([]int32, 0, 10)
	vistos := map[int32]bool{}
	for _, v := range valores {
		n, err := strconv.Atoi(v)
		id := int32(n)
		if err != nil || !existentes[id] {
			w.WriteHeader(http.StatusBadRequest)
			h.renderApuestas(w, r, usuario, "Piloto inválido")
			return
		}
		if vistos[id] {
			w.WriteHeader(http.StatusBadRequest)
			h.renderApuestas(w, r, usuario, "No podés repetir pilotos")
			return
		}
		vistos[id] = true
		prediccion = append(prediccion, id)
	}

	// 3) Guardar
	_, err = q.CrearApuesta(ctx, sqlc.CrearApuestaParams{
		NombreUsuario: usuario,
		IDGranPremio:  idGP,
		FechaCarrera:  fechaCarrera,
		Prediccion:    prediccion,
		FechaApuesta:  time.Now(),
	})
	if err != nil {
		if strings.Contains(err.Error(), "unq_usuario_gran_premio_carrera") {
			w.WriteHeader(http.StatusConflict)
			h.renderApuestas(w, r, usuario, "Ya apostaste en ese gran premio")
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		h.renderApuestas(w, r, usuario, "No se pudo guardar la apuesta")
		return
	}

	// Redirect para que F5 no reenvíe el formulario.
	http.Redirect(w, r, "/apuestas", http.StatusSeeOther)
}

func (h *ApuestaHandler) renderApuestas(w http.ResponseWriter, r *http.Request, usuario, mensaje string) {
	q := h.queries
	ctx := r.Context()
	data := views.ApuestasData{Usuario: usuario, Mensaje: mensaje}

	if gps, err := q.ListarGranPremiosProximos(ctx); err == nil {
		for _, gp := range gps {
			fecha := gp.FechaCarrera.Format("2006-01-02")
			data.GranPremios = append(data.GranPremios, views.OpcionGP{
				Valor: gp.IDGranPremio + "|" + fecha,
				Texto: gp.IDGranPremio + " — " + gp.FechaCarrera.Format("02/01/2006"),
			})
		}
	}
	// Pilotos de la temporada para los selects
	if ps, err := q.ListarPilotosParaApuesta(ctx); err == nil {
		for _, p := range ps {
			data.Pilotos = append(data.Pilotos, views.OpcionPiloto{ID: int(p.IDPiloto), Nombre: p.Nombre})
		}
	}

	// Nombres de todos los pilotos, para mostrar apuestas viejas
	nombres := map[int32]string{}
	if ps, err := q.ListarPilotosHistoricos(ctx); err == nil {
		for _, p := range ps {
			nombres[p.IDPiloto] = p.Nombre
		}
	}

	if aps, err := q.ListarApuestasPorUsuario(ctx, usuario); err == nil {
		for _, a := range aps {
			pos := make([]string, 0, len(a.Prediccion))
			for _, id := range a.Prediccion {
				pos = append(pos, nombres[id])
			}
			data.Apuestas = append(data.Apuestas, views.ApuestaVista{
				GranPremio: a.IDGranPremio,
				Fecha:      a.FechaCarrera.Format("02/01/2006"),
				Posiciones: pos,
			})
		}
	}

	views.ApuestasPage(data).Render(ctx, w)
}
