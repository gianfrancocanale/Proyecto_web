package dominio

import (
	"context"
	"errors"
	"time"
)

// Errores de dominio de GranPremioHistorico (ediciones/carreras).
var (
	ErrEdicionDuplicada    = errors.New("esa carrera ya está programada")
	ErrCarreraNoEncontrada = errors.New("carrera no encontrada")
	ErrCarreraCerrada      = errors.New("la carrera ya se cerró")
	ErrResultadoIncompleto = errors.New("el resultado de la carrera es obligatorio")
)

// GranPremioHistorico es una edición concreta de un gran premio.
// Resultado vacío significa que todavía no se corrió (o no se cargó).
type GranPremioHistorico struct {
	ID        string    `json:"id_gran_premio"`
	Fecha     time.Time `json:"fecha_carrera"`
	Resultado []int     `json:"resultado_carrera"`
}

// Abierta indica si todavía se puede apostar en esta edición.
func (g GranPremioHistorico) Abierta() bool {
	if len(g.Resultado) > 0 {
		return false
	}
	hoy := time.Now().UTC().Truncate(24 * time.Hour)
	fecha := g.Fecha.UTC().Truncate(24 * time.Hour)
	return fecha.After(hoy)
}

// Validar aplica las reglas de la entidad.
func (g *GranPremioHistorico) Validar() error {
	if g.ID == "" || g.Fecha.IsZero() {
		return ErrValidacion("El gran premio y la fecha de carrera son obligatorios")
	}
	return nil
}

// PuntosUsuario es el reparto de puntos calculado para un usuario.
type PuntosUsuario struct {
	Usuario string
	Puntos  int
}

// GranPremioHistoricoRepository es el puerto de persistencia de ediciones.
type GranPremioHistoricoRepository interface {
	Listar(ctx context.Context) ([]GranPremioHistorico, error)
	ListarProximos(ctx context.Context) ([]GranPremioHistorico, error)
	ListarPendientes(ctx context.Context) ([]GranPremioHistorico, error)
	Recuperar(ctx context.Context, id string, fecha time.Time) (GranPremioHistorico, error)
	Programar(ctx context.Context, g GranPremioHistorico) error
	Crear(ctx context.Context, g GranPremioHistorico) error
	RegistrarResultado(ctx context.Context, g GranPremioHistorico) error
	// CargarResultado guarda el orden de llegada y reparte los puntos en una transacción.
	CargarResultado(ctx context.Context, g GranPremioHistorico, puntos []PuntosUsuario) error
	Eliminar(ctx context.Context, id string, fecha time.Time) error
}
