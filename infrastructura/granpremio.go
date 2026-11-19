package infrastructura

import (
	"context"
	"database/sql"

	sqlc "Proyecto_web/db/sqlc"
	"Proyecto_web/dominio"
)

type GranPremioRepository struct {
	q *sqlc.Queries
}

func NewGranPremioRepository(db *sql.DB) *GranPremioRepository {
	return &GranPremioRepository{q: sqlc.New(db)}
}

func aDominioGranPremio(g sqlc.ListarGrandesPremiosRow) dominio.GranPremio {
	nombre := ""
	if g.NombreUltimoGanador.Valid {
		nombre = g.NombreUltimoGanador.String
	}
	return dominio.GranPremio{
		ID:                  g.IDGranPremio,
		Pais:                g.Pais,
		Longitud:            int(g.LongitudKm),
		Vueltas:             int(g.CantidadVueltas),
		UltimoGanador:       int(g.IDUltimoGanador),
		NombreUltimoGanador: nombre,
	}
}

func (r *GranPremioRepository) Listar(ctx context.Context) ([]dominio.GranPremio, error) {
	filas, err := r.q.ListarGrandesPremios(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]dominio.GranPremio, 0, len(filas))
	for _, f := range filas {
		out = append(out, aDominioGranPremio(f))
	}
	return out, nil
}

func (r *GranPremioRepository) Recuperar(ctx context.Context, id string) (dominio.GranPremio, error) {
	g, err := r.q.RecuperarGranPremio(ctx, id)
	if esNoRows(err) {
		return dominio.GranPremio{}, dominio.ErrGranPremioNoEncontrado
	}
	if err != nil {
		return dominio.GranPremio{}, err
	}
	return dominio.GranPremio{
		ID:            g.IDGranPremio,
		Pais:          g.Pais,
		Longitud:      int(g.LongitudKm),
		Vueltas:       int(g.CantidadVueltas),
		UltimoGanador: int(g.IDUltimoGanador),
	}, nil
}

func (r *GranPremioRepository) Existe(ctx context.Context, id string) (bool, error) {
	_, err := r.q.RecuperarGranPremio(ctx, id)
	if esNoRows(err) {
		return false, nil
	}
	return err == nil, err
}

func (r *GranPremioRepository) Circuitos(ctx context.Context) ([]string, error) {
	return r.q.ListarCircuitosParaAdmin(ctx)
}

func (r *GranPremioRepository) Crear(ctx context.Context, g dominio.GranPremio) error {
	_, err := r.q.CrearGranPremio(ctx, sqlc.CrearGranPremioParams{
		IDGranPremio:    g.ID,
		Pais:            g.Pais,
		LongitudKm:      int32(g.Longitud),
		CantidadVueltas: int32(g.Vueltas),
		IDUltimoGanador: int32(g.UltimoGanador),
	})
	if esDuplicado(err) {
		return dominio.ErrGranPremioDuplicado
	}
	return err
}

func (r *GranPremioRepository) Modificar(ctx context.Context, g dominio.GranPremio) error {
	_, err := r.q.ModificarGranPremio(ctx, sqlc.ModificarGranPremioParams{
		IDGranPremio:    g.ID,
		Pais:            g.Pais,
		LongitudKm:      int32(g.Longitud),
		CantidadVueltas: int32(g.Vueltas),
		IDUltimoGanador: int32(g.UltimoGanador),
	})
	if esNoRows(err) {
		return dominio.ErrGranPremioNoEncontrado
	}
	return err
}

func (r *GranPremioRepository) Eliminar(ctx context.Context, id string) error {
	return r.q.EliminarGranPremio(ctx, id)
}
