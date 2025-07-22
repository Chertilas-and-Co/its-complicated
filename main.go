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

func hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.WriteHeader(200)
	fmt.Fprintf(w, "Hello, %s!", ps.ByName("name"))
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}

func getForm(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	http.ServeFile(w, r, "form.html")
}

func registerHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Form parse error", http.StatusBadRequest)
		return
	}

	password := r.FormValue("password")
	passwordConfirm := r.FormValue("passwordConfirm")
	fmt.Println(password, " ", passwordConfirm)

	if password != passwordConfirm {
		http.Error(w, "passwords do not match", http.StatusBadRequest)
		return
	}

	w.Write([]byte("Registration success!"))
}

func main() {
	router := httprouter.New()
	router.ServeFiles("/static/*filepath", http.Dir("static"))

	router.GET("/register", handleLogger(getForm))
	router.NotFound = http.HandlerFunc(handlerFuncLogger(notFoundHandler))
	router.GET("/hello/:name", handleLogger(hello))
	router.POST("/register", handleLogger(registerHandler))

	fmt.Println("Server is listening at :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
