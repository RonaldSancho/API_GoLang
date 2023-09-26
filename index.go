package main

import ( //modulos usados
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)


type tarea struct { //tipo de dato llamado tarea
    ID        int    `json:"ID"` //`json:"ID"` es la forma en como se le responderá al cliente 
    Nombre    string `json:"Nombre"` 
    Contenido string `json:"Contenido"` 
}

type ArregloTareas []tarea //se almacenan todas las tareas

var tareas = ArregloTareas{
	{
		ID: 1,
		Nombre: "Prueba",
		Contenido: "Prueba",
	},
}

func ConsultarTareas(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json"); //Información extra que indica el tipo de contenido que se envía
	json.NewEncoder(w).Encode(tareas) //se envían todas las tareas en formato JSON
}


func CrearTarea(w http.ResponseWriter, r *http.Request){
	var NuevaTarea tarea //la variable NuevaTarea está vacia pero es una instancia de tarea

	reqBody, err := ioutil.ReadAll(r.Body) //por medio del body se obtienen los datos que envia el usuario

	if err != nil{ //si se encuentra un error en el body se muestra el siguiente mensaje
		fmt.Fprintf(w, "Inserte los datos correspondientes.")
	}

	json.Unmarshal(reqBody, &NuevaTarea) // desde el reqBody se obtienen los datos y se agregan a NuevaTarea
	NuevaTarea.ID = len(tareas) + 1 //se incrementa el ID
	tareas = append(tareas, NuevaTarea) // a las tareas que ya existen se le agrega una nueva tarea

	w.Header().Set("Content-Type", "application/json"); //Información extra que indica el tipo de contenido que se envía
	w.WriteHeader(http.StatusCreated) //Información extra que indica el estado
	json.NewEncoder(w).Encode(NuevaTarea) //se le muestra al usuario los datos que creó
}


func ConsultarTarea(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r) //la variable vars guarda los parametros desde el Request

	TareaID, err := strconv.Atoi(vars["id"]) //con el Atoi convertimos de String a int
	if err != nil{ //si el usuario envia un texto inválido se mostrará el siguiente mensaje
		fmt.Fprintf(w, "El ID ingresado es inválido.")
		return
	}

	for _, tarea := range tareas{ //recorre la lista de tareas buscando un ID que coincida
		if tarea.ID == TareaID{
			w.Header().Set("Content-Type", "application/json"); //Información extra que indica el tipo de contenido que se envía
			json.NewEncoder(w).Encode(tarea)
		}
	}
}

func EliminarTarea(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r) //la variable vars guarda los parametros desde el Request

	TareaID, err := strconv.Atoi(vars["id"]) //con el Atoi convertimos de String a int

	if err != nil{ //si el usuario envia un texto inválido se mostrará el siguiente mensaje
		fmt.Fprintf(w, "El ID ingresado es inválido.")
		return
	}

	for i, tarea := range tareas{ //recorre la lista de tareas buscando un ID que coincida
		if tarea.ID == TareaID{
			tareas = append(tareas[:i], tareas[i + 1:]...) //Mantiene las tareas antes y despues del indice, creando así una nueva lista
			fmt.Fprintf(w, "La tarea fue borrada exitosamente.")
		}
	}
}

func ActualizarTarea(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r) //la variable vars guarda los parametros desde el Request
	TareaID, err := strconv.Atoi(vars["id"]) //con el Atoi convertimos de String a int
	var TareaActualizada tarea

	if err != nil{ //si el usuario envia un texto inválido se mostrará el siguiente mensaje
		fmt.Fprintf(w, "El ID ingresado es inválido.")
	}

	reqBody, err := ioutil.ReadAll(r.Body) //por medio del body se obtienen los datos que envia el usuario
	if err != nil{ //si se encuentra un error en el body se muestra el siguiente mensaje
		fmt.Fprintf(w, "Inserte los datos correspondientes.")
	}

	json.Unmarshal(reqBody, &TareaActualizada) // desde el reqBody se obtienen los datos y se agregan a TareaActualizada

	for i, tarea := range tareas{ //recorre la lista de tareas buscando un ID que coincida
		if tarea.ID == TareaID{
			tareas = append(tareas[:i], tareas[i + 1:]...) //Mantiene las tareas antes y despues del indice, creando así una nueva lista
			TareaActualizada.ID = TareaID
			tareas = append(tareas, TareaActualizada) // a las tareas que ya existen se le agrega una tarea actualizada
			fmt.Fprintf(w, "La tarea fue modificada exitosamente.")
		}
	}
}

func IndexRoute(w http.ResponseWriter, r *http.Request){ // http.ResponseWriter es la respuesta al usuario y http.Request es la petición
	fmt.Fprintf(w, "Hola clase de paradigmas");
}

func main(){
	router := mux.NewRouter().StrictSlash(true) //la URL debe ser escrita correctamente

	router.HandleFunc("/", IndexRoute) //se define el nombre de la ruta y la función que ejecutará
	router.HandleFunc("/tareas", ConsultarTareas).Methods("GET") //se le indica que se usa el metodo GET
	router.HandleFunc("/tareas", CrearTarea).Methods("POST") //se le indica que se usa el metodo POST
	router.HandleFunc("/tareas/{id}", ConsultarTarea).Methods("GET") //se ponen {} para indicar que despues del / viene un parametro
	router.HandleFunc("/tareas/{id}", EliminarTarea).Methods("DELETE") //se ponen {} para indicar que despues del / viene un parametro
	router.HandleFunc("/tareas/{id}", ActualizarTarea).Methods("PUT") //se ponen {} para indicar que despues del / viene un parametro
	log.Fatal(http.ListenAndServe(":3000", router))
	//Se crea el servidor, se agrega el puerto donde estará escuchando y el router
	//Se coloca dentro del Log.Fatal por si llega a ocurrir un error se muestre en consola

}