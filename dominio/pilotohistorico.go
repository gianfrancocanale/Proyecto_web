package dominio

import (
	"context"
	"errors"
	"strings"
	"time"
)

// Errores de dominio de PilotoHistorico.
var (
	ErrPilotoHistoricoDuplicado    = errors.New("ya existe un piloto histórico con ese número")
	ErrPilotoHistoricoNoEncontrado = errors.New("piloto histórico no encontrado")
)

// PilotoHistorico guarda los datos biográficos de un piloto, corra o no en la temporada.
type PilotoHistorico struct {
	ID         int       `json:"id_piloto"`
	Nombre     string    `json:"nombre"`
	Pais       string    `json:"pais"`
	Nacimiento time.Time `json:"fecha_nacimiento"`
	Titulos    int       `json:"titulos_ganados"`
}

// Validar aplica las reglas de la entidad.
func (p *PilotoHistorico) Validar() error {
	p.Nombre = strings.TrimSpace(p.Nombre)
	p.Pais = strings.TrimSpace(p.Pais)

	if p.ID <= 0 {
		return ErrValidacion("Número de piloto inválido")
	}
	if p.Nombre == "" || p.Pais == "" || p.Nacimiento.IsZero() {
		return ErrValidacion("Todos los campos son obligatorios")
	}
	if p.Titulos < 0 {
		return ErrValidacion("Los títulos ganados no pueden ser negativos")
	}
	return nil
}

// PilotoHistoricoRepository es el puerto de persistencia de pilotos históricos.
type PilotoHistoricoRepository interface {
	Listar(ctx context.Context) ([]PilotoHistorico, error)
	Recuperar(ctx context.Context, id int) (PilotoHistorico, error)
	Existe(ctx context.Context, id int) (bool, error)
	Crear(ctx context.Context, p PilotoHistorico) error
	Modificar(ctx context.Context, p PilotoHistorico) error
	Eliminar(ctx context.Context, id int) error
}
