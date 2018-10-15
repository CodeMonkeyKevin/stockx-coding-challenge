// app.go

package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"strconv"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

type shoeJsonPayload struct {
	ID            int    `json:"id"`
	Shoe          string `json:"shoe"`
	TrueToSizeVal int    `json:"trueToSizeVal"`
}

func (a *App) Initialize(user, password, dbname string) {
	connectionString :=
		fmt.Sprintf("sslmode=disable dbname=%s user=%s password='%s'", dbname, user, password)

	var err error
	a.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

func (a *App) Run(addr string) {
	log.Println("Starting Server...")
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/", a.getShoes).Methods("GET")
	a.Router.HandleFunc("/shoes", a.getShoes).Methods("GET")
	a.Router.HandleFunc("/shoes", a.addShoeData).Methods("POST")
	a.Router.HandleFunc("/shoes/{id:[0-9]+}", a.getShoe).Methods("GET")
	a.Router.HandleFunc("/shoes/{id:[0-9]+}", a.deleteShoe).Methods("DELETE")
}

func (a *App) getShoe(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid  ID")
		return
	}

	s := shoe{ID: id}
	if err := s.findById(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Shoe not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, s)
}

func (a *App) getShoes(w http.ResponseWriter, r *http.Request) {
	shoes, err := getShoes(a.DB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, shoes)
}

func (a *App) addShoeData(w http.ResponseWriter, r *http.Request) {
	var s shoeJsonPayload
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&s); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	var rShoe shoe
	rShoe, err := getOrCreateShoeByName(a.DB, s.Shoe)

	if err != nil {
		log.Fatal(err)
		return
	}

	err = rShoe.updateShoe(a.DB, s.TrueToSizeVal)

	if err != nil {
		log.Fatal(err)
	}

	respondWithJSON(w, http.StatusCreated, rShoe)
}

func (a *App) deleteShoe(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Shoe ID")
		return
	}

	s := shoe{ID: id}
	if err := s.deleteShoe(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message, "result": "error"})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
