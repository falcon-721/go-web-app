package main

import (
	"bufio"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"

	"github.com/gorilla/mux"
)

func main() {
	serveWeb()
}

//struc to pass into the templage
type defaultContext struct {
	Title       string
	ErrorMsgs   string
	SuccessMsgs string
}

var themeName = getThemeName()          //late we will collect this value from config file
var staticPages = populateStaticPages() //custom function to colect all the available pages under pages folder

func serveWeb() {

	gorillaRoute := mux.NewRouter()

	gorillaRoute.HandleFunc("/", serveContent)
	gorillaRoute.HandleFunc("/{page_alias}", serveContent) //Dynamic url value

	http.HandleFunc("/img/", serveResource)
	http.HandleFunc("/css/", serveResource)
	http.HandleFunc("/js/", serveResource)

	http.Handle("/", gorillaRoute)
	http.ListenAndServe(":8080", nil)
}
func serveContent(w http.ResponseWriter, r *http.Request) {
	urlParams := mux.Vars(r)
	page_alias := urlParams["page_alias"]
	if page_alias == "" {
		page_alias = "home"
	}

	staticPage := staticPages.Lookup(page_alias + ".html")
	log.Println("page ", staticPage)

	if staticPage == nil {
		staticPage = staticPages.Lookup("404.html")
		w.WriteHeader(404)
	}

	//Value to pass into template
	context := defaultContext{}
	context.Title = page_alias
	context.ErrorMsgs = ""
	context.SuccessMsgs = ""
	log.Println(context)
	err := staticPage.Execute(w, context)
	if err != nil {
		log.Println(err)
	}
}
func getThemeName() string {
	return "bs4"
}

//--------------------------------------------------------------
// Retrieve all fils under the given folder and its subsequent folder
func populateStaticPages() *template.Template {
	result := template.New("templates")
	templatePaths := new([]string)

	basePath := "pages"
	templateFolder, _ := os.Open(basePath)
	defer templateFolder.Close()
	templatePathsRaw, _ := templateFolder.Readdir(-1)
	for _, pathInfo := range templatePathsRaw {
		log.Println(pathInfo.Name())
		*templatePaths = append(*templatePaths, basePath+"/"+pathInfo.Name())
	}

	basePath = "themes/" + themeName
	templateFolder, _ = os.Open(basePath)
	defer templateFolder.Close()
	templatePathsRaw, _ = templateFolder.Readdir(-1)
	for _, pathInfo := range templatePathsRaw {
		log.Println(pathInfo.Name())
		*templatePaths = append(*templatePaths, basePath+"/"+pathInfo.Name())
	}
	result.ParseFiles(*templatePaths...)
	return result
}

//--------------------------------------------------------------

//--------------------------------------------------------------
//Serve Resources of types js, img, css files
func serveResource(w http.ResponseWriter, req *http.Request) {
	path := "public/" + themeName + req.URL.Path
	var contentType string
	if strings.HasSuffix(path, ".css") {
		contentType = "text/css; charset=utf-8"
	} else if strings.HasSuffix(path, ".png") {
		contentType = "img/png; charset=utf-8"
	} else if strings.HasSuffix(path, ".jpg") {
		contentType = "img/jpg; charset=utf-8"
	} else if strings.HasSuffix(path, ".js") {
		contentType = "application/javascript; charset=utf-8"
	} else {
		contentType = "text/plain; charset=utf-8"
	}

	log.Println(path)
	f, err := os.Open(path)
	if err == nil {
		defer f.Close()
		w.Header().Add("Content-Type", contentType)
		br := bufio.NewReader(f)
		br.WriteTo(w)
	} else {
		w.WriteHeader(404)
	}
}

//--------------------------------------------------------------
