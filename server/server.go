package server

import (
	"google-cloud-task-processor/config"
	"google-cloud-task-processor/handlers/imageHandler"
	"google-cloud-task-processor/services/imageService"
	"net/http"
)
import "github.com/gorilla/mux"

type server struct {
	router *mux.Router
	imageHandler *imageHandler.ImageHandler
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) Start(port string) {
	http.ListenAndServe(port, s.router)
}

func (s *server) ConfigureRouter() {
	s.router.HandleFunc("/hello", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("hello world"))
	})

	s.router.HandleFunc("/storage/image", s.imageHandler.UploadImage).Methods("POST")
}

func NewServer(config *config.Config) *server {
	server := &server{
		router:       mux.NewRouter(),
		imageHandler: imageHandler.New(imageService.New(config)),
	}

	server.ConfigureRouter()
	return server
}