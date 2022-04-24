package gojsonserver

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type RequestHandler struct {
	path       string
	methods    []string
	jsonPath   string
	statusCode int
	delay      int
}

func GET(path string, jsonPath string, statusCode int, delay int) RequestHandler {
	return RequestHandler{
		path, []string{"GET"}, jsonPath, statusCode, delay,
	}
}

func POST(path string, jsonPath string, statusCode int, delay int) RequestHandler {
	return RequestHandler{
		path, []string{"POST"}, jsonPath, statusCode, delay,
	}
}

func PUT(path string, jsonPath string, statusCode int, delay int) RequestHandler {
	return RequestHandler{
		path, []string{"PUT"}, jsonPath, statusCode, delay,
	}
}

func PATCH(path string, jsonPath string, statusCode int, delay int) RequestHandler {
	return RequestHandler{
		path, []string{"PATCH"}, jsonPath, statusCode, delay,
	}
}

func DELETE(path string, jsonPath string, statusCode int, delay int) *RequestHandler {
	return &RequestHandler{
		path, []string{"DELETE"}, jsonPath, statusCode, delay,
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

func (js *JsonServer) start() {
	for _, requestHandler := range js.handlers {
		http.HandleFunc(requestHandler.path, func(res http.ResponseWriter, req *http.Request) {
			file, err := ioutil.ReadFile(requestHandler.jsonPath)
			if err != nil {
				log.Println(err)
			} else {
				time.Sleep(time.Duration(requestHandler.delay) * time.Millisecond)
				res.WriteHeader(requestHandler.statusCode)
				res.Write(file)
			}
		})
	}
	hostPort := js.config.host + ":" + fmt.Sprintf("%d", js.config.port)
	fmt.Printf("Running json server on %s \n", hostPort)
	http.ListenAndServe(hostPort, nil)
}
