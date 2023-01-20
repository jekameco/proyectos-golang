/******************/
/*comandos necesarios */

/*
go mod init (nombre de la carpeta)  -> poner carpetas o driver en la carpeta automaticamente
go get -u github.com/go-sql-driver/mysql  conectar a base de datos
*/
/******************/

/*nombre de la carpeta que esta guardado, siempre al archivo principal es main*/
package main

/*importar funciones de go para usar dentro del codigo*/
import (
	//"fmt"// mostrar mensajes
	"database/sql"
	"log"           // mosttrar informacion desde la terminal
	"net/http"      // mostrar  un sitio
	"text/template" //modulos de plantillas

	_ "github.com/go-sql-driver/mysql" //forma de insertar driver iniciar con raya al piso
)

/*funcion para conectar con la base de datos  nombre-de-la funcion(lo que se obtiene)(lo que se devuelve) */
func conexionBD() (conexion *sql.DB) {
	/*datos necesario para conectarse a la base de datos*/
	Driver := "mysql"
	Usuario := "root"
	Contraseña := ""
	Nombre := "sistema"

	/*conectandose a la base de datos*/
	/*Driver,
	Usuario+":"
	+Contraseña+"@tcp(127.0.0.1)/" -> direccion ip en este momento de localhost, pero si es diferente se cambiara
	+Nombre*/
	conexion, err := sql.Open(Driver, Usuario+":"+Contraseña+"@tcp(127.0.0.1)/"+Nombre)

	if err != nil {
		panic(err.Error())
	}

	return conexion
}

/*obtener informacion en una carpeta , y buscar la informacion de los templates*/
var plantillas = template.Must(template.ParseGlob("plantillas/*")) //acceder al folder

/*Funcion inicial*/
func main() {
	http.HandleFunc("/", Inicio) //decirle al usuario que escriba en el navegador la ruta- accceder la funcion
	http.HandleFunc("/crear", Crear)
	http.HandleFunc("/insertar", Insertar)
	http.HandleFunc("/borrar", Borrar)
	http.HandleFunc("/actualizar", Actualizar)
	log.Println("Servidor corriendo...") // impresion en consola
	http.ListenAndServe(":8080", nil)    // puerto del servidor

}

/*depositar la estructra de cada empleado es usado hasta el momento para los select*/

type Empleado struct {
	Id     int
	Nombre string
	Correo string
}

/*funcion que se encuentra dentro de HandleFunc, el primer parametro es la respuesta de la solicitud y el segundo parametro me da la informacion del usuario GET o POST*/
func Inicio(w http.ResponseWriter, r *http.Request) {
	conexionEstablecida := conexionBD() // establecer conexion con bd

	registros, err := conexionEstablecida.Query("SELECT * FROM empleados") // realizar la consulta sin prepare, directamente query (no se necesita execute)
	if err != nil {
		panic(err.Error())
	}
	empleado := Empleado{}          //union de arreglo empleado
	arregloEmpleado := []Empleado{} // guardar datos por cada fila

	/* recorrer el select fila por fila*/
	for registros.Next() {
		var id int
		var nombre, correo string
		/* asignacion de los valores*/
		err = registros.Scan(&id, &nombre, &correo)
		if err != nil {
			panic(err.Error())
		}

		empleado.Id = id
		empleado.Nombre = nombre
		empleado.Correo = correo

		arregloEmpleado = append(arregloEmpleado, empleado) // insertar al  arreglo cada fila de la bd
	}

	//fmt.Println(arregloEmpleado)

	//fmt.Fprintf(w, "Hola Develoteca") // metodo para mostrar en el navegador
	plantillas.ExecuteTemplate(w, "inicio", arregloEmpleado) // reconocer directamente el archivo- tercer parametro estamos enviando informacion al html
}

func Crear(w http.ResponseWriter, r *http.Request) {
	plantillas.ExecuteTemplate(w, "crear", nil) // reconocer directamente el archivo
}

func Insertar(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		nombre := r.FormValue("name") // forma de obtener los valores de los inputs
		correo := r.FormValue("email")

		conexionEstablecida := conexionBD() // establecer conexion con bd

		insertarRegistro, err := conexionEstablecida.Prepare("INSERT INTO `empleados`(`nombre`, `correo`) VALUES (?,?)") // realizar una insercion a la base de datos
		if err != nil {
			panic(err.Error())
		}

		insertarRegistro.Exec(nombre, correo) // poner los inputs dentro del insert, se identifica por un signo de incognito ?
		http.Redirect(w, r, "/", 301)         //forma de redireccionar

	}
}

func Actualizar(w http.ResponseWriter, r *http.Request) {
	plantillas.ExecuteTemplate(w, "crear", nil) // reconocer directamente el archivo
}

func Borrar(w http.ResponseWriter, r *http.Request) {
	idEmpleado := r.URL.Query().Get("id") //recibir accion por metodo get
	//fmt.Println(idEmpleado)
	conexionEstablecida := conexionBD() // establecer conexion con bd

	eliminarRegistro, err := conexionEstablecida.Prepare("DELETE FROM `empleados` WHERE id=?") // realizar una insercion a la base de datos
	if err != nil {
		panic(err.Error())
	}
	eliminarRegistro.Exec(idEmpleado)

	http.Redirect(w, r, "/", 301)
}
