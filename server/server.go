package server

import (
	"github.com/gorilla/mux"
	"net/http"
)

// Server for the backend
type Server struct {
}

// New server
func New() *Server {
	return &Server{}
}

// just for testing, no error checking or anything
func (s *Server) sayHello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello world"))
}

// Start the server
func (s *Server) Start(port string) error {
	router := mux.NewRouter()
	router.Path("/hello").HandlerFunc(s.sayHello).Methods("GET")

	// shouldn't ever return
	return http.ListenAndServe(port, router)
}
