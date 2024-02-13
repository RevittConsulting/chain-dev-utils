package server

import (
	"context"
	"fmt"
	"github.com/RevittConsulting/cdk-envs/config"
	"github.com/RevittConsulting/cdk-envs/pkg/atomics"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type Server struct {
	Config       *config.Config
	ShuttingDown *atomics.AtomicBool
	Router       *chi.Mux
	Deps         *dependencies
}

func NewServer(sd *atomics.AtomicBool, r *chi.Mux) *Server {
	return &Server{
		ShuttingDown: sd,
		Router:       r,
	}
}

type StartFunc func(ctx context.Context, s *Server) error

func Start(startFunc StartFunc) {
	ctx, cancel := context.WithCancel(context.Background())

	var shutdown atomics.AtomicBool
	shutdown.Set(false)

	r := chi.NewRouter()

	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"*"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	r.Use(cors.Handler)

	s := NewServer(&shutdown, r)

	err := startFunc(ctx, s)
	if err != nil {
		log.Fatal("failed to start server")
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		cancel()
	}()
}

func (s *Server) Setup(ctx context.Context, cfg *config.Config) error {
	s.Config = cfg

	if err := s.SetupDeps(); err != nil {
		return err
	}

	if err := s.SetupHandlers(); err != nil {
		return err
	}

	port := fmt.Sprintf(":%v", cfg.Port)
	server := http.Server{
		Addr:    port,
		Handler: s.Router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("error starting http: %v", err)
		}
	}()

	log.Printf("server started on port %v", port)

	<-ctx.Done()

	if err := server.Shutdown(context.Background()); err != nil {
		log.Fatalf("error shutting down server: %v", err)
	}

	return nil
}