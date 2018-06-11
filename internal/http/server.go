package http

import (
  "bufio"
  "fmt"
  "github.com/gorilla/mux"
  "github.com/kleijnweb/git-mirror-manager/internal/manager"
  "github.com/kleijnweb/git-mirror-manager/internal/util"
  log "github.com/sirupsen/logrus"
  "net/http"
  "time"
)

// Server handles request/responses and delegates to the manager
type Server struct {
  manager *manager.Manager
  addr    string
}

// NewServer creates a new Server
func NewServer(manager *manager.Manager) *Server {
  return &Server{manager: manager}
}

// Start initializes the server and makes it listen for connections
func (s *Server) Start() {
  srv, err := s.configure(manager.NewConfig())

  if err != nil {
    s.handleStartupError(err)
  }

  if err := srv.ListenAndServe(); err != nil {
    s.handleStartupError(&util.ApplicationError{err, util.ErrNet})
  }
}

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.RequestURI)
    next.ServeHTTP(w, r)
  })
}

func (s *Server) configure(config *manager.Config) (*http.Server, *util.ApplicationError) {
  router := mux.NewRouter()
  router.HandleFunc("/ping", s.ping).Methods("GET")
  router.HandleFunc("/repo", s.createMirror).Methods("POST")
  router.HandleFunc("/repo/{namespace}/{name}", s.deleteMirror).Methods("DELETE")
  router.Use(s.loggingMiddleware)

  srv := &http.Server{
    Handler:      router,
    Addr:         config.ManagerAddr,
    WriteTimeout: 15 * time.Second,
    ReadTimeout:  15 * time.Second,
  }

  log.Println("Listening on " + config.ManagerAddr)

  return srv, nil
}

func (s *Server) ping(w http.ResponseWriter, r *http.Request) {
  w.WriteHeader(http.StatusOK)
  fmt.Fprintf(w, "pong\n")
}

func (s *Server) createMirror(w http.ResponseWriter, r *http.Request) {
  scanner := bufio.NewScanner(r.Body)
  for scanner.Scan() {
    if err := s.manager.Add(scanner.Text()); err != nil {
      s.handleServingError(w, err)
    }
  }

  if err := scanner.Err(); err != nil {
    s.handleServingError(w, util.NewError("failed reading request body", util.ErrUser))
  }
}

func (s *Server) deleteMirror(w http.ResponseWriter, r *http.Request) {
  name := mux.Vars(r)["namespace"] + "/" + mux.Vars(r)["name"]
  if err := s.manager.Remove(name); err != nil {
    s.handleServingError(w, err)
  }
}

func (s *Server) handleServingError(w http.ResponseWriter, err *util.ApplicationError) {
  if err.Code == util.ErrUser {
    w.WriteHeader(http.StatusBadRequest)
  } else if err.Code == util.ErrNotFound {
    w.WriteHeader(http.StatusNotFound)
  } else {
    w.WriteHeader(http.StatusInternalServerError)
  }

  log.Print(err)
}

func (s *Server) handleStartupError(err *util.ApplicationError) {
  log.Fatal(err)
}
