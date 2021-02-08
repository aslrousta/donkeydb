package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/aslrousta/donkeydb"
	"github.com/gorilla/mux"
)

func main() {
	d := donkeydb.New()

	r := mux.NewRouter()
	r.Handle("/{key}", getHandler(d)).Methods(http.MethodGet)
	r.Handle("/{key}", setHandler(d)).Methods(http.MethodPost)
	http.ListenAndServe(":8080", r)
}

func getHandler(d *donkeydb.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pathVars := mux.Vars(r)
		value, err := d.Get(pathVars["key"])
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, &struct {
			Value interface{} `json:"value"`
		}{Value: value})
	}
}

func setHandler(d *donkeydb.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		pathVars := mux.Vars(r)
		value, err := ioutil.ReadAll(r.Body)
		if err != nil {
			writeError(w, err)
			return
		}
		if err := d.Set(pathVars["key"], string(value)); err != nil {
			writeError(w, err)
		}
	}
}

func writeError(w http.ResponseWriter, err error) {
	if err == donkeydb.ErrNothing {
		w.WriteHeader(http.StatusNotFound)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
	writeJSON(w, &struct {
		Message string `json:"message"`
	}{Message: err.Error()})
}

func writeJSON(w http.ResponseWriter, data interface{}) {
	encoder := json.NewEncoder(w)
	encoder.Encode(data)
}
