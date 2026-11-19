// handlers/granpremio.go
package handlers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"Proyecto_web/aplicacion"
	"Proyecto_web/dominio"
	vista "Proyecto_web/views"
)

type GranPremioHandler struct {
	granPremios *aplicacion.GranPremioService
	historicos  *aplicacion.PilotoHistoricoService
	usuarios    *aplicacion.UsuarioService
}

func NewGranPremioHandler(g *aplicacion.GranPremioService, h *aplicacion.PilotoHistoricoService, u *aplicacion.UsuarioService) *GranPremioHandler {
	return &GranPremioHandler{granPremios: g, historicos: h, usuarios: u}
}

func (h *GranPremioHandler) esAdmin(ctx context.Context, usuario string) bool {
	if usuario == "" {
		return false
	}
	ok, err := h.usuarios.EsAdmin(ctx, usuario)
	return err == nil && ok
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
		vista.GranPremiosPage(h.datosGranPremios(r.Context(), usuario, orden, "")).Render(r.Context(), w)

	case http.MethodPost:
		h.crearDesdeForm(w, r, usuario)

	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

func (h *GranPremioHandler) mostrarError(w http.ResponseWriter, r *http.Request, usuario, mensaje string) {
	h.renderLista(w, r, usuario, mensaje)
}

func (h *GranPremioHandler) crearDesdeForm(w http.ResponseWriter, r *http.Request, usuario string) {
	if err := r.ParseForm(); err != nil {
		h.mostrarError(w, r, usuario, "Formulario inválido")
		return
	}

	g := dominio.GranPremio{
		ID:   strings.TrimSpace(r.FormValue("id_gran_premio")),
		Pais: strings.TrimSpace(r.FormValue("pais")),
	}
	vueltas, err := strconv.Atoi(strings.TrimSpace(r.FormValue("cantidad_vueltas")))
	if err != nil {
		h.mostrarError(w, r, usuario, "Cantidad de vueltas inválida")
		return
	}
	g.Vueltas = vueltas
	g.UltimoGanador, err = strconv.Atoi(r.FormValue("id_ultimo_ganador"))
	if err != nil {
		h.mostrarError(w, r, usuario, "Elegí el último ganador")
		return
	}
	if !asignarLongitud(r.FormValue("longitud_km"), &g.Longitud) {
		h.mostrarError(w, r, usuario, "Longitud inválida")
		return
	}

	if err := h.granPremios.Crear(r.Context(), g); err != nil {
		h.mostrarError(w, r, usuario, mensajeErrorGranPremio(err))
		return
	}
	h.renderLista(w, r, usuario, "")
}

// ---------------------------------------------------------------
// Datos y render
// ---------------------------------------------------------------

func (h *GranPremioHandler) datosGranPremios(ctx context.Context, usuario, orden, mensaje string) vista.GranPremiosData {
	adm := h.esAdmin(ctx, usuario)
	data := vista.GranPremiosData{EsAdmin: adm, Orden: orden, Mensaje: mensaje}

	gps, err := h.granPremios.Listar(ctx, orden)
	if err != nil {
		log.Println("Listar grandes premios:", err)
	} else {
		data.GranPremios = gps
	}
	if adm {
		data.Pilotos = h.opcionesPilotos(ctx)
	}
	return data
}

func (h *GranPremioHandler) renderLista(w http.ResponseWriter, r *http.Request, usuario, mensaje string) {
	data := h.datosGranPremios(r.Context(), usuario, "asc", mensaje)
	vista.GranPremiosLista(data.GranPremios, data.EsAdmin, data.Mensaje).Render(r.Context(), w)
}

func (h *GranPremioHandler) opcionesPilotos(ctx context.Context) []vista.OpcionPiloto {
	pilotos, err := h.historicos.Opciones(ctx)
	if err != nil {
		log.Println("Listar pilotos históricos:", err)
		return nil
	}
	out := make([]vista.OpcionPiloto, 0, len(pilotos))
	for _, p := range pilotos {
		out = append(out, vista.OpcionPiloto{ID: p.ID, Nombre: p.Nombre})
	}
	return out
}

func mensajeErrorGranPremio(err error) string {
	var ev dominio.ErrValidacion
	switch {
	case errors.As(err, &ev):
		return ev.Error()
	case errors.Is(err, dominio.ErrGranPremioDuplicado):
		return "Ya existe un gran premio con ese nombre"
	default:
		log.Println("gran premio:", err)
		return "No se pudo crear el gran premio"
	}
}

// asignarLongitud valida el texto y lo escribe en el destino (int, float o string).
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

		// La identidad viaja en un campo oculto; no se puede cambiar.
		id := strings.TrimSpace(r.FormValue("id_gran_premio"))
		if id == "" {
			http.Redirect(w, r, "/gran_premios", http.StatusSeeOther)
			return
		}
		g := dominio.GranPremio{ID: id, Pais: strings.TrimSpace(r.FormValue("pais"))}
		vueltas, err := strconv.Atoi(strings.TrimSpace(r.FormValue("cantidad_vueltas")))
		if err != nil {
			h.renderModificar(w, r, id, "Cantidad de vueltas inválida", http.StatusBadRequest)
			return
		}
		g.Vueltas = vueltas
		g.UltimoGanador, err = strconv.Atoi(r.FormValue("id_ultimo_ganador"))
		if err != nil {
			h.renderModificar(w, r, id, "Elegí el último ganador", http.StatusBadRequest)
			return
		}
		if !asignarLongitud(r.FormValue("longitud_km"), &g.Longitud) {
			h.renderModificar(w, r, id, "Longitud inválida", http.StatusBadRequest)
			return
		}

		if err := h.granPremios.Modificar(ctx, g); err != nil {
			var ev dominio.ErrValidacion
			status := http.StatusInternalServerError
			mensaje := "No se pudo guardar el cambio"
			if errors.As(err, &ev) {
				status, mensaje = http.StatusBadRequest, ev.Error()
			}
			h.renderModificar(w, r, id, mensaje, status)
			return
		}
		http.Redirect(w, r, "/gran_premios", http.StatusSeeOther)

	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

func (h *GranPremioHandler) renderModificar(w http.ResponseWriter, r *http.Request, id, mensaje string, status int) {
	ctx := r.Context()

	g, err := h.granPremios.Recuperar(ctx, id)
	if err != nil {
		http.Redirect(w, r, "/gran_premios", http.StatusSeeOther)
		return
	}

	data := vista.ModificarGranPremioData{
		Mensaje:  mensaje,
		ID:       g.ID,
		Pais:     g.Pais,
		Longitud: fmt.Sprint(g.Longitud),
		Vueltas:  fmt.Sprint(g.Vueltas),
		Ganador:  fmt.Sprint(g.UltimoGanador),
		Pilotos:  h.opcionesPilotos(ctx),
	}
	w.WriteHeader(status)
	vista.ModificarGranPremioPage(data).Render(ctx, w)
}

// ---------------------------------------------------------------
// Eliminar un gran premio (admin)
// ---------------------------------------------------------------

func (h *GranPremioHandler) EliminarGranPremio(w http.ResponseWriter, r *http.Request) {
	usuario, _ := usuarioActual(r)

	id := r.PathValue("id")
	if id == "" {
		h.renderLista(w, r, usuario, "ID de gran premio inválido")
		return
	}
	if err := h.granPremios.Eliminar(r.Context(), id); err != nil {
		log.Println("Eliminar gran premio:", err)
		msg := "No se pudo borrar el gran premio"
		if strings.Contains(strings.ToLower(err.Error()), "foreign key") {
			msg = "No se puede borrar: tiene carreras o apuestas asociadas"
		}
		h.renderLista(w, r, usuario, msg)
		return
	}
	h.renderLista(w, r, usuario, "")
}
