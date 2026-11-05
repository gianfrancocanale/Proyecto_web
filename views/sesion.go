package views

import "context"

type ctxKey string

const UsuarioKey ctxKey = "usuario"

func logueado(ctx context.Context) bool {
	_, ok := ctx.Value(UsuarioKey).(string)
	return ok
}
