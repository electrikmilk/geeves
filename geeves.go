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

var Dir = "static"
var Root = "/"
var Host = "localhost"
var Params url.Values
var path string

// Init starts to listen and serve at port
func Init(port int) {
	var p = fmt.Sprintf(":%d", port)
	Logf(FYI, "Routes", "%v", Routes)
	Logf(FYI, "Controllers", "%v", Controllers)
	Logf(GOOD, "Ready", "Waiting for requests at http://%s:%d.", Host, port)
	go http.HandleFunc(Root, handler)
	if err := http.ListenAndServe(p, nil); err != nil {
		Logf(BAD, "FATAL", fmt.Sprintf("Failed to listen and serve - %s", err))
		log.Fatal(err)
	}
}

func handler(writer http.ResponseWriter, request *http.Request) {
	// get request info
	var method = request.Method
	path = request.URL.Path[:]
	Logf(FYI, "Accepted", "%s %s", method, path)
	if err := request.ParseForm(); err != nil {
		Logf(BAD, "Fail", "ParseForm() error: %v", err)
		serverError(writer)
		return
	}
	Params = request.URL.Query()
	if len(Params) > 0 {
		Logf(FYI, "Params", "%v", Params)
	}
	// If URL has static directory prefix
	var staticDir = fmt.Sprintf("/%s/", Dir)
	if strings.HasPrefix(path, staticDir) {
		var splitPath = strings.Split(path, "/")
		var staticFilePath = fmt.Sprintf("%s/%s", Dir, splitPath[2])
		_, err := os.OpenFile(staticFilePath, os.O_RDONLY, 0775)
		if errors.Is(err, os.ErrNotExist) {
			Logf(WARN, "404", "%s %s", method, path)
			notFound(writer)
			Logf(FYI, "Closing", "%s %s", method, path)
			return
		}
		// Make sure a route doesn't own this file
		for _, route := range Routes {
			if staticFilePath == route.file {
				Logf(WARN, "404", "%s %s - \"%s\" belongs to route \"%s\", the correct link is %s", method, path, splitPath[2], route.name, route.url)
				notFound(writer)
				Logf(FYI, "Closing", "%s %s", method, path)
				return
			}
		}
		// Make sure this file is not an HTML file
		var ext = filepath.Ext(staticFilePath)
		if ext == ".html" {
			Logf(WARN, "404", "%s %s - \"%s\" is an HTML file, please create a route", method, path, splitPath[2])
			notFound(writer)
			Logf(FYI, "Closing", "%s %s", method, path)
			return
		}
		Logf(GOOD, "Static", "%s -> %s", staticFilePath, splitPath[2])
		http.ServeFile(writer, request, staticFilePath)
		Logf(FYI, "Closing", "%s %s", method, path)
		return
	}
	// look for matching route
	for _, route := range Routes {
		if string(route.method) == method && route.url == path {
			Logf(GOOD, "Route", "%s %s - serving file %s", method, route.url, route.file)
			http.ServeFile(writer, request, route.file)
			Logf(FYI, "Closing", "%s %s", method, path)
			return
		}
		if string(route.method) != method && route.url == path {
			Logf(FYI, "Notice", "Route %s has method %s, not %s", route.url, route.method, method)
		}
	}
	// Look for matching controller
	for _, controller := range Controllers {
		if string(controller.method) == method && controller.url == path {
			// Pass request to controller
			Logf(GOOD, "Controller", "%s %s - %s", method, controller.url, controller.name)
			// fmt.Printf("Path Params: %v\n", path_params)
			controller.callback(writer, request)
			Logf(FYI, "Closing", "%s %s", method, path)
			return
		}
		if string(controller.method) != method && controller.url == path {
			Logf(WARN, "Denied", "Controller %s has method %s, not %s", controller.url, controller.method, method)
		}
	}
	// Failsafe
	Logf(WARN, "Failed", "404: %s", path)
	notFound(writer)
	Logf(FYI, "Closing", "%s %s", method, path)
}
