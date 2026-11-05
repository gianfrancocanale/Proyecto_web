package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"sync"
	"time"
)

const cookieSesion = "sesion"

var (
	sesiones   = map[string]string{} // token -> nombre de usuario
	sesionesMu sync.RWMutex
)

func crearSesion(w http.ResponseWriter, nombreUsuario string) error {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return err
	}
	token := hex.EncodeToString(b)

	sesionesMu.Lock()
	sesiones[token] = nombreUsuario
	sesionesMu.Unlock()

	http.SetCookie(w, &http.Cookie{
		Name:     cookieSesion,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(24 * time.Hour),
	})
	return nil
}

func usuarioActual(r *http.Request) (string, bool) {
	c, err := r.Cookie(cookieSesion)
	if err != nil {
		return "", false
	}
	sesionesMu.RLock()
	defer sesionesMu.RUnlock()
	u, ok := sesiones[c.Value]
	return u, ok
}

func cerrarSesion(w http.ResponseWriter, r *http.Request) {
	if c, err := r.Cookie(cookieSesion); err == nil {
		sesionesMu.Lock()
		delete(sesiones, c.Value)
		sesionesMu.Unlock()
	}
	http.SetCookie(w, &http.Cookie{
		Name: cookieSesion, Value: "", Path: "/", MaxAge: -1,
	})
}

// RequiereLogin envuelve un handler: si no hay sesión, manda al login.
func RequiereLogin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, ok := usuarioActual(r); !ok {
			http.Redirect(w, r, "/web/ingresar", http.StatusSeeOther)
			return
		}
		next(w, r)
	}
}
