package gojsonserver

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	. "github.com/gorilla/websocket"
)

type RequestHandler struct {
	path       string
	methods    []string
	jsonPath   func(*http.Request) string
	statusCode int
	delay      int
}

type WebSocketHandler struct {
	path    string
	handler func(*Conn)
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

func WebSocketProvider(path string, handler func(con *Conn)) WebSocketHandler {
	return WebSocketHandler{path, handler}
}

type ServerConfig struct {
	Host string
	Port int
}

type JsonServer struct {
	config            ServerConfig
	handlers          []RequestHandler
	webSocketHandlers []WebSocketHandler
}

func NewLocalJsonServer(port int, handlers []RequestHandler) *JsonServer {
	return NewJsonServer(
		ServerConfig{"localhost", port},
		handlers,
		nil,
	)
}

func NewJsonServer(config ServerConfig, handlers []RequestHandler, webSocketHandlers []WebSocketHandler) *JsonServer {
	return &JsonServer{
		config,
		handlers,
		webSocketHandlers,
	}
}

func (js *JsonServer) Start() {
	for _, requestHandler := range js.handlers {
		log.Printf("Registering handler %s %s ", strings.Join(requestHandler.methods, ", "), requestHandler.path)
		func(requestHandler RequestHandler) {
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
		}(requestHandler)

		if js.webSocketHandlers != nil {
			for _, webSocketHandler := range js.webSocketHandlers {
				func(webSocketHandler WebSocketHandler) {
					http.HandleFunc(webSocketHandler.path, func(res http.ResponseWriter, req *http.Request) {
						var upgrader = Upgrader{}
						// allow any origin
						// should not be used in production code
						upgrader.CheckOrigin = func(r *http.Request) bool { return true }
						conn, err := upgrader.Upgrade(res, req, nil)
						if err != nil {
							log.Print("upgrade failed: ", err)
							return
						}
						defer conn.Close()
						webSocketHandler.handler(conn)
					})

				}(webSocketHandler)
			}
		}
	}
	hostPort := js.config.Host + ":" + fmt.Sprintf("%d", js.config.Port)
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
