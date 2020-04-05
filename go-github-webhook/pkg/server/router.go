package server

import (
	"time"

	"github.com/IceflowRE/redeclipse-server-docker/pkg/server/utils"
	"github.com/IceflowRE/redeclipse-server-docker/pkg/structs"
	"github.com/IceflowRE/redeclipse-server-docker/pkg/updater"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func CreateRouter(appConfig *utils.Config, updaterConfig *updater.Config, storage *structs.HashStorage, workDir string) *chi.Mux {
	mRouter := chi.NewRouter()

	mRouter.Use(middleware.RequestID)
	mRouter.Use(middleware.RealIP)
	mRouter.Use(middleware.Logger)
	mRouter.Use(middleware.Recoverer)
	mRouter.Use(middleware.Timeout(20 * time.Second))

	mRouter.Route("/", func(router chi.Router) {
		router.Get("/ping", PingGet)
		router.With(githubMiddleware([]byte(*appConfig.WebhookSecret))).Post("/", EventHandler(updaterConfig, storage, workDir))
		router.With(apiKeyMiddleware(appConfig.ApiKeys)).Post("/hash", HashPost(storage))
	})
	return mRouter
}
