package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

// структура-обёртка для ResponseWriter — перехватывает WriteHeader
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

// middleware для логирования
func handleLogger(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		start := time.Now()
		sr := &statusRecorder{ResponseWriter: w, status: 200}
		h(sr, r, ps)
		duration := time.Since(start)
		log.Printf("%s %s %d %s", r.Method, r.URL.Path, sr.status, duration)
	}
}

func handlerFuncLogger(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		sr := &statusRecorder{ResponseWriter: w, status: 200}
		h(sr, r)
		duration := time.Since(start)
		log.Printf("%s %s %d %s", r.Method, r.URL.Path, sr.status, duration)
	}
}

// Обычный обработчик
func hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.WriteHeader(200)
	fmt.Fprintf(w, "Hello, %s!", ps.ByName("name"))
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}

func main() {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(handlerFuncLogger(notFoundHandler))
	router.GET("/hello/:name", handleLogger(hello))
	log.Fatal(http.ListenAndServe(":8080", router))
}
