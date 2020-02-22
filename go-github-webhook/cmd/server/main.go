package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	serverPkg "github.com/IceflowRE/redeclipse-server-docker/pkg/server"
)

func WaitingForClose(server *http.Server) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("\nShutdown with timeout: 5s")

	if err := server.Shutdown(ctx); err != nil {
		log.Println(err)
	}
}

func main() {
	config, updaterConfig, hashStorage, buildCtx := serverPkg.EntryPoint()
	if config == nil {
		os.Exit(1)
	}
	server := &http.Server{
		Handler:      serverPkg.CreateRouter(config, updaterConfig, hashStorage, buildCtx),
		Addr:         ":" + strconv.Itoa(*config.Port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	go WaitingForClose(server)
	log.Println("Webhook listener (HTTP) will run on port: " + strconv.Itoa(*config.Port))
	log.Fatalln(server.ListenAndServe())
}
