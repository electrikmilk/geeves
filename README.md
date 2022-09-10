# Geeves

A simple Go web framework.

My main focus was to allow the developer to write less boilerplate.

### Features

- [x] Serve static files
- [x] Create routes at a specific url and request method and serve a file in response.
- [x] Create a route controller and process the request and write a response by serving a file, redirecting or directly outputting something.
- [x] Helper functions for outputting in web specific encoding formats (HTML, JSON, XML).
- [x] Helper functions for standard logging (info, warning, success, error) in a similar format and corresponding colors. You can also format a string to be outputted in one of the defined colors alternatively.

### Missing

- [ ] Database abstraction (Model/ORM)
- [ ] View templates

## Usage

```go
import "github.com/electrikmilk/geeves"

func main() {
  geeves.Init("3000")
}
``` 

### Defaults

Set static content directory:
```go
geeves.Dir = "public" // default: "static"
```

Set path to start at:
```go
geeves.Root = "/my-geeves-site/" // default: "/"
```

### Routes

Create a route to serve a static file:
```go
geeves.Static("home","/","index.html") // defaults to GET method and serves file
```

Create a route controller and serve a processed response:
```go
// Get(), Post(), Patch(), Put(), Delete()
geeves.Get("api", "/api", func(writer http.ResponseWriter, request *http.Request) {
	var response map[string]any
	response = make(map[string]any)
	response["status"] = "success"
	response["message"] = "Hello from API!"
	geeves.Jprint(writer, response)
})
```

Redirect to a custom URL or by a route name:

```go
geeves.Redirect(writer, request, "/url")
geeves.Redirect(writer, request, "home") // redirects to url for "home" route
```

Get the URL of a route:
```go
geeves.RouteUrl("home") // returns "/"
```

### Output

Output in HTML encoding:
```go
geeves.Output(writer, "Output...")
geeves.Outputf(writer, "%s...", variable)
```

Encode and output any type of variable in JSON:
```go
var response map[string]any
response = make(map[string]any)
response["message"] = "Hello"
geeves.Jprint(writer, response)

response := [...]string{"one", "two", "three"}
geeves.Jprint(writer, response)
```

Output directly in an encoding:
```go
geeves.Eprint(writer, geeves.HTMLEncoding|JSONEncoding|XMLEncoding, "content")
geeves.Eprintf(writer, geeves.HTMLEncoding|JSONEncoding|XMLEncoding, "%s...", variable)
```

Log helper functions:
```go
geeves.Log(geeves.FYI|WARN|GOOD|BAD, "Label", "My log message")
geeves.Logf(geeves.FYI|WARN|GOOD|BAD, "Label", "My log message: %s", variable)
// FYI = BLUE, WARN = YELLOW, GOOD = GREEN, BAD = RED

fmt.Println(Color("string", geeves.Pigment))
// Pigments: RED, GREEN, YELLOW, BLUE, PURPLE, CYAN, GRAY|GREY, WHITE
// alias: Colour()
```

**Note:** `BAD` type log calls `Fatal()` or `Fatalf()` which kills the server.
