package geeves

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

var Dir string = "static"
var Root string = "/"
var Params url.Values
var path string

// Init starts to listen and serve at port
func Init(port string) {
	port = fmt.Sprintf(":%s", port)
	Logf(FYI, "Routes", "%v", Routes)
	Logf(FYI, "Controllers", "%v", Controllers)
	Logf(GOOD, "Ready", "Waiting for requests at http://localhost%s.", port)
	go http.HandleFunc(Root, handler)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}

func handler(writer http.ResponseWriter, request *http.Request) {
	// get request info
	var method string = request.Method
	path = request.URL.Path[:]
	Logf(FYI, "Request", "%s -> %s", method, path)
	if err := request.ParseForm(); err != nil {
		Logf(BAD, "Fail", "ParseForm() error: %v", err)
		return
	}
	Params = request.URL.Query()
	if len(Params) > 0 {
		fmt.Printf("Params: %v\n", Params)
	}
	// If URL has static directory prefix
	var static_dir string = fmt.Sprintf("/%s/", Dir)
	if strings.HasPrefix(path, static_dir) {
		var split_path []string = strings.Split(path, "/")
		var static_file_path string = fmt.Sprintf("%s/%s", Dir, split_path[2])
		_, err := os.OpenFile(static_file_path, os.O_RDONLY, 0775)
		if errors.Is(err, os.ErrNotExist) {
			Logf(WARN, "Failed", "404: %s", path)
			http.Error(writer, "404 not found.", http.StatusNotFound)
			return
		}
		// Make sure a route doesn't own this file
		for _, route := range Routes {
			if static_file_path == route.file {
				Logf(WARN, "Denied", "%s belongs to route \"%s\", the correct link is %s", split_path[2], route.name, route.url)
				http.Error(writer, "404 not found.", http.StatusNotFound)
				return
			}
		}
		// Make sure this file is not an HTML file
		var ext string = filepath.Ext(static_file_path)
		if ext == ".html" {
			Logf(WARN, "Denied", "%s is an HTML file, please create a route", split_path[2])
			http.Error(writer, "404 not found.", http.StatusNotFound)
			return
		}
		Logf(GOOD, "Static", "%s -> %s", static_file_path, split_path[2])
		http.ServeFile(writer, request, static_file_path)
		return
	}
	// look for matching route
	for _, route := range Routes {
		if string(route.method) == method && route.url == path {
			Logf(GOOD, "Route", "%s -> %s", route.url, route.file)
			http.ServeFile(writer, request, route.file)
			return
		}
		if string(route.method) != method && route.url == path {
			Logf(FYI, "Notice", "Route %s has method %s, not %s", route.url, route.method, method)
		}
	}
	// Look for matching controller
	for _, contrl := range Controllers {
		if string(contrl.method) == method && contrl.url == path {
			// Pass request to controller
			Logf(GOOD, "Controller", "%s -> %s", contrl.url, contrl.name)
			// fmt.Printf("Path Params: %v\n", path_params)
			contrl.callback(writer, request)
			return
		}
		if string(contrl.method) != method && contrl.url == path {
			Logf(WARN, "Denied", "Controller %s has method %s, not %s", contrl.url, contrl.method, method)
		}
	}
	// Failsafe
	Logf(WARN, "Failed", "404: %s", path)
	http.Error(writer, "404 not found.", http.StatusNotFound)
}

func checkUrl(route string, url string) string {
	if !strings.HasPrefix(url, "/") {
		Logf(FYI, "Notice", "Route \"%s\" URL should begin with a slash, added automatically.", route, url)
		url = "/" + url
	}
	if url != "/" && strings.HasSuffix(url, "/") {
		Logf(FYI, "Notice", "Route \"%s\" URL does not need a trailing slash: %s", route, url)
	}
	if strings.Contains(url, ".") {
		Logf(BAD, "Fatal", "Route \"%s\" URL contains a \".\" character: %s", route, url)
	}
	return url
}
