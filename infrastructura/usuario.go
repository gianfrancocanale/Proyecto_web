package infrastructura

import (
	"context"
	"database/sql"

	sqlc "Proyecto_web/db/sqlc"
	"Proyecto_web/dominio"
)

type UsuarioRepository struct {
	q *sqlc.Queries
}

func NewUsuarioRepository(db *sql.DB) *UsuarioRepository {
	return &UsuarioRepository{q: sqlc.New(db)}
}

func aDominioUsuario(u sqlc.Usuario) dominio.Usuario {
	return dominio.Usuario{
		Nombre:     u.NombreUsuario,
		Contrasena: u.ContrasenaHash,
		Puntos:     int(u.PuntosGanados),
		EsAdmin:    u.EsAdmin,
	}
}

func (r *UsuarioRepository) Listar(ctx context.Context) ([]dominio.Usuario, error) {
	filas, err := r.q.RecuperarUsuarios(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]dominio.Usuario, 0, len(filas))
	for _, f := range filas {
		out = append(out, dominio.Usuario{Nombre: f.NombreUsuario, Puntos: int(f.PuntosGanados)})
	}
	return out, nil
}

func (r *UsuarioRepository) Recuperar(ctx context.Context, nombre string) (dominio.Usuario, error) {
	u, err := r.q.RecuperarUsuarioPorId(ctx, nombre)
	if esNoRows(err) {
		return dominio.Usuario{}, dominio.ErrUsuarioNoEncontrado
	}
	if err != nil {
		return dominio.Usuario{}, err
	}
	return aDominioUsuario(u), nil
}

func (r *UsuarioRepository) Crear(ctx context.Context, u dominio.Usuario) error {
	_, err := r.q.CrearUsuario(ctx, sqlc.CrearUsuarioParams{
		NombreUsuario:  u.Nombre,
		ContrasenaHash: u.Contrasena,
		PuntosGanados:  int32(u.Puntos),
	})
	if esDuplicado(err) {
		return dominio.ErrUsuarioDuplicado
	}
	return err
}

func (r *UsuarioRepository) Modificar(ctx context.Context, u dominio.Usuario) error {
	_, err := r.q.ModificarUsuario(ctx, sqlc.ModificarUsuarioParams{
		NombreUsuario:  u.Nombre,
		ContrasenaHash: u.Contrasena,
		PuntosGanados:  int32(u.Puntos),
	})
	if esNoRows(err) {
		return dominio.ErrUsuarioNoEncontrado
	}
	return err
}

func (r *UsuarioRepository) Eliminar(ctx context.Context, nombre string) error {
	return r.q.EliminarUsuario(ctx, nombre)
}

func (r *UsuarioRepository) EsAdmin(ctx context.Context, nombre string) (bool, error) {
	ok, err := r.q.EsAdmin(ctx, nombre)
	if esNoRows(err) {
		return false, nil
	}
	return ok, err
}

func (r *UsuarioRepository) ObtenerPuntos(ctx context.Context, nombre string) (int, error) {
	puntos, err := r.q.ObtenerPuntosUsuario(ctx, nombre)
	return int(puntos), err
}

func (r *UsuarioRepository) SumarPuntos(ctx context.Context, nombre string, puntos int) error {
	return r.q.SumarPuntosUsuario(ctx, sqlc.SumarPuntosUsuarioParams{
		NombreUsuario: nombre,
		Puntos:        int32(puntos),
	})
}

// repository/usuario.go
func (r *UsuarioRepository) ListarJugadores(ctx context.Context) ([]dominio.Jugador, error) {
	filas, err := r.q.ListarJugadores(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]dominio.Jugador, len(filas))
	for i, f := range filas {
		out[i] = dominio.Jugador{
			Nombre: f.NombreUsuario,
			Puntos: int(f.PuntosGanados),
		}
	}
	return out, nil
}
