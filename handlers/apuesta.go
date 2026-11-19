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
	"Proyecto_web/views"
)

type ApuestaHandler struct {
	apuestas *aplicacion.ApuestaService
	usuarios *aplicacion.UsuarioService
	carreras *aplicacion.GranPremioHistoricoService
	pilotos  *aplicacion.PilotoService
}

func NewApuestaHandler(
	a *aplicacion.ApuestaService,
	u *aplicacion.UsuarioService,
	c *aplicacion.GranPremioHistoricoService,
	p *aplicacion.PilotoService,
) *ApuestaHandler {
	return &ApuestaHandler{apuestas: a, usuarios: u, carreras: c, pilotos: p}
}

// Página principal: listar y crear
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

func (h *ApuestaHandler) mostrarError(w http.ResponseWriter, r *http.Request, usuario, mensaje string) {
	h.renderLista(w, r, usuario, mensaje)
}

func (h *ApuestaHandler) crearApuestaDesdeForm(w http.ResponseWriter, r *http.Request, usuario string) {
	if err := r.ParseForm(); err != nil {
		h.mostrarError(w, r, usuario, "Formulario inválido")
		return
	}

	idGP, fechaStr, ok := strings.Cut(r.FormValue("gran_premio"), "|")
	fechaCarrera, err := time.Parse("2006-01-02", fechaStr)
	if !ok || err != nil {
		h.mostrarError(w, r, usuario, "Elegí un gran premio")
		return
	}
	prediccion := parsearPrediccion(r)

	if err := h.apuestas.Crear(r.Context(), usuario, idGP, fechaCarrera, prediccion); err != nil {
		h.mostrarError(w, r, usuario, mensajeErrorApuesta(err))
		return
	}
	h.renderLista(w, r, usuario, "")
}

// datosApuestas arma todo lo que necesitan la página y la lista.
func (h *ApuestaHandler) datosApuestas(ctx context.Context, usuario, mensaje string) views.ApuestasData {
	data := views.ApuestasData{Usuario: usuario, Mensaje: mensaje}

	abiertas := map[string]bool{}
	carreras, err := h.apuestas.CarrerasAbiertas(ctx)
	if err != nil {
		log.Println("Listar carreras próximas:", err)
	}
	for _, g := range carreras {
		clave := g.ID + "|" + g.Fecha.Format("2006-01-02")
		abiertas[clave] = true
		data.GranPremios = append(data.GranPremios, views.OpcionGP{
			Valor: clave,
			Texto: g.ID + " — " + g.Fecha.Format("02/01/2006"),
		})
	}

	pilotos, err := h.apuestas.PilotosDisponibles(ctx)
	if err != nil {
		log.Println("Listar pilotos para apuesta:", err)
	}
	for _, p := range pilotos {
		data.Pilotos = append(data.Pilotos, views.OpcionPiloto{ID: p.ID, Nombre: p.Nombre})
	}

	nombres, err := h.apuestas.NombresPilotos(ctx)
	if err != nil {
		log.Println("Listar pilotos históricos:", err)
	}
	aps, err := h.apuestas.ListarPorUsuario(ctx, usuario)
	if err != nil {
		log.Println("Listar apuestas por usuario:", err)
	}
	for _, a := range aps {
		pos := make([]string, 0, len(a.Prediccion))
		for _, id := range a.Prediccion {
			pos = append(pos, nombres[id])
		}
		clave := a.GranPremio + "|" + a.FechaCarrera.Format("2006-01-02")
		data.Apuestas = append(data.Apuestas, views.ApuestaVista{
			ID:         a.ID,
			GranPremio: a.GranPremio,
			Fecha:      a.FechaCarrera.Format("02/01/2006"),
			Posiciones: pos,
			Editable:   abiertas[clave],
		})
	}
	return data
}

func (h *ApuestaHandler) renderApuestas(w http.ResponseWriter, r *http.Request, usuario, mensaje string) {
	ctx := r.Context()
	data := h.datosApuestas(ctx, usuario, mensaje)

	puntos, err := h.usuarios.ObtenerPuntos(ctx, usuario)
	if err != nil {
		log.Println("Obtener puntos:", err)
	}
	data.Puntos = puntos

	// Ranking de jugadores (no admin), ya ordenado por el servicio
	jugadores, err := h.usuarios.ListarJugadores(ctx)
	if err != nil {
		log.Println("Listar jugadores:", err) // no rompemos la página por esto
	}
	for i, j := range jugadores {
		data.Jugadores = append(data.Jugadores, views.JugadorVista{
			Posicion: i + 1,
			Nombre:   j.Nombre,
			Puntos:   j.Puntos,
			EsYo:     j.Nombre == usuario,
		})
	}

	views.ApuestasPage(data).Render(ctx, w)
}

func (h *ApuestaHandler) renderLista(w http.ResponseWriter, r *http.Request, usuario, mensaje string) {
	data := h.datosApuestas(r.Context(), usuario, mensaje)
	views.ApuestasLista(data.Apuestas, data.Mensaje).Render(r.Context(), w)
}

// Modificar una apuesta: /apuestas/modificar
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
		h.renderModificar(w, r, usuario, id, "", http.StatusOK)

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
		if err := h.apuestas.Modificar(r.Context(), id, usuario, parsearPrediccion(r)); err != nil {
			mensaje := "No se pudo guardar el cambio"
			status := http.StatusInternalServerError
			var ev dominio.ErrValidacion
			switch {
			case errors.As(err, &ev):
				mensaje, status = ev.Error(), http.StatusBadRequest
			case errors.Is(err, dominio.ErrCarreraCerrada):
				http.Redirect(w, r, "/apuestas", http.StatusSeeOther)
				return
			case errors.Is(err, dominio.ErrApuestaNoEncontrada):
				http.Redirect(w, r, "/apuestas", http.StatusSeeOther)
				return
			}
			h.renderModificar(w, r, usuario, id, mensaje, status)
			return
		}
		http.Redirect(w, r, "/apuestas", http.StatusSeeOther)

	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

func (h *ApuestaHandler) renderModificar(w http.ResponseWriter, r *http.Request, usuario string, id int, mensaje string, status int) {
	ctx := r.Context()

	a, err := h.apuestas.RecuperarEditable(ctx, id, usuario)
	if err != nil {
		http.Redirect(w, r, "/apuestas", http.StatusSeeOther)
		return
	}

	seleccion := make([]int, 0, len(a.Prediccion))
	for _, p := range a.Prediccion {
		seleccion = append(seleccion, p)
	}
	pilotos, err := h.pilotos.ListarParaApuesta(ctx)
	if err != nil {
		log.Println("Listar pilotos para apuesta:", err)
	}
	opciones := make([]views.OpcionPiloto, 0, len(pilotos))
	for _, p := range pilotos {
		opciones = append(opciones, views.OpcionPiloto{ID: p.ID, Nombre: p.Nombre})
	}

	data := views.ModificarData{
		Usuario:         usuario,
		Mensaje:         mensaje,
		ID:              id,
		GranPremioTexto: a.GranPremio + " — " + a.FechaCarrera.Format("02/01/2006"),
		Pilotos:         opciones,
		Seleccion:       seleccion,
	}
	w.WriteHeader(status)
	views.ModificarApuestaPage(data).Render(ctx, w)
}

// parsearPrediccion convierte los valores del form a IDs (las reglas las valida el dominio).
func parsearPrediccion(r *http.Request) []int {
	valores := r.Form["piloto"]
	prediccion := make([]int, 0, len(valores))
	for _, v := range valores {
		n, err := strconv.Atoi(v)
		if err != nil {
			return nil
		}
		prediccion = append(prediccion, n)
	}
	return prediccion
}

func mensajeErrorApuesta(err error) string {
	var ev dominio.ErrValidacion
	switch {
	case errors.As(err, &ev):
		return ev.Error()
	case errors.Is(err, dominio.ErrApuestaDuplicada):
		return "Ya apostaste en ese gran premio"
	case errors.Is(err, dominio.ErrCarreraCerrada):
		return "Ese gran premio no está disponible"
	default:
		log.Println("apuesta:", err)
		return "No se pudo guardar la apuesta"
	}
}

// DELETE /apuestas/{id}
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
	if err := h.apuestas.Eliminar(r.Context(), id, usuario); err != nil {
		switch {
		case errors.Is(err, dominio.ErrApuestaNoEncontrada):
			http.Error(w, "Apuesta no encontrada", http.StatusNotFound)
		case errors.Is(err, dominio.ErrCarreraCerrada):
			http.Error(w, "La carrera ya se cerró", http.StatusConflict)
		default:
			log.Println("Eliminar apuesta:", err)
			http.Error(w, "No se pudo borrar", http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusOK)
}
