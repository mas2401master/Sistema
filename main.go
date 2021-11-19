package main

//quede en minuto 31
//https://www.youtube.com/watch?v=G58gN0lIbyI
// go mod init SISTEMA para organizar los paquetes que se usaran el el proyecto
// go get -u github.com/go-sql-driver/mysql para descargar el driver de sql
import (
	//"fmt"
	"database/sql" //operaciones SQL
	"log"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql" //Conectarse a mysql
)

var plantillas = template.Must(template.ParseGlob("plantillas/*"))

type Empleado struct {
	Id     int
	Nombre string
	Correo string
}

func NewDB() (*sql.DB, error) {
	Driver := "mysql"
	Usuario := "root"
	Contrasenia := ""
	Nombre := "golang_empleados"

	db, err := sql.Open(Driver, Usuario+":"+Contrasenia+"@tcp(127.0.0.1)/"+Nombre)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
func main() {
	log.Println("Iniciando el Servidor...")
	http.HandleFunc("/", Inicio)
	http.HandleFunc("/crear", Crear)
	http.HandleFunc("/insertar", Insertar)
	http.HandleFunc("/borrar", Borrar)
	http.HandleFunc("/editar", Editar)
	http.HandleFunc("/actualizar", Actualizar)
	http.ListenAndServe(":8080", nil)
}

func Inicio(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintln(w, "Mi sitio Web con GoLang")
	conec, error := NewDB()
	defer conec.Close()
	registros, error := conec.Query("SELECT * FROM empleados")
	if error != nil {
		panic(error.Error())
	}
	empleado := Empleado{}
	listEmpleados := []Empleado{}

	for registros.Next() {
		error = registros.Scan(&empleado.Id, &empleado.Nombre, &empleado.Correo)
		if error != nil {
			panic(error.Error())
		}
		listEmpleados = append(listEmpleados, empleado)
	}
	//fmt.Println(listEmpleados)
	plantillas.ExecuteTemplate(w, "inicio", listEmpleados)
}

func Crear(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintln(w, "Mi sitio Web con GoLang")
	plantillas.ExecuteTemplate(w, "crear", nil)
}

func Insertar(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		nombre := r.FormValue("nombre")
		correo := r.FormValue("correo")
		db, err := NewDB()
		defer db.Close()
		if err != nil {
			panic(err.Error)
		}
		insertar, error := db.Prepare("INSERT INTO empleados(nombre,correo) VALUES(?, ?)")
		if error != nil {
			panic(error.Error())
		}
		insertar.Exec(nombre, correo)
		http.Redirect(w, r, "/", 301)
	}
}

func Borrar(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	db, err := NewDB()
	defer db.Close()
	if err != nil {
		panic(err.Error)
	}
	borrar, error := db.Prepare("DELETE FROM empleados WHERE id=?")
	if error != nil {
		panic(error.Error())
	}
	borrar.Exec(id)
	http.Redirect(w, r, "/", 301)
}

func Editar(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	conec, err := NewDB()
	defer conec.Close()
	if err != nil {
		panic(err.Error)
	}
	registro, error := conec.Query("SELECT * FROM empleados WHERE id=?", id)
	if error != nil {
		panic(error.Error())
	}
	empleado := Empleado{}
	for registro.Next() {
		error = registro.Scan(&empleado.Id, &empleado.Nombre, &empleado.Correo)
		if error != nil {
			panic(error.Error())
		}
	}
	plantillas.ExecuteTemplate(w, "editar", empleado)
	//http.Redirect(w, r, "/editar", empleado)
}

func Actualizar(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		id := r.FormValue("id")
		nombre := r.FormValue("nombre")
		correo := r.FormValue("correo")
		db, err := NewDB()

		if err != nil {
			panic(err.Error)
		}
		actualizar, error := db.Prepare("UPDATE empleados SET nombre=?,correo=? WHERE id=?")
		if error != nil {
			panic(error.Error())
		}
		actualizar.Exec(nombre, correo, id)
		defer db.Close()
		http.Redirect(w, r, "/", 301)
	}
}
