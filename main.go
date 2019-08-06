package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jaredfolkins/longstorm/app"
	"github.com/urfave/negroni"
)

func main() {

	// init middleware
	n := negroni.New()
	n.Use(negroni.NewRecovery())
	n.Use(negroni.NewLogger())

	// build router
	rtr := mux.NewRouter()
	rtr.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))
	rtr.HandleFunc("/", app.IndexHandler).Methods("GET")
	rtr.HandleFunc("/tweet", app.TweetHandler).Methods("POST")
	rtr.HandleFunc("/review", app.ReviewHandler).Methods("GET")
	rtr.HandleFunc("/long-storm", app.LongStormHandler).Methods("POST")
	rtr.HandleFunc("/settings", app.SettingsHandler).Methods("GET")
	rtr.HandleFunc("/settings", app.SaveSettingsHandler).Methods("POST")

	// add primary router to negroni
	n.Use(negroni.HandlerFunc(before))
	n.UseHandler(rtr)
	n.Use(negroni.HandlerFunc(after))

	log.Fatal(http.ListenAndServe(":3000", n))

}

func before(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	// TODO: is this needed?
	next(w, r)
}

func after(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	// TODO: is this needed?
	next(w, r)
}
