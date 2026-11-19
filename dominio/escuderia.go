package dominio

import (
	"context"
	"errors"
	"strings"
	"time"
)

// Errores de dominio
var (
	ErrDuplicada    = errors.New("ya existe una escudería con ese nombre")
	ErrNoEncontrada = errors.New("escudería no encontrada")
)

// ErrValidacion lleva el mensaje que se le puede mostrar al usuario.
type ErrValidacion string

func (e ErrValidacion) Error() string { return string(e) }

// Escuderia es la entidad del dominio.
type Escuderia struct {
	ID        string
	Principal string
	Titulos   int
	Puntos    int
	Fundacion time.Time
}

// Validar aplica las reglas de la entidad (y normaliza los textos).
func (e *Escuderia) Validar() error {
	e.ID = strings.TrimSpace(e.ID)
	e.Principal = strings.TrimSpace(e.Principal)

	if e.ID == "" || e.Principal == "" {
		return ErrValidacion("El nombre y el Team Principal son obligatorios")
	}
	if e.Titulos < 0 {
		return ErrValidacion("Los títulos tienen que ser un número mayor o igual a 0")
	}
	if e.Puntos < 0 {
		return ErrValidacion("Los puntos tienen que ser un número mayor o igual a 0")
	}
	return nil
}

// EscuderiaRepository es el puerto: lo define el dominio, lo implementa la infraestructura.
type EscuderiaRepository interface {
	Listar(ctx context.Context) ([]Escuderia, error)
	Recuperar(ctx context.Context, id string) (Escuderia, error)
	Crear(ctx context.Context, e Escuderia) error
	Modificar(ctx context.Context, e Escuderia) error
	Eliminar(ctx context.Context, id string) error
}
