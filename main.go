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

	http.Handle("/", http.FileServer(http.Dir("front")))

	hUsuario := h.NewUsuarioHandler(db)
	hPilotoHistorico := h.NewPilotoHistoricoHandler(db)
	hPiloto := h.NewPilotoHandler(db)
	hEscuderia := h.NewEscuderiaHandler(db)
	hGranPremio := h.NewGranPremioHandler(db)
	hApuesta := h.NewApuestaHandler(db)
	hGranPremioHistorico := h.NewGranPremioHistoricoHandler(db)

	http.Handle("GET /{$}", templ.Handler(v.IndexPage()))
	http.HandleFunc("GET /web/escuderias", hEscuderia.Pagina)

	// MANEJO DE RUTAS DE ENTIDAD USUARIO:
	http.HandleFunc("/usuarios", hUsuario.HandlerUsuario)
	http.HandleFunc("/usuarios/{nombre}", hUsuario.HandlerUsuarioId)

	//MANEJO DE RUTAS DE ENTIDAD PILOTO HISTORICO:
	http.HandleFunc("/pilotos_historicos", hPilotoHistorico.HandlerPilotoHistorico)
	http.HandleFunc("/pilotos_historicos/{id}", hPilotoHistorico.HandlerPilotoHistoricoId)

	//MANEJO DE RUTAS DE ENTIDAD ESCUDERIA:
	http.HandleFunc("/escuderias", hEscuderia.HandlerEscuderia)
	http.HandleFunc("/escuderias/{id}", hEscuderia.HandlerEscuderiaId)

	//MANEJO DE RUTAS DE ENTIDAD PILOTO:
	http.HandleFunc("/pilotos", hPiloto.HandlerPiloto)
	http.HandleFunc("/pilotos/{id}", hPiloto.HandlerPilotoId)

	//MANEJO DE RUTAS DE ENTIDAD GRAN PREMIO:
	http.HandleFunc("/gran_premios", hGranPremio.HandlerGranPremio)
	http.HandleFunc("/gran_premios/{id}", hGranPremio.HandlerGranPremioId)

	//MANEJO DE RUTAS DE ENTIDAD APUESTA:
	http.HandleFunc("/apuestas", hApuesta.HandlerApuesta)
	http.HandleFunc("/apuestas/{id}/{usuario}", hApuesta.HandlerApuestaIdUsuario)
	http.HandleFunc("/apuestas/{usuario}", hApuesta.HandlerApuestaUsuario)

	//MANEJO DE RUTAS DE ENTIDAD GRAN PREMIO HISTORICO:
	http.HandleFunc("/gran_premios_historicos", hGranPremioHistorico.HandlerGranPremioHistorico)
	http.HandleFunc("/gran_premios_historicos/{id}/{fecha_carrera}", hGranPremioHistorico.HandlerGranPremioHistoricoId)

	fmt.Println(" ")
	fmt.Println("Servidor escuchando en :8080")
	http.ListenAndServe(":8080", nil)
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
