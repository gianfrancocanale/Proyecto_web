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

	http.Handle("GET /front/", http.StripPrefix("/front/", http.FileServer(http.Dir("front"))))

	http.Handle("GET /{$}", templ.Handler(v.IndexPage()))
	http.Handle("/", templ.Handler(v.NoEncontradoPage(), templ.WithStatus(http.StatusNotFound)))

	// Páginas HTML
	http.HandleFunc("GET /escuderias", hEscuderia.Pagina)
	http.HandleFunc("GET /pilotos", hPiloto.Pagina)
	http.HandleFunc("GET /gran_premios", hGranPremio.Pagina)

	http.HandleFunc("/ingresar", hUsuario.PaginaIngresar)
	http.HandleFunc("/registrar", hUsuario.PaginaRegistrar)
	http.HandleFunc("/salir", h.RequiereLogin(hUsuario.Salir))

	http.HandleFunc("/apuestas", h.RequiereLogin(hApuesta.PaginaApuestas))
	http.HandleFunc("GET /apuestas/modificar", h.RequiereLogin(hApuesta.ModificarApuesta))
	http.HandleFunc("POST /apuestas/modificar", h.RequiereLogin(hApuesta.ModificarApuesta))
	http.HandleFunc("DELETE /apuestas/{id}", h.RequiereLogin(hApuesta.EliminarApuesta))

	// Panel de administración
	http.HandleFunc("/admin", hAdmin.RequiereAdmin(hAdmin.PaginaAdmin))
	http.HandleFunc("/admin/carrera", hAdmin.RequiereAdmin(hAdmin.ProgramarCarrera))
	http.HandleFunc("/admin/resultado", hAdmin.RequiereAdmin(hAdmin.CargarResultado))

	http.HandleFunc("POST /pilotos", hAdmin.RequiereAdmin(hPiloto.Pagina))
	http.HandleFunc("GET /pilotos/modificar", hAdmin.RequiereAdmin(hPiloto.ModificarPiloto))
	http.HandleFunc("POST /pilotos/modificar", hAdmin.RequiereAdmin(hPiloto.ModificarPiloto))
	http.HandleFunc("DELETE /pilotos/{id}", hAdmin.RequiereAdmin(hPiloto.EliminarPiloto))

	http.HandleFunc("POST /escuderias", hAdmin.RequiereAdmin(hEscuderia.CrearEscuderia))
	http.HandleFunc("DELETE /escuderias/{id}", hAdmin.RequiereAdmin(hEscuderia.EliminarEscuderia))
	http.HandleFunc("GET /escuderias/modificar", hAdmin.RequiereAdmin(hEscuderia.PaginaModificar))
	http.HandleFunc("POST /escuderias/modificar", hAdmin.RequiereAdmin(hEscuderia.ModificarEscuderia))

	http.HandleFunc("POST /gran_premios", hAdmin.RequiereAdmin(hGranPremio.Pagina))
	http.HandleFunc("GET /gran_premios/modificar", hAdmin.RequiereAdmin(hGranPremio.ModificarGranPremio))
	http.HandleFunc("POST /gran_premios/modificar", hAdmin.RequiereAdmin(hGranPremio.ModificarGranPremio))
	http.HandleFunc("DELETE /gran_premios/{id}", hAdmin.RequiereAdmin(hGranPremio.EliminarGranPremio))

	// API REST (los handlers ya existían, ahora detrás de servicios)
	http.HandleFunc("GET /usuarios", hUsuario.ListarUsuarios)
	http.HandleFunc("POST /usuarios", hUsuario.CrearUsuario)
	http.HandleFunc("PUT /usuarios/{nombre}", hUsuario.ActualizarUsuario)
	http.HandleFunc("DELETE /usuarios/{nombre}", hUsuario.EliminarUsuario)

	http.HandleFunc("/pilotos_historicos", hPilotoHistorico.HandlerPilotoHistorico)
	http.HandleFunc("/pilotos_historicos/{id}", hPilotoHistorico.HandlerPilotoHistoricoId)

	http.HandleFunc("/gran_premios_historicos", hGranPremioHistorico.HandlerGranPremioHistorico)
	http.HandleFunc("/gran_premios_historicos/{id}/{fecha_carrera}", hGranPremioHistorico.HandlerGranPremioHistoricoId)

	fmt.Println(" ")
	log.Println("Servidor escuchando en :8080")
	// ConSesion pone el nombre de usuario en el contexto para las vistas.
	log.Fatal(http.ListenAndServe(":8080", h.ConSesion(http.DefaultServeMux)))
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
