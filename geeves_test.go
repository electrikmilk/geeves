package geeves

import (
	"net/http"
	"os"
	"testing"
	"time"
)

type TemplateData struct {
	Subject string
}

func TestServer(t *testing.T) {
	Dir = "test"
	Get("home", "/", func(writer http.ResponseWriter, request *http.Request) {
		Output(writer, "Hello, World!")
	})
	Static("test", "/test", "test.html")
	Get("template", "/template", func(writer http.ResponseWriter, request *http.Request) {
		var data = TemplateData{Subject: "Template"}
		Template("template", data, "template.html", writer)
	})
	Init(3000)
	time.Sleep(time.Second * 5)
	os.Exit(0)
}
