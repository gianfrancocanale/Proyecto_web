package infrastructura

import (
	"database/sql"
	"errors"
	"strings"
)

// Los errores de Postgres llegan como texto; acá se traducen a errores del dominio.
func esNoRows(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}

func esDuplicado(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "duplicate key") ||
		strings.Contains(msg, "unique constraint") ||
		strings.Contains(msg, "violates unique")
}

func esForeignKey(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(strings.ToLower(err.Error()), "foreign key")
}
