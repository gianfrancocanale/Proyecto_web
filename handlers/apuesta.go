package handlers

import (
	"context"
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

// ---------------------------------------------------------------
// Página principal: listar y crear
// ---------------------------------------------------------------

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

// mostrarError devuelve la lista con un mensaje arriba.
// Responde 200 a propósito: HTMX no reemplaza el contenido con 4xx/5xx.
func (h *ApuestaHandler) mostrarError(w http.ResponseWriter, r *http.Request, usuario, mensaje string) {
	h.renderLista(w, r, usuario, mensaje)
}

func (h *ApuestaHandler) crearApuestaDesdeForm(w http.ResponseWriter, r *http.Request, usuario string) {
	if err := r.ParseForm(); err != nil {
		h.mostrarError(w, r, usuario, "Formulario inválido")
		return
	}

	ctx := r.Context()

	// 1) Gran premio: "id|YYYY-MM-DD". Se valida contra los que están abiertos.
	idGP, fechaStr, ok := strings.Cut(r.FormValue("gran_premio"), "|")
	fechaCarrera, errFecha := time.Parse("2006-01-02", fechaStr)
	if !ok || errFecha != nil {
		h.mostrarError(w, r, usuario, "Elegí un gran premio")
		return
	}
	if !h.carreraAbierta(ctx, idGP, fechaCarrera) {
		h.mostrarError(w, r, usuario, "Ese gran premio no está disponible")
		return
	}

	// 2) Pilotos: exactamente 10, existentes y sin repetir.
	prediccion, msg := h.validarPrediccion(r)
	if msg != "" {
		h.mostrarError(w, r, usuario, msg)
		return
	}

	// 3) Guardar
	_, err := h.queries.CrearApuesta(ctx, sqlc.CrearApuestaParams{
		NombreUsuario: usuario,
		IDGranPremio:  idGP,
		FechaCarrera:  fechaCarrera,
		Prediccion:    prediccion,
		FechaApuesta:  time.Now(),
	})
	if err != nil {
		if strings.Contains(err.Error(), "unq_usuario_gran_premio_carrera") {
			h.mostrarError(w, r, usuario, "Ya apostaste en ese gran premio")
			return
		}
		log.Println("CrearApuesta:", err)
		h.mostrarError(w, r, usuario, "No se pudo guardar la apuesta")
		return
	}

	// Se vuelve a consultar la lista y se devuelve solo ese fragmento.
	h.renderLista(w, r, usuario, "")
}

// datosApuestas arma todo lo que necesitan la página y la lista.
func (h *ApuestaHandler) datosApuestas(ctx context.Context, usuario, mensaje string) views.ApuestasData {
	q := h.queries
	data := views.ApuestasData{Usuario: usuario, Mensaje: mensaje}

	// Carreras abiertas: para el select y para saber qué apuestas se pueden modificar
	abiertas := map[string]bool{}
	gps, err := q.ListarGranPremiosProximos(ctx)
	if err != nil {
		log.Println("ListarGranPremiosProximos:", err)
	}
	for _, gp := range gps {
		clave := gp.IDGranPremio + "|" + gp.FechaCarrera.Format("2006-01-02")
		abiertas[clave] = true
		data.GranPremios = append(data.GranPremios, views.OpcionGP{
			Valor: clave,
			Texto: gp.IDGranPremio + " — " + gp.FechaCarrera.Format("02/01/2006"),
		})
	}

	data.Pilotos = h.opcionesPilotos(ctx)

	// Nombres de todos los pilotos, para mostrar apuestas viejas
	nombres := map[int32]string{}
	hist, err := q.ListarPilotosHistoricos(ctx)
	if err != nil {
		log.Println("ListarPilotosHistoricos:", err)
	}
	for _, p := range hist {
		nombres[p.IDPiloto] = p.Nombre
	}

	aps, err := q.ListarApuestasPorUsuario(ctx, usuario)
	if err != nil {
		log.Println("ListarApuestasPorUsuario:", err)
	}
	for _, a := range aps {
		pos := make([]string, 0, len(a.Prediccion))
		for _, id := range a.Prediccion {
			pos = append(pos, nombres[id])
		}
		clave := a.IDGranPremio + "|" + a.FechaCarrera.Format("2006-01-02")
		data.Apuestas = append(data.Apuestas, views.ApuestaVista{
			ID:         int(a.IDApuesta),
			GranPremio: a.IDGranPremio,
			Fecha:      a.FechaCarrera.Format("02/01/2006"),
			Posiciones: pos,
			Editable:   abiertas[clave],
		})
	}

	return data
}

// renderApuestas devuelve la página completa (solo para el GET).
func (h *ApuestaHandler) renderApuestas(w http.ResponseWriter, r *http.Request, usuario, mensaje string) {
	data := h.datosApuestas(r.Context(), usuario, mensaje)
	views.ApuestasPage(data).Render(r.Context(), w)
}

// renderLista devuelve solo el fragmento de la lista, para HTMX.
func (h *ApuestaHandler) renderLista(w http.ResponseWriter, r *http.Request, usuario, mensaje string) {
	data := h.datosApuestas(r.Context(), usuario, mensaje)
	views.ApuestasLista(data.Apuestas, data.Mensaje).Render(r.Context(), w)
}

// ---------------------------------------------------------------
// Modificar una apuesta: /apuestas/modificar
// ---------------------------------------------------------------

func (h *ApuestaHandler) ModificarApuesta(w http.ResponseWriter, r *http.Request) {
	usuario, ok := usuarioActual(r)
	if !ok {
		http.Redirect(w, r, "/ingresar", http.StatusSeeOther)
		return
	}

	switch r.Method {
	case http.MethodGet:
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil {
			http.Redirect(w, r, "/apuestas", http.StatusSeeOther)
			return
		}
		h.renderModificar(w, r, usuario, int32(id), "", http.StatusOK)

	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			http.Redirect(w, r, "/apuestas", http.StatusSeeOther)
			return
		}
		id, err := strconv.Atoi(r.FormValue("id_apuesta"))
		if err != nil {
			http.Redirect(w, r, "/apuestas", http.StatusSeeOther)
			return
		}

		prediccion, msg := h.validarPrediccion(r)
		if msg != "" {
			h.renderModificar(w, r, usuario, int32(id), msg, http.StatusBadRequest)
			return
		}

		_, err = h.queries.ActualizarPrediccionApuesta(r.Context(), sqlc.ActualizarPrediccionApuestaParams{
			IDApuesta:     int32(id),
			NombreUsuario: usuario, // de la cookie, nunca del formulario
			Prediccion:    prediccion,
			FechaApuesta:  time.Now(),
		})
		if errors.Is(err, sql.ErrNoRows) {
			// No es tuya, no existe o la carrera ya se cerró
			http.Redirect(w, r, "/apuestas", http.StatusSeeOther)
			return
		}
		if err != nil {
			log.Println("ActualizarPrediccionApuesta:", err)
			h.renderModificar(w, r, usuario, int32(id), "No se pudo guardar el cambio", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/apuestas", http.StatusSeeOther)

	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

func (h *ApuestaHandler) renderModificar(w http.ResponseWriter, r *http.Request, usuario string, id int32, mensaje string, status int) {
	ctx := r.Context()

	// Solo trae la apuesta si es de este usuario
	apuesta, err := h.queries.RecuperarApuestaPorIdUsuario(ctx, sqlc.RecuperarApuestaPorIdUsuarioParams{
		NombreUsuario: usuario,
		IDApuesta:     id,
	})
	if err != nil {
		http.Redirect(w, r, "/apuestas", http.StatusSeeOther)
		return
	}

	// Y solo si la carrera sigue abierta
	if !h.carreraAbierta(ctx, apuesta.IDGranPremio, apuesta.FechaCarrera) {
		http.Redirect(w, r, "/apuestas", http.StatusSeeOther)
		return
	}

	seleccion := make([]int, 0, len(apuesta.Prediccion))
	for _, p := range apuesta.Prediccion {
		seleccion = append(seleccion, int(p))
	}

	data := views.ModificarData{
		Usuario:         usuario,
		Mensaje:         mensaje,
		ID:              int(id),
		GranPremioTexto: apuesta.IDGranPremio + " — " + apuesta.FechaCarrera.Format("02/01/2006"),
		Pilotos:         h.opcionesPilotos(ctx),
		Seleccion:       seleccion,
	}

	w.WriteHeader(status)
	views.ModificarApuestaPage(data).Render(ctx, w)
}

// ---------------------------------------------------------------
// Reglas del negocio
// ---------------------------------------------------------------

// carreraAbierta dice si todavía se puede apostar (o modificar) en esa edición.
func (h *ApuestaHandler) carreraAbierta(ctx context.Context, idGP string, fecha time.Time) bool {
	gps, err := h.queries.ListarGranPremiosProximos(ctx)
	if err != nil {
		log.Println("ListarGranPremiosProximos:", err)
		return false
	}
	f := fecha.Format("2006-01-02")
	for _, gp := range gps {
		if gp.IDGranPremio == idGP && gp.FechaCarrera.Format("2006-01-02") == f {
			return true
		}
	}
	return false
}

// opcionesPilotos devuelve los pilotos de la temporada para los selects.
func (h *ApuestaHandler) opcionesPilotos(ctx context.Context) []views.OpcionPiloto {
	ps, err := h.queries.ListarPilotosParaApuesta(ctx)
	if err != nil {
		log.Println("ListarPilotosParaApuesta:", err)
		return nil
	}
	out := make([]views.OpcionPiloto, 0, len(ps))
	for _, p := range ps {
		out = append(out, views.OpcionPiloto{ID: int(p.IDPiloto), Nombre: p.Nombre})
	}
	return out
}

// validarPrediccion devuelve los 10 id_piloto en orden, o un mensaje de error.
// Necesita que r.ParseForm() ya se haya llamado.
func (h *ApuestaHandler) validarPrediccion(r *http.Request) ([]int32, string) {
	pilotosDB, err := h.queries.ListarPilotosParaApuesta(r.Context())
	if err != nil {
		log.Println("ListarPilotosParaApuesta:", err)
		return nil, "Error interno"
	}
	existentes := make(map[int32]bool, len(pilotosDB))
	for _, p := range pilotosDB {
		existentes[p.IDPiloto] = true
	}

	valores := r.Form["piloto"]
	if len(valores) != 10 {
		return nil, "Tenés que elegir 10 pilotos"
	}

	prediccion := make([]int32, 0, 10)
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
		prediccion = append(prediccion, id)
	}
	return prediccion, ""
}

// ---------------------------------------------------------------
// Eliminar una apuesta
// ---------------------------------------------------------------

// DELETE /apuestas/{id}
// Responde 200 vacío para que HTMX reemplace la tarjeta por nada.
// Si no se puede borrar responde con un error 4xx: HTMX no toca el DOM y la tarjeta queda.
func (h *ApuestaHandler) EliminarApuesta(w http.ResponseWriter, r *http.Request) {
	usuario, ok := usuarioActual(r)
	if !ok {
		http.Error(w, "No autenticado", http.StatusUnauthorized)
		return
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Id inválido", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	// La apuesta tiene que ser del usuario de la cookie...
	apuesta, err := h.queries.RecuperarApuestaPorIdUsuario(ctx, sqlc.RecuperarApuestaPorIdUsuarioParams{
		NombreUsuario: usuario,
		IDApuesta:     int32(id),
	})
	if err != nil {
		http.Error(w, "Apuesta no encontrada", http.StatusNotFound)
		return
	}

	// ...y la carrera tiene que seguir abierta (regla de negocio).
	if !h.carreraAbierta(ctx, apuesta.IDGranPremio, apuesta.FechaCarrera) {
		http.Error(w, "La carrera ya se cerró", http.StatusConflict)
		return
	}

	err = h.queries.EliminarApuesta(ctx, sqlc.EliminarApuestaParams{
		IDApuesta:     int32(id),
		NombreUsuario: usuario,
	})
	if err != nil {
		log.Println("EliminarApuesta:", err)
		http.Error(w, "No se pudo borrar", http.StatusInternalServerError)
		return
	}

	// 200 con cuerpo vacío: HTMX reemplaza el target por nada.
	w.WriteHeader(http.StatusOK)
}
