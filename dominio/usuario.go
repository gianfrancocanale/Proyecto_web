package dominio

import (
	"context"
	"errors"
	"strings"
)

type Jugador struct {
	Nombre string
	Puntos int
}

// Errores de dominio de Usuario.
var (
	ErrUsuarioDuplicado    = errors.New("ya existe un usuario con ese nombre")
	ErrUsuarioNoEncontrado = errors.New("usuario no encontrado")
	ErrCredenciales        = errors.New("usuario o contraseña incorrectos")
)

// Usuario es la entidad del dominio. Los tags conservan el contrato JSON de la API.
type Usuario struct {
	Nombre     string `json:"nombre_usuario"`
	Contrasena string `json:"contrasena_hash"`
	Puntos     int    `json:"puntos_ganados"`
	EsAdmin    bool   `json:"es_admin"`
}

// Validar aplica las reglas de la entidad (y normaliza el nombre).
func (u *Usuario) Validar() error {
	u.Nombre = strings.TrimSpace(u.Nombre)
	if u.Nombre == "" || strings.TrimSpace(u.Contrasena) == "" {
		return ErrValidacion("El nombre de usuario y la contraseña son obligatorios")
	}
	if u.Puntos < 0 {
		return ErrValidacion("Los puntos ganados no pueden ser negativos")
	}
	return nil
}

// UsuarioRepository es el puerto: lo define el dominio, lo implementa la infraestructura.
type UsuarioRepository interface {
	Listar(ctx context.Context) ([]Usuario, error)
	Recuperar(ctx context.Context, nombre string) (Usuario, error)
	Crear(ctx context.Context, u Usuario) error
	Modificar(ctx context.Context, u Usuario) error
	Eliminar(ctx context.Context, nombre string) error
	EsAdmin(ctx context.Context, nombre string) (bool, error)
	ObtenerPuntos(ctx context.Context, nombre string) (int, error)
	SumarPuntos(ctx context.Context, nombre string, puntos int) error
	ListarJugadores(ctx context.Context) ([]Jugador, error)
}
