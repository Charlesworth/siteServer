package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"html/template"
	"log"
	"net/http"
)

type PageViews struct {
	Index    int
	Contact  int
	Projects int
}

func main() {
	fmt.Println("webserver started")

	router := httprouter.New()

	router.GET("/", getTmpl)

	http.Handle("/", router)

	log.Println("Listening...")
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func get(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	fmt.Fprint(w, "test")
}

func getTmpl(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	mock := PageViews{10, 20, 30}

	tmpl := template.Must(template.ParseFiles("tmpl/1.html"))

	err := tmpl.Execute(w, mock)
	if err != nil {
		fmt.Println(err)
	}

}
