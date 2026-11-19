package infrastructura

import (
	"context"
	"database/sql"
	"time"

	sqlc "Proyecto_web/db/sqlc"
	"Proyecto_web/dominio"
)

type GranPremioHistoricoRepository struct {
	db *sql.DB
	q  *sqlc.Queries
}

func NewGranPremioHistoricoRepository(db *sql.DB) *GranPremioHistoricoRepository {
	return &GranPremioHistoricoRepository{db: db, q: sqlc.New(db)}
}

func aDominioCarrera(g sqlc.GranPremioHistorico) dominio.GranPremioHistorico {
	out := dominio.GranPremioHistorico{ID: g.IDGranPremio, Fecha: g.FechaCarrera}
	for _, id := range g.ResultadoCarrera {
		out.Resultado = append(out.Resultado, int(id))
	}
	return out
}

func (r *GranPremioHistoricoRepository) Listar(ctx context.Context) ([]dominio.GranPremioHistorico, error) {
	filas, err := r.q.ListarGrandesPremiosHistoricos(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]dominio.GranPremioHistorico, 0, len(filas))
	for _, f := range filas {
		g := dominio.GranPremioHistorico{ID: f.IDGranPremio, Fecha: f.FechaCarrera}
		for _, id := range f.ResultadoCarrera {
			g.Resultado = append(g.Resultado, int(id))
		}
		out = append(out, g)
	}
	return out, nil
}

func (r *GranPremioHistoricoRepository) ListarProximos(ctx context.Context) ([]dominio.GranPremioHistorico, error) {
	filas, err := r.q.ListarGranPremiosProximos(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]dominio.GranPremioHistorico, 0, len(filas))
	for _, f := range filas {
		out = append(out, dominio.GranPremioHistorico{ID: f.IDGranPremio, Fecha: f.FechaCarrera})
	}
	return out, nil
}

func (r *GranPremioHistoricoRepository) ListarPendientes(ctx context.Context) ([]dominio.GranPremioHistorico, error) {
	filas, err := r.q.ListarEdicionesPendientes(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]dominio.GranPremioHistorico, 0, len(filas))
	for _, f := range filas {
		out = append(out, dominio.GranPremioHistorico{ID: f.IDGranPremio, Fecha: f.FechaCarrera})
	}
	return out, nil
}

func (r *GranPremioHistoricoRepository) Recuperar(ctx context.Context, id string, fecha time.Time) (dominio.GranPremioHistorico, error) {
	g, err := r.q.RecuperarGranPremioHistorico(ctx, sqlc.RecuperarGranPremioHistoricoParams{
		IDGranPremio: id,
		FechaCarrera: fecha,
	})
	if esNoRows(err) {
		return dominio.GranPremioHistorico{}, dominio.ErrCarreraNoEncontrada
	}
	if err != nil {
		return dominio.GranPremioHistorico{}, err
	}
	return aDominioCarrera(g), nil
}

func (r *GranPremioHistoricoRepository) Programar(ctx context.Context, g dominio.GranPremioHistorico) error {
	_, err := r.q.CrearEdicionGranPremio(ctx, sqlc.CrearEdicionGranPremioParams{
		IDGranPremio: g.ID,
		FechaCarrera: g.Fecha,
	})
	if esDuplicado(err) {
		return dominio.ErrEdicionDuplicada
	}
	return err
}

func (r *GranPremioHistoricoRepository) Crear(ctx context.Context, g dominio.GranPremioHistorico) error {
	_, err := r.q.CrearGranPremioHistorico(ctx, sqlc.CrearGranPremioHistoricoParams{
		IDGranPremio:     g.ID,
		FechaCarrera:     g.Fecha,
		ResultadoCarrera: aInt32(g.Resultado),
	})
	if esDuplicado(err) {
		return dominio.ErrEdicionDuplicada
	}
	return err
}

func (r *GranPremioHistoricoRepository) RegistrarResultado(ctx context.Context, g dominio.GranPremioHistorico) error {
	_, err := r.q.RegistrarResultadoCarrera(ctx, sqlc.RegistrarResultadoCarreraParams{
		IDGranPremio:     g.ID,
		FechaCarrera:     g.Fecha,
		ResultadoCarrera: aInt32(g.Resultado),
	})
	if esNoRows(err) {
		return dominio.ErrCarreraNoEncontrada
	}
	return err
}

// CargarResultado guarda el resultado y el reparto de puntos en una sola transacción.
func (r *GranPremioHistoricoRepository) CargarResultado(ctx context.Context, g dominio.GranPremioHistorico, puntos []dominio.PuntosUsuario) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	qtx := r.q.WithTx(tx)

	_, err = qtx.CargarResultadoCarrera(ctx, sqlc.CargarResultadoCarreraParams{
		IDGranPremio:     g.ID,
		FechaCarrera:     g.Fecha,
		ResultadoCarrera: aInt32(g.Resultado),
	})
	if esNoRows(err) {
		return dominio.ErrCarreraNoEncontrada
	}
	if err != nil {
		return err
	}

	for _, p := range puntos {
		if err := qtx.SumarPuntosUsuario(ctx, sqlc.SumarPuntosUsuarioParams{
			NombreUsuario: p.Usuario,
			Puntos:        int32(p.Puntos),
		}); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (r *GranPremioHistoricoRepository) Eliminar(ctx context.Context, id string, fecha time.Time) error {
	return r.q.EliminarGranPremioHistorico(ctx, sqlc.EliminarGranPremioHistoricoParams{
		IDGranPremio: id,
		FechaCarrera: fecha,
	})
}

func aInt32(valores []int) []int32 {
	out := make([]int32, len(valores))
	for i, v := range valores {
		out[i] = int32(v)
	}
	return out
}
