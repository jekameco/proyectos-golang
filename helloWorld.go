/*nombre de la carpeta que esta guardado, siempre al archivo principal es main*/
package main

/*importar funciones de go para usar dentro del codigo*/
import (
	// mostrar mensajes
	"log"           // mosttrar informacion desde la terminal
	"net/http"      // mostrar  un sitio
	"text/template" //modulos de plantillas
)

/*obtener informacion en una carpeta , y buscar la informacion de los templates*/
var plantillas = template.Must(template.ParseGlob("plantillas/*")) //acceder al folder

/*Funcion inicial*/
func main() {

	http.HandleFunc("/", Inicio) //decirle al usuario que escriba en el navegador la ruta- accceder la funcion
	http.HandleFunc("/crear", Crear)
	log.Println("Servidor corriendo...") // impresion en consola
	http.ListenAndServe(":8080", nil)    // puerto del servidor

}

/*funcion que se encuentra dentro de HandleFunc, el primer parametro es la respuesta de la solicitud y el segundo parametro me da la informacion del usuario GET o POST*/
func Inicio(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "Hola Develoteca") // metodo para mostrar en el navegador
	plantillas.ExecuteTemplate(w, "inicio", nil) // reconocer directamente el archivo
}

func Crear(w http.ResponseWriter, r *http.Request) {
	plantillas.ExecuteTemplate(w, "crear", nil) // reconocer directamente el archivo
}
