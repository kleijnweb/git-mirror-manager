package manager

import (
	"net/http"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"github.com/gorilla/mux"
	"time"
)

func NewManagerServer() *managerServer {
	return &managerServer{}
}

type managerServer struct {
	manager *manager
	addr    string
}

func (s *managerServer) Start() {
	srv, err := s.configure(newConfig())

	if err != nil {
		s.handleStartupError(err)
	}

	if err := srv.ListenAndServe(); err != nil {
		s.handleStartupError(&Error{err, errNet})
	}
}

func (s *managerServer) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func (s *managerServer) configure(config *config) (*http.Server, *Error) {
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

func (s *managerServer) ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "pong\n")
}

func (s *managerServer) createMirror(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		s.handleServingError(w, &Error{err, errUser})
		return
	}
	if err := s.manager.add(string(body)); err != nil {
		s.handleServingError(w, err)
	}
}

func (s *managerServer) deleteMirror(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["namespace"] + "/" + mux.Vars(r)["name"]
	if err := s.manager.remove(name); err != nil {
		s.handleServingError(w, err)
	}
}

func (s *managerServer) handleServingError(w http.ResponseWriter, err *Error) {
	if err.code == errUser {
		w.WriteHeader(http.StatusBadRequest)
	} else if err.code == errNotFound {
		w.WriteHeader(http.StatusNotFound)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}

	log.Print(err)
}

func (s *managerServer) handleStartupError(err *Error) {
	log.Fatal(err)
}
