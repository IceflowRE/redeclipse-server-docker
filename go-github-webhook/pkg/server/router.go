package server

import (
	"time"

	"github.com/IceflowRE/redeclipse-server-docker/pkg/server/utils"
	"github.com/IceflowRE/redeclipse-server-docker/pkg/updater"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func CreateRouter(appConfig *utils.AppConfig, updaterConfig *updater.AppConfig, storage *updater.HashStorage, buildCtx *updater.BuildContext) *chi.Mux {
	mRouter := chi.NewRouter()

	mRouter.Use(middleware.RequestID)
	mRouter.Use(middleware.RealIP)
	mRouter.Use(middleware.Logger)
	mRouter.Use(middleware.Recoverer)
	mRouter.Use(middleware.Timeout(20 * time.Second))

	mRouter.Route("/", func(router chi.Router) {
		router.With(utils.HeaderMiddle).With(utils.SignatureMiddle([]byte(*appConfig.WebhookSecret))).Post("/", EventHandler(updaterConfig, storage, buildCtx))
		router.Get("/ping", PingGet)
	})
	return mRouter
}
