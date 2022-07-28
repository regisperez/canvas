package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/BurntSushi/toml"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strings"
	"time"
)

type Database struct {
	Server       string
	Port         string
	Database     string
	User         string
	Password     string
}

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (a *App) Initialize(config string) {

	var (
		database Database
	)
	if _, err := toml.DecodeFile(config, &database); err != nil {
		fmt.Println(err)
	}

	fmt.Println(config)

	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", database.User, database.Password, database.Server, database.Port, database.Database)

	var err error
	a.DB, err = sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	a.Router = mux.NewRouter()

	a.initializeRoutes()
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func (a *App) getCanvas(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	canvas := Canvas{ID: id}
	if err := canvas.GetCanvas(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Canvas not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, canvas)
}

func (a *App) getCanvasList(w http.ResponseWriter, r *http.Request) {

	canvasList, err := GetCanvasList(a.DB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, canvasList)
}

func (a *App) createCanvas(w http.ResponseWriter, r *http.Request) {
	var canvas Canvas
	canvas.ID = uuid.New().String()
	canvas.CreationDate = time.Now().Format("2006-01-02 15:04:05")
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&canvas); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := canvas.CreateCanvas(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, canvas)
}

func (a *App) updateCanvas(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var canvas Canvas
	canvas.CreationDate = time.Now().Format("2006-01-02 15:04:05")
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&canvas); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}
	defer r.Body.Close()
	canvas.ID = id

	if err := canvas.UpdateCanvas(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, canvas)
}

func (a *App) deleteCanvas(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id:= vars["id"]

	canvas := Canvas{ID: id}
	if err := canvas.DeleteCanvas(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) createCanvasRequest(w http.ResponseWriter, r *http.Request) {
	var (
		canvas        Canvas
		canvasRequest CanvasCreateRequest
		err           error
	)
	canvas.ID = uuid.New().String()
	canvas.CreationDate = time.Now().Format("2006-01-02 15:04:05")
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&canvasRequest); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	canvas.Drawing, err = CanvasCreate(&canvasRequest)

	if err := canvas.CreateCanvas(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, err.Error())
	}else{
		respondWithJSON(w, http.StatusCreated, CanvasResponse{
			ID:           canvas.ID,
			Drawing:      strings.Split(canvas.Drawing, "\n"),
			CreationDate: canvas.CreationDate,
		})
	}

}

func (a *App) getCanvasResponse(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	canvas := Canvas{ID: id}
	if err := canvas.GetCanvas(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Canvas not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, CanvasResponse{
		ID:           canvas.ID,
		Drawing:      strings.Split(canvas.Drawing, "\n"),
		CreationDate: canvas.CreationDate,
	})
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/canvasList", a.getCanvasList).Methods("GET")
	a.Router.HandleFunc("/canvas", a.createCanvas).Methods("POST")
	a.Router.HandleFunc("/canvas/{id:[A-Za-z0-9\\W]+}", a.getCanvas).Methods("GET")
	a.Router.HandleFunc("/canvas/{id:[A-Za-z0-9\\W]+}", a.updateCanvas).Methods("PUT")
	a.Router.HandleFunc("/canvas/{id:[A-Za-z0-9\\W]+}", a.deleteCanvas).Methods("DELETE")
	a.Router.HandleFunc("/canvasCreateRequest", a.createCanvasRequest).Methods("POST")
	a.Router.HandleFunc("/canvasResponse/{id:[A-Za-z0-9\\W]+}", a.getCanvasResponse).Methods("GET")
}
