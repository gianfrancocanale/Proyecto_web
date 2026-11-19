// repository/escuderia.go
package infrastructura

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	sqlc "Proyecto_web/db/sqlc"
	"Proyecto_web/dominio"
)

type EscuderiaRepository struct {
	q *sqlc.Queries
}

func NewEscuderiaRepository(db *sql.DB) *EscuderiaRepository {
	return &EscuderiaRepository{q: sqlc.New(db)}
}

// sqlc -> dominio
func aDominio(e sqlc.Escuderium) dominio.Escuderia {
	return dominio.Escuderia{
		ID:        e.IDEscuderia,
		Principal: e.TeamPrincipal,
		Titulos:   int(e.TitulosConstructores),
		Puntos:    int(e.PuntosTemporada),
		Fundacion: e.FechaFundacion,
	}
}

func (r *EscuderiaRepository) Listar(ctx context.Context) ([]dominio.Escuderia, error) {
	filas, err := r.q.ListarEscuderias(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]dominio.Escuderia, len(filas))
	for i, f := range filas {
		out[i] = aDominio(f)
	}
	return out, nil
}

func (r *EscuderiaRepository) Recuperar(ctx context.Context, id string) (dominio.Escuderia, error) {
	f, err := r.q.RecuperarEscuderia(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return dominio.Escuderia{}, dominio.ErrNoEncontrada
	}
	if err != nil {
		return dominio.Escuderia{}, err
	}
	return aDominio(f), nil
}

func (r *EscuderiaRepository) Crear(ctx context.Context, e dominio.Escuderia) error {
	_, err := r.q.CrearEscuderia(ctx, sqlc.CrearEscuderiaParams{
		IDEscuderia:          e.ID,
		PuntosTemporada:      int32(e.Puntos),
		TitulosConstructores: int32(e.Titulos),
		FechaFundacion:       e.Fundacion,
		TeamPrincipal:        e.Principal,
	})
	if err != nil && strings.Contains(err.Error(), "duplicate key") {
		return dominio.ErrDuplicada
	}
	return err
}

func (r *EscuderiaRepository) Modificar(ctx context.Context, e dominio.Escuderia) error {
	_, err := r.q.ModificarEscuderia(ctx, sqlc.ModificarEscuderiaParams{
		IDEscuderia:          e.ID,
		TeamPrincipal:        e.Principal,
		TitulosConstructores: int32(e.Titulos),
		FechaFundacion:       e.Fundacion,
		PuntosTemporada:      int32(e.Puntos),
	})
	if errors.Is(err, sql.ErrNoRows) {
		return dominio.ErrNoEncontrada
	}
	return err
}

func (r *EscuderiaRepository) Eliminar(ctx context.Context, id string) error {
	return r.q.EliminarEscuderia(ctx, id)
}
