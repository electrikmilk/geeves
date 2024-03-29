package geeves

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"time"
)

type Pigment string
type ServerLog int
type Encoding string

const (
	RESET  Pigment = "\033[0m"
	RED            = "\033[31m"
	GREEN          = "\033[32m"
	YELLOW         = "\033[33m"
	BLUE           = "\033[34m"
	PURPLE         = "\033[35m"
	CYAN           = "\033[36m"
	GRAY           = "\033[37m"
	GREY           = "\033[37m"
	WHITE          = "\033[97m"
)

const (
	FYI  ServerLog = 0
	GOOD           = 1
	WARN           = 2
	BAD            = 3
)

const (
	HTMLEncoding Encoding = "text/html"
	XMLEncoding           = "application/xml"
	JSONEncoding          = "application/json"
)

// Color returns the string in color
func Color(str string, color Pigment) string {
	if runtime.GOOS == "windows" {
		return str
	}
	return fmt.Sprintf("%s%s%s", color, str, RESET)
}

// Colour is an alias for Color
func Colour(str string, color Pigment) string {
	return Color(str, color)
}

// Log outputs the string with the label and in the color that corresponds to the type
func Log(logType ServerLog, label string, str string) {
	generateLog(&logType, &label, &str)
	fmt.Printf(str)
}

// Logf outputs the resulting string with the label and in the color that corresponds to the type
func Logf(logType ServerLog, label string, str string, vars ...interface{}) {
	generateLog(&logType, &label, &str)
	fmt.Printf(str, vars...)
}

func generateLog(logtype *ServerLog, label *string, str *string) {
	*str = fmt.Sprintf("[%s] (%s) %s\n", time.Now(), *label, *str)
	switch *logtype {
	case 0:
		*str = Color(*str, BLUE)
	case 1:
		*str = Color(*str, GREEN)
	case 2:
		*str = Color(*str, YELLOW)
	case 3:
		*str = Color(*str, RED)
	}
}

// Output is an alias function to output the string in HTML encoding
func Output(writer http.ResponseWriter, content string) {
	Eprint(writer, HTMLEncoding, content)
}

// Outputf is an alias function to output the resulting string in HTML encoding
func Outputf(writer http.ResponseWriter, content string, vars ...interface{}) {
	Eprintf(writer, HTMLEncoding, content, vars...)
}

// Jprint is an alias function to output the content in JSON encoding
func Jprint(writer http.ResponseWriter, content any) {
	response, err := json.Marshal(content)
	if err != nil {
		Logf(BAD, "Fatal", "JSON failed to encode: %v", err)
	}
	Eprint(writer, JSONEncoding, string(response))
}

// Eprint encodes the string in encoding and writes it to writer
func Eprint(writer http.ResponseWriter, encoding Encoding, content string) {
	var contentType = fmt.Sprintf("%s; charset=utf-8", encoding)
	writer.Header().Set("Content-Type", contentType)
	fmt.Fprint(writer, content)
}

// Eprintf encodes the resulting string in encoding and writes it to writer
func Eprintf(writer http.ResponseWriter, encoding Encoding, content string, vars ...interface{}) {
	var contentType = fmt.Sprintf("%s; charset=utf-8", encoding)
	content = fmt.Sprintf(content, vars...)
	writer.Header().Set("Content-Type", contentType)
	fmt.Fprint(writer, content)
}

func notFound(writer http.ResponseWriter) {
	http.Error(writer, "404 Not Found.", http.StatusNotFound)
}

func serverError(writer http.ResponseWriter) {
	http.Error(writer, "500 Internal Server Error.", http.StatusInternalServerError)
}
