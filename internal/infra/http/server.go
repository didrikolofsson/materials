// Package http provides HTTP server functionality.
package http

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/didrikolofsson/materials/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	httpServer *http.Server
}

func New(port string, handlers *handlers.Handlers) *Server {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Routes (Possibly admin routes)
	r.Get("/teachers", handlers.ListTeachers)
	r.Get("/subjects", handlers.ListSubjects)

	// Routes (Possibly public)
	r.Get("/materials", handlers.ListMaterials)
	r.Get("/materials/{id}/versions", handlers.ListMaterialVersionsByMaterialID)
	r.Put("/materials/{id}/versions/{version_id}/main", handlers.UpdateMaterialVersionMain)

	// Teacher routes
	r.Get("/teachers/{id}", handlers.GetTeacherByID)
	r.Get("/teachers/{id}/materials", handlers.GetTeacherMaterials)
	r.Post("/teachers/{id}/materials", handlers.CreateInitialTeacherMaterial)
	r.Get("/teachers/{id}/materials/{material_id}", handlers.GetTeacherMaterialByID)
	r.Put("/teachers/{id}/materials/{material_id}", handlers.UpdateTeacherMaterialByID)
	r.Delete("/teachers/{id}/materials/{material_id}", handlers.DeleteTeacherMaterialByID)

	// Healthcheck
	r.Get("/ping", handlers.Ping)

	addr := fmt.Sprintf(":%s", port)
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	return &Server{httpServer: srv}
}

func (s *Server) Run() error {
	// Listen for termination signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("Server is running on port %s", s.httpServer.Addr)
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on %s: %v", s.httpServer.Addr, err)
		}
	}()

	<-stop
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return s.httpServer.Shutdown(ctx)
}
