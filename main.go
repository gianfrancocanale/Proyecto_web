package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"Proyecto_web/aplicacion"
	h "Proyecto_web/handlers"
	i "Proyecto_web/infrastructura"
	v "Proyecto_web/views"

	"github.com/a-h/templ"
	_ "github.com/lib/pq"
)

func main() {
	db := conectarBase()
	defer db.Close()

	// Raíz de composición: los repositorios son la única capa que conoce SQL.
	escuderiaRepo := i.NewEscuderiaRepository(db)
	usuarioRepo := i.NewUsuarioRepository(db)
	pilotoRepo := i.NewPilotoRepository(db)
	pilotoHistoricoRepo := i.NewPilotoHistoricoRepository(db)
	granPremioRepo := i.NewGranPremioRepository(db)
	granPremioHistoricoRepo := i.NewGranPremioHistoricoRepository(db)
	apuestaRepo := i.NewApuestaRepository(db)

	// Los servicios de aplicación sólo conocen puertos del dominio.
	escuderiaSvc := aplicacion.NewEscuderiaService(escuderiaRepo)
	usuarioSvc := aplicacion.NewUsuarioService(usuarioRepo)
	pilotoSvc := aplicacion.NewPilotoService(pilotoRepo, escuderiaRepo)
	pilotoHistoricoSvc := aplicacion.NewPilotoHistoricoService(pilotoHistoricoRepo)
	granPremioSvc := aplicacion.NewGranPremioService(granPremioRepo, pilotoHistoricoRepo)
	granPremioHistoricoSvc := aplicacion.NewGranPremioHistoricoService(granPremioHistoricoRepo, granPremioRepo, pilotoHistoricoRepo)
	apuestaSvc := aplicacion.NewApuestaService(apuestaRepo, granPremioHistoricoRepo, pilotoRepo, pilotoHistoricoRepo)
	adminSvc := aplicacion.NewAdminService(usuarioRepo, granPremioRepo, granPremioHistoricoRepo, apuestaRepo, pilotoRepo)

	// Handlers: capa de interfaz de entrada, sin acceso directo a la BD.
	hEscuderia := h.NewEscuderiaHandler(escuderiaSvc, usuarioSvc)
	hUsuario := h.NewUsuarioHandler(usuarioSvc)
	hPiloto := h.NewPilotoHandler(pilotoSvc, usuarioSvc)
	hGranPremio := h.NewGranPremioHandler(granPremioSvc, pilotoHistoricoSvc, usuarioSvc)
	hApuesta := h.NewApuestaHandler(apuestaSvc, usuarioSvc, granPremioHistoricoSvc, pilotoSvc)
	hAdmin := h.NewAdminHandler(adminSvc)
	hPilotoHistorico := h.NewPilotoHistoricoHandler(pilotoHistoricoSvc)
	hGranPremioHistorico := h.NewGranPremioHistoricoHandler(granPremioHistoricoSvc)

	mux := http.NewServeMux()

	mux.Handle("GET /front/", http.StripPrefix("/front/", http.FileServer(http.Dir("front"))))

	mux.Handle("GET /{$}", templ.Handler(v.IndexPage()))
	mux.Handle("/", templ.Handler(v.NoEncontradoPage(), templ.WithStatus(http.StatusNotFound)))

	// Páginas HTML
	mux.HandleFunc("GET /escuderias", hEscuderia.Pagina)
	mux.HandleFunc("GET /pilotos", hPiloto.Pagina)
	mux.HandleFunc("GET /gran_premios", hGranPremio.Pagina)

	mux.HandleFunc("/ingresar", hUsuario.PaginaIngresar)
	mux.HandleFunc("/registrar", hUsuario.PaginaRegistrar)
	mux.HandleFunc("/salir", h.RequiereLogin(hUsuario.Salir))

	mux.HandleFunc("/apuestas", h.RequiereLogin(hApuesta.PaginaApuestas))
	mux.HandleFunc("GET /apuestas/modificar", h.RequiereLogin(hApuesta.ModificarApuesta))
	mux.HandleFunc("POST /apuestas/modificar", h.RequiereLogin(hApuesta.ModificarApuesta))
	mux.HandleFunc("DELETE /apuestas/{id}", h.RequiereLogin(hApuesta.EliminarApuesta))

	// Panel de administración
	mux.HandleFunc("/admin", hAdmin.RequiereAdmin(hAdmin.PaginaAdmin))
	mux.HandleFunc("/admin/carrera", hAdmin.RequiereAdmin(hAdmin.ProgramarCarrera))
	mux.HandleFunc("/admin/resultado", hAdmin.RequiereAdmin(hAdmin.CargarResultado))

	mux.HandleFunc("POST /pilotos", hAdmin.RequiereAdmin(hPiloto.Pagina))
	mux.HandleFunc("GET /pilotos/modificar", hAdmin.RequiereAdmin(hPiloto.ModificarPiloto))
	mux.HandleFunc("POST /pilotos/modificar", hAdmin.RequiereAdmin(hPiloto.ModificarPiloto))
	mux.HandleFunc("DELETE /pilotos/{id}", hAdmin.RequiereAdmin(hPiloto.EliminarPiloto))

	mux.HandleFunc("POST /escuderias", hAdmin.RequiereAdmin(hEscuderia.CrearEscuderia))
	mux.HandleFunc("DELETE /escuderias/{id}", hAdmin.RequiereAdmin(hEscuderia.EliminarEscuderia))
	mux.HandleFunc("GET /escuderias/modificar", hAdmin.RequiereAdmin(hEscuderia.PaginaModificar))
	mux.HandleFunc("POST /escuderias/modificar", hAdmin.RequiereAdmin(hEscuderia.ModificarEscuderia))

	mux.HandleFunc("POST /gran_premios", hAdmin.RequiereAdmin(hGranPremio.Pagina))
	mux.HandleFunc("GET /gran_premios/modificar", hAdmin.RequiereAdmin(hGranPremio.ModificarGranPremio))
	mux.HandleFunc("POST /gran_premios/modificar", hAdmin.RequiereAdmin(hGranPremio.ModificarGranPremio))
	mux.HandleFunc("DELETE /gran_premios/{id}", hAdmin.RequiereAdmin(hGranPremio.EliminarGranPremio))

	// API REST (los handlers ya existían, ahora detrás de servicios)
	mux.HandleFunc("GET /usuarios", hUsuario.ListarUsuarios)
	mux.HandleFunc("POST /usuarios", hUsuario.CrearUsuario)
	mux.HandleFunc("PUT /usuarios/{nombre}", hUsuario.ActualizarUsuario)
	mux.HandleFunc("DELETE /usuarios/{nombre}", hUsuario.EliminarUsuario)

	mux.HandleFunc("/pilotos_historicos", hPilotoHistorico.HandlerPilotoHistorico)
	mux.HandleFunc("/pilotos_historicos/{id}", hPilotoHistorico.HandlerPilotoHistoricoId)

	mux.HandleFunc("/gran_premios_historicos", hGranPremioHistorico.HandlerGranPremioHistorico)
	mux.HandleFunc("/gran_premios_historicos/{id}/{fecha_carrera}", hGranPremioHistorico.HandlerGranPremioHistoricoId)

	fmt.Println(" ")
	log.Println("Servidor escuchando en :8080")
	// ConSesion pone el nombre de usuario en el contexto para las vistas.
	log.Fatal(http.ListenAndServe(":8080", h.ConSesion(mux)))
}

func conectarBase() *sql.DB {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error al conectar con la base de datos: %v", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("Error al hacer ping a la base de datos: %v", err)
	}
	return db
}
