package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/rizzza/echoserver/internal/config"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	cfg    config.Config
	router *chi.Mux
}

type Echo struct {
	Text string `json:"echo"`
}

func New(cfg config.Config) *Server {
	router := chi.NewRouter()
	srv := &Server{
		cfg:    cfg,
		router: router,
	}

	router.Get("/", srv.Get)
	router.Post("/echo", srv.Echo)

	return srv
}

func (s Server) GetRouter() *chi.Mux {
	return s.router
}

func (s Server) Get(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("hello there")); err != nil {
		log.Error("failed to write out error", err)
	}
}

func (s Server) Echo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")

	if r.Header.Get("content-type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		if _, err := w.Write([]byte(`{"error": "invalid content-type: must be application/json"}`)); err != nil {
			log.Error("failed to write out response", err)
		}
		return
	}

	reqBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Errorf("failed to read request body - %s: %v", string(reqBytes), err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !json.Valid(reqBytes) {
		w.WriteHeader(http.StatusBadRequest)
		if _, err := w.Write([]byte(`{"error": "invalid json"}`)); err != nil {
			log.Error("failed to write out response", err)
			return
		}
	}

	echo := &Echo{}
	if err := json.Unmarshal(reqBytes, echo); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Error("failed to unmarshal request", err)
		return
	}

	if len(echo.Text) <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		if _, err := w.Write([]byte(`{"error": "'echo' field is required"}`)); err != nil {
			log.Error("failed to write out response", err)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(fmt.Sprintf(`{"echo":"Echo server says -- %s"}`, echo.Text))); err != nil {
		log.Error("failed to write out response", err)
		return
	}
}

// Close the *server
func (s *Server) Close() {
	// close any relevent things
}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
