package infrastructura

import (
	"context"
	"database/sql"
	"strings"
	"time"

	sqlc "Proyecto_web/db/sqlc"
	"Proyecto_web/dominio"
)

type ApuestaRepository struct {
	q *sqlc.Queries
}

func NewApuestaRepository(db *sql.DB) *ApuestaRepository {
	return &ApuestaRepository{q: sqlc.New(db)}
}

func aDominioApuesta(a sqlc.Apuestum) dominio.Apuesta {
	out := dominio.Apuesta{
		ID:           int(a.IDApuesta),
		Usuario:      a.NombreUsuario,
		GranPremio:   a.IDGranPremio,
		FechaCarrera: a.FechaCarrera,
		FechaApuesta: a.FechaApuesta,
	}
	for _, id := range a.Prediccion {
		out.Prediccion = append(out.Prediccion, int(id))
	}
	return out
}

func (r *ApuestaRepository) ListarPorUsuario(ctx context.Context, usuario string) ([]dominio.Apuesta, error) {
	filas, err := r.q.ListarApuestasPorUsuario(ctx, usuario)
	if err != nil {
		return nil, err
	}
	out := make([]dominio.Apuesta, 0, len(filas))
	for _, f := range filas {
		a := dominio.Apuesta{
			ID:           int(f.IDApuesta),
			GranPremio:   f.IDGranPremio,
			FechaCarrera: f.FechaCarrera,
			FechaApuesta: f.FechaApuesta,
		}
		for _, id := range f.Prediccion {
			a.Prediccion = append(a.Prediccion, int(id))
		}
		out = append(out, a)
	}
	return out, nil
}

func (r *ApuestaRepository) ListarPorGranPremio(ctx context.Context, id string, fecha time.Time) ([]dominio.Apuesta, error) {
	filas, err := r.q.ListarApuestasPorGranPremio(ctx, sqlc.ListarApuestasPorGranPremioParams{
		IDGranPremio: id,
		FechaCarrera: fecha,
	})
	if err != nil {
		return nil, err
	}
	out := make([]dominio.Apuesta, 0, len(filas))
	for _, f := range filas {
		a := dominio.Apuesta{ID: int(f.IDApuesta), Usuario: f.NombreUsuario}
		for _, id := range f.Prediccion {
			a.Prediccion = append(a.Prediccion, int(id))
		}
		out = append(out, a)
	}
	return out, nil
}

func (r *ApuestaRepository) Recuperar(ctx context.Context, id int, usuario string) (dominio.Apuesta, error) {
	a, err := r.q.RecuperarApuestaPorIdUsuario(ctx, sqlc.RecuperarApuestaPorIdUsuarioParams{
		NombreUsuario: usuario,
		IDApuesta:     int32(id),
	})
	if esNoRows(err) {
		return dominio.Apuesta{}, dominio.ErrApuestaNoEncontrada
	}
	if err != nil {
		return dominio.Apuesta{}, err
	}
	out := dominio.Apuesta{
		ID:           int(a.IDApuesta),
		Usuario:      usuario,
		GranPremio:   a.IDGranPremio,
		FechaCarrera: a.FechaCarrera,
		FechaApuesta: a.FechaApuesta,
	}
	for _, id := range a.Prediccion {
		out.Prediccion = append(out.Prediccion, int(id))
	}
	return out, nil
}

func (r *ApuestaRepository) Crear(ctx context.Context, a dominio.Apuesta) error {
	_, err := r.q.CrearApuesta(ctx, sqlc.CrearApuestaParams{
		NombreUsuario: a.Usuario,
		IDGranPremio:  a.GranPremio,
		FechaCarrera:  a.FechaCarrera,
		Prediccion:    aInt32(a.Prediccion),
		FechaApuesta:  a.FechaApuesta,
	})
	if err != nil && strings.Contains(err.Error(), "unq_usuario_gran_premio_carrera") {
		return dominio.ErrApuestaDuplicada
	}
	if esDuplicado(err) {
		return dominio.ErrApuestaDuplicada
	}
	return err
}

func (r *ApuestaRepository) Actualizar(ctx context.Context, a dominio.Apuesta) error {
	_, err := r.q.ActualizarPrediccionApuesta(ctx, sqlc.ActualizarPrediccionApuestaParams{
		IDApuesta:     int32(a.ID),
		NombreUsuario: a.Usuario,
		Prediccion:    aInt32(a.Prediccion),
		FechaApuesta:  a.FechaApuesta,
	})
	if esNoRows(err) {
		return dominio.ErrApuestaNoEncontrada
	}
	return err
}

func (r *ApuestaRepository) Eliminar(ctx context.Context, id int, usuario string) error {
	return r.q.EliminarApuesta(ctx, sqlc.EliminarApuestaParams{
		IDApuesta:     int32(id),
		NombreUsuario: usuario,
	})
}
