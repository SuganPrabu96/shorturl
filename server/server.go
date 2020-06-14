package server

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/SuganPrabu96/shorturl/redis"
	"github.com/gorilla/mux"
	"github.com/morikuni/failure"
)

const (
	shortURLLength = 10
)

// Server is a simple HTTP server which routes requests
type Server struct {
	Address string
	Client  *redis.Client
}

// CreateEndpointObject is request body for create endpoint
type CreateEndpointObject struct {
	URL string
}

// GetRandomString creates a random string of given length
func GetRandomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" + "abcdefghijklmnopqrstuvwxyz" +
		"0123456789")
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	return b.String()
}

// NewServer creates a new server
func NewServer(host string, port int, client *redis.Client) *Server {
	address := host + ":" + strconv.Itoa(port)
	return &Server{address, client}
}

// CreateEndpoint creates a shortened url
func (s *Server) CreateEndpoint(w http.ResponseWriter, req *http.Request) {
	var c CreateEndpointObject
	err := json.NewDecoder(req.Body).Decode(&c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	shortURL := GetRandomString(shortURLLength)

	// If URL doesn't start with http:// or https:// attach a prefix http://
	// This is needed for routing endpoints
	if !(strings.HasPrefix(c.URL, "http://") || strings.HasPrefix(c.URL, "https://")) {
		c.URL = "http://" + c.URL
	}
	err = s.Client.SetValueIfNotExists(shortURL, c.URL)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Println("Shortened URL ", shortURL)
}

// RouteEndpoint routes endpoint
func (s *Server) RouteEndpoint(w http.ResponseWriter, req *http.Request) {
	shortURL := mux.Vars(req)["shortURL"]
	targetURL, err := s.Client.GetValue(shortURL)
	if err != nil || targetURL == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	fmt.Println("Target URL : ", targetURL)
	http.Redirect(w, req, targetURL, 301)
}

// Serve starts the server
func (s *Server) Serve() {
	router := mux.NewRouter()
	router.HandleFunc("/{shortURL}", s.RouteEndpoint).Methods("GET")
	router.HandleFunc("/create", s.CreateEndpoint).Methods("PUT")
	err := http.ListenAndServe(s.Address, router)
	if err != nil {
		failure.Wrap(err)
	}
}
