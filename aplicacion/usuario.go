package aplicacion

import (
	"context"
	"errors"
	"sort"

	"Proyecto_web/dominio"

	"golang.org/x/crypto/bcrypt"
)

type UsuarioService struct {
	repo dominio.UsuarioRepository
}

func NewUsuarioService(repo dominio.UsuarioRepository) *UsuarioService {
	return &UsuarioService{repo: repo}
}

func (s *UsuarioService) Listar(ctx context.Context) ([]dominio.Usuario, error) {
	return s.repo.Listar(ctx)
}

func (s *UsuarioService) Recuperar(ctx context.Context, nombre string) (dominio.Usuario, error) {
	return s.repo.Recuperar(ctx, nombre)
}

// Crear hashea la contraseña antes de persistir: el hash es regla de negocio.
func (s *UsuarioService) Crear(ctx context.Context, u dominio.Usuario) error {
	if err := u.Validar(); err != nil {
		return err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Contrasena), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Contrasena = string(hash)
	u.Puntos = 0
	u.EsAdmin = false
	return s.repo.Crear(ctx, u)
}

func (s *UsuarioService) Modificar(ctx context.Context, u dominio.Usuario) error {
	if err := u.Validar(); err != nil {
		return err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Contrasena), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Contrasena = string(hash)
	return s.repo.Modificar(ctx, u)
}

func (s *UsuarioService) Eliminar(ctx context.Context, nombre string) error {
	if nombre == "" {
		return dominio.ErrValidacion("Nombre de usuario inválido")
	}
	return s.repo.Eliminar(ctx, nombre)
}

// Autenticar verifica credenciales. Nunca revela si el usuario existe o no.
func (s *UsuarioService) Autenticar(ctx context.Context, nombre, contraseña string) (dominio.Usuario, error) {
	if nombre == "" || contraseña == "" {
		return dominio.Usuario{}, dominio.ErrCredenciales
	}
	u, err := s.repo.Recuperar(ctx, nombre)
	if errors.Is(err, dominio.ErrUsuarioNoEncontrado) {
		return dominio.Usuario{}, dominio.ErrCredenciales
	}
	if err != nil {
		return dominio.Usuario{}, err
	}
	if bcrypt.CompareHashAndPassword([]byte(u.Contrasena), []byte(contraseña)) != nil {
		return dominio.Usuario{}, dominio.ErrCredenciales
	}
	return u, nil
}

func (s *UsuarioService) EsAdmin(ctx context.Context, nombre string) (bool, error) {
	if nombre == "" {
		return false, nil
	}
	return s.repo.EsAdmin(ctx, nombre)
}

func (s *UsuarioService) ObtenerPuntos(ctx context.Context, nombre string) (int, error) {
	return s.repo.ObtenerPuntos(ctx, nombre)
}

func (s *UsuarioService) SumarPuntos(ctx context.Context, nombre string, puntos int) error {
	if puntos <= 0 {
		return nil
	}
	return s.repo.SumarPuntos(ctx, nombre, puntos)
}

func (s *UsuarioService) ListarJugadores(ctx context.Context) ([]dominio.Jugador, error) {
	js, err := s.repo.ListarJugadores(ctx)
	if err != nil {
		return nil, err
	}
	// Más puntos primero; si empatan, por nombre.
	sort.Slice(js, func(i, j int) bool {
		if js[i].Puntos != js[j].Puntos {
			return js[i].Puntos > js[j].Puntos
		}
		return js[i].Nombre < js[j].Nombre
	})
	return js, nil
}
