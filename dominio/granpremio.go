package dominio

import (
	"context"
	"errors"
	"strings"
)

// Errores de dominio de GranPremio.
var (
	ErrGranPremioDuplicado    = errors.New("ya existe un gran premio con ese nombre")
	ErrGranPremioNoEncontrado = errors.New("gran premio no encontrado")
)

// GranPremio es un circuito del catálogo histórico.
type GranPremio struct {
	ID                  string `json:"id_gran_premio"`
	Pais                string `json:"pais"`
	Longitud            int    `json:"longitud_km"`
	Vueltas             int    `json:"cantidad_vueltas"`
	UltimoGanador       int    `json:"id_ultimo_ganador"`
	NombreUltimoGanador string `json:"-"`
}

// Validar aplica las reglas de la entidad (y normaliza los textos).
func (g *GranPremio) Validar() error {
	g.ID = strings.TrimSpace(g.ID)
	g.Pais = strings.TrimSpace(g.Pais)

	if g.ID == "" || g.Pais == "" {
		return ErrValidacion("Completá el nombre del gran premio y el país")
	}
	if g.Longitud <= 0 {
		return ErrValidacion("La longitud tiene que ser un número mayor a 0")
	}
	if g.Vueltas <= 0 {
		return ErrValidacion("La cantidad de vueltas tiene que ser un número mayor a 0")
	}
	if g.UltimoGanador <= 0 {
		return ErrValidacion("Elegí el último ganador")
	}
	return nil
}

// GranPremioRepository es el puerto de persistencia de circuitos.
type GranPremioRepository interface {
	Listar(ctx context.Context) ([]GranPremio, error)
	Recuperar(ctx context.Context, id string) (GranPremio, error)
	Existe(ctx context.Context, id string) (bool, error)
	Circuitos(ctx context.Context) ([]string, error)
	Crear(ctx context.Context, g GranPremio) error
	Modificar(ctx context.Context, g GranPremio) error
	Eliminar(ctx context.Context, id string) error
}
