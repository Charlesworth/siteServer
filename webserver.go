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

	router.GET("/Glob", testGlob)
	router.GET("/Files", testFiles)

	http.Handle("/", router)

	log.Println("Listening...")
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func testFiles(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	mock := PageViews{10, 20, 30}

	tmpl := template.Must(template.ParseFiles("tmpl/index.html", "tmpl/head.html"))

	//ExecuteTemplate writes the template to w, writing "indexPage" as the main as defined
	//in index.html, and with a data interface
	err := tmpl.ExecuteTemplate(w, "indexPage", mock)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func testGlob(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	templates := template.Must(template.ParseGlob("tmpl/*"))

	err := templates.ExecuteTemplate(w, "indexPage", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
