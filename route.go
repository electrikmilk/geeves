package geeves

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"text/template"
)

type HTTPMethod string

const (
	GET    HTTPMethod = "GET"
	POST              = "POST"
	PATCH             = "PATCH"
	PUT               = "PUT"
	DELETE            = "DELETE"
)

type route struct {
	name   string
	method HTTPMethod
	url    string
	file   string
}

type controllerFunc func(writer http.ResponseWriter, request *http.Request)

type controller struct {
	name     string
	method   HTTPMethod
	url      string
	callback controllerFunc
}

var Controllers []controller
var Routes []route

// Static creates a route when accessed at url by method will serve file
func Static(name string, url string, file string) {
	for _, route := range Routes {
		if route.name == name {
			Logf(BAD, "Fatal", "Failed to create route with name \"%s\", already exists", name)
			panic("unable to create route")
		}
	}
	for _, controller := range Controllers {
		if controller.name == name {
			Logf(BAD, "Fatal", "Failed to create route with name \"%s\", controller \"%s\" already exists", name, controller.name)
			panic("unable to create route")
		}
	}
	checkUrl(&name, &url)
	file = fmt.Sprintf("%s/%s", Dir, file)
	Routes = append(Routes, route{name: name, method: GET, url: url, file: file})
}

// Route creates route controller when accessed at url by method will call callback
func Route(name *string, method HTTPMethod, url *string, callback *controllerFunc) {
	checkUrl(*&name, &*url)
	for _, controller := range Controllers {
		if controller.name == *name {
			Logf(BAD, "Fatal", "Failed to create route controller with name \"%s\", already exists", name)
			panic("unable to create controller")
		}
	}
	for _, route := range Routes {
		if route.name == *name {
			Logf(BAD, "Fatal", "Failed to create route controller with name \"%s\", route \"%s\" already exists", name, route.name)
			panic("unable to create controller")
		}
	}
	Controllers = append(Controllers, controller{name: *name, method: method, url: *url, callback: *callback})
}

// Get is an alias for Route
func Get(name string, url string, callback controllerFunc) {
	Route(&name, GET, &url, &callback)
}

// Post is an alias for Route
func Post(name string, url string, callback controllerFunc) {
	Route(&name, POST, &url, &callback)
}

// Put is an alias for Route
func Put(name string, url string, callback controllerFunc) {
	Route(&name, PUT, &url, &callback)
}

// Patch is an alias for Route
func Patch(name string, url string, callback controllerFunc) {
	Route(&name, PATCH, &url, &callback)
}

// Delete is an alias for Route
func Delete(name string, url string, callback controllerFunc) {
	Route(&name, DELETE, &url, &callback)
}

// Redirect is a helper function to redirect to a new url or route by name
func Redirect(writer http.ResponseWriter, request *http.Request, newUrl string) {
	for _, route := range Routes {
		if route.name == newUrl {
			newUrl = route.url
		}
	}
	Logf(FYI, "Redirect", "%s -> %s", path, newUrl)
	http.Redirect(writer, request, newUrl, http.StatusSeeOther)
}

// RouteUrl returns the url of a route by name
func RouteUrl(name string) string {
	var url string
	for _, route := range Routes {
		if route.name == name {
			url = route.url
		}
	}
	if url == "" {
		Logf(BAD, "Error", "Route %s does not exist", name)
		panic("route does not exist")
	}
	return url
}

func checkUrl(route *string, url *string) {
	if !strings.HasPrefix(*url, "/") {
		Logf(FYI, "Notice", "Route \"%s\" URL should begin with a slash, added automatically.", *route, *url)
		*url = "/" + *url
	}
	if *url != "/" && strings.HasSuffix(*url, "/") {
		Logf(FYI, "Notice", "Route \"%s\" URL does not need a trailing slash: %s", *route, *url)
	}
	if strings.Contains(*url, ".") {
		Logf(BAD, "Fatal", "Route \"%s\" URL contains a \".\" character: %s", *route, *url)
		panic("unable to create route: invalid route url")
	}
}

// Template uses Go's template feature to replace the html in `file` with `data`.
// The output is written to `writer`.
func Template(name string, data interface{}, file string, writer http.ResponseWriter) {
	var staticFilePath = fmt.Sprintf("%s/%s", Dir, file)
	if _, err := os.Stat(staticFilePath); os.IsNotExist(err) {
		Logf(BAD, "Template", "File \"%s\" does not exist!", file)
		serverError(writer)
		return
	}
	var bytes, err = os.ReadFile(staticFilePath)
	if err != nil {
		Logf(BAD, "Template", "Unable to read file \"%s\"", file)
		serverError(writer)
		return
	}
	t, err := template.New(name).Parse(string(bytes))
	if err != nil {
		Logf(BAD, "Template", "Failed to parse template \"%s\"", file)
		serverError(writer)
		return
	}
	writer.Header().Set("Content-Type", "text/html")
	err = t.Execute(writer, data)
	if err != nil {
		Logf(BAD, "Template", "Failed to execute template \"%s\"", file)
		serverError(writer)
	}
}
