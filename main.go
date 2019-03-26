package main

import (
	"context"
	"expvar"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sovikc/neb/authoring"
	"github.com/sovikc/neb/postgres"
	"github.com/sovikc/neb/rde"
	"github.com/sovikc/neb/server"
)

func main() {

	// Setup repository
	var (
		projects rde.ProjectRepository
	)

	// create connection pool
	pool, err := getConnPool()
	if err != nil {
		panic(err)
	}
	projects = postgres.NewProjectRepository(pool)

	var as authoring.Service
	as = authoring.NewService(projects)
	as = authoring.NewLoggingService(as)
	as = authoring.NewInstrumentingService(
		expvar.NewInt("addProject"),
		expvar.NewInt("addFeature"),
		expvar.NewInt("addWireframe"),
		expvar.NewInt("updateWireframeTitle"),
		as)

	srv := server.New(as)

	httpServer := &http.Server{Addr: httpAddr,
		Handler:      srv,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second}

	// registers a function to call on Shutdown
	httpServer.RegisterOnShutdown(func() {
		log.Println("Call shutdown hooks")
	})

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		signal.Notify(sigint, syscall.SIGINT)
		signal.Notify(sigint, syscall.SIGTERM)
		<-sigint

		// We received an interrupt signal, shut down.
		if err := httpServer.Shutdown(context.Background()); err != nil {
			// Error from closing listeners, or context timeout:
			log.Printf("HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()

	if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
		// Error starting or closing listener:
		log.Printf("HTTP server ListenAndServe: %v", err)
	}

	<-idleConnsClosed

}
