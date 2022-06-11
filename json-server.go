package gojsonserver

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type RequestHandler struct {
	path       string
	methods    []string
	jsonPath   func(*http.Request) string
	statusCode int
	delay      int
}

func Get(path string, jsonPath string, statusCode int, delay int) RequestHandler {
	return NewRequestHandler("GET", path, func(*http.Request) string { return jsonPath }, statusCode, delay)
}

func Post(path string, jsonPath string, statusCode int, delay int) RequestHandler {
	return NewRequestHandler("POST", path, func(*http.Request) string { return jsonPath }, statusCode, delay)
}

func Put(path string, jsonPath string, statusCode int, delay int) RequestHandler {
	return NewRequestHandler("PUT", path, func(*http.Request) string { return jsonPath }, statusCode, delay)
}

func Patch(path string, jsonPath string, statusCode int, delay int) RequestHandler {
	return NewRequestHandler("PATCH", path, func(*http.Request) string { return jsonPath }, statusCode, delay)
}

func Delete(path string, jsonPath string, statusCode int, delay int) RequestHandler {
	return NewRequestHandler("DELETE", path, func(*http.Request) string { return jsonPath }, statusCode, delay)
}

func GetProvider(path string, jsonPathProvider func(*http.Request) string, statusCode int, delay int) RequestHandler {
	return NewRequestHandler("GET", path, jsonPathProvider, statusCode, delay)
}

func PostProvider(path string, jsonPathProvider func(*http.Request) string, statusCode int, delay int) RequestHandler {
	return NewRequestHandler("POST", path, jsonPathProvider, statusCode, delay)
}

func PutProvider(path string, jsonPathProvider func(*http.Request) string, statusCode int, delay int) RequestHandler {
	return NewRequestHandler("PUT", path, jsonPathProvider, statusCode, delay)
}

func PatchProvider(path string, jsonPathProvider func(*http.Request) string, statusCode int, delay int) RequestHandler {
	return NewRequestHandler("PATCH", path, jsonPathProvider, statusCode, delay)
}

func DeleteProvider(path string, jsonPathProvider func(*http.Request) string, statusCode int, delay int) RequestHandler {
	return NewRequestHandler("DELETE", path, jsonPathProvider, statusCode, delay)
}

func NewRequestHandler(method string, path string, jsonProvider func(*http.Request) string, statusCode int, delay int) RequestHandler {
	return RequestHandler{
		path, []string{method}, jsonProvider, statusCode, delay,
	}
}

type ServerConfig struct {
	host string
	port int
}

type JsonServer struct {
	config   ServerConfig
	handlers []RequestHandler
}

func NewLocalJsonServer(port int, handlers []RequestHandler) *JsonServer {
	return NewJsonServer(
		ServerConfig{"localhost", port},
		handlers,
	)
}

func NewJsonServer(config ServerConfig, handlers []RequestHandler) *JsonServer {
	return &JsonServer{
		config,
		handlers,
	}
}

func (js *JsonServer) Start() {
	for _, requestHandler := range js.handlers {
		log.Printf("Registering handler %s %s ", strings.Join(requestHandler.methods, ", "), requestHandler.path)
		http.HandleFunc(requestHandler.path, func(res http.ResponseWriter, req *http.Request) {
			if !contains(requestHandler.methods, req.Method) {
				http.NotFound(res, req)
			} else {
				file, err := ioutil.ReadFile(requestHandler.jsonPath(req))
				if err != nil {
					log.Println(err)
				} else {
					time.Sleep(time.Duration(requestHandler.delay) * time.Millisecond)
					res.Header().Set("Content-Type", "application/json")
					res.Header().Set("Server", "GoJsonServer")
					res.WriteHeader(requestHandler.statusCode)
					res.Write(file)
				}
			}
		})
	}
	hostPort := js.config.host + ":" + fmt.Sprintf("%d", js.config.port)
	fmt.Printf("Running json server on %s \n", hostPort)
	http.ListenAndServe(hostPort, nil)
}

func contains(methods []string, method string) bool {
	for _, m := range methods {
		if m == method {
			return true
		}
	}
	return false
}
