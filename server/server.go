package server

import (
	"encoding/json"
	"net/http"

	"github.com/SuganPrabu96/shorturl/redis"
	"github.com/gorilla/mux"
	"github.com/morikuni/failure"
)

// Server is a simple HTTP server which routes requests
type Server struct {
	r *redis.Client
}

// CreateEndpointObject is request body for create endpoint
type CreateEndpointObject struct {
	url string
}

// NewServer creates a new server
func NewServer(r *redis.Client) *Server {
	return &Server{r}
}

// CreateEndpoint creates a shortened url
func (s *Server) CreateEndpoint(w http.ResponseWriter, req *http.Request) {
	var c CreateEndpointObject
	err := json.NewDecoder(req.Body).Decode(&c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// TODO: Create the shortURL
	shortURL := "aaa"

	// TODO: Retry if shortURL already exists
	s.r.SetValueIfNotExists(shortURL, c.url)

	w.WriteHeader(http.StatusOK)
}

// RouteEndpoint routes endpoint
func (s *Server) RouteEndpoint(w http.ResponseWriter, req *http.Request) {
	shortURL := mux.Vars(req)["shortURL"]
	targetURL, err := s.r.GetValue(shortURL)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	}
	http.Redirect(w, req, targetURL, 301)
}

// Serve starts the server
func (s *Server) Serve() {
	router := mux.NewRouter()
	router.HandleFunc("/{shortURL}", s.RouteEndpoint).Methods("GET")
	router.HandleFunc("/create", s.CreateEndpoint).Methods("PUT")
	err := http.ListenAndServe(":5000", router)
	if err != nil {
		failure.Wrap(err)
	}
}
