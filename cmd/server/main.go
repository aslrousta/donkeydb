package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"github.com/aslrousta/donkeydb"
)

var (
	dbPath = flag.String("db", "", "database file path")
	port   = flag.Int("port", 8080, "client port")
)

func main() {
	flag.Parse()
	if *dbPath == "" {
		flag.Usage()
		os.Exit(1)
	}
	file, isNew, err := openOrCreate(*dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	db, err := initDB(file, isNew)
	if err != nil {
		log.Fatal(err)
	}
	r := mux.NewRouter()
	r.Handle("/{key}", getHandler(db)).Methods(http.MethodGet)
	r.Handle("/{key}", setHandler(db)).Methods(http.MethodPost)
	r.Handle("/{key}", delHandler(db)).Methods(http.MethodDelete)
	http.ListenAndServe(fmt.Sprintf(":%d", *port), r)
}

func openOrCreate(path string) (*os.File, bool, error) {
	file, err := os.OpenFile(path, os.O_RDWR, 0)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, false, err
		}
		file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			return nil, false, err
		}
		return file, true, nil
	}
	return file, false, nil
}

func initDB(f *os.File, isNew bool) (*donkeydb.Database, error) {
	if isNew {
		return donkeydb.Create(f)
	}
	return donkeydb.Open(f)
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

func delHandler(d *donkeydb.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pathVars := mux.Vars(r)
		if err := d.Del(pathVars["key"]); err != nil {
			writeError(w, err)
			return
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
