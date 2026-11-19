package infrastructura

import (
	"context"
	"database/sql"

	sqlc "Proyecto_web/db/sqlc"
	"Proyecto_web/dominio"
)

type PilotoRepository struct {
	db *sql.DB
	q  *sqlc.Queries
}

func NewPilotoRepository(db *sql.DB) *PilotoRepository {
	return &PilotoRepository{db: db, q: sqlc.New(db)}
}

func aDominioPiloto(p sqlc.ListarPilotosRow) dominio.Piloto {
	return dominio.Piloto{
		ID:         int(p.IDPiloto),
		Nombre:     p.Nombre,
		Pais:       p.Pais,
		Nacimiento: p.FechaNacimiento,
		Titulos:    int(p.TitulosGanados),
		Escuderia:  p.IDEscuderia,
		Puntos:     int(p.PuntosTemporada),
	}
}

func (r *PilotoRepository) Listar(ctx context.Context) ([]dominio.Piloto, error) {
	filas, err := r.q.ListarPilotos(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]dominio.Piloto, 0, len(filas))
	for _, f := range filas {
		out = append(out, aDominioPiloto(f))
	}
	return out, nil
}

func (r *PilotoRepository) Recuperar(ctx context.Context, id int) (dominio.Piloto, error) {
	p, err := r.q.RecuperarPiloto(ctx, int32(id))
	if esNoRows(err) {
		return dominio.Piloto{}, dominio.ErrPilotoNoEncontrado
	}
	if err != nil {
		return dominio.Piloto{}, err
	}
	h, err := r.q.RecuperarPilotoHistorico(ctx, int32(id))
	if err != nil {
		return dominio.Piloto{}, err
	}
	return dominio.Piloto{
		ID:         int(p.IDPiloto),
		Nombre:     h.Nombre,
		Pais:       h.Pais,
		Nacimiento: h.FechaNacimiento,
		Titulos:    int(h.TitulosGanados),
		Escuderia:  p.IDEscuderia,
		Puntos:     int(p.PuntosTemporada),
	}, nil
}

func (r *PilotoRepository) ListarParaApuesta(ctx context.Context) ([]dominio.Piloto, error) {
	filas, err := r.q.ListarPilotosParaApuesta(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]dominio.Piloto, 0, len(filas))
	for _, f := range filas {
		out = append(out, dominio.Piloto{ID: int(f.IDPiloto), Nombre: f.Nombre})
	}
	return out, nil
}

func (r *PilotoRepository) Escuderias(ctx context.Context) ([]string, error) {
	filas, err := r.q.ListarEscuderias(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]string, 0, len(filas))
	for _, e := range filas {
		out = append(out, e.IDEscuderia)
	}
	return out, nil
}

func (r *PilotoRepository) ExisteEscuderia(ctx context.Context, id string) (bool, error) {
	_, err := r.q.RecuperarEscuderia(ctx, id)
	if esNoRows(err) {
		return false, nil
	}
	return err == nil, err
}

func (r *PilotoRepository) ExisteEnTemporada(ctx context.Context, id int) (bool, error) {
	_, err := r.q.RecuperarPiloto(ctx, int32(id))
	if esNoRows(err) {
		return false, nil
	}
	return err == nil, err
}

func (r *PilotoRepository) ExisteHistorico(ctx context.Context, id int) (bool, error) {
	_, err := r.q.RecuperarPilotoHistorico(ctx, int32(id))
	if esNoRows(err) {
		return false, nil
	}
	return err == nil, err
}

// Inscribir mantiene la atomicidad de "histórico + temporada" dentro de una transacción.
func (r *PilotoRepository) Inscribir(ctx context.Context, p dominio.Piloto, crearHistorico bool) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	qtx := r.q.WithTx(tx)

	if crearHistorico {
		_, err = qtx.CrearPilotoHistorico(ctx, sqlc.CrearPilotoHistoricoParams{
			IDPiloto:        int32(p.ID),
			Nombre:          p.Nombre,
			Pais:            p.Pais,
			FechaNacimiento: p.Nacimiento,
			TitulosGanados:  int32(p.Titulos),
		})
		if err != nil {
			return err
		}
	}

	_, err = qtx.CrearPilotoTemporada(ctx, sqlc.CrearPilotoTemporadaParams{
		IDPiloto:        int32(p.ID),
		IDEscuderia:     p.Escuderia,
		PuntosTemporada: int32(p.Puntos),
	})
	if err != nil {
		if esDuplicado(err) {
			return dominio.ErrPilotoDuplicado
		}
		return err
	}
	return tx.Commit()
}

func (r *PilotoRepository) Modificar(ctx context.Context, p dominio.Piloto) error {
	_, err := r.q.ModificarPiloto(ctx, sqlc.ModificarPilotoParams{
		IDPiloto:        int32(p.ID),
		IDEscuderia:     p.Escuderia,
		PuntosTemporada: int32(p.Puntos),
	})
	if esNoRows(err) {
		return dominio.ErrPilotoNoEncontrado
	}
	return err
}

func (r *PilotoRepository) Eliminar(ctx context.Context, id int) error {
	return r.q.EliminarPilotoTemporada(ctx, int32(id))
}
