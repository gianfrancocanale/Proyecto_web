package main

import (
	"fmt"
	"net/http"
)

func main() {
	// 1. Define el directorio que contiene los archivos estáticos.
	staticDir := "./static"

	// 2. Crea un manejador de archivos estáticos usando http.FileServer.

	fileServer := http.FileServer(http.Dir(staticDir))
	// el FileServer maneja el Content-type automáticamente según la extensión del archivo.

	// 3. Registra el manejador para que atienda todas las peticiones ("/").

	http.Handle("/", fileServer)

	// 4. Define el puerto y muestra un mensaje.
	port := ":8080"
	fmt.Printf("Servidor ESTÁTICO escuchando en http://localhost%s\n", port)
	fmt.Printf("Sirviendo archivos desde: %s\n", staticDir)

	// 5. Inicia el servidor.
	err := http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Printf("Error al iniciar el servidor: %s\n", err)
	}
	// aaaaz
}
