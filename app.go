package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	//"strings"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

type api struct {
	DB *sql.DB
}

func (a *App) Initialize(user, password, dbname string) {
	connectionString := fmt.Sprintf("%s:%s@/%s", user, password, dbname)
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

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (a *App) lawFirmActions_createNew(w http.ResponseWriter, r *http.Request) {
	var lfa lawFirmActions

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&lfa); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()
	if err := lfa.lawFirmAction_createNew(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusCreated, lfa)
}

func (a *App) lawFirmActions_getOne(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}
	lfa := lawFirmActions{ID: id}
	if err := lfa.lawFirmActions_getOne(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "User not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	respondWithJSON(w, http.StatusOK, lfa)
}

func (a *App) lawFirmActions_getAll(w http.ResponseWriter, r *http.Request) {
	count, _ := strconv.ParseInt(r.FormValue("count"), 10, 64)
	start, _ := strconv.ParseInt(r.FormValue("start"), 10, 64)

	if count > 10 || count < 1 {
		count = 10
	}
	if start < 0 {
		start = 0
	}
	res, err := lawFirmActions_getAll(a.DB, start, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, res)
}

func (a *App) lawFirmActions_update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}
	var lfa lawFirmActions
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&lfa); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}
	defer r.Body.Close()

	lfa.ID = id
	if err := lfa.lawFirmActions_update(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, lfa)
}

func (a *App) lawFirmActions_delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid User ID")
		return
	}

	lfa := lawFirmActions{ID: id}
	if err := lfa.lawFirmActions_delete(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

//-----------------------------------------------------------------------------------------------------
func (a *App) budgetCategory_createNew(w http.ResponseWriter, r *http.Request) {
	var b budgetCategory

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&b); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()
	if err := b.budgetCategory_createNew(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusCreated, b)
}

func (a *App) budgetCategory_getOne(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}
	b := budgetCategory{ID: id}
	if err := b.budgetCategory_getOne(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "User not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	respondWithJSON(w, http.StatusOK, b)
}

func (a *App) budgetCategory_getAll(w http.ResponseWriter, r *http.Request) {
	count, _ := strconv.ParseInt(r.FormValue("count"), 10, 64)
	start, _ := strconv.ParseInt(r.FormValue("start"), 10, 64)

	if count > 10 || count < 1 {
		count = 10
	}
	if start < 0 {
		start = 0
	}
	res, err := budgetCategory_getAll(a.DB, start, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, res)
}

func (a *App) budgetCategory_update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}
	var b budgetCategory
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&b); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}
	defer r.Body.Close()

	b.ID = id
	if err := b.budgetCategory_update(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, b)
}

func (a *App) budgetCategory_delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid User ID")
		return
	}

	b := budgetCategory{ID: id}
	if err := b.budgetCategory_delete(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

//------------------------------------------------------------------------------------------
func (a *App) lawFirmCost_createNew(w http.ResponseWriter, r *http.Request) {
	var lfc lawFirmCost

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&lfc); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	defer r.Body.Close()
	if err := lfc.lawFirmCost_createNew(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusCreated, lfc)
}

func (a *App) lawFirmCost_getOne(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}
	lfc := lawFirmCost{ID: id}
	if err := lfc.lawFirmCost_getOne(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "User not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	respondWithJSON(w, http.StatusOK, lfc)
}

func (a *App) lawFirmCost_getAll(w http.ResponseWriter, r *http.Request) {
	count, _ := strconv.ParseInt(r.FormValue("count"), 10, 64)
	start, _ := strconv.ParseInt(r.FormValue("start"), 10, 64)

	if count > 10 || count < 1 {
		count = 10
	}
	if start < 0 {
		start = 0
	}
	res, err := lawFirmCost_getAll(a.DB, start, count)
	if err != nil {

		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, res)
}

func (a *App) lawFirmCost_update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}
	var lfc lawFirmCost
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&lfc); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}
	defer r.Body.Close()

	lfc.ID = id
	if err := lfc.lawFirmCost_update(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, lfc)
}

func (a *App) lawFirmCost_delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid User ID")
		return
	}

	lfc := lawFirmCost{ID: id}
	if err := lfc.lawFirmCost_delete(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

//-----------------------------------------------------------------------------------------------------
func (a *App) lawFirmCostUsptoFees_createNew(w http.ResponseWriter, r *http.Request) {
	var lfcuf lawFirmCostUsptoFees

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&lfcuf); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()
	if err := lfcuf.lawFirmCostUsptoFees_createNew(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusCreated, lfcuf)
}

func (a *App) lawFirmCostUsptoFees_getOne(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}
	lfcuf := lawFirmCostUsptoFees{ID: id}
	if err := lfcuf.lawFirmCostUsptoFees_getOne(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "User not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	respondWithJSON(w, http.StatusOK, lfcuf)
}

func (a *App) lawFirmCostUsptoFees_getAll(w http.ResponseWriter, r *http.Request) {
	count, _ := strconv.ParseInt(r.FormValue("count"), 10, 64)
	start, _ := strconv.ParseInt(r.FormValue("start"), 10, 64)

	if count > 10 || count < 1 {
		count = 10
	}
	if start < 0 {
		start = 0
	}
	res, err := lawFirmCostUsptoFees_getAll(a.DB, start, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, res)
}

func (a *App) lawFirmCostUsptoFees_update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}
	var lfcuf lawFirmCostUsptoFees
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&lfcuf); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}
	defer r.Body.Close()

	lfcuf.ID = id
	if err := lfcuf.lawFirmCostUsptoFees_update(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, lfcuf)
}

func (a *App) lawFirmCostUsptoFees_delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid User ID")
		return
	}

	lfcuf := lawFirmCostUsptoFees{ID: id}
	if err := lfcuf.lawFirmCostUsptoFees_delete(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/budgetcategories", a.budgetCategory_createNew).Methods("POST")
	a.Router.HandleFunc("/budgetcategories", a.budgetCategory_getAll).Methods("GET")
	a.Router.HandleFunc("/budgetcategories/{id:[0-9]+}", a.budgetCategory_getOne).Methods("GET")
	a.Router.HandleFunc("/budgetcategories/{id:[0-9]+}", a.budgetCategory_update).Methods("PUT")
	a.Router.HandleFunc("/budgetcategories/{id:[0-9]+}", a.budgetCategory_delete).Methods("DELETE")
	a.Router.HandleFunc("/lawfirmcost", a.lawFirmCost_createNew).Methods("POST")
	a.Router.HandleFunc("/lawfirmcost", a.lawFirmCost_getAll).Methods("GET")
	a.Router.HandleFunc("/lawfirmcost/{id:[0-9]+}", a.lawFirmCost_getOne).Methods("GET")
}
