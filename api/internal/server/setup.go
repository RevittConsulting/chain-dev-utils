package server

import (
	"github.com/RevittConsulting/cdk-envs/internal/buckets"
	"github.com/RevittConsulting/cdk-envs/internal/buckets/db/mdbx"
	"github.com/RevittConsulting/cdk-envs/internal/chains"
	"github.com/RevittConsulting/cdk-envs/internal/health"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

type dependencies struct {
	chain   *chains.Service
	buckets *buckets.Service
}

func (s *Server) SetupDeps() error {
	var deps dependencies

	deps.chain = chains.NewService(&s.Config.Chain)

	mdbxdb := mdbx.New()
	//if err := mdbxdb.Open("chaindata/mdbx.dat"); err != nil {
	//	return err
	//}
	deps.buckets = buckets.NewService(&s.Config.Buckets, mdbxdb)

	s.Deps = &deps
	return nil
}

func (s *Server) SetupHandlers() error {

	s.Router = chi.NewRouter()

	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"*"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	s.Router.Use(cors.Handler)

	s.Router.Route("/api/v1", func(r chi.Router) {
		health.NewHandler(r, s.ShuttingDown)
		chains.NewHandler(r, s.Deps.chain)
		buckets.NewHandler(r, s.Deps.buckets)
	})
	return nil
}
