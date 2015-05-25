package main

import (
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/Charlesworth/viewsLib"
	"github.com/julienschmidt/httprouter"
)

//Post struct is used to record information about each blog post
//Each Post should be named using the post date
type Post struct {
	Title string
	Views int
}

func main() {
	var addr = flag.String("addr", ":3000", "The port address of the application server")
	flag.Parse()

	fmt.Println("webserver started")

	http.Handle("/", newRouter())
	log.Println("Listening...")
	log.Fatal(http.ListenAndServe(*addr, nil))
}

func handlePost(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	mock := Post{"Test Post", 20}
	//need to check if params.ByName("post") is in the list of posts and if not 404

	viewLib.Counter.RLock()
	_, ok := viewLib.Counter.M[params.ByName("post")+".html"]
	viewLib.Counter.RUnlock()

	if ok == true {

		tmpl := template.Must(template.ParseFiles("tmpl/wrapper.html", "posts/"+params.ByName("post")+".html"))

		//ExecuteTemplate writes the template to w, writing "indexPage" as the main as defined
		//in index.html, and with a data interface
		err := tmpl.ExecuteTemplate(w, "wrapper", mock)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	} else {
		w.WriteHeader(404)
	}

}

func handleRefresh(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	refresh()
}

func refresh() error {
	dir, err := os.Getwd()

	files, err := ioutil.ReadDir(dir + "/posts")

	for _, f := range files {
		fmt.Println(f.Name())
		//check if in the view counter
		viewLib.Counter.RLock()
		_, ok := viewLib.Counter.M[f.Name()] //change to param
		viewLib.Counter.RUnlock()

		if ok == true {
			fmt.Println(f.Name() + " is already present")
		} else {
			viewLib.Counter.Lock()
			viewLib.Counter.M[f.Name()] = 0
			viewLib.Counter.Unlock()
			fmt.Println(f.Name() + " added to posts")
		}
	}
	return err

	//make the index page, need the date and post title to be shown in order
}

func testFiles(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	mock := Post{"Test Post", 20}

	tmpl := template.Must(template.ParseFiles("tmpl/wrapper.html", "tmpl/post.html"))

	//ExecuteTemplate writes the template to w, writing "indexPage" as the main as defined
	//in index.html, and with a data interface
	err := tmpl.ExecuteTemplate(w, "wrapper", mock)
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

//newRouter returns a httprouter.Router complete with the routes
func newRouter() *httprouter.Router {

	router := httprouter.New()
	router.GET("/Glob", testGlob)
	router.GET("/Files", testFiles)
	//password := input args for password
	//router.GET("/", handleIndex)
	//router.GET("/contact", handleContact)
	//router.GET("/projects", handleProjects)
	//router.GET("/posts", handlePostsIndex)
	router.GET("/posts/:post", handlePost)
	//router.GET("/stats", handleStats)
	router.GET("/refresh", handleRefresh)

	return router
}
