package main

import (
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/Charlesworth/viewsLib"
	"github.com/julienschmidt/httprouter"
)

//Post struct is used to record information about each blog post
//Each Post should be named using the post date
// type Post struct {
// 	Title string
// 	Views int
// 	Date  int
// }

func main() {
	var addr = flag.String("addr", ":3000", "The port address of the application server")
	flag.Parse()

	//test code
	// posts := Posts{
	// 	{"United States", 0, 10},
	// 	{"Bahamas", 0, 51},
	// 	{"Japan", 0, 12},
	// }
	//
	// sort.Sort(posts)
	//
	// for _, c := range posts {
	// 	fmt.Println(c.Date, c.Name)
	// }
	//test code

	refreshPosts()

	fmt.Println("webserver started")

	http.Handle("/", newRouter())
	log.Println("Listening...")
	log.Fatal(http.ListenAndServe(*addr, nil))
}

func handlePost(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	mock := Post{"Test Post", 20, 1}

	//need to check if params.ByName("post") is in the list of posts and if not 404
	ok := viewLib.PageExists(params.ByName("post") + ".html")

	if ok == true {

		tmpl := template.Must(template.ParseFiles("tmpl/wrapper.html", "posts/"+params.ByName("post")+".html"))

		viewLib.ViewInc(r.RemoteAddr, params.ByName("post")+".html")

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
	refreshPosts()
}

func refreshPosts() {
	dir, err := os.Getwd()

	postsDir, err := ioutil.ReadDir(dir + "/posts")

	var posts Posts

	for _, post := range postsDir {
		//check if the post is in the view counter
		viewLib.Counter.RLock()
		views, postExists := viewLib.Counter.M[post.Name()]
		viewLib.Counter.RUnlock()

		if postExists == true {
			//fmt.Println(post.Name() + " is present")
		} else {
			viewLib.Counter.Lock()
			viewLib.Counter.M[post.Name()] = 0
			viewLib.Counter.Unlock()
			//fmt.Println(post.Name() + " added to posts")
		}

		//fmt.Println(post.Name())
		//add to a list of posts
		name := strings.Split(post.Name(), "-")
		fmt.Print(name[0], "/", name[1], "/", name[2], " ")
		l := len(name) - 1
		var fullName string
		for i := 3; i <= l; i++ {
			if i == l {
				fullName = fullName + strings.Split(name[i], ".")[0]
			} else {
				fullName = fullName + name[i] + " "
			}
		}

		order, _ := strconv.Atoi(name[2] + name[1] + name[0])
		fmt.Println(order)
		a := Post{fullName, views, order} //get the views from line 80
		posts = append(posts, a)
	}

	//make that list into the index page
	sort.Sort(posts)

	for i := range posts {
		fmt.Print(posts[i].Name, ", Date:", posts[i].Date, " Views:", posts[i].Views, "\n")
	}

	if err != nil {
		fmt.Println(err)
	}
}

func testFiles(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	mock := Post{"Test Post", 20, 1}

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
