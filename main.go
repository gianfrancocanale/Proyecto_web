package main

import (
	h "Proyecto_web/handlers"
	v "Proyecto_web/views"
	"database/sql" //para manejar errores(id no encontrados)
	"fmt"
	"net/http"
	"os" //obtener datos de variables de entorno(bases de datos)

	"github.com/a-h/templ"

	//convertir string a int

	_ "github.com/lib/pq"
)

func main() {
	db := conectarBase()
	fmt.Printf("La base de datos se inicio")
	defer db.Close()

	hUsuario := h.NewUsuarioHandler(db)
	hAdmin := h.NewAdminHandler(db)
	hPiloto := h.NewPilotoHandler(db)
	hEscuderia := h.NewEscuderiaHandler(db)
	hGranPremio := h.NewGranPremioHandler(db)
	hApuesta := h.NewApuestaHandler(db)

	http.Handle("GET /front/", http.StripPrefix("/front/", http.FileServer(http.Dir("front"))))

	http.Handle("GET /{$}", templ.Handler(v.IndexPage()))
	http.Handle("/", templ.Handler(v.NoEncontradoPage(), templ.WithStatus(http.StatusNotFound)))

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

	//ENDPINT DE ADMINISTRADOR

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

	fmt.Println(" ")
	fmt.Println("Servidor escuchando en :8080")
	http.ListenAndServe(":8080", h.ConSesion(http.DefaultServeMux)) //la parte del mux es por que lo obliga a pasar por el middleware de sesion,
	//para que pueda acceder a la cookie y ver si esta logueado o no
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
		fmt.Errorf("Error al conectar con la base de datos: %v", err)
	}

	err = db.Ping()
	if err != nil {
		fmt.Errorf("Error al hacer ping a la base de datos: %v", err)
	}

	return db
}
