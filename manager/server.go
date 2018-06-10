package manager

import (
	"bufio"
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

// NewServer creates a new Server
func NewServer() *Server {
	return &Server{}
}

// Server handles request/responses and delegates to the manager
type Server struct {
	manager *manager
	addr    string
}

// Start initialized the server and makes it listen for connections
func (s *Server) Start() {
	srv, err := s.configure(newConfig())

	if err != nil {
		s.handleStartupError(err)
	}

	if err := srv.ListenAndServe(); err != nil {
		s.handleStartupError(&Error{err, errNet})
	}
}

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func (s *Server) configure(config *config) (*http.Server, *Error) {
	router := mux.NewRouter()
	router.HandleFunc("/ping", s.ping).Methods("GET")
	router.HandleFunc("/repo", s.createMirror).Methods("POST")
	router.HandleFunc("/repo/{namespace}/{name}", s.deleteMirror).Methods("DELETE")
	router.Use(s.loggingMiddleware)

	srv := &http.Server{
		Handler:      router,
		Addr:         config.managerAddr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("Listening on " + config.managerAddr)

	s.manager = &manager{}
	if err := s.manager.configure(config); err != nil {
		return nil, err
	}
	return srv, nil
}

func (s *Server) ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "pong\n")
}

func (s *Server) createMirror(w http.ResponseWriter, r *http.Request) {
	scanner := bufio.NewScanner(r.Body)
	for scanner.Scan() {
		if err := s.manager.add(scanner.Text()); err != nil {
			s.handleServingError(w, err)
		}
	}

	if err := scanner.Err(); err != nil {
		s.handleServingError(w, newError("failed reading request body", errUser))
	}
}

func (s *Server) deleteMirror(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["namespace"] + "/" + mux.Vars(r)["name"]
	if err := s.manager.remove(name); err != nil {
		s.handleServingError(w, err)
	}
}

func (s *Server) handleServingError(w http.ResponseWriter, err *Error) {
	if err.code == errUser {
		w.WriteHeader(http.StatusBadRequest)
	} else if err.code == errNotFound {
		w.WriteHeader(http.StatusNotFound)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}

	log.Print(err)
}

func (s *Server) handleStartupError(err *Error) {
	log.Fatal(err)
}
