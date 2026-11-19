package infrastructura

import (
	"context"
	"database/sql"

	sqlc "Proyecto_web/db/sqlc"
	"Proyecto_web/dominio"
)

type PilotoHistoricoRepository struct {
	q *sqlc.Queries
}

func NewPilotoHistoricoRepository(db *sql.DB) *PilotoHistoricoRepository {
	return &PilotoHistoricoRepository{q: sqlc.New(db)}
}

func aDominioPilotoHistorico(p sqlc.PilotoHistorico) dominio.PilotoHistorico {
	return dominio.PilotoHistorico{
		ID:         int(p.IDPiloto),
		Nombre:     p.Nombre,
		Pais:       p.Pais,
		Nacimiento: p.FechaNacimiento,
		Titulos:    int(p.TitulosGanados),
	}
}

func (r *PilotoHistoricoRepository) Listar(ctx context.Context) ([]dominio.PilotoHistorico, error) {
	filas, err := r.q.ListarPilotosHistoricos(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]dominio.PilotoHistorico, 0, len(filas))
	for _, f := range filas {
		out = append(out, dominio.PilotoHistorico{ID: int(f.IDPiloto), Nombre: f.Nombre})
	}
	return out, nil
}

func (r *PilotoHistoricoRepository) Recuperar(ctx context.Context, id int) (dominio.PilotoHistorico, error) {
	p, err := r.q.RecuperarPilotoHistorico(ctx, int32(id))
	if esNoRows(err) {
		return dominio.PilotoHistorico{}, dominio.ErrPilotoHistoricoNoEncontrado
	}
	if err != nil {
		return dominio.PilotoHistorico{}, err
	}
	return aDominioPilotoHistorico(p), nil
}

func (r *PilotoHistoricoRepository) Existe(ctx context.Context, id int) (bool, error) {
	_, err := r.q.RecuperarPilotoHistorico(ctx, int32(id))
	if esNoRows(err) {
		return false, nil
	}
	return err == nil, err
}

func (r *PilotoHistoricoRepository) Crear(ctx context.Context, p dominio.PilotoHistorico) error {
	_, err := r.q.CrearPilotoHistorico(ctx, sqlc.CrearPilotoHistoricoParams{
		IDPiloto:        int32(p.ID),
		Nombre:          p.Nombre,
		Pais:            p.Pais,
		FechaNacimiento: p.Nacimiento,
		TitulosGanados:  int32(p.Titulos),
	})
	if esDuplicado(err) {
		return dominio.ErrPilotoHistoricoDuplicado
	}
	return err
}

func (r *PilotoHistoricoRepository) Modificar(ctx context.Context, p dominio.PilotoHistorico) error {
	_, err := r.q.ModificarPilotoHistorico(ctx, sqlc.ModificarPilotoHistoricoParams{
		IDPiloto:        int32(p.ID),
		Nombre:          p.Nombre,
		Pais:            p.Pais,
		FechaNacimiento: p.Nacimiento,
		TitulosGanados:  int32(p.Titulos),
	})
	if esNoRows(err) {
		return dominio.ErrPilotoHistoricoNoEncontrado
	}
	return err
}

func (r *PilotoHistoricoRepository) Eliminar(ctx context.Context, id int) error {
	return r.q.EliminarPilotoHistorico(ctx, int32(id))
}
