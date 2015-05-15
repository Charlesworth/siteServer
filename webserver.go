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

	router.GET("/index", testGlob)

	http.Handle("/", router)

	log.Println("Listening...")
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func getTest(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	mock := PageViews{10, 20, 30}

	tmpl := template.Must(template.ParseFiles("tmpl/comments.html", "tmpl/foot.html")) //"tmpl/head.html", "tmpl/post.html", "tmpl/comments.html", "tmpl/foot.html"))

	err := tmpl.Execute(w, mock)
	if err != nil {
		fmt.Println(err)
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
